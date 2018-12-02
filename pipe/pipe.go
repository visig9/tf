// Package pipe offer ability to chain multiple function as a pipeline.
package pipe

import (
	"sync"
)

type taggedInput struct {
	id   int
	data interface{}
}

type taggedOutput struct {
	id   int
	data interface{}
}

// Pipe is core object of this package.
//
// All data in a pipe are processing asynchronously, and those data out
// from a pipe will keep it original order.
//
// Example:
//   workerNumber := 4
//
//   addOne := func(x interface{}) interface{} {
//   	return x.(int) + 1
//   }
//   square := func(x interface{}) interface{} {
//   	return x.(int) * x.(int)
//   }
//
//   pip := pipe.New(workerNumber, addOne, square)
//
//   for i := 0; i < 10; i++ {
//     pip.In <- i  // inject the data
//   }
//   close(pip.In)  // after all data inject, close this pipe.
//                  // it can free some goroutine resource and propagate
//                  // the close to pip.Out at right time.
//
//   for i := 0; i < 10; i++ {
//     out := <-pip.Out
//     fmt.Println("(%v+1)^2 == %v", i, out)  // output keep it order.
//   }
type Pipe struct {
	// In is data input channel.
	// close(pipe.In) for close whole pipe properly.
	In chan interface{}
	// Out is data output channel.
	// It will be closed automatically when In be closed and all
	// pending data retrived from this channel.
	Out     chan interface{}
	convert func(interface{}) interface{}
}

func chainFns(
	fns ...func(interface{}) interface{},
) func(interface{}) interface{} {
	chainTwoFunc := func(
		fn,
		next func(interface{}) interface{},
	) func(interface{}) interface{} {
		return func(x interface{}) interface{} {
			return next(fn(x))
		}
	}

	fn := fns[0]

	for _, next := range fns[1:] {
		fn = chainTwoFunc(fn, next)
	}
	return fn
}

// New create a Pipe instance.
func New(workers int, converts ...func(interface{}) interface{}) *Pipe {
	if workers < 1 {
		panic("Pipe's workers cannot < 1")
	}

	cSize := workers * 10

	pipe := Pipe{
		In:  make(chan interface{}, cSize),
		Out: make(chan interface{}, cSize),
	}

	convert := chainFns(converts...)

	tic := make(chan taggedInput, cSize)
	toc := make(chan taggedOutput, cSize)

	pipe.startWorkers(workers, convert, tic, toc)
	go pipe.inputRoutine(tic)
	go pipe.outputRoutine(toc, cSize)

	return &pipe
}

func (pipe *Pipe) startWorkerRoutine(
	wg *sync.WaitGroup,
	convert func(interface{}) interface{},
	tic chan taggedInput,
	toc chan taggedOutput,
) {
	wg.Add(1)

	// Worker routines
	go func() {
		defer wg.Done()
		for tq := range tic {
			toc <- taggedOutput{
				id:   tq.id,
				data: convert(tq.data),
			}
		}
	}()
}

func (pipe *Pipe) startWorkers(
	workers int,
	convert func(interface{}) interface{},
	tic chan taggedInput,
	toc chan taggedOutput,
) {
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		pipe.startWorkerRoutine(&wg, convert, tic, toc)
	}

	// toc closing routine
	go func() {
		wg.Wait()

		close(toc)
	}()
}

func (pipe *Pipe) inputRoutine(tic chan taggedInput) {
	count := 0
	for data := range pipe.In {
		tic <- taggedInput{
			id:   count,
			data: data,
		}
		count++
	}

	close(tic)
}

func (pipe *Pipe) outputRoutine(toc chan taggedOutput, cSize int) {
	var next int
	buf := make(map[int]interface{}, cSize)

	for ta := range toc {
		buf[ta.id] = ta.data

		for data, hit := buf[next]; hit; data, hit = buf[next] {
			pipe.Out <- data

			delete(buf, next)

			next++
		}
	}

	close(pipe.Out)
}
