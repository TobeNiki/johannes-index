package scraping

import (
	"testing"
)

func TestRequestGet(t *testing.T) {
	url := "https://qiita.com/KamikawaTakato/items/73b748414567d27a9c52"
	_, err := requestGet(url)
	if err != nil {
		t.Errorf("正常系エラー %v", err)
	}
	url = "https://www.youtube.com/watch?v=00H3ZzWV0kw"
	_, err = requestGet(url)
	if err != nil {
		t.Error("sns 失敗")
	}
	url = "https://morioh.com/p/05c0b1ec326b"
	_, err = requestGet(url)
	if err != nil {
		t.Error("海外サイト")
	}
	// 異常系エラー404
	url = "https://qiita.com/KamikawaTakato/items/73b748414567d27a9c52/dsds"
	_, err = requestGet(url)
	if err == nil {
		t.Error("404 errorが失敗")
	}
	//スクレイピング禁止サイト
	url = "https://kakuyomu.jp/works/16816452219389978278/episodes/16816700427262585662"
	_, err = requestGet(url)
	if err != nil {
		t.Error("403 errorが失敗")
	}
}
func TestScrapingData(t *testing.T) {
	crower := New()

	//正常系
	url := "https://qiita.com/KamikawaTakato/items/73b748414567d27a9c52"
	result, err := crower.ScrapingData(url)
	if err != nil {
		t.Errorf("正常系エラー %v", err)
	}
	if result.Title != "無料、独学で機械学習エンジニアになる！~機械学習が学べる無料サイト、書籍~ - Qiita" {
		t.Error("failed title get")
	}
	if result.Text == "" {
		t.Error("failed text get")
	}
	//異常系404
	url += "https://qiita.com/KamikawaTakato/items/73b748414567d27a9c52/dsds"
	_, err = crower.ScrapingData(url)
	if err == nil {
		t.Error("404 error が失敗")
	}
	url = "https://www.youtube.com/watch?v=00H3ZzWV0kw"
	_, err = crower.ScrapingData(url)
	if err != nil {
		t.Error("sns が取得失敗")
	}

	url = "https://morioh.com/p/05c0b1ec326b"
	_, err = crower.ScrapingData(url)
	if err != nil {
		t.Error("sns が取得失敗")
	}
	//スクレイピング禁止サイト
	//今はUserAgentをつけると回避できてしまう=>後々改良する
	//今は暫定的にスクレイピングする
	url = "https://kakuyomu.jp/works/16816452219389978278/episodes/16816700427262585662"
	_, err = crower.ScrapingData(url)

	if err != nil {
		t.Error("403 errorが失敗")
	}

}
