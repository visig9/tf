// Package textrel provide text relevance calculation.
package textrel

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func dedup(ss []string) []string {
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

func toLowerCase(texts ...string) (lterms []string) {
	lterms = make([]string, len(texts))
	for idx, text := range texts {
		lterms[idx] = strings.ToLower(text)
	}

	return
}

// Flag can modify Score operation.
type Flag int

// Those flag can change the behavior of FileByTerms
const (
	// case insensitive scoring.
	CaseInsensitive Flag = 1 << iota

	// score of file name in FileScore.
	Filename
)

type scoreFunc func(text string, terms []string) (score float64)

// ByTerms calculate and return the total relevance between the text
// and terms.
//
// The larger number meaning higher relevance. The 0 mean no relevance.
func ByTerms(text string, terms []string) (score float64) {
	terms = dedup(terms)

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

// ByTermsCI calculate and return the total relevance between the text
// and terms. (case-insensitive)
func ByTermsCI(text string, terms []string) (score float64) {
	text = strings.ToLower(text)
	terms = toLowerCase(terms...)

	return ByTerms(text, terms)
}

// FileByTerms calculate and the relevance between the terms and the file.
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

	text := string(content)
	fname := filepath.Base(fpath)

	var scoreFn scoreFunc
	if flag&CaseInsensitive != 0 {
		scoreFn = ByTermsCI
	} else {
		scoreFn = ByTerms
	}

	if flag&Filename != 0 {
		score += scoreFn(fname, terms)
	}

	score += scoreFn(text, terms)

	return score, nil
}
