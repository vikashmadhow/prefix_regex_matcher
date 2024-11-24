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

func (r *CompiledRegex) Matcher() *Matcher {
	return &Matcher{Start, "", r, r.Dfa.start}
}

func (m *Matcher) Reset() {
	m.LastMatch = Start
	m.Matched = ""
	m.state = m.compiled.Dfa.start
}

func (m *Matcher) MatchNext(r rune) MatchType {
	trans := m.compiled.Dfa.trans[m.state]
	for c, t := range trans {
		if c.match(r) {
			m.state = t
			if slices.Index(m.compiled.Dfa.final, t) == -1 {
				m.LastMatch = PartialMatch
				m.Matched += string(r)
			} else {
				m.LastMatch = FullMatch
				m.Matched += string(r)
			}
			return m.LastMatch
		}
	}
	m.LastMatch = NoMatch
	return m.LastMatch
}
