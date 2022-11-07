package scraping

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/TobeNiki/Index/backend/bookmark"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/saintfish/chardet"
	"github.com/sclevine/agouti"
	"golang.org/x/net/html/charset"
)

// chorme(selenium)を使用してクロールするのは最終手段
type ChormeDriver struct {
	driver   *agouti.WebDriver
	headless bool
}

func (chrome *ChormeDriver) Start() {
	//chromedriver set up
	var optionsArgs []string
	optionsArgs = append(optionsArgs, "--disable-gpu")
	if chrome.headless {
		optionsArgs = append(optionsArgs, "--headless")
	}
	option := agouti.ChromeOptions(
		"args", optionsArgs,
	)
	chrome.driver = agouti.ChromeDriver(option)
	defer chrome.driver.Stop()
}
func (chorme *ChormeDriver) chormeDriverRequest(url string) (*strings.Reader, error) {
	// execute Start function before executing this function
	if err := chorme.driver.Start(); err != nil {
		return nil, fmt.Errorf("chorme driver error is %s", err)
	}
	page, err := chorme.driver.NewPage()
	if err != nil {
		return nil, fmt.Errorf("chorme driver new page error %s", err)
	}
	if err = page.Navigate(url); err != nil {
		return nil, fmt.Errorf("failed to navigate: %v", err)
	}
	dom, err := page.HTML()
	if err != nil {
		return nil, fmt.Errorf("failed to get html : %v", err)
	}
	return strings.NewReader(dom), nil
}

func (chrome *ChormeDriver) Stop() {
	chrome.driver.Stop()
}

//フロントからpostで送られてくるjson形式の定義
type BookmarkRequest struct {
	ID                string `json:"id"`
	URL               string `json:"url"`
	IsUseChromeDriver bool   `json:"isUseChromeDriver"`
	FolderId          string `json:"folderId"`
}
type Crower struct {
	ChormeDriver      ChormeDriver
	Sleeptime         int
	IsUseChromeDriver bool
}

func (crower *Crower) ScrapingData(url string) (*bookmark.BookmarkSource, error) {
	var document *goquery.Document
	if crower.IsUseChromeDriver {
		reader, err := crower.ChormeDriver.chormeDriverRequest(url)
		if err != nil {
			return nil, err
		}
		document, err = goquery.NewDocumentFromReader(reader)
		if err != nil {
			return nil, err
		}
	} else {
		reader, err := requestGet(url)
		if err != nil {
			return nil, err
		}
		document, err = goquery.NewDocumentFromReader(reader)
		if err != nil {
			return nil, err
		}
		// ogp, err := GetOGP(reader)
		// opg 保存処理
	}
	bookmarkData := bookmark.BookmarkSource{}
	bookmarkData.URL = url
	bookmarkData.Title = document.Find("title").Text()
	bookmarkData.FolderID = ""
	bodySelection := document.Find("body")
	pTagTexts := bodySelection.Find("p").Text()
	bookmarkData.Text = strip.StripTags(pTagTexts) //主にpタグのテキストを取得する
	bookmarkData.Trashed = false
	bookmarkData.Date = time.Now().Format("2006-01-02")
	faviconStr, err := GetFavicon(url)
	if err != nil {
		//faviconが取得できなかった場合はデフォルトのfavicon(google)を格納しておく
		faviconStr = DefaultBase64Code
	}
	bookmarkData.Favicon = faviconStr
	bookmarkData.OmitStopWord()

	return &bookmarkData, nil
}

func New() *Crower {
	return &Crower{
		ChormeDriver: ChormeDriver{
			headless: true,
		},
		Sleeptime:         5,
		IsUseChromeDriver: false,
	}
}
func requestGet(url string) (io.Reader, error) {
	client := &http.Client{Timeout: time.Duration(10) * time.Second}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed http client %v", err)
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	response, err := client.Do(request)
	//404の場合は登録しないようにする(登録しても無駄なので)
	//403及び400はスクレイピング禁止サイトの可能性がある
	//存在しないわけではないはずなので、テキスト情報を空にして登録する
	if err != nil || response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed get response %v", err)
	}

	defer response.Body.Close()
	buffer, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response body")
	}
	reader, err := charDetectAndTransform(buffer)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func charDetectAndTransform(buffer []byte) (io.Reader, error) {
	detector := chardet.NewTextDetector()
	detectResult, err := detector.DetectBest(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed charctor detect: %v", err)
	}
	bufReader := bytes.NewReader(buffer)
	reader, err := charset.NewReaderLabel(detectResult.Charset, bufReader)
	if err != nil {
		return nil, fmt.Errorf("failed charcotr transform: %v", err)
	}
	return reader, nil
}
