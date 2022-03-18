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

// CheckRules checks if the RuleReaders have a Name and a Reader set.
func CheckRules(g interface{}) error {
	t := reflect.TypeOf(g).Elem()
	v := reflect.ValueOf(g).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type != reflect.TypeOf(RuleReader{}) {
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
type RuleReader struct {
	Name   string
	Reader Reader
}

// Map connects the move of a Reader to the function f.
func (r *RuleReader) Map(f MapFunc) {
	r.Reader = Map(r.Reader, f)
}

// PickAs collects the Segments if a Reader was moven and sets the Info field.
func (r *RuleReader) PickAs(basket *Basket, info string) {
	r.Reader = Pick(r.Reader, basket, info)
}

// Pick collects the Segments if a Reader was moven and sets the Info field with the Reader Name.
func (r *RuleReader) Pick(basket *Basket) {
	r.PickAs(basket, r.Name)
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
