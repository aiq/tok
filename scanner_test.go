package tok

import (
	"fmt"
	"testing"
	"unicode"
)

func TestScanString(t *testing.T) {
	checkPart := func(p, expP string, ok, expOk bool) {
		if p != expP || ok != expOk {
			t.Errorf("unexpected ScanString result: %s / %v", p, ok)
		}
	}
	sca := NewScanner("Hi, 世界! I saw a ☃ in w2s34")
	part, ok := sca.ScanString(5)
	checkPart(part, "Hi, 世", ok, true)
	part, ok = sca.ScanString(3)
	checkPart(part, "界! ", ok, true)
	part, ok = sca.ScanString(8)
	checkPart(part, "I saw a ", ok, true)
	part, ok = sca.ScanString(1)
	checkPart(part, "☃", ok, true)
	part, ok = sca.ScanString(20)
	checkPart(part, " in w2s34", ok, false)
}

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
type stringCase struct {
	inp string
	str string
	exp headTail
}

func TestIf(t *testing.T) {
	cases := []stringCase{
		{"select *", "select", headTail{"select", " *"}},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.If(c.str)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestIfAny(t *testing.T) {
	sca := NewScanner("SELECT *")
	sca.IfAny("select", "Select", "SELECT")
	exp := headTail{"SELECT", " *"}
	if e := exp.check(1, sca); e != nil {
		t.Errorf("%v", e)
	}
}

func TestTo(t *testing.T) {
	cases := []stringCase{
		{"select * from events", "from", headTail{"select * ", "from events"}},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.To(c.str)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

// ------------------------------------------------------------------------ fold
type foldCase stringCase

func TestIfFold(t *testing.T) {
	cases := []foldCase{
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
	cases := []foldCase{
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

type runeCase struct {
	inp string
	r   rune
	exp headTail
}

func TestToRune(t *testing.T) {
	cases := []runeCase{
		{"i am", ' ', headTail{"i", " am"}},
		{"...?", '?', headTail{"...", "?"}},
		{"1. end", '1', headTail{"", "1. end"}},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.ToRune(c.r)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestWhileRune(t *testing.T) {
	cases := []runeCase{
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
type anyRuneCase stringCase

func TestIfAnyRune(t *testing.T) {
	cases := []anyRuneCase{
		{"-123", "-+", headTail{"-", "123"}},
		{"+123", "-+", headTail{"+", "123"}},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.IfAnyRune(c.str)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestToAnyRune(t *testing.T) {
	cases := []anyRuneCase{
		{"123,56", " ,", headTail{"123", ",56"}},
		{"123 56", " ,", headTail{"123", " 56"}},
		{"12356;", ";.", headTail{"12356", ";"}},
		{"12356.", ";.", headTail{"12356", "."}},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.ToAnyRune(c.str)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestWhileAnyRune(t *testing.T) {
	cases := []anyRuneCase{
		{"-123", "-", headTail{"-", "123"}},
		{"256,0", "1234567890", headTail{"256", ",0"}},
		{"12344512", "1234567890", headTail{"", "12344512"}},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.WhileAnyRune(c.str)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

// --------------------------------------------------------------------- between
type betweenCase struct {
	inp string
	min rune
	max rune
	exp headTail
}

func TestIfBetween(t *testing.T) {
	cases := []betweenCase{
		{"abba", 'a', 'z', headTail{"a", "bba"}},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.IfBetween(c.min, c.max)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestToBetween(t *testing.T) {
	cases := []betweenCase{
		{"pinball2000", '0', '9', headTail{"pinball", "2000"}},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.ToBetween(c.min, c.max)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestWhileBetween(t *testing.T) {
	cases := []betweenCase{
		{"pinball2000", 'a', 'z', headTail{"pinball", "2000"}},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.WhileBetween(c.min, c.max)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

// ---------------------------------------------------------------------- match
type matchCase struct {
	inp string
	f   MatchFunc
	exp headTail
}

func TestIfMatch(t *testing.T) {
	cases := []matchCase{
		{"PIN", unicode.IsUpper, headTail{"P", "IN"}},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.IfMatch(c.f)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestToMatch(t *testing.T) {
	cases := []matchCase{
		{"PINball", unicode.IsLower, headTail{"PIN", "ball"}},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.ToMatch(c.f)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestWhileMatch(t *testing.T) {
	cases := []matchCase{
		{"123...", unicode.IsDigit, headTail{"123", "..."}},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.WhileMatch(c.f)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

// ---------------------------------------------------------------------- state
func TestLineCol(t *testing.T) {
	cases := []struct {
		diff int
		inp  string
		tab  int
		line int
		col  int
	}{
		{3, "abcdefgh", 1, 1, 4},
		{0, "abcdefgh", 1, 1, 1},
		{0, "", 1, 1, 1},
		{3, "abc", 1, 1, 4},
		{8, "\nabcd\n\nefgh\n\n", 1, 4, 2},
		{0, "abc\ndef\nghi", 1, 1, 1},
		{8, "\nabcd\n\n\tefgh\n\n", 4, 4, 5},
		{8, "\nabcd\n\n\tefgh\n\n", 4, 4, 5},
		{3, "\t\tabcdefgh", 4, 1, 10},
		{0, "\t\tabcdefgh", 4, 1, 1},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.Move(c.diff)
		line, col := sca.LineCol(c.tab)
		if line != c.line {
			t.Errorf("%d unexpected line: %d != %d", i, line, c.line)
		}
		if col != c.col {
			t.Errorf("%d unexpected col: %d != %d", i, col, c.col)
		}
	}
}
