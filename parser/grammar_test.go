package grammar

import (
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
				Sentence: &OneOrMore{&ProductionRef{"Stmt"}},
			},
			{
				Name: "Stmt",
				Sentence: &Choice{
					Alternates: []Sentence{
						&Sequence{Elements: []Sentence{
							&TokenRef{"LET"},
							&TokenRef{"ID"},
							&TokenRef{":="},
							&ProductionRef{"Expr"},
							&TokenRef{";"},
						}},
						&Sequence{Elements: []Sentence{
							&TokenRef{"ID"},
							&TokenRef{"="},
							&ProductionRef{"Expr"},
							&TokenRef{";"},
						}},
					},
				},
			},
			{
				Name: "Expr",
				Sentence: &Sequence{Elements: []Sentence{
					&ProductionRef{"Term"},
					&Optional{&Sequence{
						Elements: []Sentence{
							&TokenRef{"ADD"},
							&ProductionRef{"Expr"},
						},
					}},
				}},
			},
			{
				Name:     "Term",
				Sentence: &OneOrMore{&ProductionRef{"Factor"}},
			},
			{
				Name: "Factor",
				Sentence: &Sequence{Elements: []Sentence{
					&ProductionRef{"Base"},
					&Optional{&Sequence{
						Elements: []Sentence{
							&TokenRef{"MUL"},
							&ProductionRef{"Expr"},
						},
					}},
				}},
			},
			{
				Name: "Base",
				Sentence: &Choice{Alternates: []Sentence{
					&Sequence{Elements: []Sentence{
						&TokenRef{"("},
						&ProductionRef{"Expr"},
						&TokenRef{")"},
					}},
					&TokenRef{"INT"},
					&TokenRef{"ID"},
				},
				},
			},
		},
	)
}
