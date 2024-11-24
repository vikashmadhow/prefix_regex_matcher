package regex

import (
	"testing"
)

func TestEmpty(t *testing.T) {
	r := NewRegex("")
	println(r.Dfa.ToGraphViz("Empty string"))
	if !r.match("") {
		t.Error("'' did not match ''")
	}
	if r.match("a") {
		t.Error("'' matched 'a'")
	}
}

func TestSingleChar(t *testing.T) {
	r := NewRegex("a")
	if !r.match("a") {
		t.Error("'a' did not match 'a'")
	}
	if r.match("b") {
		t.Error("'a' matched 'b'")
	}
	if r.match("aa") {
		t.Error("'a' matched 'aa'")
	}
	if r.match("") {
		t.Error("'a' matched ''")
	}
}

func TestSequence(t *testing.T) {
	r := NewRegex("abc")
	if !r.match("abc") {
		t.Error("'abc' did not match 'abc'")
	}
	if r.match("ab") {
		t.Error("'abc' matched 'ab'")
	}
	if r.match("abcabc") {
		t.Error("'abc' matched 'abcabc'")
	}
	if r.match("") {
		t.Error("'abc' matched ''")
	}
}

func TestChoice(t *testing.T) {
	r := NewRegex("a|b")
	if !r.match("a") {
		t.Error("'a|b' did not match 'a'")
	}
	if !r.match("b") {
		t.Error("'a|b' did not match 'b'")
	}
	if r.match("ab") {
		t.Error("'a|b' matched 'ab'")
	}
}

func TestSequenceChoice(t *testing.T) {
	r := NewRegex("ab|ac")
	if !r.match("ab") {
		t.Error("'ab|ac' did not match 'ab'")
	}
	if !r.match("ac") {
		t.Error("'ab|ac' did not match 'ac'")
	}
	if r.match("abac") {
		t.Error("'ab|ac' matched 'abac'")
	}
}

func TestOpt(t *testing.T) {
	r := NewRegex("a?")
	if !r.match("a") {
		t.Error("'a?' did not match 'a'")
	}
	if !r.match("") {
		t.Error("'a?' did not match ''")
	}
	if r.match("aa") {
		t.Error("'a?' matched 'aa'")
	}
}

func TestSequenceOpt(t *testing.T) {
	r := NewRegex("(ab)?")
	if !r.match("ab") {
		t.Error("'(ab)?' did not match 'ab'")
	}
	if !r.match("") {
		t.Error("'(ab)?' did not match ''")
	}
	if r.match("abab") {
		t.Error("'(ab)?' matched 'abab'")
	}
}

func TestSequenceOpt2(t *testing.T) {
	r := NewRegex("ab?")
	if !r.match("ab") {
		t.Error("'ab?' did not match 'ab'")
	}
	if !r.match("a") {
		t.Error("'ab?' did not match 'a'")
	}
	if r.match("") {
		t.Error("'ab?' matched ''")
	}
	if r.match("abab") {
		t.Error("'(ab)?' matched 'abab'")
	}
}

func TestZeroOrMore(t *testing.T) {
	r := NewRegex("a*")
	if !r.match("a") {
		t.Error("'a*' did not match 'a'")
	}
	if !r.match("aa") {
		t.Error("'a*' did not match 'aa'")
	}
	if !r.match("aaa") {
		t.Error("'a*' did not match 'aaa'")
	}
	if !r.match("") {
		t.Error("'a*' did not match ''")
	}
}

func TestZeroOrMoreSequence(t *testing.T) {
	r := NewRegex("(ab)*")
	if !r.match("ab") {
		t.Error("'(ab)*' did not match 'ab'")
	}
	if !r.match("abab") {
		t.Error("'(ab)*' did not match 'abab'")
	}
	if !r.match("ababab") {
		t.Error("'(ab)*' did not match 'ababab'")
	}
	if !r.match("") {
		t.Error("'(ab)*' did not match ''")
	}
	if r.match("a") {
		t.Error("'(ab)*' matched 'a'")
	}
	if r.match("b") {
		t.Error("'(ab)*' matched 'b'")
	}
	if r.match("aba") {
		t.Error("'(ab)*' matched 'aba'")
	}
}

