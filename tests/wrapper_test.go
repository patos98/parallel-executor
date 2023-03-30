package tests

import (
	"parallel-executor/executor"
	"parallel-executor/master"
	"parallel-executor/worker"
	"testing"
)

func TestWrapper(t *testing.T) {
	todo := make(chan struct{})
	defer close(todo)
	ready := make(chan master.TaskChannels[taskWithFn])
	defer close(ready)

	var status string
	workerCtx := worker.Context[taskWithFn]{Todo: todo, Ready: ready}
	startWorker(workerCtx, func(task taskWithFn) {
		if status != "Task started" {
			t.Error("Wrapper function has not run!")
		}
	})

	masterCtx := master.Context[taskWithFn]{Todo: todo, Ready: ready}
	original := executor.NewSingle(taskWithFn{name: "Single Task"})
	wrapperFn := func(task taskWithFn, executorFn master.ExecutableFn[taskWithFn]) {
		status = "Task started"
		executorFn(task)
		status = "Task ended"
	}

	master.StartNew[taskWithFn](masterCtx, executor.Wrapping[taskWithFn](original, wrapperFn))

	if status != "Task ended" {
		t.Error("Wrapper function has not run!")
	}
}
