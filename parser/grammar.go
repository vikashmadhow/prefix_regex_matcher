// package grammar

package grammar

import (
	"errors"
	"fmt"
	"github.com/vikashmadhow/prefix_regex_matcher/lexer"
	"io"
	"maps"
	"os"
	"strconv"
	"strings"
)

type Grammar struct {
	Id          string
	Lexer       *lexer.Lexer
	Productions []*Production
	ProdByName  map[string]*Production
}

type Production struct {
	Name     string
	Sentence Sentence
	follow   map[string]bool
}

type LanguageElement interface {
	Terminal() bool
	MatchEmpty(*Grammar) bool
	First(*Grammar, CycleDetector) (map[string]bool, error)
	Recognise(*Grammar, LanguageElement, *lexer.TokenSeq, CycleDetector) (*SyntaxTree, error)
	ToString() string
}

type Sentence interface {
	Follow(*Grammar, string, CycleDetector) (map[string]bool, bool, error)
	LanguageElement
}

type TokenLanguageElement struct {
	Token *lexer.Token
}

type TokenRef struct {
	Ref string
}

type ProductionRef struct {
	Ref string
}

type Sequence struct {
	Elements []Sentence
	first    map[string]bool
}

type Choice struct {
	Alternates []Sentence
	first      map[string]bool
}

type Optional struct {
	Sentence Sentence
}

type ZeroOrMore struct {
	Sentence Sentence
}

type OneOrMore struct {
	Sentence Sentence
}

type Repeat struct {
	Min, Max int
	Sentence Sentence
	first    map[string]bool
	follow   map[string]bool
}

type CycleDetector interface {
	add(LanguageElement) error
	remove(LanguageElement)
}

type CycleDetectorSet struct {
	Seen map[LanguageElement]bool
}

func (c *CycleDetectorSet) add(a LanguageElement) error {
	if _, ok := c.Seen[a]; ok {
		return errors.New("cycle detected containing " + a.ToString())
	}
	c.Seen[a] = true
	return nil
}

func (c *CycleDetectorSet) remove(a LanguageElement) {
	delete(c.Seen, a)
}

// --- TOKEN REFERENCE --- //

func (t *TokenRef) Terminal() bool {
	return true
}

func (t *TokenRef) First(_ *Grammar, _ CycleDetector) (map[string]bool, error) {
	return map[string]bool{t.Ref: true}, nil
}

func (t *TokenRef) Follow(_ *Grammar, _ string, _ CycleDetector) (map[string]bool, bool, error) {
	return map[string]bool{}, false, nil
}

func (t *TokenRef) MatchEmpty(g *Grammar) bool {
	tokenType := g.Lexer.TokenTypes[t.Ref]
	return tokenType.Compiled.MatchEmpty()
}

func (t *TokenRef) Recognise(_ *Grammar, _ LanguageElement, tokens *lexer.TokenSeq, _ CycleDetector) (*SyntaxTree, error) {
	token, err, valid := tokens.Next()
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("")
	}
	if t.Ref == token.Type {
		return &SyntaxTree{&TokenLanguageElement{&token}, nil}, nil
	}
	return nil, errors.New("token type " + token.Type + " does not match expected type " + t.Ref)
}

func (t *TokenRef) ToString() string {
	return t.Ref
}

// --- Token as a language element --- //

func (t *TokenLanguageElement) Terminal() bool {
	return true
}

func (t *TokenLanguageElement) First(_ *Grammar, _ CycleDetector) (map[string]bool, error) {
	return map[string]bool{t.Token.Type: true}, nil
}

func (t *TokenLanguageElement) Follow(_ *Grammar, _ string, _ CycleDetector) (map[string]bool, bool, error) {
	return map[string]bool{}, false, nil
}

func (t *TokenLanguageElement) MatchEmpty(g *Grammar) bool {
	tokenType := g.Lexer.TokenTypes[t.Token.Type]
	return tokenType.Compiled.MatchEmpty()
}

func (t *TokenLanguageElement) Recognise(_ *Grammar, _ LanguageElement, tokens *lexer.TokenSeq, _ CycleDetector) (*SyntaxTree, error) {
	token, err, valid := tokens.Next()
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("")
	}
	if t.Token.Type == token.Type {
		return &SyntaxTree{&TokenLanguageElement{&token}, nil}, nil
	}
	return nil, errors.New("token type " + token.Type + " does not match expected type " + t.Token.Type)
}

