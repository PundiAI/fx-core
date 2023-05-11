#!/usr/bin/env bash

set -eo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

readonly solidity_dir="${PROJECT_DIR}/solidity"
readonly bridge_config_file="${PROJECT_DIR}/tests/data/bridge.json"
readonly bridge_config_out_file="${OUT_DIR}/bridge_contract.json"

export LOCAL_PORT=${LOCAL_PORT:-"8545"}
export LOCAL_URL=${LOCAL_URL:-"http://127.0.0.1:$LOCAL_PORT"}
export REST_RPC=${REST_RPC:-"http://127.0.0.1:1317"}
export MNEMONIC=${MNEMONIC:-"test test test test test test test test test test test junk"}

function start() {
  (
    cd "$solidity_dir" || exit 1
    yarn install >/dev/null 2>&1

    nohup npx hardhat node --port "$LOCAL_PORT" >"${OUT_DIR}/hardhat.log" 2>&1 &
  )
}

function deploy_bridge_contract() {
  (
    cd "$solidity_dir" || exit 1
    yarn install >/dev/null 2>&1

    export BRIDGE_CONFIG_FILE="${bridge_config_file}"
    export CONFIG_OUT_FILE="${bridge_config_out_file}"
    npx hardhat run scripts/deploy_bridge.ts
  )
}

function deploy_staking_contract() {
  (
    cd "$project_dir/solidity" || exit 1
    yarn install
    npx hardhat run scripts/deploy_staking.ts
  )
}

function init_bridge_contract() {
  (
    cd "$solidity_dir" || exit 1
    yarn install >/dev/null 2>&1

    export CONFIG_FILE="${bridge_config_out_file}"
    npx hardhat run scripts/init_bridge.ts
  )
}

function send() {
  index=${1:-"0"}
  shift
  (
    cd "$solidity_dir" || exit 1
    yarn install >/dev/null 2>&1

    npx hardhat send "$@" --mnemonic "$MNEMONIC" --index "$index"
  )
}

function call() {
  (
    cd "$solidity_dir" || exit 1
    yarn install >/dev/null 2>&1

    npx hardhat call "$@"
  )
}

function stop() {
  pkill -f node
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
