package tok

import (
	"testing"
	"unicode/utf8"
)

func TestReader(t *testing.T) {
	cases := []struct {
		inp    string
		reader Reader
		tail   string
	}{
		{
			"var i =  \n 456;",
			Seq(Fold("VAR"), WS(), Lit("i ="), Many(WS()), Many(Between('0', '9')), Lit(";")),
			"",
		},
		{
			"var i =  \n 456;",
			Seq(Lit("var i ="), Many(WS()), Many(Digit()), Lit(";")),
			"",
		},
	}
	for i, c := range cases {
		sca := NewScanner(c.inp)
		err := sca.Use(c.reader)
		if err != nil {
			t.Errorf("%d unexpected error: %v", i, err)
		}
		if c.tail != sca.Tail() {
			t.Errorf("%d unexpected tail value: %q != %q", i, c.tail, sca.Tail())
		}
	}
}

func TestReaderWhat(t *testing.T) {
	cases := []struct {
		r   Reader
		exp string
	}{
		{Any(Rune('!'), Lit("abc")), `[ '!' "abc" ]`},
		{Between('!', '☃'), `<!☃>`},
		{Between(0, utf8.MaxRune), `<\x00\U0010ffff>`},
		{BetweenAny("a-zA-Z"), `[< az AZ >]`},
		{Fold("true"), `~"true"`},
		{Holey('a', 'z', "ox"), `(<az> - "ox")`},
		{Seq(Rune('!'), Many(AnyRune(" +-")), Lit("abc")), `> '!' +[" +-"] "abc" >`},
		{To(Bool("")), `->bool{}`},
		{Uint(16, 64), `uint{16,64}`},
		{WS(), `[" \r\n\t"]`},
		{Zom(WS()), `*[" \r\n\t"]`},
	}
	for i, c := range cases {
		if c.r.What() != c.exp {
			t.Errorf("%d unexpected what message: %s", i, c.r.What())
		}
	}
}
