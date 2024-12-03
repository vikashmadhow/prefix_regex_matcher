// author: Vikash Madhow (vikash.madhow@gmail.com)

package lexer

import (
	"errors"
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

func TestMultiLine(t *testing.T) {
	l := New(
		&TokenDefinition{Type: "LET", Pattern: "let"},
		&TokenDefinition{Type: "INT", Pattern: "\\d+"},
		&TokenDefinition{Type: "ID", Pattern: "[_a-zA-Z][_a-zA-Z0-9]*"},
		&TokenDefinition{Type: "EQ", Pattern: "="},
		&TokenDefinition{Type: "PLUS", Pattern: "\\+|-"},
		&TokenDefinition{Type: "TIME", Pattern: "\\*|/"},
		&TokenDefinition{Type: "SPC", Pattern: "\\s+"},
	)

	var tokens []Token
	for token := range l.lex(`let x = 1000
							 let y =x+y*-2000`) {
		if token.Type != "SPC" {
			tokens = append(tokens, token)
		}
	}

	_, err := matchTokens(tokens, []Token{
		{"LET", "let", 1, 1},
		//{"SPC", " ", 1, 4},
		{"ID", "x", 1, 5},
		//{"SPC", " ", 1, 6},
		{"EQ", "=", 1, 7},
		//{"SPC", " ", 1, 8},
		{"INT", "1000", 1, 9},
		//{"SPC", "\n\t\t\t\t\t\t\t ", 2, 0},
		{"LET", "let", 2, 9},
		//{"SPC", " ", 2, 12},
		{"ID", "y", 2, 13},
		//{"SPC", " ", 2, 14},
		{"EQ", "=", 2, 15},
		{"ID", "x", 2, 16},
		{"PLUS", "+", 2, 17},
		{"ID", "y", 2, 18},
		{"TIME", "*", 2, 19},
		{"PLUS", "-", 2, 20},
		{"INT", "2000", 2, 21},
	})

	if err != nil {
		t.Error(err)
	}
}

func TestEndError(t *testing.T) {
	l := New(
		&TokenDefinition{Type: "LET", Pattern: "let"},
		&TokenDefinition{Type: "INT", Pattern: "\\d+"},
		&TokenDefinition{Type: "ID", Pattern: "[_a-zA-Z][_a-zA-Z0-9]*"},
		&TokenDefinition{Type: "EQ", Pattern: ":="},
		&TokenDefinition{Type: "EQ_PLUS", Pattern: ":\\+"},
		&TokenDefinition{Type: "PLUS", Pattern: "\\+|-"},
		&TokenDefinition{Type: "TIME", Pattern: "\\*|/"},
		&TokenDefinition{Type: "SPC", Pattern: "\\s+"},
	)

	var tokens []Token
	for token, e := range l.lex(`let x : 1000 :`) {
		if e != nil {
			println(e.Error())
			return
		}
		if token.Type != "SPC" {
			tokens = append(tokens, token)
		}
	}

	_, err := matchTokens(tokens, []Token{
		{"LET", "let", 1, 1},
		{"ID", "x", 1, 5},
		{"EQ", ":=", 1, 7},
		{"INT", "1000", 1, 10},
	})

	if err != nil {
		t.Error(err)
	}
}

func matchTokens(t1 []Token, t2 []Token) (bool, error) {
	if len(t1) != len(t2) {
		return false, errors.New(fmt.Sprint("comparing different number of tokens:", len(t1), ",", len(t2)))
	}
	for i, token := range t1 {
		if t2[i] != token {
			return false, errors.New(fmt.Sprint("failed at position:", i, ",", token, "!=", t2[i]))
		}
	}
	return true, nil
}
