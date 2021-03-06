# sheila

S.H.E.I.L.A. stands for "Sampling Hue Environment Information Leveraging APIs"
and is meant solely as a personal project. Currently, it does the following:

* Samples Hue light and sensor information and exposes this as an endpoint.
* Lauches side car containers Prometheus and Grafana to visualize this
  information.




## Build and Test

Follow these steps to build the Docker image for this application.

1. Change the working directory to the root of this repo.
0. Create the executable:  
   `go build` 
0. Run the tests:  
   `go test -count=1 ./... -coverprofile cover.out`  
   `go tool cover -func cover.out`  
   `go tool cover -html cover.out`


## Run Locally

Simply run these steps to run the application locally as the stand-alone application.

1. Run the application itself, set the environment variables and start the application.  
   `export SHEILA_USER=[authorized local Hue user]`  
   `./sheila`  
0. The application is instrumented with the Prometheus client go-lang library and exposes its metrics endpoint.  
   `curl http://0.0.0.0:2112/metrics`

Follow these steps to run the application locally alongside Prometheus and Grafana containers.

1. Create the container:  
   `docker build -f ./build/package/Dockerfile -t sheila .`
0. Launch the containers for the application, Prometheus, and Grafana:  
   `docker-compose -f ./deployments/docker-compose.yaml up -d`
0. Look up the dynamic ports in use:
   `docker-compose -f ./deployments/docker-compose.yaml ps`
0. Open Grafana in your browser to view the metrics dashboard:  
   `http://0.0.0.0:${GRAFANA_PORT}`
0. Once done exploring, tear down your environment:  
   `docker-compose -f ./deployments/docker-compose.yaml down`




#### References

1. https://github.com/golang-standards/project-layout
1. Hue API at https://developers.meethue.com/