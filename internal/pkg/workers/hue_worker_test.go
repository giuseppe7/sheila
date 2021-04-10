package workers

import (
	"os"
	"testing"
)

const testApplicationNamespace = "testShiela"
const envKeyUsername = "SHEILA_USER"

func TestDoWork(t *testing.T) {
	// Initialize the client used by all tests in this package.
	authorizedUser, ok := os.LookupEnv(envKeyUsername)
	if !ok {
		panic("environment variable not set")
	}

	hueWorker, err := NewHueWorker(testApplicationNamespace, authorizedUser)
	if err != nil {
		t.Errorf("not expecting error on worker initialization, %s", err)
		return
	}

	err = hueWorker.ObserveLights()
	if err != nil {
		t.Errorf("not expecting error on observing light , %s", err)
		return
	}

	err = hueWorker.ObserveSensors()
	if err != nil {
		t.Errorf("not expecting error on observing sensors , %s", err)
		return
	}
}
