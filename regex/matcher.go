package regex

import "slices"

type MatchType int

const (
	// NoMatch Last match was unsuccessful, the matcher must be reset for valid subsequent matching
	NoMatch MatchType = iota

	// PartialMatch Text matched a prefix of the regular expression, but has not reached a final state
	// to be considered a full match.
	PartialMatch

	// FullMatch A full match was achieved; subsequent supplied characters can still result in a full
	// match if the longer string is part of the regular language of the regular expression.
	FullMatch

	// Start Matching has not started yet. The matcher is set to this state on creation or reset.
	Start
)

type Matcher struct {
	LastMatch MatchType
	Matched   string
	compiled  *CompiledRegex
	state     state
}

func (compiled *CompiledRegex) Matcher() *Matcher {
	return &Matcher{Start, "", compiled, compiled.Dfa.start}
}

func (matcher *Matcher) Reset() {
	matcher.LastMatch = Start
	matcher.Matched = ""
	matcher.state = matcher.compiled.Dfa.start
}

func (matcher *Matcher) MatchNext(r rune) MatchType {
	trans := matcher.compiled.Dfa.trans[matcher.state]
	for c, t := range trans {
		if c.match(r) {
			matcher.state = t
			if slices.Index(matcher.compiled.Dfa.final, t) == -1 {
				matcher.LastMatch = PartialMatch
				matcher.Matched += string(r)
			} else {
				matcher.LastMatch = FullMatch
				matcher.Matched += string(r)
			}
			return matcher.LastMatch
		}
	}
	matcher.LastMatch = NoMatch
	return matcher.LastMatch
}
