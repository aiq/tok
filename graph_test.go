package tok

import (
	"testing"
)

func TestGraphAppend(t *testing.T) {
	V := func(class string, from int, to int) Value {
		return Value{class, MakeToken(Marker(from), Marker(to))}
	}
	N := func(class string, from int, to int, nodes ...*Node) *Node {
		return &Node{V(class, from, to), nodes}
	}
	cases := []struct {
		base *Graph
		v    Value
		g    *Graph
		ok   bool
	}{
		{NewGraph("root"), V("val", 10, 15), &Graph{N("root", 10, 15, N("val", 10, 15))}, true},
		{
			&Graph{N("root", 10, 15, N("key", 10, 15))}, V("val", 17, 25),
			&Graph{N("root", 10, 25, N("key", 10, 15), N("val", 17, 25))}, true,
		},
		{
			&Graph{N("root", 10, 25, N("key", 10, 15), N("val", 17, 25))}, V("obj", 8, 28),
			&Graph{N("root", 8, 28, N("obj", 8, 28, N("key", 10, 15), N("val", 17, 25)))}, true,
		},
		{ // clash
			&Graph{N("root", 10, 25, N("key", 10, 15), N("val", 17, 25))}, V("obj", 8, 20),
			&Graph{N("root", 10, 25, N("key", 10, 15), N("val", 17, 25))}, false,
		},
	}
	for i, c := range cases {
		g, ok := c.base.Append(c.v)
		if ok != c.ok {
			t.Errorf("%d unexpected ok value: %v", i, ok)
		}
		if !g.Equal(c.g) {
			t.Errorf("%d unexpected graph: %v", i, g)
		}
	}
}

func TestBuildGraph(t *testing.T) {
	V := func(class string, from int, to int) Value {
		return Value{class, MakeToken(Marker(from), Marker(to))}
	}
	N := func(class string, from int, to int, nodes ...*Node) *Node {
		return &Node{V(class, from, to), nodes}
	}
	cases := []struct {
		inp []Value
		exp *Graph
	}{
		{
			[]Value{V("text", 0, 20), V("obj", 2, 18), V("id", 3, 8), V("val", 10, 16)},
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
