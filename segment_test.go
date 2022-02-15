package tok

import "testing"

func TestSortSegments(t *testing.T) {
	cases := []struct {
		inp []Segment
		exp []Segment
	}{
		{
			[]Segment{S("obj", 2, 18), S("val", 10, 16), S("id", 3, 8), S("text", 0, 20)},
			[]Segment{S("text", 0, 20), S("obj", 2, 18), S("id", 3, 8), S("val", 10, 16)},
		},
	}
	for i, c := range cases {
		SortSegments(c.inp)
		for j, v := range c.inp {
			if v != c.exp[j] {
				t.Errorf("%d unexpected values at %d: %v != %v", i, j, v, c.exp[j])
			}
		}
	}
}

func TestSortSegmentsByOrder(t *testing.T) {
	cases := []struct {
		inp   []Segment
		order []string
		exp   []Segment
	}{
		{
			[]Segment{S("obj", 2, 18), S("val", 10, 16), S("id", 3, 8), S("text", 0, 20), S("member", 2, 18)},
			[]string{"member", "obj"},
			[]Segment{S("text", 0, 20), S("member", 2, 18), S("obj", 2, 18), S("id", 3, 8), S("val", 10, 16)},
		},
	}
	for i, c := range cases {
		SortSegmentsByOrder(c.inp, c.order)
		for j, v := range c.inp {
			if v != c.exp[j] {
				t.Errorf("%d unexpected values at %d: %v != %v", i, j, v, c.exp[j])
			}
		}
	}
}

func TestSegmentate(t *testing.T) {
	cases := []struct {
		str      string
		segments []Segment
		exp      []Segment
	}{
		{
			"abcdefgh", []Segment{S("v", 0, 1), S("v", 2, 4), S("v", 7, 8)}, []Segment{
				S("v", 0, 1), S("", 1, 2), S("v", 2, 4),
				S("", 4, 7), S("v", 7, 8),
			},
		},
	}
	for i, c := range cases {
		sca := NewScanner(c.str)
		res, err := sca.Segmentate(c.segments)
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
