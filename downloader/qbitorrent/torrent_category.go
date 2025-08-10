package qbitorrent

import (
	"encoding/json"
	"hermes/tools"
	"net/http"
	"strings"

	commmon_errors "github.com/dokidokikoi/go-common/errors"
	"github.com/pkg/errors"
)

const (
	API_TORRENT_CATEGORY_SET  = apiVer + "/torrents/setCategory"
	API_TORRENT_CATEGORY_LIST = apiVer + "/torrents/categories"
)

type Category struct {
	Name     string `json:"name"`
	SavePath string `json:"savePath"`
}

func (c *Clinet) SetTorrentCategory(category string, hashes ...string) error {
	if len(hashes) == 0 || category == "" {
		return nil
	}
	rsp, err := tools.Req(
		http.MethodPost,
		addr+API_TORRENT_CATEGORY_SET,
		nil,
		tools.SetFromWithOption(map[string]string{
			"hashes":   strings.Join(hashes, "|"),
			"category": category,
		}),
		tools.SetCookiesWithOption(c.cookies...),
	)
	if err != nil {
		return err
	}
	if !rsp.IsSuccess() {
		if rsp.StatusCode() == http.StatusForbidden {
			return errors.Wrapf(commmon_errors.ErrAccessDenied, "set torrent category failed, status: %s", rsp.Status())
		}
		return errors.Errorf("set torrent category failed, status: %s", rsp.Status())
	}
	return nil
}

func (c *Clinet) ListTorrentCategory() (map[string]Category, error) {
	rsp, err := tools.Req(
		http.MethodPost,
		addr+API_TORRENT_CATEGORY_LIST,
		nil,
		tools.SetCookiesWithOption(c.cookies...),
	)
	if err != nil {
		return nil, err
	}
	if !rsp.IsSuccess() {
		if rsp.StatusCode() == http.StatusForbidden {
			return nil, errors.Wrapf(commmon_errors.ErrAccessDenied, "list torrent category failed, status: %s", rsp.Status())
		}
		return nil, errors.Errorf("list torrent category failed, status: %s", rsp.Status())
	}

	result := make(map[string]Category)
	body := rsp.Bytes()
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.Wrapf(err, "json unmarshal error, body: %s", string(body))
	}
	return result, nil
}
