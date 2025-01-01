// author: Vikash Madhow (vikash.madhow@gmail.com)

package lexer

import (
	"errors"
	"fmt"
	"github.com/vikashmadhow/prefix_regex_matcher/seq"
	"slices"
	"testing"
)

func TestBasicLexer(t *testing.T) {
	l := New(
		&TokenType{Id: "LET", Pattern: "let"},
		&TokenType{Id: "INT", Pattern: "[0-9]+"},
		&TokenType{Id: "ID", Pattern: "[_a-zA-Z][_a-zA-Z0-9]*"},
		&TokenType{Id: "EQ", Pattern: "="},
		&TokenType{Id: "SPC", Pattern: "\\s+"},
	)

	var tokens []Token
	//tokenSeq := l.LexText("let x =  1000")
	//defer tokenSeq.Stop()
	//for token := range seq.Push2(tokenSeq.Next, tokenSeq.Stop) {
	//	tokens = append(tokens, token)
	//}

	for token := range l.LexTextSeq("let x =  1000") {
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
		{EOF, "", 1, 14},
	}) {
		t.Error("Invalid output", tokens)
	}
}

func TestLexerError(t *testing.T) {
	l := New(
		&TokenType{Id: "LET", Pattern: "let"},
		&TokenType{Id: "INT", Pattern: "[0-9]+"},
		&TokenType{Id: "ID", Pattern: "[_a-zA-Z][_a-zA-Z0-9]*"},
		&TokenType{Id: "EQ", Pattern: "="},
		&TokenType{Id: "SPC", Pattern: "\\s+"},
	)

	var tokens []Token
	tokenSeq := l.LexText("let x? =  1000")
	for token, e := range seq.Push2(tokenSeq.Next, tokenSeq.Stop) {
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
		{EOF, "", 1, 14},
	}) {
		t.Error("Invalid output", tokens)
	}
}

func TestMultiline(t *testing.T) {
	l := New(
		&TokenType{Id: "LET", Pattern: "let"},
		&TokenType{Id: "INT", Pattern: "\\d+"},
		&TokenType{Id: "ID", Pattern: "[_a-zA-Z][_a-zA-Z0-9]*"},
		&TokenType{Id: "EQ", Pattern: "="},
		&TokenType{Id: "PLUS", Pattern: "\\+|-"},
		&TokenType{Id: "TIME", Pattern: "\\*|/"},
		&TokenType{Id: "SPC", Pattern: "\\s+"},
	)
	//l.Modulator(func(token regex.Token, err error) []seq.Pair[regex.Token, error] {
	//	if token.Type == "SPC" {
	//		return nil
	//	} else {
	//		return []seq.Pair[regex.Token, error]{{token, err}}
	//	}
	//})

	l.Modulator(Ignore("SPC"))

	var tokens []Token
	//tokenSeq := l.LexText(`let x = 1000
	//						 let y =x+y*-2000`)
	//for token := range seq.Push2(tokenSeq.Next, tokenSeq.Stop) {
	//	//if token.Type != "SPC" {
	//	tokens = append(tokens, token)
	//	//}
	//}

	for token := range l.LexTextSeq(`let x = 1000
							 let y =x+y*-2000`) {
		tokens = append(tokens, token)
	}

	fmt.Println(tokens)
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
		{EOF, "", 2, 25},
	})

	if err != nil {
		t.Error(err)
	}
}

func TestUnicode(t *testing.T) {
	l := New(
		&TokenType{Id: "LET", Pattern: "let"},
		&TokenType{Id: "INT", Pattern: "\\d+"},
		&TokenType{Id: "ID", Pattern: "[_a-zA-Z]\\S*"},
		&TokenType{Id: "EQ", Pattern: "="},
		&TokenType{Id: "PLUS", Pattern: "\\+|-"},
		&TokenType{Id: "TIME", Pattern: "\\*|/"},
		&TokenType{Id: "SPC", Pattern: "\\s+"},
	)
	l.Buffer(3)
	l.Modulator(Ignore("SPC"))

	var tokens []Token
	for token := range l.LexTextSeq(`let A日本語 = 1000`) {
		tokens = append(tokens, token)
	}

	fmt.Println(tokens)
	_, err := matchTokens(tokens, []Token{
		{"LET", "let", 1, 1},
		{"ID", "A日本語", 1, 5},
		{"EQ", "=", 1, 10},
		{"INT", "1000", 1, 12},
		{EOF, "", 1, 16},
	})

	if err != nil {
		t.Error(err)
	}
}

