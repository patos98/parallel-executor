package executor

import (
	"sync"

	"github.com/patos98/parallel-executor/master"
)

type ParallelOrderedExecutor[T any] struct {
	executors []master.Executor[T]
}

var _ master.Executor[any] = (*ParallelOrderedExecutor[any])(nil)

func NewParallelOrdered[T any](executors []master.Executor[T]) *ParallelOrderedExecutor[T] {
	return &ParallelOrderedExecutor[T]{executors: executors}
}

func (pe *ParallelOrderedExecutor[T]) Execute(params master.ExecutorParams[T]) {
	var wg sync.WaitGroup
	for _, executor := range pe.executors {
		wg.Add(1)
		executor := executor // avoid re-use of the same value in each goroutine closure
		started := make(chan struct{})
		go func() {
			defer wg.Done()
			executor.Execute(master.ExecutorParams[T]{
				ExecutableFn: params.ExecutableFn,
				AfterTaskStarted: func() {
					if params.AfterTaskStarted != nil {
						params.AfterTaskStarted()
					}

					started <- struct{}{}
				},
			})
		}()
		<-started
	}
	wg.Wait()
}
