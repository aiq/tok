package tok

import (
	"testing"
)

func S(info string, from int, to int) Segment {
	return Segment{info, MakeToken(Marker(from), Marker(to))}
}

func N(info string, from int, to int, nodes ...*Node) *Node {
	return &Node{S(info, from, to), nodes}
}

func TestGraphAppend(t *testing.T) {
	cases := []struct {
		base *Graph
		s    Segment
		g    *Graph
		ok   bool
	}{
		{NewGraph("root"), S("val", 10, 15), &Graph{N("root", 10, 15, N("val", 10, 15))}, true},
		{
			&Graph{N("root", 10, 15, N("key", 10, 15))}, S("val", 17, 25),
			&Graph{N("root", 10, 25, N("key", 10, 15), N("val", 17, 25))}, true,
		},
		{
			&Graph{N("root", 10, 25, N("key", 10, 15), N("val", 17, 25))}, S("obj", 8, 28),
			&Graph{N("root", 8, 28, N("obj", 8, 28, N("key", 10, 15), N("val", 17, 25)))}, true,
		},
		{ // clash
			&Graph{N("root", 10, 25, N("key", 10, 15), N("val", 17, 25))}, S("obj", 8, 20),
			&Graph{N("root", 10, 25, N("key", 10, 15), N("val", 17, 25))}, false,
		},
	}
	for i, c := range cases {
		g, ok := c.base.Append(c.s)
		if ok != c.ok {
			t.Errorf("%d unexpected ok value: %v", i, ok)
		}
		if !g.Equal(c.g) {
			t.Errorf("%d unexpected graph: %v", i, g)
		}
	}
}

func TestGraphLeafs(t *testing.T) {
	cases := []struct {
		g   *Graph
		exp []Segment
	}{
		{
			&Graph{N("root", 8, 28, N("obj", 8, 28, N("key", 10, 15), N("val", 17, 25)))},
			[]Segment{S("key", 10, 15), S("val", 17, 25)},
		},
	}
	for i, c := range cases {
		leafs := c.g.Leafs()
		if len(leafs) != len(c.exp) {
			t.Errorf("%d unexpected number of leafs: %d != %d", i, len(leafs), len(c.exp))
		} else {
			for j, l := range leafs {
				if l != c.exp[j] {
					t.Errorf("%d unexpected value at %d: %v", i, j, l)
				}
			}
		}
	}
}

func TestBuildGraph(t *testing.T) {
	cases := []struct {
		inp []Segment
		exp *Graph
	}{
		{
			[]Segment{S("text", 0, 20), S("obj", 2, 18), S("id", 3, 8), S("val", 10, 16)},
			&Graph{N("root", 0, 20, N("text", 0, 20, N("obj", 2, 18, N("id", 3, 8), N("val", 10, 16))))},
		},
	}
	for i, c := range cases {
		g := BuildGraph("root", c.inp)
		if !g.Equal(c.exp) {
			t.Errorf("%d unexpected graph: %v", i, g)
		}
	}
}
