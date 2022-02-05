package grammer

import (
	"fmt"
	"unicode/utf8"

	. "github.com/aiq/tok"
)

type JSONReader struct {
	Value      RefReader
	Object     RefReader
	Members    RefReader
	Member     RefReader
	Array      RefReader
	Elements   RefReader
	Element    RefReader
	String     RefReader
	Characters RefReader
	Character  RefReader
	Escape     RefReader
	Hex        RefReader
	Number     RefReader
	Integer    RefReader
	Fraction   RefReader
	Exponent   RefReader
	OneNine    RefReader
	Digit      RefReader
	Digits     RefReader
	Sign       RefReader
	WS         RefReader
}

func JSON() *JSONReader {
	g := &JSONReader{
		Value:      Ref("value"),
		Object:     Ref("object"),
		Members:    Ref("members"),
		Member:     Ref("member"),
		Array:      Ref("array"),
		Elements:   Ref("elements"),
		Element:    Ref("element"),
		String:     Ref("string"),
		Characters: Ref("characters"),
		Character:  Ref("character"),
		Escape:     Ref("escape"),
		Hex:        Ref("hex"),
		Number:     Ref("number"),
		Integer:    Ref("integer"),
		Fraction:   Ref("fraction"),
		Exponent:   Ref("exponent"),
		OneNine:    Ref("onenine"),
		Digit:      Ref("digit"),
		Digits:     Ref("digits"),
		Sign:       Ref("sign"),
		WS:         Ref("ws"),
	}
	g.WS.Sub = Opt(Many(WS()))
	g.Sign.Sub = Opt(AnyRune("+-"))
	g.OneNine.Sub = Between('1', '9')
	g.Digit.Sub = Digit()
	g.Digits.Sub = Many(Digit())
	g.Exponent.Sub = Opt(Seq(AnyRune("eE"), &g.Sign, &g.Digits))
	g.Fraction.Sub = Opt(Seq(Rune('.'), &g.Digits))
	g.Integer.Sub = Seq(Opt(Rune('-')), Any(Rune('0'), Seq(&g.OneNine, Opt(&g.Digits))))
	g.Number.Sub = Seq(&g.Integer, &g.Fraction, &g.Exponent)
	g.Hex.Sub = HexDigit()
	g.Escape.Sub = Any(AnyRune(`"\/bfnrt`), Seq(Rune('u'), Times(4, &g.Hex)))
	g.Character.Sub = Any(Holey(' ', utf8.MaxRune, `"\`), Seq(Rune('\\'), &g.Escape))
	g.Characters.Sub = Opt(Many(&g.Character))
	g.String.Sub = Seq(Rune('"'), &g.Characters, Rune('"'))
	g.Element.Sub = Seq(&g.WS, &g.Value, &g.WS)
	g.Elements.Sub = Seq(&g.Element, Opt(Many(Seq(Rune(','), &g.Element))))
	g.Array.Sub = Seq(Rune('['), Any(&g.Elements, &g.WS), Rune(']'))
	g.Member.Sub = Seq(&g.WS, &g.String, &g.WS, Rune(':'), &g.Element)
	g.Members.Sub = Seq(&g.Member, Opt(Many(Seq(Rune(','), &g.Member))))
	g.Object.Sub = Seq(Rune('{'), Any(&g.Members, &g.WS), Rune('}'))
	g.Value.Sub = Any(&g.Object, &g.Array, &g.String, &g.Number, Lit("true"), Lit("false"), Lit("null"))
	return g
}

func (r *JSONReader) Read(s *Scanner) error {
	err := r.Element.Read(s)
	if err != nil {
		l, c := s.LineCol()
		return fmt.Errorf("json parse error at %d:%d: %v", l, c, err)
	}
	return nil
}

func (r *JSONReader) What() string {
	return "json"
}

func (r *JSONReader) Grammer() Rules {
	rules := []Rule{}
	rules = append(rules, r.WS.Rule())
	rules = append(rules, r.Sign.Rule())
	rules = append(rules, r.OneNine.Rule())
	rules = append(rules, r.Digit.Rule())
	rules = append(rules, r.Digits.Rule())
	rules = append(rules, r.Exponent.Rule())
	rules = append(rules, r.Fraction.Rule())
	rules = append(rules, r.Integer.Rule())
	rules = append(rules, r.Number.Rule())
	rules = append(rules, r.Hex.Rule())
	rules = append(rules, r.Escape.Rule())
	rules = append(rules, r.Character.Rule())
	rules = append(rules, r.Characters.Rule())
	rules = append(rules, r.String.Rule())
	rules = append(rules, r.Element.Rule())
	rules = append(rules, r.Elements.Rule())
	rules = append(rules, r.Array.Rule())
	rules = append(rules, r.Member.Rule())
	rules = append(rules, r.Members.Rule())
	rules = append(rules, r.Object.Rule())
	rules = append(rules, r.Value.Rule())
	return rules
}
