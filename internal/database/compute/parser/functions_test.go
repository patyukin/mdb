package parser

import (
	"fmt"
	"testing"
)

// Тест функции isWhitespace
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
			result := isWhitespace(tt.input)
			if result != tt.expected {
				t.Errorf("isWhitespace(%q) = %v; ожидалось %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Тест функции isValidArgChar
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
		{"Символ", '$', true},
		{"Непечатаемый", '\n', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidArgChar(tt.input)
			if result != tt.expected {
				t.Errorf("isValidArgChar(%q) = %v; ожидалось %v", tt.input, result, tt.expected)
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
		{"Запятая", ',', false},
		{"Восклицательный знак", '!', true},
		{"Вопросительный знак", '?', true},
		{"Дефис", '-', true},
		{"Кавычка", '"', false},
		{"Скобка", '(', false},
		{"Скобка 2", ')', false},
		{"Буква", 'A', false},
		{"Цифра", '1', false},
		{"Пробел", ' ', false},
		{"Символ", '$', true},
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
			result := containsWildcard(tt.s)
			if result != tt.expected {
				t.Errorf("containsWildcard(%q) = %v; want %v", tt.s, result, tt.expected)
			}
		})
	}
}
