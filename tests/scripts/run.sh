#!/usr/bin/env bash

set -eo pipefail

PROJECT_DIR="${PROJECT_DIR:-"$(git rev-parse --show-toplevel)"}"
export PROJECT_DIR
export OUT_DIR="${PROJECT_DIR}/out"

readonly script_dir="${PROJECT_DIR}/tests/scripts"

function run() {
  "${script_dir}/fxcore.sh" init
  "${script_dir}/fxcore.sh" start

  export IBC_CHANNEL="channel-0"
  export PURSE_ADDRESS="0x0000000000000000000000000000000000000000"
  export PUNDIX_ADDRESS="0x0000000000000000000000000000000000000000"
  "${script_dir}/pundix.sh" init
  "${script_dir}/pundix.sh" start

  "${script_dir}/ibcrelayer.sh" transfer
  "${script_dir}/ibcrelayer.sh" init
  "${script_dir}/ibcrelayer.sh" start
}

function close() {
  "${script_dir}/fxcore.sh" stop
  "${script_dir}/pundix.sh" stop
  "${script_dir}/ibcrelayer.sh" stop
  rm -rf "${PROJECT_DIR}/out"
}

#trap close EXIT

run
