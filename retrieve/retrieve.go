// Retrieve.exe path/to/index
package main

import (
	"app/bm25"
	"app/ext"
	"app/getDoc"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func SplitSents(text *string) []string {
	// fmt.Println("\u002E")
	// fmt.Println("\u003F")
	// fmt.Println("\u0021")
	// runes := []rune(*text)
	var sents []string
	start := 0
	for i := 0; i < len(*text); i++ {
		if strings.ContainsRune(string((*text)[i]), 46) {
			// fmt.Print((*text)[i])
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

func QBSummary(query *[]string, text *string) []string {
	// returns up to nMax most relevant sentences in doc
	sents := SplitSents(text)
	scores := make([]int, len(sents))
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
			kTmp := 1
			for w, word := range words {
				if word == token {
					// count matches
					c += 1
					if !found {
						// count unique matches
						found = true
						d += 1
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
		scores[i] = l + c + d + k
	}
	sort.Slice(sents, func(i, j int) bool {
		return scores[i] > scores[j]
	})
	return sents
}

func retrieve() error {
	// load data from index
	os.Args[1] = strings.TrimSuffix(os.Args[1], "/")
	fmt.Print("loading lexicon ... ")
	var lex ext.Lexicon
	lex.Word2ID = map[string]int{}
	if err := ext.JSONLoad(os.Args[1]+"/lexicon.json", &lex); err != nil {
		return err
	}
	fmt.Println("done")
	fmt.Print("loading inverted index ... ")
	var postings ext.InvIndex
	if err := ext.JSONLoad(os.Args[1]+"/postings.json", &postings); err != nil {
		return err
	}
	fmt.Println("done")
	fmt.Print("loading metadata ... ")
	var meta ext.Meta
	if err := ext.JSONLoad(os.Args[1]+"/metadata.json", &meta); err != nil {
		return err
	}
	fmt.Println("done")

	var wg sync.WaitGroup
queryLoop:
	for {
		// input query
		fmt.Print("query: ")
		var query string
		fmt.Scan(&query)

		// BM25
		tokens := ext.Tokenize(query)
		results, err := bm25.BM25(&tokens, &lex, &postings, &meta)
		ext.Check(err)

		// display top 10
		disp := ""
		for i, result := range results {
			wg.Add(1)

			text, err := getDoc.GetBody(os.Args[1], "docno", result.DocNo, &meta)
			if err != nil {
				return err
			}
			headline := meta.Main[result.DocNo].Headline
			if headline == "" {
				headline = text[:50] + "..."
			}
			sents := QBSummary(&tokens, &text)
			disp += fmt.Sprintf("%2d.\t%v\n", i+1, headline)
			for i, sent := range sents {
				disp += sent + " [...] "
				if i == 1 {
					break
				}
			}
			disp += "\n\n"
			wg.Done()
			if i >= 9 {
				break
			}
		}
		fmt.Print(disp)
		// report time
		var action string
		for {
			// prompt for rank, N, Q
			fmt.Print("make a selection: ")
			fmt.Scan(&action)
			switch action {
			case "Q", "q":
				fmt.Println("quitting ...")
				break queryLoop
			case "N", "n":
				fmt.Println("new query ...")
				continue queryLoop
			case "R", "r":
				fmt.Print(disp)
			case "1", "2", "3", "4", "5", "6", "7", "8", "9", "10":
				// show full doc
				rank, _ := strconv.Atoi(action)
				r, _ := getDoc.GetRaw(os.Args[1], "docno", results[rank+1].DocNo, &meta)
				wg.Add(1)
				fmt.Println(r)
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
	// input verification
	err := retrieve()
	if err != nil {
		fmt.Printf("\n%v\n", err)
	}
}
