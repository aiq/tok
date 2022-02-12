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
		{To(Bool("")), `->bool{""}`},
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

func TestJanus(t *testing.T) {
	check := func(str, expStr string, e, expE error) {
		if str != expStr {
			t.Errorf("unexpected string: %q != %q", str, expStr)
		}
		if e != expE {
			t.Errorf("unexpected error value: %v", e)
		}
	}
	beg, end := Janus("i", Many(Between('a', 'z')))
	str, err := NewScanner("two two").CaptureUse(Seq(beg, WS(), end))
	check(str, "two two", err, nil)

	beg, end = Janus("c", Many(Rune('=')))
	comBeg := Seq(Rune('['), beg, Rune('['))
	comEnd := Seq(Rune(']'), end, Rune(']'))
	str, err = NewScanner("[==[long lua string]==] ~=").CaptureUse(Seq(comBeg, To(comEnd)))
	check(str, "[==[long lua string", err, nil)

	str, err = NewScanner("[==[long lua string]==] ~=").CaptureUse(Seq(comBeg, Past(comEnd)))
	check(str, "[==[long lua string]==]", err, nil)
}
