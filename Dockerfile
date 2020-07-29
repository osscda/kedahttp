# taken from Athens 
# https://github.com/gomods/athens/blob/main/cmd/proxy/Dockerfile
ARG GOLANG_VERSION=1.14
ARG ALPINE_VERSION=3.11.5

FROM golang:${GOLANG_VERSION}-alpine AS builder

WORKDIR $GOPATH/src/github.com/arschles/containerscaler

COPY . .

ARG VERSION="unset"

RUN GO111MODULE=on CGO_ENABLED=0 GOPROXY="https://proxy.golang.org" go build -o /bin/containerscalerproxy .

FROM alpine:${ALPINE_VERSION}

ENV GO111MODULE=on

COPY --from=builder /bin/containerscalerproxy /bin/containerscalerproxy

EXPOSE 8080

ENTRYPOINT ["/bin/containerscalerproxy"]
