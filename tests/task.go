package tests

type task struct {
	name string
}

type taskWithFn struct {
	name string
	fn   func()
}

func (t taskWithFn) execute() {
	if t.fn != nil {
		t.fn()
	}
}
