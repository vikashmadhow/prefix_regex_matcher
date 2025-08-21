package grammar

import (
	"fmt"
	"testing"

	"github.com/vikashmadhow/prefix_regex_matcher/lexer"
)

// Seq[characters] -> Lexer -> Seq[Token] -> Modulator... -> Seq[Token] -> Syntax Analyser -> ST -> Semantic Processors... -> AT -> Translators... -> Translation

func TestSimpleParser(t *testing.T) {
	g := testGrammar()
	tree, err := g.ParseTextFromStart(
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
							&TokenRef{"LET", Promote1},
							&TokenRef{"ID", Retain},
							&TokenRef{":=", Drop},
							&ProductionRef{"Expr", Retain},
							&TokenRef{";", Drop},
						}},
						&Sequence{Elements: []Sentence{
							&TokenRef{"ID", Retain},
							&TokenRef{"=", Promote1},
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
							&TokenRef{"ADD", Promote2},
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
							&TokenRef{"MUL", Promote2},
							&ProductionRef{"Expr", Retain},
						}, TreeRetention: Retain,
					}, Retain},
				}},
			},
			{
				Name: "Base",
				Sentence: &Choice{Alternates: []Sentence{
					&Sequence{Elements: []Sentence{
						&TokenRef{"(", Promote1},
						&ProductionRef{"Expr", Retain},
						&TokenRef{")", Drop},
					}},
					&TokenRef{"INT", Retain},
					&TokenRef{"ID", Retain},
				}},
			},
		},
	)
}
