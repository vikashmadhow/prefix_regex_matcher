// author: Vikash Madhow (vikash.madhow@gmail.com)

package lexer

import (
	"fmt"
	"slices"
	"testing"
)

func TestBasicLexer(t *testing.T) {
	l := New(
		&TokenDefinition{Type: "LET", Pattern: "let"},
		&TokenDefinition{Type: "INT", Pattern: "[0-9]+"},
		&TokenDefinition{Type: "ID", Pattern: "[_a-zA-Z][_a-zA-Z0-9]*"},
		&TokenDefinition{Type: "EQ", Pattern: "="},
		&TokenDefinition{Type: "SPC", Pattern: "\\s+"},
	)

	var tokens []Token
	for token := range l.lex("let x =  1000") {
		tokens = append(tokens, token)
	}

	if !slices.Equal(tokens, []Token{
		{"LET", "let", 1, 1},
		{"SPC", " ", 1, 4},
		{"ID", "x", 1, 5},
		{"SPC", " ", 1, 6},
		{"EQ", "=", 1, 7},
		{"SPC", "  ", 1, 8},
		{"INT", "1000", 1, 10},
	}) {
		t.Error("Invalid output", tokens)
	}
}

func TestLexerError(t *testing.T) {
	l := New(
		&TokenDefinition{Type: "LET", Pattern: "let"},
		&TokenDefinition{Type: "INT", Pattern: "[0-9]+"},
		&TokenDefinition{Type: "ID", Pattern: "[_a-zA-Z][_a-zA-Z0-9]*"},
		&TokenDefinition{Type: "EQ", Pattern: "="},
		&TokenDefinition{Type: "SPC", Pattern: "\\s+"},
	)

	var tokens []Token
	for token, e := range l.lex("let x? =  1000") {
		if e != nil {
			println(e.Error())
			return
		}
		fmt.Println(token)
		tokens = append(tokens, token)
	}

	if !slices.Equal(tokens, []Token{
		{"LET", "let", 1, 1},
		{"SPC", " ", 1, 4},
		{"ID", "x", 1, 5},
		{"SPC", " ", 1, 6},
		{"EQ", "=", 1, 7},
		{"SPC", "  ", 1, 8},
		{"INT", "1000", 1, 10},
	}) {
		t.Error("Invalid output", tokens)
	}
}
