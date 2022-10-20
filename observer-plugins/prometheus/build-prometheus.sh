#!/bin/sh

GOPROXY=https://goproxy.io CGO_ENABLED=$1 GOOS=$3 GOARCH=$2 go build -o observer-prometheus-server .
