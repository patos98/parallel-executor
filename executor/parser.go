package executor

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/patos98/parallel-executor/master"
)

const (
	EXECUTOR_TYPE_SEQUENTIAL       = "SEQUENTIAL"
	EXECUTOR_TYPE_PARALLEL         = "PARALLEL"
	EXECUTOR_TYPE_PARALLEL_ORDERED = "PARALLEL_ORDERED"
)

type genericExecutor[T any] struct {
	ExecutorType string               `json:"type"`
	Task         T                    `json:"task"`
	Executors    []genericExecutor[T] `json:"executors"`
}

func FromJson[T comparable](content []byte) (master.Executor[T], error) {
	var executor genericExecutor[T]
	err := json.Unmarshal(content, &executor)
	if err != nil {
		return nil, err
	}

	return convertGenericExecutor(executor)
}

func convertGenericExecutor[T comparable](executor genericExecutor[T]) (result master.Executor[T], err error) {
	var emptyTask T
	if executor.Task != emptyTask {
		result = NewSingle(executor.Task)
		return
	}

	executors, err := convertGenericExecutors(executor.Executors)
	if err != nil {
		return
	}

	switch executor.ExecutorType {
	case EXECUTOR_TYPE_SEQUENTIAL:
		result = NewSequential(executors)
	case EXECUTOR_TYPE_PARALLEL:
		result = NewParallel(executors)
	case EXECUTOR_TYPE_PARALLEL_ORDERED:
		result = NewParallelOrdered(executors)
	default:
		err = errors.New(fmt.Sprint("Unknown executor type: ", executor.ExecutorType))
	}

	return
}

func convertGenericExecutors[T comparable](executors []genericExecutor[T]) (result []master.Executor[T], err error) {
	for _, executor := range executors {
		var e master.Executor[T]
		e, err = convertGenericExecutor(executor)
		if err != nil {
			return
		}
		result = append(result, e)
	}
	return
}
