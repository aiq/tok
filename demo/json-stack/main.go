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

	sca := tok.NewScanner(string(inp))
	reader := grammar.JSON()
	basket := sca.NewBasket()
	reader.Key.Pick(basket)
	reader.Object.Pick(basket)
	reader.Array.Pick(basket)
	reader.String.Pick(basket)
	reader.Number.Pick(basket)
	reader.Bool.Pick(basket)
	reader.Null.Pick(basket)
	err = sca.Use(reader)
	if err != nil {
		log.Fatalf("invalid log json file %q: %v", filename, err)
	}
	g := tok.BuildGraph(filename, basket.Picked())
	fmt.Print(g.FlameStack())
}
