package tests

import (
	"reflect"
	"testing"

	"github.com/patos98/parallel-executor/executor"
	"github.com/patos98/parallel-executor/master"
)

type comparableString string

func (s comparableString) Equals(other comparableString) bool {
	return s == other
}

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

	expectedExecutor := executor.NewSequential([]master.Executor[comparableString]{
		executor.NewSingle(comparableString("Task 1")),
		executor.NewSingle(comparableString("Task 2")),
		executor.NewParallel([]master.Executor[comparableString]{
			executor.NewSingle(comparableString("Task 3")),
			executor.NewSingle(comparableString("Task 4")),
		}),
		executor.NewSingle(comparableString("Task 5")),
		executor.NewSingle(comparableString("Task 6")),
		executor.NewParallelOrdered([]master.Executor[comparableString]{
			executor.NewSingle(comparableString("Task 7")),
			executor.NewSingle(comparableString("Task 8")),
		}),
	})

	parsedExecutor, err := executor.FromJson[comparableString]([]byte(executorJson))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedExecutor, parsedExecutor) {
		t.Fatalf("Parsed executor does not match expected!\nexpected: %#v\nparsed: %#v", expectedExecutor, parsedExecutor)
	}
}
