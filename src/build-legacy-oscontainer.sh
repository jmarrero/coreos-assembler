#!/usr/bin/env bash
# shellcheck disable=SC1091
set -euo pipefail
# Start VM and call buildah
. /usr/lib/coreos-assembler/cmdlib.sh; prepare_build
. /usr/lib/coreos-assembler/cmdlib.sh
runvm -- /usr/lib/coreos-assembler/oscontainer-deprecated-legacy-format.py "$@"
