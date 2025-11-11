package utils

import (
	"net/url"
	"path/filepath"
	"strings"
)

func GetFileName(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		_, name := filepath.Split(uri)
		return name
	}

	arr := strings.Split(u.Path, "/")
	return strings.Split(arr[len(arr)-1], "?")[0]
}
