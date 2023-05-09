#!/usr/bin/env bash

set -eo pipefail

readonly project_dir="$(git rev-parse --show-toplevel)"
readonly bridge_config_file="${project_dir}/tests/data/bridge.json"
readonly bridge_config_out_file="${project_dir}/out/bridge_contract.json"

readonly json_rpc_port="8535"
readonly hardhat_network="localhost"
readonly rest_rpc_url="http://127.0.0.1:1317"

export LOCAL_URL="http://127.0.0.1:$json_rpc_port"

function start() {
  (
    cd "$project_dir/solidity" || exit 1
    yarn install

    nohup npx hardhat node --port "$json_rpc_port" &
  )
}

function deploy_bridge_contract() {
  (
    cd "$project_dir/solidity" || exit 1
    yarn install

    export BRIDGE_CONFIG_FILE="${bridge_config_file}"
    export CONFIG_OUT_FILE="${bridge_config_out_file}"
    npx hardhat run scripts/deploy_bridge.ts --network "$hardhat_network"
  )
}

function init_bridge_contract() {
  (
    cd "$project_dir/solidity" || exit 1
    yarn install

    export CONFIG_FILE="${bridge_config_out_file}"
    export REST_RPC="${rest_rpc_url}"
    npx hardhat run scripts/init_bridge.ts --network "$hardhat_network"
  )
}

. "${project_dir}/tests/scripts/setup-env.sh"
