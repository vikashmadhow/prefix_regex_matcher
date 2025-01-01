// author: Vikash Madhow (vikash.madhow@gmail.com)

// Package lexer implements a simple fast lexer using the prefix regular expression matcher.
package lexer

import (
	"bufio"
	"errors"
	"github.com/vikashmadhow/prefix_regex_matcher/regex"
	"github.com/vikashmadhow/prefix_regex_matcher/seq"
	"io"
	"iter"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Lexer struct {
	Definition []*TokenType
	matchers   []*tokenMatcher
	bufferSize int
	modulators []Modulator
}

type tokenMatcher struct {
	def     *TokenType
	matcher *regex.Matcher
}

func SimpleTokenType(id string) *TokenType {
	return NewTokenType(id, regex.Escape(id))
}

func NewTokenType(id string, pattern string) *TokenType {
	return &TokenType{id, pattern, regex.NewRegex(pattern)}
}

func New(definition ...*TokenType) *Lexer {
	var matchers []*tokenMatcher
	for _, d := range definition {
		if d.Compiled == nil {
			d.Compiled = regex.NewRegex(d.Pattern)
		}
		matchers = append(matchers, &tokenMatcher{d, d.Compiled.Matcher()})
	}
	return &Lexer{definition, matchers, 1024, nil}
}

func (lexer *Lexer) Buffer(size int) {
	lexer.bufferSize = size
}

func (lexer *Lexer) Modulator(modulator ...Modulator) {
	lexer.modulators = append(lexer.modulators, modulator...)
}

func (lexer *Lexer) LexText(input string) *TokenSeq {
	return lexer.Lex(strings.NewReader(input))
}

func (lexer *Lexer) LexTextSeq(input string) iter.Seq2[Token, error] {
	return lexer.LexSeq(strings.NewReader(input))
}

func (lexer *Lexer) Lex(in io.Reader) *TokenSeq {
	next, stop := iter.Pull2(lexer.lex(in))
	if lexer.modulators != nil {
		for _, m := range lexer.modulators {
			next = seq.FlatMap2(next, m)
		}
	}
	return &TokenSeq{next: next, stop: stop}
}

func (lexer *Lexer) LexSeq(in io.Reader) iter.Seq2[Token, error] {
	next := lexer.lex(in)
	if lexer.modulators != nil {
		for _, m := range lexer.modulators {
			next = seq.FlatMapSeq2(next, m)
		}
	}
	return next
}

func (lexer *Lexer) lex(in io.Reader) iter.Seq2[Token, error] {
	column, line := 1, 1
	scanner := bufio.NewReader(in)

	return func(yield func(t Token, e error) bool) {
		var matching int
		var previousMatches []*tokenMatcher
		var previousPartialMatches []*tokenMatcher

		start := 0
		bufferSize := lexer.bufferSize
		if bufferSize < 8 {
			bufferSize = 8
		}
		input := make([]byte, bufferSize)
		for {
			read, err := scanner.Read(input[start:])
			read += start
			start = 0

			//fmt.Println(read)
			//fmt.Println(string(input[:read]))

			for position := 0; position < read; {
				matching = 0
				previousMatches = nil
				previousPartialMatches = nil

				r, n := utf8.DecodeRune(input[position:read])
				if r == utf8.RuneError {
					start = copy(input[0:], input[position:read])
					break
				}
				//fmt.Println("  >>", string(r))

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
			if err != nil {
				//fmt.Println(err)
				yield(lexer.produceToken(previousMatches, previousPartialMatches, line, column))
				break
			}
		}
		yield(Token{Type: EOF, Text: "", Line: line, Column: column}, nil)
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
		token = Token{match.def.Id, match.matcher.Matched, line, column - utf8.RuneCountInString(match.matcher.Matched)}
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
				msg += m.def.Id + " (next expected character(s): "
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
