#!/usr/bin/env bash

set -eo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

a_chain_name=${IBC_A_CHAIN_NAME:-"fxcore"}
b_chain_name=${IBC_B_CHAIN_NAME:-"pundix"}
account_index=${IBC_ACCOUNT_INDEX:-"1"}
readonly ibc_from="ibc-$FROM"

readonly docker_name="ibc-relay"
readonly ibc_home_dir="$OUT_DIR/.hermes"

export SCRIPT_DIR=${SCRIPT_DIR:-"$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"}
export IBC_DOCKER_IMAGE=${IBC_DOCKER_IMAGE:-"ghcr.io/informalsystems/hermes:1.3.0"}

function transfer() {
  for chain_name in $a_chain_name $b_chain_name; do
    (
      "${SCRIPT_DIR}/$chain_name.sh" add_key "$ibc_from" "${account_index}"

      "${SCRIPT_DIR}/$chain_name.sh" cosmos_transfer "$ibc_from" 200
    )
  done
}

function docker_run() {
  local opts=$1
  shift
  local args=("$@")
  local name=$docker_name

  [[ "$opts" == "--rm" ]] && name="$docker_name-tmp"
  IFS=' ' read -r -a opts_ary <<<"$opts"
  docker run "${opts_ary[@]}" --name "$name" --network "$DOCKER_NETWORK" -v "${ibc_home_dir}":/home/hermes/.hermes \
    "$IBC_DOCKER_IMAGE" "${args[@]}"
}

