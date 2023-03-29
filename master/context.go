package master

type TaskChannels[T any] interface {
	Tasks() chan<- T
	Done() <-chan struct{}
}

type Context[T any] struct {
	Todo  chan<- struct{}
	Ready <-chan TaskChannels[T]
}
