package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type HueClient struct {
	httpClient     *http.Client
	histogram      *prometheus.HistogramVec
	FoundServer    HueDiscoveryResult
	AuthorizedUser string
}

type HueDiscoveryResult struct {
	ID                string `json:"id"`
	InternalIpAddress string `json:"internalipaddress"`
}

// HueLight structure with a subset of fields from the response.
type HueLight struct {
	Name  string `json:"name"`
	State struct {
		On        bool `json:"on"`
		Reachable bool `json:"reachable"`
	}
	UniqueID string `json:"uniqueid"`
}

// HueSensor structure with a subset of fields from the response.
type HueSensor struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	State struct {
		LightLevel  int    `json:"lightlevel"`
		Dark        bool   `json:"dark"`
		Daylight    bool   `json:"daylight"`
		Temperature int    `json:"temperature"`
		Presence    bool   `json:"presence"`
		LastUpdated string `json:"lastupdated"`
	}
	Config struct {
		On        bool `json:"on"`
		Battery   int  `json:"battery"`
		Reachable bool `json:"reachable"`
	}
	UniqueID string `json:"uniqueid"`
}

func NewHueClient(applicationNamespace string, authorizedUser string) *HueClient {
	hueClient := new(HueClient)

	hueClient.AuthorizedUser = authorizedUser

	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
	}
	hueClient.httpClient = &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	// Capture metrics on the command execution times.
	hueClient.histogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: applicationNamespace,
			Name:      "hue_client_command_duration_seconds",
			Help:      "Histogram of client calls in seconds.",
			Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5, 10, 20, 30},
		},
		[]string{"url", "method", "statusCode"},
	)
	prometheus.MustRegister(hueClient.histogram)

	return hueClient
}

// DiscoverHue returns the first result found or an error with an empty structure if none.
func (hueClient *HueClient) DiscoverHue() (HueDiscoveryResult, error) {
	// Constant referenced by the developer API.
	const discoveryUrl = "http://discovery.meethue.com"
	var results []HueDiscoveryResult
	var err error

	req, _ := http.NewRequest("GET", discoveryUrl, nil)
	start := time.Now()
	resp, err := hueClient.httpClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		hueClient.histogram.WithLabelValues(discoveryUrl, req.Method, strconv.Itoa(-1)).Observe(duration.Seconds())
		return results[0], err
	} else if resp.StatusCode != 200 {
		hueClient.histogram.WithLabelValues(discoveryUrl, req.Method, strconv.Itoa(resp.StatusCode)).Observe(duration.Seconds())
		return results[0], fmt.Errorf("connected but received non-200 status code %d", resp.StatusCode)
	}
	hueClient.histogram.WithLabelValues(discoveryUrl, req.Method, strconv.Itoa(resp.StatusCode)).Observe(duration.Seconds())
	b, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(b, &results)
	if err != nil {
		log.Println("error in umarshaling", err)
	}
	defer resp.Body.Close()
	//log.Print(fmt.Sprintf("obtained: %v", results))

	hueClient.FoundServer = results[0]
	return results[0], nil
}

func (hueClient *HueClient) IsReachable() bool {
	if hueClient.FoundServer.ID == "" {
		log.Println("server was not discovered")
		return false
	}

	apiEndpoint := hueClient.FoundServer.InternalIpAddress
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s", apiEndpoint), nil)
	start := time.Now()
	resp, err := hueClient.httpClient.Do(req)
	duration := time.Since(start)
	//log.Print(fmt.Sprintf("duration: %d", duration))
	hueClient.histogram.WithLabelValues(apiEndpoint, req.Method, strconv.Itoa(-1)).Observe(duration.Seconds())

	if err != nil {
		return false
	} else if resp.StatusCode != 200 {
		return false
	} else {
		return true
	}
}

// Returns the map of HueLight objects obtained from the Hue System.
func (hueClient *HueClient) GetLights() (map[int]HueLight, error) {
	lights := make(map[int]HueLight)
	apiEndpoint := "/lights"
	url := fmt.Sprintf("http://%s/api/%s%s", hueClient.FoundServer.InternalIpAddress, hueClient.AuthorizedUser, apiEndpoint)

	req, _ := http.NewRequest("GET", url, nil)
	start := time.Now()
	resp, err := hueClient.httpClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		hueClient.histogram.WithLabelValues(apiEndpoint, req.Method, strconv.Itoa(-1)).Observe(duration.Seconds())
		return lights, err
	} else if resp.StatusCode != 200 {
		hueClient.histogram.WithLabelValues(apiEndpoint, req.Method, strconv.Itoa(resp.StatusCode)).Observe(duration.Seconds())
		return lights, fmt.Errorf("connected but received non-200 status code %d", resp.StatusCode)
	}
	hueClient.histogram.WithLabelValues(apiEndpoint, req.Method, strconv.Itoa(resp.StatusCode)).Observe(duration.Seconds())
	b, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(b, &lights)
	if err != nil {
		log.Println("error in umarshaling", err)
	}
	defer resp.Body.Close()
	//log.Print(fmt.Sprintf("obtained: %v", lights))

	return lights, err
}

// Returns the map of HueLight objects obtained from the Hue System.
func (hueClient *HueClient) GetSensors() (map[int]HueSensor, error) {
	sensors := make(map[int]HueSensor)
	apiEndpoint := "/sensors"
	url := fmt.Sprintf("http://%s/api/%s%s", hueClient.FoundServer.InternalIpAddress, hueClient.AuthorizedUser, apiEndpoint)

	req, _ := http.NewRequest("GET", url, nil)
	start := time.Now()
	resp, err := hueClient.httpClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		hueClient.histogram.WithLabelValues(apiEndpoint, req.Method, strconv.Itoa(-1)).Observe(duration.Seconds())
		return sensors, err
	} else if resp.StatusCode != 200 {
		hueClient.histogram.WithLabelValues(apiEndpoint, req.Method, strconv.Itoa(resp.StatusCode)).Observe(duration.Seconds())
		return sensors, fmt.Errorf("connected but received non-200 status code %d", resp.StatusCode)
	}
	hueClient.histogram.WithLabelValues(apiEndpoint, req.Method, strconv.Itoa(resp.StatusCode)).Observe(duration.Seconds())
	b, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(b, &sensors)
	if err != nil {
		log.Println("error in umarshaling", err)
	}
	defer resp.Body.Close()
	//log.Print(fmt.Sprintf("obtained: %v", sensors))

	return sensors, err
}
