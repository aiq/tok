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
