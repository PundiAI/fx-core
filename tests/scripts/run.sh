#!/usr/bin/env bash

set -eo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

function run() {
  (
    # shellcheck source=/dev/null
    . "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"
    create_docker_network
  )

  (
    export LOCAL_PORT=8535
    "${script_dir}/ethereum.sh" start
    "${script_dir}/ethereum.sh" deploy_bridge_token
  )

  (
    "${script_dir}/fxcore.sh" init
    "${script_dir}/fxcore.sh" start
  )

  (
    export IBC_CHANNEL="channel-0"
    PURSE_ADDRESS="$("${script_dir}/ethereum.sh" get_token_address bsc PURSE)"
    export PURSE_ADDRESS
    PUNDIX_ADDRESS="$("${script_dir}/ethereum.sh" get_token_address eth PUNDIX)"
    export PUNDIX_ADDRESS
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
