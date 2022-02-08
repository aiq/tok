package tok

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//------------------------------------------------------------------------------

type Node struct {
	Value
	Nodes []*Node
}

func (n *Node) Equal(oth *Node) bool {
	if n.Value != oth.Value {
		return false
	}
	if len(n.Nodes) != len(oth.Nodes) {
		return false
	}
	for i, sub := range n.Nodes {
		if !sub.Equal(oth.Nodes[i]) {
			return false
		}
	}
	return true
}

type nodeByPos []*Node

func (s nodeByPos) Len() int {
	return len(s)
}

func (s nodeByPos) Less(i, j int) bool {
	return s[i].Before(s[j].Token)
}

func (s nodeByPos) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//------------------------------------------------------------------------------

type Graph struct {
	Root *Node
}

func rankNodes(a *Node, b *Node) bool {
	if a.Clashes(b.Token) {
		return false
	}
	nodes := []*Node{}
	var ok bool
	for i := 0; i < len(a.Nodes); i++ {
		c := a.Nodes[i]
		if b.Covers(c.Token) {
			ok = rankNodes(b, c)
			if !ok {
				return false
			}
		} else if b.Clashes(c.Token) {
			return false
		} else {
			if c.Covers(b.Token) {
				return rankNodes(c, b)
			}
			nodes = append(nodes, c)
		}
	}
	a.Nodes = nodes
	a.Nodes = append(a.Nodes, b)

	sort.Stable(nodeByPos(a.Nodes))
	return true
}

func (g *Graph) Append(v Value) (*Graph, bool) {
	n := &Node{Value: v}
	bkp := g.Root.Token
	if !g.Root.Covers(n.Token) {
		if len(g.Root.Nodes) > 0 {
			g.Root.Token = g.Root.Token.Merge(n.Value.Token)
		} else {
			g.Root.Token = n.Value.Token
		}
	}
	ok := rankNodes(g.Root, n)
	if !ok {
		g.Root.Token = bkp
	}
	return g, ok
}

type jsonNode struct {
	Name     string      `json:"name"`
	Value    int         `json:"value"`
	Children []*jsonNode `json:"children,omitempty"`
}

func makeJSONNode(n *Node) *jsonNode {
	jn := &jsonNode{
		Name:  n.Info,
		Value: n.Len(),
	}
	for _, sub := range n.Nodes {
		jn.Children = append(jn.Children, makeJSONNode(sub))
	}
	return jn
}

func makeStackLines(b *strings.Builder, prefix string, n *Node) {
	prefix += n.Info
	b.WriteString(prefix)
	b.WriteRune(' ')
	b.WriteString(strconv.Itoa(n.Len()))
	b.WriteRune('\n')
	for _, sub := range n.Nodes {
		makeStackLines(b, prefix+";", sub)
	}
}

func (g *Graph) Format(format string) string {
	if format == "json" {
		jn := makeJSONNode(g.Root)
		out, _ := json.Marshal(jn)
		return string(out)
	} else if format == "json-i" {
		jn := makeJSONNode(g.Root)
		out, _ := json.MarshalIndent(jn, "", "\t")
		return string(out)
	} else if format == "stack" {
		b := &strings.Builder{}
		makeStackLines(b, "", g.Root)
		return b.String()
	}

	return fmt.Sprintf("unknown graph format %q", format)
}

func (g *Graph) Equal(oth *Graph) bool {
	return g.Root.Equal(oth.Root)
}

func BuildGraph(name string, values []Value) *Graph {
	g := NewGraph(name)
	for _, v := range values {
		g.Append(v)
	}
	return g
}

func NewGraph(name string) *Graph {
	g := &Graph{&Node{}}
	g.Root.Value.Info = name
	return g
}
