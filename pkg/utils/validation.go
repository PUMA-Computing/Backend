package utils

import "strings"

func RemoveWhitespace(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func SplitEmail(s string) (string, string) {
	parts := strings.Split(s, "@")
	return parts[0], parts[1]
}

func IsEmail(s string) bool {
	return strings.Contains(s, "@")
}
