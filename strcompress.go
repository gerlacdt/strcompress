package strcompress

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

// Decompresses a given compressed string.
// You can find a detailed problem description here:
// https://techdevguide.withgoogle.com/paths/advanced/compress-decompression/#!

// Examples:
// 3[a]          ->  aaa
// 3[abc]4[ab]c  ->  abcabcabcababababc

// The solution implements a recursive-descent parser with the
// following grammar:

// <Exp> ::= <Number> '[' <Exp> ']' <Letter> | <Letter>
// <Number> ::= [0..9]+
// <Letter> ::= [a-z]*

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
	lookahead *Token
}

func (p *Parser) match(tok *Token) error {
	if tok.value != p.lookahead.value {
		return fmt.Errorf("Error token does not match: expected: %v, got: %v", p.lookahead, tok)
	}
	token, err := p.tokenizer.nextToken()
	if err != nil {
		return fmt.Errorf("Error in expression getting token: %v", err)
	}
	p.lookahead = token
	return nil
}

func (p *Parser) expression() (string, error) {
	if p.lookahead.kind == Number {
		n, err := strconv.Atoi(p.lookahead.value)
		if err != nil {
			return "", fmt.Errorf("Number could not be parsed, value: %s", p.lookahead.value)
		}

		// match number
		err = p.match(p.lookahead)
		if err != nil {
			return "", fmt.Errorf("Error matching number: %v", err)
		}

		// match opening bracket
		err = p.match(&Token{kind: Bracket, value: "["})
		if err != nil {
			return "", fmt.Errorf("Error matching opening bracket: %v", err)
		}

		// recursive call
		result, err := p.expression()
		if err != nil {
			return "", fmt.Errorf("Error in expression: %v", err)
		}

		final := ""
		for i := 0; i < n; i++ {
			final += result
		}
		// match closing bracket
		err = p.match(&Token{kind: Bracket, value: "]"})
		if err != nil {
			return "", fmt.Errorf("Error matching closing bracket: %v", err)
		}

		if p.lookahead.value == "]" {
			return final, nil
		}

		if p.lookahead.kind == Letter {
			value := p.lookahead.value
			err := p.match(p.lookahead)
			if err != nil {
				return "", fmt.Errorf("Error matching inner letter: %v", err)
			}
			return final + value, nil
		}

		if p.lookahead.kind == Number {
			result2, err := p.expression()
			if err != nil {
				return "", fmt.Errorf("Error parsing Number for after expression, lookahead: %v", p.lookahead)
			}
			return final + result2, nil
		}

		return final, nil
	}

	if p.lookahead.kind == Letter {
		value := p.lookahead.value
		err := p.match(p.lookahead)
		if err != nil {
			return "", fmt.Errorf("Error matching letter: %v", err)
		}
		return value, nil
	}

	return "", fmt.Errorf("Parser error, not in first, got token: %s", p.lookahead.value)
}

func (p *Parser) decompress() (string, error) {
	token, err := p.tokenizer.nextToken()
	if err != nil {
		return "", fmt.Errorf("Error in expression getting token: %v", err)
	}
	p.lookahead = token
	result, err := p.expression()
	if err != nil {
		return "", err
	}
	return result, nil
}
