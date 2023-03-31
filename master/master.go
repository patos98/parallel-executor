package master

import "parallel-executor/logger"

type ExecutableFn[T any] func(T) T

type Executor[T any] interface {
	Execute(ExecutableFn[T])
}

func StartNew[T any](ctx Context[T], executor Executor[T]) {
	executor.Execute(func(task T) T {
		for {
			logger.Info("Task todo:", task)
			ctx.Todo <- struct{}{}

			taskChannels := <-ctx.Ready
			logger.Info("Worker ready for task:", task)

			taskChannels.Tasks() <- task
			logger.Info("Task in progress:", task)

			doneMsg := <-taskChannels.Done()
			if doneMsg.Err != nil {
				logger.Error(doneMsg.Err)

				// retry to execute task; may introduce retry strategies in the future
				continue
			}

			logger.Info("Task done:", doneMsg.Task)
			return doneMsg.Task
		}
	})
}
