package bm25

import (
	"app/ext"
	"math"
	"sort"
)

func UniqueInts(input []int) []int {
	// return unique set of integers, maintains order
	var u []int
	m := map[int]bool{}
	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

func BM25(tokenStr *[]string, lex *ext.Lexicon, postings *ext.InvIndex, meta *ext.Meta) ([]ext.ScoreItem, error) {
	// BM25 retrieval, no stemming
	var tokens []int
	for _, token := range *tokenStr {
		if v, ok := (*lex).Word2ID[token]; ok {
			tokens = append(tokens, v)
		}
	}
	// paraters
	k1 := 1.2
	b := 0.75
	k2 := 7.0
	N := float64(len((*meta).DocNos))
	avdl := 0.0
	for _, d := range (*meta).Main {
		avdl += float64(d.Length) / N
	}
	scores := map[int]float64{}
	for _, token := range UniqueInts(tokens) {
		for docID, count := range (*postings)[token] {
			f := float64(count)
			qf := 0.0
			for _, w := range tokens {
				if w == token {
					qf += 1
				}
			}
			dl := float64((*meta).Main[(*meta).DocNos[docID]].Length)
			K := k1 * ((1 - b) + b*dl/avdl)
			n := float64(len((*postings)[token]))
			partial := (k1 + 1) * f / (K + f) * (k2 + 1) * qf / (k2 + qf) * math.Log((N-n+0.5)/(n+0.5))
			if _, ok := scores[docID]; !ok {
				scores[docID] = partial
			} else {
				scores[docID] += partial
			}
		}
	}
	// create result structure
	var results []ext.ScoreItem
	for docID, score := range scores {
		results = append(results, ext.ScoreItem{Name: (*meta).DocNos[docID], Score: score})
	}
	// sort results by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return results, nil
}
