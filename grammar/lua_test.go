package grammar

import (
	"testing"

	"github.com/aiq/tok"
)

func TestLuaParts(t *testing.T) {
	lua := Lua()
	l := tok.MonitorGrammar(lua)
	cases := []struct {
		inp string
		r   tok.Reader
	}{
		{`""`, &lua.LiteralString},
		{`"a\tbc\n"`, &lua.LiteralString},
		{`'a\tbc\n'`, &lua.LiteralString},
		{`[[abc]]`, &lua.LiteralString},
		{`assert( dir and dir ~= "", "directory parameter is missing or empty" )`, &lua.FuncCall},
		{`if not isdodd( base ) then base = doSomething( base ) end`, &lua.IfElse},
		{`return coroutine.wrap( function() yieldtree( dir ) end )`, &lua.RetStat},
	}
	for i, c := range cases {
		l.Reset()
		sca := tok.NewScanner(c.inp)
		err := sca.Use(c.r)
		l.PrintWithPreview(c.inp, 10)
		if err != nil {
			t.Errorf("%d unexpected error: %v", i, err)
			l.PrintWithPreview(c.inp, 10)
		} else {
			tok.BuildGraph("lua", sca.NewBasket().Picked())
		}
	}
}

func TestLua(t *testing.T) {
	posCases := []struct {
		lua       string
		funcnames []string
	}{
		{``, []string{}},
	}
	for i, c := range posCases {
		sca := tok.NewScanner(c.lua)
		r := Lua()
		err := sca.Use(r)
		if err != nil {
			t.Errorf("%d unexpected error: %v", i, err)
		}
		if !sca.AtEnd() {
			t.Errorf("did not read the whole lua")
		}
	}
}
