package tok

import (
	"fmt"
	"reflect"
)

//------------------------------------------------------------------------------

// Rule
type Rule interface {
	// Returns a string with the following syntax: <name>: <what>.
	Rule() string
}

// Rules
type Rules []Rule

// Lines calls the Rule function on all rules and returns the result.
func (rs Rules) Lines() []string {
	res := []string{}
	for _, r := range rs {
		res = append(res, r.Rule())
	}
	return res
}

//------------------------------------------------------------------------------

// Grammer is a Reader that has a Grammar.
type Grammar interface {
	Reader
	Grammar() Rules
}

// SetRuleNames sets the rule names via the associated Field-Tag.
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

// CollectRules collects the RuleReaders that the grammar g has.
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
// RuleReader can be used to set the rules of a grammar.
type RuleReader struct {
	Name   string
	Reader Reader
}

// Map connects the move of a Reader to the function f.
func (r *RuleReader) Map(f MapFunc) {
	r.Reader = Map(r.Reader, f)
}

// Pick collects the Tokens if a Reader was moven and sets the Info field.
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
