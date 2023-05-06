#!/usr/bin/env bash

readonly project_dir="$(git rev-parse --show-toplevel)"

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

function close() {
  "${project_dir}/tests/scripts/fxcore.sh" close
  "${project_dir}/tests/scripts/pundix.sh" close
}

trap close EXIT

run