func (t *TokenLanguageElement) ToString() string {
	return fmt.Sprint(*t.Token)
}

// --- Production reference (in a sentence) --- //

func (p *ProductionRef) Terminal() bool {
	return false
}

func (p *ProductionRef) First(g *Grammar, cd CycleDetector) (map[string]bool, error) {
	err := cd.add(p)
	if err != nil {
		return nil, err
	}
	f, e := g.ProdByName[p.Ref].Sentence.First(g, cd)
	cd.remove(p)
	return f, e
}

func (p *ProductionRef) Follow(_ *Grammar, _ string, _ CycleDetector) (map[string]bool, bool, error) {
	return map[string]bool{}, false, nil
}

func (p *ProductionRef) MatchEmpty(g *Grammar) bool {
	prod := g.ProdByName[p.Ref]
	return prod.MatchEmpty(g)
}

func (p *ProductionRef) Recognise(g *Grammar, production LanguageElement, tokens *lexer.TokenSeq, cd CycleDetector) (*SyntaxTree, error) {
	prod := g.ProdByName[p.Ref]
	return prod.Recognise(g, production, tokens, cd)
}

func (p *ProductionRef) ToString() string {
	return p.Ref
}

// --- Production  --- //

// s : a b
// x : b s A
// y : s t?
// t : X
// z : y s

// follow(s) = {A, X, V}
// b s A
// s X
// s V

func (p *Production) Terminal() bool {
	return false
}

func (p *Production) First(g *Grammar, cd CycleDetector) (map[string]bool, error) {
	err := cd.add(p)
	if err != nil {
		return nil, err
	}
	f, e := p.Sentence.First(g, cd)
	cd.remove(p)
	return f, e
}

func (p *Production) Follow(g *Grammar, cd CycleDetector) (map[string]bool, error) {
	if p.follow != nil {
		return p.follow, nil
	}
	//if _, ok := (*cd)[p.Name]; ok {
	//	return nil, fmt.Errorf("recursive reference through production %s", p.Name)
	//}
	p.follow = make(map[string]bool)
	for _, prod := range g.Productions {
		if prod.Name != p.Name {
			f, emptyTillEnd, err := prod.Sentence.Follow(g, p.Name, cd)
			if err != nil {
				return nil, err
			}
			maps.Insert(p.follow, maps.All(f))
			if emptyTillEnd {
				f, err = prod.Follow(g, cd)
				if err != nil {
					return nil, err
				}
				maps.Insert(p.follow, maps.All(f))
			}
		}
	}
	//delete(*cd, p.Name)
	return p.follow, nil
}

func (p *Production) MatchEmpty(g *Grammar) bool {
	return p.Sentence.MatchEmpty(g)
}

func next(tokens *lexer.TokenSeq) (*lexer.Token, error) {
	token, err, valid := tokens.Next()
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("lexer returned an invalid token")
	}
	return &token, nil
}

func peek(tokens *lexer.TokenSeq) (*lexer.Token, error) {
	token, err := next(tokens)
	if err != nil {
		return nil, err
	}
	tokens.Pushback(token)
	return token, nil
}

func (p *Production) Recognise(g *Grammar, _ LanguageElement, tokens *lexer.TokenSeq, cd CycleDetector) (*SyntaxTree, error) {
	token, err := peek(tokens)
	if err != nil {
		return nil, err
	}
	first, err := p.First(g, cd)
	if err != nil {
		return nil, err
	}
	if _, ok := first[token.Type]; ok {
		return p.Sentence.Recognise(g, p, tokens, cd)
	} else if p.Sentence.MatchEmpty(g) {
		follow, err := p.Follow(g, cd)
		if err != nil {
			return nil, err
		}
		if _, ok := follow[token.Type]; ok {
			return nil, nil
		} else {
			return nil, fmt.Errorf("unexpected token %v", token.Type)
		}
	} else {
		return nil, fmt.Errorf("unexpected token %v", token.Type)
	}
}

func (p *Production) ToString() string {
	return p.Name + ": " + p.Sentence.ToString()
}

// --- Choice --- //

func (c *Choice) Terminal() bool {
	return false
}

