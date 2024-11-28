// author: Vikash Madhow (vikash.madhow@gmail.com)

// Package lexer implements a simple fast lexer using the prefix regular expression matcher.
package lexer

import (
	"errors"
	"github.com/vikashmadhow/prefix_regex_matcher/regex"
	"iter"
	"unicode/utf8"
)

type Token struct {
	Type   string
	Text   string
	Line   int
	Column int
}

type TokenDefinition struct {
	Type     string
	Pattern  string
	compiled regex.CompiledRegex
}

type tokenMatcher struct {
	def     *TokenDefinition
	matcher *regex.Matcher
}

type Lexer struct {
	definition []*TokenDefinition
	matchers   []*tokenMatcher
}

func New(definition ...*TokenDefinition) *Lexer {
	var matchers []*tokenMatcher
	for _, d := range definition {
		d.compiled = regex.NewRegex(d.Pattern)
		matchers = append(matchers, &tokenMatcher{d, d.compiled.Matcher()})
	}
	return &Lexer{definition, matchers}
}

func (lexer *Lexer) lex(input string) iter.Seq2[Token, error] {
	position := 0
	var line = 1
	var column = 1
	return func(yield func(t Token, e error) bool) {
		var matching int
		var previousMatches []*tokenMatcher
		for position < len(input) {
			matching = 0
			previousMatches = nil
			r, n := utf8.DecodeRuneInString(input[position:])
			if r == '\n' {
				line++
				column = 1
			}
			for _, m := range lexer.matchers {
				if m.matcher.LastMatch == regex.FullMatch {
					previousMatches = append(previousMatches, m)
				}
				if m.matcher.LastMatch != regex.NoMatch {
					match := m.matcher.MatchNext(r)
					if match != regex.NoMatch {
						matching++
					}
				}
			}
			if matching == 0 {
				t, e := lexer.produceToken(previousMatches, line, column)
				if !yield(t, e) || e != nil {
					return
				}
			} else {
				position += n
				column++
			}
		}
		yield(lexer.produceToken(previousMatches, line, column))
	}
}

func (lexer *Lexer) produceToken(previousMatches []*tokenMatcher, line int, column int) (Token, error) {
	var token Token
	var err error
	if len(previousMatches) > 0 {
		match := previousMatches[0]
		token = Token{match.def.Type, match.matcher.Matched, line, column - len(match.matcher.Matched)}
		err = nil
	} else {
		token = Token{}
		err = errors.New("no match")
	}
	for _, m := range lexer.matchers {
		m.matcher.Reset()
	}
	return token, err
}
