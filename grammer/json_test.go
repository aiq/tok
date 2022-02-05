package grammer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aiq/tok"
)

func TestJSON(t *testing.T) {
	fmt.Println(strings.Join(JSON().Grammer().Lines(), "\n"))

	hitCases := []string{
		`{}`,
		`{"key":"value"}`,
		`{
			"glossary": {
				"title": "example glossary",
				"GlossDiv": {
					"title": "S",
					"GlossList": {
						"GlossEntry": {
							"ID": "SGML",
							"SortAs": "SGML",
							"GlossTerm": "Standard Generalized Markup Language",
							"Acronym": "SGML",
							"Abbrev": "ISO 8879:1986",
							"GlossDef": {
								"para": "A meta-markup language, used to create markup languages such as DocBook.",
								"GlossSeeAlso": ["GML", "XML"]
							},
							"GlossSee": "markup"
						}
					}
				}
			}
		}`,
	}
	for i, c := range hitCases {
		sca := tok.NewScanner(c)
		err := sca.Use(JSON())
		if err != nil {
			t.Errorf("%d unexpected error: %v", i, err)
		}
	}
}
