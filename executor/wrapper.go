package executor

import "github.com/patos98/parallel-executor/master"

type WrapperFn[T any] func(executableFnParams master.ExecutableFnParams[T], executorParams master.ExecutorParams[T]) T

type wrapper[T any] struct {
	executor  master.Executor[T]
	wrapperFn WrapperFn[T]
}

func (w *wrapper[T]) Execute(executorParams master.ExecutorParams[T]) {
	w.executor.Execute(master.ExecutorParams[T]{
		ExecutableFn: func(executableFnParams master.ExecutableFnParams[T]) T {
			return w.wrapperFn(executableFnParams, executorParams)
		},
		AfterTaskStarted: executorParams.AfterTaskStarted,
	})
}

func Wrapping[T any](executor master.Executor[T], wrapperFn WrapperFn[T]) *wrapper[T] {
	return &wrapper[T]{
		executor:  executor,
		wrapperFn: wrapperFn,
	}
}
