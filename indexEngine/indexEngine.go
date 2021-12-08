// ./indexEngine.exe path/to/latimes.gz path/to/index
package main

import (
	"app/ext"
	"bufio"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func updateLex(word string, docID int, lex *ext.Lexicon, postings *ext.InvIndex) {
	// add new word or update count of existing word per doc
	if _, ok := lex.Word2ID[word]; ok {
		(*postings)[lex.Word2ID[word]][docID]++
	} else {
		lex.Words = append(lex.Words, word)
		lex.Word2ID[word] = len(lex.Words) - 1
		*postings = append(*postings, map[int]int{docID: 1})
	}
}

func Index(srcGZ string, indexDir string, rerun bool) {
	start := time.Now()

	indexDir = strings.TrimSuffix(indexDir, "/")
	f, err := os.Open(srcGZ)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	ext.Check(err)
	defer gz.Close()

	if !rerun {
		err = os.Mkdir(indexDir, 0777)
		if err != nil {
			fmt.Print(err)
			return
		}
	}

	scanner := bufio.NewScanner(gz)
	id := 0
	var doc string
	var meta ext.Meta
	meta.DocNos = []string{}
	meta.Main = map[string]ext.DocEntry{}
	var lex ext.Lexicon
	lex.Word2ID = map[string]int{}
	var postings ext.InvIndex
	fmt.Println("ID\tDocNo\t\tElapsed\t\tRAM")

	// iterates through each line in gzip file
	for scanner.Scan() {
		line := scanner.Text()
		doc += line + "\n"
		if line == "</DOC>" {
			var docXML ext.DocTree
			xml.Unmarshal([]byte(doc), &docXML)
			docno := strings.TrimSpace(docXML.DOCNO)
			mmddyy, _ := strconv.Atoi(docno[2:8])
			yy := mmddyy % 100
			dd := ((mmddyy - yy) % 10000) / 100
			mm := (mmddyy - yy - dd*100) / 10000
			yyStr := docno[6:8]
			mmStr := docno[2:4]
			ddStr := docno[4:6]
			date := fmt.Sprintf("%v %v, %v", time.Month(mm).String(), dd, 1900+yy)
			headline := strings.TrimSpace(strings.Join(docXML.HEADLINE, ""))
			re := regexp.MustCompile(`\n`)
			headline = re.ReplaceAllString(headline, "")

			textContent := strings.Join(docXML.HEADLINE, "") +
				strings.Join(docXML.TEXT, "") +
				strings.Join(docXML.TABLE, "") +
				strings.Join(docXML.GRAPHIC, "")

			meta.DocNos = append(meta.DocNos, docno)

			tokens := ext.Tokenize(textContent)
			for _, token := range tokens {
				updateLex(token, id, &lex, &postings)
			}
			meta.Main[docno] = ext.DocEntry{
				Date: ext.DocDate{
					Year:  yyStr,
					Month: mmStr,
					Day:   ddStr,
					Full:  date,
				},
				Headline: headline,
				ID:       id,
				Length:   len(tokens),
			}

			if !rerun {
				// create appropriate path for doc based on date
				_ = os.Mkdir(indexDir+"/"+yyStr, 0777)
				_ = os.Mkdir(indexDir+"/"+yyStr+"/"+mmStr, 0777)
				_ = os.Mkdir(indexDir+"/"+yyStr+"/"+mmStr+"/"+ddStr, 0777)
				err = os.WriteFile(indexDir+"/"+yyStr+"/"+mmStr+"/"+ddStr+"/"+docno+".xml", []byte(doc), 0644)
				ext.Check(err)
			}

			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf("\r%6d\t%13v\t%12v\t%4dMb", id, docno, time.Since(start), m.Alloc/1024/1024)
			doc = ""
			id++
		}
	}
	b, err := json.MarshalIndent(meta, "", "  ")
	ext.Check(err)
	err = ioutil.WriteFile(fmt.Sprintf("%v/metadata.json", indexDir), b, 0644)
	ext.Check(err)

	b, err = json.MarshalIndent(lex, "", "  ")
	ext.Check(err)
	err = ioutil.WriteFile(fmt.Sprintf("%v/lexicon.json", indexDir), b, 0644)
	ext.Check(err)

	b, err = json.MarshalIndent(postings, "", "  ")
	ext.Check(err)
	err = ioutil.WriteFile(fmt.Sprintf("%v/postings.json", indexDir), b, 0644)
	ext.Check(err)

	fmt.Printf("\nindex created at: %v\n", indexDir)
}

func main() {
	rerun := false
	if len(os.Args) == 4 && os.Args[3] == "-r" {
		rerun = true
	} else if len(os.Args) != 3 {
		msg := "invalid number of input arguments\n" +
			"usage:\t./indexEngine.exe sourceFile indexDir\n" +
			"\tsourceFile:\tgzip file containing documents\n" +
			"\tindexDir:\tdirectory where documents and metadata will be stored"
		fmt.Println(msg)
		return
	}
	Index(os.Args[1], os.Args[2], rerun)

}
