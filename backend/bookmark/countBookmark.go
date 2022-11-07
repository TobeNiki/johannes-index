package bookmark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func (es *ElasticSearchBookmaek) CountBookmarkFromUnorganized() (float64, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"trashed": false,
						},
					},
				},
				"must_not": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"folderId": "",
						},
					},
				},
			},
		},
	}
	return es.countRequest(query)
}
func (es *ElasticSearchBookmaek) CountBookmarkFromFolder(folderId string) (float64, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"trashed": false,
						},
					},
					{
						"match": map[string]interface{}{
							"folderId": folderId,
						},
					},
				},
			},
		},
	}
	return es.countRequest(query)
}
func (es *ElasticSearchBookmaek) CountBookmarkFromTrashed() (float64, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"trashed": true,
			},
		},
	}
	return es.countRequest(query)
}
func (es *ElasticSearchBookmaek) CountAllBookmark() (float64, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"trashed": false,
			},
		},
	}
	return es.countRequest(query)
}
func (es *ElasticSearchBookmaek) countRequest(query map[string]interface{}) (float64, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return 0.0, fmt.Errorf("error encoding query: %s", err)
	}
	req := esapi.CountRequest{
		Index: []string{es.index},
		Body:  &buf,
	}
	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return 0.0, fmt.Errorf("error count response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return 0.0, fmt.Errorf("error count response: %s", err)
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		fmt.Println(err)
		return 0.0, fmt.Errorf("error parsing the response body: %s", err)
	}
	return r["count"].(float64), nil

}
