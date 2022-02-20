package grammar

import (
	"fmt"
	"unicode/utf8"

	. "github.com/aiq/tok"
)

type MXTReader struct {
	Chunks       RuleReader `name:"chunks"`
	Chunk        RuleReader `name:"chunk"`
	Header       RuleReader `name:"header"`
	Marker       RuleReader `name:"marker"`
	NextMarker   RuleReader `name:"next-marker"`
	Name         RuleReader `name:"name"`
	Comment      RuleReader `name:"comment"`
	Arrow        RuleReader `name:"arrow"`
	Salt         RuleReader `name:"salt"`
	EmptyContent RuleReader `name:"empty-content"`
	Content      RuleReader `name:"content"`
	Word         RuleReader `name:"word"`
	WordChar     RuleReader `name:"wordchar"`
	NL           RuleReader `name:"nl"`
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

func (r *MXTReader) Grammar() Rules {
	return CollectRules(r)
}
