package regex

import (
	"math"
)

type CompiledRegex struct {
	Regex Regex
	Nfa   *automata
	Dfa   *automata
}

// NewRegex creates a new regular expression from the input
func NewRegex(input string) CompiledRegex {
	parser := parser{[]rune(input), 0}
	r := parser.regex()
	n := r.nfa()
	d := dfa(n)
	return CompiledRegex{r, n, d}
}

// ----------------Regex top-down parsing----------------//
type parser struct {
	input    []rune
	position int
}

func (r *parser) peek() rune {
	if r.position < len(r.input) {
		return r.input[r.position]
	} else {
		return 0
	}
}

func (r *parser) next() rune {
	c := r.input[r.position]
	r.position++
	return c
}

func (r *parser) hasMore() bool {
	return len(r.input) > r.position
}

func (r *parser) regex() Regex {
	term := r.term()
	if r.hasMore() && r.peek() == '|' {
		r.next()
		right := r.regex()
		return &choice{term, right}
	} else {
		return term
	}
}

func (r *parser) term() Regex {
	var factors []Regex
	for r.hasMore() && r.peek() != ')' && r.peek() != '|' {
		factors = append(factors, r.factor())
	}
	return &sequence{factors}
}

func (r *parser) factor() Regex {
	base := r.base()
	if r.hasMore() {
		switch r.peek() {
		case '*':
			r.next()
			return &zeroOrMore{base}
		case '+':
			r.next()
			return &oneOrMore{base}
		case '?':
			r.next()
			return &zeroOrOne{base}
		}
	}
	return base
}

func (r *parser) base() Regex {
	if r.peek() == '(' {
		r.next()
		re := r.regex()
		// lenient parsing: don't break if no closing bracket, read to the end
		if r.hasMore() {
			r.next()
		}
		return &group{re}
	} else {
		return r.ch()
	}
}

func (r *parser) ch() Regex {
	if r.peek() == '[' {
		r.next()

		exclude := false
		if r.peek() == '^' {
			r.next()
			exclude = true
		}

		var charSets []char
		for r.hasMore() && r.peek() != ']' {
			from := r.next()
			if r.peek() == '-' {
				r.next()
				if r.hasMore() && r.peek() != ']' {
					to := r.next()
					charSets = append(charSets, &charRange{from, to})
				} else {
					charSets = append(charSets, &charRange{from, math.MaxUint8})
				}
			} else {
				charSets = append(charSets, &singleChar{from})
			}
		}
		// lenient parsing: don't break if no closing square bracket, read to the end
		if r.hasMore() {
			r.next()
		}
		return &characterSet{exclude, charSets}

	} else if r.peek() == '\\' {
		r.next()
		// lenient parsing: single backlash at the end is interpreted as escaping itself
		if r.hasMore() {
			return &singleChar{r.next()}
		} else {
			return &singleChar{'\\'}
		}

	} else {
		return &singleChar{r.next()}
	}
}
