package qbitorrent

import (
	"izumi/tools"
	"net/http"
	"strings"

	commmon_errors "github.com/dokidokikoi/go-common/errors"
	"github.com/pkg/errors"
)

const (
	API_TORRENT_TAG_ADD = apiVer + "/torrents/addTags"
	API_TORRENT_TAG_REM = apiVer + "/torrents/removeTags"
)

func (c *Clinet) AddTorrentTags(tags []string, hashes ...string) error {
	if len(hashes) == 0 || len(tags) == 0 {
		return nil
	}
	rsp, err := tools.Req(
		http.MethodPost,
		addr+API_TORRENT_TAG_ADD,
		nil,
		tools.SetFromWithOption(map[string]string{
			"hashes": strings.Join(hashes, "|"),
			"tags":   strings.Join(tags, ","),
		}),
		tools.SetCookiesWithOption(c.cookies...),
	)
	if err != nil {
		return err
	}
	if !rsp.IsSuccess() {
		if rsp.StatusCode() == http.StatusForbidden {
			return errors.Wrapf(commmon_errors.ErrAccessDenied, "add torrent tags failed, status: %s", rsp.Status())
		}
		return errors.Errorf("add torrent tags failed, status: %s", rsp.Status())
	}
	return nil
}

func (c *Clinet) RemTorrentTags(tags []string, hashes ...string) error {
	if len(hashes) == 0 || len(tags) == 0 {
		return nil
	}
	rsp, err := tools.Req(
		http.MethodPost,
		addr+API_TORRENT_TAG_REM,
		nil,
		tools.SetFromWithOption(map[string]string{
			"hashes": strings.Join(hashes, "|"),
			"tags":   strings.Join(tags, ","),
		}),
		tools.SetCookiesWithOption(c.cookies...),
	)
	if err != nil {
		return err
	}
	if !rsp.IsSuccess() {
		if rsp.StatusCode() == http.StatusForbidden {
			return errors.Wrapf(commmon_errors.ErrAccessDenied, "remove torrent tags failed, status: %s", rsp.Status())
		}
		return errors.Errorf("remove torrent tags failed, status: %s", rsp.Status())
	}
	return nil
}
