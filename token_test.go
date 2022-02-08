package tok

import (
	"testing"
)

func TestTokenCovers(t *testing.T) {
	t_02_23 := MakeToken(Marker(2), Marker(23))
	t_10_30 := MakeToken(Marker(10), Marker(30))
	t_11_19 := MakeToken(Marker(11), Marker(19))

	cases := []struct {
		t   Token
		oth Token
		exp bool
	}{
		{t_02_23, t_10_30, false},
		{t_10_30, t_02_23, false},
		{t_02_23, t_11_19, true},
		{t_11_19, t_02_23, false},
		{t_10_30, t_11_19, true},
		{t_11_19, t_10_30, false},
	}

	for _, c := range cases {
		res := c.t.Covers(c.oth)
		if res != c.exp {
			t.Errorf("unexpected Covers result between %v and %v: %v", c.t, c.oth, res)
		}
	}
}

func TestTokenClash(t *testing.T) {
	t_03_08 := MakeToken(Marker(3), Marker(8))
	t_02_23 := MakeToken(Marker(2), Marker(23))
	t_10_30 := MakeToken(Marker(10), Marker(30))
	t_11_19 := MakeToken(Marker(11), Marker(19))

	cases := []struct {
		t   Token
		oth Token
		exp bool
	}{
		{t_02_23, t_10_30, true},
		{t_10_30, t_02_23, true},
		{t_02_23, t_11_19, false},
		{t_11_19, t_02_23, false},
		{t_10_30, t_11_19, false},
		{t_11_19, t_10_30, false},
		{t_03_08, t_10_30, false},
		{t_10_30, t_03_08, false},
	}

	for _, c := range cases {
		res := c.t.Clashes(c.oth)
		if res != c.exp {
			t.Errorf("unexpected Clashes result between %v and %v: %v", c.t, c.oth, res)
		}
	}
}
