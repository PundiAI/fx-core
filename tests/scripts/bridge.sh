#!/usr/bin/env bash

set -eo pipefail

PROJECT_DIR="${PROJECT_DIR:-"$(git rev-parse --show-toplevel)"}"
export PROJECT_DIR
export OUT_DIR="${PROJECT_DIR}/out"

readonly bridger_start_index=2
readonly bridger_oracle_number=3
readonly bridge_image="functionx/fx-bridge-golang:latest"

export NODE_HOME="$OUT_DIR/.fxcore"

function create_oracles() {
  local chain_name=("$@")
  local index=${bridger_start_index}

  for chain in "${chain_name[@]}"; do
    local oracle_file="$OUT_DIR/$chain-bridge-oracle.json"
    if [ -f "$oracle_file" ]; then
      continue
    fi
    oracles=$($DAEMON query "$chain" oracles | jq -r '.oracles[]')
    if [ ${#oracles} -gt 0 ]; then
      echo "oracles already exist"
      continue
    fi

    local oracle_file="$OUT_DIR/$chain-bridge-oracle.json"
    echo "[]" >"$oracle_file"
    for ((i = 0; i < "$bridger_oracle_number"; i++)); do
      oracle_address=$(add_key "$chain-oracle-$i" "$index" | jq -r ".address")
      bridger_address=$(add_key "$chain-bridger-$i" "$((index + 1))" | jq -r ".address")
      external_address=$(show_address "$chain-bridger-$i" -e)

      jq -cs add "$oracle_file" <(echo "[{\"oracle_name\":\"$chain-oracle-$i\",\"oracle_address\":\"$oracle_address\",\"bridge_name\":\"$chain-bridger-$i\",\"bridge_address\":\"$bridger_address\",\"external_address\":\"$external_address\",\"oracle_index\":\"$index\",\"bridge_index\":\"$((index + 1))\"}]") >"$oracle_file.tmp" &&
        mv "$oracle_file.tmp" "$oracle_file"

      index=$((index + 2))
    done
  done
}

function update_crosschain_oracles() {
  local chain_name=("$@")
  for chain in "${chain_name[@]}"; do
    local oracle_file="$OUT_DIR/$chain-bridge-oracle.json"
    if [ ! -f "$oracle_file" ]; then
      continue
    fi

    oracles=()
    while read -r oracle_address; do
      oracles+=("$oracle_address")
    done < <(jq -r '.[] | "\(.oracle_address)"' "$oracle_file")

    oracles_list=$(
      IFS=,
      echo "${oracles[*]}"
    )
    min_deposit=$("${PROJECT_DIR}/tests/scripts/tx-proposal.sh" query_min_deposit)
    proposal_id=$(cosmos_tx "$chain" update-crosschain-oracles "$oracles_list" --deposit="$min_deposit" --title="Update $chain chain oracles" --description="oracles description" --from "$FROM" --home "$NODE_HOME" -y | jq -r '.logs[0].events[]|select(.type=="proposal_deposit")|.attributes[1].value')
    "${PROJECT_DIR}/tests/scripts/tx-proposal.sh" vote yes "$proposal_id"
  done
}

function create_oracle_bridger() {
  local chain_name=("$@")
  for chain in "${chain_name[@]}"; do
    local oracle_file="$OUT_DIR/$chain-bridge-oracle.json"
    if [ ! -f "$oracle_file" ]; then
      continue
    fi
    min_deposit=$("${PROJECT_DIR}/tests/scripts/tx-proposal.sh" query_min_deposit)
    validator_address=$(show_val_address "$FROM" -a)

    while read -r oracle_address oracle_name bridge_address external_address; do
      cosmos_tx "$chain" create-oracle-bridger "$validator_address" "$bridge_address" "$external_address" "$min_deposit" --from "$oracle_name" --home "$NODE_HOME" -y
    done < <(jq -r '.[] | "\(.oracle_address) \(.oracle_name) \(.bridge_address) \(.external_address)"' "$oracle_file")
  done
}

function setup_bridge_server() {
  LOCAL_IP=$(ifconfig | grep 'inet ' | grep -v '\.0\.' | awk '(NR==1){print $2}')

  local bridge_contract_file="${OUT_DIR}/bridge_contract.json"
  local chain_name=("$@")

  cat >"$OUT_DIR/bridge-docker-compose.yml" <<EOF
version: "3"

services:
EOF

  for chain in "${chain_name[@]}"; do
    local oracle_file="$OUT_DIR/$chain-bridge-oracle.json"
    bridge_contract_address=$(jq -r '.[]|select(.chain_name=="$chain")|.bridge_contract_address' "$bridge_contract_file")

    while read -r bridge_index bridge_name; do
      cat >>"$OUT_DIR/bridge-docker-compose.yml" <<EOF
    fx-$chain-bridge-$bridge_name:
      container_name: fx-$chain-bridge
      image: $bridge_image
      hostname: fx-$chain-bridge
      command: --chain-name="$chain" --external-jsonrpc="http://$LOCAL_IP:8535" --external-key="$TEST_MNEMONIC" --external-index="$bridge_index" --fx-bridge-addr="$bridge_contract_address" --fx-gas-price=4000000000000FX --fx-grpc="$LOCAL_IP:9090" --fx-key="$TEST_MNEMONIC" --fx-index="$bridge_index"
      networks:
        - bridge

EOF

    done < <(jq -r '.[] | "\(.bridge_index) \(.bridge_name)"' "$oracle_file")
  done
}

function run_test() {
  "$PROJECT_DIR"/tests/scripts/fxcore.sh init

  "$PROJECT_DIR"/tests/scripts/fxcore.sh start

  LOCAL_PORT=8535 "$PROJECT_DIR"/tests/scripts/contract.sh start

  LOCAL_PORT=8535 "$PROJECT_DIR"/tests/scripts/contract.sh deploy_bridge_contract

  create_oracles "eth" "bsc"

  update_crosschain_oracles "eth" "bsc"

  create_oracle_bridger "eth" "bsc"

  LOCAL_PORT=8535 "$PROJECT_DIR"/tests/scripts/contract.sh init_bridge_contract

  setup_bridge_server "eth" "bsc"
}

function end_test() {
  "$PROJECT_DIR"/tests/scripts/fxcore.sh stop
  "$PROJECT_DIR"/tests/scripts/contract.sh stop
}

. "${PROJECT_DIR}/tests/scripts/setup-env.sh"
