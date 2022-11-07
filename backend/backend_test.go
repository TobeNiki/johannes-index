package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TobeNiki/Index/database"
)

func TestMain(t *testing.T) {
	//バックエンドのシナリオテスト
	router := SetUpServer()
	w := httptest.NewRecorder()
	//ユーザー登録
	scenarioUser := database.Register{
		LoginData: database.Login{
			UserID:   "scenario",
			Password: "test",
		},
		DisplayName: "シナリオテスト用",
	}
	data, err := json.Marshal(scenarioUser)
	if err != nil {
		log.Fatal("json failed")
	}

	req, _ := http.NewRequest("POST", "/regist", bytes.NewReader(data))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status code not 201, message:%s", w.Body.String())
	}
	//ログイン処理
	loginUser := database.Login{
		UserID:   "scenario",
		Password: "test",
	}
	data, err = json.Marshal(loginUser)
	if err != nil {
		log.Fatal("json failed")
	}
	req, _ = http.NewRequest("GET", "/login", bytes.NewBuffer(data))
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf(" [login]status code not 201, message:%s", w.Body.String())
	}
	var jsonMapData map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&jsonMapData); err != nil {
		t.Errorf(" [login]status code not 201, message:%s", w.Body.String())
	}
	token := jsonMapData["token"].(string)
	fmt.Println(token)
}
