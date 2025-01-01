package lexer

import (
	"github.com/vikashmadhow/prefix_regex_matcher/seq"
)

// Modulator is a function that can modify the token stream produced by a lexer. It
// takes a token and maps it to a, potentially empty, sequence of tokens. Mapping to
// an empty sequence (or nil) takes out the token from the stream. Otherwise, the
// mapped sequence of tokens is merged into the token stream. This is similar to a
// flat-map operation on the token stream. Multiple modulators can be set on a lexer
// and they are invoked in the same the sequence that they were installed.
type Modulator func(Token, error) []seq.Pair[Token, error]

// IgnoreTokens is a Modulator that takes out the specified token types from the token
// stream. It is useful to remove syntactically useless tokens such as white space.
func IgnoreTokens(types ...string) Modulator {
	ignore := map[string]bool{}
	for _, t := range types {
		ignore[t] = true
	}
	return func(t Token, e error) []seq.Pair[Token, error] {
		if _, ok := ignore[t.Type]; ok {
			return nil
		} else {
			return []seq.Pair[Token, error]{{t, e}}
		}
	}
}
