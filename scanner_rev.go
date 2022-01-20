package tok

import (
	"strings"
	"unicode/utf8"
)

type revItr struct {
	Str  string
	init int
	Val  rune
	vi   int
}

func makeRevItr(str string) revItr {
	return revItr{
		Str:  str,
		init: len(str),
		Val:  utf8.RuneError,
		vi:   0,
	}
}

func (itr *revItr) next() bool {
	if itr.Val != utf8.RuneError {
		itr.Str = itr.Str[0 : len(itr.Str)-itr.vi]
	}
	itr.Val, itr.vi = utf8.DecodeLastRuneInString(itr.Str)
	return itr.Val != utf8.RuneError
}

func (itr *revItr) pos() int {
	return itr.init - len(itr.Str)
}

// -----------------------------------------------------------------------------
func NewRevScanner(str string) *Scanner {
	sca := NewScanner(str)
	sca.ToEnd()
	return sca
}

// ---------------------------------------------------------------------- string
func (s *Scanner) RevIf(str string) bool {
	if strings.HasSuffix(s.Head(), str) {
		return s.move(-len(str))
	}
	return false
}

func (s *Scanner) RevIfAny(strs []string) bool {
	for _, str := range strs {
		if s.RevIf(str) {
			return true
		}
	}
	return false
}

func (s *Scanner) RevTo(str string) bool {
	itr := makeRevItr(s.Head())
	for itr.next() {
		if strings.HasSuffix(itr.Str, str) {
			return s.move(-itr.pos())
		}
	}

	return false
}

// ------------------------------------------------------------------------ fold
func (s *Scanner) RevIfFold(str string) bool {
	i := len(str)
	head := s.Head()
	suffix := head[len(head)-i:]
	if strings.EqualFold(suffix, str) {
		return s.move(-i)
	}
	return false
}

func (s *Scanner) RevToFold(str string) bool {
	i := len(str)
	itr := makeRevItr(s.Head())
	for itr.next() {
		suffix := itr.Str[len(itr.Str)-i:]
		if strings.EqualFold(suffix, str) {
			return s.move(-itr.pos())
		}
	}

	return false
}

// ------------------------------------------------------------------------ rune
func (s *Scanner) RevIfRune(r rune) bool {
	last, i := utf8.DecodeLastRuneInString(s.Head())
	if i == -1 || last != r {
		return false
	}
	return s.move(-i)
}

func (s *Scanner) RevToRune(r rune) bool {
	itr := makeRevItr(s.Head())
	for itr.next() {
		if itr.Val == r {
			return s.move(-itr.pos())
		}
	}

	return false
}

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
			return s.move(-itr.pos())
		}
	}

	return s.move(-itr.pos())
}

// --------------------------------------------------------------------- anyrune
func (s *Scanner) RevIfAnyRune(str string) bool {
	last, i := utf8.DecodeLastRuneInString(s.Head())
	if i != -1 && strings.ContainsRune(str, last) {
		return s.move(-i)
	}
	return false
}

func (s *Scanner) RevToAnyRune(str string) bool {
	itr := makeRevItr(s.Head())
	for itr.next() {
		if strings.ContainsRune(str, itr.Val) {
			return s.move(-itr.pos())
		}
	}

	return false
}

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
			return s.move(-itr.pos())
		}
	}
	return s.move(-itr.pos())
}

// --------------------------------------------------------------------- between
func (s *Scanner) RevIfBetween(min, max rune) bool {
	return false
}

func (s *Scanner) RevToBetween(min, max rune) bool {
	return false
}

func (s *Scanner) RevWhileBetween(min, max rune) bool {
	return false
}

// ---------------------------------------------------------------------- match
func (s *Scanner) RevIfMatch(check MatchFunc) bool {
	return false
}

func (s *Scanner) RevToMatch(check MatchFunc) bool {
	return false
}

func (s *Scanner) RevWhileMatch(check MatchFunc) bool {
	return false
}
