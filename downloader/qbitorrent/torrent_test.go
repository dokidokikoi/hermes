package qbitorrent

import (
	"fmt"
	"testing"
)

func Test_GetTorrents(t *testing.T) {
	cli := Clinet{}
	cookies, err := cli.Auth("admin", "password")
	if err != nil {
		t.Fatal(err)
	}
	cli.cookies = cookies

	data, err := cli.GetTorrents(GetTorrentsParam{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data)
}

func Test_PauseTorrent(t *testing.T) {
	cli := Clinet{}
	cookies, err := cli.Auth("admin", "password")
	if err != nil {
		t.Fatal(err)
	}
	cli.cookies = cookies

	err = cli.PauseTorrents("6e07791f695a218fae3088e6ca1d85ccc2a1c494")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ResumeTorrent(t *testing.T) {
	cli := Clinet{}
	cookies, err := cli.Auth("admin", "password")
	if err != nil {
		t.Fatal(err)
	}
	cli.cookies = cookies

	err = cli.ResumeTorrents("6e07791f695a218fae3088e6ca1d85ccc2a1c494")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_AddTorrents(t *testing.T) {
	cli := Clinet{}
	cookies, err := cli.Auth("admin", "password")
	if err != nil {
		t.Fatal(err)
	}
	cli.cookies = cookies

	err = cli.AddTorrents(AddTorrentParam{
		Tags: []string{"mt"},
		// Torrents: []string{
		// 	"/home/doki/Downloads/[M-TEAM]GVG204.mkv.torrent",
		// 	"/home/doki/Downloads/[M-TEAM]GVH-751.torrent",
		// },
		Urls: []string{
			"magnet:?xt=urn:btih:2d1c123c3509fe9cee8609d5fe0d29d449059cb2&dn=%28%E6%88%90%E5%B9%B4%E3%82%B3%E3%83%9F%E3%83%83%E3%82%AF%29%20%5B%E4%BA%8C%E4%B8%89%E6%9C%88%E3%81%9D%E3%81%86%5D%201LDK%2BJK%20%E3%81%84%E3%81%8D%E3%81%AA%E3%82%8A%E5%90%8C%E5%B1%85%EF%BC%9F%20%E5%AF%86%E7%9D%80%EF%BC%81%EF%BC%9F%20%E5%88%9D%E3%82%A8%E3%83%83%E3%83%81%21%EF%BC%81%EF%BC%9F%20%E7%AC%AC1-53%E8%A9%B1&tr=http%3A%2F%2Fsukebei.tracker.wf%3A8888%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce",
			"magnet:?xt=urn:btih:e72f82522cc417e1e096026b96934fb83a8a3be2&dn=%28C97%29%20%5BMONE%E3%81%91%E3%81%97%E3%81%94%E3%82%80%20%28%E3%82%82%E3%81%AD%E3%81%A6%E3%81%83%29%5D%20TOHO%20Illustrations%202019%20winter%20%28%E6%9D%B1%E6%96%B9Project%29&tr=http%3A%2F%2Fsukebei.tracker.wf%3A8888%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ListTorrentCategory(t *testing.T) {
	cli := Clinet{}
	cookies, err := cli.Auth("admin", "password")
	if err != nil {
		t.Fatal(err)
	}
	cli.cookies = cookies

	categories, err := cli.ListTorrentCategory()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(categories)
}
