#!/usr/bin/env bash

set -eo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

export BRIDGE_TOKENS_OUT_DIR=${BRIDGE_TOKEN_OUT_DIR:-"${OUT_DIR}/bridge_tokens_out.json"}
export BRIDGE_CONTRACTS_OUT_DIR=${BRIDGE_CONTRACTS_OUT_DIR:-"${OUT_DIR}/bridge_contracts_out.json"}
export LOCAL_URL=${LOCAL_URL:-"http://localhost:8545"}

readonly crosschain_config_file="${PROJECT_DIR}/tests/data/crosschain.json"

function ethereum_tx() {
  "${PROJECT_DIR}/tests/scripts/ethereum.sh" "$@"
}

function send_to_fx() {
  chain_name="$1"
  destination_address=$(add_key destination 99 | jq -r '.address')
  (
    bridge_contract_address=$(jq -r '.[] | select(.chain_name == "'"$chain_name"'") | "\(.bridge_proxy_address)"' "$BRIDGE_CONTRACTS_OUT_DIR")
    crosschain_list=$(jq -r '.[] | select(.chain_name == "'"$chain_name"'") | "\(.crosschain_list)"' "$crosschain_config_file")
    targets=$(echo "$crosschain_list" | jq -r '.[] | select(.method == "send_to_fx") | "\(.targets[])"')

    while read -r bridge_token_address; do
      for target in $targets; do
        ethereum_tx send_to_fx "$bridge_contract_address" "$bridge_token_address" "111111111" "$destination_address" "$target"
      done
    done < <(jq -r '.[] | select(.chain_name == "'"$chain_name"'") | "\(.bridge_token_address)"' "$BRIDGE_TOKENS_OUT_DIR")
  )
}

function send_to_external() {
  chain_name="$1"
  from_address=$(add_key destination 99 | jq -r '.address')
  destination_address=$(show_address destination -e)

  balances=$(cosmos_query bank balances "$from_address")
  echo "$balances"
}

## ARGS: <src_port> <src_channel> <receiver> <amount> [options...]
function ibc_transfer() {
  echo ""

}

## ARGS: <src_port> <src_channel> <receiver> <amount> [options...]
function fx_ibc_transfer() {
  echo ""
}

function convert_coin() {
  echo ""

}

function convert_erc20() {
  echo ""

}

function convert_denom() {
  echo ""

}

function cross_chain() {
  echo ""

}

function fip20_cross_chain() {
  echo ""

}

function crosschain() {
  while read -r chain crosschain_list; do
    method_name=$(echo "$crosschain_list" | jq -r '.method')
    $method_name "$chain"
  done < <(jq -r '.[] | "\(.chain) \(.crosschain_list[])"' "$crosschain_config_file")
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
