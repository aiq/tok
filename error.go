package tok

import "fmt"

// Error type that ReadFunc and the Reader here return.
type ReadError struct {
	Marker
	What string
}

// Later checks if e occurred later as oth.
func (e ReadError) Later(oth ReadError) bool {
	return e.Marker >= oth.Marker
}

// Error function to match the error interface.
func (e ReadError) Error() string {
	return fmt.Sprintf("not able to read %s at %d", e.What, e.Marker)
}

// Generates a ReadError for name.
func (s *Scanner) ErrorFor(name string) error {
	return ReadError{s.Mark(), name}
}

// Generates a ErrorFor if ok is false, otherwise returns the function nil.
func (s *Scanner) ErrorIfFalse(ok bool, name string) error {
	if !ok {
		return ReadError{s.Mark(), name}
	}
	return nil
}