function init() {
  [[ ! -d "${ibc_home_dir}" ]] && docker_stop "$docker_name" && rm -rf "${ibc_home_dir}"
  mkdir -p "$OUT_DIR/.hermes/keys"

  a_chain_id=$(jq -r '.chain_id' "$OUT_DIR/$a_chain_name.json")
  b_chain_id=$(jq -r '.chain_id' "$OUT_DIR/$b_chain_name.json")

  a_gas_price=$(jq -r '.gas_prices' "$OUT_DIR/$a_chain_name.json")
  a_fee_denom=$(jq -r '.fee_denom' "$OUT_DIR/$a_chain_name.json")
  [[ -z "$a_fee_denom" || "$a_fee_denom" == "null" ]] && a_fee_denom=$(jq -r '.staking_denom' "$OUT_DIR/$a_chain_name.json")
  a_gas_price=${a_gas_price%"${a_fee_denom}"}
  b_gas_price=$(jq -r '.gas_prices' "$OUT_DIR/$b_chain_name.json")
  b_fee_denom=$(jq -r '.fee_denom' "$OUT_DIR/$b_chain_name.json")
  [[ -z "$b_fee_denom" || "$b_fee_denom" == "null" ]] && b_fee_denom=$(jq -r '.staking_denom' "$OUT_DIR/$b_chain_name.json")
  b_gas_price=${b_gas_price%"${b_fee_denom}"}

  a_rpc_addr=$(jq -r ".node_rpc" "$OUT_DIR/$a_chain_name.json")
  a_websocket_addr=${a_rpc_addr/http/ws}/websocket
  b_rpc_addr=$(jq -r ".node_rpc" "$OUT_DIR/$b_chain_name.json")
  b_websocket_addr=${b_rpc_addr/http/ws}/websocket

  a_grpc_addr=$(jq -r ".node_grpc" "$OUT_DIR/$a_chain_name.json")
  b_grpc_addr=$(jq -r ".node_grpc" "$OUT_DIR/$b_chain_name.json")

  a_trusting_period=$(jq -r ".trusting_period" "$OUT_DIR/$a_chain_name.json")
  b_trusting_period=$(jq -r ".trusting_period" "$OUT_DIR/$b_chain_name.json")
  [[ -z "$a_trusting_period" || "$a_trusting_period" == "null" ]] && a_trusting_period="20days"
  [[ -z "$b_trusting_period" || "$b_trusting_period" == "null" ]] && b_trusting_period="20days"

  # config: https://hermes.informal.systems/documentation/configuration/description.html
  cat >"${ibc_home_dir}"/config.toml <<EOF
[global]
log_level = 'info'

[mode]

[mode.clients]
enabled = true
refresh = true
misbehaviour = false

[mode.connections]
enabled = false

[mode.channels]
enabled = false

[mode.packets]
enabled = true
clear_interval = 100
clear_on_start = true
tx_confirmation = false
auto_register_counterparty_payee = false

[rest]
enabled = false
host = '127.0.0.1'
port = 3000

[telemetry]
enabled = false
host = '127.0.0.1'
port = 3001

[[chains]]

# Specify the chain ID. Required
id = '$a_chain_id'
rpc_addr = '$a_rpc_addr'
grpc_addr = '$a_grpc_addr'
websocket_addr = '$a_websocket_addr'
rpc_timeout = '10s'
account_prefix = '$(jq -r ".bech32_prefix" "$OUT_DIR/$a_chain_name.json")'
key_name = "$ibc_from"
address_type = { derivation = 'ethermint', proto_type = { pk_type = '/ethermint.crypto.v1.ethsecp256k1.PubKey' } }
store_prefix = 'ibc'
default_gas = 100000
max_gas = 800000
gas_price = { price = ${a_gas_price}, denom = '$a_fee_denom' }
gas_multiplier = 1.1
max_msg_num = 30
max_tx_size = 1048576
clock_drift = '5s'
max_block_time = '30s'
trusting_period = '$a_trusting_period'
trust_threshold = { numerator = '2', denominator = '3' }
memo_prefix = ''
[chains.packet_filter]
policy = 'allow'
list = [
 ['transfer', 'channel-0'],
]

[[chains]]
id = '$b_chain_id'
rpc_addr = '$b_rpc_addr'
grpc_addr = '$b_grpc_addr'
websocket_addr = '$b_websocket_addr'
rpc_timeout = '10s'
account_prefix = '$(jq -r ".bech32_prefix" "$OUT_DIR/$b_chain_name.json")'
key_name = "$ibc_from"
address_type = { derivation = 'cosmos' }
store_prefix = 'ibc'
default_gas = 100000
max_gas = 400000
gas_price = { price = ${b_gas_price}, denom = '$b_fee_denom' }
gas_multiplier = 1.1
max_msg_num = 30
max_tx_size = 2097152
clock_drift = '5s'
max_block_time = '30s'
trusting_period = '$b_trusting_period'
trust_threshold = { numerator = '2', denominator = '3' }
[chains.packet_filter]
policy = 'allow'
list = [
 ['transfer', 'channel-0'],
]

EOF

  import_key "${a_chain_id}" "${ibc_from}" "m/44'/60'/0'/0/${account_index}"
  import_key "${b_chain_id}" "${ibc_from}" "m/44'/118'/0'/0/${account_index}"

  config_check

  a_chain_account=$(jq -r ".account" "$OUT_DIR/.hermes/keys/$a_chain_id/keyring-test/$ibc_from.json")
  b_chain_account=$(jq -r ".account" "$OUT_DIR/.hermes/keys/$b_chain_id/keyring-test/$ibc_from.json")

  echo -e "\n====================== ibc relay account ======================"
  echo -e "$a_chain_id \t $a_chain_account"
  echo -e "$b_chain_id \t $b_chain_account"
  echo -e "====================== ibc relay account ======================\n"
  echo "please ensure the above addresses have tokens before running the create channel and start"
}

function import_key() {
  local chain_name=${1} key_name=${2} hd_path=${3}

  docker_run --rm keys delete --chain "${chain_name}" --key-name "${key_name}" >/dev/null
  local mnemonic_path="${ibc_home_dir}/mnemonic"
  echo "${TEST_MNEMONIC}" >"${mnemonic_path}"
  docker_run --rm keys add --chain "${chain_name}" --key-name "${key_name}" --hd-path="${hd_path}" --mnemonic-file ./.hermes/mnemonic
  rm "${mnemonic_path}"
}

function create_channel() {
  a_chain_id=$(jq -r '.chain_id' "$OUT_DIR/$a_chain_name.json")
  b_chain_id=$(jq -r '.chain_id' "$OUT_DIR/$b_chain_name.json")
  docker_run --rm create channel --a-chain "${a_chain_id}" --b-chain "${b_chain_id}" --a-port transfer --b-port transfer --new-client-connection --yes
}

function config_check() {
  docker_run --rm config validate
}

function health_check() {
  docker_run --rm health-check
}

function start() {
  health_check
  docker_run -d start
}

function stop() {
  docker rm -f $docker_name
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
