package tests

import (
	"parallel-executor/executor"
	"parallel-executor/master"
	"parallel-executor/worker"
	"sync"
	"testing"
	"time"
)

func TestParallel(t *testing.T) {
	todo := make(chan struct{})
	defer close(todo)
	ready := make(chan master.TaskChannels[taskWithFn])
	defer close(ready)

	workerCount := 2
	startedMutex := sync.Mutex{}
	started := make(map[string]struct{}, workerCount)
	doneMutex := sync.Mutex{}
	done := make(map[string]struct{}, workerCount)
	workerCtx := worker.Context[taskWithFn]{Todo: todo, Ready: ready}
	startWorkers(workerCount, workerCtx, func(t taskWithFn) (result taskWithFn, err error) {
		startedMutex.Lock()
		started[t.name] = struct{}{}
		startedMutex.Unlock()

		time.Sleep(1 * time.Second)

		t.execute()

		time.Sleep(1 * time.Second)

		doneMutex.Lock()
		done[t.name] = struct{}{}
		doneMutex.Unlock()

		return
	})

	masterCtx := master.Context[taskWithFn]{Todo: todo, Ready: ready}
	executor := executor.NewParallel([]master.Executor[taskWithFn]{
		executor.NewSingle(taskWithFn{name: "First Task", fn: firstParallelTaskFn(t, started, done)}),
		executor.NewSingle(taskWithFn{name: "Second Task", fn: secondParallelTaskFn(t, started, done)}),
	})
	master.StartNew[taskWithFn](masterCtx, executor)
}

func firstParallelTaskFn(t *testing.T, started, done map[string]struct{}) func() {
	return func() {
		if _, has := started["Second Task"]; !has {
			t.Error("Second task has not started yet!")
		}
		if _, has := done["Second Task"]; has {
			t.Error("Second task has already been done!")
		}
	}
}

func secondParallelTaskFn(t *testing.T, started, done map[string]struct{}) func() {
	return func() {
		if _, has := started["First Task"]; !has {
			t.Error("First task has not started yet!")
		}
		if _, has := done["First Task"]; has {
			t.Error("First task has already been done!")
		}
	}
}
