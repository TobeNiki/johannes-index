package scraping

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const faviconAPI = "http://www.google.com/s2/favicons"
const DefaultBase64Code = "VBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAACXBIWXMAAAsSAAALEgHS3X78AAACiElEQVQ4EaVTzU8TURCf2tJuS7tQtlRb6UKBIkQwkRRSEzkQgyEc6lkOKgcOph78Y+CgjXjDs2i44FXY9AMTlQRUELZapVlouy3d7kKtb0Zr0MSLTvL2zb75eL838xtTvV6H/xELBptMJojeXLCXyobnyog4YhzXYvmCFi6qVSfaeRdXdrfaU1areV5KykmX06rcvzumjY/1ggkR3Jh+bNf1mr8v1D5bLuvR3qDgFbvbBJYIrE1mCIoCrKxsHuzK+Rzvsi29+6DEbTZz9unijEYI8ObBgXOzlcrx9OAlXyDYKUCzwwrDQx1wVDGg089Dt+gR3mxmhcUnaWeoxwMbm/vzDFzmDEKMMNhquRqduT1KwXiGt0vre6iSeAUHNDE0d26NBtAXY9BACQyjFusKuL2Ry+IPb/Y9ZglwuVscdHaknUChqLF/O4jn3V5dP4mhgRJgwSYm+gV0Oi3XrvYB30yvhGa7BS70eGFHPoTJyQHhMK+F0ZesRVVznvXw5Ixv7/C10moEo6OZXbWvlFAF9FVZDOqEABUMRIkMd8GnLwVWg9/RkJF9sA4oDfYQAuzzjqzwvnaRUFxn/X2ZlmGLXAE7AL52B4xHgqAUqrC1nSNuoJkQtLkdqReszz/9aRvq90NOKdOS1nch8TpL555WDp49f3uAMXhACRjD5j4ykuCtf5PP7Fm1b0DIsl/VHGezzP1KwOiZQobFF9YyjSRYQETRENSlVzI8iK9mWlzckpSSCQHVALmN9Az1euDho9Xo8vKGd2rqooA8yBcrwHgCqYR0kMkWci08t/R+W4ljDCanWTg9TJGwGNaNk3vYZ7VUdeKsYJGFNkfSzjXNrSX20s4/h6kB81/271ghG17l+rPTAAAAAElFTkSuQmCC"

func GetFavicon(urlStr string) (string, error) {
	// urlを解析してドメイン名を取得する
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("failed url parse %v", err)
	}
	//ドメイン名が存在しない場合はデフォルトのBASE64コードを返却
	if u.Hostname() == "" {
		return "", errors.New("failed get errors")
	}
	parts := strings.Split(u.Hostname(), ".")
	var domain string
	if len(parts) > 1 {
		domain = parts[len(parts)-2] + "." + parts[len(parts)-1]
	} else {
		domain = u.Hostname()
	}
	// faviconAPI(非公式？)にアクセス
	request, err := http.NewRequest("GET", faviconAPI, nil)
	if err != nil {
		return "", fmt.Errorf("error  request object %v", err)
	}
	params := request.URL.Query()
	params.Add("domain", domain)
	request.URL.RawQuery = params.Encode()
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed request favicon api %v", err)
	}

	defer response.Body.Close()
	// byteデータとして取得し base64文字列にする
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed read response body %v", err)
	}
	base64Data := base64.StdEncoding.EncodeToString(body)
	return base64Data, nil
}
