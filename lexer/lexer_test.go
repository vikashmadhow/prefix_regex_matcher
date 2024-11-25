// author: Vikash Madhow (vikash.madhow@gmail.com)

package lexer

import (
	"slices"
	"testing"
)

func TestLexer(t *testing.T) {
	l := New(
		&TokenDefinition{Type: "LET", Pattern: "let"},
		&TokenDefinition{Type: "INT", Pattern: "[0-9]+"},
		&TokenDefinition{Type: "ID", Pattern: "[_a-zA-Z][_a-zA-Z0-9]*"},
		&TokenDefinition{Type: "EQ", Pattern: "="},
		&TokenDefinition{Type: "SPC", Pattern: "[ \t\r\n]+"},
	)

	var tokens []Token
	for token := range l.lex("let x =  1000") {
		tokens = append(tokens, token)
	}

	if !slices.Equal(tokens, []Token{
		{"LET", "let", 0, 0},
		{"SPC", " ", 0, 0},
		{"ID", "x", 0, 0},
		{"SPC", " ", 0, 0},
		{"EQ", "=", 0, 0},
		{"SPC", "  ", 0, 0},
		{"INT", "1000", 0, 0},
	}) {
		t.Error("Invalid output", tokens)
	}
}
