package master

type TaskDoneMessage[T any] struct {
	Task T
	Err  error
}

type TaskChannels[T any] interface {
	Tasks() chan<- T
	Done() <-chan TaskDoneMessage[T]
}

type Context[T any] struct {
	Todo  chan<- struct{}
	Ready <-chan TaskChannels[T]
}
