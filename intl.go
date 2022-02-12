package tok

import (
	"unicode/utf8"
)

//------------------------------------------------------------------------------

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

//------------------------------------------------------------------------------

func getPrefix(str string, n int) (string, bool) {
	count := 0
	for i := range str {
		if count == n {
			return str[:i], true
		}
		count++
	}
	return str, false
}
