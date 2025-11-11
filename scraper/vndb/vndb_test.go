package vndb

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/dokidokikoi/go-common/tools"
)

func init() {
	os.Setenv("https_proxy", "socks://127.0.0.1:7890")
	os.Setenv("http_proxy", "socks://127.0.0.1:7890")
	os.Setenv("all_proxy", "socks://127.0.0.1:7890")
}

func Test_SearchUri(t *testing.T) {
	rsp, err := tools.ReqWithProxy(http.MethodPost, VNDBSearchUri, map[string]any{
		"fields": strings.Join(DetailFields, ","),
	}, "")
	if err != nil {
		t.Fatal(err)
	}
	if rsp.StatusCode()/100 != 2 {
		t.Error(rsp.StatusCode())
	}
	bs := rsp.Bytes()
	var resp BaseResponse[VN]
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}
