#!/usr/bin/env bash

set -eo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

readonly bridger_start_index=100
readonly bridger_oracle_number=3
readonly bridge_image="functionx/fx-bridge-golang:3.1.0"
readonly bridge_contract_file="${OUT_DIR}/bridge_contract.json"

export NODE_HOME="$OUT_DIR/.fxcore"
export LOCAL_PORT=${LOCAL_PORT:-"8545"}

function proposal() {
  "${PROJECT_DIR}/tests/scripts/proposal.sh" "$@"
}

function create_oracles() {
  local chain_name=("$@")
  local index=${bridger_start_index}

  for chain in "${chain_name[@]}"; do
    local oracle_file="$OUT_DIR/$chain-bridge-oracle.json"
    echo "[]" >"$oracle_file"
    for ((i = 0; i < "$bridger_oracle_number"; i++)); do
      add_key "$chain-oracle-$i" "$index"
      add_key "$chain-bridger-$i" "$((index + 1))"

      oracle_address=$(show_address "$chain-oracle-$i" -a)
      bridger_address=$(show_address "$chain-bridger-$i" -a)
      external_address=$(show_address "$chain-bridger-$i" -e)

      cat >"$oracle_file.new" <<EOF
[
  {
    "oracle_name": "$chain-oracle-$i",
    "oracle_address": "$oracle_address",
    "bridge_name": "$chain-bridger-$i",
    "bridge_address": "$bridger_address",
    "external_address": "$external_address",
    "oracle_index": "$index",
    "bridge_index": "$((index + 1))"
  }
]
EOF
      jq -cs add "$oracle_file" "$oracle_file.new" >"$oracle_file.tmp" &&
        mv "$oracle_file.tmp" "$oracle_file"
      index=$((index + 2))
    done
    rm -r "$oracle_file.new"
  done
}

function update_crosschain_oracles() {
  local chain_name=("$@")
  for chain in "${chain_name[@]}"; do
    local oracle_file="$OUT_DIR/$chain-bridge-oracle.json"
    [[ ! -f "$oracle_file" ]] && continue

    oracles_list=$(jq -r '. | map(.oracle_address) | join(",")' "$oracle_file")

    min_deposit=$(proposal query_min_deposit)
    cosmos_tx "$chain" update-crosschain-oracles "$oracles_list" --deposit="$min_deposit" --title="Update $chain chain oracles" --description="oracles description" --from "$FROM"
    proposal vote yes
  done
}

function create_oracle_bridger() {
  local chain_name=("$@")
  for chain in "${chain_name[@]}"; do
    local oracle_file="$OUT_DIR/$chain-bridge-oracle.json"
    [[ ! -f "$oracle_file" ]] && continue

    min_deposit=$(proposal query_min_deposit)
    validator_address=$(show_address "$FROM" -a --bech val)

    while read -r oracle_name oracle_address oracle_index bridge_name bridge_address bridge_index external_address; do
      add_key "$oracle_name" "$oracle_index"
      add_key "$bridge_name" "$bridge_index"

      cosmos_tx bank send "$FROM" "$oracle_address" "$min_deposit" --from "$FROM"
      cosmos_transfer "$oracle_name" 100
      cosmos_transfer "$bridge_name" 500

      cosmos_tx "$chain" create-oracle-bridger "$validator_address" "$bridge_address" "$external_address" "$min_deposit" --from "$oracle_name"
    done < <(jq -r '.[] | "\(.oracle_name) \(.oracle_address) \(.oracle_index) \(.bridge_name) \(.bridge_address) \(.bridge_index) \(.external_address)"' "$oracle_file")
  done
}

function setup_bridge_server() {
  local external_json_rpc_url="http://host.docker.internal:$LOCAL_PORT"
  local chain_name=("$@")

  cat >"$OUT_DIR/bridge-docker-compose.yml" <<EOF
version: "3"

services:
EOF

  for chain in "${chain_name[@]}"; do
    local oracle_file="$OUT_DIR/$chain-bridge-oracle.json"

    bridge_contract_address=$(jq --arg chain_name "$chain" -r '.external_chain_list[]|select(.chain_name==$chain_name).bridge_contract_address' "$bridge_contract_file")

    while read -r bridge_index bridge_name; do
      cat >>"$OUT_DIR/bridge-docker-compose.yml" <<EOF
    fx-$bridge_name:
      container_name: fx-$bridge_name
      image: $bridge_image
      hostname: fx-$bridge_name
      command: --chain-name="$chain" --external-jsonrpc="$external_json_rpc_url" --external-key="$TEST_MNEMONIC" --external-index="$bridge_index" --fx-bridge-addr="$bridge_contract_address" --fx-gas-price=4000000000000FX --fx-grpc="http://fxcore:9090" --fx-key="$TEST_MNEMONIC" --fx-index="$bridge_index"
      networks:
        - $DOCKER_NETWORK
EOF

    done < <(jq -r '.[] | "\(.bridge_index) \(.bridge_name)"' "$oracle_file")
  done

  cat >>"$OUT_DIR/bridge-docker-compose.yml" <<EOF
networks:
  $DOCKER_NETWORK:
    external: true
EOF
}

function request_batch() {
  local chain_name=("$@")
  for chain in "${chain_name[@]}"; do
    local oracle_file="$OUT_DIR/$chain-bridge-oracle.json"
    [[ ! -f "$oracle_file" ]] && continue

    length=$(cosmos_query "$chain" batch-fees | jq '.batch_fees | length')
    [[ "$length" -eq 0 ]] && continue

    bridge_index=$(jq -r '.[0].bridge_index' "$oracle_file")
    add_key "$chain-bridger-0" "$bridge_index"

    while read -r token_contract; do
      denom=$(cosmos_query "$chain" denom "$token_contract" | jq -r '.denom')
      cosmos_tx "$chain" build-batch "$denom" "1" "1" "$(show_address "$chain-bridger-0" -e)" --from "$chain-bridger-0"
    done < <(cosmos_query "$chain" batch-fees | jq -r '.batch_fees[] | "\(.token_contract)"')
  done
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
