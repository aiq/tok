package tok

import (
	"strings"
	"unicode/utf8"
)

// -----------------------------------------------------------------------------
// Creates a new Scanner to scan the str string and moves the scanner to the end.
func NewRevScanner(str string) *Scanner {
	sca := NewScanner(str)
	sca.ToEnd()
	return sca
}

// ---------------------------------------------------------------------- string
// Moves s the length of str backward if Head() has str as the prefix.
// Returns true if s was moved, otherwise false.
func (s *Scanner) RevIf(str string) bool {
	if strings.HasSuffix(s.Head(), str) {
		return s.Move(-len(str))
	}
	return false
}

// Moves s the length fo the value in strs backward if Head() is the suffix.
// Returns true if s was moved, otherwise false.
func (s *Scanner) RevIfAny(strs []string) bool {
	for _, str := range strs {
		if s.RevIf(str) {
			return true
		}
	}
	return false
}

// Moves s to the last appearance of str in Head().
// Returns true if s was moved, otherwise false.
func (s *Scanner) RevTo(str string) bool {
	itr := makeRevItr(s.Head())
	for itr.next() {
		if strings.HasSuffix(itr.Str, str) {
			return s.Move(-itr.pos())
		}
	}

	return false
}

// ------------------------------------------------------------------------ fold
// Moves s the length of str backward if Head() has str as the suffix under Unicode case-folding.
// Returns true if s was moved, otherwise false.
func (s *Scanner) RevIfFold(str string) bool {
	i := len(str)
	head := s.Head()
	suffix := head[len(head)-i:]
	if strings.EqualFold(suffix, str) {
		return s.Move(-i)
	}
	return false
}

// Moves s to the last appearance of str in Head() under Unicode case-folding.
// Returns true if s was moved, otherwise false.
func (s *Scanner) RevToFold(str string) bool {
	i := len(str)
	itr := makeRevItr(s.Head())
	for itr.next() {
		suffix := itr.Str[len(itr.Str)-i:]
		if strings.EqualFold(suffix, str) {
			return s.Move(-itr.pos())
		}
	}

	return false
}

// ------------------------------------------------------------------------ rune
// Moves s one rune value backward if the last rune in Head() equals r.
// Returns true if s was moved, otherwise false.
func (s *Scanner) RevIfRune(r rune) bool {
	last, i := utf8.DecodeLastRuneInString(s.Head())
	if i == -1 || last != r {
		return false
	}
	return s.Move(-i)
}

// Moves s to the last rune in Head() that matches with r.
// The function does not move s if not values matches with r.
// Returns true if a match was found, otherwise false.
func (s *Scanner) RevToRune(r rune) bool {
	itr := makeRevItr(s.Head())
	for itr.next() {
		if itr.Val == r {
			return s.Move(-itr.pos())
		}
	}

	return false
}

// Moves s to the last rune in Head() that does not match with r.
// Returns true if s was moved, otherwise false.
func (s *Scanner) RevWhileRune(r rune) bool {
	if len(s.Head()) == 0 {
		return false
	}

	itr := makeRevItr(s.Head())
	for itr.next() {
		if itr.Val != r {
			if itr.pos() == len(s.Head()) {
				return false
			}
			return s.Move(-itr.pos())
		}
	}

	return s.Move(-itr.pos())
}

// --------------------------------------------------------------------- anyrune
// Moves s one rune backward if the last rune in Head() is any of the runes in str.
// Returns true if s was moved, otherwise false.
func (s *Scanner) RevIfAnyRune(str string) bool {
	last, i := utf8.DecodeLastRuneInString(s.Head())
	if i != -1 && strings.ContainsRune(str, last) {
		return s.Move(-i)
	}
	return false
}

// Moves s to the last rune in Head() that matches any of the runes in str.
// Returns true if a match was found, otherwise false.
func (s *Scanner) RevToAnyRune(str string) bool {
	itr := makeRevItr(s.Head())
	for itr.next() {
		if strings.ContainsRune(str, itr.Val) {
			return s.Move(-itr.pos())
		}
	}

	return false
}

// Moves s to the last rune in Head() that does not match any of the runes in str.
// Returns true if s was moved, otherwise false.
func (s *Scanner) RevWhileAnyRune(str string) bool {
	if len(s.Head()) == 0 {
		return false
	}

	itr := makeRevItr(s.Head())
	for itr.next() {
		if !strings.ContainsRune(str, itr.Val) {
			if itr.pos() == len(s.Head()) {
				return false
			}
			return s.Move(-itr.pos())
		}
	}
	return s.Move(-itr.pos())
}

// --------------------------------------------------------------------- between
// Moves s one rune backward if the last rune in Head() is >= min and <= max.
// Returns true if s was moved, otherwise false.
func (s *Scanner) RevIfBetween(min, max rune) bool {
	return false
}

// Moves s to the last rune in Head() that is >= min and <= max.
// Returns true if s was moved, otherwise false.
func (s *Scanner) RevToBetween(min, max rune) bool {
	return false
}

// Moves s to the last rune in Head() that is < min or > max.
// Returns true if s was moved, otherwise false.
func (s *Scanner) RevWhileBetween(min, max rune) bool {
	return false
}

// ---------------------------------------------------------------------- match
// Moves s one rune backward if the last rune in Head() passes the check.
// Returns true if s was moved, otherwise false.
func (s *Scanner) RevIfMatch(check MatchFunc) bool {
	return false
}

// Moves s to the last rune in Head() that passes the ckeck.
// Returns true if s was moved, otherwise false.
func (s *Scanner) RevToMatch(check MatchFunc) bool {
	return false
}

// Moves s to the last rune in Head() that does not pass the ckeck.
// Returns true if s was moved, otherwise false.
func (s *Scanner) RevWhileMatch(check MatchFunc) bool {
	return false
}
