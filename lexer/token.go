package lexer

import (
	"errors"
	"github.com/vikashmadhow/prefix_regex_matcher/regex"
	"github.com/vikashmadhow/prefix_regex_matcher/seq"
)

var (
	Empty = "∅"
	EOF   = "Ω"
)

type Token struct {
	Type   string
	Text   string
	Line   int
	Column int
}

type TokenType struct {
	Id       string
	Pattern  string
	Compiled *regex.CompiledRegex
}

type TokenSeq struct {
	next seq.Seq2[Token, error]
	stop func()
	//pushedBack []*Token
	pushedBack chan *Token
}

func (t *TokenSeq) Next() (*Token, error) {
	if len(t.pushedBack) > 0 {
		//token := <- t.pushedBack[len(t.pushedBack)-1]
		//t.pushedBack = t.pushedBack[:len(t.pushedBack)-1]
		return <-t.pushedBack, nil
	}
	//return t.next()
	token, err, valid := t.next()
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("lexer returned an invalid token")
	}
	return &token, nil

}

func (t *TokenSeq) Peek() (*Token, error) {
	token, err := t.Next()
	if err != nil {
		return nil, err
	}
	t.Pushback(token)
	return token, nil
}

func (t *TokenSeq) Pushback(token *Token) {
	//t.pushedBack = append(t.pushedBack, token)
	t.pushedBack <- token
}
func (t *TokenSeq) Stop() {
	t.stop()
}
