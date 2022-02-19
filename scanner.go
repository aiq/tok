package tok

import (
	"strings"
	"unicode/utf8"
)

type Scanner struct {
	full string
	pos  int
}

// Creates a new Scanner to scan the str string.
func NewScanner(str string) *Scanner {
	return &Scanner{
		full: str,
		pos:  0,
	}
}

// ----------------------------------------------------------------------- scope

// Type of functions to mark the use of the scanner.
// The main use case are lambda functions.
type ScopeFunc func() bool

// Scopes a lambda functions and moves s back if f does not success.
func (s *Scanner) Scope(f ScopeFunc) bool {
	m := s.Mark()
	res := f()
	if !res {
		s.ToMarker(m)
	}
	return res
}

// Returns the sub string that was scanned by f.
func (s *Scanner) Trace(f ScopeFunc) (string, bool) {
	m := s.Mark()
	res := f()
	str := s.Since(m)
	if !res {
		s.ToMarker(m)
	}
	return str, res
}

// Repeats f until it returns false.
func (s *Scanner) While(f ScopeFunc) bool {
	start := s.Mark()
	for f() {
	}
	end := s.Mark()
	return start < end
}

// ScanString reads n runes as string from the scanner.
// The scanner will only be moved if all n runes can be read from the scanner.
// Returns true if s was moved, otherwise false.
func (s *Scanner) ScanString(n int) (string, bool) {
	sub, ok := getPrefix(s.Tail(), n)
	if ok {
		ok = s.Move(len(sub))
	}
	return sub, ok
}

// ---------------------------------------------------------------------- string
// Moves s the length of str forward if Tail() has str as the prefix.
// Returns true if s was moved, otherwise false.
func (s *Scanner) If(str string) bool {
	if strings.HasPrefix(s.Tail(), str) {
		return s.Move(len(str))
	}
	return false
}

// Moves s the length of the value in strs forward if Tail() that is the prefix.
// Returns true if s was moved, otherwise false.
func (s *Scanner) IfAny(strs ...string) bool {
	for _, str := range strs {
		if s.If(str) {
			return true
		}
	}
	return false
}

// Moves s to the first appearance of str in Tail().
// Returns true if s was moved, otherwise false.
func (s *Scanner) To(str string) bool {
	i := strings.Index(s.Tail(), str)
	if i != -1 {
		return s.Move(i)
	}
	return false
}

// ------------------------------------------------------------------------ fold
// Moves s the length of str forward if Tail() has str as the prefix under Unicode case-folding.
// Returns true if s was moved, otherwise false.
func (s *Scanner) IfFold(str string) bool {
	i := len(str)
	prefix := s.Tail()[:i]
	if strings.EqualFold(prefix, str) {
		return s.Move(i)
	}
	return false
}

// Moves s to the first appearance of str in Tail() under Unicode case-folding.
// Returns true if s was moved, otherwise false.
func (s *Scanner) ToFold(str string) bool {
	i := len(str)
	sub := s.Tail()
	for n := 0; len(sub) >= i+n; n++ {
		prefix := sub[n : n+i]
		if strings.EqualFold(prefix, str) {
			return s.Move(n)
		}
	}
	return false
}

// ------------------------------------------------------------------------ rune
// Moves s one rune value forward if the first rune in Tail() equals r.
// Returns true if s was moved, otherwise false.
func (s *Scanner) IfRune(r rune) bool {
	first, i := utf8.DecodeRuneInString(s.Tail())
	if i == -1 || first != r {
		return false
	}
	return s.Move(i)
}

// Moves s to the first rune in Tail() that matches with r.
// The function does not move s if no value matches with r.
// Returns true if a match was found, otherwise false.
func (s *Scanner) ToRune(r rune) bool {
	for i, v := range s.Tail() {
		if v == r {
			return s.Move(i)
		}
	}
	return false
}

// Moves s to the first rune in Tail() that does not match with r.
// Returns true if s was moved, otherwise false.
func (s *Scanner) WhileRune(r rune) bool {
	if len(s.Tail()) == 0 {
		return false
	}

	for i, v := range s.Tail() {
		if v != r {
			return s.Move(i)
		}
	}
	return s.ToEnd()
}

// --------------------------------------------------------------------- anyrune
// Moves s one rune forward if the first rune in Tail() is any of the runes in str.
// Returns true if s was moved, otherwise false.
func (s *Scanner) IfAnyRune(str string) bool {
	for _, r := range str {
		if s.IfRune(r) {
			return true
		}
	}
	return false
}

// Moves s to the first rune in Tail() that matches any of the runes in str.
// Returns true if a match was found, otherwise false.
func (s *Scanner) ToAnyRune(str string) bool {
	i := strings.IndexAny(s.Tail(), str)
	if i != -1 {
		return s.Move(i)
	}
	return false
}

// Moves s to the first rune in Tail() that does not match any of the runes in str.
// Returns true if s was moved, otherwise false.
func (s *Scanner) WhileAnyRune(str string) bool {
	if len(s.Tail()) == 0 {
		return false
	}

	for i, v := range s.Tail() {
		if !strings.ContainsRune(str, v) {
			if i == 0 {
				return false
			}
			return s.Move(i)
		}
	}
	return true
}

