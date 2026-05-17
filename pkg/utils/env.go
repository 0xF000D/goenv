package utils

import "strings"

func Unquote(s string) string {
	if len(s) < 2 {
		return s
	}

	quotes := "'\"`"
	firstChar := s[0]
	lastChar := s[len(s)-1]

	if firstChar == lastChar && strings.ContainsAny(string(firstChar), quotes) {
		return s[1 : len(s)-1]
	}

	return s
}
