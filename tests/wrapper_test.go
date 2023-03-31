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
	ready := make(chan master.TaskChannels[task])
	defer close(ready)

	var status string
	workerCtx := worker.Context[task]{Todo: todo, Ready: ready}
	startWorker(workerCtx, func(task task) (task, error) {
		if status != "Task started" {
			t.Error("Wrapper function has not run!")
		}
		task.name = "Single Task Done"
		return task, nil
	})

	masterCtx := master.Context[task]{Todo: todo, Ready: ready}
	original := executor.NewSingle(task{name: "Single Task"})
	wrapperFn := func(task task, executorFn master.ExecutableFn[task]) task {
		status = "Task started"
		task = executorFn(task)
		if task.name != "Single Task Done" {
			t.Errorf("ExecutorFn returned task in wrong state: %#v", task)
		}

		status = "Task ended"
		return task
	}

	master.StartNew[task](masterCtx, executor.Wrapping[task](original, wrapperFn))

	if status != "Task ended" {
		t.Error("Wrapper function has not run!")
	}
}
