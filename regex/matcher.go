package regex

import "slices"

type MatchType int

const (
	NoMatch MatchType = iota
	PartialMatch
	FullMatch
)

type Matcher struct {
	compiled *CompiledRegex
	state    state
}

func (compiled *CompiledRegex) Matcher() Matcher {
	return Matcher{compiled, compiled.Dfa.start}
}

func (matcher *Matcher) Reset() {
	matcher.state = matcher.compiled.Dfa.start
}

func (matcher *Matcher) MatchNext(r rune) MatchType {
	trans := matcher.compiled.Dfa.trans[matcher.state]
	for c, t := range trans {
		if c.match(r) {
			matcher.state = t
			if slices.Index(matcher.compiled.Dfa.final, t) == -1 {
				return PartialMatch
			} else {
				return FullMatch
			}
		}
	}
	return NoMatch
}
