package main

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/visig9/tf/logger"
	"github.com/visig9/tf/readline"
	"github.com/visig9/tf/textrel"
)

var version string

var files []string

var caseInsensitive bool
var ignoreFilename bool

var printZero bool

var rootCmd = &cobra.Command{
	Use:   os.Args[0] + " TERM...",
	Short: "Calculate the relevance between FILEs and TERMs.",
	Long: `Calculate the relevance between FILEs and TERMs.

  This program are designd for evaluate the relevance between
  FILEs and TERMs.
  Accept FILEs from stdin.

  Source Code: https://github.com/visig9/tf
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
		var flag textrel.Flag
		if caseInsensitive {
			flag = flag | textrel.CaseInsensitive
		}
		if !ignoreFilename {
			flag = flag | textrel.Filename
		}

		printScore(files, args, flag, printZero)
	},
	Version: version,
}

func init() {
	rootCmd.Flags().StringArrayVarP(
		&files, "file", "f", []string{},
		"the files want be evaluated",
	)
	rootCmd.Flags().BoolVarP(
		&caseInsensitive, "case-insensitive", "C", false,
		"scoring case-insensitive",
	)
	rootCmd.Flags().BoolVarP(
		&ignoreFilename, "ignore-name", "N", false,
		"don't scoring the file's name",
	)
	rootCmd.Flags().BoolVarP(
		&printZero, "include-zero", "z", false,
		"show the files which relevance is zero",
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Err.Println(err)
		os.Exit(1)
	}
}
