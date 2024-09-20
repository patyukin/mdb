package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

func IsWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

func IsValidArgChar(r rune) bool {
	return isLetter(r) || isDigit(r) || isPunctuation(r)
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r)
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

// ConvertPatternToRegex конвертирует шаблон с * в регулярное выражение.
func ConvertPatternToRegex(pattern string) string {
	pattern = regexp.QuoteMeta(pattern)
	pattern = strings.ReplaceAll(pattern, "\\*", ".*")

	return "^" + pattern + "$"
}

// MatchPattern сравнивает строку с шаблоном.
func MatchPattern(s, pattern string) (bool, error) {
	regexPattern := ConvertPatternToRegex(pattern)
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return false, fmt.Errorf("failed to compile regex pattern: %w", err)
	}

	return regex.MatchString(s), nil
}

func ContainsWildcard(s string) bool {
	return strings.Contains(s, "*")
}

func IsUppercase(ch rune) bool {
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

// IsValidArgumentChar проверяет, допустим ли символ в аргументе
func IsValidArgumentChar(ch rune) bool {
	return isLetter(ch) || isDigit(ch) || isPunctuation(ch) || ch == '/'
}
