package main

import (
	"os"
	"runtime"

	"gitlab.com/visig/tf/logger"
	"gitlab.com/visig/tf/pipe"
	"gitlab.com/visig/tf/readline"
	"gitlab.com/visig/tf/strescape"
	"gitlab.com/visig/tf/tfreq"
)

func fscore(
	path string,
	terms []string,
	flag tfreq.ScoreFlag,
) (ans float64) {
	fi, err := os.Stat(path)

	if err != nil {
		pe := err.(*os.PathError)
		logger.Err.Printf(
			"ignore '%v': %v\n",
			strescape.SingleQuote(pe.Path), pe.Err,
		)

		return
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		logger.Err.Printf(
			"ignore '%v': is a directory\n",
			strescape.SingleQuote(path),
		)
	case mode.IsRegular():
		ans = tfreq.FileScore(path, terms, flag)
	}

	return
}

type question struct {
	path  string
	terms []string
	flag  tfreq.ScoreFlag
}
type answer struct {
	q     question
	score float64
}

func getPipe() *pipe.Pipe {
	convert := func(x interface{}) interface{} {
		q := x.(question)
		return answer{
			q:     q,
			score: fscore(q.path, q.terms, q.flag),
		}
	}

	return pipe.New(runtime.NumCPU(), convert)
}

func printScore(
	paths,
	terms []string,
	flag tfreq.ScoreFlag,
	printZero bool,
) {
	pip := getPipe()

	for _, path := range paths {
		pip.In <- question{
			path:  path,
			terms: terms,
			flag:  flag,
		}
	}

	if readline.IsPipe(os.Stdin) {
		ch := readline.Channel(os.Stdin)
		go func() {
			for path := range ch {
				pip.In <- question{
					path:  path,
					terms: terms,
					flag:  flag,
				}
			}
			close(pip.In)
		}()
	} else {
		close(pip.In)
	}

	for i := range pip.Out {
		ans := i.(answer)

		if printZero || ans.score > 0 {
			logger.Std.Printf(
				"%11.9f %v\n",
				ans.score,
				ans.q.path,
			)
		}
	}
}
