# taken from Athens 
# https://github.com/gomods/athens/blob/main/cmd/proxy/Dockerfile
ARG GOLANG_VERSION=1.14
ARG ALPINE_VERSION=3.11.5

FROM golang:${GOLANG_VERSION}-alpine AS builder

COPY . .

ARG VERSION="unset"

RUN go build -o /bin/containerscalerproxy .

FROM alpine:${ALPINE_VERSION}

ENV GO111MODULE=on

COPY --from=builder /bin/containerscalerproxy /bin/containerscalerproxy

EXPOSE 8080

ENTRYPOINT ["/bin/containerscalerproxy"]
