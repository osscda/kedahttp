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

.PHONY: helminstall
helminstall:
	helm install cscaler ./charts/cscaler-proxy

.PHONY: helmupgrade
helmupgrade:
	helm upgrade cscaler ./charts/cscaler-proxy


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
