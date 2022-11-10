# Build executor-default-plugins image

FROM golang:1.18 as builder

ARG ARCH=amd64
ARG CGO=0
ARG GOOS=linux

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
COPY build-executor-default-plugins.sh build-executor-default-plugins.sh
COPY staging/ staging/

## Copy the go source and vendor
COPY cmd/server/server.go cmd/server/server.go
COPY pkg/ pkg/

RUN ./build-executor-default-plugins.sh ${CGO} ${ARCH} ${GOOS}

## Build binary
#RUN GOPROXY=https://goproxy.io CGO_ENABLED=${CGO} GOOS=${GOOS} GOARCH=${ARCH} go build -a -o executor-default-plugins cmd/server/server.go

## Build image
FROM alpine:3.15
WORKDIR /
COPY --from=builder /workspace/executor-default-plugins /usr/local/bin/executor-default-plugins

CMD ["executor-default-plugins"]
