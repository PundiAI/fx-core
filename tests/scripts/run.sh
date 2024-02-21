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
    LOCAL_PORT=8535
    export LOCAL_PORT
    "${script_dir}/ethereum.sh" start
    "${script_dir}/ethereum.sh" deploy_bridge_token
    "${script_dir}/ethereum.sh" deploy_bridge_contract
  )

  (
    "${script_dir}/fxcore.sh" init
    "${script_dir}/fxcore.sh" start
  )

  (
    # shellcheck source=/dev/null
    . "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"
    if [[ "$ENABLE_IBC" == "true" ]]; then
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
    fi
  )

  (
    "${script_dir}/bridge.sh" create_oracles eth
    "${script_dir}/bridge.sh" update_crosschain_oracles eth
    "${script_dir}/bridge.sh" create_oracle_bridger eth

    LOCAL_PORT=8535
    export LOCAL_PORT
    "${script_dir}/ethereum.sh" init_bridge
    "${script_dir}/ethereum.sh" add_bridge_tokens
    "${script_dir}/bridge.sh" setup_bridge_server eth
    "${script_dir}/bridge.sh" register_coin
  )

  (
    LOCAL_PORT=8545
    export LOCAL_PORT
    "${script_dir}/ethereum.sh" deploy_bridge_call_contract
    LOCAL_PORT=8535
    export LOCAL_PORT
    "${script_dir}/ethereum.sh" bridge_erc20_call_test eth
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
