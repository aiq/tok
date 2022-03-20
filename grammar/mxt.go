package grammar

import (
	"fmt"
	"unicode/utf8"

	. "github.com/aiq/tok"
)

type MXTReader struct {
	Chunks       Rule `name:"chunks"`
	Chunk        Rule `name:"chunk"`
	Header       Rule `name:"header"`
	Marker       Rule `name:"marker"`
	NextMarker   Rule `name:"next-marker"`
	Name         Rule `name:"name"`
	Comment      Rule `name:"comment"`
	Arrow        Rule `name:"arrow"`
	Salt         Rule `name:"salt"`
	EmptyContent Rule `name:"empty-content"`
	Content      Rule `name:"content"`
	Word         Rule `name:"word"`
	WordChar     Rule `name:"wordchar"`
	NL           Rule `name:"nl"`
}

// MXT creates a Grammar to Read a MXT file.
// The implementation is based on https://mxt.aiq.dk/
func MXT() *MXTReader {
	g := &MXTReader{}
	SetRuleNames(g)
	g.NL.Reader = NL()
	g.WordChar.Reader = Between(0x21, utf8.MaxRune)
	g.Word.Reader = Many(&g.WordChar)
	g.Marker.Reader = Seq(Lit("//"), Zom(&g.WordChar))
	g.Name.Reader = &g.Word
	arrow := BodyTail(Zom(&g.WordChar), Lit("-->"))
	g.Comment.Reader = To(Named("arrow", arrow))
	g.Arrow.Reader = arrow
	saltHead, saltTail := Janus("", Opt(&g.Word))
	g.Salt.Reader = Seq(Many(' '), saltHead, Zom(' '))
	g.Header.Reader = Seq(&g.Marker, Many(' '), &g.Name, &g.Comment, &g.Arrow, Opt(&g.Salt))
	g.NextMarker.Reader = Seq(&g.NL, "//", saltTail)
	g.EmptyContent.Reader = Any(AtEnd(), At(&g.NextMarker))
	g.Content.Reader = To(Any(&g.NextMarker, AtEnd()))
	g.Chunk.Reader = Seq(&g.Header, Any(&g.EmptyContent, Seq(&g.NL, &g.Content)))
	g.Chunks.Reader = Seq(&g.Chunk, Zom(Seq(&g.NL, &g.Chunk)))
	return g
}

func (r *MXTReader) Read(s *Scanner) error {
	err := r.Chunks.Read(s)
	if err != nil {
		return fmt.Errorf("mxt parse error: %v", err)
	}
	return nil
}

func (r *MXTReader) What() string {
	return "mxt"
}

func (r *MXTReader) Grammar() []*Rule {
	return CollectRules(r)
}
