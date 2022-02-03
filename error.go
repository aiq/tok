package tok

import "fmt"

type ReadError struct {
	What string
	Line int
	Col  int
}

func (e ReadError) Error() string {
	return fmt.Sprintf("not able to read %q at line:%d col:%d", e.What, e.Line, e.Col)
}

func MakeReadError(what string, line int, col int) ReadError {
	return ReadError{what, line, col}
}

func (s *Scanner) ErrorFor(name string) error {
	l, c := s.LineCol()
	return fmt.Errorf("not able to read %q at line:%d col:%d", name, l, c)
}

func (s *Scanner) BoolErrorFor(ok bool, name string) error {
	if !ok {
		l, c := s.LineCol()
		return MakeReadError(name, l, c)
	}
	return nil
}