// --------------------------------------------------------------------- between

func inRange(min, val, max rune) bool {
	return min <= val && val <= max
}

// Moves s one rune forward if the first rune in Tail() is >= min and <= max.
// Returns true if s was moved, otherwise false.
func (s *Scanner) IfBetween(min, max rune) bool {
	if len(s.Tail()) == 0 {
		return false
	}
	val, i := utf8.DecodeRuneInString(s.Tail())
	if !inRange(min, val, max) {
		return false
	}
	return s.Move(i)
}

// Moves s to the first rune in Tail() that is >= min and <= max.
// Returns true if a match was found, otherwise false.
func (s *Scanner) ToBetween(min, max rune) bool {
	if len(s.Tail()) == 0 {
		return false
	}

	for i, v := range s.Tail() {
		if inRange(min, v, max) {
			return s.Move(i)
		}
	}
	return false
}

// Moves s to the first rune in Tail() that is < min or > max.
// Returns true if s was moved, otherwise false.
func (s *Scanner) WhileBetween(min, max rune) bool {
	if len(s.Tail()) == 0 {
		return false
	}

	for i, v := range s.Tail() {
		if !inRange(min, v, max) {
			return s.Move(i)
		}
	}
	return s.ToEnd()
}

// ---------------------------------------------------------------------- match
// Type of functions to check a rune.
type MatchFunc func(rune) bool

// Moves s one rune value forward if the first rune in Tail() passes the check.
// Returns true if s was moved, otherwise false.
func (s *Scanner) IfMatch(check MatchFunc) bool {
	first, i := utf8.DecodeRuneInString(s.Tail())
	if i != 0 && check(first) {
		return s.Move(i)
	}
	return false
}

// Moves s to the first rune in Tail() that passes the check.
// Returns true if s was moved, otherwise false.
func (s *Scanner) ToMatch(check MatchFunc) bool {
	for i, v := range s.Tail() {
		if check(v) {
			return s.Move(i)
		}
	}
	return false
}

// Moves s to the first rune in Tail() that does not pass the check.
// Returns true if s was moved, otherwise false.
func (s *Scanner) WhileMatch(check MatchFunc) bool {
	if len(s.Tail()) == 0 {
		return false
	}

	for i, v := range s.Tail() {
		if !check(v) {
			if i == 0 {
				return false
			}
			return s.Move(i)
		}
	}

	return false
}

// ----------------------------------------------------------------------- state
// Returns the right side from the current position in s.
func (s *Scanner) Tail() string {
	return s.full[s.pos:]
}

// Returns the left side from the current position in s.
func (s *Scanner) Head() string {
	return s.full[:s.pos]
}

// Returns the current line and column in the full string.
func (s *Scanner) LineCol(tab int) (int, int) {
	lines := strings.Split(strings.ReplaceAll(s.Head(), "\r\n", "\n"), "\n")
	if len(lines) == 0 {
		return 0, 0
	}
	last := lines[len(lines)-1]
	tabs := strings.Count(last, "\t")
	n := (len(last) - tabs) + tabs*tab
	return len(lines), n + 1
}

// A positive value moves s n bytes to the right, a negative value moves s n bytes to the left.
func (s *Scanner) Move(n int) bool {
	npos := s.pos + n
	if 0 > npos || npos > len(s.full) {
		return false
	}
	s.pos = npos
	return true
}

// A positve value moves s n runes to the right, a negative value moves s n runes to the left.
func (s *Scanner) MoveRunes(n int) bool {
	c := 0
	if n > 0 {
		for i, r := range s.Tail() {
			c++
			if n == c {
				return s.Move(i + utf8.RuneLen(r))
			}
		}
		return false
	} else if n < 0 {
		itr := makeRevItr(s.Head())
		for itr.next() {
			c--
			if n == c {
				return s.Move(-(itr.pos() + utf8.RuneLen(itr.Val)))
			}
		}
		return false
	}
	return true
}

// Returns true if s is at the end, otherwise false.
func (s *Scanner) AtEnd() bool {
	return len(s.full) == s.pos
}

// Returns true if s is at the start, otherwise false.
func (s *Scanner) AtStart() bool {
	return s.pos == 0
}

// ---------------------------------------------------------------------- Marker
// Marker represents a position in the text that will be scanned.
type Marker int

// Moves s to the marked position.
// Returns true if s was moved, otherwise false.
func (s *Scanner) ToMarker(m Marker) bool {
	if m < 0 && Marker(len(s.full)) <= m {
		return false
	}
	s.pos = int(m)
	return true
}

// Moves s to the end of the text that should be scanned.
func (s *Scanner) ToEnd() bool {
	return s.ToMarker(Marker(len(s.full)))
}

// Moves s to the start of the text that should be scanned.
func (s *Scanner) ToStart() bool {
	return s.ToMarker(Marker(0))
}

// Returns a Marker to mark the current positon in the text.
func (s *Scanner) Mark() Marker {
	return Marker(s.pos)
}
