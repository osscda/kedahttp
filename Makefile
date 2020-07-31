BIN_DIR := ./bin

.PHONY: proxy
proxy:
	go build -o proxy ./cmd/proxy

.PHONY: controller
controller:
	go build -o controller ./cmd/controller

.PHONY: cli
cli:
	go build -v -o ${BIN_DIR}/cscaler ./cmd/cli

.PHONY: clean-cli
clean-cli: 
	rm -rf ${BIN_DIR}/cli

.PHONY: clean-bin
clean-bin: clean-cli

.PHONY: clean
clean: clean-bin