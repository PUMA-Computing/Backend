package utils

import "strings"

func RemoveWhitespace(s string) string {
	return strings.ReplaceAll(s, " ", "")
}
