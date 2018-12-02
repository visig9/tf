package tfreq_test

import (
	"math"
	"testing"

	"gitlab.com/visig/tf/tfreq"
)

func TestScore(t *testing.T) {
	cases := []struct {
		text  string
		terms []string
		flag  tfreq.ScoreFlag
		want  float64
	}{
		{"", []string{""}, 0, 0},
		{"", []string{"apple"}, 0, 0},
		{"My apple", []string{}, 0, 0},
		{"My apple", []string{"apple"}, 0, 5.0 / 8},
		{"My apple",
			[]string{"apple", "apple"},
			0,
			5.0 / 8},
		{"My apple is the best!",
			[]string{"apple"},
			0,
			5.0 / 21},
		{"My apple is the best apple",
			[]string{"apple "},
			0,
			6.0 / 26},
		{"My apple is the best apple",
			[]string{"apple"},
			0,
			10.0 / 26},
		{"My Apple is the best apple",
			[]string{"app"},
			0,
			6.0 / 26},
		{"My Apple is the best apple",
			[]string{"app"},
			tfreq.ScoreCaseSensitive,
			3.0 / 26},
		{"My Apple is the best apple",
			[]string{"apple", "is"},
			0,
			12.0 / 26},
	}

	for _, c := range cases {
		if ans := tfreq.Score(c.text, c.terms, c.flag); ans != c.want {
			t.Errorf(
				"Score(%q, %q) == %v, want: %v",
				c.text,
				c.terms,
				ans,
				c.want,
			)
		}
	}
}

func TestFileScore(t *testing.T) {
	cases := []struct {
		fpath string
		terms []string
		flag  tfreq.ScoreFlag
		want  float64
	}{
		{"testdata/apple-not-exists.txt", []string{"apple"}, 0, 0},
		{"testdata/apple.txt", []string{""}, 0, 0},
		{"testdata/apple.txt",
			[]string{"apple"},
			tfreq.ScoreFilename,
			5.0/9 + 10.0/28},
		{"testdata/apple.txt",
			[]string{"apple"},
			0,
			10.0 / 28},
		{"testdata/apple.txt",
			[]string{"app"},
			tfreq.ScoreFilename,
			3.0/9 + 6.0/28},
		{"testdata/apple.txt",
			[]string{"best"},
			tfreq.ScoreFilename,
			0.0/9 + 4.0/28},
	}
	tolerance := 0.00001

	for _, c := range cases {
		s := tfreq.FileScore(c.fpath, c.terms, c.flag)
		if math.Abs(c.want-s) > tolerance {
			t.Errorf(
				"FileScore(%q, %q) == %v, want: %v",
				c.fpath,
				c.terms,
				s, c.want,
			)
		}
	}
}
