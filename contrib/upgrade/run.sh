#!/usr/bin/env bash

set -o errexit -o nounset

export CUR_VERSION=${CUR_VERSION:-"v3"}
export NEXT_VERSION=${NEXT_VERSION:-"v4"}

export NODE_HOME=${NODE_HOME:-"./out/.fxcore"}
export BINARY=${BINARY:-"$NODE_HOME/cosmovisor/genesis/bin/fxcored"}

export CHAIN_ID=${CHAIN_ID:-"fxcore"}

cur_path=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
readonly cur_path

"${cur_path}/run.sh" > "./run.log" 2>&1 &
readonly pid=$!
trap 'kill -9 $pid' SIGINT SIGTERM EXIT

"${cur_path}/run-upgrade.sh" 10

wait $pid