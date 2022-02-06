package tok

import "fmt"

// Token marks a sub string in a Scanner.
type Token struct {
	from Marker
	to   Marker
}

// Norm mormalizes a Token.
func (t Token) Norm() Token {
	if t.from > t.to {
		t.from, t.to = t.to, t.from
	}
	return t
}

// Equal reports if t and oth mark the same sub string in a Scanner.
func (t Token) Equal(oth Token) bool {
	t, oth = t.Norm(), oth.Norm()
	return t.from == oth.from && t.to == oth.to
}

// Has returns true if sub is a sub Token of t.
func (t Token) Has(sub Token) bool {
	t, sub = t.Norm(), sub.Norm()
	return t.from <= sub.from && t.to >= sub.to
}

// Split splits a Token into two parts via sep.
func (t Token) Split(sep Token) (Token, Token) {
	t, sep = t.Norm(), sep.Norm()
	return Token{t.from, sep.from}, Token{sep.to, t.to}
}

// String returns a readable representation of a Token.
func (t Token) String() string {
	return fmt.Sprintf("[%d-%d)", t.from, t.to)
}

// Returns the sub string that t represents.
func (s *Scanner) Get(t Token) string {
	t = t.Norm()
	return s.full[t.from:t.to]
}

// Marks the sub string that was scanned by f.
func (s *Scanner) Tokenize(f ScopeFunc) (Token, bool) {
	a := s.Mark()
	res := f()
	b := s.Mark()
	if !res {
		s.ToMarker(a)
	}
	return Token{a, b}.Norm(), res
}

// Segment groups a Token with the sub string from a Scanner.
type Segment struct {
	FromToken bool
	Token
	Value string
}

// Equal reports if seg and oth mark are equal.
func (seg Segment) Equal(oth Segment) bool {
	return seg.FromToken == oth.FromToken &&
		seg.Token.Equal(oth.Token) &&
		seg.Value == oth.Value
}

// String returns a readable representation of a Segment.
func (seg Segment) String() string {
	from := ""
	if seg.FromToken {
		from = ">"
	}
	return fmt.Sprintf("%s%s%q", seg.Token.String(), from, seg.Value)
}

// Segmentate splits the full string of a Scanner into segments.
func (s *Scanner) Segmentate(tokens []Token) ([]Segment, error) {
	res := []Segment{}
	left := Token{}
	rest := Token{0, Marker(len(s.full))}
	for _, t := range tokens {
		t = t.Norm()
		if !rest.Has(t) {
			return res, fmt.Errorf("invalid token %s", t.String())
		}
		left, rest = rest.Split(t)
		res = append(res, Segment{false, left, s.Get(left)})
		res = append(res, Segment{true, t, s.Get(t)})
	}
	res = append(res, Segment{false, rest, s.Get(rest)})
	return res, nil
}
