package executor

import "parallel-executor/master"

type SequentialExecutor[T any] struct {
	executors []master.Executor[T]
}

var _ master.Executor[any] = (*SequentialExecutor[any])(nil)

func NewSequential[T any](executors []master.Executor[T]) *SequentialExecutor[T] {
	return &SequentialExecutor[T]{executors: executors}
}

func (se *SequentialExecutor[T]) Execute(executorFn master.ExecutableFn[T]) {
	for _, executor := range se.executors {
		executor.Execute(executorFn)
	}
}
