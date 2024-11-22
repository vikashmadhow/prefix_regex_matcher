package regex

// The char interface represents the different type of character and
// character sets used in regular expressions. They become the label
// in the NFA and DFA generated to recognise strings with the regular
// expressions.

import (
	"container/list"
	"math"
)

// -------------Character and character sets parsing-------------//
type char interface {
	match(c rune) bool
	isEmpty() bool
	Regex
}

type empty struct{ _ uint8 }

type singleChar struct {
	char rune
}

type charRange struct {
	from rune
	to   rune
}

type characterSet struct {
	exclude  bool
	charSets *list.List // []char
}

//------------- The empty character -------------//

func (c empty) Pattern() string {
	return ""
}

func (c empty) isEmpty() bool {
	return true
}

func (c empty) nfa() *automata {
	return nil
}

func (c empty) match(rune) bool {
	return false
}

//------------- A single character match -------------//

func (c singleChar) Pattern() string {
	return string(c.char)
}

func (c singleChar) isEmpty() bool {
	return false
}

func (c singleChar) nfa() *automata {
	return charNfa(c)
}

func (c singleChar) match(char rune) bool {
	return c.char == char
}

//------------- A character range match -------------//

func (c charRange) Pattern() string {
	if c.to < math.MaxUint8 {
		return string(c.from) + "-" + string(c.to)
		//return "Range(" + string(c.from) + "-" + string(c.to) + ")"
	} else {
		return string(c.from) + "-"
		//return "Range(" + string(c.from) + "-)"
	}
}

func (c charRange) isEmpty() bool {
	return false
}

func (c charRange) nfa() *automata {
	return charNfa(c)
}

func (c charRange) match(char rune) bool {
	return c.from <= char && char <= c.to
}

//------------- A character set combines different characters (and ranges) -------------//

func (c characterSet) Pattern() string {
	ret := "["
	//ret := "CharSet("
	if c.exclude {
		ret += "^"
	}
	first := true
	//for _, cs := range c.charSets {
	for cs := c.charSets.Front(); cs != nil; cs = cs.Next() {
		if first {
			first = false
		} else {
			ret += "|"
		}
		ret += cs.Value.(char).Pattern()
	}
	ret += "]"
	//ret += ")"
	return ret
}

func (c characterSet) isEmpty() bool {
	return false
}

func (c characterSet) nfa() *automata {
	return charNfa(c)
}

func (c characterSet) match(ch rune) bool {
	//for _, cs := range c.charSets {
	for cs := c.charSets.Front(); cs != nil; cs = cs.Next() {
		if (c.exclude && !cs.Value.(char).match(ch)) ||
			(!c.exclude && cs.Value.(char).match(ch)) {
			return true
		}
	}
	return false
}

func charNfa(c char) *automata {
	a := automata{
		trans: make(transitions),
		start: &stateObj{},
		final: []state{&stateObj{}},
	}
	addTransitions(&a, a.start, map[char]state{c: a.final[0]})
	return &a
}
