package grammar

import (
	"fmt"
	"github.com/vikashmadhow/prefix_regex_matcher/lexer"
	"testing"
)

func TestSimpleParser(t *testing.T) {
	g := testGrammar()
	tree, err := g.ParseText(
		`let x := 1000;
		 let y := 2000;
	     x = x + 5 * (4 + y / 2);
		 y = y + x;`,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf(tree.ToGraphViz("First program"))
}

func testGrammar() *Grammar {
	lex := lexer.New(
		lexer.NewTokenType("LET", "let"),
		lexer.NewTokenType("INT", "\\d+"),
		lexer.NewTokenType("ID", "[_a-zA-Z][_a-zA-Z0-9]*"),
		lexer.SimpleTokenType("="),
		lexer.SimpleTokenType(":="),
		lexer.SimpleTokenType(">"),
		lexer.SimpleTokenType(">="),
		lexer.SimpleTokenType("<"),
		lexer.SimpleTokenType("<="),
		lexer.SimpleTokenType("("),
		lexer.SimpleTokenType(")"),
		lexer.SimpleTokenType(";"),
		lexer.NewTokenType("ADD", "\\+|-"),
		lexer.NewTokenType("MUL", "\\*|/"),
		lexer.NewTokenType("SPC", "\\s+"),
	)
	lex.Modulator(lexer.Ignore("SPC"))
	return New(
		"test_language",
		lex,
		[]*Production{
			{
				Name:     "Program",
				Sentence: &OneOrMore{&ProductionRef{"Stmt", Retain}, Retain},
			},
			{
				Name: "Stmt",
				Sentence: &Choice{
					Alternates: []Sentence{
						&Sequence{Elements: []Sentence{
							&TokenRef{"LET", Promote},
							&TokenRef{"ID", Retain},
							&TokenRef{":=", Drop},
							&ProductionRef{"Expr", Retain},
							&TokenRef{";", Drop},
						}},
						&Sequence{Elements: []Sentence{
							&TokenRef{"ID", Retain},
							&TokenRef{"=", Promote},
							&ProductionRef{"Expr", Retain},
							&TokenRef{";", Drop},
						}},
					},
				},
			},
			{
				Name: "Expr",
				Sentence: &Sequence{Elements: []Sentence{
					&ProductionRef{"Term", Retain},
					&Optional{&Sequence{
						Elements: []Sentence{
							&TokenRef{"ADD", Retain},
							&ProductionRef{"Expr", Retain},
						}, TreeRetention: Retain,
					}, Retain},
				}},
			},
			{
				Name:     "Term",
				Sentence: &OneOrMore{&ProductionRef{"Factor", Retain}, Retain},
			},
			{
				Name: "Factor",
				Sentence: &Sequence{Elements: []Sentence{
					&ProductionRef{"Base", Retain},
					&Optional{&Sequence{
						Elements: []Sentence{
							&TokenRef{"MUL", Retain},
							&ProductionRef{"Expr", Retain},
						}, TreeRetention: Retain,
					}, Retain},
				}},
			},
			{
				Name: "Base",
				Sentence: &Choice{Alternates: []Sentence{
					&Sequence{Elements: []Sentence{
						&TokenRef{"(", Retain},
						&ProductionRef{"Expr", Retain},
						&TokenRef{")", Retain},
					}},
					&TokenRef{"INT", Retain},
					&TokenRef{"ID", Retain},
				}},
			},
		},
	)
}