func TestReverse(t *testing.T) {
	l := New(
		&TokenType{Id: "LET", Pattern: "let"},
		&TokenType{Id: "INT", Pattern: "\\d+"},
		&TokenType{Id: "ID", Pattern: "[_a-zA-Z]\\S*"},
		&TokenType{Id: "EQ", Pattern: "="},
		&TokenType{Id: "PLUS", Pattern: "\\+|-"},
		&TokenType{Id: "TIME", Pattern: "\\*|/"},
		&TokenType{Id: "SPC", Pattern: "\\s+"},
	)
	//l.Buffer(3)
	l.Modulator(Ignore("SPC"), Reverse())

	var tokens []Token
	for token := range l.LexTextSeq(`let A日本語 = 1000`) {
		tokens = append(tokens, token)
	}

	fmt.Println(tokens)
	_, err := matchTokens(tokens, []Token{
		{"INT", "1000", 1, 12},
		{"EQ", "=", 1, 10},
		{"ID", "A日本語", 1, 5},
		{"LET", "let", 1, 1},
	})

	if err != nil {
		t.Error(err)
	}
}

func TestReverseAlternate(t *testing.T) {
	l := New(
		&TokenType{Id: "LET", Pattern: "let"},
		&TokenType{Id: "INT", Pattern: "\\d+"},
		&TokenType{Id: "ID", Pattern: "[_a-zA-Z]\\S*"},
		&TokenType{Id: "EQ", Pattern: "="},
		&TokenType{Id: "PLUS", Pattern: "\\+|-"},
		&TokenType{Id: "TIME", Pattern: "\\*|/"},
		&TokenType{Id: "SPC", Pattern: "\\s+"},
	)
	//l.Buffer(3)
	l.Modulator(Ignore("SPC"))

	l.Modulator(func() Modulator {
		var stream []seq.Pair[Token, error] = nil
		return func(t Token, err error) []seq.Pair[Token, error] {
			if t.Type == EOF {
				slices.Reverse(stream)
				return stream
			} else {
				stream = append(stream, seq.Pair[Token, error]{A: t, B: err})
				if len(stream) >= 2 {
					ret := make([]seq.Pair[Token, error], len(stream))
					copy(ret, stream)
					slices.Reverse(ret)
					stream = nil
					return ret
				}
				return nil
			}
		}
	}())

	var tokens []Token
	for token := range l.LexTextSeq(`let A日本語 = 1000 +`) {
		tokens = append(tokens, token)
	}

	fmt.Println(tokens)
	_, err := matchTokens(tokens, []Token{
		{"ID", "A日本語", 1, 5},
		{"LET", "let", 1, 1},
		{"INT", "1000", 1, 12},
		{"EQ", "=", 1, 10},
		{"PLUS", "+", 1, 17},
	})

	if err != nil {
		t.Error(err)
	}
}

func TestEndError(t *testing.T) {
	l := New(
		&TokenType{Id: "LET", Pattern: "let"},
		&TokenType{Id: "INT", Pattern: "\\d+"},
		&TokenType{Id: "ID", Pattern: "[_a-zA-Z][_a-zA-Z0-9]*"},
		&TokenType{Id: "EQ", Pattern: ":="},
		&TokenType{Id: "EQ_PLUS", Pattern: ":\\+"},
		&TokenType{Id: "PLUS", Pattern: "\\+|-"},
		&TokenType{Id: "TIME", Pattern: "\\*|/"},
		&TokenType{Id: "SPC", Pattern: "\\s+"},
	)

	var tokens []Token
	tokenSeq := l.LexText(`let x : 1000 :`)
	for token, e := range seq.Push2(tokenSeq.Next, tokenSeq.Stop) {
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
		{EOF, "", 1, 14},
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
