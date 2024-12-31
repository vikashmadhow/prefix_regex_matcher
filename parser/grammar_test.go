package grammar

import (
	"github.com/vikashmadhow/prefix_regex_matcher/lexer"
	"testing"
)

func TestSimpleParser(t *testing.T) {

}

func testGrammar() *Grammar {
	return NewGrammar(
		"test_language",
		[]*lexer.TokenType{
			lexer.NewTokenType("LET", "let"),
			lexer.NewTokenType("INT", "\\d+"),
			lexer.NewTokenType("ID", "[_a-zA-Z][_a-zA-Z0-9]*"),
			lexer.NewTokenType("EQ", "="),
			lexer.NewTokenType("GT", ">"),
			lexer.NewTokenType("GTE", ">="),
			lexer.NewTokenType("LT", "<"),
			lexer.NewTokenType("LTE", "<="),
			lexer.NewTokenType("ADD", "\\+|-"),
			lexer.NewTokenType("MUL", "\\*|/"),
			lexer.NewTokenType("SPC", "\\s+"),
			lexer.SimpleTokenType("("),
			lexer.SimpleTokenType(")"),
		},

		[]*Production{{
			Name: "S",
			Sentence: &Choice{Alternates: []Sentence{
				&Sequence{Elements: []Sentence{
					&TokenRef{"INT"},
					&ProductionRef{"F"},
					&ZeroOrMore{&Sequence{Elements: []Sentence{&TokenRef{"X"}}}},
				}}},
			}},
		},
	)
}
