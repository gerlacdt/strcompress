package decompress

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type Kind int

const (
	Number  Kind = 0
	Letter  Kind = 1
	Bracket Kind = 2
	Empty   Kind = 3
)

// Token for the string decompression
type Token struct {
	kind  Kind
	value string
}

type Tokenizer struct {
	istring string
}

func isBracket(s string) bool {
	if s == "[" || s == "]" {
		return true
	}
	return false
}

func (t *Tokenizer) nextToken() (*Token, error) {
	if t.istring == "" {
		return &Token{kind: Empty, value: ""}, nil
	}
	first := t.istring[0]
	firstRune := rune(first)

	if unicode.IsDigit(firstRune) {
		buffer := ""
		index := 0
		for i := 0; i < len(t.istring) && unicode.IsDigit(rune(t.istring[i])); i++ {
			buffer += string(t.istring[i])
			index++
		}
		t.istring = t.istring[index:]
		return &Token{kind: Number, value: buffer}, nil
	}

	if unicode.IsLetter(firstRune) {
		buffer := ""
		index := 0
		for i := 0; i < len(t.istring) && unicode.IsLetter(rune(t.istring[i])); i++ {
			buffer += string(t.istring[i])
			index++
		}
		t.istring = t.istring[index:]
		return &Token{kind: Letter, value: buffer}, nil
	}

	if isBracket(string(first)) {
		t.istring = t.istring[1:]
		return &Token{kind: Bracket, value: string(first)}, nil
	}

	return nil, errors.New("Tokenizer error")
}

type Parser struct {
	tokenizer *Tokenizer
}

func (p *Parser) expression() (string, error) {
	token, err := p.tokenizer.nextToken()
	if err != nil {
		return "", fmt.Errorf("Error in expression getting token: %v", err)
	}

	if token.kind == Number {
		n, err := strconv.Atoi(token.value)
		if err != nil {
			return "", fmt.Errorf("Number could not be parsed, value: %s", token.value)
		}

		// match opening bracket
		token, err = p.tokenizer.nextToken()
		if err != nil {
			return "", fmt.Errorf("Parsing error, expect opening bracket, error %v", err)
		}
		if token.kind != Bracket && token.value != "[" {
			return "", fmt.Errorf("Parsing error, expected opening bracked, got %s", token.value)
		}

		// recursive call
		result, err := p.expression()

		final := ""
		for i := 0; i < n; i++ {
			final += result
		}

		// match closing bracket
		token, err = p.tokenizer.nextToken()
		if err != nil {
			return "", fmt.Errorf("Parsing error, expected closing bracket, error: %v", err)
		}
		if token.kind != Bracket && token.value != "]" {
			return "", fmt.Errorf("Parsing error, expected closing bracked, got %s", token.value)
		}

		// match letter or empty
		token, err = p.tokenizer.nextToken()
		if err != nil {
			return "", fmt.Errorf("Parsing error, expected letter or empty, error: %v", err)
		}
		if token.kind == Empty {
			return final, nil
		}
		if token.kind == Letter {
			return final + token.value, nil
		}
		if token.kind == Number {
			result2, err := p.expression()
			if err != nil {
				return "", fmt.Errorf("Parsing error, in following expression: %v", err)
			}
			return final + result2, nil
		}
	}

	if token.kind == Letter {
		return token.value, nil
	}

	return "", fmt.Errorf("Parser error, not match in expression, got token: %s", token.value)
}

func (p *Parser) decompress() (string, error) {
	result, err := p.expression()
	if err != nil {
		return "", err
	}
	return result, nil
}
