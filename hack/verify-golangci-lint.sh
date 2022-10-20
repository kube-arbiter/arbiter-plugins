#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_PATH=$(dirname "${BASH_SOURCE[0]}")/..
source "${ROOT_PATH}/hack/lib/init.sh"

go::setup_env

echo "installing golangci-lint"
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0
cd "${ROOT_PATH}"

echo "running golangci-lint"

for path in "${CODE_PATH[@]}"; do
	if (cd $path && golangci-lint run ./... && cd -); then
		echo "in $path golangci-lint verified."
	else
		echo "in $path golangci-lint failed."
		exit 1
	fi
done
