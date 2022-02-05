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

func (s *Scanner) TracedUse(r Reader) (string, error) {
	m := s.Mark()
	err := r.Read(s)
	if err != nil {
		s.ToMarker(m)
	}
	return s.Since(m), err
}

func (s *Scanner) TokenizeUse(r Reader) (Token, error) {
	a := s.Mark()
	err := r.Read(s)
	if err != nil {
		s.ToMarker(a)
	}
	b := s.Mark()
	return Token{a, b}, err
}

//----------------------------------------------------------
// AnyReader
type AnyReader struct {
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

func (r *AnyReader) Read(s *Scanner) error {
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

func (r *AnyReader) What() string {
	sub := []string{}
	for _, sr := range r.readers {
		sub = append(sub, sr.What())
	}
	return "{ " + strings.Join(sub, " | ") + " }"
}

func Any(list ...Reader) Reader {
	return &AnyReader{list}
}

func AnyString(list ...string) Reader {
	readers := []Reader{}
	for _, s := range list {
		readers = append(readers, Lit(s))
	}
	return &AnyReader{readers}
}

//----------------------------------------------------------
type AnyRuneReader struct {
	str string
}

func (r *AnyRuneReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.IfAnyRune(r.str), r.What())
}

func (r *AnyRuneReader) What() string {
	return "{" + strconv.Quote(r.str) + "}"
}

func AnyRune(str string) Reader {
	return &AnyRuneReader{str}
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
type BetweenAnyReader struct {
	min []rune
	max []rune
}

func (r BetweenAnyReader) Read(s *Scanner) error {
	m := s.Mark()
	for i := 0; i < len(r.min); i++ {
		if s.IfBetween(r.min[i], r.max[i]) {
			return nil
		}
	}
	s.ToMarker(m)
	return s.ErrorFor(r.What())
}

func (r BetweenAnyReader) What() string {
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

func BetweenAny(str string) Reader {
	runes := []rune(str)
	if len(runes)%3 != 0 {
		return Any()
	}

	r := &BetweenAnyReader{}
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
type CollectReader struct {
	Bag []Token
	sub Reader
}

func (r *CollectReader) Read(s *Scanner) error {
	t, err := s.TokenizeUse(r.sub)
	if err == nil {
		r.Bag = append(r.Bag, t)
	}
	return err
}

func (r *CollectReader) What() string {
	return "collect(" + r.sub.What() + ")"
}

func Collect(r Reader) *CollectReader {
	return &CollectReader{
		sub: r,
	}
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
type FoldReader struct {
	val string
}

func (r FoldReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.IfFold(r.val), r.What())
}

func (r FoldReader) What() string {
	return "f" + r.val + ""
}

func Fold(str string) Reader {
	return &FoldReader{str}
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
	return HexDigitReader{}
}

//----------------------------------------------------------
type HoleyReader struct {
	min   rune
	max   rune
	holes string
}

func (r HoleyReader) Read(s *Scanner) error {
	val, i := utf8.DecodeRuneInString(s.Tail())
	if inRange(r.min, val, r.max) && !strings.ContainsRune(r.holes, val) {
		return s.BoolErrorFor(s.Move(i), r.What())
	}
	return s.ErrorFor(r.What())
}

func (r HoleyReader) What() string {
	return fmt.Sprintf("(u%d-u%d - %q)", r.min, r.max, r.holes)
}

func Holey(min rune, max rune, holes string) Reader {
	return HoleyReader{min, max, holes}
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
type LitReader struct {
	str string
}

func (r LitReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.If(r.str), r.What())
}

func (r LitReader) What() string {
	return strconv.Quote(r.str)
}

func Lit(str string) Reader {
	return LitReader{str}
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
	return &ManyReader{r}
}

//----------------------------------------------------------
type MapFunc func(Token)

type MapReader struct {
	sub Reader
	f   MapFunc
}

func (r MapReader) Read(s *Scanner) error {
	t, err := s.TokenizeUse(r.sub)
	if err == nil {
		r.f(t)
	}
	return err
}

func (r MapReader) What() string {
	return "map(" + r.sub.What() + ")"
}

func Map(r Reader, f MapFunc) MapReader {
	return MapReader{r, f}
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
type NamedReader struct {
	Name string
	Sub  Reader
}

func (r *NamedReader) Read(s *Scanner) error {
	return r.Sub.Read(s)
}

func (r *NamedReader) What() string {
	return r.Name
}

func Named(name string, r Reader) *NamedReader {
	return &NamedReader{name, r}
}

//----------------------------------------------------------
type OptReader struct {
	sub Reader
}

func (r *OptReader) Read(s *Scanner) error {
	r.sub.Read(s)
	return nil
}

func (r *OptReader) What() string {
	return "?" + r.sub.What()
}

func Opt(r Reader) Reader {
	return &OptReader{r}
}

//----------------------------------------------------------
type RuneReader struct {
	r rune
}

func (r RuneReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.IfRune(r.r), r.What())
}

func (r RuneReader) What() string {
	return strconv.QuoteRune(r.r)
}

func Rune(r rune) Reader {
	return RuneReader{r}
}

//----------------------------------------------------------
type SeqReader struct {
	readers []Reader
}

func (r *SeqReader) Read(s *Scanner) error {
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

func (r *SeqReader) What() string {
	sub := []string{}
	for _, sr := range r.readers {
		sub = append(sub, sr.What())
	}
	return "( " + strings.Join(sub, " > ") + " )"
}

func Seq(list ...Reader) Reader {
	return &SeqReader{list}
}

//----------------------------------------------------------
type TimesReader struct {
	N   int
	sub Reader
}

func (r *TimesReader) Read(s *Scanner) error {
	for i := 0; i < r.N; i++ {
		if e := r.sub.Read(s); e != nil {
			return e
		}
	}
	return nil
}

func (r *TimesReader) What() string {
	return fmt.Sprintf("%d*%s", r.N, r.sub.What())
}

func Times(n int, r Reader) *TimesReader {
	return &TimesReader{n, r}
}

//----------------------------------------------------------
type ToReader struct {
	sub Reader
}

func (r *ToReader) Read(s *Scanner) error {
	m := s.Mark()
	for ; !s.AtEnd(); s.Move(1) {
		if e := r.sub.Read(s); e != nil {
			return nil
		}
	}
	s.ToMarker(m)
	return s.ErrorFor(r.What())
}

func (r *ToReader) What() string {
	return "->" + r.sub.What()
}

func To(r Reader) Reader {
	return &ToReader{r}
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

func WS() Reader {
	return AnyRune(" \r\n\t")
}
