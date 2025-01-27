package utils

import "strings"

func EscapeMarkdown(text string) string {
	var builder strings.Builder
	// List of characters to escape according to Telegram MarkdownV2
	specialChars := map[rune]bool{
		'_': true, '*': true, '[': true, ']': true, '(': true,
		')': true, '~': true, '`': true, '>': true, '#': true,
		'+': true, '-': true, '=': true, '|': true, '{': true,
		'}': true, '.': true, '!': true, '\\': true,
	}

	for _, r := range text {
		if specialChars[r] {
			builder.WriteRune('\\')
		}
		builder.WriteRune(r)
	}

	return builder.String()
}
