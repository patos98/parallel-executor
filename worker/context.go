package worker

import "parallel-executor/master"

type TaskChannels[T any] struct {
	tasks chan<- T
	done  <-chan error
}

func (tc TaskChannels[T]) Tasks() chan<- T    { return tc.tasks }
func (tc TaskChannels[T]) Done() <-chan error { return tc.done }

type Context[T any] struct {
	Todo  <-chan struct{}
	Ready chan<- master.TaskChannels[T]
}
