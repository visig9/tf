package pipe_test

import (
	"testing"

	"gitlab.com/visig/tf/pipe"
)

func echo(x interface{}) interface{} {
	return x
}

func TestPipe(t *testing.T) {
	inCounts := []int{0, 1, 5, 10}

	for _, inCount := range inCounts {
		pip := pipe.New(4, echo)

		for i := 0; i < inCount; i++ {
			pip.In() <- "hi"
		}
		close(pip.In())

		outCount := 0
		for range pip.Out() {
			outCount++
		}

		if inCount != outCount {
			t.Errorf(
				"%v (In) != %v (Out)",
				inCount,
				outCount,
			)
		}
	}
}

func TestPipeDataOrder(t *testing.T) {
	cases := []string{"hi", ",", "I", "am", "a", "good", "guy"}

	pip := pipe.New(4, echo)

	for _, c := range cases {
		pip.In() <- c
	}
	close(pip.In())

	for idx, c := range cases {
		get := <-pip.Out()

		if get != c {
			t.Errorf(
				"(<-pipe.Out) [%v] == %q, want: %q",
				idx, get, c,
			)
		}
	}
}

func TestPipePanic(t *testing.T) {
	testfunc := func(chansize int, wantPanic bool) {
		defer func() {
			err := recover()
			if isPanic := err != nil; isPanic != wantPanic {
				t.Errorf(
					"chansize: %v; panic: %v",
					chansize, isPanic,
				)
			}
		}()
		pipe.New(chansize, echo)
	}

	cases := []struct {
		chansize  int
		wantPanic bool
	}{
		{2, false},
		{1, false},
		{0, false},
		{-1, true},
		{-2, true},
	}

	for _, c := range cases {
		testfunc(c.chansize, c.wantPanic)
	}
}

func TestChain(t *testing.T) {
	cases := [][]int{
		{1, 2, 3, 4, 5},
		{10, 9, 5, 7, 8},
		{2, 4, 1},
	}

	addOne := func(x interface{}) interface{} {
		return x.(int) + 1
	}
	square := func(x interface{}) interface{} {
		return x.(int) * x.(int)
	}

	pip := pipe.New(1, addOne, square)

	for _, c := range cases {
		for _, in := range c {
			pip.In() <- in
		}
		for _, in := range c {
			out := <-pip.Out()
			if square(addOne(in)) != out {
				t.Errorf(
					"x == %v, but (x+1)^2 == %v",
					in, out,
				)
			}
		}
	}
}
