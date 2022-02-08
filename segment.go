package tok

import "fmt"

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
		res = append(res, Segment{left, s.Get(left.Token)})
		res = append(res, Segment{v, s.Get(v.Token)})
	}
	res = append(res, Segment{rest, s.Get(rest.Token)})
	return res, nil
}
