package tok

// Token marks a sub string.
type Token struct {
	a Marker
	b Marker
}

// Returns the sub string that t represents.
func (s *Scanner) Get(t Token) string {
	if t.a > t.b {
		t.a, t.b = t.b, t.a
	}
	return s.full[t.a:t.b]
}

// Marks the sub string that was scanned by f.
func (s *Scanner) Tokenize(f ScopeFunc) (Token, bool) {
	a := s.Mark()
	res := f()
	b := s.Mark()
	if !res {
		s.ToMarker(a)
	}
	return Token{a, b}, res
}

type Segment struct {
	FromToken bool
	Token
	Value string
}

func (s *Scanner) Segmentate(tokens []Token) []Segment {
	res := []Segment{}

	return res
}