func (c *Choice) First(g *Grammar, cd CycleDetector) (map[string]bool, error) {
	if c.first != nil {
		return c.first, nil
	}
	err := cd.add(c)
	if err != nil {
		return nil, err
	}
	c.first = map[string]bool{}
	for _, a := range c.Alternates {
		f, e := a.First(g, cd)
		if e != nil {
			return nil, e
		}
		maps.Insert(c.first, maps.All(f))
	}
	cd.remove(c)
	return c.first, nil
}

func (c *Choice) Follow(g *Grammar, production string, cd CycleDetector) (map[string]bool, bool, error) {
	follow := make(map[string]bool)
	emptyTillEnd := false
	for _, a := range c.Alternates {
		f, empty, err := a.Follow(g, production, cd)
		if err != nil {
			return nil, false, err
		}
		maps.Insert(follow, maps.All(f))
		emptyTillEnd = emptyTillEnd || empty
	}
	return follow, emptyTillEnd, nil
}

func (c *Choice) MatchEmpty(g *Grammar) bool {
	for _, a := range c.Alternates {
		if a.MatchEmpty(g) {
			return true
		}
	}
	return false
}

func (c *Choice) Recognise(g *Grammar, production LanguageElement, tokens *lexer.TokenSeq, cd CycleDetector) (*SyntaxTree, error) {
	token, err := peek(tokens)
	if err != nil {
		return nil, err
	}

	var alternate Sentence
	for _, a := range c.Alternates {
		first, err := a.First(g, cd)
		if err != nil {
			return nil, err
		}
		if _, ok := first[token.Type]; ok {
			alternate = a
			break
		}
	}
	if alternate == nil {
		return nil, fmt.Errorf("no alternates found for choice %q on token %q", c.ToString(), token.Type)
	}
	return alternate.Recognise(g, production, tokens, cd)
}

func (c *Choice) ToString() string {
	s := ""
	first := true
	for _, a := range c.Alternates {
		if first {
			first = false
		} else {
			s += " | "
		}
		s += a.ToString()
	}
	return s
}

// --- SEQUENCE --- //

func (s *Sequence) Terminal() bool {
	return false
}

func (s *Sequence) First(g *Grammar, cd CycleDetector) (map[string]bool, error) {
	if s.first != nil {
		return s.first, nil
	}
	err := cd.add(s)
	if err != nil {
		return nil, err
	}
	s.first = map[string]bool{}
	for _, e := range s.Elements {
		f, er := e.First(g, cd)
		if er != nil {
			return nil, er
		}
		maps.Insert(s.first, maps.All(f))
		if !e.MatchEmpty(g) {
			break
		}
	}
	cd.remove(s)
	return s.first, nil
}

func (s *Sequence) Follow(g *Grammar, production string, cd CycleDetector) (map[string]bool, bool, error) {
	found := -1
search:
	for i, e := range s.Elements {
		switch p := e.(type) {
		case *ProductionRef:
			if p.Ref == production {
				found = i
				break search
			}
		}
	}

	follow := make(map[string]bool)
	emptyTillEnd := true
	if found != -1 {
		for _, e := range s.Elements[found+1:] {
			first, err := e.First(g, cd)
			if err != nil {
				return nil, false, err
			}
			maps.Insert(follow, maps.All(first))
			if !e.MatchEmpty(g) {
				emptyTillEnd = false
				break
			}
		}
	}
	return follow, emptyTillEnd, nil
}

func (s *Sequence) MatchEmpty(g *Grammar) bool {
	for _, e := range s.Elements {
		if !e.MatchEmpty(g) {
			return false
		}
	}
	return true
}

func (s *Sequence) Recognise(g *Grammar, production LanguageElement, tokens *lexer.TokenSeq, cd CycleDetector) (*SyntaxTree, error) {
	tree := SyntaxTree{production, nil}
	for _, e := range s.Elements {
		token, err := peek(tokens)
		if err != nil {
			return nil, err
		}
		first, err := e.First(g, cd)
		if err != nil {
			return nil, err
		}
		if _, ok := first[token.Type]; ok {
			child, err := e.Recognise(g, production, tokens, cd)
			if err != nil {
				return nil, err
			}
			tree.Children = append(tree.Children, child)
		} else if !e.MatchEmpty(g) {
			return nil, fmt.Errorf("token %q cannot start %q", token.Type, e.ToString())
		}
	}
	if len(tree.Children) == 1 {
		return tree.Children[0], nil
	} else {
		return &tree, nil
	}
}

