// Package textrel provide text relevance calculation.
package textrel

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

func toLowerCase(text string, terms []string) (
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

// Flag can modify Score operation.
type Flag int

// Those flag can change the behavior of Score and FileScore
const (

	// case insensitive scoring.
	CaseInsensitive Flag = 1 << iota

	// add score of file name in FileScore.
	Filename
)

// ByTerms calculate and return the total relevance between the text
// and terms.
//
// The larger number meaning higher relevance. The 0 mean no relevance.
func ByTerms(text string, terms []string, flag Flag) (score float64) {
	if flag&CaseInsensitive != 0 {
		text, terms = toLowerCase(text, terms)
	}

	terms = dedupStrings(terms)

	textLength := float64(strings.Count(text, "") - 1)

	for _, term := range terms {
		termCount := float64(strings.Count(text, term))
		termsLength := float64(strings.Count(term, "")-1) * termCount

		if termsLength != 0 && textLength != 0 {
			score += (termsLength / textLength) * termCount
		}
	}

	return
}

// FileByTerm calculate and the relevance between the terms and the file.
//
// The err != nil if file read fail.
func FileByTerms(fpath string, terms []string, flag Flag) (
	score float64,
	err error,
) {
	content, err := ioutil.ReadFile(fpath)
	if err != nil {
		return 0, err
	}
	score += ByTerms(string(content), terms, flag)

	if flag&Filename != 0 {
		fname := filepath.Base(fpath)
		score += ByTerms(fname, terms, flag)
	}

	return score, nil
}
