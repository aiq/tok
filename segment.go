package tok

import (
	"fmt"
	"sort"
)

//------------------------------------------------------------------------------

type Value struct {
	Info string
	Token
}

func (v Value) Split(sep Value) (Value, Value) {
	l, r := v.Token.Split(sep.Token)
	return Value{v.Info, l}, Value{v.Info, r}
}

func (v Value) String() string {
	return v.Info + v.Token.String()
}

type valueSorter struct {
	values  []Value
	cmpFunc func(Value, Value) bool
}

func (s *valueSorter) Len() int {
	return len(s.values)
}

func (s *valueSorter) Less(i, j int) bool {
	return s.cmpFunc(s.values[i], s.values[j])
}

func (s *valueSorter) Swap(i, j int) {
	s.values[i], s.values[j] = s.values[j], s.values[i]
}

func SortValues(values []Value) {
	sorter := &valueSorter{
		values: values,
		cmpFunc: func(a, b Value) bool {
			return a.Covers(b.Token) || a.Before(b.Token)
		},
	}
	sort.Stable(sorter)
}

func SortValuesByOrder(values []Value, order []string) {
	sorter := &valueSorter{
		values: values,
		cmpFunc: func(a, b Value) bool {
			if a.Token != b.Token {
				return a.Covers(b.Token) || a.Before(b.Token)
			}
			ai, bi := -1, -1
			for i, info := range order {
				if a.Info == info {
					ai = i
				} else if b.Info == info {
					bi = i
				} else if ai != -1 && bi != -1 {
					break
				}
			}

			return ai < bi
		},
	}
	sort.Stable(sorter)
}

//------------------------------------------------------------------------------

// Segment groups a Token with the sub string from a Scanner.
type Segment struct {
	Value
	Sub string
}

// Known reports if this segment is identified.
func (seg Segment) Known() bool {
	return seg.Info != ""
}

// String returns a readable representation of a Segment.
func (seg Segment) String() string {
	return fmt.Sprintf("%s%q", seg.Value, seg.Sub)
}

// Segmentate splits the full string of a Scanner into segments.
func (s *Scanner) Segmentate(values []Value) ([]Segment, error) {
	res := []Segment{}
	left := Value{}
	rest := Value{"", MakeToken(0, Marker(len(s.full)))}
	for _, v := range values {
		if !rest.Covers(v.Token) {
			return res, fmt.Errorf("invalid token %s", v.String())
		}
		left, rest = rest.Split(v)
		if left.Len() > 0 {
			res = append(res, Segment{left, s.Get(left.Token)})
		}
		res = append(res, Segment{v, s.Get(v.Token)})
	}
	if rest.Len() > 0 {
		res = append(res, Segment{rest, s.Get(rest.Token)})
	}
	return res, nil
}
