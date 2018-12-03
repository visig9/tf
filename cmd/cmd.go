package main

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"gitlab.com/visig/tf/logger"
	"gitlab.com/visig/tf/readline"
	"gitlab.com/visig/tf/tfreq"
)

var version string

var files []string

var caseInsensitive bool
var ignoreFilename bool

var printZero bool

var rootCmd = &cobra.Command{
	Use:   os.Args[0] + " TERM...",
	Short: "Calculate term-frequency of files.",
	Long: `Calculate term-frequency of files.

  Calculate the term-frequency between TERMs and FILEs.
  This program also accept the FILEs from stdin.

  Source Code: https://gitlab.com/visig/tf
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one term")
		}

		if !readline.IsPipe(os.Stdin) && len(files) == 0 {
			return errors.New("requires at least one file")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var flag tfreq.ScoreFlag
		if !caseInsensitive {
			flag = flag | tfreq.ScoreCaseSensitive
		}
		if !ignoreFilename {
			flag = flag | tfreq.ScoreFilename
		}

		printScore(files, args, flag, printZero)
	},
	Version: version,
}

func init() {
	rootCmd.Flags().StringArrayVarP(
		&files, "file", "f", []string{},
		"files want be evaluated",
	)
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
