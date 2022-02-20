package grammar

import (
	"fmt"

	. "github.com/aiq/tok"
)

type LuaReader struct {
	Name             RuleReader `name:"Name"`
	Numeral          RuleReader `name:"Numeral"`
	LiteralString    RuleReader `name:"LiteralString"`
	UnOp             RuleReader `name:"unop"`
	BinOp            RuleReader `name:"binop"`
	FieldSep         RuleReader `name:"fieldsep"`
	Field            RuleReader `name:"field"`
	FieldList        RuleReader `name:"fieldlist"`
	TableConstructor RuleReader `name:"tableconstructor"`
	ParList          RuleReader `name:"parlist"`
	FuncBody         RuleReader `name:"funcbody"`
	FunctionDef      RuleReader `name:"functiondef"`
	Args             RuleReader `name:"args"`
	FunctionCall     RuleReader `name:"functioncall"`
	PrefixExp        RuleReader `name:"prefixexp"`
	Exp              RuleReader `name:"exp"`
	ExpList          RuleReader `name:"explist"`
	NameList         RuleReader `name:"namelist"`
	Var              RuleReader `name:"var"`
	VarList          RuleReader `name:"varlist"`
	FuncName         RuleReader `name:"funcname"`
	Label            RuleReader `name:"label"`
	RetStat          RuleReader `name:"retstat"`
	Attrib           RuleReader `name:"attrib"`
	AttNameList      RuleReader `name:"attnamelist"`
	Stat             RuleReader `name:"stat"`
	Block            RuleReader `name:"block"`
	Chunk            RuleReader `name:"chunk"`
}

// Lua creates a Grammar to Read a Lua file.
// The implementation is based on https://www.lua.org/manual/5.4/manual.html#9
func Lua() *LuaReader {
	r := Rune
	g := &LuaReader{}
	SetRuleNames(g)
	g.Name.Reader = Lit("a")
	g.Numeral.Reader = Lit("f")
	g.LiteralString.Reader = Lit("f")
	g.UnOp.Reader = Any("-", "not", "#", "~")
	g.BinOp.Reader = Any(
		"+", "-", "*", "/", "//", "^", "%", "&", "~", "|", ">>", "<<",
		"..", "<", "<=", ">", ">=", "==", "~=", "and", "or",
	)
	g.FieldSep.Reader = AnyRune(",;")
	g.Field.Reader = Any(
		Seq(r('['), &g.Exp, r(']'), r('='), &g.Exp),
		Seq(&g.Name, r('='), &g.Exp),
		&g.Exp,
	)
	g.FieldList.Reader = Seq(&g.Field, Zom(Seq(&g.FieldSep, &g.Field)), Opt(&g.FieldSep))
	g.TableConstructor.Reader = Seq(r('{'), Opt(&g.FieldList), r('}'))
	g.ParList.Reader = Any(Seq(&g.NameList, Opt(Seq(r(','), Lit("...")))), Lit("..."))
	g.FuncBody.Reader = Seq(r('('), &g.ParList, r(')'), &g.Block, Lit("end"))
	g.FunctionDef.Reader = Seq(Lit("function"), &g.FuncBody)
	g.Args.Reader = Any(
		Seq(r('('), &g.ExpList, r(')')),
		&g.TableConstructor,
		&g.LiteralString,
	)
	g.FunctionCall.Reader = Any(
		Seq(&g.PrefixExp, &g.Args),
		Seq(&g.PrefixExp, r(':'), &g.Name, &g.Args),
	)
	g.PrefixExp.Reader = Any(&g.Var, &g.FunctionCall, Seq(r('('), &g.Exp, r(')')))
	g.Exp.Reader = Any(
		Lit("nil"), Lit("false"), Lit("true"), &g.Numeral, &g.LiteralString, Lit("..."),
		&g.FunctionDef, &g.PrefixExp, &g.TableConstructor,
		Seq(&g.Exp, &g.BinOp, &g.Exp), Seq(&g.UnOp, &g.Exp),
	)
	g.ExpList.Reader = Seq(&g.Exp, Zom(Seq(r(','), &g.Exp)))
	g.NameList.Reader = Seq(&g.Name, Zom(Seq(r(','), &g.Name)))
	g.Var.Reader = Any(
		&g.Name,
		Seq(&g.PrefixExp, r('['), &g.Exp, r(']')),
		Seq(&g.PrefixExp, r('.'), &g.Name),
	)
	g.VarList.Reader = Seq(&g.Var, Zom(Seq(r(','), &g.Var)))
	g.FuncName.Reader = Seq(&g.Name, Zom(Seq(r('.'), &g.Name)), Opt(Seq(r(':'), &g.Name)))
	g.Label.Reader = Seq(Lit("::"), &g.Name, Lit("::"))
	g.RetStat.Reader = Seq(Lit("return"), Opt(&g.ExpList), Opt(r(';')))
	g.Attrib.Reader = Opt(Seq(r('<'), &g.Name, r('>')))
	g.AttNameList.Reader = Seq(&g.Name, &g.Attrib, Zom(Seq(r(','), &g.Name, &g.Attrib)))
	g.Stat.Reader = Any(
		r(';'),
		Seq(&g.VarList, r('='), &g.ExpList),
		&g.FunctionCall,
		&g.Label,
		Lit("break"),
		Seq(Lit("goto"), &g.Name),
		Seq(Lit("do"), &g.Block, Lit("end")),
		Seq(Lit("while"), &g.Exp, Lit("do"), &g.Block, Lit("end")),
		Seq(Lit("repeat"), &g.Block, Lit("until"), &g.Exp),
		Seq(Lit("if"), &g.Exp, Lit("then"), &g.Block, Zom(Seq(Lit("elseif"), &g.Exp, Lit("then"), &g.Block)), Opt(Seq(Lit("else"), &g.Block)), Lit("end")),
		Seq(Lit("for"), &g.Name, r('='), &g.Exp, r(','), &g.Exp, Opt(Seq(r(','), &g.Exp)), Lit("do"), &g.Block, Lit("end")),
		Seq(Lit("for"), &g.NameList, Lit("in"), &g.ExpList, Lit("do"), &g.Block, Lit("end")),
		Seq(Lit("function"), &g.FuncName, &g.FuncBody),
		Seq(Lit("local"), Lit("function"), &g.Name, &g.FuncBody),
		Seq(Lit("local"), &g.AttNameList, Opt(Seq(r('='), &g.ExpList))),
	)
	g.Block.Reader = Seq(Zom(&g.Stat), Opt(&g.RetStat))
	g.Chunk.Reader = &g.Block

	return g
}

func (r *LuaReader) Read(s *Scanner) error {
	var err error
	if err != nil {
		return fmt.Errorf("mxt parse error: %v", err)
	}
	return nil
}

func (r *LuaReader) What() string {
	return "mxt"
}

func (r *LuaReader) Grammar() Rules {
	return CollectRules(r)
}
