package grammar

import (
	"strconv"
)

type SyntaxTree struct {
	Node     LanguageElement
	Children []*SyntaxTree
}

func (tree *SyntaxTree) ToGraphViz(title string) string {
	spec := "digraph G {\n"
	if len(title) > 0 {
		spec += "\tlabel=\"" + title + "\"\n"
	}
	spec += tree.graphVizNode(1, "0")
	spec += "}"
	return spec
}

func (tree *SyntaxTree) graphVizNode(level int, position string) string {
	spec := ""
	for i, c := range tree.Children {
		pos := position + "/" + strconv.Itoa(i)
		spec += "\t\"" + tree.Node.ToString() + " [" + strconv.Itoa(level-1) + ":" + position + "]" +
			"\" -> \"" + c.Node.ToString() + " [" + strconv.Itoa(level) + ":" + pos + "]" + "\"\n"
		spec += c.graphVizNode(level+1, pos)
	}
	return spec
}
