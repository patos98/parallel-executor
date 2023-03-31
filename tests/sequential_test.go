package tests

import (
	"testing"
	"time"

	"github.com/patos98/parallel-executor/executor"
	"github.com/patos98/parallel-executor/master"
	"github.com/patos98/parallel-executor/worker"
)

func TestSequential(t *testing.T) {
	todo := make(chan struct{})
	defer close(todo)
	ready := make(chan master.TaskChannels[taskWithFn])
	defer close(ready)

	workerCount := 2
	workerCtx := worker.Context[taskWithFn]{Todo: todo, Ready: ready}
	done := make(map[string]struct{}, workerCount)
	startWorkers(workerCount, workerCtx, func(t taskWithFn) (result taskWithFn, err error) {
		time.Sleep(1 * time.Second)
		t.execute()
		done[t.name] = struct{}{}
		return
	})

	masterCtx := master.Context[taskWithFn]{Todo: todo, Ready: ready}
	executor := executor.NewSequential([]master.Executor[taskWithFn]{
		executor.NewSingle(taskWithFn{name: "First Task"}),
		executor.NewSingle(taskWithFn{name: "Second Task", fn: secondSequentialTaskFn(t, done)}),
	})
	master.StartNew[taskWithFn](masterCtx, executor)
}

func secondSequentialTaskFn(t *testing.T, done map[string]struct{}) func() {
	return func() {
		if _, has := done["First Task"]; !has {
			t.Error("First task has not finished yet!")
		}
	}
}
