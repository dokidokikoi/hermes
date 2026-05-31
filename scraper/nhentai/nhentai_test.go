package main

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("https_proxy", "http://127.0.0.1:7890")
}

func Test_NHentaiSearch(t *testing.T) {
	client := NewNHentai()

	items, err := client.search("ぺどりあ! プリンセスフランドール", "date", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(items)
}

func Test_NHentaiGetDetail(t *testing.T) {
	client := NewNHentai()

	items, err := client.getDetail(58645)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(items)
}
