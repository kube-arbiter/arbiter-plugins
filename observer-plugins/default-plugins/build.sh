#!/bin/sh

CGO=$1
ARCH=$2
GOOS=$3

echo "args: cgo=${CGO}, arch=${ARCH} goos=${GOOS}"

GOPROXY=https://goproxy.io CGO_ENABLED=$1 GOOS=$3 GOARCH=$2 go build -o observer-default-plugins cmd/server/server.go
