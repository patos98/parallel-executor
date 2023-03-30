package master

type ExecutableFn[T any] func(T)

type Executor[T any] interface {
	Execute(ExecutableFn[T])
}

func StartNew[T any](ctx Context[T], executor Executor[T]) {
	executor.Execute(func(todoTask T) {
		ctx.Todo <- struct{}{}
		taskChannels := <-ctx.Ready
		taskChannels.Tasks() <- todoTask
		<-taskChannels.Done()
	})
}
