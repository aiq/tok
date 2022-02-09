package tok

import "testing"

func TestReverseValues(t *testing.T) {
	cases := []struct {
		inp []Value
		exp []Value
	}{
		{
			[]Value{V("id", 3, 8), V("val", 10, 16), V("obj", 2, 18), V("text", 0, 20)},
			[]Value{V("text", 0, 20), V("obj", 2, 18), V("val", 10, 16), V("id", 3, 8)},
		},
	}
	for i, c := range cases {
		ReverseValues(c.inp)
		for j, v := range c.inp {
			if v != c.exp[j] {
				t.Errorf("%d unexpected values at %d: %v != %v", i, j, v, c.exp[j])
			}
		}
	}
}

func TestSegmentate(t *testing.T) {
	cases := []struct {
		str    string
		values []Value
		exp    []Segment
	}{
		{
			"abcdefgh", []Value{V("v", 0, 1), V("v", 2, 4), V("v", 7, 8)}, []Segment{
				{V("v", 0, 1), "a"}, {V("", 1, 2), "b"}, {V("v", 2, 4), "cd"},
				{V("", 4, 7), "efg"}, {V("v", 7, 8), "h"},
			},
		},
	}
	for i, c := range cases {
		sca := NewScanner(c.str)
		res, err := sca.Segmentate(c.values)
		if err != nil {
			t.Error(unexpError(i, err))
		}
		if len(res) != len(c.exp) {
			t.Errorf("%d wrong results: %d != %d", i, len(res), len(c.exp))
		} else {
			for j, seg := range res {
				if seg != c.exp[j] {
					t.Errorf("%d unexpected segment at %d: %v != %v", i, j, seg, c.exp[j])
				}
			}
		}
	}
}
