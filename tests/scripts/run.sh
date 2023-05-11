#!/usr/bin/env bash

set -eo pipefail

function start() {
  "${SCRIPT_DIR}/fxcore.sh" start

  export IBC_CHANNEL="channel-0"
  export PURSE_ADDRESS="0x0000000000000000000000000000000000000000"
  export PUNDIX_ADDRESS="0x0000000000000000000000000000000000000000"
  "${SCRIPT_DIR}/pundix.sh" start

  "${SCRIPT_DIR}/ibcrelayer.sh" start
}

function stop() {
  "${SCRIPT_DIR}/fxcore.sh" stop
  "${SCRIPT_DIR}/pundix.sh" stop
  "${SCRIPT_DIR}/ibcrelayer.sh" stop
  rm -rf "${PROJECT_DIR}/out"
}

#trap stop EXIT

start

sleep 100