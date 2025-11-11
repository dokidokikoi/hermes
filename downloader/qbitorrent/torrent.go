package qbitorrent

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	commmon_errors "github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/tools"
	"github.com/pkg/errors"
)

const (
	API_TORRENT_LIST         = apiVer + "/torrents/info"
	API_TORRENT_PAUSE        = apiVer + "/torrents/stop"
	API_TORRENT_RESUME       = apiVer + "/torrents/start"
	API_TORRENT_SET_LOCATION = apiVer + "/torrents/setLocation"
	API_TORRENT_DELETE       = apiVer + "/torrents/delete"
	API_TORRENT_ADD          = apiVer + "/torrents/add"
)

type GetTorrentsParam struct {
	Filter   string
	Category string
	Tag      string
	Sort     string
	Reverse  bool
	Limit    int64
	Offset   int64
	Hashes   []string
}

type Torrent struct {
	AddedOn                  int     `json:"added_on"`
	AmountLeft               int     `json:"amount_left"`
	AutoTmm                  bool    `json:"auto_tmm"`
	Availability             float64 `json:"availability"`
	Category                 string  `json:"category"`
	Comment                  string  `json:"comment"`
	Completed                int     `json:"completed"`
	CompletionOn             int     `json:"completion_on"`
	ContentPath              string  `json:"content_path"`
	DlLimit                  int     `json:"dl_limit"`
	Dlspeed                  int     `json:"dlspeed"`
	DownloadPath             string  `json:"download_path"`
	Downloaded               int     `json:"downloaded"`
	DownloadedSession        int     `json:"downloaded_session"`
	Eta                      int     `json:"eta"`
	FLPiecePrio              bool    `json:"f_l_piece_prio"`
	ForceStart               bool    `json:"force_start"`
	HasMetadata              bool    `json:"has_metadata"`
	Hash                     string  `json:"hash"`
	InactiveSeedingTimeLimit int     `json:"inactive_seeding_time_limit"`
	InfohashV1               string  `json:"infohash_v1"`
	InfohashV2               string  `json:"infohash_v2"`
	LastActivity             int     `json:"last_activity"`
	MagnetURI                string  `json:"magnet_uri"`
	MaxInactiveSeedingTime   int     `json:"max_inactive_seeding_time"`
	MaxRatio                 int     `json:"max_ratio"`
	MaxSeedingTime           int     `json:"max_seeding_time"`
	Name                     string  `json:"name"`
	NumComplete              int     `json:"num_complete"`
	NumIncomplete            int     `json:"num_incomplete"`
	NumLeechs                int     `json:"num_leechs"`
	NumSeeds                 int     `json:"num_seeds"`
	Popularity               float64 `json:"popularity"`
	Priority                 int     `json:"priority"`
	Private                  bool    `json:"private"`
	Progress                 float64 `json:"progress"`
	Ratio                    float64 `json:"ratio"`
	RatioLimit               int     `json:"ratio_limit"`
	Reannounce               int     `json:"reannounce"`
	RootPath                 string  `json:"root_path"`
	SavePath                 string  `json:"save_path"`
	SeedingTime              int     `json:"seeding_time"`
	SeedingTimeLimit         int     `json:"seeding_time_limit"`
	SeenComplete             int     `json:"seen_complete"`
	SeqDl                    bool    `json:"seq_dl"`
	Size                     int     `json:"size"`
	State                    string  `json:"state"`
	SuperSeeding             bool    `json:"super_seeding"`
	Tags                     string  `json:"tags"`
	TimeActive               int     `json:"time_active"`
	TotalSize                int     `json:"total_size"`
	Tracker                  string  `json:"tracker"`
	TrackersCount            int     `json:"trackers_count"`
	UpLimit                  int     `json:"up_limit"`
	Uploaded                 int64   `json:"uploaded"`
	UploadedSession          int64   `json:"uploaded_session"`
	Upspeed                  int     `json:"upspeed"`
}

func (c *Clinet) GetTorrents(param GetTorrentsParam) ([]Torrent, error) {
	ps := make(map[string]string)
	if param.Filter != "" {
		ps["filter"] = param.Filter
	}
	if param.Category != "" {
		ps["category"] = param.Category
	}
	if param.Tag != "" {
		ps["tag"] = param.Tag
	}
	if param.Sort != "" {
		ps["sort"] = param.Sort
	}
	if param.Reverse {
		ps["reverse"] = "true"
	}
	if param.Limit > 0 {
		ps["limit"] = strconv.FormatInt(param.Limit, 10)
	}
	if param.Offset > 0 {
		ps["offset"] = strconv.FormatInt(param.Offset, 10)
	}
	if len(param.Hashes) > 0 {
		ps["hashs"] = strings.Join(param.Hashes, "|")
	}
	rsp, err := tools.Req(
		http.MethodGet,
		addr+API_TORRENT_LIST,
		nil,
		tools.SetQueryParamsWithOption(ps),
		tools.SetCookiesWithOption(c.cookies...),
	)
	if err != nil {
		return nil, err
	}
	if !rsp.IsSuccess() {
		if rsp.StatusCode() == http.StatusForbidden {
			return nil, errors.Wrapf(commmon_errors.ErrAccessDenied, "maindata failed, status: %s", rsp.Status())
		}
		return nil, errors.Errorf("maindata failed, status: %s", rsp.Status())
	}

	result := []Torrent{}
	body := rsp.Bytes()
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.Wrapf(err, "json unmarshal error, body: %s", string(body))
	}
	return result, nil
}

