package tools

import "regexp"

func DateExtarct(str string) string {
	re := regexp.MustCompile(`[1-2]\d\d\d-[0-1]?\d-[0-3]?\d`)
	return re.FindString(str)
}
