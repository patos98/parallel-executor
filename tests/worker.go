package tests

import "parallel-executor/worker"

func startWorkers[T any](workerCount int, ctx worker.Context[T], workerFn func(T)) {
	for i := 0; i < workerCount; i++ {
		startWorker(ctx, workerFn)
	}
}

func startWorker[T any](ctx worker.Context[T], workerFn func(T)) {
	go worker.StartNew[T](ctx, worker.WorkerFn[T](workerFn))
}
