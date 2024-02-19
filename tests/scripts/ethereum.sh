#!/usr/bin/env bash

set -eo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

readonly solidity_dir="${PROJECT_DIR}/solidity"
readonly bridge_contracts_file="${PROJECT_DIR}/tests/data/bridge_contracts.json"
readonly bridge_tokens_file="${PROJECT_DIR}/tests/data/bridge_tokens.json"

export LOCAL_PORT=${LOCAL_PORT:-"8545"}
export LOCAL_URL=${LOCAL_URL:-"http://127.0.0.1:$LOCAL_PORT"}
export REST_RPC=${REST_RPC:-"http://127.0.0.1:1317"}
export MNEMONIC=${MNEMONIC:-"test test test test test test test test test test test junk"}

export BRIDGE_TOKENS_OUT_DIR=${BRIDGE_TOKEN_OUT_DIR:-"${OUT_DIR}/bridge_tokens_out.json"}
export BRIDGE_CONTRACTS_OUT_DIR=${BRIDGE_CONTRACTS_OUT_DIR:-"${OUT_DIR}/bridge_contracts_out.json"}

function start() {
  stop
  (
    cd "$solidity_dir" || exit 1
    yarn install >/dev/null 2>&1

    nohup npx hardhat node --port "$LOCAL_PORT" >"${OUT_DIR}/hardhat.log" 2>&1 &
    sleep 2
  )
}

function hardhat_task() {
  (
    cd "$solidity_dir" || exit 1
    yarn install >/dev/null 2>&1

    npx hardhat "$@" --network localhost
  )
}

function stop() {
  pgrep -f "hardhat node" | xargs kill -9
}

## ARGS: <contract-name> [constructor-params...] Example: deploy_contract ERC20TokenTest "TestToken" "TT" "18" "10000000"
function deploy_contract() {
  hardhat_task deploy-contract --contract-name "$@" --mnemonic "$MNEMONIC" --disable-confirm "true"
}

## ARGS: <contract-logic> <contract-proxy> <rest-rpc> <chain-name>
function init_bridge_contract() {
  local logic=${1} proxy=${2} rest_url=${3} chain_name=${4}
  shift 4
  hardhat_task init-bridge --bridge-logic "$logic" --bridge-contract "$proxy" \
    --rest-url "$rest_url" --chain-name "$chain_name" --mnemonic "$MNEMONIC" --disable-confirm "true" "$@"
}

## ARGS: <bridge-contract> <bridge-token> <is-original> <target-ibc>
# shellcheck disable=SC2317  # Don't warn about unreachable commands in this function
function add_bridge_token() {
  local contract=${1} token=${2} is_original=${3} target_ibc=${4}
  shift 4
  hardhat_task add-bridge-token --bridge-contract "$contract" --token-contract "$token" \
    --is-original "$is_original" --target-ibc "$target_ibc" --mnemonic "$MNEMONIC" --disable-confirm "true" "$@"
}

## ARGS: <to> <function> [params...] Example: send 0x.... transfer(address,uint256) 0x.... 1
function send() {
  hardhat_task send "$@" --mnemonic "$MNEMONIC" --disable-confirm "true"
}

## ARGS: <contract> <function> [params...] Example: call 0x.... balanceOf(address) 0x....
function call() {
  hardhat_task call "$@"
}

## ARGS: <bridge-contract> <bridge-token> <amount> <destination> <target-ibc> [opts...]
function send_to_fx() {
  local bridge_contract=${1} bridge_token=${2} amount=${3} destination=${4} target_ibc=${5:-""}
  shift 5
  hardhat_task send-to-fx --bridge-contract "$bridge_contract" --bridge-token "$bridge_token" --amount "$amount" --destination "$destination" --target-ibc "$target_ibc" --mnemonic "$MNEMONIC" --disable-confirm "true" "$@"
}

