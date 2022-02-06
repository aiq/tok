package grammer

import (
	"fmt"
	"unicode/utf8"

	. "github.com/aiq/tok"
)

type JSONReader struct {
	Value      RuleReader `name:"value"`
	Object     RuleReader `name:"object"`
	Members    RuleReader `name:"members"`
	Member     RuleReader `name:"member"`
	Array      RuleReader `name:"array"`
	Elements   RuleReader `name:"elements"`
	Element    RuleReader `name:"element"`
	String     RuleReader `name:"string"`
	Characters RuleReader `name:"characters"`
	Character  RuleReader `name:"character"`
	Escape     RuleReader `name:"escape"`
	Hex        RuleReader `name:"hex"`
	Number     RuleReader `name:"number"`
	Integer    RuleReader `name:"integer"`
	Fraction   RuleReader `name:"fraction"`
	Exponent   RuleReader `name:"exponent"`
	OneNine    RuleReader `name:"onenine"`
	Digit      RuleReader `name:"digit"`
	Digits     RuleReader `name:"digits"`
	Sign       RuleReader `name:"sign"`
	WS         RuleReader `name:"ws"`
}

// based on https://www.crockford.com/mckeeman.html
func JSON() *JSONReader {
	g := &JSONReader{}
	SetRuleNames(g)
	g.WS.Sub = Zom(WS())
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
	g.Characters.Sub = Zom(&g.Character)
	g.String.Sub = Seq(Rune('"'), &g.Characters, Rune('"'))
	g.Element.Sub = Seq(&g.WS, &g.Value, &g.WS)
	g.Elements.Sub = Seq(&g.Element, Zom(Seq(Rune(','), &g.Element)))
	g.Array.Sub = Seq(Rune('['), Any(&g.Elements, &g.WS), Rune(']'))
	g.Member.Sub = Seq(&g.WS, &g.String, &g.WS, Rune(':'), &g.Element)
	g.Members.Sub = Seq(&g.Member, Zom(Seq(Rune(','), &g.Member)))
	g.Object.Sub = Seq(Rune('{'), Any(&g.Members, &g.WS), Rune('}'))
	g.Value.Sub = Any(&g.Object, &g.Array, &g.String, &g.Number, Lit("true"), Lit("false"), Lit("null"))
	return g
}

func (r *JSONReader) Read(s *Scanner) error {
	err := r.Element.Read(s)
	if err != nil {
		l, c := s.LineCol(1)
		return fmt.Errorf("json parse error at %d:%d: %v", l, c, err)
	}
	return nil
}

func (r *JSONReader) What() string {
	return "json"
}

func (r *JSONReader) Grammar() Rules {
	return CollectRules(r)
}
