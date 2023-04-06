package executor

import "github.com/patos98/parallel-executor/master"

type SingleExecutor[T any] struct {
	task T
}

var _ master.Executor[any] = (*SingleExecutor[any])(nil)

func NewSingle[T any](task T) *SingleExecutor[T] {
	return &SingleExecutor[T]{task: task}
}

func (se *SingleExecutor[T]) Execute(params master.ExecutorParams[T]) {
	params.ExecutableFn(master.ExecutableFnParams[T]{
		Task:             se.task,
		AfterTaskStarted: params.AfterTaskStarted,
	})
}
