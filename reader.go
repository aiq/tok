package tok

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ReadFunc represents the prototype of a read function.
type ReadFunc func(*Scanner) error

// Reader can be used by the scanner to read from the scanner.
type Reader interface {
	Read(s *Scanner) error
	What() string
}

// Use uses r on the scanner.
// The scanner is only moved if no error occurs.
func (s *Scanner) Use(r Reader) error {
	m := s.Mark()
	err := r.Read(s)
	if err != nil {
		s.ToMarker(m)
	}
	return err
}

// UseFunc uses f on the scanner.
// The scanner is only moved if no error occurs.
func (s *Scanner) UseFunc(f ReadFunc) error {
	m := s.Mark()
	err := f(s)
	if err != nil {
		s.ToMarker(m)
	}
	return err
}

// TraceUse traces the readed sub string.
func (s *Scanner) TraceUse(r Reader) (string, error) {
	m := s.Mark()
	err := r.Read(s)
	if err != nil {
		s.ToMarker(m)
	}
	return s.Since(m), err
}

// TraceUseFunc traces the via f traced sub string.
func (s *Scanner) TraceUseFunc(f ReadFunc) (string, error) {
	m := s.Mark()
	err := f(s)
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
		if ok && re.Later(deepest) {
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
	return "[ " + strings.Join(sub, " ") + " ]"
}

// Any creates a Reader that tries to Read with any of the given Reader.
// The first Reader that reads without an error will be used.
func Any(list ...Reader) Reader {
	return &anyReader{list}
}

// AnyFold creates a Reader that tries to Read any of the strings in list.
// The first string that matches under Unicode case-folding will be used.
func AnyFold(list ...string) Reader {
	readers := []Reader{}
	for _, s := range list {
		readers = append(readers, Fold(s))
	}
	return &anyReader{readers}
}

// AnyLit creates a Reader that tries to Read any of the strings in list.
func AnyLit(list ...string) Reader {
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
	return "[" + strconv.QuoteToGraphic(r.str) + "]"
}

// AnyRune creates a Reader that tries to Read any of the runes in list.
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

func quoteRune(r rune) string {
	return strings.Trim(strconv.QuoteRuneToGraphic(r), "'")
}

func quoteBetween(min, max rune) string {
	return fmt.Sprintf("<%s%s>", quoteRune(min), quoteRune(max))
}

func (r betweenReader) What() string {
	return quoteBetween(r.min, r.max)
}

// Between creates a Reader that tries to Read a rune that is >= min and <= max.
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
	b.WriteString("[<")
	for i := 0; i < len(r.min); i++ {
		b.WriteRune(' ')
		b.WriteString(quoteRune(r.min[i]))
		b.WriteString(quoteRune(r.max[i]))
	}
	b.WriteString(" >]")
	return b.String()
}

// BetweenAny creates a Reader that tries to Read a rune that is between any of the ranges that str descibes.
// A range can look like "a-z", multible ranges can look like "a-zA-Z0-9".
// Invalid str values that can't be interpreted lead to an invalid reader.
func BetweenAny(str string) Reader {
	runes := []rune(str)
	if len(runes)%3 != 0 {
		return InvalidReader("invalid init string for BetweenAny: %q", str)
	}

	r := &betweenAnyReader{}
	for i := 0; i < len(runes); i += 3 {
		min := runes[i]
		sep := runes[i+1]
		max := runes[i+2]
		if sep != '-' {
			return InvalidReader("invalid init string for BetweenAny: %q", str)
		}
		r.min = append(r.min, min)
		r.max = append(r.max, max)
	}
	return r
}

// BetweenAny creates a Reader that tries to Read a rune that is between any of the ranges that str descibes.
// The number minMax arguments must be even, an invalid number of minMax values lead to an invalid reader.
func BuildBetweenAny(minMax ...rune) Reader {
	if len(minMax)%2 != 0 {
		return InvalidReader("odd number of min-max values for BetweenAny: %d", len(minMax))
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
// BoolReader is a Reader that stores the readed bool value in the field Value.
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
	return "bool{" + strconv.QuoteToGraphic(r.Format) + "}"
}

// Bool creates a Reader to Read bool values from the scanner.
// Valid format values are
// - "l" for true and false
// - "U" for TRUE and FALSE
// - "Cc" for True and False
// - "*" for all cases
// An empty format string will be interpreted as "*".
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
	return "~" + strconv.QuoteToGraphic(r.val)
}

// AnyFold creates a Reader that tries to Read a string under Unicode case-folding.
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
	return fmt.Sprintf("(%s - %s)", quoteBetween(r.min, r.max), strconv.QuoteToGraphic(r.holes))
}

// Holey creates a Reader that tries to Read a rune that is >= min and <= max without the runes in holes.
func Holey(min rune, max rune, holes string) Reader {
	return holeyReader{min, max, holes}
}

//------------------------------------------------------------------------------
// IntReader is a Reader that stores the readed int value in the field Value.
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
	return fmt.Sprintf("int{%d,%d}", r.Base, r.BitSize)
}

