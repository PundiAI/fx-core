#!/usr/bin/env bash

set -eo pipefail

current_path=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
readonly current_path

export NODE_HOME=${NODE_HOME:-"./out/.upgrade"}
[[ -d "$NODE_HOME" ]] && rm -r "$NODE_HOME" && mkdir -p "$NODE_HOME"

[[ ! -f "${current_path}/run-upgrade.sh" ]] && echo "run-upgrade.sh not found" && exit 1
nohup "${current_path}/run-upgrade.sh" 10 >&1 &

[[ ! -f "${current_path}/run-cosmovisor.sh" ]] && echo "run-cosmovisor.sh not found" && exit 1
"${current_path}/run-cosmovisor.sh" init "v3.1.x" "v4.1.x"
