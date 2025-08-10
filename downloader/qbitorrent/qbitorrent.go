package qbitorrent

import (
	"hermes/tools"
	"net/http"

	"github.com/pkg/errors"
)

var addr = "http://192.168.1.5:8999"

const apiVer = "/api/v2"
const (
	API_AUTH = apiVer + "/auth/login"

	API_MAIN_DATA = apiVer + "/sync/maindata"
)

type Clinet struct {
	cookies []*http.Cookie
}

func (c *Clinet) Auth(username, password string) ([]*http.Cookie, error) {
	rsp, err := tools.Req(
		http.MethodPost,
		addr+API_AUTH,
		nil,
		tools.SetHeadersWithOption(map[string]string{
			"Referer": addr,
		}),
		tools.SetFromWithOption(map[string]string{
			"username": username,
			"password": password,
		}))
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode() != http.StatusOK {
		return nil, errors.Errorf("auth failed, status: %s", rsp.Status())
	}
	return rsp.Cookies(), nil
}