func (c *Clinet) PauseTorrents(hashes ...string) error {
	if len(hashes) == 0 {
		return nil
	}
	rsp, err := tools.Req(
		http.MethodPost,
		addr+API_TORRENT_PAUSE,
		nil,
		tools.SetFromWithOption(map[string]string{
			"hashes": strings.Join(hashes, "|"),
		}),
		tools.SetCookiesWithOption(c.cookies...),
	)
	if err != nil {
		return err
	}
	if !rsp.IsSuccess() {
		if rsp.StatusCode() == http.StatusForbidden {
			return errors.Wrapf(commmon_errors.ErrAccessDenied, "pause torrent failed, status: %s", rsp.Status())
		}
		return errors.Errorf("pause torrent failed, status: %s", rsp.Status())
	}
	return nil
}

func (c *Clinet) ResumeTorrents(hashes ...string) error {
	if len(hashes) == 0 {
		return nil
	}
	rsp, err := tools.Req(
		http.MethodPost,
		addr+API_TORRENT_RESUME,
		nil,
		tools.SetFromWithOption(map[string]string{
			"hashes": strings.Join(hashes, "|"),
		}),
		tools.SetCookiesWithOption(c.cookies...),
	)
	if err != nil {
		return err
	}
	if !rsp.IsSuccess() {
		if rsp.StatusCode() == http.StatusForbidden {
			return errors.Wrapf(commmon_errors.ErrAccessDenied, "resume torrent failed, status: %s", rsp.Status())
		}
		return errors.Errorf("resume torrent failed, status: %s", rsp.Status())
	}
	return nil
}

func (c *Clinet) SetTorrentLocation(location string, hashes ...string) error {
	if len(hashes) == 0 || location == "" {
		return nil
	}
	rsp, err := tools.Req(
		http.MethodPost,
		addr+API_TORRENT_SET_LOCATION,
		nil,
		tools.SetFromWithOption(map[string]string{
			"hashes":   strings.Join(hashes, "|"),
			"location": location,
		}),
		tools.SetCookiesWithOption(c.cookies...),
	)
	if err != nil {
		return err
	}
	if !rsp.IsSuccess() {
		if rsp.StatusCode() == http.StatusForbidden {
			return errors.Wrapf(commmon_errors.ErrAccessDenied, "set torrent location failed, status: %s", rsp.Status())
		}
		return errors.Errorf("set torrent location failed, status: %s", rsp.Status())
	}
	return nil
}

func (c *Clinet) DelTorrents(deleteFiles bool, hashes ...string) error {
	if len(hashes) == 0 {
		return nil
	}
	rsp, err := tools.Req(
		http.MethodPost,
		addr+API_TORRENT_DELETE,
		nil,
		tools.SetFromWithOption(map[string]string{
			"hashes":      strings.Join(hashes, "|"),
			"deleteFiles": strconv.FormatBool(deleteFiles),
		}),
		tools.SetCookiesWithOption(c.cookies...),
	)
	if err != nil {
		return err
	}
	if !rsp.IsSuccess() {
		if rsp.StatusCode() == http.StatusForbidden {
			return errors.Wrapf(commmon_errors.ErrAccessDenied, "delete torrent failed, status: %s", rsp.Status())
		}
		return errors.Errorf("delete torrent failed, status: %s", rsp.Status())
	}
	return nil
}

type AddTorrentParam struct {
	Urls         []string // URLs separated with newlines
	Torrents     []string // Raw data of torrent file. torrents can be presented multiple times.
	Savepath     string   // Download folder
	Category     string   // Category for the torrent
	Tags         []string // Tags for the torrent, split by ','
	SkipChecking bool     // Skip hash checking. Possible values are true, false (default)
	Paused       bool     // Add torrents in the paused state. Possible values are true, false (default)
	Rename       string   // Rename torrent
}

func (c *Clinet) AddTorrents(param AddTorrentParam) error {
	ps := make(map[string]string)
	if len(param.Urls) > 0 {
		ps["urls"] = strings.Join(param.Urls, "\n")
	}
	if len(param.Savepath) > 0 {
		ps["savepath"] = param.Savepath
	} else if param.Category != "" {
		categories, err := c.ListTorrentCategory()
		if err != nil {
			return err
		}
		cate, ok := categories[param.Category]
		if ok && cate.SavePath != "" {
			ps["savepath"] = cate.SavePath
		}
	}
	if param.Category != "" {
		ps["category"] = param.Category
	}
	if len(param.Tags) > 0 {
		ps["tags"] = strings.Join(param.Tags, ",")
	}
	if param.SkipChecking {
		ps["skip_checking"] = "true"
	}
	if param.Paused {
		ps["paused"] = "true"
	}
	if param.Rename != "" {
		ps["rename"] = param.Rename
	}
	files := make(map[string][]string)
	if len(param.Torrents) > 0 {
		files["fileselect[]"] = param.Torrents
	}

	rsp, err := tools.Req(
		http.MethodPost,
		addr+API_TORRENT_ADD,
		nil,
		tools.SetMultiFileWithOption(ps, files),
		tools.SetCookiesWithOption(c.cookies...),
		tools.SetHeadersWithOption(map[string]string{
			"referer": "http://192.168.1.5:8999/upload.html",
		}),
	)
	if err != nil {
		return err
	}
	if !rsp.IsSuccess() {
		if rsp.StatusCode() == http.StatusForbidden {
			return errors.Wrapf(commmon_errors.ErrAccessDenied, "add torrent failed, status: %s", rsp.Status())
		}
		return errors.Errorf("add torrent failed, status: %s", rsp.Status())
	}
	return nil
}
