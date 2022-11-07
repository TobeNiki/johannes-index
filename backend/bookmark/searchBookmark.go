package bookmark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TobeNiki/Index/backend/utils"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func topValidate(topNumber int) int {
	if topNumber > 10000 {
		topNumber = 10000
	}
	if topNumber < 0 {
		topNumber = 10
	}
	return topNumber
}
func sortTypeValidate(sortType string) string {
	sortTypeValues := []string{"asc, desc"}
	if !utils.IncludeStringSlice(sortTypeValues, sortType) {
		sortType = "desc" // default
	}
	return sortType
}
func (es *ElasticSearchBookmaek) searchRequest(query map[string]interface{}) ([]Bookmark, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %s", err)
	}
	res, err := es.client.Search(
		es.client.Search.WithContext(context.Background()),
		es.client.Search.WithIndex(es.index),
		es.client.Search.WithBody(&buf),
		es.client.Search.WithTrackTotalHits(true),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response : %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %s", err)
		} else {
			return nil, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	var b map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&b); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}
	responseData := []Bookmark{}
	for _, hit := range b["hits"].(map[string]interface{})["hits"].([]interface{}) {
		result := new(Bookmark)
		result.SetData(hit)
		responseData = append(responseData, *result)
	}
	return responseData, nil
}
func (es *ElasticSearchBookmaek) GetBookmark(topNumber int, sortType string) ([]Bookmark, error) {
	topNumber = topValidate(topNumber)
	sortType = sortTypeValidate(sortType)
	query := map[string]interface{}{
		"size": topNumber,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"trashed": false,
						},
					},
					{
						"range": map[string]interface{}{
							"date": map[string]interface{}{
								"lte": time.Now().Format("2006-01-02"),
							},
						},
					},
				},
			},
		},
		"sort": map[string]interface{}{
			"date": map[string]interface{}{
				"order": sortType,
			},
		},
	}

	return es.searchRequest(query)
}
func (es *ElasticSearchBookmaek) GetBookmarkFromID(id string) (Bookmark, error) {
	req := esapi.IndexRequest{
		Index:      es.index,
		DocumentID: id,
	}
	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return Bookmark{}, fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return Bookmark{}, fmt.Errorf("[%s] error indexing document ID=%s", res.Status(), id)
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return Bookmark{}, fmt.Errorf("error parsing the response body: %s", err)
	}
	responseData := Bookmark{}
	responseData.SetData(r)
	return responseData, nil
}
func (es *ElasticSearchBookmaek) GetBookmarkFromFolderID(folderId string, topNumber int, sortType string) ([]Bookmark, error) {
	topNumber = topValidate(topNumber)
	sortType = sortTypeValidate(sortType)
	query := map[string]interface{}{
		"size": topNumber,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"folderId": folderId,
						},
					},
					{
						"range": map[string]interface{}{
							"date": map[string]interface{}{
								"lte": time.Now().Format("2006-01-02"),
							},
						},
					},
				},
			},
		},
		"sort": map[string]interface{}{
			"date": map[string]interface{}{
				"order": sortType,
			},
		},
	}
	return es.searchRequest(query)
}
func (es *ElasticSearchBookmaek) GetBookmarkFromTrash(topNumber int, sortType string) ([]Bookmark, error) {
	topNumber = topValidate(topNumber)
	sortType = sortTypeValidate(sortType)
	query := map[string]interface{}{
		"size": topNumber,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"trashed": true,
						},
					},
					{
						"range": map[string]interface{}{
							"date": map[string]interface{}{
								"lte": time.Now().Format("2006-01-02"),
							},
						},
					},
				},
			},
		},
		"sort": map[string]interface{}{
			"date": map[string]interface{}{
				"order": sortType,
			},
		},
	}
	return es.searchRequest(query)
}

type SearchBookmarkQuery struct {
	SearchWord          string `json:"word"`
	SearchWordTypeIsAnd bool   `json:"search_word_type_is_and"`
	SearchTarget        string `json:"search_target"`
	FolderID            string `json:"folderid"`
	Top                 int    `json:"top"`
	SortType            string `json:"sort_type"`
}

func (es *ElasticSearchBookmaek) GetBookmarkFromUnorganized(topNumber int, sortType string) ([]Bookmark, error) {
	result, err := es.GetBookmark(topNumber, sortType)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(result); i++ {
		if result[i].Source.FolderID != "" {
			//folderidが空ではない物は取り出す
			result = result[:i+copy(result[i:], result[i+1:])]
		}
	}
	return result, nil
}
func (es *ElasticSearchBookmaek) SearchBookmark(searchData *SearchBookmarkQuery) ([]Bookmark, error) {
	searchData.Top = topValidate(searchData.Top)

	var mustOfBoolQueryStr = []map[string]interface{}{}
	//ゴミ箱に入っているとこからは検索しない
	mustOfBoolQueryStr = append(mustOfBoolQueryStr, map[string]interface{}{
		"match": map[string]interface{}{
			"trashed": false,
		},
	})
	searchTragetPattern := []string{"title", "text"}
	var textMatchQueryStr map[string]interface{}
	if searchData.SearchWord != "" {
		if !utils.IncludeStringSlice(searchTragetPattern, searchData.SearchTarget) {
			searchData.SearchTarget = "title"
		}
		targetField := utils.GetHitStringFromSlice(searchTragetPattern, searchData.SearchTarget)
		if searchData.SearchWordTypeIsAnd {
			textMatchQueryStr = map[string]interface{}{
				"match": map[string]interface{}{
					targetField: searchData.SearchWord,
					"operator":  "and",
				},
			}
		} else {
			textMatchQueryStr = map[string]interface{}{
				"match": map[string]interface{}{
					targetField: searchData.SearchWord,
				},
			}
		}
		mustOfBoolQueryStr = append(mustOfBoolQueryStr, textMatchQueryStr)
	}

	var folderMatchQueryStr map[string]interface{}
	if searchData.FolderID != "" {
		folderMatchQueryStr = map[string]interface{}{
			"match": map[string]interface{}{
				"folderId": searchData.FolderID,
			},
		}
		mustOfBoolQueryStr = append(mustOfBoolQueryStr, folderMatchQueryStr)
	}
	sortType := sortTypeValidate(searchData.SortType)
	var query map[string]interface{}
	if searchData.SearchWord != "" || searchData.FolderID != "" {
		query = map[string]interface{}{
			"size": searchData.Top,
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": mustOfBoolQueryStr,
				},
			},
			"sort": map[string]interface{}{
				"date": map[string]interface{}{
					"order": sortType,
				},
			},
		}
	}
	return es.searchRequest(query)
}
