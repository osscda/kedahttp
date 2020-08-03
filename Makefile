.PHONY: proxy
proxy:
	go build -o proxy ./cmd/proxy

.PHONY: runproxy
runproxy:
	go run ./cmd/proxy

.PHONY: controller
controller:
	go build -o controller ./cmd/controller

.PHONY: runcontroller
runcontroller:
	go run ./cmd/controller


.PHONY: cli
cli:
	go build -o cli ./cmd/cli