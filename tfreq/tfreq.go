// Package tfreq provide term-frequency score calculation.
package tfreq

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func dedupStrings(ss []string) []string {
	set := make(map[string]struct{})

	for _, s := range ss {
		_, exists := set[s]

		if !exists {
			set[s] = struct{}{}
		}
	}

	result := make([]string, 0, len(set))

	for s := range set {
		result = append(result, s)
	}

	return result
}

func modToLowerCase(text string, terms []string) (
	ltext string,
	lterms []string,
) {
	ltext = strings.ToLower(text)

	lterms = make([]string, len(terms))
	for idx, term := range terms {
		lterms[idx] = strings.ToLower(term)
	}

	return
}

// ScoreFlag can modify Score operation.
type ScoreFlag int

// Those flag can change the behavior of Score and FileScore
const (

	// do not lowercase before scoring.
	ScoreCaseSensitive ScoreFlag = 1 << iota

	// add score of file name in FileScore.
	ScoreFilename
)

// Score calculate and return the total term-frequency between the text
// and terms.
//
// The larger number meaning higher relevance. The 0 mean no relevance.
func Score(text string, terms []string, flag ScoreFlag) float64 {
	if flag&ScoreCaseSensitive == 0 {
		text, terms = modToLowerCase(text, terms)
	}

	terms = dedupStrings(terms)

	var termsLength float64
	textLength := float64(strings.Count(text, "") - 1)

	for _, term := range terms {
		termCount := float64(strings.Count(text, term))
		termsLength += float64(strings.Count(term, "")-1) * termCount
	}

	if termsLength == 0 || textLength == 0 {
		return 0
	}

	return termsLength / textLength
}

// FileScore calculate and return the total term-grequency between the text
// and the file.
//
// The returned score was the sum of two Score() of file content and
// filename. The 0 meaning no relevance or ReadFile() failed.
func FileScore(fpath string, terms []string, flag ScoreFlag) (score float64) {
	content, err := ioutil.ReadFile(fpath)
	if err != nil {
		return 0
	}
	score += Score(string(content), terms, flag)

	if ScoreFilename&flag != 0 {
		fname := filepath.Base(fpath)
		score += Score(fname, terms, flag)
	}

	return score
}
