package lexer

import (
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

func (t *TokenSeq) Next() (Token, error, bool) {
	if len(t.pushedBack) > 0 {
		//token := <- t.pushedBack[len(t.pushedBack)-1]
		//t.pushedBack = t.pushedBack[:len(t.pushedBack)-1]
		return *<-t.pushedBack, nil, true
	}
	return t.next()
}

func (t *TokenSeq) Stop() {
	t.stop()
}

func (t *TokenSeq) Pushback(token *Token) {
	//t.pushedBack = append(t.pushedBack, token)
	t.pushedBack <- token
}
