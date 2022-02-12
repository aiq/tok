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

	values := []tok.Value{}
	sca := tok.NewScanner(string(inp))
	reader := grammar.JSON()
	reader.Key.Pick("key", &values)
	reader.Object.Pick("object", &values)
	reader.Array.Pick("array", &values)
	reader.String.Pick("string", &values)
	reader.Number.Pick("number", &values)
	reader.Bool.Pick("bool", &values)
	reader.Null.Pick("null", &values)
	err = sca.Use(reader)
	if err != nil {
		log.Fatalf("invalid log json file %q: %v", filename, err)
	}
	g := tok.BuildGraph(filename, values)
	fmt.Print(g.FlameStack())
}
