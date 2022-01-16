package tok

import (
	"fmt"
	"testing"
)

type headTail struct {
	head string
	tail string
}

func (x headTail) check(i int, s *Scanner) error {
	if x.head != s.Head() {
		return fmt.Errorf("test %d head: expected %q got %q", i, x.head, s.Head())
	}
	if x.tail != s.Tail() {
		return fmt.Errorf("test %d tail: expected %q got %q", i, x.tail, s.Tail())
	}
	return nil
}

// ---------------------------------------------------------------------- string

// ------------------------------------------------------------------------ fold
func TestIfFold(t *testing.T) {
	cases := []struct {
		inp string
		str string
		exp headTail
	}{
		{"select *", "select", headTail{"select", " *"}},
		{"select *", "Select", headTail{"select", " *"}},
		{"select *", "SELECT", headTail{"select", " *"}},
	}

	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.IfFold(c.str)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestToFold(t *testing.T) {
	cases := []struct {
		inp string
		str string
		exp headTail
	}{
		{"select * from events", "from", headTail{"select * ", "from events"}},
		{"select * from events", "From", headTail{"select * ", "from events"}},
		{"select * from events", "FROM", headTail{"select * ", "from events"}},
		{"select * from", "FROM", headTail{"select * ", "from"}},
	}

	for i, c := range cases {
		sca := NewScanner(c.inp)

		sca.ToFold(c.str)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

// ------------------------------------------------------------------------ rune
func TestIfRune(t *testing.T) {
	sca := NewScanner("a世z")

	cases := []struct {
		r   rune
		res bool
		exp headTail
	}{
		{'x', false, headTail{"", "a世z"}},
		{'a', true, headTail{"a", "世z"}},
		{'世', true, headTail{"a世", "z"}},
		{'>', false, headTail{"a世", "z"}},
		{'z', true, headTail{"a世z", ""}},
		{'-', false, headTail{"a世z", ""}},
	}

	for i, c := range cases {
		res := sca.IfRune(c.r)
		if c.res != res {
			t.Errorf("test %d res: expected %v got %v", i, c.res, res)
		}
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestWhileRune(t *testing.T) {
	cases := []struct {
		inp string
		r   rune
		exp headTail
	}{
		{"-123", '-', headTail{"-", "123"}},
		{"....a", '.', headTail{"....", "a"}},
		{"   ", ' ', headTail{"   ", ""}},
	}

	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.WhileRune(c.r)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

// --------------------------------------------------------------------- anyrune

// --------------------------------------------------------------------- between

// ---------------------------------------------------------------------- match
