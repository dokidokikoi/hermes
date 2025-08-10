package ggbases

import (
	"fmt"
	"os"
	"testing"

	"resty.dev/v3"
)

func Test_GetCover(t *testing.T) {
	fmt.Println(GetCover("RJ01312275"))
	fmt.Println(GetCover("3351716_b0ac9b4834_"))
	fmt.Println(GetCover("https://cdn.cloudflare.steamstatic.com/steam/apps/3072130/header.jpg"))
	fmt.Println(GetCover("RJ308908"))
	fmt.Println(GetCover("VJ01002420"))
	fmt.Println(GetCover("VJ01002890"))
}

func Test_GenDLCover(t *testing.T) {
	fmt.Println(GenDLCover("01312275", "RJ"))
}

func Test_CoverUrl(t *testing.T) {
	fmt.Println(CoverUrl("3351716_b0ac9b4834_"))
}

func Test_MarkdownImg(t *testing.T) {
	markdown := `[![FastPic.Ru](https://i125.fastpic.org/thumb/2025/0507/31/86cbef33cff24b8b4ede6dc45a0c7531.jpeg)](https://fastpic.org/view/125/2025/0507/86cbef33cff24b8b4ede6dc45a0c7531.webp.html) [![FastPic.Ru](https://i125.fastpic.org/thumb/2025/0507/cf/b4997c79264c3ce45ebecb4f7d7816cf.jpeg)](https://fastpic.org/view/125/2025/0507/b4997c79264c3ce45ebecb4f7d7816cf.webp.html) [![FastPic.Ru](https://i125.fastpic.org/thumb/2025/0507/87/88edfa1deeec37f3971f650c8a4b0487.jpeg)](https://fastpic.org/view/125/2025/0507/88edfa1deeec37f3971f650c8a4b0487.webp.html) [![FastPic.Ru](https://i125.fastpic.org/thumb/2025/0507/fd/_fb7798221ee570c5341dfe172402b5fd.jpeg)](https://fastpic.org/view/125/2025/0507/_fb7798221ee570c5341dfe172402b5fd.png.html) [![FastPic.Ru](https://i125.fastpic.org/thumb/2025/0507/b6/_b545a7ec41d9d41790f29400db341db6.jpeg)](https://fastpic.org/view/125/2025/0507/_b545a7ec41d9d41790f29400db341db6.png.html)`

	imgs := MarkdownImg(markdown)
	for _, img := range imgs {
		fmt.Println(img)
	}
}

func Test_Get(t *testing.T) {
	rsp, err := resty.New().R().Get("https://ggbases.dlgal.com/search.so?p=0&title=%E5%BD%BC%E5%A5%B3&advanced=0")
	if err != nil {
		t.Fatal(err)
	}
	f, _ := os.Create("index.html")
	f.Write(rsp.Bytes())
	f.Close()
}
