package grammar

import (
	"testing"

	"github.com/aiq/tok"
)

func TestJSON(t *testing.T) {

	posCases := []string{
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
		`{"menu": {
			"id": "file",
			"value": "File",
			"popup": {
			  "menuitem": [
				{"value": "New", "onclick": "CreateNewDoc()"},
				{"value": "Open", "onclick": "OpenDoc()"},
				{"value": "Close", "onclick": "CloseDoc()"}
			  ]
			}
		  }}`,
	}
	for i, c := range posCases {
		sca := tok.NewScanner(c)
		err := sca.Use(JSON())
		if err != nil {
			t.Errorf("%d unexpected error: %v", i, err)
		}
		if !sca.AtEnd() {
			t.Errorf("did not read the whole json")
		}
	}
}
