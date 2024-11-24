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

func TestDigits(t *testing.T) {
	r := NewRegex("\\d{3,5}")
	if !r.match("569") {
		t.Error("'\\d{3,5}' did not match '569'")
	}
	if !r.match("5697") {
		t.Error("'\\d{3,5}' did not match '5697'")
	}
	if !r.match("56975") {
		t.Error("'\\d{3,5}' did not match '56975'")
	}
	if r.match("569751") {
		t.Error("'\\d{3,5}' matched '569751'")
	}
	if r.match("5bc") {
		t.Error("'\\d{3,5}' matched '5bc'")
	}
}

func TestNonDigits(t *testing.T) {
	r := NewRegex("\\D{3,5}")
	if !r.match("abF") {
		t.Error("'\\D{3,5}' did not match 'abF'")
	}
	if !r.match("abFs") {
		t.Error("'\\D{3,5}' did not match 'abFs'")
	}
	if !r.match("abFs?") {
		t.Error("'\\D{3,5}' did not match 'abFs?'")
	}
	if r.match("abFs?;") {
		t.Error("'\\D{3,5}' matched 'abFs?;'")
	}
	if r.match("5bc") {
		t.Error("'\\D{3,5}' matched '5bc'")
	}
}

func TestSpaces(t *testing.T) {
	r := NewRegex("\\s{3,5}")
	if !r.match(" 	 ") {
		t.Error("'\\s{3,5}' did not match ' 	 '")
	}
	if !r.match("  		") {
		t.Error("'\\s{3,5}' did not match '  		'")
	}
	if !r.match("  		 ") {
		t.Error("'\\s{3,5}' did not match '  		 '")
	}
	if r.match("  		  ") {
		t.Error("'\\s{3,5}' matched '  		  '")
	}
	if r.match("5  ") {
		t.Error("'\\s{3,5}' matched '5  '")
	}
}

func TestNonSpaces(t *testing.T) {
	r := NewRegex("\\S{3,5}")
	if !r.match("abc") {
		t.Error("'\\S{3,5}' did not match 'abc'")
	}
	if !r.match("abcd") {
		t.Error("'\\S{3,5}' did not match 'abcd'")
	}
	if !r.match("abcde") {
		t.Error("'\\S{3,5}' did not match 'abcde'")
	}
	if r.match("abcdef") {
		t.Error("'\\S{3,5}' matched 'abcdef'")
	}
	if r.match("   ") {
		t.Error("'\\S{3,5}' matched '   '")
	}
}

func TestWords(t *testing.T) {
	r := NewRegex("\\w{3,5}")
	if !r.match("ab0") {
		t.Error("'\\w{3,5}' did not match 'ab0'")
	}
	if !r.match("ab01") {
		t.Error("'\\w{3,5}' did not match 'ab01'")
	}
	if !r.match("abc01") {
		t.Error("'\\w{3,5}' did not match 'abc01'")
	}
	if r.match("abc012") {
		t.Error("'\\w{3,5}' matched 'abc012'")
	}
	if r.match("?bc") {
		t.Error("'\\w{3,5}' matched '?bc'")
	}
}

func TestNonWords(t *testing.T) {
	r := NewRegex("\\W{3,5}")
	if !r.match("<>?") {
		t.Error("'\\W{3,5}' did not match '<>?'")
	}
	if !r.match("<>?,") {
		t.Error("'\\W{3,5}' did not match '<>?,'")
	}
	if !r.match("<>?,.") {
		t.Error("'\\W{3,5}' did not match '<>?,.'")
	}
	if r.match("<>?,./") {
		t.Error("'\\W{3,5}' matched '<>?,./'")
	}
	if r.match("A<>") {
		t.Error("'\\W{3,5}' matched 'A<>'")
	}
}

func TestDot(t *testing.T) {
	r := NewRegex(".{3,5}")
	if !r.match("^*k") {
		t.Error("'.{3,5}' did not match '^*k'")
	}
	if !r.match("^*k)") {
		t.Error("'.{3,5}' did not match '^*k'")
	}
	if !r.match("^*k)$") {
		t.Error("'.{3,5}' did not match '^*k)$'")
	}
	if r.match("^*k)$d") {
		t.Error("'.{3,5}' matched '^*k)$d'")
	}
	if r.match("") {
		t.Error("'.{3,5}' matched ''")
	}
}
