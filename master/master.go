package master

type ExecutorFn[T any] func(T)

type Executor[T any] interface {
	Execute(ExecutorFn[T])
}

func StartNew[T any](ctx Context[T], executor Executor[T]) {
	executor.Execute(func(todoTask T) {
		ctx.Todo <- struct{}{}
		taskChannels := <-ctx.Ready
		taskChannels.Tasks() <- todoTask
		<-taskChannels.Done()
	})
}
