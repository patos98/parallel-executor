package worker

type Worker[T any] interface {
	work(T) error
}

type WorkerFn[T any] func(task T) error

func (workerFn WorkerFn[T]) work(task T) error {
	return workerFn(task)
}

func StartNew[T any](ctx Context[T], w Worker[T]) {
	tasks := make(chan T)
	defer close(tasks)
	done := make(chan error)
	defer close(done)

	for range ctx.Todo {
		ctx.Ready <- TaskChannels[T]{
			tasks: tasks,
			done:  done,
		}
		currentTask := <-tasks
		err := w.work(currentTask)
		done <- err
	}
}
