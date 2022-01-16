package tok

import (
	"strings"
	"unicode/utf8"
)

type Scanner struct {
	full string
	pos  int
}

func NewScanner(str string) *Scanner {
	return &Scanner{
		full: str,
		pos:  0,
	}
}

func (s *Scanner) move(i int) bool {
	i += s.pos
	if i > len(s.full) {
		return false
	}
	s.pos = i
	return true
}

// ----------------------------------------------------------------------- scope

type ScopeFunc func() bool

func (s *Scanner) Scope(f ScopeFunc) bool {
	marker := s.Mark()
	res := f()
	if !res {
		s.ToMarker(marker)
	}
	return res
}

func (s *Scanner) Traced(f ScopeFunc) (string, bool) {
	m := s.Mark()
	res := f()
	return s.Since(m), res
}

func (s *Scanner) While(f ScopeFunc) bool {
	start := s.Mark()
	for f() {
	}
	end := s.Mark()
	return start < end
}

func (s *Scanner) Opt(val bool) bool {
	return true
}

// ---------------------------------------------------------------------- string
func (s *Scanner) If(str string) bool {
	if strings.HasPrefix(s.Tail(), str) {
		return s.move(len(str))
	}
	return false
}

func (s *Scanner) IfAny(strs []string) bool {
	for _, str := range strs {
		if s.If(str) {
			return true
		}
	}
	return false
}

func (s *Scanner) To(str string) bool {
	i := strings.Index(s.Tail(), str)
	if i != -1 {
		return s.move(i)
	}
	return false
}

// ------------------------------------------------------------------------ fold
func (s *Scanner) IfFold(str string) bool {
	i := len(str)
	prefix := s.Tail()[:i]
	if strings.EqualFold(prefix, str) {
		return s.move(i)
	}
	return false
}

func (s *Scanner) ToFold(str string) bool {
	i := len(str)
	sub := s.Tail()
	for n := 0; len(sub) >= i+n; n++ {
		prefix := sub[n : n+i]
		if strings.EqualFold(prefix, str) {
			return s.move(n)
		}
	}
	return false
}

// ------------------------------------------------------------------------ rune
func (s *Scanner) IfRune(r rune) bool {
	first, i := utf8.DecodeRuneInString(s.Tail())
	if i == -1 || first != r {
		return false
	}
	return s.move(i)
}

func (s *Scanner) ToRune(r rune) bool {
	for i, v := range s.Tail() {
		if v == r {
			return s.move(i)
		}
	}
	return false
}

func (s *Scanner) WhileRune(r rune) bool {
	if len(s.Tail()) == 0 {
		return false
	}

	for i, v := range s.Tail() {
		if v != r {
			return s.move(i)
		}
	}
	return s.ToEnd()
}

// --------------------------------------------------------------------- anyrune
func (s *Scanner) IfAnyRune(str string) bool {
	for _, r := range str {
		if s.IfRune(r) {
			return true
		}
	}
	return false
}

func (s *Scanner) ToAnyRune(str string) bool {
	i := strings.IndexAny(s.Tail(), str)
	if i != -1 {
		return s.move(i)
	}
	return false
}

func (s *Scanner) WhileAnyRune(str string) bool {
	if len(s.Tail()) == 0 {
		return false
	}

	for i, v := range s.Tail() {
		if !strings.ContainsRune(str, v) {
			if i == 0 {
				return false
			}
			return s.move(i)
		}
	}
	return true
}

// --------------------------------------------------------------------- between

func inRange(min, val, max rune) bool {
	return min <= val && val <= max
}

func (s *Scanner) IfBetween(min, max rune) bool {
	val, i := utf8.DecodeRuneInString(s.Tail())
	if !inRange(min, val, max) {
		return false
	}
	return s.move(i)
}

func (s *Scanner) ToBetween(min, max rune) bool {
	for i, v := range s.Tail() {
		if inRange(min, v, max) {
			return s.move(i)
		}
	}
	return false
}

func (s *Scanner) WhileBetween(min, max rune) bool {
	if len(s.Tail()) == 0 {
		return false
	}

	for i, v := range s.Tail() {
		if !inRange(min, v, max) {
			return s.move(i)
		}
	}
	return s.ToEnd()
}

// ---------------------------------------------------------------------- match

type MatchFunc func(rune) bool

func (s *Scanner) IfMatch(check MatchFunc) bool {
	first, i := utf8.DecodeRuneInString(s.Tail())
	if i != 0 && check(first) {
		return s.move(i)
	}
	return false
}

func (s *Scanner) ToMatch(check MatchFunc) bool {
	for i, v := range s.Tail() {
		if check(v) {
			return s.move(i)
		}
	}
	return false
}

func (s *Scanner) WhileMatch(check MatchFunc) bool {
	if len(s.Tail()) == 0 {
		return false
	}

	for i, v := range s.Tail() {
		if !check(v) {
			if i == 0 {
				return false
			}
			return s.move(i)
		}
	}

	return false
}

// ----------------------------------------------------------------------- state

func (s *Scanner) Tail() string {
	return s.full[s.pos:]
}

func (s *Scanner) Head() string {
	return s.full[:s.pos]
}

func (s *Scanner) ScannedRune() rune {
	r, _ := utf8.DecodeLastRuneInString(s.Head())
	return r
}

func (s *Scanner) LineCol() (int, int) {
	lines := strings.Split(strings.ReplaceAll(s.Head(), "\r\n", "\n"), "\n")
	if len(lines) == 0 {
		return 0, 0
	}
	last := lines[len(lines)-1]
	return len(lines), len(last)
}

func (s *Scanner) AtEnd() bool {
	return len(s.Tail()) == 0
}

func (s *Scanner) ToEnd() bool {
	return s.ToMarker(Marker(len(s.full)))
}

func (s *Scanner) ToStart() bool {
	return s.ToMarker(Marker(0))
}

// ---------------------------------------------------------------------- Marker

type Marker int

func (s *Scanner) Since(m Marker) string {
	return s.full[m:s.pos]
}

func (s *Scanner) ToMarker(m Marker) bool {
	if m < 0 && Marker(len(s.full)) <= m {
		return false
	}
	s.pos = int(m)
	return true
}

func (s *Scanner) Mark() Marker {
	return Marker(s.pos)
}
