#!/bin/bash


CGO=$1
ARCH=$2
GOOS=$3

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/..

echo "build metric server"

cd $ROOT/observer-plugins/metric-server;
./build-metric-server.sh $1 $2 $3

echo "build prometheus server"
cd $ROOT/observer-plugins/prometheus;
./build-prometheus.sh $1 $2 $3

echo "build resource tagger"
cd $ROOT/executor-plugins/resource-tagger
./build-resource-tagger.sh $1 $2 $3
