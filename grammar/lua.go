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
		Seq('[', &g.Exp, ']', '=', &g.Exp),
		Seq(&g.Name, '=', &g.Exp),
		&g.Exp,
	)
	g.FieldList.Reader = Seq(&g.Field, Zom(Seq(&g.FieldSep, &g.Field)), Opt(&g.FieldSep))
	g.TableConstructor.Reader = Seq('{', Opt(&g.FieldList), '}')
	g.ParList.Reader = Any(Seq(&g.NameList, Opt(Seq(',', "..."))), "...")
	g.FuncBody.Reader = Seq('(', &g.ParList, ')', &g.Block, "end")
	g.FunctionDef.Reader = Seq(Lit("function"), &g.FuncBody)
	g.Args.Reader = Any(
		Seq('(', &g.ExpList, ')'),
		&g.TableConstructor,
		&g.LiteralString,
	)
	g.FunctionCall.Reader = Any(
		Seq(&g.PrefixExp, &g.Args),
		Seq(&g.PrefixExp, ':', &g.Name, &g.Args),
	)
	g.PrefixExp.Reader = Any(&g.Var, &g.FunctionCall, Seq('(', &g.Exp, ')'))
	g.Exp.Reader = Any(
		"nil", "false", "true", &g.Numeral, &g.LiteralString, "...",
		&g.FunctionDef, &g.PrefixExp, &g.TableConstructor,
		Seq(&g.Exp, &g.BinOp, &g.Exp), Seq(&g.UnOp, &g.Exp),
	)
	g.ExpList.Reader = Seq(&g.Exp, Zom(Seq(',', &g.Exp)))
	g.NameList.Reader = Seq(&g.Name, Zom(Seq(',', &g.Name)))
	g.Var.Reader = Any(
		&g.Name,
		Seq(&g.PrefixExp, '[', &g.Exp, ']'),
		Seq(&g.PrefixExp, '.', &g.Name),
	)
	g.VarList.Reader = Seq(&g.Var, Zom(Seq(',', &g.Var)))
	g.FuncName.Reader = Seq(&g.Name, Zom(Seq('.', &g.Name)), Opt(Seq(':', &g.Name)))
	g.Label.Reader = Seq(Lit("::"), &g.Name, Lit("::"))
	g.RetStat.Reader = Seq(Lit("return"), Opt(&g.ExpList), Opt(';'))
	g.Attrib.Reader = Opt(Seq('<', &g.Name, '>'))
	g.AttNameList.Reader = Seq(&g.Name, &g.Attrib, Zom(Seq(',', &g.Name, &g.Attrib)))
	g.Stat.Reader = Any(
		';',
		Seq(&g.VarList, '=', &g.ExpList),
		&g.FunctionCall,
		&g.Label,
		"break",
		Seq("goto", &g.Name),
		Seq("do", &g.Block, "end"),
		Seq("while", &g.Exp, "do", &g.Block, "end"),
		Seq("repeat", &g.Block, "until", &g.Exp),
		Seq("if", &g.Exp, "then", &g.Block,
			Zom(Seq("elseif", &g.Exp, "then", &g.Block)),
			Opt(Seq("else", &g.Block)),
			"end",
		),
		Seq("for", &g.Name, '=', &g.Exp, ',', &g.Exp, Opt(Seq(',', &g.Exp)), "do", &g.Block, "end"),
		Seq("for", &g.NameList, "in", &g.ExpList, "do", &g.Block, "end"),
		Seq("function", &g.FuncName, &g.FuncBody),
		Seq("local", "function", &g.Name, &g.FuncBody),
		Seq("local", &g.AttNameList, Opt(Seq('=', &g.ExpList))),
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