function deploy_bridge_contract() {
  echo "[]" >"$BRIDGE_CONTRACTS_OUT_DIR"
  add_key "$FROM" 0
  while read -r chain_name contract_class_name; do
    external_address=$(show_address "$FROM" -e)

    logic_address=$(deploy_contract "$contract_class_name")
    proxy_address=$(deploy_contract "TransparentUpgradeableProxy" "$logic_address" "$external_address" "0x")

    cat >"$BRIDGE_CONTRACTS_OUT_DIR.new" <<EOF
[
  {
    "chain_name": "$chain_name",
    "bridge_logic_address": "$logic_address",
    "bridge_proxy_address": "$proxy_address"
  }
]
EOF
    jq -cs add "$BRIDGE_CONTRACTS_OUT_DIR" "$BRIDGE_CONTRACTS_OUT_DIR.new" >"$BRIDGE_CONTRACTS_OUT_DIR.tmp" &&
      mv "$BRIDGE_CONTRACTS_OUT_DIR.tmp" "$BRIDGE_CONTRACTS_OUT_DIR"
  done < <(jq -r '.[] | "\(.chain_name) \(.contract_class_name)"' "$bridge_contracts_file")
  rm -r "$BRIDGE_CONTRACTS_OUT_DIR.new"
}

function deploy_bridge_token() {
  echo "[]" >"$BRIDGE_TOKENS_OUT_DIR"

  while read -r bridge_chains symbol decimals total_supply is_original target_ibc name; do
    for bridge_chain in "${bridge_chains[@]}"; do
      for chain_name in $(echo "$bridge_chain" | jq -r '.[]'); do
        erc20_address=$(deploy_contract "ERC20TokenTest" "$name" "$symbol" "$decimals" "$total_supply")

        cat >"$BRIDGE_TOKENS_OUT_DIR.new" <<EOF
[
  {
    "chain_name": "$chain_name",
    "symbol": "$symbol",
    "bridge_token_address": "$erc20_address",
    "target_ibc": "$target_ibc",
    "is_original": "$is_original"
  }
]
EOF
        jq -cs add "$BRIDGE_TOKENS_OUT_DIR" "$BRIDGE_TOKENS_OUT_DIR.new" >"$BRIDGE_TOKENS_OUT_DIR.tmp" &&
          mv "$BRIDGE_TOKENS_OUT_DIR.tmp" "$BRIDGE_TOKENS_OUT_DIR"
      done
    done
  done < <(jq -r '.[] | "\(.bridge_chains) \(.symbol) \(.decimals) \(.total_supply) \(.is_original) \(.target_ibc) \(.name)"' "$bridge_tokens_file")
  rm -r "$BRIDGE_TOKENS_OUT_DIR.new"
}

function get_token_address() {
  chain_name=$1 symbol=$2
  jq -r '.[] | select(.chain_name == "'"$chain_name"'") | select(.symbol == "'"$symbol"'") | .bridge_token_address' "$BRIDGE_TOKENS_OUT_DIR"
}

function init_bridge() {
  while read -r chain_name bridge_logic_address bridge_proxy_address; do
    init_bridge_contract "$bridge_logic_address" "$bridge_proxy_address" "$REST_RPC" "$chain_name"
  done < <(jq -r '.[] | "\(.chain_name) \(.bridge_logic_address) \(.bridge_proxy_address)"' "$BRIDGE_CONTRACTS_OUT_DIR")
}

function add_bridge_token() {
  while read -r chain_name bridge_proxy_address; do
    while read -r bridge_token_address is_original target_ibc; do
      if [ "$target_ibc" == "null" ]; then
        target_ibc=""
      fi
      add_bridge_token "$bridge_proxy_address" "$bridge_token_address" "$is_original" "$target_ibc"
    done < <(jq -r '.[] | select(.chain_name == "'"$chain_name"'") | "\(.bridge_token_address) \(.is_original) \(.target_ibc)"' "$BRIDGE_TOKENS_OUT_DIR")
  done < <(jq -r '.[] | "\(.chain_name) \(.bridge_proxy_address)"' "$BRIDGE_CONTRACTS_OUT_DIR")
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
