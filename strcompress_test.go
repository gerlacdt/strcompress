package strcompress

import (
	"errors"
	"testing"
)

func TestNextToken(t *testing.T) {
	tt := []struct {
		input    string
		expected string
		rest     string
		err      error
	}{
		{"123", "123", "", nil},
		{"abc", "abc", "", nil},
		{"456abc", "456", "abc", nil},
		{"abc123", "abc", "123", nil},
		{"[]", "[", "]", nil},
		{"[abc]", "[", "abc]", nil},
		{"", "", "", errors.New("EOF")},
	}
	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			tokenizer := Tokenizer{istring: tc.input}
			result, err := tokenizer.nextToken()
			if err != nil && tc.err == nil {
				t.Fatalf("Error nextToken: %v", err)
			}
			if result == nil {
				result = &Token{kind: Number, value: ""}
			}
			if result.value != tc.expected {
				t.Errorf("expected: %s, got: %s", tc.expected, result.value)
			}
			if tokenizer.istring != tc.rest {
				t.Errorf("expected rest: %s, got: %s", tc.rest, tokenizer.istring)
			}
		})
	}
}

func TestDecompress(t *testing.T) {
	tt := []struct {
		input    string
		expected string
		err      error
	}{
		{"a", "a", nil},
		{"3[a]", "aaa", nil},
		{"2[3[a]b]", "aaabaaab", nil},
		{"3[abc]4[ab]c", "abcabcabcababababc", nil},
		{"1[1[1[xx]]]", "xx", nil},
		{"0[a]", "", nil},
		{"1[]", "", nil},
	}
	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			tokenizer := &Tokenizer{istring: tc.input}
			parser := &Parser{tokenizer: tokenizer}
			result, err := parser.decompress()
			if err != nil {
				t.Fatalf("Error in decompress, %v", err)
			}
			if result != tc.expected {
				t.Errorf("expected: %s, got: %s", tc.expected, result)
			}
		})
	}
}
