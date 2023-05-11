#!/usr/bin/env bash

set -eo pipefail

PROJECT_DIR="${PROJECT_DIR:-"$(git rev-parse --show-toplevel)"}"
export PROJECT_DIR
export OUT_DIR="${PROJECT_DIR}/out"

readonly a_chain_name="fxcore"
readonly b_chain_name="pundix"
readonly docker_image_ibc="functionx/ibc-relay:1.4.8"
readonly ibc_from="ibc-test1"
readonly ibc_home_dir="$OUT_DIR/.ibcrelayer"
readonly script_dir="${PROJECT_DIR}/tests/scripts"

function transfer() {
  for chain_name in $a_chain_name $b_chain_name; do
    (
      "${script_dir}/$chain_name.sh" add_key "$ibc_from" 1

      "${script_dir}/$chain_name.sh" cosmos_transfer "$ibc_from" 200
    )
  done
}

function docker_run() {
  local opts=$1
  local args=("$@")
  docker run "$opts" --name ibc-relay -v "${ibc_home_dir}":/root/.relayer --network bridge "$docker_image_ibc" \
    "${args[@]}" --home=/root/.relayer
}

function init() {
    [[ ! -d "${ibc_home_dir}" ]] && rm -rf "${ibc_home_dir}"

  a_chain_id=$(jq -r '.chain_id' "$OUT_DIR/$a_chain_name.json")
  b_chain_id=$(jq -r '.chain_id' "$OUT_DIR/$b_chain_name.json")

  docker_run --rm config init

  cat >"${ibc_home_dir}"/config/config.yaml <<EOF
global:
  api-listen-addr: :5183
  timeout: 10s
  light-cache-size: 20
chains:
-
  chain-id: "$a_chain_id"
  key-type: "mnemonic"
  key-value: "$TEST_MNEMONIC"
  pub-key-type: "ethermint/PubKeyEthSecp256k1"
  hd-path: "m/44'/60'/0'/0/1"

  rpc-addr: "http://${a_chain_name}:26657"
  account-prefix: $(jq -r '.bech32_prefix' "$OUT_DIR/$a_chain_name.json")
  gas-adjustment: 1.01
  gas-prices: "$(jq -r '.gas_prices' "$OUT_DIR/$a_chain_name.json")"
  trusting-period: 29m
  skip-un-relay-sequences: []
  iterator-block-config:
    delay-second: 3
    handler-block-count: 100
    batch-handle-block-count: 15
-
  chain-id: "$b_chain_id"
  key-type: "mnemonic"
  key-value: "$TEST_MNEMONIC"
  pub-key-type: "tendermint/PubKeySecp256k1"
  hd-path: "m/44'/118'/0'/0/1"

  rpc-addr: "http://${b_chain_name}:26657"
  account-prefix: $(jq -r '.bech32_prefix' "$OUT_DIR/$b_chain_name.json")
  gas-adjustment: 1.01
  gas-prices: "$(jq -r '.gas_prices' "$OUT_DIR/$b_chain_name.json")"
  trusting-period: 29m
  skip-un-relay-sequences: []
  iterator-block-config:
    delay-second: 3
    handler-block-count: 100
    batch-handle-block-count: 15

paths: {}
EOF

  docker_run --rm paths generate "$a_chain_id" "$b_chain_id" transfer --port=transfer

  docker_run --rm light init "$a_chain_id" -f

  docker_run --rm light init "$b_chain_id" -f

  docker_run --rm tx link transfer -d
}

function start() {
  docker_run -itd start transfer --time-threshold=19m --notify.enable=false --debug=true
}

function stop() {
  docker stop ibc-relay
  docker rm ibc-relay
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"
