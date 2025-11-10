package tools

import "strings"

func EqualFoldNoSpace(a, b string) bool {
	a = strings.ReplaceAll(a, " ", "")
	b = strings.ReplaceAll(b, " ", "")
	return strings.EqualFold(a, b)
}

func ToLowerNoSpace(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", ""))
}
