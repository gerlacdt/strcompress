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

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return true
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
		// match opening bracket
		token, err = p.tokenizer.nextToken()
		if err != nil {
			return "", fmt.Errorf("Parsing error, expect opening bracket, error %v", err)
		}
		if token.kind != Bracket && token.value != "[" {
			return "", fmt.Errorf("Parsing error, expected opening bracked, got %s, token.value")
		}

		// recusive call
		result, err := p.expression()

		// match closing bracket
		token, err = p.tokenizer.nextToken()
		if err != nil {
			return "", fmt.Errorf("Parsing error, expected closing bracket, error: %v", err)
		}
		if token.kind != Bracket && token.value != "]" {
			return "", fmt.Errorf("Parsing error, expected closing bracked, got %s, token.value")
		}

		// match letter or empty
		token, err = p.tokenizer.nextToken()
		if err != nil {
			return "", fmt.Errorf("Parsing error, expected letter or empty, error: %v", err)
		}
		if token.kind == Empty {
			return result, nil
		}
		if token.kind == Letter {
			return result + token.value, nil
		}
	}

	if token.kind == Letter {
		return token.value, nil
	}

	return "", fmt.Errorf("Parser error, not match in expression, got token: %s", token.value)
}

func (p *Parser) decompress(s string) (string, error) {
	result, err := p.expression()
	if err != nil {
		return "", err
	}
	return result, nil
}
