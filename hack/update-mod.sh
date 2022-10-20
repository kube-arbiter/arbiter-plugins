#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_PATH=$(dirname "${BASH_SOURCE[0]}")/..
source "${ROOT_PATH}/hack/lib/init.sh"

for path in "${CODE_PATH[@]}"; do
	echo "in $path run go mod tidy"
	cd $path && go mod tidy && cd -
done

echo "go mod updated done."
