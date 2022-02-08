package tok

import "fmt"

//------------------------------------------------------------------------------

// Token marks a sub string in a Scanner.
type Token struct {
	from Marker
	to   Marker
}

func MakeToken(a, b Marker) Token {
	t := Token{a, b}
	if t.from > t.to {
		t.from, t.to = t.to, t.from
	}
	return t
}

// After reports whether t is after oth.
func (t Token) After(oth Token) bool {
	return t.from >= oth.to
}

// Before reports whether t is before oth.
func (t Token) Before(oth Token) bool {
	return t.to <= oth.from
}

func (t Token) Len() int {
	return int(t.to - t.from)
}

// Merge creates a Token that convers both Tokens t and oth.
func (t Token) Merge(oth Token) Token {
	if t.from > oth.from {
		t.from = oth.from
	}
	if t.to < oth.to {
		t.to = oth.to
	}
	return t
}

// Clash returns true if both tokens overlap, but no one covers the other.
func (t Token) Clashes(oth Token) bool {
	return !t.Before(oth) && !t.After(oth) && !t.Covers(oth) && !oth.Covers(t)
}

// Covers returns true if sub is a sub Token of t.
func (t Token) Covers(sub Token) bool {
	return t.from <= sub.from && t.to >= sub.to
}

// Split splits a Token into two parts via sep.
func (t Token) Split(sep Token) (Token, Token) {
	return MakeToken(t.from, sep.from), MakeToken(sep.to, t.to)
}

// String returns a readable representation of a Token.
func (t Token) String() string {
	return fmt.Sprintf("[%d-%d)", t.from, t.to)
}

// Returns the sub string that t represents.
func (s *Scanner) Get(t Token) string {
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
	return MakeToken(a, b), res
}
