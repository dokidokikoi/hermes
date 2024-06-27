package tools

import "strings"

func GetFileName(url string) string {
	if url == "" {
		return ""
	}
	arr := strings.Split(url, "/")
	return arr[len(arr)-1]
}
