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

var (
	Empty = TokenType{Id: "∅", Pattern: ""}
	EOF   = TokenType{Id: "Ω", Pattern: "$"}
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
	next       seq.Seq2[Token, error]
	stop       func()
	pushedBack []*Token
}

func (t *TokenSeq) Next() (Token, error, bool) {
	if len(t.pushedBack) > 0 {
		token := t.pushedBack[len(t.pushedBack)-1]
		t.pushedBack = t.pushedBack[:len(t.pushedBack)-1]
		return *token, nil, true
	}
	return t.next()
}

func (t *TokenSeq) Stop() {
	t.stop()
}

func (t *TokenSeq) Pushback(token *Token) {
	t.pushedBack = append(t.pushedBack, token)
}

func IgnoreTokens(types ...string) seq.FlatMap2Func[Token, error, Token, error] {
	ignore := map[string]bool{}
	for _, t := range types {
		ignore[t] = true
	}
	return func(t Token, e error) []seq.KeyValue[Token, error] {
		if _, ok := ignore[t.Type]; ok {
			return nil
		} else {
			return []seq.KeyValue[Token, error]{{t, e}}
		}
	}
}

type Lexer struct {
	Definition []*TokenType
	matchers   []*tokenMatcher
	bufferSize int
	modulators []seq.FlatMap2Func[Token, error, Token, error]
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

func (lexer *Lexer) Modulator(modulator ...seq.FlatMap2Func[Token, error, Token, error]) {
	lexer.modulators = append(lexer.modulators, modulator...)
}

//func (lexer *Lexer) Lex(input string) iter.Seq2[Token, error] {
//	position := 0
//	var line = 1
//	var column = 1
//	return func(yield func(t Token, e error) bool) {
//		var matching int
//		var previousMatches []*tokenMatcher
//		var previousPartialMatches []*tokenMatcher
//		for position < len(input) {
//			matching = 0
//			previousMatches = nil
//			previousPartialMatches = nil
//			r, n := utf8.DecodeRuneInString(input[position:])
//			for _, m := range lexer.matchers {
//				fillPrevious(m, &previousMatches, &previousPartialMatches)
//				if m.matcher.LastMatch != regex.NoMatch {
//					match := m.matcher.MatchNext(r)
//					if match != regex.NoMatch {
//						matching++
//					}
//				}
//			}
//			if matching == 0 {
//				t, e := lexer.produceToken(previousMatches, previousPartialMatches, line, column)
//				if !yield(t, e) || e != nil {
//					return
//				}
//			} else {
//				position += n
//				if r == '\n' {
//					line++
//					column = 1
//				} else {
//					column++
//				}
//			}
//		}
//		for _, m := range lexer.matchers {
//			fillPrevious(m, &previousMatches, &previousPartialMatches)
//		}
//		yield(lexer.produceToken(previousMatches, previousPartialMatches, line, column))
//	}
//}

func (lexer *Lexer) LexText(input string) *TokenSeq {
	return lexer.Lex(strings.NewReader(input))
}

//func (lexer *Lexer) Lex(in io.Reader) iter.Seq2[Token, error] {
//	var position int
//	var column int
//	line := 0
//
//	scanner := bufio.NewReader(in)
//	return func(yield func(t Token, e error) bool) {
//		var matching int
//		var previousMatches []*tokenMatcher
//		var previousPartialMatches []*tokenMatcher
//
//		for {
//			input, err := scanner.ReadString('\n')
//			// fmt.Print(input)
//
//			line++
//			position = 0
//			column = 1
//
//			for position < len(input) {
//				matching = 0
//				previousMatches = nil
//				previousPartialMatches = nil
//				r, n := utf8.DecodeRuneInString(input[position:])
//				for _, m := range lexer.matchers {
//					fillPrevious(m, &previousMatches, &previousPartialMatches)
//					if m.matcher.LastMatch != regex.NoMatch {
//						match := m.matcher.MatchNext(r)
//						if match != regex.NoMatch {
//							matching++
//						}
//					}
//				}
//				if matching == 0 {
//					t, e := lexer.produceToken(previousMatches, previousPartialMatches, line, column)
//					if !yield(t, e) || e != nil {
//						return
//					}
//				} else {
//					position += n
//					column++
//				}
//			}
//			for _, m := range lexer.matchers {
//				fillPrevious(m, &previousMatches, &previousPartialMatches)
//			}
//			yield(lexer.produceToken(previousMatches, previousPartialMatches, line, column))
//
//			if err != nil {
//				break
//			}
//		}
//	}
//}

func (lexer *Lexer) Lex(in io.Reader) *TokenSeq {
	next, stop := iter.Pull2(lexer.lex(in))
	if lexer.modulators != nil {
		for _, m := range lexer.modulators {
			next = seq.FlatMap2(next, m)
		}
	}
	return &TokenSeq{next: next, stop: stop}
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
