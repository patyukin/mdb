package utils

import (
	"fmt"
	"testing"
)

// Тест функции IsWhitespace
func TestIsWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{"Пробел", ' ', true},
		{"Табуляция", '\t', true},
		{"Новая строка", '\n', true},
		{"Возврат каретки", '\r', true},
		{"Неразрывный пробел", '\u00A0', true},
		{"Буква", 'A', false},
		{"Цифра", '1', false},
		{"Пунктуация", ',', false},
		{"Символ", '$', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsWhitespace(tt.input)
			if result != tt.expected {
				t.Errorf("IsWhitespace(%q) = %v; ожидалось %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Тест функции IsValidArgChar
func TestIsValidArgChar(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{"Буква заглавная", 'A', true},
		{"Буква строчная", 'z', true},
		{"Цифра", '5', true},
		{"Пунктуация", '.', true},
		{"Подчеркивание", '_', true},
		{"Дефис", '-', true},
		{"Звездочка", '*', true},
		{"Пробел", ' ', false},
		{"Символ", '$', false},
		{"Непечатаемый", '\n', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidArgChar(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidArgChar(%q) = %v; ожидалось %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Тест функции isLetter
func TestIsLetter(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{"Буква заглавная", 'A', true},
		{"Буква строчная", 'z', true},
		{"Буква кириллица", 'Б', true},
		{"Цифра", '0', false},
		{"Пунктуация", '.', false},
		{"Символ", '$', false},
		{"Пробел", ' ', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isLetter(tt.input)
			if result != tt.expected {
				t.Errorf("isLetter(%q) = %v; ожидалось %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Тест функции isDigit
func TestIsDigit(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{"Цифра ноль", '0', true},
		{"Цифра девять", '9', true},
		{"Буква", 'A', false},
		{"Пунктуация", '.', false},
		{"Символ", '$', false},
		{"Пробел", ' ', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDigit(tt.input)
			if result != tt.expected {
				t.Errorf("isDigit(%q) = %v; ожидалось %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Тест функции isPunctuation
func TestIsPunctuation(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{"Точка", '.', true},
		{"Запятая", ',', true},
		{"Восклицательный знак", '!', true},
		{"Вопросительный знак", '?', true},
		{"Дефис", '-', true},
		{"Кавычка", '"', true},
		{"Скобка", '(', true},
		{"Буква", 'A', false},
		{"Цифра", '1', false},
		{"Пробел", ' ', false},
		{"Символ", '$', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPunctuation(tt.input)
			if result != tt.expected {
				t.Errorf("isPunctuation(%q) = %v; ожидалось %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestConvertPatternToRegex tests the ConvertPatternToRegex function.
func TestConvertPatternToRegex(t *testing.T) {
	tests := []struct {
		pattern       string
		expectedRegex string
	}{
		{
			pattern:       "abc",
			expectedRegex: "^abc$",
		},
		{
			pattern:       "a*c",
			expectedRegex: "^a.*c$",
		},
		{
			pattern:       "*abc*",
			expectedRegex: "^.*abc.*$",
		},
		{
			pattern:       "a?c",
			expectedRegex: "^a\\?c$",
		},
		{
			pattern:       "a.c",
			expectedRegex: "^a\\.c$",
		},
		{
			pattern:       "a*c*d",
			expectedRegex: "^a.*c.*d$",
		},
		{
			pattern:       "",
			expectedRegex: "^$",
		},
		{
			pattern:       "*",
			expectedRegex: "^.*$",
		},
		{
			pattern:       "a[bc]d",
			expectedRegex: "^a\\[bc\\]d$",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Pattern: %q", tt.pattern), func(t *testing.T) {
			regex := ConvertPatternToRegex(tt.pattern)
			if regex != tt.expectedRegex {
				t.Errorf("ConvertPatternToRegex(%q) = %q; want %q", tt.pattern, regex, tt.expectedRegex)
			}
		})
	}
}

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		str     string
		pattern string
		match   bool
	}{
		{
			str:     "abc",
			pattern: "abc",
			match:   true,
		},
		{
			str:     "abc",
			pattern: "a*c",
			match:   true,
		},
		{
			str:     "abbc",
			pattern: "a*c",
			match:   true,
		},
		{
			str:     "ac",
			pattern: "a*c",
			match:   true,
		},
		{
			str:     "abcd",
			pattern: "a*d",
			match:   true,
		},
		{
			str:     "abcd",
			pattern: "a*c",
			match:   false,
		},
		{
			str:     "xyz",
			pattern: "*y*",
			match:   true,
		},
		{
			str:     "xyz",
			pattern: "x*z",
			match:   true,
		},
		{
			str:     "xyz",
			pattern: "x*y*z",
			match:   true,
		},
		{
			str:     "xyz",
			pattern: "a*",
			match:   false,
		},
		{
			str:     "abc",
			pattern: "",
			match:   false,
		},
		{
			str:     "",
			pattern: "",
			match:   true,
		},
		{
			str:     "",
			pattern: "*",
			match:   true,
		},
		{
			str:     "abc",
			pattern: "*",
			match:   true,
		},
		{
			str:     "a.c",
			pattern: "a.c",
			match:   true,
		},
		{
			str:     "a.c",
			pattern: "a\\.*c",
			match:   false,
		},
		{
			str:     "abc",
			pattern: "a\\*c",
			match:   false,
		},
		{
			str:     "a*c",
			pattern: "a\\*c",
			match:   false,
		},
		{
			str:     "abc",
			pattern: "a?c",
			match:   false,
		},
		{
			str:     "a?c",
			pattern: "a?c",
			match:   true,
		},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("String: %q, Pattern: %q", tt.str, tt.pattern)
		t.Run(testName, func(t *testing.T) {
			match, err := MatchPattern(tt.str, tt.pattern)
			if err != nil {
				t.Errorf("MatchPattern(%q, %q) returned error: %v", tt.str, tt.pattern, err)
			}
			if match != tt.match {
				t.Errorf("MatchPattern(%q, %q) = %v; want %v", tt.str, tt.pattern, match, tt.match)
			}
		})
	}
}

func TestContainsWildcard(t *testing.T) {
	tests := []struct {
		s        string
		expected bool
	}{
		{
			s:        "abc",
			expected: false,
		},
		{
			s:        "a*c",
			expected: true,
		},
		{
			s:        "*abc*",
			expected: true,
		},
		{
			s:        "a\\*c",
			expected: true,
		},
		{
			s:        "",
			expected: false,
		},
		{
			s:        "*",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("String: %q", tt.s), func(t *testing.T) {
			result := ContainsWildcard(tt.s)
			if result != tt.expected {
				t.Errorf("ContainsWildcard(%q) = %v; want %v", tt.s, result, tt.expected)
			}
		})
	}
}
