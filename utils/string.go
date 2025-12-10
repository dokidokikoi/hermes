package utils

import "strings"

func EqualFoldNoSpace(a, b string) bool {
	a = strings.ReplaceAll(a, " ", "")
	b = strings.ReplaceAll(b, " ", "")
	return strings.EqualFold(a, b)
}

func ToLowerNoSpace(s string) string {
	replacer := strings.NewReplacer(" ", "", "\u3000", "")
	return strings.ToLower((replacer.Replace(s)))
}
