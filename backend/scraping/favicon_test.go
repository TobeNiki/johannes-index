package scraping

import (
	"fmt"
	"testing"
)

//base64の文字を画像にデコードするサイト=>https://rakko.tools/tools/71/
func TestGetFavicon(t *testing.T) {
	//正常系
	var urlStr = "https://qiita.com/KamikawaTakato/items/73b748414567d27a9c52"
	result, err := GetFavicon(urlStr)
	if err != nil {
		t.Errorf("favicon get failed %v", err)
	}
	fmt.Println(result)

	//正常系２(sns)
	urlStr = "https://www.youtube.com/results?search_query=%E5%BD%B1%E3%81%AE%E5%AE%9F%E5%8A%9B%E8%80%85%E3%81%AB%E3%81%AA%E3%82%8A%E3%81%9F%E3%81%8F%E3%81%A6+op"
	result, err = GetFavicon(urlStr)
	if err != nil {
		t.Errorf("favicon get failed from sns url %v", err)
	}
	fmt.Println(result)

	//正常系３(グーグル検索)
	urlStr = "https://www.google.com/search?q=base64+%E7%94%BB%E5%83%8F%E5%8C%96&rlz=1C1FQRR_jaJP976JP976&oq=&aqs=chrome.1.35i39i362l8.224698j0j7&sourceid=chrome&ie=UTF-8"
	result, err = GetFavicon(urlStr)
	if err != nil {
		t.Errorf("favicon get failed from google %v", err)
	}
	fmt.Println(result)

	//異常系(urlじゃない)=>デフォルトの文字列が返されること
	//異常系だけど、error値を返さない
	urlStr = "localhostperfume"
	_, err = GetFavicon(urlStr)
	if err == nil {
		t.Error("異常系 ホストネームがないのにエラーが帰ってきていない")
	}
	//異常系(localhost)=>APIの使用的にlocalhostはグーグルのやつ
	urlStr = "http://localhost:9000/#/login"
	result, err = GetFavicon(urlStr)
	if err != nil {
		t.Errorf("favicon get failed from localhost %v", err)
	}
	fmt.Println(result)
}
