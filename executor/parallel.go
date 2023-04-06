package executor

import (
	"sync"

	"github.com/patos98/parallel-executor/master"
)

type ParallelExecutor[T any] struct {
	executors []master.Executor[T]
}

var _ master.Executor[any] = (*ParallelExecutor[any])(nil)

func NewParallel[T any](executors []master.Executor[T]) *ParallelExecutor[T] {
	return &ParallelExecutor[T]{executors: executors}
}

func (pe *ParallelExecutor[T]) Execute(params master.ExecutorParams[T]) {
	var wg sync.WaitGroup
	for _, executor := range pe.executors {
		wg.Add(1)
		executor := executor // avoid re-use of the same value in each goroutine closure
		go func() {
			defer wg.Done()
			executor.Execute(params)
		}()
	}
	wg.Wait()
}
