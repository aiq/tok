package grammar

import (
	"fmt"
	"unicode/utf8"

	. "github.com/aiq/tok"
)

func LuaNumeral() Reader {
	hexExponent := Seq(AnyRune("pP"), Opt(AnyRune("+-")), Many(Digit()))
	hexFloat := Seq('0', AnyRune("xX"), Many(HexDigit()), Opt(Seq('.', Many(HexDigit()))), Opt(hexExponent))
	exponent := Seq(AnyRune("eE"), Opt(AnyRune("+-")), Many(Digit()))
	floatValue := Seq(Many(Digit()), Opt(Seq('.', Many(Digit()))), Opt(exponent))
	intValue := Many(Digit())
	hexValue := Seq('0', AnyRune("xX"), Many(HexDigit()))
	return Any(hexFloat, hexValue, floatValue, intValue)
}

func LuaString() Reader {
	utfEscape := Seq("\\u{", Many(HexDigit()), "}")
	hexEscape := Seq("\\x", Times(2, HexDigit()))
	decimalEscape := Any(
		Seq('\\', Digit()),
		Seq('\\', Times(2, Digit())),
		Seq('\\', Between('0', '2'), Times(2, Digit())),
	)
	escapeSequence := Any(
		Seq('\\', AnyRune("abfnrtvz\"'\\")),
		decimalEscape,
		hexEscape,
		utfEscape,
	)
	normalString := Seq('"', Any(Holey(' ', utf8.MaxRune, `"\`), escapeSequence), '"')
	charString := Seq('\'', Any(Holey(' ', utf8.MaxRune, `'\`), escapeSequence), '\'')
	strHead, strTail := Janus("", Zom("="))
	longString := Seq('[', strHead, '[', Between(' ', utf8.MaxRune), ']', strTail, ']')
	return Any(normalString, charString, longString)
}

func LuaComment() Reader {
	line := Seq("--", To('\n'))
	cmtHead, cmtTail := Janus("", Zom("="))
	long := Seq("--[", cmtHead, '[', Between(' ', utf8.MaxRune), ']', cmtTail, ']')
	return Any(long, line)
}

type LuaReader struct {
	SheBang       RuleReader `name:"SheBang"`
	Name          RuleReader `name:"Name"`
	Numeral       RuleReader `name:"Numeral"`
	LiteralString RuleReader `name:"LiteralString"`
	Comment       RuleReader `name:"Comment"`

	UnOp  RuleReader `name:"unop"`
	BinOp RuleReader `name:"binop"`

	FieldSep         RuleReader `name:"fieldsep"`
	Field            RuleReader `name:"field"`
	FieldList        RuleReader `name:"fieldlist"`
	TableConstructor RuleReader `name:"tableconstructor"`

	FuncParams RuleReader `name:"funcparams"`
	FuncBody   RuleReader `name:"funcbody"`
	FuncDef    RuleReader `name:"funcdef"`
	FuncArgs   RuleReader `name:"funcargs"`
	FuncCall   RuleReader `name:"funccall"`

	PrefixExp RuleReader `name:"prefixexp"`
	FinalExp  RuleReader `name:"finalexp"`
	Exp       RuleReader `name:"exp"`
	ExpList   RuleReader `name:"explist"`
	NameList  RuleReader `name:"namelist"`
	Var       RuleReader `name:"var"`
	VarList   RuleReader `name:"varlist"`

	FuncName    RuleReader `name:"funcname"`
	Label       RuleReader `name:"label"`
	RetStat     RuleReader `name:"retstat"`
	Attrib      RuleReader `name:"attrib"`
	AttNameList RuleReader `name:"attnamelist"`
	Break       RuleReader `name:"break"`
	GoTo        RuleReader `name:"goto"`
	Do          RuleReader `name:"do"`
	While       RuleReader `name:"while"`
	Repeat      RuleReader `name:"repeat"`
	IfElse      RuleReader `name:"ifelse"`
	For         RuleReader `name:"for"`
	ForEach     RuleReader `name:"foreach"`
	Func        RuleReader `name:"func"`
	LocalFunc   RuleReader `name:"localfunc"`
	LocalAtt    RuleReader `name:"localatt"`
	Stat        RuleReader `name:"stat"`

	Block  RuleReader `name:"block"`
	Chunk  RuleReader `name:"chunk"`
	Script RuleReader `name:"script"`
}

// Lua creates a Grammar to Read a Lua file.
// The implementation is based on https://www.lua.org/manual/5.4/manual.html#9
func Lua() *LuaReader {
	g := &LuaReader{}
	MustSetRuleNames(g)
	g.SheBang.Reader = Seq("#!", To(Any('\r', '\n')))
	g.Name.Reader = Seq(Set("a-zA-Z", "_"), Zom(Set("a-zA-Z0-9", "_")))
	g.NameList.Reader = SkipWSSeq(&g.Name, Zom(Seq(',', &g.Name)))
	g.Numeral.Reader = LuaNumeral()
	g.LiteralString.Reader = LuaString()
	g.Comment.Reader = LuaComment()

	g.UnOp.Reader = Any("-", "not", "#", "~")
	g.BinOp.Reader = Any(
		"+", "-", "*", "/", "//", "^", "%", "&", "~", "|", ">>", "<<",
		"..", "<", "<=", ">", ">=", "==", "~=", "and", "or",
	)

	skiper := Zom(Any(WS(), &g.Comment))
	skipSeq := func(list ...interface{}) Reader {
		return SkipSeq(skiper, list...)
	}

	g.FieldSep.Reader = AnyRune(",;")
	g.Field.Reader = Any(
		skipSeq('[', &g.Exp, ']', '=', &g.Exp),
		skipSeq(&g.Name, '=', &g.Exp),
		&g.Exp,
	)
	g.FieldList.Reader = skipSeq(&g.Field, Zom(Seq(&g.FieldSep, &g.Field)), Opt(&g.FieldSep))
	g.TableConstructor.Reader = skipSeq('{', Opt(&g.FieldList), '}')

	g.FuncParams.Reader = Any(skipSeq(&g.NameList, Opt(skipSeq(',', "..."))), "...")
	g.FuncBody.Reader = SkipWSSeq('(', &g.FuncParams, ')', &g.Block, "end")
	g.FuncDef.Reader = SkipWSSeq(Lit("function"), &g.FuncBody)
	g.FuncArgs.Reader = Any(
		skipSeq('(', &g.ExpList, ')'),
		&g.TableConstructor,
		&g.LiteralString,
	)
	g.FuncCall.Reader = Any(
		skipSeq(&g.PrefixExp, &g.FuncArgs),
		skipSeq(&g.PrefixExp, ':', &g.Name, &g.FuncArgs),
	)

	g.PrefixExp.Reader = Any(
		&g.FuncCall,
		&g.Var,
		skipSeq('(', &g.Exp, ')'),
	)
	g.FinalExp.Reader = Any(
		"nil", "false", "true",
		LuaNumeral(), &g.LiteralString, "...",
		&g.FuncDef,
		&g.PrefixExp,
		&g.TableConstructor,
		Seq(&g.UnOp, &g.Exp),
	)
	g.Exp.Reader = Any(
		SkipWSSeq(&g.FinalExp, &g.BinOp, &g.Exp),
		&g.FinalExp,
	)
	g.ExpList.Reader = SkipWSSeq(&g.Exp, Zom(Seq(',', &g.Exp)))
	g.Var.Reader = Any(
		&g.Name,
		skipSeq(&g.PrefixExp, '[', &g.Exp, ']'),
		skipSeq(&g.PrefixExp, '.', &g.Name),
	)
	g.VarList.Reader = SkipWSSeq(&g.Var, Zom(SkipWSSeq(',', &g.Var)))

	g.FuncName.Reader = SkipWSSeq(&g.Name, Zom(SkipWSSeq('.', &g.Name)), Opt(SkipWSSeq(':', &g.Name)))
	g.Label.Reader = Seq(Lit("::"), &g.Name, Lit("::"))
	g.RetStat.Reader = SkipWSSeq(Lit("return"), Opt(&g.ExpList), Opt(';'))
	g.Attrib.Reader = Opt(Seq('<', &g.Name, '>'))
	g.AttNameList.Reader = SkipWSSeq(&g.Name, &g.Attrib, Zom(SkipWSSeq(',', &g.Name, &g.Attrib)))
	g.Break.Reader = Lit("break")
	g.GoTo.Reader = SkipWSSeq("goto", &g.Name)
	g.Do.Reader = SkipWSSeq("do", &g.Block, "end")
	g.While.Reader = SkipWSSeq("while", &g.Exp, "do", &g.Block, "end")
	g.Repeat.Reader = SkipWSSeq("repeat", &g.Block, "until", &g.Exp)
	g.IfElse.Reader = SkipWSSeq(
		"if", &g.Exp, "then", &g.Block,
		Zom(SkipWSSeq("elseif", &g.Exp, "then", &g.Block)),
		Opt(SkipWSSeq("else", &g.Block)),
		"end",
	)
	g.For.Reader = SkipWSSeq("for", &g.Name, '=', &g.Exp, ',', &g.Exp, Opt(SkipWSSeq(',', &g.Exp)), "do", &g.Block, "end")
	g.ForEach.Reader = SkipWSSeq("for", &g.NameList, "in", &g.ExpList, "do", &g.Block, "end")
	g.Func.Reader = SkipWSSeq("function", &g.FuncName, &g.FuncBody)
	g.LocalFunc.Reader = SkipWSSeq("local", "function", &g.Name, &g.FuncBody)
	g.LocalAtt.Reader = SkipWSSeq("local", &g.AttNameList, Opt(SkipWSSeq('=', &g.ExpList)))
	g.Stat.Reader = Any(
		';',
		Seq(&g.VarList, '=', &g.ExpList),
		&g.FuncCall,
		&g.Label,
		&g.Break,
		&g.GoTo,
		&g.Do,
		&g.While,
		&g.Repeat,
		&g.IfElse,
		&g.For,
		&g.ForEach,
		&g.Func,
		&g.LocalFunc,
		&g.LocalAtt,
	)
	g.Block.Reader = SkipWSSeq(Zom(&g.Stat), Opt(&g.RetStat))
	g.Chunk.Reader = &g.Block
	g.Script.Reader = Seq(Opt(&g.SheBang), &g.Block)

	MustCheckRules(g)
	return g
}

func (r *LuaReader) Read(s *Scanner) error {
	err := r.Chunk.Read(s)
	if err != nil {
		return fmt.Errorf("lua parse error: %v", err)
	}
	return nil
}

func (r *LuaReader) What() string {
	return "lua"
}

func (r *LuaReader) Grammar() Rules {
	return CollectRules(r)
}
