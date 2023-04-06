package master

import "github.com/patos98/parallel-executor/logger"

type ExecutableFnParams[T any] struct {
	Task             T
	AfterTaskStarted func()
}

type ExecutableFn[T any] func(ExecutableFnParams[T]) T

type ExecutorParams[T any] struct {
	ExecutableFn     ExecutableFn[T]
	AfterTaskStarted func()
}

type Executor[T any] interface {
	Execute(ExecutorParams[T])
}

func StartNew[T any](ctx Context[T], executor Executor[T]) {
	executor.Execute(ExecutorParams[T]{
		ExecutableFn: func(params ExecutableFnParams[T]) T {
			for {
				logger.Info("Task todo:", params.Task)
				ctx.Todo <- struct{}{}

				taskChannels := <-ctx.Ready
				logger.Info("Worker ready for task:", params.Task)

				taskChannels.Tasks() <- params.Task
				logger.Info("Task in progress:", params.Task)

				if params.AfterTaskStarted != nil {
					params.AfterTaskStarted()
				}

				doneMsg := <-taskChannels.Done()
				if doneMsg.Err != nil {
					logger.Error(doneMsg.Err)

					// retry to execute task; may introduce retry strategies in the future
					continue
				}

				logger.Info("Task done:", doneMsg.Task)
				return doneMsg.Task
			}
		},
	})

}
