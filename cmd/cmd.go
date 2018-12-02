package main

import (
	"os"

	"github.com/spf13/cobra"

	"gitlab.com/visig/tf/logger"
	"gitlab.com/visig/tf/tfreq"
)

var version string

var terms []string

var caseInsensitive bool
var ignoreFilename bool

var printZero bool

var rootCmd = &cobra.Command{
	Use:   os.Args[0] + " FILE...",
	Short: "Calculate term-frequency of files.",
	Long: `Calculate term-frequency of files.

  Calculate the term-frequency between FILE and TERMs.
  Accept multiple -t TERMs and multiple FILEs at one time.
  This program also accept the FILEs from stdin.

  Source Code: https://gitlab.com/visig/tf
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var flag tfreq.ScoreFlag
		if !caseInsensitive {
			flag = flag | tfreq.ScoreCaseSensitive
		}
		if !ignoreFilename {
			flag = flag | tfreq.ScoreFilename
		}

		printScore(args, terms, flag, printZero)
	},
	Version: version,
}

func init() {
	rootCmd.Flags().StringArrayVarP(
		&terms, "term", "t", []string{},
		"calculate by this key terms",
	)
	rootCmd.MarkFlagRequired("term")
	rootCmd.Flags().BoolVarP(
		&caseInsensitive, "case-insensitive", "C", false,
		"scoring by case-insensitive",
	)
	rootCmd.Flags().BoolVarP(
		&ignoreFilename, "ignore-name", "N", false,
		"don't scoring the file's name",
	)
	rootCmd.Flags().BoolVarP(
		&printZero, "include-zero", "z", false,
		"don't omit the file which relevance is zero",
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Err.Println(err)
		os.Exit(1)
	}
}
