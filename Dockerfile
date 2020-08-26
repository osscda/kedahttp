# taken from Athens 
# https://github.com/gomods/athens/blob/main/cmd/proxy/Dockerfile
ARG GOLANG_VERSION=1.14
ARG ALPINE_VERSION=3.11.5

FROM golang:${GOLANG_VERSION}-alpine AS builder

WORKDIR $GOPATH/src/github.com/arschles/containerscaler

COPY . .
COPY kedascaler.external.proto /bin/proto.proto

ARG VERSION="unset"

RUN GO111MODULE=on CGO_ENABLED=0 GOPROXY="https://proxy.golang.org" go build -o /bin/containerscalerproxy ./cmd/proxy

FROM alpine:${ALPINE_VERSION}

RUN apk add -U curl

ENV GO111MODULE=on

COPY --from=builder /bin/containerscalerproxy /bin/containerscalerproxy
COPY --from=builder /bin/proto.proto /bin/proto.proto

RUN curl -o grpcurl.tgz -L https://github.com/fullstorydev/grpcurl/releases/download/v1.7.0/grpcurl_1.7.0_linux_x86_64.tar.gz && \
    tar -xzf grpcurl.tgz && \
    mv grpcurl /bin/grpcurl

EXPOSE 8080

ENTRYPOINT ["/bin/containerscalerproxy"]
