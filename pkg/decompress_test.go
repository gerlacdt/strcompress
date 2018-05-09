package decompress

import (
	"errors"
	"fmt"
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
			fmt.Printf("result: %v, error: %v\n", result, err)
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
