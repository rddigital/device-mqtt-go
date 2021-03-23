.PHONY: build test clean prepare update docker docker-push docker-arm64

GO = CGO_ENABLED=0 GO111MODULE=on go

MICROSERVICES=cmd/device-mqtt

.PHONY: $(MICROSERVICES)

DOCKERS=docker_device_mqtt_go

.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
GIT_SHA=$(shell git rev-parse HEAD)

GOFLAGS=-ldflags "-X github.com/rddigital/device-mqtt-go.Version=$(VERSION)"

build: $(MICROSERVICES)
	$(GO) build ./...

cmd/device-mqtt:
	$(GO) build $(GOFLAGS) -o $@ ./cmd

test:
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) vet ./...
	gofmt -l .
	["`gofmt -l .`" = ""]
	./bin/test-attribution-txt.sh
	./bin/test-go-mod-tidy.sh

clean:
	rm -f $(MICROSERVICES)

run:
	cd cmd && ./device-mqtt

docker: $(DOCKERS)

docker_device_mqtt_go:
	docker build \
		--label "git_sha=$(GIT_SHA)" \
		-t rddigital/docker-device-mqtt-go:$(VERSION) \
		.
docker-push:
	docker push \
	rddigital/docker-device-mqtt-go:$(VERSION)

docker-arm64:
	docker buildx build  --platform linux/arm64,linux/arm/v7 -t rddigital/docker-device-mqtt-go-arm64:$(VERSION) --push .