#!/bin/sh

CGO=$1
ARCH=$2
GOOS=$3

GOPROXY=https://goproxy.io CGO_ENABLED=${CGO} GOOS=${GOOS} GOARCH=${ARCH} go build -a -o resource-tagger cmd/server/server.go
