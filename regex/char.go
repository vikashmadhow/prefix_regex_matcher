// author: Vikash Madhow (vikash.madhow@gmail.com)

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
	groups() *list.List // [int]
	setGroups(g *list.List)
	Regex
}

type empty struct{ _ uint8 }

type anyChar struct {
	group list.List
}

type singleChar struct {
	char  rune
	group list.List
}

type charRange struct {
	from  rune
	to    rune
	group list.List
}

type characterSet struct {
	exclude  bool
	charSets list.List // [char]
	group    list.List // [int]
}

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
	return c.char == char
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
	return c.from <= char && char <= c.to
}

//------------- A character set combines different characters (and ranges) -------------//

func (c *characterSet) Pattern() string {
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
	//ret += "]:" + label(c.groups())
	//ret += ")"
	return ret
}

func (c *characterSet) isEmpty() bool {
	return false
}

func (c *characterSet) groups() *list.List {
	return &c.group
}

func (c *characterSet) setGroups(g *list.List) {
	c.group = *g
}

func (c *characterSet) nfa() *automata {
	return charNfa(c)
}

func (c *characterSet) match(ch rune) bool {
	if c.exclude {
		for cs := c.charSets.Front(); cs != nil; cs = cs.Next() {
			if cs.Value.(char).match(ch) {
				return false
			}
		}
		return true
	} else {
		for cs := c.charSets.Front(); cs != nil; cs = cs.Next() {
			if cs.Value.(char).match(ch) {
				return true
			}
		}
		return false
	}
}

func charNfa(c char) *automata {
	a := automata{
		Trans: make(transitions),
		start: &stateObj{},
		final: []state{&stateObj{}},
	}
	addTransitions(&a, a.start, map[char]state{c: a.final[0]})
	return &a
}
