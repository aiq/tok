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

type setRuleNamesGrammar struct {
	Rule1  Rule `name:"first"`
	Ignore Reader
	Rule2  Rule `name:"num2"`
	Rule3  Rule
}

func (r *setRuleNamesGrammar) Read(s *Scanner) error {
	return nil
}

func (r *setRuleNamesGrammar) What() string {
	return ""
}

func (g *setRuleNamesGrammar) Grammar() []*Rule {
	return []*Rule{&g.Rule1, &g.Rule2, &g.Rule3}
}

func TestSetRuleNames(t *testing.T) {
	{
		g := setRuleNamesGrammar{}
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
