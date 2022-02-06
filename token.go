package tok

import "fmt"

// Token marks a sub string.
type Token struct {
	from Marker
	to   Marker
}

func (t Token) Norm() Token {
	if t.from > t.to {
		t.from, t.to = t.to, t.from
	}
	return t
}

func (t Token) Equal(oth Token) bool {
	t, oth = t.Norm(), oth.Norm()
	return t.from == oth.from && t.to == oth.to
}

func (t Token) Has(oth Token) bool {
	t, oth = t.Norm(), oth.Norm()
	return t.from <= oth.from && t.to >= oth.to
}

func (t Token) Split(oth Token) (Token, Token) {
	t, oth = t.Norm(), oth.Norm()
	return Token{t.from, oth.from}, Token{oth.to, t.to}
}

func (t Token) Len() int {
	t = t.Norm()
	return int(t.to - t.from)
}

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

type Segment struct {
	FromToken bool
	Token
	Value string
}

func (seg Segment) Equal(oth Segment) bool {
	return seg.FromToken == oth.FromToken &&
		seg.Token.Equal(oth.Token) &&
		seg.Value == oth.Value
}

func (seg Segment) String() string {
	from := ""
	if seg.FromToken {
		from = ">"
	}
	return fmt.Sprintf("%s%s%q", seg.Token.String(), from, seg.Value)
}

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