func (s *Sequence) ToString() string {
	text := ""
	first := true
	for _, e := range s.Elements {
		if first {
			first = false
		} else {
			text += " "
		}
		text += e.ToString()
	}
	return text
}

// --- Optional (?) --- //

func (o *Optional) Terminal() bool {
	return false
}

func (o *Optional) First(g *Grammar, cd CycleDetector) (map[string]bool, error) {
	return o.Sentence.First(g, cd)
}

func (o *Optional) Follow(g *Grammar, production string, cd CycleDetector) (map[string]bool, bool, error) {
	f, _, err := o.Sentence.Follow(g, production, cd)
	if err != nil {
		return nil, false, err
	}
	return f, true, nil
}

func (o *Optional) MatchEmpty(*Grammar) bool {
	return true
}

func (o *Optional) Recognise(g *Grammar, production LanguageElement, tokens *lexer.TokenSeq, cd CycleDetector) (*SyntaxTree, error) {
	token, err := peek(tokens)
	if err != nil {
		return nil, err
	}
	first, err := o.Sentence.First(g, cd)
	if err != nil {
		return nil, err
	}
	if _, ok := first[token.Type]; ok {
		return o.Sentence.Recognise(g, production, tokens, cd)
	} else if !o.Sentence.MatchEmpty(g) {
		return nil, fmt.Errorf("token %q cannot start %q", token.Type, o.Sentence.ToString())
	}
	return nil, nil
}

func (o *Optional) ToString() string {
	return "(" + o.Sentence.ToString() + ")?"
}

// --- Zero or more (*) --- //

func (o *ZeroOrMore) Terminal() bool {
	return false
}

func (o *ZeroOrMore) First(g *Grammar, cd CycleDetector) (map[string]bool, error) {
	return o.Sentence.First(g, cd)
}

func (o *ZeroOrMore) Follow(g *Grammar, production string, cd CycleDetector) (map[string]bool, bool, error) {
	return o.Sentence.Follow(g, production, cd)
}

func (o *ZeroOrMore) MatchEmpty(*Grammar) bool {
	return true
}

func (o *ZeroOrMore) Recognise(g *Grammar, production LanguageElement, tokens *lexer.TokenSeq, cd CycleDetector) (*SyntaxTree, error) {
	first, err := o.Sentence.First(g, cd)
	if err != nil {
		return nil, err
	}
	tree := SyntaxTree{production, nil}
	matchedOnce := false
	for {
		token, err := peek(tokens)
		if err != nil {
			return nil, err
		}
		if _, ok := first[token.Type]; ok {
			matchedOnce = true
			child, err := o.Sentence.Recognise(g, production, tokens, cd)
			if err != nil {
				return nil, err
			}
			tree.Children = append(tree.Children, child)

		} else if !matchedOnce && !o.Sentence.MatchEmpty(g) {
			return nil, fmt.Errorf("token %q cannot start %q", token.Type, o.Sentence.ToString())

		} else {
			break
		}
	}
	return &tree, nil
}

func (o *ZeroOrMore) ToString() string {
	return "(" + o.Sentence.ToString() + ")*"
}

// --- One or more (*) --- //

func (o *OneOrMore) Terminal() bool {
	return false
}

func (o *OneOrMore) First(g *Grammar, cd CycleDetector) (map[string]bool, error) {
	return o.Sentence.First(g, cd)
}

func (o *OneOrMore) Follow(g *Grammar, production string, cd CycleDetector) (map[string]bool, bool, error) {
	return o.Sentence.Follow(g, production, cd)
}

func (o *OneOrMore) MatchEmpty(g *Grammar) bool {
	return o.Sentence.MatchEmpty(g)
}

func (o *OneOrMore) Recognise(g *Grammar, production LanguageElement, tokens *lexer.TokenSeq, cd CycleDetector) (*SyntaxTree, error) {
	first, err := o.Sentence.First(g, cd)
	if err != nil {
		return nil, err
	}
	tree := SyntaxTree{production, nil}
	matchedOnce := false
	for {
		token, err := peek(tokens)
		if err != nil {
			return nil, err
		}
		if _, ok := first[token.Type]; ok {
			matchedOnce = true
			child, err := o.Sentence.Recognise(g, production, tokens, cd)
			if err != nil {
				return nil, err
			}
			tree.Children = append(tree.Children, child)

		} else if !matchedOnce && !o.Sentence.MatchEmpty(g) {
			return nil, fmt.Errorf("token %q cannot start %q", token.Type, o.Sentence.ToString())

		} else {
			break
		}
	}
	return &tree, nil
}

