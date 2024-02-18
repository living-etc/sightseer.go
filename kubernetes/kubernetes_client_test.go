package kubernetes

import (
	"errors"
	"testing"
)

const kubectloutput = `worker-1   Ready   <none>   15h   v1.21.0
worker-2   Ready   <none>   15h   v1.21.0
worker-3   Ready   <none>   15h   v1.21.0
`

var kubectlErrorOutput = errors.New(
	`E0216 18:28:19.876478   93619 memcache.go:265] couldn't get current server API group list: Get "https://127.0.0.1:6443/api?timeout=32s": dial tcp 127.0.0.1:6443: connect: connection refused
E0216 18:28:19.876749   93619 memcache.go:265] couldn't get current server API group list: Get "https://127.0.0.1:6443/api?timeout=32s": dial tcp 127.0.0.1:6443: connect: connection refused
E0216 18:28:19.877870   93619 memcache.go:265] couldn't get current server API group list: Get "https://127.0.0.1:6443/api?timeout=32s": dial tcp 127.0.0.1:6443: connect: connection refused
E0216 18:28:19.878989   93619 memcache.go:265] couldn't get current server API group list: Get "https://127.0.0.1:6443/api?timeout=32s": dial tcp 127.0.0.1:6443: connect: connection refused
E0216 18:28:19.880078   93619 memcache.go:265] couldn't get current server API group list: Get "https://127.0.0.1:6443/api?timeout=32s": dial tcp 127.0.0.1:6443: connect: connection refused
The connection to the server 127.0.0.1:6443 was refused - did you specify the right host or port?
`,
)

type MockCommandExecutor struct {
	MockResponse string
	MockError    error
}

func (executor MockCommandExecutor) executeCommand(
	command string,
) (string, error) {
	return executor.MockResponse, executor.MockError
}

func TestWorkers(t *testing.T) {
	mockKubectlCommandExecutor := MockCommandExecutor{
		MockResponse: kubectloutput,
		MockError:    nil,
	}

	kubernetesClient := KubernetesClient{
		kubeConfigPath:         "",
		caCertPath:             "",
		apiServierHost:         "",
		kubectlCommandExecutor: mockKubectlCommandExecutor,
	}

	workers, _ := kubernetesClient.Workers()

	t.Run("number of workers", func(t *testing.T) {
		got := len(workers)
		want := 3
		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("state of workers", func(t *testing.T) {
		notReadyWorkers := []Worker{}
		for _, worker := range workers {
			if worker.Status != "Ready" {
				notReadyWorkers = append(notReadyWorkers, worker)
			}
		}

		noOfNotReadyWorkers := len(notReadyWorkers)
		if noOfNotReadyWorkers > 0 {
			t.Errorf(
				"%v workers not in READY state: %v",
				noOfNotReadyWorkers,
				notReadyWorkers,
			)
		}
	})
}

func TestWorkers_AuthenticationError(t *testing.T) {
	mockKubectlCommandExecutor := MockCommandExecutor{
		MockResponse: "",
		MockError:    kubectlErrorOutput,
	}

	kubernetesClient := KubernetesClient{
		kubeConfigPath:         "",
		caCertPath:             "",
		apiServierHost:         "",
		kubectlCommandExecutor: mockKubectlCommandExecutor,
	}

	t.Run("authentication error", func(t *testing.T) {
		_, err := kubernetesClient.Workers()
		if err == nil {
			t.Error("wanted and error, got nil")
		}
	})
}
