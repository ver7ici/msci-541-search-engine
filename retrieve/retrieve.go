// ./retrieve.exe path/to/index
package main

import (
	"app/bm25"
	"app/ext"
	"app/getDoc"
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func SplitSents(text *string) []string {
	// split text into sentences
	var sents []string
	start := 0
	for i := 0; i < len(*text); i++ {
		if strings.ContainsRune(string((*text)[i]), 46) ||
			strings.ContainsRune(string((*text)[i]), 63) ||
			strings.ContainsRune(string((*text)[i]), 33) {
			if start < i {
				sents = append(sents, strings.TrimSpace((*text)[start:i+1]))
			}
			start = i + 1
		} else if i == len(*text)-1 {
			sents = append(sents, strings.TrimSpace((*text)[start:]))
		}
	}
	return sents
}

func QBSummary(query *[]string, text *string) []ext.ScoreItem {
	// ranks sentences by relevance to query
	sents := SplitSents(text)
	scores := make([]ext.ScoreItem, len(sents))
	for i, sent := range sents {
		// 1st: l = 2, 2nd: l = 1, else: l = 0
		l := 0
		if i < 2 {
			l += 2 - i
		}
		// c: number of sentence words that are query terms including repeats
		c := 0
		// d: number of distinct query terms in sentence
		d := 0
		// k: number of contiguous words in query order
		k := 0
		words := ext.Tokenize(strings.ToLower(sent))
		for t, token := range *query {
			found := false
			kTmp := 0
			for w, word := range words {
				if word == token {
					// count matches
					c += 1
					if !found {
						// count unique matches
						found = true
						d += 1
						if kTmp == 0 {
							kTmp += 1
						}
					}
					if t > 0 && w > 0 && words[w-1] == (*query)[t-1] {
						// count contiguous matches
						kTmp += 1
					}
				}
			}
			if kTmp > k {
				k = kTmp
			}
		}
		// no weighting
		scores[i].Name = sent
		scores[i].Score = float64(l + c + d + k)
	}
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})
	return scores
}

func retrieve(indexDir string) error {
	// load data from index
	fmt.Print("loading lexicon ... ")
	var lex ext.Lexicon
	lex.Word2ID = map[string]int{}
	if err := ext.JSONLoad(indexDir+"/lexicon.json", &lex); err != nil {
		return err
	}
	fmt.Println("done")
	fmt.Print("loading inverted index ... ")
	var postings ext.InvIndex
	if err := ext.JSONLoad(indexDir+"/postings.json", &postings); err != nil {
		return err
	}
	fmt.Println("done")
	fmt.Print("loading metadata ... ")
	var meta ext.Meta
	if err := ext.JSONLoad(indexDir+"/metadata.json", &meta); err != nil {
		return err
	}
	fmt.Println("done")

	var wg sync.WaitGroup
queryLoop:
	for {
		// input query
		fmt.Print("\nquery: ")
		inQuery := bufio.NewReader(os.Stdin)
		query, err := inQuery.ReadString('\n')
		ext.Check(err)
		start := time.Now()

		// BM25 retrieval
		tokens := ext.Tokenize(query)
		results, err := bm25.BM25(&tokens, &lex, &postings, &meta)
		ext.Check(err)

		// display top 10
		disp := "\n"
		for i, result := range results {
			wg.Add(1)
			text, err := getDoc.GetBody(indexDir, "docno", result.Name, &meta)
			if err != nil {
				return err
			}
			headline := meta.Main[result.Name].Headline
			if headline == "" {
				headline = text[:50] + "..."
			}
			sents := QBSummary(&tokens, &text)
			disp += fmt.Sprintf(
				"%2d. %v (%v/%v/%v)\n",
				i+1,
				headline,
				meta.Main[result.Name].Date.Day,
				meta.Main[result.Name].Date.Month,
				meta.Main[result.Name].Date.Year,
			)
			for i, sent := range sents {
				disp += sent.Name + " "
				if i == 1 {
					break
				}
			}
			disp += fmt.Sprintf("(%v)\n\n", result.Name)
			wg.Done()
			if i >= 9 {
				break
			}
		}
		fmt.Print(disp)
		fmt.Printf("Retrieval took %v\n\n", time.Since(start))

		// prompt for rank, N, Q
		var action string
		for {
			fmt.Print("make a selection: ")
			fmt.Scan(&action)
			ext.Check(err)
			switch action {
			case "Q", "q":
				fmt.Println("quit")
				break queryLoop
			case "N", "n":
				fmt.Println("new query ...")
				continue queryLoop
			case "R", "r":
				fmt.Print(disp)
			case "1", "2", "3", "4", "5", "6", "7", "8", "9", "10":
				// show full doc
				rank, _ := strconv.Atoi(action)
				r, _ := getDoc.GetRaw(indexDir, "docno", results[rank-1].Name, &meta)
				wg.Add(1)
				fmt.Println("\n" + r)
				wg.Done()
			default:
				fmt.Print("Q:\tquit\n" +
					"N:\tenter new query\n" +
					"R:\tshow results list\n" +
					"1-10:\tshow full document\n")
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		msg := "invalid number of input arguments\n" +
			"usage:\t./retrieve.exe indexDir\n" +
			"\tindexDir:\tpath to index directory"
		fmt.Println(msg)
		return
	}
	err := retrieve(strings.TrimSuffix(os.Args[1], "/"))
	if err != nil {
		fmt.Printf("\n%v\n", err)
	}
}
