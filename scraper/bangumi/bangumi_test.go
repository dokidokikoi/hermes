package bangumi_test

import (
	"encoding/json"
	"fmt"
	"hermes/config"
	"hermes/scraper/bangumi"
	"testing"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
)

func init() {
	config.SetConfig("../../conf/application.yaml")
	zaplog.SetLogger(config.GetConfig().LogConfig)
}

func TestSearch(t *testing.T) {
	items, err := bangumi.BangumiScraper.Search("ボクの彼女はガテン系", 1)
	if err != nil {
		panic(err)
	}
	for _, item := range items {
		fmt.Printf("%+v\n", item.Name)
	}
}

func TestGetItem(t *testing.T) {
	item, err := bangumi.BangumiScraper.GetItem("https://api.bgm.tv/v0/subjects/259061")
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
