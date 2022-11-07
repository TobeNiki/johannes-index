package scraping

import (
	"fmt"
	"io"

	"github.com/dyatlov/go-opengraph/opengraph"
)

/*
og := opengraph.NewOpenGraph()
err = og.ProcessHTML(strings.NewReader(string(body)))
で得られる情報は以下(例)
{
    "type":"article",
    "url":"https://qiita.com/7110/items/4ece0ce9be0ce910ee90",
    "title":"Python datetime 日付の計算、文字列変換をする方法 strftime, strptime【決定版】 - Qiita",
    "description":"\n\nPython datetime 日付の計算、文字列変換をする方法 strftime, strptime\n\n\n\n2019/08/02 最終更新\n\n\nここではPython標準のdatetimeによる、日付の取得・計算 timedelt...",
    "determiner":"",
    "site_name":"Qiita",
    "locale":"",
    "locales_alternate":null,
    "images":[
        {
            "url":"https://qiita-user-contents.imgix.net/https%3A%2F%2Fcdn.qiita.com%2Fassets%2Fpublic%2Farticle-ogp-background-9f5428127621718a910c8b63951390ad.png?ixlib=rb-4.0.0\u0026w=1200\u0026mark64=aHR0cHM6Ly9xaWl0YS11c2VyLWNvbnRlbnRzLmltZ2l4Lm5ldC9-dGV4dD9peGxpYj1yYi00LjAuMCZ3PTkxNiZ0eHQ9UHl0aG9uJTIwZGF0ZXRpbWUlMjAlRTYlOTclQTUlRTQlQkIlOTglRTMlODElQUUlRTglQTglODglRTclQUUlOTclRTMlODAlODElRTYlOTYlODclRTUlQUQlOTclRTUlODglOTclRTUlQTQlODklRTYlOEYlOUIlRTMlODIlOTIlRTMlODElOTklRTMlODIlOEIlRTYlOTYlQjklRTYlQjMlOTUlMjBzdHJmdGltZSUyQyUyMHN0cnB0aW1lJUUzJTgwJTkwJUU2JUIxJUJBJUU1JUFFJTlBJUU3JTg5JTg4JUUzJTgwJTkxJnR4dC1jb2xvcj0lMjMyMTIxMjEmdHh0LWZvbnQ9SGlyYWdpbm8lMjBTYW5zJTIwVzYmdHh0LXNpemU9NTYmdHh0LWNsaXA9ZWxsaXBzaXMmdHh0LWFsaWduPWxlZnQlMkN0b3Amcz0xZGRhYmFiZDNjN2U5YzM3NzgwZDAwZTQ2MjMwNDY4Zg\u0026mark-x=142\u0026mark-y=112\u0026blend64=aHR0cHM6Ly9xaWl0YS11c2VyLWNvbnRlbnRzLmltZ2l4Lm5ldC9-dGV4dD9peGxpYj1yYi00LjAuMCZ3PTYxNiZ0eHQ9JTQwNzExMCZ0eHQtY29sb3I9JTIzMjEyMTIxJnR4dC1mb250PUhpcmFnaW5vJTIwU2FucyUyMFc2JnR4dC1zaXplPTM2JnR4dC1hbGlnbj1sZWZ0JTJDdG9wJnM9NjJhYmJlYmFjM2IyNWY1OTUzNjlmNWFkMzU1YmE1MjQ\u0026blend-x=142\u0026blend-y=491\u0026blend-mode=normal\u0026s=b92420d541ff103e7f5a11ace414b719",
            "secure_url":"",
            "type":"",
            "width":0,
            "height":0
        }
    ],
    "audios":null,
    "videos":null,
    "article":{
        "published_time":null,
        "modified_time":null,
        "expiration_time":null,
        "section":"",
        "tags":null,
        "authors":null
    }
}
*/
//以下の構造対は後にDBに保存するようにする？
type Ogp struct {
	Type        string      `json:"type"`
	Url         string      `json:"url"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Determiner  string      `json:"determiner"`
	SiteName    string      `json:"site_name"`
	Locale      string      `json:"locale"`
	Image       []ogp_Image `json:"images"`
	//以下audioなど
}
type ogp_Image struct {
	Url       string `json:"url"`
	SecureUrl string `json:"secure_url"`
	Type      string `json:"type"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

func GetOGP(reader io.Reader) (*opengraph.OpenGraph, error) {
	og := opengraph.NewOpenGraph()
	if err := og.ProcessHTML(reader); err != nil {
		return nil, fmt.Errorf("ogp process error %v", err)
	}
	return og, nil
}
