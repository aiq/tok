package tok

type Token struct {
	a Marker
	b Marker
}

func (s *Scanner) Get(t Token) string {
	if t.a > t.b {
		t.a, t.b = t.b, t.a
	}
	return s.full[t.a:t.b]
}

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
