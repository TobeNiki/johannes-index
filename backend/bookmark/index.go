package bookmark

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func (es *ElasticSearchBookmaek) CreateIndeice() error {
	if es.index == "" {
		return errors.New("elastic search indcies name is string empty")
	}
	// analyzerの詳しい設定内容はhttps://qiita.com/hatsu/items/dacbbba02d72947df435

	query := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   3,
			"number_of_replicas": 1,
			"analyzer": map[string]interface{}{
				"index_analyzer": map[string]interface{}{
					"type":                "custom",
					"char_filter":         []string{"html_strip"}, //html タグを排除
					"discard_punctuation": false,                  //句読点を出力しない
					"tokenizer":           "kuromoji_tokenizer",   //kuromoji default
					"filter": []string{
						"kuromoji_baseform",
						"kuromoji_part_of_speech",
						"ja_stop",
						"kuromoji_number",
						"kuromoji_stemmer",
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type":     "text",
					"analyzer": "index_analyzer",
				},
				"url": map[string]interface{}{
					"type": "text",
				},
				"text": map[string]interface{}{
					"type":     "text",
					"analyzer": "index_analyzer",
				},
				"folderId": map[string]interface{}{
					"type": "text",
				},
				"trashed": map[string]interface{}{
					"type": "boolean",
				},
				"date": map[string]interface{}{
					"type":   "date",
					"format": "yyyy-MM-dd",
				},
				"favicon": map[string]interface{}{
					"type": "text",
				},
			},
		},
	}
	data, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("error marshaling documnt: %s", err)
	}
	req := esapi.IndexRequest{
		Index:   es.index,
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}
	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error indexing document ID=%s", res.Status())
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return fmt.Errorf("error parsing the response body: %s", err)
		} else {
			fmt.Println(r)
			// Print the response status and indexed document version.
			if res.StatusCode == http.StatusCreated && r["result"] == "created" {
				return nil
			} else {
				return fmt.Errorf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
			}
		}
	}
}

func (es *ElasticSearchBookmaek) DeleteIndex() error {
	req := esapi.IndicesDeleteRequest{
		Index: []string{es.index},
	}
	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()
	fmt.Println(res)
	if res.IsError() {
		return errors.New("failed delete index ")
	}
	return nil
}
