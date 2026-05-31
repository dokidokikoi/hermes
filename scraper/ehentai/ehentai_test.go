package main

import "testing"

func Test_EHentaiSearch(t *testing.T) {
	client := NewEHentai()

	items, err := client.search("ふたり、ひと夏のあやまち -呂500-")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(items)
}

func Test_EHentaiGetDetail(t *testing.T) {
	client := NewEHentai()

	item, err := client.getDetail("1124192", "d53e3a31be")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(item)
}
