package tracker

import (
	"fmt"
	"gotracker/config"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"
)

var (
	reqClient   *http.Client
	baseUrl     string
	onceClient  sync.Once
	onceBaseUrl sync.Once
)

// Клиент запросов к трекеру
func ReqClient() *http.Client {
	onceClient.Do(func() {
		cfg := config.Get()
		proxyUrl, _ := url.Parse(cfg.Tracker.ProxyUrl)
		jar, _ := cookiejar.New(nil)
		reqClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
			Jar: jar,
		}

		cookie := &http.Cookie{
			Name:  "bb_session",
			Value: cfg.Tracker.BBSession,
		}
		cookies := [1]*http.Cookie{cookie}
		baseUrl, _ := url.Parse(cfg.Tracker.BaseUrl)
		reqClient.Jar.SetCookies(baseUrl, cookies[:])
	})
	return reqClient
}

// Получение домена трекера
func BaseUrl() string {
	onceBaseUrl.Do(func() {
		cfg := config.Get()
		baseUrl = cfg.Tracker.BaseUrl
	})
	return baseUrl
}

// Поиск торрентов по запросу с указанием страницы поиска
func SearchQuery(q string, page int) (*SearchResult, error) {
	params := url.Values{
		"start": []string{strconv.Itoa(page * 50)},
		"nm":    []string{q},
		"o":     []string{"4"},
		"s":     []string{"2"},
	}.Encode()
	reqURL := BaseUrl() + "/forum/tracker.php" + "?" + params
	req, _ := http.NewRequest("POST", reqURL, nil)
	resp, err := ReqClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := charmap.Windows1251.NewDecoder()
	node, err := html.Parse(decoder.Reader(resp.Body))
	if err != nil {
		return nil, err
	}
	return ParseSearchResult(node)
}

// Загрузка .torrent файла по id элемента
func ItemTorrentFile(id int) (file []byte, err error) {
	reqURL := BaseUrl() + "/forum/dl.php?t=" + strconv.Itoa(id)
	req, _ := http.NewRequest("POST", reqURL, nil)
	resp, err := ReqClient().Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	file, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

// Запрос к трекеру для получения пиров
func TrackerReq(p *TrackerReqParam) (trackerResp *BencodeTrackerResp, err error) {
	reqUrl, err := BuildTrackerURL(p)
	if err != nil {
		return
	}
	client := http.Client{}
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return
	}
	req.Header.Add("User-Agent", "uTorrent/2210(25110)")
	fmt.Println("...resp")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	fmt.Println("resp...")
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	trackerResp, err = TrackerRespFromBytes(data)
	if err != nil {
		return
	}
	return
}