// Int creates a Reader to Read int values from the scanner.
// Valid base values are 8. 10 and 16.
// Valid bitSize values are 8, 16, 32 and 64.
func Int(base int, bitSize int) *IntReader {
	return &IntReader{
		Base:    base,
		BitSize: bitSize,
	}
}

//------------------------------------------------------------------------------
type janusBeginReader struct {
	reader Reader
	end    *janusEndReader
	name   string
}

func (r *janusBeginReader) Read(s *Scanner) error {
	t, err := s.TokenizeUse(r.reader)
	if err == nil {
		r.end.reader.str = s.Get(t)
	}
	return err
}

func (r *janusBeginReader) What() string {
	return "!" + r.name + "<" + r.reader.What()
}

type janusEndReader struct {
	reader litReader
	name   string
}

func (r *janusEndReader) Read(s *Scanner) error {
	return r.reader.Read(s)
}

func (r *janusEndReader) What() string {
	return "!" + r.name + ">"
}

// Janus creates two Reader.
// The first one tries to match with r.
// If the first matches expects the second the matched sub string.
func Janus(name string, r Reader) (Reader, Reader) {
	end := &janusEndReader{
		reader: litReader{""},
		name:   name,
	}
	beg := &janusBeginReader{
		reader: r,
		end:    end,
		name:   name,
	}
	return beg, end
}

//------------------------------------------------------------------------------
type invalidReader struct {
	err error
}

func (r invalidReader) Read(s *Scanner) error {
	return fmt.Errorf("invalid reader: %v", r.err)
}

func (r invalidReader) What() string {
	return fmt.Sprintf("(invalid reader: %v)", r.err)
}

// InvalidReader is a Reader that allways fails.
// The arguments will be passed to fmt.Errorf.
func InvalidReader(format string, a ...interface{}) invalidReader {
	return invalidReader{fmt.Errorf(format, a...)}
}

//------------------------------------------------------------------------------
type litReader struct {
	str string
}

func (r litReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.If(r.str), r.What())
}

func (r litReader) What() string {
	return strconv.QuoteToGraphic(r.str)
}

// Lit creates a Reader that tries to read the string str.
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

// Many creates a Reader that expects that r matches one or more times.
// See Zom for a Reader that expects zero or more.
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
	return r.sub.What()
}

// Map creates a Reader that passed the Token that r reads forward to f.
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

// Match creates a Reader that tries to read a rune that matches by f.
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

// Named creates a Reader with a custom name that the function What returns.
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

// Opt creates Reader that catches the error that r can produce and returns allways nil
func Opt(r Reader) Reader {
	return &optReader{r}
}

//------------------------------------------------------------------------------
type pastReader struct {
	sub Reader
}

func (r *pastReader) Read(s *Scanner) error {
	m := s.Mark()
	for ; !s.AtEnd(); s.Move(1) {
		if e := r.sub.Read(s); e == nil {
			return nil
		}
	}
	s.ToMarker(m)
	return s.ErrorFor(r.What())
}

func (r *pastReader) What() string {
	return "-->" + r.sub.What()
}

// Past creates a Reader that reads until r matches, with the matched part.
func Past(r Reader) Reader {
	return &pastReader{r}
}

//------------------------------------------------------------------------------
type runeReader struct {
	r rune
}

func (r runeReader) Read(s *Scanner) error {
	return s.BoolErrorFor(s.IfRune(r.r), r.What())
}

func (r runeReader) What() string {
	return strconv.QuoteRuneToGraphic(r.r)
}

// Rune creates a Reader that tries to read r.
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
	return "> " + strings.Join(sub, " ") + " >"
}

// Seq creates a Reader that tries to Read with all readers in list sequential.
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

// Times creates a Reader that tries to Read n times with r.
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
		subM := s.Mark()
		if e := r.sub.Read(s); e == nil {
			s.ToMarker(subM)
			return nil
		}
	}
	s.ToMarker(m)
	return s.ErrorFor(r.What())
}

func (r *toReader) What() string {
	return "->" + r.sub.What()
}

// To creates a Reader that reads until r matches, without the matched part.
func To(r Reader) Reader {
	return &toReader{r}
}

//------------------------------------------------------------------------------
// UintReader is a Reader that stores the readed uint value in the field Value.
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
	return fmt.Sprintf("uint{%d,%d}", r.Base, r.BitSize)
}

// Uint creates a Reader to Read uint values from the scanner.
// Valid base values are 8, 10 and 16.
// Valid bitSize values are 8, 16, 32 and 64.
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

// Wrap creates a Reader that wraps f.
func Wrap(what string, f ReadFunc) Reader {
	return wrapReader{what, f}
}

//------------------------------------------------------------------------------
// WS creates a Reader to read one whitespace character(" \r\n\t").
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

// Zom creates a Reader that expects that r matches zero or more times.
// See Many for a Reader that expects one or more.
func Zom(r Reader) Reader {
	return &zomReader{r}
}
