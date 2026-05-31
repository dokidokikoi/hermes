package main

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"

	"resty.dev/v3"
)

const (
	nhentaiApiPrefix = "https://nhentai.net/api/v2"
)

type NHentai struct {
	httpClient *resty.Client
}

func NewNHentai() *NHentai {
	return &NHentai{
		httpClient: resty.New(),
	}
}

// func (h *NHentai) GetItem(title string) (*Comic, error) {
// 	items, err := h.search(title, "date", 1)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(items) == 0 {
// 		return nil, errors.New("not found")
// 	}
// 	detail, err := h.getDetail(items[0].ID)
// 	if err != nil {
// 		return nil, err
// 	}

// }

type (
	nhentaiSearchItem struct {
		ID              int64   `json:"id"`
		MediaID         string  `json:"media_id"`
		EnglishTitle    string  `json:"english_title"`
		JapaneseTitle   string  `json:"japanese_title"`
		Thumbnail       string  `json:"thumbnail"`
		ThumbnailWidth  int64   `json:"thumbnail_width"`
		ThumbnailHeight int64   `json:"thumbnail_height"`
		NumPages        int     `json:"num_pages"`
		TagIds          []int64 `json:"tag_ids"`
		Blacklisted     bool    `json:"blacklisted"`
	}
	nhentaiDetail struct {
		Loc   any    `json:"loc"`
		Msg   string `json:"msg"`
		Type  string `json:"type"`
		Input string `json:"input"`
		Ctx   any    `json:"ctx"`
	}
	nhentaiSearchResp struct {
		Result  []nhentaiSearchItem `json:"result"`
		NumPage int                 `json:"num_pages"`
		PrePage int                 `json:"pre_page"`
		Total   int                 `json:"total"`
		Detail  nhentaiDetail       `json:"detail"`
		Error   string              `json:"error"`
	}
)

func (h *NHentai) search(query, sort string, page int) ([]nhentaiSearchItem, error) {
	resp, err := h.httpClient.R().SetQueryParams(map[string]string{
		"query": query,
		"sort":  sort,
		"page":  strconv.Itoa(page),
	}).Get(nhentaiApiPrefix + "/search")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}
	var result nhentaiSearchResp
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, err
	}
	if result.Error != "" {
		return nil, errors.New(result.Error)
	}
	return result.Result, nil
}

type (
	nhentaiDetailResp struct {
		ID      int64  `json:"id"`
		MediaID string `json:"media_id"`
		Title   struct {
			English  string `json:"english"`
			Japanese string `json:"japanese"`
			Pretty   string `json:"pretty"`
		} `json:"title"`
		Cover struct {
			Path string `json:"path"`
		} `json:"cover"`
		Thumbnail struct {
			Path string `json:"path"`
		} `json:"thumbnail"`
		Tags []struct {
			ID    int64  `json:"id"`
			Type  string `json:"type"`
			Name  string `json:"name"`
			Slug  string `json:"slug"`
			Url   string `json:"url"`
			Count int    `json:"count"`
		} `json:"tags"`
		NumPages     int    `json:"num_pages"`
		NumFavorites int    `json:"num_favorites"`
		Error        string `json:"error"`
	}
)

func (h *NHentai) getDetail(id int64) (*nhentaiDetailResp, error) {
	resp, err := h.httpClient.R().Get(nhentaiApiPrefix + "/galleries/" + strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("status code: %d, body: %s", resp.StatusCode(), resp.String())
	}
	var result nhentaiDetailResp
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return nil, err
	}
	if result.Error != "" {
		return nil, errors.New(result.Error)
	}
	return &result, nil
}
