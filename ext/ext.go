// custom types and helper functions
package ext

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"strings"
	"unicode"
)

// #################################################
// type definitions
// #################################################

type DocTree struct {
	XMLName  xml.Name `xml:"DOC"`
	DOCNO    string   `xml:"DOCNO"`
	HEADLINE []string `xml:"HEADLINE>P"`
	TEXT     []string `xml:"TEXT>P"`
	TABLE    []string `xml:"TEXT>TABLE>TABLEROW>TABLECELL"`
	GRAPHIC  []string `xml:"GRAPHIC>P"`
}
type Meta struct {
	Main   map[string]DocEntry
	DocNos []string
}
type DocEntry struct {
	Date     DocDate
	Headline string
	ID       int
	Length   int
}
type DocDate struct {
	Day   string
	Month string
	Year  string
	Full  string
}
type Lexicon struct {
	Words   []string
	Word2ID map[string]int
}
type InvIndex []map[int]int
type BM struct {
	DocNo string
	Score float64
}
type ScoreItem struct {
	Name  string
	Score float64
}

// #################################################
// helper functions
// #################################################

func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
func JSONLoad(file string, v interface{}) error {
	// load JSON file content into v
	j, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(j, &v)
	return err
}
func Tokenize(text string) []string {
	// break text into words
	runes := []rune(strings.ToLower(text))
	var tokens []string
	start := 0
	for i := 0; i < len(runes); i++ {
		if !unicode.In(runes[i], unicode.Digit, unicode.Letter) {
			if start < i {
				tokens = append(tokens, string(runes[start:i]))
			}
			start = i + 1
		} else if i == len(runes)-1 {
			tokens = append(tokens, string(runes[start:]))
		}
	}
	return tokens
}
