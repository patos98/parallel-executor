package worker

type Worker[T any] interface {
	work(T)
}

type WorkerFn[T any] func(task T)

func (workerFn WorkerFn[T]) work(task T) {
	workerFn(task)
}

func StartNew[T any](ctx Context[T], w Worker[T]) {
	tasks := make(chan T)
	defer close(tasks)
	done := make(chan struct{})
	defer close(done)

	for range ctx.Todo {
		ctx.Ready <- TaskChannels[T]{
			tasks: tasks,
			done:  done,
		}
		currentTask := <-tasks
		w.work(currentTask)
		done <- struct{}{}
	}
}
