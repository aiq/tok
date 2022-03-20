package tok

import (
	"strings"
)

// Tracker
type Tracker interface {
	Update(m Marker)
}

// Returns a new empty Basket that is set as Tracker on the scanner.
func (s *Scanner) NewBasket() *Basket {
	b := &Basket{}
	s.Tracker = b
	return b
}

func (s *Scanner) NewBasketFor(g Grammar) *Basket {
	b := &Basket{}
	b.PickWith(g.Grammar()...)
	s.Tracker = b
	return b
}

//------------------------------------------------------------------------------

// Basket can be used to Pick readed Segments.
type Basket struct {
	segments []Segment
}

func (b *Basket) Add(seg Segment) {
	b.segments = append(b.segments, seg)
}

func (b *Basket) Update(m Marker) {
	for i := len(b.segments); i > 0; i-- {
		seg := b.segments[i-1]
		if seg.to <= m {
			b.segments = b.segments[:i]
			return
		}
	}
	b.segments = []Segment{}
}

func (b *Basket) Picked() []Segment {
	return b.segments
}

func (b *Basket) String() string {
	segs := []string{}
	for _, seg := range b.segments {
		segs = append(segs, seg.String())
	}
	return strings.Join(segs, ";")
}

// PickWith calls Pick on all rules with the Basket as paramter.
func (b *Basket) PickWith(rules ...*Rule) {
	for _, r := range rules {
		r.Pick(b)
	}
}

//------------------------------------------------------------------------------
type pickReader struct {
	info   string
	basket *Basket
	sub    Reader
}

func (r *pickReader) Read(s *Scanner) error {
	t, err := s.TokenizeUse(r.sub)
	if err == nil {
		r.basket.Add(Segment{
			Info:  r.info,
			Token: t,
		})
	}
	return err
}

func (r *pickReader) What() string {
	return r.sub.What()
}

// Pick creates a Reader that appends the Segments that r reads forward to the Basket with info as Info value.
func Pick(r Reader, b *Basket, info string) Reader {
	return &pickReader{
		info:   info,
		basket: b,
		sub:    r,
	}
}
