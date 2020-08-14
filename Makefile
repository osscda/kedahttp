BIN_DIR := ./bin

.PHONY: proxy
proxy:
	go build -o proxy ./cmd/proxy

.PHONY: runproxy
runproxy:
	go run ./cmd/proxy

.PHONY: dockerbuild
dockerbuild:
	docker build -t arschles/cscaler .

.PHONY: dockerbuild
dockerpush: dockerbuild
	docker push arschles/cscaler

.PHONY: dockerbuild
dockerpush: dockerbuild
	docker push arschles/cscaler

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