func TestOneOrMore(t *testing.T) {
	r := NewRegex("a+")
	if !r.match("a") {
		t.Error("'a+' did not match 'a'")
	}
	if !r.match("aa") {
		t.Error("'a+' did not match 'aa'")
	}
	if !r.match("aaa") {
		t.Error("'a+' did not match 'aaa'")
	}
	if r.match("") {
		t.Error("'a+' matched ''")
	}
}

func TestOneOrMoreSequence(t *testing.T) {
	r := NewRegex("(ab)+")
	if !r.match("ab") {
		t.Error("'(ab)+' did not match 'ab'")
	}
	if !r.match("abab") {
		t.Error("'(ab)+' did not match 'abab'")
	}
	if !r.match("ababab") {
		t.Error("'(ab)+' did not match 'ababab'")
	}
	if r.match("") {
		t.Error("'(ab)+' matched ''")
	}
	if r.match("a") {
		t.Error("'(ab)+' matched 'a'")
	}
	if r.match("b") {
		t.Error("'(ab)+' matched 'b'")
	}
	if r.match("aba") {
		t.Error("'(ab)+' matched 'aba'")
	}
}

func TestRepeat(t *testing.T) {
	r := NewRegex("(ab|ac){5,3}")
	if !r.match("abacab") {
		t.Error("'(ab|ac){3,5}' did not match 'abacab'")
	}
	if !r.match("abacabab") {
		t.Error("'(ab|ac){3,5}' did not match 'abacabab'")
	}
	if !r.match("abacababac") {
		t.Error("'(ab|ac){3,5}' did not match 'abacababac'")
	}
	if r.match("abacababacab") {
		t.Error("'(ab|ac){3,5}' matched 'abacababacab'")
	}
	if r.match("") {
		t.Error("'(ab|ac){3,5}' matched ''")
	}
}

func TestRepeatExact(t *testing.T) {
	r := NewRegex("(ab|ac){3}")
	if !r.match("abacab") {
		t.Error("'(ab|ac){3}' did not match 'abacab'")
	}
	if r.match("abacabab") {
		t.Error("'(ab|ac){3}' matched 'abacabab'")
	}
	if r.match("abacababac") {
		t.Error("'(ab|ac){3}' matched 'abacababac'")
	}
	if r.match("abacababacab") {
		t.Error("'(ab|ac){3}' matched 'abacababacab'")
	}
	if r.match("") {
		t.Error("'(ab|ac){3}' matched ''")
	}
}

func TestRepeatNoUpperLimit(t *testing.T) {
	r := NewRegex("(ab|ac){3,}")
	if !r.match("abacab") {
		t.Error("'(ab|ac){3,}' did not match 'abacab'")
	}
	if !r.match("abacabab") {
		t.Error("'(ab|ac){3,}' did not match 'abacabab'")
	}
	if !r.match("abacababac") {
		t.Error("'(ab|ac){3,}' did not match 'abacababac'")
	}
	if !r.match("abacababacab") {
		t.Error("'(ab|ac){3,}' did not match 'abacababacab'")
	}
	if r.match("") {
		t.Error("'(ab|ac){3,}' matched ''")
	}
}

func TestRepeatNoLowerLimit(t *testing.T) {
	r := NewRegex("(ab|ac){,3}")
	if !r.match("") {
		t.Error("'(ab|ac){,3}' did not match ''")
	}
	if !r.match("abac") {
		t.Error("'(ab|ac){,3}' did not match 'abac'")
	}
	if !r.match("abacab") {
		t.Error("'(ab|ac){,3}' did not match 'abacab'")
	}
	if r.match("abacabab") {
		t.Error("'(ab|ac){,3}' matched 'abacabab'")
	}
	if r.match("abacababac") {
		t.Error("'(ab|ac){,3}' matched 'abacababac'")
	}
}
