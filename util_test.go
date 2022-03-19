package tok

import "testing"

func TestSubStringFrom(t *testing.T) {
	cases := []struct {
		str  string
		from int
		n    int
		exp  string
	}{
		{"abcdef", 0, 3, "abc"},
		{"äbcdef", 0, 3, "äbc"},
		{"äöü", 0, 4, "äöü"},
		{"äöü", 2, 4, "ü"},
	}
	for i, c := range cases {
		res := subStringFrom(c.str, c.from, c.n)
		if res != c.exp {
			t.Errorf("%d unexpected preview value: %q != %q", i, res, c.exp)
		}
	}
}
