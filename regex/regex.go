// Package regex implements a regular expression library that can be
// supplied with a sequence of characters and will return if the characters at
// that point is a prefix matching the regular expression.
//
// Supports regular expressions conforming to this EBNF definition:
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

// Regex is the base visible interface of regular expressions
type Regex interface {
    Pattern() string
    nfa() *automata
}

// choice represents the regex | regex rule
type choice struct {
    left  Regex
    right Regex
}

// sequence represents a sequence of regular expressions (a, b, ...)
type sequence struct {
    sequence []Regex
}

// zeroOrOne is for an optional regular expression (re?)
type zeroOrOne struct {
    opt Regex
}

// zeroOrMore is for the Kleene closure (re*)
type zeroOrMore struct {
    repeat Regex
}

// oneOrMore is for positive closure (re+)
type oneOrMore struct {
    repeat Regex
}

// group is for grouping regular expressions inside brackets, i.e., (re)
type group struct {
    re Regex
}

//-----------------Regex interface methods------------//

func (c *choice) Pattern() string {
    //return c.left.Pattern() + "|" + c.right.Pattern()
    return "Or(" + c.left.Pattern() + ", " + c.right.Pattern() + ")"
}

// automata constructs a finite automaton for the choice (union) of two regular expressions.
//
//			      left
//			     ∧    \
//		        /      v
//			start      final
//			     \     ∧
//	              v   /
//			      right
func (c *choice) nfa() *automata {
    a := automata{
        trans: make(transitions),
        start: &stateObj{},
        final: []state{&stateObj{}},
    }

    left := c.left.nfa()
    right := c.right.nfa()

    merge(&a, left)
    merge(&a, right)

    addTransitions(&a, a.start, map[char]state{&empty{}: left.start, &empty{}: right.start})
    addTransitions(&a, left.final[0], map[char]state{&empty{}: a.final[0]})
    addTransitions(&a, right.final[0], map[char]state{&empty{}: a.final[0]})

    return &a
}

func (s *sequence) Pattern() string {
    //ret := ""
    ret := "Seq("
    first := true
    for _, re := range s.sequence {
        if first {
            first = false
        } else {
            ret += ", "
        }
        ret += re.Pattern()
    }
    ret += ")"
    return ret
}

// automata constructs a finite-state automaton for the sequence of regular expressions.
// It merges the individual automata of each regular expression in the sequence, connecting
// the final state of one to the start state of the next. It returns a pointer to the resulting automata.
//
//	start --> re1 in sequence --> re2 --> .... --> final
func (s *sequence) nfa() *automata {
    a := automata{
        trans: make(transitions),
        start: &stateObj{},
        final: []state{&stateObj{}},
    }

    first := true
    for _, re := range s.sequence {
        reAutomata := re.nfa()
        merge(&a, reAutomata)
        if first {
            a.start = reAutomata.start
            first = false
        } else {
            addTransitions(&a, a.final[0], map[char]state{&empty{}: reAutomata.start})
        }
        a.final = reAutomata.final
    }
    return &a
}

func (r *zeroOrOne) Pattern() string {
    //return r.opt.Pattern() + "?"
    return "?(" + r.opt.Pattern() + ")"
}

// automata constructs and returns an NFA for an optional subpattern.
//
//			 _______________
//			/               \
//		   /                 v
//	     start --> ... --> final
func (r *zeroOrOne) nfa() *automata {
    opt := r.opt.nfa()
    addTransitions(opt, opt.start, map[char]state{&empty{}: opt.final[0]})
    return opt
}

func merge(target *automata, source *automata) *automata {
    for k, v := range source.trans {
        target.trans[k] = v
    }
    return target
}

func addTransitions(target *automata, from state, to map[char]state) *automata {
    existing, ok := target.trans[from]
    if !ok {
        target.trans[from] = to
    } else {
        for k, v := range to {
            existing[k] = v
        }
    }
    return target
}

func (r *zeroOrMore) Pattern() string {
    //return r.repeat.Pattern() + "*"
    return "*(" + r.repeat.Pattern() + ")"
}

// automata generates a finite automaton for a zero-or-more repetition (Kleene closure) of the pattern.
//
//			 ______________
//			^              \
//		   /                v
//	     start --> ... --> final
//	       ^                /
//		    \              v
//		     --------------
func (r *zeroOrMore) nfa() *automata {
    repeat := r.repeat.nfa()
    addTransitions(repeat, repeat.start, map[char]state{&empty{}: repeat.final[0]})
    addTransitions(repeat, repeat.final[0], map[char]state{&empty{}: repeat.start})
    return repeat
}

func (r *oneOrMore) Pattern() string {
    //return r.repeat.Pattern() + "+"
    return "+(" + r.repeat.Pattern() + ")"
}

// automata generates a finite automaton for a zero-or-more repetition (Kleene closure) of the pattern.
//
//	     start --> ... --> final
//	      ^                  /
//		   \                v
//			 ---------------
func (r *oneOrMore) nfa() *automata {
    repeat := r.repeat.nfa()
    addTransitions(repeat, repeat.final[0], map[char]state{&empty{}: repeat.start})
    return repeat
}

func (r *group) Pattern() string {
    //return "(" + r.re.Pattern() + ")"
    return "Grp(" + r.re.Pattern() + ")"
}

func (r *group) nfa() *automata {
    return r.re.nfa()
}
