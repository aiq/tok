package tok

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//------------------------------------------------------------------------------

// Node
type Node struct {
	Segment
	Nodes []*Node
}

// Equal
func (n *Node) Equal(oth *Node) bool {
	if n.Segment != oth.Segment {
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

// Graph allows to arrange the picked values hierarchically.
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

// Append appends a Segment to a graph.
func (g *Graph) Append(seg Segment) (*Graph, bool) {
	n := &Node{Segment: seg}
	bkp := g.Root.Token
	if !g.Root.Covers(n.Token) {
		if len(g.Root.Nodes) > 0 {
			g.Root.Token = g.Root.Token.Merge(n.Segment.Token)
		} else {
			g.Root.Token = n.Segment.Token
		}
	}
	ok := rankNodes(g.Root, n)
	if !ok {
		g.Root.Token = bkp
	}
	return g, ok
}

// Equal
func (g *Graph) Equal(oth *Graph) bool {
	return g.Root.Equal(oth.Root)
}

func makeStackLines(b *strings.Builder, prefix string, i int, n *Node) {
	prefix = fmt.Sprintf("%s%d.%s", prefix, i, n.String())
	b.WriteString(prefix)
	b.WriteRune(' ')
	b.WriteString(strconv.Itoa(n.Len()))
	b.WriteRune('\n')
	for j, sub := range n.Nodes {
		makeStackLines(b, prefix+";", j+1, sub)
	}
}

// FlameStack
func (g *Graph) FlameStack() string {
	b := &strings.Builder{}
	makeStackLines(b, "", 1, g.Root)
	return b.String()
}

func (n *Node) appendLeafs(leafs *[]Segment) {
	if len(n.Nodes) == 0 {
		*leafs = append(*leafs, n.Segment)
	} else {
		for _, sub := range n.Nodes {
			sub.appendLeafs(leafs)
		}
	}
}

// Leafs returns the values of a graph that hang on the leafs.
func (g *Graph) Leafs() []Segment {
	leafs := []Segment{}
	g.Root.appendLeafs(&leafs)
	return leafs
}

// BuildGraph creates a graph with the segments.
// The name argument will be used as the info of the root node.
func BuildGraph(name string, segments []Segment) *Graph {
	g := NewGraph(name)
	for _, v := range segments {
		g.Append(v)
	}
	return g
}

// NewGraph creates an empty graph with name as the info of the root node.
func NewGraph(name string) *Graph {
	g := &Graph{&Node{}}
	g.Root.Segment.Info = name
	return g
}
