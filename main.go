package main

import (
	"log"
	"parallel-executor/executor"
	"parallel-executor/master"
	"parallel-executor/worker"
	"time"
)

const (
	MSG_START = "[START]"
	MSG_END   = "[END]  "
)

type Task struct {
	name string
	ctx  any
}

func main() {
	todo := make(chan struct{})
	defer close(todo)
	ready := make(chan master.TaskChannels[Task])
	defer close(ready)

	workerCount := 5
	workerCtx := worker.Context[Task]{Todo: todo, Ready: ready}
	createWorkers(workerCount, workerCtx)

	masterCtx := master.Context[Task]{Todo: todo, Ready: ready}
	executors := createExecutors()
	master.StartNew(masterCtx, executors)
}

func createExecutors() (tasks master.Executor[Task]) {
	tasks = executor.NewSequential([]master.Executor[Task]{
		executor.NewSingle(Task{name: "Task1", ctx: nil}),
		executor.NewSingle(Task{name: "Task2", ctx: nil}),
		executor.NewParallel([]master.Executor[Task]{
			executor.NewSequential([]master.Executor[Task]{
				executor.NewSingle(Task{name: "Task3", ctx: nil}),
				executor.NewSingle(Task{name: "Task4", ctx: nil}),
			}),
			executor.NewSequential([]master.Executor[Task]{
				executor.NewSingle(Task{name: "Task5", ctx: nil}),
				executor.NewSingle(Task{name: "Task6", ctx: nil}),
			}),
		}),
		executor.NewParallel([]master.Executor[Task]{
			executor.NewSingle(Task{name: "Task7", ctx: nil}),
			executor.NewSingle(Task{name: "Task8", ctx: nil}),
			executor.NewSingle(Task{name: "Task9", ctx: nil}),
			executor.NewSingle(Task{name: "Task10", ctx: nil}),
		}),
	})
	return
}

func createWorkers(workerCount int, ctx worker.Context[Task]) {
	var w worker.Worker[Task] = worker.WorkerFn[Task](doTask)
	for i := 0; i < workerCount; i++ {
		go worker.StartNew(ctx, w)
	}
}

func doTask(task Task) {
	log.Println(MSG_START, task.name)
	time.Sleep(1 * time.Second)
	log.Println(MSG_END, task.name)
}
