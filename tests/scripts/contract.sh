#!/usr/bin/env bash

set -eo pipefail

project_dir="$(git rev-parse --show-toplevel)"
readonly project_dir

solidity_dir="${project_dir}/solidity"
readonly bridge_config_file="${project_dir}/tests/data/bridge.json"
readonly bridge_config_out_file="${project_dir}/out/bridge_contract.json"

export REST_RPC="http://127.0.0.1:1317"

export LOCAL_PORT=${LOCAL_PORT:-"8545"}
export LOCAL_URL=${LOCAL_URL:-"http://127.0.0.1:$LOCAL_PORT"}

function start() {
  (
    cd "$solidity_dir" || exit 1
    yarn install

    nohup npx hardhat node --port "$LOCAL_PORT" &
  )
}

function deploy_bridge_contract() {
  (
    cd "$project_dir/solidity" || exit 1
    yarn install

    export BRIDGE_CONFIG_FILE="${bridge_config_file}"
    export CONFIG_OUT_FILE="${bridge_config_out_file}"
    npx hardhat run scripts/deploy_bridge.ts
  )
}

function init_bridge_contract() {
  (
    cd "$project_dir/solidity" || exit 1
    yarn install

    export CONFIG_FILE="${bridge_config_out_file}"
    npx hardhat run scripts/init_bridge.ts
  )
}

function send() {
  (
    cd "$project_dir/solidity" || exit 1
    yarn install

    npx hardhat send "$@"
  )
}

function call() {
  (
    cd "$project_dir/solidity" || exit 1
    yarn install

    npx hardhat call "$@"
  )
}

. "${project_dir}/tests/scripts/setup-env.sh"
