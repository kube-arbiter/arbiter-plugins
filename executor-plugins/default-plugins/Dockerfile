# Build resource-tagger plugin image

FROM golang:1.18 as builder

ARG ARCH=amd64
ARG CGO=0
ARG GOOS=linux

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
COPY build-resource-tagger.sh build-resource-tagger.sh
COPY staging/ staging/

## Copy the go source and vendor
COPY cmd/server/server.go cmd/server/server.go
COPY pkg/ pkg/

RUN ./build-resource-tagger.sh ${CGO} ${ARCH} ${GOOS}

## Build binary
#RUN GOPROXY=https://goproxy.io CGO_ENABLED=${CGO} GOOS=${GOOS} GOARCH=${ARCH} go build -a -o resource-tagger cmd/server/server.go
#RUN GOPROXY=https://goproxy.io CGO_ENABLED=${CGO} GOOS=${GOOS} GOARCH=${ARCH} go build -a -o resource-tagger-client cmd/client/client.go

## Build image
FROM alpine:3.15
WORKDIR /
COPY --from=builder /workspace/resource-tagger /usr/local/bin/resource-tagger

CMD ["resource-tagger"]