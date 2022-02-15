package tok

import (
	"fmt"
	"sort"
)

//------------------------------------------------------------------------------
// Segment is a Token with additional Meta-Information that can be stored in Info.
type Segment struct {
	Info string
	Token
}

// Known reports if information about this segment exist.
func (seg Segment) Known() bool {
	return seg.Info != ""
}

// Split splits a Segment into two parts via sep.
func (v Segment) Split(sep Segment) (Segment, Segment) {
	l, r := v.Token.Split(sep.Token)
	return Segment{v.Info, l}, Segment{v.Info, r}
}

// String returns a readable representation of a Segment.
func (v Segment) String() string {
	return v.Info + v.Token.String()
}

type segSorter struct {
	values  []Segment
	cmpFunc func(Segment, Segment) bool
}

func (s *segSorter) Len() int {
	return len(s.values)
}

func (s *segSorter) Less(i, j int) bool {
	return s.cmpFunc(s.values[i], s.values[j])
}

func (s *segSorter) Swap(i, j int) {
	s.values[i], s.values[j] = s.values[j], s.values[i]
}

// SortSegments sorts the segments in a slice.
// Segments that cover other segments will appear before the covered values.
func SortSegments(values []Segment) {
	sorter := &segSorter{
		values: values,
		cmpFunc: func(a, b Segment) bool {
			return a.Covers(b.Token) || a.Before(b.Token)
		},
	}
	sort.Stable(sorter)
}

// SortSegmentsByOrder sorts like SortSegments with an order as additional tiebreaker.
// The appearence of an information in the order slice determines the order for segments with equal tokens.
func SortSegmentsByOrder(values []Segment, order []string) {
	sorter := &segSorter{
		values: values,
		cmpFunc: func(a, b Segment) bool {
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

// Segmentate splits a segment into subsegments.
func (s Segment) Segmentate(segments []Segment) ([]Segment, error) {
	res := []Segment{}
	left := Segment{}
	rest := s
	for _, seg := range segments {
		if !rest.Covers(seg.Token) {
			return res, fmt.Errorf("invalid token %s", seg.String())
		}
		left, rest = rest.Split(seg)
		if left.Len() > 0 {
			res = append(res, left)
		}
		res = append(res, seg)
	}
	if rest.Len() > 0 {
		res = append(res, rest)
	}
	return res, nil
}

// Segmentate splits the full string of a Scanner into segments.
func (s *Scanner) Segmentate(segments []Segment) ([]Segment, error) {
	rest := Segment{"", MakeToken(0, Marker(len(s.full)))}
	return rest.Segmentate(segments)
}
