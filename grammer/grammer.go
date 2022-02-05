package grammer

import (
	"fmt"

	"github.com/aiq/tok"
)

type Rule struct {
	Name string
	What string
}

func (r Rule) String() string {
	return fmt.Sprintf("%s: %s", r.Name, r.What)
}

type Rules []Rule

func (rs Rules) Lines() []string {
	res := []string{}
	for _, r := range rs {
		res = append(res, r.String())
	}
	return res
}

type Grammar interface {
	Grammar() Rules
}

type RefReader struct {
	Name string
	Sub  tok.Reader
}

func (r *RefReader) Read(s *tok.Scanner) error {
	return r.Sub.Read(s)
}

func (r *RefReader) What() string {
	return r.Name
}

func (r RefReader) Rule() Rule {
	return Rule{r.Name, r.Sub.What()}
}

func Ref(name string) RefReader {
	return RefReader{name, nil}
}
