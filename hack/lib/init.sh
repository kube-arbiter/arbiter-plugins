#!/usr/bin/env bash

unset CDPATH

ROOT_PATH=$(dirname "${BASH_SOURCE[0]}")/../..

declare -A CODE_PATH=(
	["observer-metric-server"]="observer-plugins/metric-server"
	["observer-prometheus-server"]="observer-plugins/prometheus"
	["executor-resource-tagger"]="executor-plugins/resource-tagger"
)

declare -A CODE_MAIN_PATH=(
	["observer-metric-server"]="main.go"
	["observer-prometheus-server"]="main.go"
	["executor-resource-tagger"]="cmd/server/server.go"
)

declare -A DOCKERFILE_PATH=(
	["observer-metric-server"]="observer-plugins/metric-server/Dockerfile"
	["observer-prometheus-server"]="observer-plugins/prometheus/Dockerfile"
	["executor-resource-tagger"]="executor-plugins/resource-tagger/Dockerfile"
)

source "${ROOT_PATH}/hack/lib/util.sh"
source "${ROOT_PATH}/hack/lib/version.sh"
source "${ROOT_PATH}/hack/lib/golang.sh"
source "${ROOT_PATH}/hack/lib/docker.sh"
