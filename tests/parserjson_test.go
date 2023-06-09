package tests

import (
	"reflect"
	"testing"

	"github.com/patos98/parallel-executor/executor"
	"github.com/patos98/parallel-executor/master"
)

func TestParserJson(t *testing.T) {
	executorJson := `{
		"type": "SEQUENTIAL",
		"executors": [
			{ "task": "Task 1" },
			{ "task": "Task 2" },
			{
				"type": "PARALLEL",
				"executors": [
					{ "task": "Task 3" },
					{ "task": "Task 4" }
				]
			},
			{ "task": "Task 5" },
			{ "task": "Task 6" },
			{
				"type": "PARALLEL_ORDERED",
				"executors": [
					{ "task": "Task 7" },
					{ "task": "Task 8" }
				]
			}
		]
	}`

	expectedExecutor := executor.NewSequential([]master.Executor[string]{
		executor.NewSingle("Task 1"),
		executor.NewSingle("Task 2"),
		executor.NewParallel([]master.Executor[string]{
			executor.NewSingle("Task 3"),
			executor.NewSingle("Task 4"),
		}),
		executor.NewSingle("Task 5"),
		executor.NewSingle("Task 6"),
		executor.NewParallelOrdered([]master.Executor[string]{
			executor.NewSingle("Task 7"),
			executor.NewSingle("Task 8"),
		}),
	})

	parsedExecutor, err := executor.FromJson[string]([]byte(executorJson))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedExecutor, parsedExecutor) {
		t.Fatalf("Parsed executor does not match expected!\nexpected: %#v\nparsed: %#v", expectedExecutor, parsedExecutor)
	}
}
