// author: Vikash Madhow (vikash.madhow@gmail.com)

package regex

// The char interface represents the different type of character and
// character sets used in regular expressions. They become the label
// in the NFA and DFA generated to recognize strings with the regular
// expressions.

import (
	"container/list"
	"math"
	"unicode"
)

// -------------Character and character sets parsing-------------//
type (
	char interface {
		match(c rune) bool
		isEmpty() bool
		groups() *list.List // [int]
		setGroups(g *list.List)

		// spanSet returns the range of characters that can be matched by this char.
		spanSet() spanSet

		Regex
	}

	empty struct{ _ uint8 }

	anyChar struct {
		mod   *modifier
		group list.List
	}

	singleChar struct {
		mod   *modifier
		char  rune
		group list.List
	}

	charRange struct {
		mod   *modifier
		from  rune
		to    rune
		group list.List
	}

	charSet struct {
		mod     *modifier
		exclude bool
		sets    list.List // [char]
		group   list.List // [int]
	}
)

//------------- The empty character -------------//

func (c *empty) Pattern() string {
	return ""
}

func (c *empty) isEmpty() bool {
	return true
}

func (c *empty) groups() *list.List {
	return nil
}

func (c *empty) setGroups(g *list.List) {}

func (c *empty) nfa() *automata {
	return nil
}

func (c *empty) match(rune) bool {
	return false
}

func (c *empty) spanSet() spanSet {
	return nil
}

//------------- Any character -------------//

func (c *anyChar) Pattern() string {
	return "."
	//return ".:" + label(c.groups())
}

func (c *anyChar) isEmpty() bool {
	return false
}

func (c *anyChar) groups() *list.List {
	return &c.group
}

func (c *anyChar) setGroups(g *list.List) {
	c.group = *g
}

func (c *anyChar) nfa() *automata {
	return charNfa(c)
}

func (c *anyChar) match(rune) bool {
	return true
}

func (c *anyChar) spanSet() spanSet {
	if c.mod.unicode {
		return allUnicode
	} else {
		return asciiPrintable
	}
}

//------------- A single character match -------------//

func (c *singleChar) Pattern() string {
	return string(c.char)
	//return string(c.char) + ":" + label(c.groups())
}

func (c *singleChar) isEmpty() bool {
	return false
}

func (c *singleChar) groups() *list.List {
	return &c.group
}

func (c *singleChar) setGroups(g *list.List) {
	c.group = *g
}

func (c *singleChar) nfa() *automata {
	return charNfa(c)
}

func (c *singleChar) match(char rune) bool {
	if c.mod.caseInsensitive {
		return unicode.ToLower(char) == unicode.ToLower(c.char)
	} else {
		return c.char == char
	}
}

func (c *singleChar) spanSet() spanSet {
	if c.mod.caseInsensitive {
		l := unicode.ToLower(c.char)
		u := unicode.ToUpper(c.char)
		if l != u {
			return spanSet{
				{l, l},
				{u, u},
			}
		}
	}
	return spanSet{
		{c.char, c.char},
	}
}

//------------- A character range match -------------//

func (c *charRange) Pattern() string {
	if c.to < math.MaxUint8 {
		return string(c.from) + "-" + string(c.to)
		//return string(c.from) + "-" + string(c.to) + ":" + label(c.groups())
		//return "Range(" + string(c.from) + "-" + string(c.to) + ")"
	} else {
		return string(c.from) + "-"
		//return string(c.from) + "-" + ":" + label(c.groups())
		//return "Range(" + string(c.from) + "-)"
	}
}

func (c *charRange) isEmpty() bool {
	return false
}

func (c *charRange) groups() *list.List {
	return &c.group
}

func (c *charRange) setGroups(g *list.List) {
	c.group = *g
}

func (c *charRange) nfa() *automata {
	return charNfa(c)
}

func (c *charRange) match(char rune) bool {
	if c.mod.caseInsensitive {
		lf := unicode.ToLower(c.from)
		uf := unicode.ToUpper(c.from)

		lt := unicode.ToLower(c.to)
		ut := unicode.ToUpper(c.to)

		if lf != uf || lt != ut {
			return (lf <= char && char <= lt) || (uf <= char && char <= ut)
		}
	}
	return c.from <= char && char <= c.to
}

func (c *charRange) spanSet() spanSet {
	if c.mod.caseInsensitive {
		lf := unicode.ToLower(c.from)
		uf := unicode.ToUpper(c.from)

		lt := unicode.ToLower(c.to)
		ut := unicode.ToUpper(c.to)

		if lf != uf && lt != ut {
			return spanSet{
				{lf, lt},
				{uf, ut},
			}
		}
	}
	return spanSet{
		{c.from, c.to},
	}
}

//------------- A character set combines different characters (and ranges) -------------//

func (c *charSet) Pattern() string {
	ret := "["
	//ret := "CharSet("
	if c.exclude {
		ret += "^"
	}
	first := true
	//for _, cs := range c.charSets {
	for cs := c.sets.Front(); cs != nil; cs = cs.Next() {
		if first {
			first = false
		} else {
			ret += "|"
		}
		ret += cs.Value.(char).Pattern()
	}
	ret += "]"
	//ret += "]:" + label(c.groups())
	//ret += ")"
	return ret
}

func (c *charSet) isEmpty() bool {
	return false
}

func (c *charSet) groups() *list.List {
	return &c.group
}

func (c *charSet) setGroups(g *list.List) {
	c.group = *g
}

func (c *charSet) nfa() *automata {
	return charNfa(c)
}

func (c *charSet) match(ch rune) bool {
	matched := false
	for cs := c.sets.Front(); cs != nil; cs = cs.Next() {
		if cs.Value.(char).match(ch) {
			matched = true
			break
		}
	}
	if c.exclude {
		return !matched
	}
	return matched
}

func (c *charSet) spanSet() spanSet {
	var span spanSet
	for cs := c.sets.Front(); cs != nil; cs = cs.Next() {
		span = append(span, cs.Value.(char).spanSet()...)
	}
	if c.exclude {
		if c.mod.unicode {
			return span.invertUnicode()
		} else {
			return span.invertAsciiPrintable()
		}
	} else {
		return span.compact()
	}
}
