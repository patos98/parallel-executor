package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/patos98/parallel-executor/executor"
	"github.com/patos98/parallel-executor/master"
	"github.com/patos98/parallel-executor/worker"
)

func TestParallelOrdered(t *testing.T) {
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

		t.execute()

		time.Sleep(2 * time.Second)

		doneMutex.Lock()
		done[t.name] = struct{}{}
		doneMutex.Unlock()

		return
	})

	masterCtx := master.Context[taskWithFn]{Todo: todo, Ready: ready}
	original := executor.NewParallelOrdered([]master.Executor[taskWithFn]{
		executor.NewSingle(taskWithFn{name: "First Task", fn: firstParallelOrderedTaskFn(t, started)}),
		executor.NewSingle(taskWithFn{name: "Second Task", fn: secondParallelOrderedTaskFn(t, started, done)}),
	})

	wrapperFn := func(executableFnParams master.ExecutableFnParams[taskWithFn], executorParams master.ExecutorParams[taskWithFn]) taskWithFn {
		task := executorParams.ExecutableFn(master.ExecutableFnParams[taskWithFn]{
			Task: executableFnParams.Task,
			AfterTaskStarted: func() {
				if executableFnParams.AfterTaskStarted != nil {
					executableFnParams.AfterTaskStarted()
				}

				time.Sleep(1 * time.Second)
			},
		})
		return task
	}

	master.StartNew[taskWithFn](masterCtx, executor.Wrapping[taskWithFn](original, wrapperFn))
}

func firstParallelOrderedTaskFn(t *testing.T, started map[string]struct{}) func() {
	return func() {
		if _, has := started["Second Task"]; has {
			t.Error("Second task has already started!")
		}
	}
}

func secondParallelOrderedTaskFn(t *testing.T, started, done map[string]struct{}) func() {
	return func() {
		if _, has := started["First Task"]; !has {
			t.Error("First task has not started yet!")
		}
		if _, has := done["First Task"]; has {
			t.Error("First task has already been done!")
		}
	}
}
