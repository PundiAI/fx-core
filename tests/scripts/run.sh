#!/usr/bin/env bash

set -eo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

function run() {
  (
    . "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"
    create_docker_network
  )

  (
    "${script_dir}/fxcore.sh" init
    "${script_dir}/fxcore.sh" start
  )

  (
    export IBC_CHANNEL="channel-0"
    export PURSE_ADDRESS="0x0000000000000000000000000000000000000000"
    export PUNDIX_ADDRESS="0x0000000000000000000000000000000000000000"
    "${script_dir}/pundix.sh" init
    "${script_dir}/pundix.sh" start

    "${script_dir}/ibcrelayer.sh" transfer
    "${script_dir}/ibcrelayer.sh" init
    "${script_dir}/ibcrelayer.sh" create_channel
    "${script_dir}/ibcrelayer.sh" start
  )

}

#function close() {
#  "${script_dir}/fxcore.sh" docker_stop
#  "${script_dir}/pundix.sh" docker_stop
#  "${script_dir}/ibcrelayer.sh" stop
#  rm -rf "${PROJECT_DIR}/out"
#}
#
#trap close EXIT SIGINT SIGTERM

run
