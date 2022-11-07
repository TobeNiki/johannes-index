package bookmark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

/*
グルーピング方法：
ユーザ単語：
	・Words
	・含む含まない
	・複数可で、複数条件のANDかORの選択も可能

	create table M_Group {
		sebsetWord varchar() //含む単語
		complementWord varchar() //含まない単語

	}

	上記は今後追加する予定
	フォルダー分けでルールベースでできるようにする、
	そのために以下必要
	・条件の保存先
	・M＿Folderに条件の有無カラムと、条件との紐づけ
*/
type FolderConditions struct {
	Conditions      []Condition `json:"conditions"`
	SubsetWordIsAnd bool        `json:"subsetWordIsAnd"`
}
type Condition struct {
	SubsetWord     string `json:"subsetWord"`
	ComplementWord string `json:"complementWord"`
}

//特定の単語に引っかかったものを一気にフォルダーに入れる処理
func (es *ElasticSearchBookmaek) GroupingBookmark(folderConditionsVal FolderConditions, folderId string) error {
	var boolQueryType string
	if folderConditionsVal.SubsetWordIsAnd {
		boolQueryType = "must"
	} else {
		boolQueryType = "should"
	}
	var matchQueries []map[string]interface{}
	var mustNotQueries []map[string]interface{}
	for _, cond := range folderConditionsVal.Conditions {
		if cond.SubsetWord != "" {
			matchQueries = append(matchQueries, map[string]interface{}{
				"match": map[string]interface{}{
					"text": cond.SubsetWord,
				},
			})
		}
		if cond.ComplementWord != "" {
			mustNotQueries = append(mustNotQueries, map[string]interface{}{
				"match": map[string]interface{}{
					"text": cond.ComplementWord,
				},
			})
		}
	}
	var boolQueries map[string]interface{}
	if len(matchQueries) > 0 && len(mustNotQueries) > 0 {
		boolQueries = map[string]interface{}{
			boolQueryType: matchQueries,
			"must_not":    mustNotQueries,
		}
	} else if len(matchQueries) > 0 {
		boolQueries = map[string]interface{}{
			boolQueryType: matchQueries,
		}
	} else if len(mustNotQueries) > 0 {
		boolQueries = map[string]interface{}{
			"must_not": mustNotQueries,
		}
	} else {
		return fmt.Errorf("target word is noting")
	}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": boolQueries,
		},
		"script": map[string]interface{}{
			"inline": "ctx._source.folderId = " + folderId,
		},
	}
	return es.updateByQueryRequest(query)
}

func (es *ElasticSearchBookmaek) updateByQueryRequest(query map[string]interface{}) error {
	data, err := json.Marshal(query)

	if err != nil {
		return fmt.Errorf("error marshaling documnt: %s", err)
	}
	req := esapi.UpdateByQueryRequest{
		Index: []string{es.index},
		Body:  bytes.NewReader(data),
	}
	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error indexing document ID=%s", res.Status())
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error update document")
	}
	return nil
}
