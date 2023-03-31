package tests

import (
	"errors"
	"parallel-executor/executor"
	"parallel-executor/master"
	"parallel-executor/worker"
	"testing"
)

func TestError(t *testing.T) {
	todo := make(chan struct{})
	defer close(todo)
	ready := make(chan master.TaskChannels[task])
	defer close(ready)

	counter := 0
	workerCtx := worker.Context[task]{Todo: todo, Ready: ready}
	startWorker(workerCtx, func(t task) (result task, err error) {
		counter++
		if counter < 4 {
			// by default master will retry to execute task with worker
			err = errors.New("counter is less than 4")
		}

		return
	})

	masterCtx := master.Context[task]{Todo: todo, Ready: ready}
	executor := executor.NewSingle(task{name: "Single Task"})
	master.StartNew(masterCtx, master.Executor[task](executor))

	if counter != 4 {
		t.Error("Single Task has not been done!")
	}
}
