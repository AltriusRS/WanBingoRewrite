package utils

import (
	"testing"
)

func TestCheckKeywordModeration(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{"Clean message", "Hello everyone!", true},
		{"Contains slur", "That guy is a nigger", false},
		{"Contains bypass attempt", "That guy is a n1gger", false},
		{"Contains fuck", "This is fucking awesome", false},
		{"Repeated chars", "Heeeeeeellooooooo", false},
		{"Excessive caps", "THIS IS ALL CAPS AND VERY LOUD", false},
		{"Normal caps", "This is a normal message", true},
		{"Short caps", "HI", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckKeywordModeration(tt.content)
			if result.Allowed != tt.expected {
				t.Errorf("CheckKeywordModeration(%q) = %v, expected %v", tt.content, result.Allowed, tt.expected)
			}
		})
	}
}

func TestCheckMarkdownModeration(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{"Clean message", "Hello everyone!", true},
		{"Inline formatting", "*italic* and **bold** text", true},
		{"Links allowed", "[click here](https://example.com)", true},
		{"Bare URL", "Check out https://example.com", true},
		{"Code inline", "Use `code` for inline code", true},
		{"Strikethrough", "~~strikethrough~~ text", true},
		{"Header rejected", "# This is a header", false},
		{"Blockquote rejected", "> This is a blockquote", false},
		{"Unordered list", "- This is a list item", false},
		{"Ordered list", "1. This is numbered", false},
		{"Horizontal rule", "---", false},
		{"Code block", "```\ncode here\n```", false},
		{"Table", "| Header | Header |\n|--------|--------|", false},
		{"Image rejected", "![alt text](image.jpg)", false},
		{"Image ref rejected", "![alt text][ref]", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckMarkdownModeration(tt.content)
			if result.Allowed != tt.expected {
				t.Errorf("CheckMarkdownModeration(%q) = %v, expected %v", tt.content, result.Allowed, tt.expected)
			}
		})
	}
}

func TestModerateContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{"Clean message", "Hello everyone!", true},
		{"Contains slur", "That guy is a nigger", false},
		{"Contains bypass attempt", "That guy is a n1gger", false},
		{"Repeated chars", "Heeeeeeellooooooo", false},
		{"Excessive caps", "THIS IS ALL CAPS AND VERY LOUD", false},
		{"Markdown header", "# Header", false},
		{"Markdown image", "![image](url)", false},
		{"Allowed link", "[link](url)", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ModerateContent(tt.content)
			if result.Allowed != tt.expected {
				t.Errorf("ModerateContent(%q) = %v, expected %v", tt.content, result.Allowed, tt.expected)
			}
		})
	}
}
