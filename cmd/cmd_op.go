package main

import (
	"os"

	"gitlab.com/visig/tf/logger"
	"gitlab.com/visig/tf/pipe"
	"gitlab.com/visig/tf/readline"
	"gitlab.com/visig/tf/textrel"
)

type answer struct {
	path  string
	score float64
	err   error
}

func getPipe(terms []string, flag textrel.Flag) *pipe.Pipe {
	convert := func(x interface{}) interface{} {
		path := x.(string)
		score, err := textrel.FileByTerms(path, terms, flag)

		return answer{
			path,
			score,
			err,
		}
	}

	return pipe.New(10, convert)
}

func printScore(
	paths,
	terms []string,
	flag textrel.Flag,
	printZero bool,
) {
	pip := getPipe(terms, flag)

	for _, path := range paths {
		pip.In() <- path
	}

	if readline.IsPipe(os.Stdin) {
		ch := readline.Channel(os.Stdin)
		go func() {
			for path := range ch {
				pip.In() <- path
			}
			close(pip.In())
		}()
	} else {
		close(pip.In())
	}

	for i := range pip.Out() {
		ans := i.(answer)

		switch {
		case ans.err != nil:
			logger.Err.Printf(
				"ignore, %v\n",
				ans.err,
			)
		case printZero || ans.score > 0:
			logger.Std.Printf(
				"%13.8f %v\n",
				ans.score,
				ans.path,
			)
		}
	}
}
