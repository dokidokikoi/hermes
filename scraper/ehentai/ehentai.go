package main

import (
	"bytes"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"resty.dev/v3"
)

const (
	ehentaiApiUrl = "https://api.e-hentai.org/api.php"
	ehentaiUrl    = "https://e-hentai.org/"
)

type EHentai struct {
	httpClient *resty.Client
}

func NewEHentai() *EHentai {
	return &EHentai{
		httpClient: resty.New(),
	}
}

type ehentaiSearchItem struct {
	GID   string `json:"gid"`
	Token string `json:"token"`
	Title string `json:"title"`
}

func (e *EHentai) search(keyword string) ([]ehentaiSearchItem, error) {
	resp, err := e.httpClient.R().SetHeaders(map[string]string{
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36",
	}).SetQueryParams(map[string]string{
		"f_search": keyword,
	}).Get(ehentaiUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}
	body := resp.Bytes()
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	items := []ehentaiSearchItem{}
	doc.Find("table.gltc tr").Each(func(i int, s *goquery.Selection) {
		nameNode := s.Find("td.glname")
		gUrl := nameNode.Find("a").AttrOr("href", "")
		if gUrl != "" {
			u, err := url.Parse(gUrl)
			if err != nil {
				return
			}
			paths := strings.Split(strings.Trim(u.Path, "/"), "/")
			if len(paths) < 3 {
				return
			}
			items = append(items, ehentaiSearchItem{
				GID:   paths[len(paths)-2],
				Token: paths[len(paths)-1],
				Title: nameNode.Text(),
			})
		}
	})
	return items, nil
}

type ehentaiDetail struct {
	GID          int64    `json:"gid"`
	Token        string   `json:"token"`
	Title        string   `json:"title"`
	TitleJpn     string   `json:"title_jpn"`
	Category     string   `json:"category"`
	Uploader     string   `json:"uploader"`
	Posted       string   `json:"posted"`
	FileCount    string   `json:"filecount"`
	FileSize     int64    `json:"filesize"`
	Expunged     bool     `json:"expunged"`
	Rating       string   `json:"rating"`
	TorrentCount string   `json:"torrentcount"`
	Tags         []string `json:"tags"`
	ParentGID    string   `json:"parent_gid"`
	ParentKey    string   `json:"parent_key"`
	CurrentGid   string   `json:"current_gid"`
	CurrentKey   string   `json:"current_key"`
	FirstGID     string   `json:"first_gid"`
	FirstKey     string   `json:"first_key"`
	Error        string   `json:"error"`
}

func (e *EHentai) getDetail(gid, token string) (*ehentaiDetail, error) {
	resp, err := e.httpClient.R().SetBody(map[string]any{
		"method": "gdata",
		"gidlist": [][]string{
			{gid, token},
		},
		"namespace": 1,
	}).Post(ehentaiApiUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}
	var result struct {
		Gmetadata []ehentaiDetail `json:"gmetadata"`
	}
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, err
	}
	if len(result.Gmetadata) == 0 {
		return nil, errors.New("no metadata")
	}
	metadata := result.Gmetadata[0]
	if metadata.Error != "" {
		return nil, errors.New(metadata.Error)
	}
	return &metadata, nil
}
