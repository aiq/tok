package tok

import (
	"fmt"
	"strings"
)

type ReadFunc func(*Scanner) error

type Reader interface {
	Read(s *Scanner) error
	What() string
}

func (s *Scanner) Use(r Reader) error {
	return r.Read(s)
}

func (s *Scanner) UseFunc(f ReadFunc) error {
	return f(s)
}

func (s *Scanner) TracedUse(r Reader) (string, error) {
	m := s.Mark()
	err := r.Read(s)
	return s.Since(m), err
}

func (s *Scanner) TokenizeUse(r Reader) (Token, error) {
	a := s.Mark()
	err := r.Read(s)
	b := s.Mark()
	return Token{a, b}, err
}

//----------------------------------------------------------
// AnyReader
type AnyReader struct {
	readers []Reader
}

func (r AnyReader) Read(s *Scanner) error {
	m := s.Mark()
	for _, sub := range r.readers {
		if e := sub.Read(s); e == nil {
			return nil
		}
	}
	s.ToMarker(m)
	return s.ErrorFor(r.What())
}

func (r AnyReader) What() string {
	sub := []string{}
	for _, sr := range r.readers {
		sub = append(sub, sr.What())
	}
	return "any{ " + strings.Join(sub, " | ") + " }"
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

func (r SeqReader) What() string {
	sub := []string{}
	for _, sr := range r.readers {
		sub = append(sub, sr.What())
	}
	return "seq{ " + strings.Join(sub, " > ") + " }"
}

func Seq(list ...Reader) Reader {
	return SeqReader{list}
}

//----------------------------------------------------------
type FoldReader struct {
	val string
}

func (r FoldReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.IfFold(r.val), r.What())
}

func (r FoldReader) What() string {
	return "fold(" + r.val + ")"
}

func Fold(str string) Reader {
	return FoldReader{str}
}

//----------------------------------------------------------
type WSReader struct {
}

func (r WSReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.IfAnyRune(" \r\n\t"), r.What())
}

func (r WSReader) What() string {
	return "WS"
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

func (r OptReader) What() string {
	return "?" + r.sub.What()
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
	return s.BoolErrorFor(start < end, r.What())
}

func (r ManyReader) What() string {
	return "+" + r.sub.What()
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
	return s.BoolErrorFor(s.IfBetween(r.min, r.max), r.What())
}

func (r BetweenReader) What() string {
	return fmt.Sprintf("[%c-%c]", r.min, r.max)
}

func Between(min rune, max rune) Reader {
	return BetweenReader{min, max}
}

//----------------------------------------------------------
type MatchReader struct {
	what string
	f    MatchFunc
}

func (r MatchReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.IfMatch(r.f), r.what)
}

func (r MatchReader) What() string {
	return r.what
}

func Match(what string, f MatchFunc) Reader {
	return MatchReader{what, f}
}

//----------------------------------------------------------
type LitReader struct {
	str string
}

func (r LitReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.If(r.str), r.What())
}

func (r LitReader) What() string {
	return `"` + r.str + `"`
}

func Lit(str string) Reader {
	return LitReader{str}
}

//----------------------------------------------------------
type DigitReader struct {
}

func (r DigitReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.IfBetween('0', '9'), r.What())
}

func (r DigitReader) What() string {
	return "[0-9]"
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
	return s.BoolErrorFor(res, r.What())
}

func (r HexDigitReader) What() string {
	return "[0-9a-fA-F]"
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
	for ; !s.AtEnd(); s.Move(1) {
		if e := r.sub.Read(s); e != nil {
			return nil
		}
	}
	s.ToMarker(m)
	return s.ErrorFor(r.What())
}

func (r ToReader) What() string {
	return "->" + r.sub.What()
}

func To(r Reader) Reader {
	return ToReader{r}
}

//----------------------------------------------------------
type CollectReader struct {
	Bag []Token
	Reader
}

func (r *CollectReader) Read(s *Scanner) error {
	t, err := s.TokenizeUse(r.Reader)
	if err == nil {
		r.Bag = append(r.Bag, t)
	}
	return err
}

func Collect(r Reader) Reader {
	return &CollectReader{
		Reader: r,
	}
}

//----------------------------------------------------------
type WrapReader struct {
	what string
	f    ReadFunc
}

func (r WrapReader) Read(s *Scanner) error {
	return r.f(s)
}

func (r WrapReader) What() string {
	return r.what
}

func Wrap(what string, f ReadFunc) Reader {
	return WrapReader{what, f}
}

//----------------------------------------------------------
type SetReader struct {
	sub Reader
}

func (r SetReader) Read(s *Scanner) error {
	return nil
}

func (r SetReader) What() string {
	return ""
}

func Set(str string) Reader {

	return SetReader{}
}

//----------------------------------------------------------
type BoolReader struct {
	Value  bool
	Format string
}

func (r *BoolReader) Read(s *Scanner) error {
	v, err := s.ReadBool(r.Format)
	r.Value = v
	return err
}

func (r *BoolReader) What() string {
	if r.Format == "" {
		return "bool"
	}
	return "bool(" + r.Format + ")"
}

func Bool(format string) *BoolReader {
	return &BoolReader{
		Format: format,
	}
}

//----------------------------------------------------------
type IntReader struct {
	Value   int64
	Base    int
	BitSize int
}

func (r *IntReader) Read(s *Scanner) error {
	v, err := s.ReadInt(r.Base, r.BitSize)
	r.Value = v
	return err
}

func (r *IntReader) What() string {
	return fmt.Sprintf("int%d", r.BitSize)
}

func Int(base int, bitSize int) *IntReader {
	return &IntReader{
		Base:    base,
		BitSize: bitSize,
	}
}

//----------------------------------------------------------
type UintReader struct {
	Value   uint64
	Base    int
	BitSize int
}

func (r *UintReader) Read(s *Scanner) error {
	v, err := s.ReadUint(r.Base, r.BitSize)
	r.Value = v
	return err
}

func (r *UintReader) What() string {
	return fmt.Sprintf("uint%d", r.BitSize)
}

func Uint(base int, bitSize int) *UintReader {
	return &UintReader{
		Base:    base,
		BitSize: bitSize,
	}
}
