package tok

import (
	"fmt"
	"reflect"
)

//------------------------------------------------------------------------------

// GrammarRules calls the Rule function on all rules and returns the result.
func GrammarLines(g []*Rule) []string {
	res := []string{}
	for _, r := range g {
		res = append(res, r.Rule())
	}
	return res
}

//------------------------------------------------------------------------------

// Grammer is a Reader that has a Grammar.
type Grammar interface {
	Reader
	Grammar() []*Rule
}

// SetRuleNames sets the rule names via the associated Field-Tag.
func SetRuleNames(g interface{}) error {
	t := reflect.TypeOf(g).Elem()
	v := reflect.ValueOf(g).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type != reflect.TypeOf(Rule{}) {
			continue
		}
		name := field.Tag.Get("name")
		if name == "" {
			continue
		}
		if e := CheckRuleName(name); e != nil {
			return fmt.Errorf("invalid name for %s: %v", field.Name, e)
		}
		rule := v.Field(i)
		ruleName := rule.FieldByName("Name")
		ruleName.Set(reflect.ValueOf(name))
	}
	return nil
}

func MustSetRuleNames(g interface{}) {
	err := SetRuleNames(g)
	if err != nil {
		panic(err)
	}
}

func CollectRules(g interface{}) []*Rule {
	rules := []*Rule{}
	v := reflect.ValueOf(g).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		i := field.Addr().Interface()
		ptr, ok := i.(*Rule)
		if !ok {
			continue
		}
		rules = append(rules, ptr)
	}
	return rules
}

// CheckRules checks if the RuleReaders have a Name and a Reader set.
func CheckRules(g interface{}) error {
	t := reflect.TypeOf(g).Elem()
	v := reflect.ValueOf(g).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type != reflect.TypeOf(Rule{}) {
			continue
		}
		rule := v.Field(i)
		name, ok := rule.FieldByName("Name").Interface().(string)
		if !ok {
			return fmt.Errorf("the Name of %s is not a string", field.Name)
		}
		if e := CheckRuleName(name); e != nil {
			return e
		}

		reader := rule.FieldByName("Reader").Interface()
		if reader == nil {
			return fmt.Errorf("the Reader of %s is a nil value", field.Name)
		}
		if _, ok := reader.(Reader); !ok {
			return fmt.Errorf("the Reader of %s is not a Reader", field.Name)
		}
	}
	return nil
}

// MustCheckRules panics if an error occurs during CheckRules.
func MustCheckRules(g interface{}) {
	err := CheckRules(g)
	if err != nil {
		panic(err)
	}
}

//------------------------------------------------------------------------------
// RuleReader can be used to set the rules of a grammar.
type Rule struct {
	Name   string
	Reader Reader
}

// Map connects the move of a Reader to the function f.
func (r *Rule) Map(f MapFunc) {
	r.Reader = Map(r.Reader, f)
}

func (r *Rule) Monitor(l *Log) {
	r.Reader = Monitor(r.Reader, l, r.Name)
}

// Pick collects the Segments if a Reader was moven and sets the Info field with the Reader Name.
func (r *Rule) Pick(basket *Basket) {
	r.Reader = Pick(r.Reader, basket, r.Name)
}

func (r *Rule) Read(s *Scanner) error {
	return r.Reader.Read(s)
}

func (r *Rule) What() string {
	return r.Name
}

func (r *Rule) Rule() string {
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

// CheckRuleName validates the Name field of a RuleReader.
func CheckRuleName(name string) error {
	return NewScanner(name).Use(RuleName())
}

// RuleName returns a Reader to validate a Name field of a RuleReader.
func RuleName() Reader {
	return &ruleNameReader{
		sub: Seq(BetweenAny("a-zA-Z"), Zom(Any(BetweenAny("a-zA-Z0-9"), AnyRune("+-._")))),
	}
}