func (o *OneOrMore) ToString() string {
	return "(" + o.Sentence.ToString() + ")+"
}

// --- Repeat match {m,n} --- //

func (r *Repeat) Terminal() bool {
	return false
}

func (r *Repeat) First(g *Grammar, cd CycleDetector) (map[string]bool, error) {
	return r.Sentence.First(g, cd)
}

func (r *Repeat) Follow(g *Grammar, production string, cd CycleDetector) (map[string]bool, bool, error) {
	return r.Sentence.Follow(g, production, cd)
}

func (r *Repeat) MatchEmpty(g *Grammar) bool {
	return r.Min == 0 || r.Sentence.MatchEmpty(g)
}

func (r *Repeat) Recognise(g *Grammar, production LanguageElement, tokens *lexer.TokenSeq, cd CycleDetector) (*SyntaxTree, error) {
	first, err := r.Sentence.First(g, cd)
	if err != nil {
		return nil, err
	}
	tree := SyntaxTree{production, nil}

	for matched := 0; matched < r.Max; matched++ {
		token, err := peek(tokens)
		if err != nil {
			return nil, err
		}
		if _, ok := first[token.Type]; ok {
			child, err := r.Sentence.Recognise(g, production, tokens, cd)
			if err != nil {
				return nil, err
			}
			tree.Children = append(tree.Children, child)

		} else if matched < r.Min && !r.Sentence.MatchEmpty(g) {
			return nil, fmt.Errorf("token %q cannot start %q", token.Type, r.Sentence.ToString())

		} else if matched >= r.Min {
			break
		}
	}
	return &tree, nil
}

func (r *Repeat) ToString() string {
	return "(" + r.Sentence.ToString() + "){" + strconv.Itoa(r.Min) + "," + strconv.Itoa(r.Max) + "}"
}

func New(name string, l *lexer.Lexer, productions []*Production) *Grammar {
	prodByName := map[string]*Production{}
	for _, p := range productions {
		prodByName[p.Name] = p
	}

	//tokenTypeByName := map[string]*lexer.TokenType{}
	//for _, t := range tokenTypes {
	//	tokenTypeByName[t.Id] = t
	//}

	return &Grammar{
		Id:    name,
		Lexer: l,
		//TokenTypes:      tokenTypes,
		//TokenTypeByName: tokenTypeByName,
		Productions: productions,
		ProdByName:  prodByName,
	}

	/*
		Algorithm to compute FIRST, FOLLOW, and nullable.
		Initialize FIRST and FOLLOW to all empty sets, and nullable to all false.
		for each terminal symbol Z
			FIRST[Z] ← {Z}
		repeat
			for each production X → Y1Y2 · · · Y k
				if Y1 . . . Y k are all nullable (or if k = 0)
				then nullable[X] ← true
				for each i from 1 to k, each j from i + 1 to k
					if Y1 · · · Y i−1 are all nullable (or if i = 1)
					then FIRST[X] ← FIRST[X] ∪ FIRST[Y i ]

					if Y i+1 · · · Y k are all nullable (or if i = k)
					then FOLLOW[Y i ] ← FOLLOW[Y i ] ∪ FOLLOW[X]

					if Y i+1 · · · Y j −1 are all nullable (or if i + 1 = j )
					then FOLLOW[Y i ] ← FOLLOW[Y i ] ∪ FIRST[Y j ]

		until FIRST, FOLLOW, and nullable did not change in this iteration.
	*/
}

func (g *Grammar) Parse(input io.Reader) (*SyntaxTree, error) {
	//l := lexer.New(g.TokenTypes...)
	tokenSeq := g.Lexer.Lex(input)
	defer tokenSeq.Stop()
	prod := g.Productions[0]
	cd := &CycleDetectorSet{make(map[LanguageElement]bool)}
	return prod.Recognise(g, prod, tokenSeq, cd)
}

func (g *Grammar) ParseFile(filename string) (*SyntaxTree, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	return g.Parse(file)
}

func (g *Grammar) ParseText(input string) (*SyntaxTree, error) {
	return g.Parse(strings.NewReader(input))
}
