package stringutil

import (
	"strings"
	"unicode"
)

// Reverse returns the string reversed (rune-safe).
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Capitalize uppercases the first letter of each word.
func Capitalize(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		runes := []rune(w)
		runes[0] = unicode.ToUpper(runes[0])
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// WordCount returns the number of words in a string.
func WordCount(s string) int {
	return len(strings.Fields(s))
}
