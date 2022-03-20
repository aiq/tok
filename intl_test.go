package tok

import "testing"

func TestGetPrefix(t *testing.T) {
	cases := []struct {
		inp   string
		n     int
		exp   string
		isSub bool
	}{
		{"Hi, 世界", 5, "Hi, 世", true},
	}
	for i, c := range cases {
		prefix, ok := getPrefix(c.inp, c.n)
		if prefix != c.exp {
			t.Errorf("%d unexpected prefix: %q != %q", i, prefix, c.exp)
		}
		if ok != c.isSub {
			t.Errorf("%d unexpected ok value: %v", i, ok)
		}
	}
}

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
