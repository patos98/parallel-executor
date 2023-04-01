package worker

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type RemoteConfig struct {
	WorkerUrl string
}

func StartNewRemote[T any](ctx Context[T], cfg RemoteConfig) {
	remoteWorker := WorkerFn[T](func(task T) (result T, err error) {
		payload, err := json.Marshal(task)
		if err != nil {
			return
		}

		resp, err := http.Post(cfg.WorkerUrl, "application/json", bytes.NewReader(payload))
		if err != nil {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}

		err = json.Unmarshal(body, &result)

		return
	})

	StartNew[T](ctx, remoteWorker)
}
