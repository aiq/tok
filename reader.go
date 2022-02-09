package tok

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Type
type ReadFunc func(*Scanner) error

type Reader interface {
	Read(s *Scanner) error
	What() string
}

func (s *Scanner) Use(r Reader) error {
	m := s.Mark()
	err := r.Read(s)
	if err != nil {
		s.ToMarker(m)
	}
	return err
}

func (s *Scanner) UseFunc(f ReadFunc) error {
	m := s.Mark()
	err := f(s)
	if err != nil {
		s.ToMarker(m)
	}
	return err
}

func (s *Scanner) TraceUse(r Reader) (string, error) {
	m := s.Mark()
	err := r.Read(s)
	if err != nil {
		s.ToMarker(m)
	}
	return s.Since(m), err
}

//------------------------------------------------------------------------------
type anyReader struct {
	readers []Reader
}

func getDeepest(errs []error) error {
	var deepest ReadError
	for _, e := range errs {
		re, ok := e.(ReadError)
		if ok && re.DeeperAs(deepest) {
			deepest = re
		}
	}
	return deepest
}

func (r *anyReader) Read(s *Scanner) error {
	m := s.Mark()
	errs := []error{}
	for _, sub := range r.readers {
		if e := sub.Read(s); e == nil {
			return nil
		} else {
			errs = append(errs, e)
		}
	}
	s.ToMarker(m)
	e := getDeepest(errs)
	if e != nil {
		return e
	}
	return s.ErrorFor(r.What())
}

func (r *anyReader) What() string {
	sub := []string{}
	for _, sr := range r.readers {
		sub = append(sub, sr.What())
	}
	return "{ " + strings.Join(sub, " | ") + " }"
}

// Any creates a Reader that tries to Read with any of the given Reader.
// The first Reader that reads without an error will be used.
func Any(list ...Reader) Reader {
	return &anyReader{list}
}

func AnyString(list ...string) Reader {
	readers := []Reader{}
	for _, s := range list {
		readers = append(readers, Lit(s))
	}
	return &anyReader{readers}
}

//------------------------------------------------------------------------------
type anyRuneReader struct {
	str string
}

func (r *anyRuneReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.IfAnyRune(r.str), r.What())
}

func (r *anyRuneReader) What() string {
	return "{" + strconv.Quote(r.str) + "}"
}

// AnyRune
func AnyRune(str string) Reader {
	return &anyRuneReader{str}
}

//------------------------------------------------------------------------------
type betweenReader struct {
	min rune
	max rune
}

func (r betweenReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.IfBetween(r.min, r.max), r.What())
}

func (r betweenReader) What() string {
	return fmt.Sprintf("[%c-%c]", r.min, r.max)
}

// Between
func Between(min rune, max rune) Reader {
	return betweenReader{min, max}
}

//------------------------------------------------------------------------------
type betweenAnyReader struct {
	min []rune
	max []rune
}

func (r *betweenAnyReader) Read(s *Scanner) error {
	m := s.Mark()
	for i := 0; i < len(r.min); i++ {
		if s.IfBetween(r.min[i], r.max[i]) {
			return nil
		}
	}
	s.ToMarker(m)
	return s.ErrorFor(r.What())
}

func (r *betweenAnyReader) What() string {
	var b strings.Builder
	b.WriteRune('[')
	for i := 0; i < len(r.min); i++ {
		b.WriteRune(r.min[i])
		b.WriteRune('-')
		b.WriteRune(r.max[i])
	}
	b.WriteRune(']')
	return b.String()
}

// BetweenAny
func BetweenAny(str string) Reader {
	runes := []rune(str)
	if len(runes)%3 != 0 {
		return Any()
	}

	r := &betweenAnyReader{}
	for i := 0; i < len(runes); i += 3 {
		min := runes[i]
		sep := runes[i+1]
		max := runes[i+2]
		if sep != '-' {
			return Any()
		}
		r.min = append(r.min, min)
		r.max = append(r.max, max)
	}
	return r
}

// BuildBetweenAny
func BuildBetweenAny(minMax ...rune) Reader {
	if len(minMax)%2 != 0 {
		return Any()
	}

	r := &betweenAnyReader{}
	for i := 0; i < len(minMax); i += 2 {
		min := minMax[i]
		max := minMax[i+1]
		r.min = append(r.min, min)
		r.max = append(r.max, max)
	}
	return r
}

//------------------------------------------------------------------------------
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

// Bool
func Bool(format string) *BoolReader {
	return &BoolReader{
		Format: format,
	}
}

//------------------------------------------------------------------------------
// Digit creates a Reader that reads the digit runes '0'-'9'
func Digit() Reader {
	return Between('0', '9')
}

//------------------------------------------------------------------------------
type foldReader struct {
	val string
}

func (r foldReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.IfFold(r.val), r.What())
}

func (r foldReader) What() string {
	return "f" + r.val + ""
}

// Fold
func Fold(str string) Reader {
	return &foldReader{str}
}

//------------------------------------------------------------------------------
// HexDigit creates a Reader that reads the hex digit runes '0'-'9'/'a'-'f'/'A'-'F'.
func HexDigit() Reader {
	return BuildBetweenAny('0', '9', 'a', 'f', 'A', 'F')
}

//------------------------------------------------------------------------------
type holeyReader struct {
	min   rune
	max   rune
	holes string
}

