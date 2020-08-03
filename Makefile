BIN_DIR := ./bin

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
	go build -v -o ${BIN_DIR}/cscaler ./cmd/cli

.PHONY: clean-cli
clean-cli: 
	rm -rf ${BIN_DIR}/cscaler

.PHONY: clean-bin
clean-bin: clean-cli

.PHONY: clean
clean: clean-bin