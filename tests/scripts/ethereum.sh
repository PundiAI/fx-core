#!/usr/bin/env bash

set -eo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

readonly solidity_dir="${PROJECT_DIR}/solidity"
readonly bridge_config_file="${PROJECT_DIR}/tests/data/bridge.json"
readonly bridge_config_out_file="${OUT_DIR}/bridge_contract.json"

export LOCAL_PORT=${LOCAL_PORT:-"8545"}
export LOCAL_URL="http://127.0.0.1:$LOCAL_PORT"
export REST_RPC=${REST_RPC:-"http://127.0.0.1:1317"}
export MNEMONIC=${MNEMONIC:-"test test test test test test test test test test test junk"}

function start() {
  (
    cd "$solidity_dir" || exit 1
    yarn install >/dev/null 2>&1

    nohup npx hardhat node --port "$LOCAL_PORT" >"${OUT_DIR}/hardhat.log" 2>&1 &
  )
}

function exec() {
  (
    cd "$solidity_dir" || exit 1
    yarn install >/dev/null 2>&1

    npx hardhat "$@"
  )
}

function stop() {
  pkill -f node
}

function deploy_bridge_contract() {
  export BRIDGE_CONFIG_FILE="${bridge_config_file}"
  export CONFIG_OUT_FILE="${bridge_config_out_file}"
  exec run scripts/deploy_bridge.ts
}

function deploy_staking_contract() {
  exec run scripts/deploy_staking.ts
}

function init_bridge_contract() {
  export CONFIG_FILE="${bridge_config_out_file}"
  exec run scripts/init_bridge.ts
}

## ARGS: <to> <function> [params...] Example: send 0x.... transfer(address,uint256) 0x.... 1
function send() {
  exec send "$@" --mnemonic "$MNEMONIC"
}

## ARGS: <contract> <function> [params...] Example: call 0x.... balanceOf(address) 0x....
function call() {
  exec call "$@"
}

## ARGS: <bridge-contract> <bridge-token> <amount> <destination> <target-ibc> [opts...]
function send_to_fx() {
  local bridge_contract=${1} bridge_token=${2} amount=${3} destination=${4} target_ibc=${5:-""}
  shift 5
  exec send-to-fx --bridge-contract "$bridge_contract" --bridge-token "$bridge_token" --amount "$amount" --destination "$destination" --target-ibc "$target_ibc" --mnemonic "$MNEMONIC" --disable-confirm true "$@"
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
