package tok

import (
	"testing"
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
