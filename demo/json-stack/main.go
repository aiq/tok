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

	segs := []tok.Segment{}
	sca := tok.NewScanner(string(inp))
	reader := grammar.JSON()
	reader.Key.PickAs(&segs, "key")
	reader.Object.PickAs(&segs, "object")
	reader.Array.PickAs(&segs, "array")
	reader.String.PickAs(&segs, "string")
	reader.Number.PickAs(&segs, "number")
	reader.Bool.PickAs(&segs, "bool")
	reader.Null.PickAs(&segs, "null")
	err = sca.Use(reader)
	if err != nil {
		log.Fatalf("invalid log json file %q: %v", filename, err)
	}
	g := tok.BuildGraph(filename, segs)
	fmt.Print(g.FlameStack())
}
