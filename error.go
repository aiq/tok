package tok

import "fmt"

// Error type that ReadFunc and the Reader here return.
type ReadError struct {
	What string
	Line int
	Col  int
}

func (e ReadError) DeeperAs(oth ReadError) bool {
	return e.Line >= oth.Line && e.Col > oth.Col
}

func (e ReadError) Error() string {
	return fmt.Sprintf("not able to read %s at line:%d col:%d", e.What, e.Line, e.Col)
}

// Generates a ReadError for name.
func (s *Scanner) ErrorFor(name string) error {
	l, c := s.LineCol()
	return fmt.Errorf("not able to read %s at line:%d col:%d", name, l, c)
}

// Generates a ErrorFor if ok is false, otherwise returns the function nil.
func (s *Scanner) BoolErrorFor(ok bool, name string) error {
	if !ok {
		l, c := s.LineCol()
		return ReadError{name, l, c}
	}
	return nil
}
