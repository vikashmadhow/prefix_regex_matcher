// author: Vikash Madhow (vikash.madhow@gmail.com)

// Parses regular expression to this grammar:
//
//	re -> re '|' re
//	    | re ('*' | '+' | '?')
//	    | re re
//	    | '(' re ')'
//	    | ch
//
//	ch -> '[' '^'? (c ['-' c])+ ']'
//	    | c
//	    | '\' ('*' | '+' | '?' | '|' | '(' | ')' | '[' | ']')
//
//	Refactored to remove left-recursion and ambiguity:
//	using: A = Aa|B  =>  A  = BA'
//	                     A' = aA'|e
//
//	regex  = term ['|' regex]
//	term   = { factor }
//	factor = base [('*' | '+' | '?')]
//	base   = '(' regex ')'
//	       | ch

package regex

import (
	"container/list"
	"math"
	"strconv"
	"strings"
	"sync/atomic"
)

// ----------------Regex top-down parsing----------------//
type parser struct {
	input    []rune
	position int
	group    *int
	groups   *list.List
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
		case '{':
			r.next()
			m := ""
			n := ""
			first := true
			if r.hasMore() {
				for r.hasMore() {
					c := r.next()
					if c == '}' {
						break
					}
					if c == ',' {
						first = false
					} else if first {
						m += string(c)
					} else {
						n += string(c)
					}
				}
				var mi, ma int32
				if len(strings.TrimSpace(m)) == 0 {
					mi = 0
				} else {
					x, err := strconv.Atoi(m)
					if err != nil {
						mi = 0
					} else {
						mi = int32(x)
					}
				}
				if len(strings.TrimSpace(n)) == 0 {
					ma = math.MaxUint8
				} else {
					x, err := strconv.Atoi(n)
					if err != nil {
						ma = math.MaxUint8
					} else {
						ma = int32(x)
					}
				}
				if first {
					ma = mi
				}
				mi = min(math.MaxUint8, max(0, mi))
				ma = min(math.MaxUint8, max(0, ma))
				if mi > ma {
					ma = atomic.SwapInt32(&mi, ma)
				}
				return &repeat{base, uint8(mi), uint8(ma)}
			} else {
				return &singleChar{'{', cp(r.groups)}
			}
		}
	}
	return base
}

func (r *parser) base() Regex {
	if r.peek() == '(' {
		r.next()

		*r.group++
		r.groups.PushBack(*r.group)

		re := r.regex()
		r.groups.Remove(r.groups.Back())

		// lenient parsing: don't break if no closing bracket, read to the end
		if r.hasMore() {
			r.next()
		}
		return &captureGroup{re}
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

		charSets := list.New()
		for r.hasMore() && r.peek() != ']' {
			from := r.next()
			if r.peek() == '-' {
				r.next()
				if r.hasMore() && r.peek() != ']' {
					to := r.next()
					charSets.PushBack(&charRange{from, to, cp(r.groups)})
				} else {
					charSets.PushBack(&charRange{from, math.MaxUint8, cp(r.groups)})
				}
			} else {
				charSets.PushBack(&singleChar{from, cp(r.groups)})
			}
		}
		// lenient parsing: don't break if no closing square bracket, read to the end
		if r.hasMore() {
			r.next()
		}
		return &characterSet{exclude, *charSets, cp(r.groups)}

	} else if r.peek() == '\\' {
		r.next()
		// lenient parsing: single backlash at the end is interpreted as escaping itself
		if r.hasMore() {
			switch c := r.next(); c {
			case 'd':
				return &charRange{'0', '9', cp(r.groups)}
			case 'D':
				cs := list.New()
				cs.PushBack(&charRange{'0', '9', cp(r.groups)})
				return &characterSet{true, *cs, cp(r.groups)}
			case 's':
				cs := list.New()
				cs.PushBack(&singleChar{' ', cp(r.groups)})
				cs.PushBack(&singleChar{'\t', cp(r.groups)})
				cs.PushBack(&singleChar{'\n', cp(r.groups)})
				cs.PushBack(&singleChar{'\f', cp(r.groups)})
				cs.PushBack(&singleChar{'\r', cp(r.groups)})
				return &characterSet{false, *cs, cp(r.groups)}
			case 'S':
				cs := list.New()
				cs.PushBack(&singleChar{' ', cp(r.groups)})
				cs.PushBack(&singleChar{'\t', cp(r.groups)})
				cs.PushBack(&singleChar{'\n', cp(r.groups)})
				cs.PushBack(&singleChar{'\f', cp(r.groups)})
				cs.PushBack(&singleChar{'\r', cp(r.groups)})
				return &characterSet{true, *cs, cp(r.groups)}
			case 'w':
				cs := list.New()
				cs.PushBack(&charRange{'0', '9', cp(r.groups)})
				cs.PushBack(&charRange{'a', 'z', cp(r.groups)})
				cs.PushBack(&charRange{'A', 'Z', cp(r.groups)})
				cs.PushBack(&singleChar{'_', cp(r.groups)})
				return &characterSet{false, *cs, cp(r.groups)}
			case 'W':
				cs := list.New()
				cs.PushBack(&charRange{'0', '9', cp(r.groups)})
				cs.PushBack(&charRange{'a', 'z', cp(r.groups)})
				cs.PushBack(&charRange{'A', 'Z', cp(r.groups)})
				cs.PushBack(&singleChar{'_', cp(r.groups)})
				return &characterSet{true, *cs, cp(r.groups)}
			default:
				return &singleChar{c, cp(r.groups)}
			}
		} else {
			return &singleChar{'\\', cp(r.groups)}
		}
	} else if r.peek() == '.' {
		r.next()
		return &anyChar{}
	} else {
		return &singleChar{r.next(), cp(r.groups)}
	}
}

func cp(groups *list.List) list.List {
	cp := list.New()
	for g := groups.Front(); g != nil; g = g.Next() {
		cp.PushBack(g.Value)
	}
	return *cp
}
