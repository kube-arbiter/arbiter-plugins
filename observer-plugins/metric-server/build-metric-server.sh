#!/bin/sh

CGO=$1
ARCH=$2
GOOS=$3

echo "args: cgo=${CGO}, arch=${ARCH} goos=${GOOS}"

sleep 3
#GOPROXY=https://goproxy.io CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o arbiter-metric-server .
GOPROXY=https://goproxy.io CGO_ENABLED=$1 GOOS=$3 GOARCH=$2 go build -o observer-metric-server .
