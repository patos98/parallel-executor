package worker

import "github.com/patos98/parallel-executor/master"

type Worker[T any] interface {
	work(T) (T, error)
}

type WorkerFn[T any] func(task T) (T, error)

func (workerFn WorkerFn[T]) work(task T) (T, error) {
	return workerFn(task)
}

func StartNew[T any](ctx Context[T], w Worker[T]) {
	tasks := make(chan T)
	defer close(tasks)
	done := make(chan master.TaskDoneMessage[T])
	defer close(done)

	for range ctx.Todo {
		ctx.Ready <- TaskChannels[T]{
			tasks: tasks,
			done:  done,
		}

		currentTask := <-tasks
		currentTask, err := w.work(currentTask) // TODO: decide if error is fatal and stop worker if so

		done <- master.TaskDoneMessage[T]{
			Task: currentTask,
			Err:  err,
		}
	}
}
