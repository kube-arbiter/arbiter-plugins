#!/bin/bash

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/..
rm -rf $ROOT/observer-plugins/metric-server/observer-metric-server;
rm -rf $ROOT/observer-plugins/prometheus/observer-prometheus-server;
rm -rf $ROOT/executor-plugins/resource-tagger/resource-tagger;
docker rmi kubearbiter/observer-metric-server:$1 kubearbiter/observer-prometheus-server:$1 kubearbiter/executor-resource-tagger:$1 2> /dev/null || true;
