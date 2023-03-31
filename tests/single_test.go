package tests

import (
	"parallel-executor/executor"
	"parallel-executor/master"
	"parallel-executor/worker"
	"testing"
)

func TestSingle(t *testing.T) {
	todo := make(chan struct{})
	defer close(todo)
	ready := make(chan master.TaskChannels[task])
	defer close(ready)

	done := false
	workerCtx := worker.Context[task]{Todo: todo, Ready: ready}
	startWorker(workerCtx, func(t task) (result task, err error) {
		done = true
		return
	})

	masterCtx := master.Context[task]{Todo: todo, Ready: ready}
	executor := executor.NewSingle(task{name: "Single Task"})
	master.StartNew(masterCtx, master.Executor[task](executor))

	if !done {
		t.Error("Single Task has not been done!")
	}
}
