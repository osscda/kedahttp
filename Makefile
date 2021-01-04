BIN_DIR := ./bin

.PHONY: proxy
proxy:
	go build -o ${BIN_DIR}/proxy ./cmd/proxy

.PHONY: runproxy
runproxy:
	go run ./cmd/proxy


.PHONY: proto
proto:
	protoc kedascaler.external.proto --go_out=plugins=grpc:externalscaler

.PHONY: dockerbuild
dockerbuild:
	docker build -t arschles/cscaler .

.PHONY: dockerpush
dockerpush: dockerbuild
	docker push arschles/cscaler

.PHONY: cli
cli:
	cd cli && cargo build
	cp cli/target/debug/cli cscaler

.PHONY: clean-cli
clean-cli:
	rm -rf ${BIN_DIR}/cscaler

.PHONY: clean-bin
clean-bin: clean-cli

.PHONY: clean
clean: clean-bin
