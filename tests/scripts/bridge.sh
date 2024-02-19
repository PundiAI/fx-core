#!/usr/bin/env bash

set -eo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

readonly bridger_start_index=100
readonly bridger_oracle_number=3
readonly bridge_tokens_file="${PROJECT_DIR}/tests/data/bridge_tokens.json"

export NODE_HOME="$OUT_DIR/.fxcore"
export LOCAL_PORT=${LOCAL_PORT:-"8545"}

export BRIDGE_IMAGE=${BRIDGE_IMAGE:-"functionx/fx-bridge-robot:latest"}
export BRIDGE_CONTRACTS_OUT_DIR=${BRIDGE_CONTRACTS_OUT_DIR:-"${OUT_DIR}/bridge_contracts_out.json"}
export BRIDGE_TOKENS_OUT_DIR=${BRIDGE_TOKENS_OUT_DIR:-"${OUT_DIR}/bridge_tokens_out.json"}
export RBIDGER_OUT_DIR=${RBIDGER_OUT_DIR:-"${OUT_DIR}/bridger_key/"}

mkdir -p "$RBIDGER_OUT_DIR"

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
      echo "12345678" | export_key "$chain-bridger-$i" "$RBIDGER_OUT_DIR/fx-$chain-bridger-$i.key" --ascii-armor
      echo "12345678" | export_key "$chain-bridger-$i" "$RBIDGER_OUT_DIR/eth-$chain-bridger-$i.key"
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

      cosmos_tx "$chain" bounded-oracle "$validator_address" "$bridge_address" "$external_address" "$min_deposit" --from "$oracle_name"
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
    bridge_contract_address=$(jq --arg chain_name "$chain" -r '.[] | select(.chain_name=="'"$chain"'") | .bridge_proxy_address' "$BRIDGE_CONTRACTS_OUT_DIR")

    while read -r bridge_index bridge_name; do
      cat >>"$OUT_DIR/bridge-docker-compose.yml" <<EOF
    fx-$bridge_name:
      container_name: fx-$bridge_name
      image: $BRIDGE_IMAGE
      hostname: fx-$bridge_name
      command: bridge --bridge-chain-name="$chain" --bridge-node-url="$external_json_rpc_url" --bridge-addr="$bridge_contract_address" --fx-gas-price=4000000000000FX --fx-grpc="http://fxcore:9090" --bridge-key=/root/eth-$bridge_name.key --bridge-pwd=12345678 --fx-key=/root/fx-$bridge_name.key --fx-pwd=12345678
      networks:
        - $DOCKER_NETWORK
      volumes:
        - $RBIDGER_OUT_DIR:/root
EOF

    done < <(jq -r '.[] | "\(.bridge_index) \(.bridge_name)"' "$oracle_file")
  done

  cat >>"$OUT_DIR/bridge-docker-compose.yml" <<EOF
networks:
  $DOCKER_NETWORK:
    external: true
EOF

  docker-compose -f "$OUT_DIR/bridge-docker-compose.yml" up -d
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

function register_coin() {
  while read -r bridge_chains symbol decimals base_denom target_ibc name; do
    [[ "$symbol" == FX ]] && continue
    for bridge_chain in "${bridge_chains[@]}"; do
      local alias_list=()
      for chain_name in $(echo "$bridge_chain" | jq -r '.[]'); do
        token_address=$(jq -r '.[] | select(.chain_name=="'"$chain_name"'") | select(.symbol=="'"$symbol"'") |.bridge_token_address' "$BRIDGE_TOKENS_OUT_DIR")
        denom="${chain_name}${token_address}"

        [[ "$target_ibc" != "null" ]] && denom=$(convert_ibc_denom "${target_ibc}/${denom}")

        if [[ "$base_denom" == "null" ]]; then
          base_denom="$denom"
          break
        fi
        alias_list+=("\"$denom\"")
      done
      proposal register_coin "$base_denom" "$name" "$symbol" "$decimals" "${alias_list[@]}"
    done
  done < <(jq -r '.[] | "\(.bridge_chains) \(.symbol) \(.decimals) \(.base_denom) \(.target_ibc) \(.name)"' "$bridge_tokens_file")
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
