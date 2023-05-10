#!/usr/bin/env bash

set -eo pipefail

project_dir="$(git rev-parse --show-toplevel)"
readonly project_dir

function start() {
  "${project_dir}/tests/scripts/fxcore.sh" start

  (
    export IBC_CHAINNEL="channel-0"
    export PURSE_ADDRESS="0x0000000000000000000000000000000000000000"
    export PUNDIX_ADDRESS="0x0000000000000000000000000000000000000000"
    "${project_dir}/tests/scripts/pundix.sh" start
  )

  "${project_dir}/tests/scripts/run-ibc.sh"
}

function stop() {
  "${project_dir}/tests/scripts/fxcore.sh" stop
  "${project_dir}/tests/scripts/pundix.sh" stop
}

trap stop EXIT

run
