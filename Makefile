.PHONY: build

build:
	@echo
	@echo "⋮⋮ Building..."
	go build -ldflags "-X main.version=`cat build_number``date -u +.%Y%m%d%H%M%S`"

test: build
	@echo
	@echo "⋮⋮ Testing..."
	go test -count=1 ./... -coverprofile cover.out

review: test
	@echo
	@echo "⋮⋮ Reviewing tests..."
	go tool cover -html cover.out

container: test
	@echo
	@echo "⋮⋮ Creating container..."
	docker build -f ./build/package/Dockerfile -t sheila .

local: container
	@echo
	@echo "⋮⋮ Launching containers for local use..."
	./scripts/config-remote-prometheus.sh
	@sleep 5
	docker compose -f ./deployments/docker-compose.yaml --project-name sheila up -d --force-recreate
	@sleep 5
	./scripts/config-revert-prometheus.sh
	docker ps | grep sheila

clean-local:
	docker compose -f ./deployments/docker-compose.yaml --project-name sheila down

all: build test container
	@echo

