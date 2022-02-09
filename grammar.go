package tok

import (
	"fmt"
	"reflect"
)

//------------------------------------------------------------------------------
type Rule interface {
	Rule() string
}

type Rules []Rule

func (rs Rules) Lines() []string {
	res := []string{}
	for _, r := range rs {
		res = append(res, r.Rule())
	}
	return res
}

//------------------------------------------------------------------------------
type Grammar interface {
	Grammar() Rules
}

func SetRuleNames(g interface{}) error {
	t := reflect.TypeOf(g).Elem()
	v := reflect.ValueOf(g).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type != reflect.TypeOf(RuleReader{}) {
			continue
		}
		name := field.Tag.Get("name")
		if name == "" {
			continue
		}
		if e := CheckRuleName(name); e != nil {
			return e
		}
		rule := v.Field(i)
		ruleName := rule.FieldByName("Name")
		ruleName.Set(reflect.ValueOf(name))
	}
	return nil
}

func CollectRules(g interface{}) Rules {
	rules := Rules{}
	v := reflect.ValueOf(g).Elem()
	for i := 0; i < v.NumField(); i++ {
		rule, ok := v.Field(i).Interface().(Rule)
		if ok {
			rules = append(rules, rule)
		}
	}
	return rules
}

//------------------------------------------------------------------------------
type RuleReader struct {
	Name   string
	Reader Reader
}

func (r *RuleReader) Map(f MapFunc) {
	r.Reader = Map(r.Reader, f)
}

func (r *RuleReader) Pick(info string, values *[]Value) {
	r.Map(func(t Token) {
		*values = append(*values, Value{
			Info:  info,
			Token: t,
		})
	})
}

func (r *RuleReader) Read(s *Scanner) error {
	return r.Reader.Read(s)
}

func (r *RuleReader) What() string {
	return r.Name
}

func (r RuleReader) Rule() string {
	return fmt.Sprintf("%s: %s", r.Name, r.Reader.What())
}

//------------------------------------------------------------------------------
type ruleNameReader struct {
	sub Reader
}

func (r ruleNameReader) Read(s *Scanner) error {
	err := r.sub.Read(s)
	return s.BoolErrorFor(err == nil, r.What())
}

func (r ruleNameReader) What() string {
	return "rulename"
}

func CheckRuleName(name string) error {
	return NewScanner(name).Use(RuleName())
}

func RuleName() Reader {
	return &ruleNameReader{
		sub: Seq(BetweenAny("a-zA-Z"), Zom(Any(BetweenAny("a-zA-Z0-9"), AnyRune("+-._")))),
	}
}
