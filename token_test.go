package tok

import "testing"

func TestSegmentate(t *testing.T) {
	T := func(a int, b int) Token {
		return Token{Marker(a), Marker(b)}
	}
	cases := []struct {
		str    string
		tokens []Token
		exp    []Segment
	}{
		{
			"abcdefgh", []Token{{0, 1}, {2, 4}, {7, 8}}, []Segment{
				{false, T(0, 0), ""}, {true, T(0, 1), "a"}, {false, T(1, 2), "b"}, {true, T(2, 4), "cd"},
				{false, T(4, 7), "efg"}, {true, T(7, 8), "h"}, {false, T(8, 8), ""},
			},
		},
	}
	for i, c := range cases {
		sca := NewScanner(c.str)
		res, err := sca.Segmentate(c.tokens)
		if err != nil {
			t.Error(unexpError(i, err))
		}
		if len(res) != len(c.exp) {
			t.Errorf("%d wrong results: %d != %d", i, len(res), len(c.exp))
		}
		for j, seg := range res {
			if !seg.Equal(c.exp[j]) {
				t.Errorf("%d unexpected segment at %d: %v != %v", i, j, seg, c.exp[j])
			}
		}
	}
}
