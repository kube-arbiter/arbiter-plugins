#!/bin/bash

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/..
rm -rf $ROOT/observer-plugins/metric-server/arbiter-metric-server;
rm -rf $ROOT/observer-plugins/prometheus/arbiter-prometheus-server;
rm -rf $ROOT/resource-tagger-plugin/resource-tagger;
rm -rf $ROOT/resource-tagger-plugin/resource-tagger-client;
docker rmi arbiter/arbiter-metric-server:$1 arbiter/arbiter-prometheus-server:$1 arbiter/resource-tagger-plugin:$1 2> /dev/null || true;

