package tok

import "testing"

// ---------------------------------------------------------------------- string
func TestRevIf(t *testing.T) {

}

func TestRefIfAny(t *testing.T) {

}

func TestRevTo(t *testing.T) {

}

// ------------------------------------------------------------------------ fold
func TestRevIfFold(t *testing.T) {
}

func TestRevToFold(t *testing.T) {

}

// ------------------------------------------------------------------------ rune
func TestRevIfRune(t *testing.T) {
	cases := []struct {
		inp string
		r   rune
		exp headTail
	}{
		{"xyz", 'z', headTail{"xy", "z"}},
	}

	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.ToEnd()
		sca.RevIfRune(c.r)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestRevToRune(t *testing.T) {
	cases := []struct {
		inp string
		r   rune
		exp headTail
	}{
		{"i am", ' ', headTail{"i ", "am"}},
		{"?...", '.', headTail{"?...", ""}},
		{"1. end", '1', headTail{"1", ". end"}},
	}

	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.ToEnd()
		sca.RevToRune(c.r)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

func TestRevWhileRune(t *testing.T) {
	cases := []struct {
		inp string
		r   rune
		exp headTail
	}{
		{"123-", '-', headTail{"123", "-"}},
		{"a....", '.', headTail{"a", "...."}},
		{"   ", ' ', headTail{"", "   "}},
	}

	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.ToEnd()
		sca.RevWhileRune(c.r)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

// --------------------------------------------------------------------- anyrune
func TestRevIfAnyRune(t *testing.T) {

}

func TestRevToAnyRune(t *testing.T) {

}

func TestRevWhileAnyRune(t *testing.T) {
	cases := []struct {
		inp string
		str string
		exp headTail
	}{
		{"123-", "-", headTail{"123", "-"}},
		{"-256", "1234567890", headTail{"-", "256"}},
		{"12344512", "1234567890", headTail{"", "12344512"}},
	}

	for i, c := range cases {
		sca := NewScanner(c.inp)
		sca.ToEnd()
		sca.RevWhileAnyRune(c.str)
		if e := c.exp.check(i, sca); e != nil {
			t.Errorf("%v", e)
		}
	}
}

// --------------------------------------------------------------------- between
func TestRevIfBetween(t *testing.T) {

}

func TestRevToBetween(t *testing.T) {

}

func TestRevWhileBetween(t *testing.T) {

}

// ---------------------------------------------------------------------- match
func TestRevIfMatch(t *testing.T) {

}

func TestRevToMatch(t *testing.T) {

}

func TestRevWhileMatch(t *testing.T) {

}
