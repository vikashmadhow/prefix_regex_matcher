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
				return &multiple{base, uint8(mi), uint8(ma)}
			} else {
				return &singleChar{'{'}
			}
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

		charSets := list.New()
		for r.hasMore() && r.peek() != ']' {
			from := r.next()
			if r.peek() == '-' {
				r.next()
				if r.hasMore() && r.peek() != ']' {
					to := r.next()
					charSets.PushBack(charRange{from, to})
				} else {
					charSets.PushBack(charRange{from, math.MaxUint8})
				}
			} else {
				charSets.PushBack(singleChar{from})
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
			switch c := r.next(); c {
			case 'd':
				return charRange{'0', '9'}
			case 'D':
				cs := list.New()
				cs.PushBack(charRange{'0', '9'})
				return &characterSet{true, cs}
			case 's':
				cs := list.New()
				cs.PushBack(singleChar{' '})
				cs.PushBack(singleChar{'\t'})
				cs.PushBack(singleChar{'\n'})
				cs.PushBack(singleChar{'\f'})
				cs.PushBack(singleChar{'\r'})
				return &characterSet{false, cs}
			case 'S':
				cs := list.New()
				cs.PushBack(singleChar{' '})
				cs.PushBack(singleChar{'\t'})
				cs.PushBack(singleChar{'\n'})
				cs.PushBack(singleChar{'\f'})
				cs.PushBack(singleChar{'\r'})
				return &characterSet{true, cs}
			case 'w':
				cs := list.New()
				cs.PushBack(charRange{'0', '9'})
				cs.PushBack(charRange{'a', 'z'})
				cs.PushBack(charRange{'A', 'Z'})
				cs.PushBack(singleChar{'_'})
				return &characterSet{false, cs}
			case 'W':
				cs := list.New()
				cs.PushBack(charRange{'0', '9'})
				cs.PushBack(charRange{'a', 'z'})
				cs.PushBack(charRange{'A', 'Z'})
				cs.PushBack(singleChar{'_'})
				return &characterSet{true, cs}
			default:
				return &singleChar{r.next()}
			}
		} else {
			return &singleChar{'\\'}
		}
	} else if r.peek() == '.' {
		r.next()
		return &anyChar{}
	} else {
		return &singleChar{r.next()}
	}
}
