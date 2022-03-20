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
	normalString := Seq('"', Zom(Any(Holey(' ', utf8.MaxRune, `"\`), escapeSequence)), '"')
	charString := Seq('\'', Zom(Any(Holey(' ', utf8.MaxRune, `'\`), escapeSequence)), '\'')
	strHead, strTail := Janus("", Zom("="))
	longBeg := Seq('[', strHead, '[')
	longEnd := Seq(']', strTail, ']')
	longString := Seq(longBeg, Past(longEnd))
	return Any(normalString, charString, longString)
}

func LuaComment() Reader {
	line := Seq("--", To('\n'))
	cmtHead, cmtTail := Janus("", Zom("="))
	longBeg := Seq("--[", cmtHead, '[')
	longEnd := Seq(']', cmtTail, ']')
	long := Seq(longBeg, Past(longEnd))
	return Any(long, line)
}

type LuaReader struct {
	SheBang       Rule `name:"SheBang"`
	Name          Rule `name:"Name"`
	Numeral       Rule `name:"Numeral"`
	LiteralString Rule `name:"LiteralString"`
	Comment       Rule `name:"Comment"`

	UnOp  Rule `name:"unop"`
	BinOp Rule `name:"binop"`

	FieldSep         Rule `name:"fieldsep"`
	Field            Rule `name:"field"`
	FieldList        Rule `name:"fieldlist"`
	TableConstructor Rule `name:"tableconstructor"`

	FuncParams Rule `name:"funcparams"`
	FuncBody   Rule `name:"funcbody"`
	FuncDef    Rule `name:"funcdef"`
	FuncArgs   Rule `name:"funcargs"`
	FuncCall   Rule `name:"funccall"`

	PrefixExp Rule `name:"prefixexp"`
	FinalExp  Rule `name:"finalexp"`
	Exp       Rule `name:"exp"`
	ExpList   Rule `name:"explist"`
	NameList  Rule `name:"namelist"`
	VarSuffix Rule `name:"varsuffix"`
	Var       Rule `name:"var"`
	VarList   Rule `name:"varlist"`

	FuncName    Rule `name:"funcname"`
	Label       Rule `name:"label"`
	RetStat     Rule `name:"retstat"`
	Attrib      Rule `name:"attrib"`
	AttNameList Rule `name:"attnamelist"`
	Break       Rule `name:"break"`
	GoTo        Rule `name:"goto"`
	Do          Rule `name:"do"`
	While       Rule `name:"while"`
	Repeat      Rule `name:"repeat"`
	IfElse      Rule `name:"ifelse"`
	For         Rule `name:"for"`
	ForEach     Rule `name:"foreach"`
	Func        Rule `name:"func"`
	LocalFunc   Rule `name:"localfunc"`
	LocalAtt    Rule `name:"localatt"`
	Stat        Rule `name:"stat"`

	Block  Rule `name:"block"`
	Chunk  Rule `name:"chunk"`
	Script Rule `name:"script"`
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
		"..", "<=", "<", ">=", ">", "==", "~=", "and", "or",
		"+", "-", "*", "/", "//", "^", "%", "&", "~", "|", ">>", "<<",
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

	nameAndArgs := skipSeq(Opt(Seq(':', &g.Name)), &g.FuncArgs)
	varOrExp := Any(&g.Var, skipSeq('(', &g.Exp, ')'))

	g.FuncParams.Reader = Any(skipSeq(&g.NameList, Opt(skipSeq(',', "..."))), "...")
	g.FuncBody.Reader = skipSeq('(', Opt(&g.FuncParams), ')', &g.Block, "end")
	g.FuncDef.Reader = SkipWSSeq(Lit("function"), &g.FuncBody)
	g.FuncArgs.Reader = Any(
		skipSeq('(', Opt(&g.ExpList), ')'),
		&g.TableConstructor,
		&g.LiteralString,
	)
	g.FuncCall.Reader = Seq(varOrExp, Many(nameAndArgs))

	g.PrefixExp.Reader = Seq(varOrExp, Zom(nameAndArgs))
	g.FinalExp.Reader = Any(
		"nil", "false", "true",
		LuaNumeral(), &g.LiteralString, "...",
		&g.FuncDef,
		skipSeq(&g.UnOp, &g.Exp),
		&g.TableConstructor,
		&g.PrefixExp,
	)
	g.Exp.Reader = Any(
		SkipWSSeq(&g.FinalExp, &g.BinOp, &g.Exp),
		&g.FinalExp,
	)
	g.ExpList.Reader = skipSeq(&g.Exp, Zom(skipSeq(',', &g.Exp)))
	g.VarSuffix.Reader = skipSeq(Zom(nameAndArgs), Any(
		skipSeq('[', &g.Exp, ']'),
		Seq('.', &g.Name),
	))
	g.Var.Reader = Seq(Any(
		&g.Name,
		skipSeq('(', &g.Exp, ')', &g.VarSuffix),
	), Zom(&g.VarSuffix))
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
		&g.FuncCall,
	)
	g.Block.Reader = skipSeq(Zom(&g.Stat), Opt(&g.RetStat))
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

func (r *LuaReader) Grammar() []*Rule {
	return CollectRules(r)
}
