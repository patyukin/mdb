package parser

import (
	"strings"
	"unicode"
)

func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

func isValidArgChar(r rune) bool {
	return isLetter(r) || isDigit(r) || isPunctuation(r)
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r)
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func containsWildcard(s string) bool {
	return strings.Contains(s, "*")
}

func isUppercase(ch rune) bool {
	return 'A' <= ch && ch <= 'Z'
}

// isPunctuation проверяет, является ли символ допустимым знаком пунктуации
func isPunctuation(ch rune) bool {
	switch ch {
	case '*', '/', '_', '-', '.', '+', '=', '?', '&', '%', '$', '#', '@', '!':
		return true
	default:
		return false
	}
}

// isValidArgumentChar проверяет, допустим ли символ в аргументе
func isValidArgumentChar(ch rune) bool {
	return isLetter(ch) || isDigit(ch) || isPunctuation(ch) || ch == '/'
}
