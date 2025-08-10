package ggbases

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	reg1 *regexp.Regexp
	reg2 *regexp.Regexp
	reg3 *regexp.Regexp
	reg4 *regexp.Regexp
	reg5 *regexp.Regexp
	reg6 *regexp.Regexp
)

var markdownRe *regexp.Regexp

func init() {
	reg1 = regexp.MustCompile("^d[0-9]{6,8}$")
	reg2 = regexp.MustCompile("^RJ[0-9]{6,8}$")
	reg3 = regexp.MustCompile("^v[0-9]{6,8}$")
	reg4 = regexp.MustCompile("^VJ[0-9]{6,8}$")
	reg5 = regexp.MustCompile("^g[0-9]{6,7}$")
	reg6 = regexp.MustCompile("^gc[0-9]{6,7}$")
	markdownRe = regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
}

func GetCover(cover string) (src string) {
	if len(cover) < 4 {
		return ""
	}
	switch {
	case reg1.MatchString(cover):
		src = GenDLCover(cover[1:], "RJ")
	case reg2.MatchString(cover):
		src = GenDLCover(cover[2:], "RJ")
	case reg3.MatchString(cover):
		src = GenDLCover(cover[1:], "VJ")
	case reg4.MatchString(cover):
		src = GenDLCover(cover[2:], "VJ")
	case reg5.MatchString(cover):
		src = GenGCCover(cover[1:])
	case reg6.MatchString(cover):
		src = GenGCCover(cover[2:])
	default:
		if cover[:4] == "http" {
			return cover
		} else if cover[:2] == "//" {
			return "https:" + cover
		} else {
			return CoverUrl(cover)
		}
	}
	return
}

func GenDLCover(didStr, ty string) string {
	didLen := len(didStr)
	if didLen == 0 {
		return ""
	}
	template := fmt.Sprintf("%%0%dd", didLen)

	did, err := strconv.ParseInt(didStr, 10, 64)
	if err != nil {
		return ""
	}
	rid := did/1000*1000 + 1000

	cover, _ := url.JoinPath(
		"https://img.dlsite.jp/modpub/images2/work/",
		func() string {
			if ty == "RJ" {
				return "doujin/RJ"
			}
			return "professional/VJ"
		}()+fmt.Sprintf(template, rid),
		fmt.Sprintf("%s%s_img_main.jpg", ty, fmt.Sprintf(template, did)),
	)
	return cover
}

func GenGCCover(did string) string {
	return fmt.Sprintf("https://cover.ydgal.com/_300_cover/getchu/gc%s.jpg", did)
}

func CoverUrl(cover string) string {
	arr := strings.Split(cover, "_")
	if len(arr) == 0 {
		return ""
	}
	num, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		return ""
	}
	ty := "old"
	if num > 1360000 {
		ty = "new"
	}

	src, _ := url.JoinPath("https://cover.ydgal.com/_200_cover/", ty, cover)
	return src
}

func MarkdownImg(text string) []string {
	matches := markdownRe.FindAllStringSubmatch(text, -1)

	imgs := []string{}
	for _, match := range matches {
		if len(match) > 1 {
			imgs = append(imgs, match[1])
		}
	}
	return imgs
}
