package bangumi_test

import (
	"encoding/json"
	"fmt"
	"izumi/scraper"
	"izumi/scraper/bangumi"
	"testing"
)

var bangumiScraper scraper.IGameScraper

func init() {
	// zaplog.SetLogger(config.GetConfig().LogConfig)
	bangumiScraper = bangumi.NewBangumi(map[string]string{
		"Authorization": "Bearer a3mi4euoJlUYRjoedVNd55Dqj4tHVpIbgTMEoiHp",
		"User-Agent":    bangumi.DefaultHeader_UserAgent,
	}, "")
}

func TestSearch(t *testing.T) {
	items, err := bangumiScraper.Search("ボクの彼女はガテン系", 1)
	if err != nil {
		panic(err)
	}
	for _, item := range items {
		fmt.Printf("%+v\n", item.Name)
	}
}

func TestGetItem(t *testing.T) {
	item, err := bangumiScraper.GetItem("https://api.bgm.tv/v0/subjects/259061")
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
