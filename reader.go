package tok

import "fmt"

type Reader interface {
	Read(s *Scanner) error
}

func (s *Scanner) Use(r Reader) error {
	return r.Read(s)
}

func (s *Scanner) TracedUse(r Reader) (string, error) {
	m := s.Mark()
	res := r.Read(s)
	return s.Since(m), res
}

func (s *Scanner) errorMsg() error {
	l, c := s.LineCol()
	return fmt.Errorf("not able to read at line %d col %d", l, c)
}

func (s *Scanner) boolError(ok bool) error {
	if !ok {
		return s.errorMsg()
	}
	return nil
}

//----------------------------------------------------------
type AnyReader struct {
	readers []Reader
}

func (r AnyReader) Read(s *Scanner) error {
	m := s.Mark()
	for _, sub := range r.readers {
		if err := sub.Read(s); err == nil {
			return nil
		}
	}
	s.ToMarker(m)
	return s.errorMsg()
}

func Any(list ...Reader) Reader {
	return AnyReader{list}
}

//----------------------------------------------------------
type SeqReader struct {
	readers []Reader
}

func (r SeqReader) Read(s *Scanner) error {
	m := s.Mark()
	var err error
	for _, sub := range r.readers {
		if e := sub.Read(s); e != nil {
			s.ToMarker(m)
			err = e
			break
		}
	}
	return err
}

func Seq(list ...Reader) Reader {
	return SeqReader{list}
}

//----------------------------------------------------------
type FoldReader struct {
	val string
}

func (r FoldReader) Read(s *Scanner) error {
	return s.boolError(s.IfFold(r.val))
}

func Fold(str string) Reader {
	return FoldReader{str}
}

//----------------------------------------------------------
type WSReader struct {
}

func (r WSReader) Read(s *Scanner) error {
	return s.boolError(s.IfAnyRune(" \r\n\t"))
}

func WS() Reader {
	return WSReader{}
}

//----------------------------------------------------------
type OptReader struct {
	sub Reader
}

func (r OptReader) Read(s *Scanner) error {
	r.sub.Read(s)
	return nil
}

func Opt(r Reader) Reader {
	return OptReader{r}
}

//----------------------------------------------------------
type ManyReader struct {
	sub Reader
}

func (r ManyReader) Read(s *Scanner) error {
	start := s.Mark()
	for r.sub.Read(s) == nil {
	}
	end := s.Mark()
	return s.boolError(start < end)
}

func Many(r Reader) Reader {
	return ManyReader{r}
}

//----------------------------------------------------------
type BetweenReader struct {
	min rune
	max rune
}

func (r BetweenReader) Read(s *Scanner) error {
	return s.boolError(s.IfBetween(r.min, r.max))
}

func Between(min rune, max rune) Reader {
	return BetweenReader{min, max}
}

//----------------------------------------------------------
type MatchReader struct {
	f MatchFunc
}

func (r MatchReader) Read(s *Scanner) error {
	return s.boolError(s.IfMatch(r.f))
}

func Match(f MatchFunc) Reader {
	return MatchReader{f}
}

//----------------------------------------------------------
type LitReader struct {
	str string
}

func (r LitReader) Read(s *Scanner) error {
	return s.boolError(s.If(r.str))
}

func Lit(str string) Reader {
	return LitReader{str}
}

//----------------------------------------------------------
type DigitReader struct {
}

func (r DigitReader) Read(s *Scanner) error {
	return s.boolError(s.IfBetween('0', '9'))
}

func Digit() Reader {
	return DigitReader{}
}

//----------------------------------------------------------
type HexDigitReader struct {
}

func (r HexDigitReader) Read(s *Scanner) error {
	res := s.IfBetween('0', '9') ||
		s.IfBetween('a', 'f') ||
		s.IfBetween('A', 'F')
	return s.boolError(res)
}

func HexDigit() Reader {
	return DigitReader{}
}

//----------------------------------------------------------
type ToReader struct {
	sub Reader
}

func (r ToReader) Read(s *Scanner) error {
	m := s.Mark()
	for ; !s.AtEnd(); s.move(1) {
		if e := r.sub.Read(s); e != nil {
			return nil
		}
	}
	s.ToMarker(m)
	return s.errorMsg()
}

func To(r Reader) Reader {
	return ToReader{r}
}
