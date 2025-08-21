package grammar

import (
	"strconv"
)

const (
	Drop TreeRetention = iota
	Retain
	Promoted
	Promote1
	Promote2
)

// TreeRetention defines whether a language element is kept in the SyntaxTree
// and where it is positioned. It has 3 values:
//  1. Retain: kept in the syntax tree at the default position;
//  2. Drop: not kept in the syntax tree;
//  3. Promote: the language element is promoted to the parent position in the tree.
type TreeRetention int

type SyntaxTree struct {
	Node     LanguageElement
	Children []*SyntaxTree
}

func (tree *SyntaxTree) ToGraphViz(title string) string {
	spec := "digraph G {\n"
	if len(title) > 0 {
		spec += "\tlabel=\"" + title + "\"\n"
	}
	spec += tree.graphVizNode("0")
	spec += "}"
	return spec
}

func (tree *SyntaxTree) graphVizNode(position string) string {
	spec := ""
	for i, c := range tree.Children {
		pos := position + strconv.Itoa(i)
		spec += "\t\"" + tree.Node.ToString() + " [" + position + "]" +
			"\" -> \"" + c.Node.ToString() + " [" + pos + "]" + "\"\n"
		spec += c.graphVizNode(pos)
	}
	return spec
}
