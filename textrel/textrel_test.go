package textrel_test

import (
	"math"
	"testing"

	"github.com/visig9/tf/textrel"
)

var tolerance = 0.00001

func TestByTerms(t *testing.T) {
	cases := []struct {
		text  string
		terms []string
		want  float64
	}{
		{"", []string{""}, 0},
		{"", []string{"apple"}, 0},
		{"My apple", []string{}, 0},
		{"My apple", []string{"apple"}, 5.0 / 8},
		{"My apple", []string{"apple", "apple"}, (5.0 / 8) * 1},
		{"My apple is the best!",
			[]string{"apple"},
			(5.0 / 21) * 1},
		{"My apple is the best apple",
			[]string{"apple "},
			(6.0 / 26) * 1},
		{"My apple is the best apple",
			[]string{"apple"},
			(10.0 / 26) * 2},
		{"My Apple is the best apple",
			[]string{"app"},
			(3.0 / 26) * 1},
		{"My Apple is the best apple",
			[]string{"apple", "is"},
			(5.0/26)*1 + (2.0/26)*1},
	}

	for _, c := range cases {
		s := textrel.ByTerms(c.text, c.terms)
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

func TestFileByTerms(t *testing.T) {
	cases := []struct {
		fpath   string
		terms   []string
		flag    textrel.Flag
		want    float64
		wantErr bool
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

		if c.wantErr && (err == nil) || !c.wantErr && (err != nil) {
			t.Errorf(
				"FileScore(%q, %q, %b) want err",
				c.fpath,
				c.terms,
				c.flag,
			)
		}
	}
}
