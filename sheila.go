package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/giuseppe7/sheila/internal/pkg/workers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const applicationNamespace = "sheila"
const envKeyUsername = "SHEILA_USER"

// Variable to be set by the Go linker at build time.
var version string

// Set up observability with Prometheus handler for metrics.
func initObservability() {

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	// Register a version gauge.
	versionGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: applicationNamespace,
			Name:      applicationNamespace + "_info",
			Help:      "Version of the application.",
		},
	)
	prometheus.MustRegister(versionGauge)
	versionValue, err := strconv.ParseFloat(version, 64)
	if err != nil {
		versionValue = 0.0
	}
	versionGauge.Set(versionValue)
}

func main() {
	log.Println("Coming online...")
	log.Print(fmt.Sprintf("Version: %v\n", version))

	// Channel to be aware of an OS interrupt like Control-C.
	var waiter sync.WaitGroup
	waiter.Add(1)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Set up observability.
	initObservability()

	// Obtain enviroment variables if required.
	authorizedUser, ok := os.LookupEnv(envKeyUsername)
	if !ok {
		log.Fatalf("environment variable %s not set", envKeyUsername)
	}

	// Do the work.
	hueWorker, err := workers.NewHueWorker(applicationNamespace, authorizedUser)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		for {
			err = hueWorker.ObserveLights()
			if err != nil {
				log.Printf(err.Error())
			}
			err = hueWorker.ObserveSensors()
			if err != nil {
				log.Printf(err.Error())
			}
			time.Sleep(time.Second)
		}
	}()

	// Function and waiter to wait for the OS interrupt and do any clean-up.
	go func() {
		<-c
		fmt.Println("\r")
		log.Println("Interrupt captured.")
		waiter.Done()
	}()
	waiter.Wait()

	// Shut down the application.
	log.Println("Shutting down.")
	return
}
