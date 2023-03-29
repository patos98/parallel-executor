package executor

import "parallel-executor/master"

type SingleExecutor[T any] struct {
	task T
}

var _ master.Executor[any] = (*SingleExecutor[any])(nil)

func NewSingle[T any](task T) *SingleExecutor[T] {
	return &SingleExecutor[T]{task: task}
}

func (se *SingleExecutor[T]) Execute(executorFn master.ExecutorFn[T]) {
	executorFn(se.task)
}
