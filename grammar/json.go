package grammar

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
	Key        RuleReader `name:"key"`
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
	Bool       RuleReader `name:"bool"`
	Null       RuleReader `name:"null"`
	WS         RuleReader `name:"ws"`
}

// JSON creates a Grammar to Read a JSON file.
// The implementation is based on https://www.crockford.com/mckeeman.html
func JSON() *JSONReader {
	g := &JSONReader{}
	SetRuleNames(g)
	g.WS.Reader = Zom(WS())
	g.Null.Reader = Lit("null")
	g.Bool.Reader = Any("true", "false")
	g.Sign.Reader = Opt(AnyRune("+-"))
	g.OneNine.Reader = Between('1', '9')
	g.Digit.Reader = Digit()
	g.Digits.Reader = Many(Digit())
	g.Exponent.Reader = Opt(Seq(AnyRune("eE"), &g.Sign, &g.Digits))
	g.Fraction.Reader = Opt(Seq('.', &g.Digits))
	g.Integer.Reader = Seq(Opt(Rune('-')), Any(Rune('0'), Seq(&g.OneNine, Opt(&g.Digits))))
	g.Number.Reader = Seq(&g.Integer, &g.Fraction, &g.Exponent)
	g.Hex.Reader = HexDigit()
	g.Escape.Reader = Any(AnyRune(`"\/bfnrt`), Seq('u', Times(4, &g.Hex)))
	g.Character.Reader = Any(Holey(' ', utf8.MaxRune, `"\`), Seq('\\', &g.Escape))
	g.Characters.Reader = Zom(&g.Character)
	g.String.Reader = Seq('"', &g.Characters, '"')
	g.Key.Reader = g.String.Reader
	g.Element.Reader = Seq(&g.WS, &g.Value, &g.WS)
	g.Elements.Reader = Seq(&g.Element, Zom(Seq(Rune(','), &g.Element)))
	g.Array.Reader = Seq('[', Any(&g.Elements, &g.WS), ']')
	g.Member.Reader = Seq(&g.WS, &g.Key, &g.WS, ':', &g.Element)
	g.Members.Reader = Seq(&g.Member, Zom(Seq(',', &g.Member)))
	g.Object.Reader = Seq('{', Any(&g.Members, &g.WS), '}')
	g.Value.Reader = Any(&g.Object, &g.Array, &g.String, &g.Number, &g.Bool, &g.Null)
	return g
}

func (r *JSONReader) Read(s *Scanner) error {
	err := r.Element.Read(s)
	if err != nil {
		return fmt.Errorf("json parse error: %v", err)
	}
	return nil
}

func (r *JSONReader) What() string {
	return "json"
}

func (r *JSONReader) Grammar() Rules {
	return CollectRules(r)
}
