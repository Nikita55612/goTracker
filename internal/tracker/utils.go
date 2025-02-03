package tracker

import (
	"fmt"
	"math/rand"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	casc "github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

var (
	selectors map[string]casc.Sel
	onceSel   sync.Once
)

type PeerId [20]byte

func HashesFromPieces(p []byte) PieceHashes {
	hashes := make([][20]byte, 0, len(p)/20)
	for c := range slices.Chunk(p, 20) {
		hash := [20]byte{}
		copy(hash[:], c)
		hashes = append(hashes, hash)
	}
	return hashes
}

// Поиск атрибута
func FindAttr(attrs []html.Attribute, k string) *string {
	for _, a := range attrs {
		if a.Key == k {
			return &a.Val
		}
	}
	return nil
}

// Генерация случайного PeerId
func GenPeerId() PeerId {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	peerId := PeerId{}
	prefix := []byte("-qB5001-")
	copy(peerId[:8], prefix)
	for i := 8; i < 20; i++ {
		peerId[i] = charset[seededRand.Intn(len(charset))]
	}
	return peerId
}

// Получение ссылки для загрузки пиров
func BuildTrackerURL(p *TrackerReqParam) (string, error) {
	base, err := url.Parse(p.BaseUrl)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(p.InfoHash[:])},
		"peer_id":    []string{string(p.PeerId[:])},
		"port":       []string{strconv.Itoa(int(p.Port))},
		"uploaded":   []string{strconv.Itoa(p.Uploaded)},
		"downloaded": []string{strconv.Itoa(p.Downloaded)},
		"compact":    []string{strconv.Itoa(p.Compact)},
		"left":       []string{strconv.Itoa(p.Left)},
		//"event":      []string{"started"},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func getSel(k string) casc.Sel {
	onceSel.Do(func() {
		srSel, _ := casc.Parse("#search-results table tbody")
		trSel, _ := casc.Parse("tr")
		tdSel, _ := casc.Parse("td")
		fnSel, _ := casc.Parse(".f-name a")
		ttSel, _ := casc.Parse(".t-title a")
		unSel, _ := casc.Parse(".u-name a")
		selectors = map[string]casc.Sel{
			"sr": srSel,
			"tr": trSel,
			"td": tdSel,
			"fn": fnSel,
			"tt": ttSel,
			"un": unSel,
		}
	})
	return selectors[k]
}

// Парсинг элементов на странице поска
func ParseSearchResult(n *html.Node) (res *SearchResult, err error) {
	n = casc.Query(n, getSel("sr"))
	if n == nil {
		err = fmt.Errorf("parse search result error")
		return
	}
	searchRes := make(SearchResult, 0, 50)
	items := casc.QueryAll(n, getSel("tr"))
	for _, item := range items {
		cols := casc.QueryAll(item, getSel("td"))
		if len(cols) < 9 {
			continue
		}
		item := new(SearchItem)
		if ptrId := FindAttr(cols[0].Attr, "id"); ptrId != nil {
			intId, err := strconv.Atoi(*ptrId)
			if err == nil {
				item.Id = &intId
			}
		}
		if ptrMark := FindAttr(cols[1].Attr, "title"); ptrMark != nil {
			item.Mark = ptrMark
		}
		if fn := casc.Query(cols[2], getSel("fn")); fn != nil {
			if fn = fn.FirstChild; fn != nil {
				item.Forum = &fn.Data
			}
		}
		if tt := casc.Query(cols[3], getSel("tt")); tt != nil {
			if tt = tt.FirstChild; tt != nil {
				item.Topic = &tt.Data
			}
		}
		if un := casc.Query(cols[4], getSel("un")); un != nil {
			if un = un.FirstChild; un != nil {
				item.Author = &un.Data
			}
		}
		if sz := cols[5].FirstChild; sz != nil {
			if sz = sz.NextSibling; sz != nil {
				if sz = sz.FirstChild; sz != nil {
					ptrSize := strings.Replace(sz.Data, " ↓", "", 1)
					item.Size = &ptrSize
				}
			}
		}
		if se := FindAttr(cols[6].Attr, "data-ts_text"); se != nil {
			if intSeeds, err := strconv.Atoi(*se); err == nil {
				item.Seeds = &intSeeds
			}
		}
		if ls := cols[7].FirstChild; ls != nil {
			if intLeechs, err := strconv.Atoi(ls.Data); err == nil {
				item.Leechs = &intLeechs
			}
		}
		if dl := cols[8].FirstChild; dl != nil {
			if intDl, err := strconv.Atoi(dl.Data); err == nil {
				item.Downloads = &intDl
			}
		}
		if dt := cols[9].FirstChild; dt != nil {
			if dt = dt.NextSibling; dt != nil {
				if dt = dt.FirstChild; dt != nil {
					item.Date = &dt.Data
				}
			}
		}
		searchRes = append(searchRes, *item)
	}
	return &searchRes, nil
}
