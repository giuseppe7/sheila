package clients

import (
	"os"
	"testing"
)

const testApplicationNamespace = "testShiela"
const envKeyUsername = "SHEILA_USER"

// Object to be used by the tests so initialize in the TestMain below.
var hueClient *HueClient

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags

	// Initialize the client used by all tests in this package.
	authorizedUser, ok := os.LookupEnv(envKeyUsername)
	if !ok {
		panic("environment variable not set")
	}

	hueClient = NewHueClient(testApplicationNamespace, authorizedUser)
	discoveryResult, err := hueClient.DiscoverHue()

	// Iterate to check connectivity before failing outright.
	if err != nil {
		panic("expected to discover a registered hue system: '%s'")
	} else if discoveryResult.ID == "" {
		panic("expected at least one hue system found")
	}

	os.Exit(m.Run())
}

func TestIsReachable(t *testing.T) {
	if !hueClient.IsReachable() {
		t.Errorf("expected to be able to reach the hue system at %s", hueClient.FoundServer.InternalIpAddress)
	}
}

func TestGetLights(t *testing.T) {
	if hueClient.IsReachable() {
		_, err := hueClient.GetLights()
		if err != nil {
			t.Errorf("failed to get lights, %s", err)
		}
	}
}

func TestGetSensors(t *testing.T) {
	if hueClient.IsReachable() {
		_, err := hueClient.GetSensors()
		if err != nil {
			t.Errorf("failed to get sensors, %s", err)
		}
	}
}
