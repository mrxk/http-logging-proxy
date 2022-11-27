.PHONY: build
build:
	go build

.PHONY: docker
docker:
	docker build -t http-logging-proxy .
