package worker

import "github.com/patos98/parallel-executor/master"

type TaskChannels[T any] struct {
	tasks chan<- T
	done  <-chan master.TaskDoneMessage[T]
}

func (tc TaskChannels[T]) Tasks() chan<- T                        { return tc.tasks }
func (tc TaskChannels[T]) Done() <-chan master.TaskDoneMessage[T] { return tc.done }

type Context[T any] struct {
	Todo  <-chan struct{}
	Ready chan<- master.TaskChannels[T]
}
