#!/usr/bin/env bash

set -o errexit -o nounset

NODE_HOME=$(realpath ./build/.fxcore)
export FX_RUN_HOME=$NODE_HOME
export FX_RUN_BINARY=$NODE_HOME/cosmovisor/genesis/bin/fxcored


CUR_PATH=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
"$CUR_PATH"/run-upgrade-commands.sh 10