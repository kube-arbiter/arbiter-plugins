#!/usr/bin/env bash

unset CDPATH

ROOT_PATH=$(dirname "${BASH_SOURCE[0]}")/../..

declare -A CODE_PATH=(
	["observer-default-plugins"]="observer-plugins/default-plugins"
	["executor-default-plugins"]="executor-plugins/default-plugins"
)

declare -A CODE_MAIN_PATH=(
	["observer-default-plugins"]="cmd/server/server.go"
	["executor-default-plugins"]="cmd/server/server.go"
)

declare -A DOCKERFILE_PATH=(
	["observer-default-plugins"]="observer-plugins/default-plugins/Dockerfile"
	["executor-default-plugins"]="executor-plugins/default-plugins/Dockerfile"
)

source "${ROOT_PATH}/hack/lib/util.sh"
source "${ROOT_PATH}/hack/lib/version.sh"
source "${ROOT_PATH}/hack/lib/golang.sh"
source "${ROOT_PATH}/hack/lib/docker.sh"
