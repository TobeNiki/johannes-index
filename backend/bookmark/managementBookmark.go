package bookmark

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/TobeNiki/Index/backend/database"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"golang.org/x/exp/utf8string"

	jwt "github.com/appleboy/gin-jwt/v2"
)

//返却用(json)構造体
type Bookmark struct {
	Id     string         `json:"_id"`
	Score  float64        `json:"_score"`
	Source BookmarkSource `json:"_source"`
}

func (bm *Bookmark) SetData(data interface{}) {
	doc := data.(map[string]interface{})
	bm.Id = doc["_id"].(string)

	bm.Source.Title = doc["_source"].(map[string]interface{})["title"].(string)
	bm.Source.URL = doc["_source"].(map[string]interface{})["url"].(string)
	bm.Source.Text = doc["_source"].(map[string]interface{})["text"].(string)
	if utf8.RuneCountInString(bm.Source.Text) > 100 {
		text := utf8string.NewString(bm.Source.Text)
		bm.Source.Text = text.Slice(0, 101)
		bm.Source.Text += "......"
	}
	bm.Source.FolderID = doc["_source"].(map[string]interface{})["folderId"].(string)
	bm.Source.Trashed = doc["_source"].(map[string]interface{})["trashed"].(bool)
	bm.Source.Date = doc["_source"].(map[string]interface{})["date"].(string)
	bm.Source.Favicon = doc["_source"].(map[string]interface{})["favicon"].(string)
}

type BookmarkSource struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Text     string `json:"text"`
	FolderID string `json:"folderId"`
	Trashed  bool   `json:"trashed"`
	Date     string `json:"date"`
	Favicon  string `json:"favicon"`
}

func (bs *BookmarkSource) OmitStopWord() {
	stopWord := []string{
		"\n", "\t", "	", "\u003e",
	}
	result := bs.Text
	for _, target := range stopWord {
		result = strings.Replace(result, target, "", -1)
	}
	bs.Text = result
}

type ElasticSearchBookmaek struct {
	client *elasticsearch.Client
	index  string
}

func New() *ElasticSearchBookmaek {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://es:9200", //docker es コンテナのアドレス
		},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err.Error())
	}
	return &ElasticSearchBookmaek{
		client: client,
	}
}

func (es *ElasticSearchBookmaek) SetIndexNameFromMUser(user *database.M_User) error {
	if user.ESIndexName != "" {
		es.index = user.ESIndexName
		return nil
	} else {
		return errors.New("elastic search indcies name is string empty")
	}
}
func (es *ElasticSearchBookmaek) SetIndexName(claims jwt.MapClaims) error {
	if _, ok := claims["index"]; ok {
		if claims["index"] != "" {
			es.index = claims["index"].(string)
			return nil
		}
	}
	return errors.New("index key or value is none")
}

func (es *ElasticSearchBookmaek) AddBookmarkGoRoutine(id string, bookmark *BookmarkSource) error {
	var (
		wg sync.WaitGroup
	)
	ch := make(chan error)
	wg.Add(1)
	go func(id string, bookmark BookmarkSource) {
		defer wg.Done()
		ch <- es.AddBookmark(id, &bookmark)
	}(id, *bookmark)
	wg.Wait()
	return <-ch
}
func (es *ElasticSearchBookmaek) AddBookmark(bookmarkid string, bookmark *BookmarkSource) error {

	data, err := json.Marshal(bookmark)
	if err != nil {
		return fmt.Errorf("error marshaling documnt: %s", err)
	}
	req := esapi.IndexRequest{
		Index:      es.index,
		DocumentID: bookmarkid,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
		Timeout:    5 * time.Second,
	}
	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("[%s] error indexing document ID=%s", res.Status(), bookmarkid)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return fmt.Errorf("error parsing the response body: %s", err)
		}
	}
	return nil

}
func (es *ElasticSearchBookmaek) DeleteBookmark(bookmarkId string) error {
	if bookmarkId == "" {
		return errors.New("bookmark id is string empty")
	}
	req := esapi.DeleteRequest{
		Index:      es.index,
		DocumentID: bookmarkId,
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error indexing document ID=%s", res.Status())
	}
	// Deserialize the response into a map.
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return fmt.Errorf("error parsing the response body: %s", err)
	}
	fmt.Println(r)
	if r["result"] != "deleted" {
		return errors.New("result is not deleted")
	}
	return nil
}

func (es *ElasticSearchBookmaek) TrashBookmarkWithFolder(folderId string) error {
	if folderId == "" {
		return errors.New("folder id is string empty")
	}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"folderId": folderId,
			},
		},
		"script": map[string]interface{}{
			"source": "ctx._source.trashed = params.isTrash",
			"lang":   "painless",
			"params": map[string]interface{}{"isTrash": true},
		},
	}
	return es.updateByQueryRequest(query)

}

func (es *ElasticSearchBookmaek) TrashBookmark(bookmarkId string) error {
	if bookmarkId == "" {
		return errors.New("bookmark id is string empty")
	}
	query := map[string]interface{}{
		"doc": map[string]interface{}{
			"trashed":  true,
			"folderId": "",
		},
	}
	data, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("error marshaling documnt: %s", err)
	}
	req := esapi.UpdateRequest{
		Index:      es.index,
		DocumentID: bookmarkId,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error indexing document ID=%s", res.Status())
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return fmt.Errorf("error parsing the response body: %s", err)
	}
	if r["result"] != "updated" {
		return errors.New("result is not deleted")
	}
	return nil
}

func (es *ElasticSearchBookmaek) UpdateBookmark(bookmarkId string, bookmark *BookmarkSource) error {
	if bookmarkId == "" {
		return errors.New("bookmark id is string empty")
	}
	query := map[string]interface{}{
		"doc": map[string]interface{}{
			"folderId": bookmark.FolderID,
			"url":      bookmark.URL,
			"text":     bookmark.Text,
			"favicon":  bookmark.Favicon,
			"trashed":  false, // ゴミ箱に入っているブクマの場合は更新するともとに戻る
		},
	}
	data, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("error marshaling documnt: %s", err)
	}
	req := esapi.UpdateRequest{
		Index:      es.index,
		DocumentID: bookmarkId,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() || res.StatusCode != http.StatusOK {
		return fmt.Errorf("error indexing document ID=%s", res.Status())
	}
	return nil
}
