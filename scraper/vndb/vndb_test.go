package vndb_test

import (
	"izumi/scraper/vndb"
	"os"
	"testing"
)

var vndbScraper *vndb.VNDB

func init() {
	os.Setenv("https_proxy", "socks://127.0.0.1:7890")
	os.Setenv("http_proxy", "socks://127.0.0.1:7890")
	os.Setenv("all_proxy", "socks://127.0.0.1:7890")
	vndbScraper = vndb.NewVNDB(nil, "").(*vndb.VNDB)
}

func Test_GetItem(t *testing.T) {
	item, err := vndbScraper.GetItem("v2002")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(item)
}
