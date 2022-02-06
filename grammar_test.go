package tok

import (
	"fmt"
	"testing"
)

func unexpError(i int, e error) string {
	return fmt.Sprintf("%d unexpected error: %v", i, e)
}

func unexpRuleName(name string, exp string) string {
	return fmt.Sprintf("unexpected rule name: %q != %q", name, exp)
}

func TestSetRuleNames(t *testing.T) {
	{
		g := struct {
			Rule1  RuleReader `name:"first"`
			Ignore Reader
			Rule2  RuleReader `name:"num2"`
			Rule3  RuleReader
		}{}
		err := SetRuleNames(&g)
		if err != nil {
			t.Error(unexpError(0, err))
		}
		if g.Rule1.Name != "first" {
			t.Error(unexpRuleName(g.Rule1.Name, "first"))
		}
		if g.Rule2.Name != "num2" {
			t.Error(unexpRuleName(g.Rule2.Name, "num2"))
		}
		if g.Rule3.Name != "" {
			t.Error(unexpRuleName(g.Rule2.Name, ""))
		}
	}
}
