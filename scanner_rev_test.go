package tok

import "testing"

// ---------------------------------------------------------------------- string
func TestRevIf(t *testing.T) {
	cases := []stringCase{
		{"* from", "from", headTail{"* ", "from"}},
	}
	for i, c := range cases {
		sca := NewRevScanner(c.inp)
		sca.RevIf(c.str)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestRefIfAny(t *testing.T) {
	sca := NewRevScanner("* From")
	sca.RevIfAny([]string{"from", "From", "FROM"})
	exp := headTail{"* ", "From"}
	if e := exp.check(1, sca); e != nil {
		t.Errorf("%v", e)
	}
}

func TestRevTo(t *testing.T) {
	cases := []stringCase{
		{"select * from events", "from", headTail{"select * from", " events"}},
	}
	for i, c := range cases {
		sca := NewRevScanner(c.inp)
		sca.RevTo(c.str)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

// ------------------------------------------------------------------------ fold
func TestRevIfFold(t *testing.T) {
	cases := []foldCase{
		{"* from", "from", headTail{"* ", "from"}},
		{"* from", "From", headTail{"* ", "from"}},
		{"* from", "FROM", headTail{"* ", "from"}},
	}
	for i, c := range cases {
		sca := NewRevScanner(c.inp)
		sca.RevIfFold(c.str)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestRevToFold(t *testing.T) {
	cases := []foldCase{
		{"select * from events", "from", headTail{"select * from", " events"}},
		{"select * from events", "From", headTail{"select * from", " events"}},
		{"select * from events", "FROM", headTail{"select * from", " events"}},
		{"select * from events", "SELECT", headTail{"select", " * from events"}},
	}
	for i, c := range cases {
		sca := NewRevScanner(c.inp)
		sca.RevToFold(c.str)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

// ------------------------------------------------------------------------ rune
func TestRevIfRune(t *testing.T) {
	cases := []runeCase{
		{"xyz", 'z', headTail{"xy", "z"}},
	}
	for i, c := range cases {
		sca := NewRevScanner(c.inp)
		sca.RevIfRune(c.r)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestRevToRune(t *testing.T) {
	cases := []runeCase{
		{"i am", ' ', headTail{"i ", "am"}},
		{"?...", '.', headTail{"?...", ""}},
		{"1. end", '1', headTail{"1", ". end"}},
	}
	for i, c := range cases {
		sca := NewRevScanner(c.inp)
		sca.RevToRune(c.r)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestRevWhileRune(t *testing.T) {
	cases := []runeCase{
		{"123-", '-', headTail{"123", "-"}},
		{"a....", '.', headTail{"a", "...."}},
		{"   ", ' ', headTail{"", "   "}},
	}
	for i, c := range cases {
		sca := NewRevScanner(c.inp)
		sca.RevWhileRune(c.r)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

// --------------------------------------------------------------------- anyrune
func TestRevIfAnyRune(t *testing.T) {
	cases := []anyRuneCase{
		{"123,", ":,", headTail{"123", ","}},
		{"123:", ":,", headTail{"123", ":"}},
	}
	for i, c := range cases {
		sca := NewRevScanner(c.inp)
		sca.RevIfAnyRune(c.str)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestRevToAnyRune(t *testing.T) {
	cases := []anyRuneCase{
		{"123,56", " ,", headTail{"123,", "56"}},
		{"123 56", " ,", headTail{"123 ", "56"}},
		{"+123", "+-", headTail{"+", "123"}},
		{"-123", "+-", headTail{"-", "123"}},
	}
	for i, c := range cases {
		sca := NewRevScanner(c.inp)
		sca.RevToAnyRune(c.str)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestRevWhileAnyRune(t *testing.T) {
	cases := []anyRuneCase{
		{"123-", "-", headTail{"123", "-"}},
		{"-256", "1234567890", headTail{"-", "256"}},
		{"12344512", "1234567890", headTail{"", "12344512"}},
	}
	for i, c := range cases {
		sca := NewRevScanner(c.inp)
		sca.RevWhileAnyRune(c.str)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

// --------------------------------------------------------------------- between
func TestRevIfBetween(t *testing.T) {
	cases := []betweenCase{}
	for i, c := range cases {
		sca := NewRevScanner(c.inp)
		sca.IfBetween(c.min, c.max)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestRevToBetween(t *testing.T) {
	cases := []betweenCase{}
	for i, c := range cases {
		sca := NewRevScanner(c.inp)
		sca.ToBetween(c.min, c.max)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestRevWhileBetween(t *testing.T) {
	cases := []betweenCase{}
	for i, c := range cases {
		sca := NewRevScanner(c.inp)
		sca.WhileBetween(c.min, c.max)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

// ---------------------------------------------------------------------- match
func TestRevIfMatch(t *testing.T) {

}

func TestRevToMatch(t *testing.T) {

}

func TestRevWhileMatch(t *testing.T) {

}
