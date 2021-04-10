package workers

import (
	"fmt"

	"github.com/giuseppe7/sheila/internal/pkg/clients"
	"github.com/prometheus/client_golang/prometheus"
)

type HueWorker struct {
	hueClient    *clients.HueClient
	gaugeLights  *prometheus.GaugeVec
	gaugeSensors *prometheus.GaugeVec
}

func NewHueWorker(applicationNamespace string, authorizedUser string) (*HueWorker, error) {
	worker := new(HueWorker)

	worker.hueClient = clients.NewHueClient(applicationNamespace, authorizedUser)
	_, err := worker.hueClient.DiscoverHue()
	if err != nil {
		return nil, err
	}
	if !worker.hueClient.IsReachable() {
		return nil, fmt.Errorf("hue system is not reachable")
	}

	labels := []string{"uniqueid", "name", "state"}
	worker.gaugeLights = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: applicationNamespace,
			Name:      "hue_worker_light_gauge",
			Help:      "Gauge of lights known after fetched from various calls.",
		},
		labels,
	)
	prometheus.MustRegister(worker.gaugeLights)

	labels = []string{"uniqueid", "name", "type"}
	worker.gaugeSensors = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: applicationNamespace,
			Name:      "hue_worker_sensor_gauge",
			Help:      "Gauge of sensors known after fetched from various calls.",
		},
		labels,
	)
	prometheus.MustRegister(worker.gaugeSensors)

	return worker, nil
}

func (worker *HueWorker) ObserveLights() error {
	hueLights, err := worker.hueClient.GetLights()
	if err != nil {
		return err
	}

	for _, light := range hueLights {
		isOn := 0.0
		isReachable := 0.0
		if light.State.On {
			isOn = 1
		}
		if light.State.Reachable {
			isReachable = 1
		}

		worker.gaugeLights.WithLabelValues(light.UniqueID, light.Name, "on").Set(isOn)
		worker.gaugeLights.WithLabelValues(light.UniqueID, light.Name, "reachable").Set(isReachable)
	}
	return nil
}

func (worker *HueWorker) ObserveSensors() error {
	hueSensors, err := worker.hueClient.GetSensors()
	if err != nil {
		return err
	}

	for _, sensor := range hueSensors {
		if sensor.Type == "ZLLTemperature" {
			worker.gaugeSensors.WithLabelValues(sensor.UniqueID, sensor.Name, sensor.Type).Set(float64(sensor.State.Temperature))

		} else if sensor.Type == "ZLLPresence" {
			presence := 0.0
			if sensor.State.Presence {
				presence = 1.0
			}
			reachable := 0.0
			if sensor.Config.Reachable {
				reachable = 1.0
			}
			worker.gaugeSensors.WithLabelValues(sensor.UniqueID, sensor.Name, sensor.Type).Set(presence)
			worker.gaugeSensors.WithLabelValues(sensor.UniqueID, sensor.Name, "battery").Set(float64(sensor.Config.Battery))
			worker.gaugeSensors.WithLabelValues(sensor.UniqueID, sensor.Name, "reachable").Set(reachable)

		} else if sensor.Type == "ZLLLightLevel" {
			worker.gaugeSensors.WithLabelValues(sensor.UniqueID, sensor.Name, sensor.Type).Set(float64(sensor.State.LightLevel))

		} else {
			worker.gaugeSensors.WithLabelValues(sensor.UniqueID, sensor.Name, sensor.Type).Set(0)
		}
	}

	return nil
}

/*
// Returns the DownDetectorResults object, including the raw JSON, if no parsing errors.
func (worker *DownDetectorWorker) getStatus(phrase string, statusChannel chan clients.DownDetectorResults) {
	results, err := worker.downDetectorClient.GetStatus(phrase)
	if err != nil {
		// Error or timeout in fetching the status.
		log.Printf("error in getting status for %v, %v", phrase, err)
		worker.gaugeWorker.WithLabelValues(phrase).Set(-1.0)
	} else {
		// Grab the last entry and track it.
		lastReport := results.Series.Reports.Data[len(results.Series.Reports.Data)-1].Y
		if worker.debug {
			log.Printf(fmt.Sprintf("%s has %d", phrase, lastReport))
		}
		worker.gaugeWorker.WithLabelValues(phrase).Set(float64(lastReport))

		// TODO: Do something with the baseline?
		// lastBaseline := results.Series.Reports.Baseline[len(results.Series.Reports.Baseline)-1].Y
	}
	statusChannel <- results
}

// DoWork is the actual main method for this worker.
func (worker *DownDetectorWorker) DoWork() {
	statusChannel := make(chan clients.DownDetectorResults, len(worker.phrases))
	go func() {
		for {
			time.Sleep(1 * time.Second)
			worker.gaugeChannel.WithLabelValues("status").Set(float64(len(statusChannel)))
		}
	}()

	for {
		for _, phrase := range worker.phrases {
			go worker.getStatus(phrase, statusChannel)
		}
		for i := 0; i < len(worker.phrases); i++ {
			<-statusChannel
		}
		time.Sleep(1 * time.Minute)
	}
}

*/
