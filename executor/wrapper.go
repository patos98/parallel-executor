package executor

import "parallel-executor/master"

type WrapperFn[T any] func(task T, executorFn master.ExecutableFn[T]) T

type wrapper[T any] struct {
	executor  master.Executor[T]
	wrapperFn WrapperFn[T]
}

func (w *wrapper[T]) Execute(executorFn master.ExecutableFn[T]) {
	w.executor.Execute(func(task T) T {
		return w.wrapperFn(task, executorFn)
	})
}

func Wrapping[T any](executor master.Executor[T], wrapperFn WrapperFn[T]) *wrapper[T] {
	return &wrapper[T]{
		executor:  executor,
		wrapperFn: wrapperFn,
	}
}
