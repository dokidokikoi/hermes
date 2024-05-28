package getchu

import (
	"fmt"
	"hermes/config"
	"hermes/tools"
	"net/http"
	"os"
	"testing"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
)

func init() {
	config.SetConfig("../../conf/application.yaml")
	zaplog.SetLogger(config.GetConfig().LogConfig)
}

func TestReq(t *testing.T) {
	data, err := GetChuScraper.DoReq(http.MethodGet, "https://www.getchu.com/php/nsearch.phtml?genre=pc_soft&search_keyword=ボクの彼女はガテン系&check_key_dtl=1&submit=", nil, nil)
	if err != nil {
		panic(err)
	}
	f, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(f, string(tools.Jp2Utf8(data)))
}
