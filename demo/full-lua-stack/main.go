package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aiq/tok"
	"github.com/aiq/tok/grammar"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("demo requires the name of a json file as input")
	}
	filename := os.Args[1]

	inp, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("not able to read %q: %v", filename, err)
	}

	sca := tok.NewScanner(string(inp))
	lua := grammar.Lua()
	basket := sca.NewBasket()
	basket.PickWith(
		&lua.Name, &lua.Numeral, &lua.LiteralString, &lua.Comment,
		&lua.UnOp, &lua.BinOp, &lua.Field, &lua.FieldList, &lua.TableConstructor,
		&lua.FuncParams, &lua.FuncBody, &lua.FuncDef, &lua.FuncArgs, &lua.FuncCall,
		&lua.PrefixExp, &lua.Exp, &lua.ExpList, &lua.NameList, &lua.Var, &lua.VarList,
		&lua.FuncName, &lua.Label, &lua.RetStat, &lua.Attrib, &lua.AttNameList, &lua.Break,
		&lua.GoTo, &lua.Do, &lua.While, &lua.Repeat, &lua.IfElse, &lua.For, &lua.ForEach,
		&lua.Func, &lua.LocalFunc, &lua.LocalAtt, &lua.Stat, &lua.Block, &lua.Chunk,
	)
	err = sca.Use(lua)
	if err != nil {
		log.Fatalf("invalid lua file %q: %v", filename, err)
	}
	g := tok.BuildGraph(filename, basket.Picked())
	fmt.Print(g.FlameStack())
}
