package textrel_test

import (
	"math"
	"testing"

	"gitlab.com/visig/tf/textrel"
)

var tolerance = 0.00001

func TestScore(t *testing.T) {
	cases := []struct {
		text  string
		terms []string
		flag  textrel.Flag
		want  float64
	}{
		{"", []string{""}, 0, 0},
		{"", []string{"apple"}, 0, 0},
		{"My apple", []string{}, 0, 0},
		{"My apple", []string{"apple"}, 0, 5.0 / 8},
		{"My apple",
			[]string{"apple", "apple"},
			0,
			(5.0 / 8) * 1},
		{"My apple is the best!",
			[]string{"apple"},
			0,
			(5.0 / 21) * 1},
		{"My apple is the best apple",
			[]string{"apple "},
			0,
			(6.0 / 26) * 1},
		{"My apple is the best apple",
			[]string{"apple"},
			0,
			(10.0 / 26) * 2},
		{"My Apple is the best apple",
			[]string{"app"},
			textrel.CaseInsensitive,
			(6.0 / 26) * 2},
		{"My Apple is the best apple",
			[]string{"app"},
			0,
			(3.0 / 26) * 1},
		{"My Apple is the best apple",
			[]string{"apple", "is"},
			0,
			(5.0/26)*1 + (2.0/26)*1},
	}

	for _, c := range cases {
		s := textrel.ByTerms(c.text, c.terms, c.flag)
		if math.Abs(c.want-s) > tolerance {
			t.Errorf(
				"Score(%q, %q) == %v, want: %v",
				c.text,
				c.terms,
				s,
				c.want,
			)
		}
	}
}

func TestFileScore(t *testing.T) {
	cases := []struct {
		fpath  string
		terms  []string
		flag   textrel.Flag
		want   float64
		haserr bool
	}{
		{"testdata/apple-not-exists.txt", []string{"apple"},
			0,
			0,
			true},
		{"testdata/apple.txt", []string{""},
			0,
			0,
			false},
		{"testdata/apple.txt", []string{"apple"},
			textrel.Filename,
			(5.0/9)*1 + (10.0/28)*2,
			false},
		{"testdata/apple.txt", []string{"apple"},
			0,
			(10.0 / 28) * 2,
			false},
		{"testdata/apple.txt", []string{"app"},
			textrel.Filename,
			(3.0/9)*1 + (6.0/28)*2,
			false},
		{"testdata/apple.txt", []string{"best"},
			textrel.Filename,
			0.0/9 + (4.0/28)*1,
			false},
	}

	for _, c := range cases {
		s, err := textrel.FileByTerms(c.fpath, c.terms, c.flag)
		if math.Abs(c.want-s) > tolerance {
			t.Errorf(
				"FileScore(%q, %q, %b) == %v, want: %v",
				c.fpath,
				c.terms,
				c.flag,
				s,
				c.want,
			)
		}

		if c.haserr && (err == nil) || !c.haserr && (err != nil) {
			t.Errorf(
				"FileScore(%q, %q, %b) want err",
				c.fpath,
				c.terms,
				c.flag,
			)
		}
	}
}
