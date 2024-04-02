package utils

import (
	"strings"
	"unicode"
)

/**
 * GenerateFriendlyURL generates a friendly URL from a title
 */
func GenerateFriendlyURL(title string) string {
	var friendlyURL string
	for _, char := range title {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			friendlyURL += string(char)
		} else {
			friendlyURL += "-"
		}
	}
	return strings.ToLower(friendlyURL)
}
