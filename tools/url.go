package tools

import "strings"

func GetFileName(url string) string {
	if url == "" {
		return ""
	}
	arr := strings.Split(url, "/")
	return strings.Split(arr[len(arr)-1], "?")[0]
}