func (r holeyReader) Read(s *Scanner) error {
	val, i := utf8.DecodeRuneInString(s.Tail())
	if inRange(r.min, val, r.max) && !strings.ContainsRune(r.holes, val) {
		return s.BoolErrorFor(s.Move(i), r.What())
	}
	return s.ErrorFor(r.What())
}

func (r holeyReader) What() string {
	return fmt.Sprintf("(u%d-u%d - %q)", r.min, r.max, r.holes)
}

// Holey
func Holey(min rune, max rune, holes string) Reader {
	return holeyReader{min, max, holes}
}

//------------------------------------------------------------------------------
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

//------------------------------------------------------------------------------
type litReader struct {
	str string
}

func (r litReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.If(r.str), r.What())
}

func (r litReader) What() string {
	return strconv.Quote(r.str)
}

// Lit
func Lit(str string) Reader {
	return litReader{str}
}

//------------------------------------------------------------------------------
type manyReader struct {
	sub Reader
}

func (r manyReader) Read(s *Scanner) error {
	start := s.Mark()
	for r.sub.Read(s) == nil {
	}
	end := s.Mark()
	return s.BoolErrorFor(start < end, r.What())
}

func (r manyReader) What() string {
	return "+" + r.sub.What()
}

// Many
func Many(r Reader) Reader {
	return &manyReader{r}
}

//------------------------------------------------------------------------------
type MapFunc func(Token)

type mapReader struct {
	sub Reader
	f   MapFunc
}

func (r *mapReader) Read(s *Scanner) error {
	t, err := s.TokenizeUse(r.sub)
	if err == nil {
		r.f(t)
	}
	return err
}

func (r *mapReader) What() string {
	return "map(" + r.sub.What() + ")"
}

// Map
func Map(r Reader, f MapFunc) Reader {
	return &mapReader{r, f}
}

//------------------------------------------------------------------------------
type matchReader struct {
	what string
	f    MatchFunc
}

func (r matchReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.IfMatch(r.f), r.what)
}

func (r matchReader) What() string {
	return r.what
}

// Match
func Match(what string, f MatchFunc) Reader {
	return matchReader{what, f}
}

//------------------------------------------------------------------------------
type namedReader struct {
	name string
	sub  Reader
}

func (r *namedReader) Read(s *Scanner) error {
	return r.sub.Read(s)
}

func (r *namedReader) What() string {
	return r.name
}

func Named(name string, r Reader) Reader {
	return &namedReader{name, r}
}

//------------------------------------------------------------------------------
type optReader struct {
	sub Reader
}

func (r *optReader) Read(s *Scanner) error {
	r.sub.Read(s)
	return nil
}

func (r *optReader) What() string {
	return "?" + r.sub.What()
}

// Opt
func Opt(r Reader) Reader {
	return &optReader{r}
}

//------------------------------------------------------------------------------
type runeReader struct {
	r rune
}

func (r runeReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.IfRune(r.r), r.What())
}

func (r runeReader) What() string {
	return strconv.QuoteRune(r.r)
}

// Rune
func Rune(r rune) Reader {
	return runeReader{r}
}

//------------------------------------------------------------------------------
type seqReader struct {
	readers []Reader
}

func (r *seqReader) Read(s *Scanner) error {
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

func (r *seqReader) What() string {
	sub := []string{}
	for _, sr := range r.readers {
		sub = append(sub, sr.What())
	}
	return "( " + strings.Join(sub, " > ") + " )"
}

func Seq(list ...Reader) Reader {
	return &seqReader{list}
}

//------------------------------------------------------------------------------
type timesReader struct {
	n   int
	sub Reader
}

func (r *timesReader) Read(s *Scanner) error {
	for i := 0; i < r.n; i++ {
		if e := r.sub.Read(s); e != nil {
			return e
		}
	}
	return nil
}

func (r *timesReader) What() string {
	return fmt.Sprintf("%d*%s", r.n, r.sub.What())
}

// Times
func Times(n int, r Reader) Reader {
	return &timesReader{n, r}
}

//------------------------------------------------------------------------------
type toReader struct {
	sub Reader
}

func (r *toReader) Read(s *Scanner) error {
	m := s.Mark()
	for ; !s.AtEnd(); s.Move(1) {
		if e := r.sub.Read(s); e != nil {
			return nil
		}
	}
	s.ToMarker(m)
	return s.ErrorFor(r.What())
}

func (r *toReader) What() string {
	return "->" + r.sub.What()
}

// To
func To(r Reader) Reader {
	return &toReader{r}
}

//------------------------------------------------------------------------------
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

// Uint
func Uint(base int, bitSize int) *UintReader {
	return &UintReader{
		Base:    base,
		BitSize: bitSize,
	}
}

//------------------------------------------------------------------------------
type wrapReader struct {
	what string
	f    ReadFunc
}

func (r wrapReader) Read(s *Scanner) error {
	return r.f(s)
}

func (r wrapReader) What() string {
	return r.what
}

// Wrap
func Wrap(what string, f ReadFunc) Reader {
	return wrapReader{what, f}
}

//------------------------------------------------------------------------------
// WS
func WS() Reader {
	return AnyRune(" \r\n\t")
}

//------------------------------------------------------------------------------
type zomReader struct {
	sub Reader
}

func (r zomReader) Read(s *Scanner) error {
	for r.sub.Read(s) == nil {
	}
	return nil
}

func (r zomReader) What() string {
	return "*" + r.sub.What()
}

// Zom
func Zom(r Reader) Reader {
	return &zomReader{r}
}
