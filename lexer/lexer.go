// author: Vikash Madhow (vikash.madhow@gmail.com)

// Package lexer implements a simple fast lexer using the prefix regular expression matcher.
package lexer

import (
	"errors"
	"github.com/vikashmadhow/prefix_regex_matcher/regex"
	"iter"
	"strconv"
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
		var previousPartialMatches []*tokenMatcher
		for position < len(input) {
			matching = 0
			previousMatches = nil
			previousPartialMatches = nil
			r, n := utf8.DecodeRuneInString(input[position:])
			for _, m := range lexer.matchers {
				fillPrevious(m, &previousMatches, &previousPartialMatches)
				if m.matcher.LastMatch != regex.NoMatch {
					match := m.matcher.MatchNext(r)
					if match != regex.NoMatch {
						matching++
					}
				}
			}
			if matching == 0 {
				t, e := lexer.produceToken(previousMatches, previousPartialMatches, line, column)
				if !yield(t, e) || e != nil {
					return
				}
			} else {
				position += n
				if r == '\n' {
					line++
					column = 1
				} else {
					column++
				}
			}
		}
		for _, m := range lexer.matchers {
			fillPrevious(m, &previousMatches, &previousPartialMatches)
		}
		yield(lexer.produceToken(previousMatches, previousPartialMatches, line, column))
	}
}

func fillPrevious(m *tokenMatcher, previousMatches *[]*tokenMatcher, previousPartialMatches *[]*tokenMatcher) {
	if m.matcher.LastMatch == regex.FullMatch {
		*previousMatches = append(*previousMatches, m)
	} else if m.matcher.LastMatch == regex.PartialMatch {
		*previousPartialMatches = append(*previousPartialMatches, m)
	}
}

func (lexer *Lexer) produceToken(
	previousMatches []*tokenMatcher,
	previousPartialMatches []*tokenMatcher,
	line int, column int) (Token, error) {
	var token Token
	var err error
	if len(previousMatches) > 0 {
		match := previousMatches[0]
		token = Token{match.def.Type, match.matcher.Matched, line, column - len(match.matcher.Matched)}
		err = nil
	} else {
		token = Token{}
		msg := "error at " + strconv.Itoa(line) + ":" + strconv.Itoa(column)
		if len(previousPartialMatches) > 0 {
			msg += ": potential partial match(es): "
			for i, m := range previousPartialMatches {
				if i > 0 {
					msg += ", "
				}
				trans := m.matcher.Compiled.Dfa.Trans[m.matcher.State]
				msg += m.def.Type + " (next expected character(s): "
				first := true
				for k := range trans {
					if first {
						first = false
					} else {
						msg += ", "
					}
					msg += k.Pattern()
				}
				msg += ")"
			}
		}
		err = errors.New(msg)
	}
	for _, m := range lexer.matchers {
		m.matcher.Reset()
	}
	return token, err
}
