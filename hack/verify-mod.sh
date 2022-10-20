#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT_PATH=$(dirname "${BASH_SOURCE[0]}")/..
source "${ROOT_PATH}/hack/lib/init.sh"

hack/update-mod.sh

if ! _out="$(
	git --no-pager diff -I"edited\smanually" --exit-code \
		$(find . -mindepth 2 -not -path "*/_output/*" -type f \( -name go.mod -o -name go.sum \))
)"; then
	echo "Generated output differs" >&2
	echo "${_out}" >&2
	echo "Verification for go mod failed."
	exit 1
fi

echo "go mod verified."
