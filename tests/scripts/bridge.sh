#!/usr/bin/env bash

set -eo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

readonly bridger_start_index=2
readonly bridger_oracle_number=3
readonly bridge_image="functionx/fx-bridge-golang:latest"
readonly docker_network="test-net"

export NODE_HOME="$OUT_DIR/.fxcore"

function create_oracles() {
  local chain_name=("$@")
  local index=${bridger_start_index}

  for chain in "${chain_name[@]}"; do
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
    cosmos_tx "$chain" update-crosschain-oracles "$oracles_list" --deposit="$min_deposit" --title="Update $chain chain oracles" --description="oracles description" --from "$FROM" -y | jq -r '.logs[0].events[]|select(.type=="proposal_deposit")|.attributes[1].value'
    "${PROJECT_DIR}/tests/scripts/tx-proposal.sh" vote yes
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
    validator_address=$(show_address "$FROM" -a --bech val)

    while read -r oracle_address oracle_name bridge_address external_address; do
      if ! show_address "$oracle_name" -a --bech val; then
        add_key "$oracle_name" "$index"

      fi
      cosmos_tx bank send "$FROM" "$oracle_address" "$min_deposit" --from "$FROM"
      cosmos_tx bank send "$FROM" "$oracle_address" "$(to_18 "100")$STAKING_DENOM" --from "$FROM"
      cosmos_tx bank send "$FROM" "$bridge_address" "$(to_18 "500")$STAKING_DENOM" --from "$FROM"
      cosmos_tx "$chain" create-oracle-bridger "$validator_address" "$bridge_address" "$external_address" "$min_deposit" --from "$oracle_name" --node "http://127.0.0.1:26657" --home "$NODE_HOME" -y
    done < <(jq -r '.[] | "\(.oracle_address) \(.oracle_name) \(.bridge_address) \(.external_address)"' "$oracle_file")
  done
}

function setup_bridge_server() {
  local bridge_contract_file="${OUT_DIR}/bridge_contract.json"
  local external_json_rpc_url="http://host.docker.internal:$LOCAL_PORT"
  local chain_name=("$@")

  cat >"$OUT_DIR/bridge-docker-compose.yml" <<EOF
version: "3"

services:
EOF

  for chain in "${chain_name[@]}"; do
    local oracle_file="$OUT_DIR/$chain-bridge-oracle.json"

    bridge_contract_address=$(jq --arg chain_name "$chain" -r '.[]|select(.chain_name==$chain_name).bridge_contract_address' "$bridge_contract_file")

    while read -r bridge_index bridge_name; do
      cat >>"$OUT_DIR/bridge-docker-compose.yml" <<EOF
    fx-$chain-bridge-$bridge_name:
      container_name: fx-$chain-bridge-$bridge_name
      image: $bridge_image
      hostname: fx-$chain-bridge
      command: --chain-name="$chain" --external-jsonrpc="$external_json_rpc_url" --external-key="$TEST_MNEMONIC" --external-index="$bridge_index" --fx-bridge-addr="$bridge_contract_address" --fx-gas-price=4000000000000FX --fx-grpc="http://fxcore:9090" --fx-key="$TEST_MNEMONIC" --fx-index="$bridge_index"
      networks:
        - $docker_network
EOF

    done < <(jq -r '.[] | "\(.bridge_index) \(.bridge_name)"' "$oracle_file")
  done

  cat >>"$OUT_DIR/bridge-docker-compose.yml" <<EOF
networks:
  $docker_network:
    external: true
EOF
}

function send_to_fx() {
  local bridge_contract_file="${OUT_DIR}/bridge_contract.json"
  destination_address=$(add_key "destination" "99" | jq -r ".address")

  while read -r bridge_info; do
    bridge_contract_address=$(echo "$bridge_info" | jq -r '.bridge_contract_address')

    addresses=$(echo "$bridge_info" | jq -r '.bridge_token[].address')
    while IFS= read -r address; do
      LOCAL_PORT=8535 "$PROJECT_DIR"/tests/scripts/contract.sh send_to_fx 0 "$bridge_contract_address" "$address" "$(to_18 "111")" "$destination_address" ""
      LOCAL_PORT=8535 "$PROJECT_DIR"/tests/scripts/contract.sh send_to_fx 0 "$bridge_contract_address" "$address" "$(to_18 "111")" "$destination_address" "erc20"
      LOCAL_PORT=8535 "$PROJECT_DIR"/tests/scripts/contract.sh send_to_fx 0 "$bridge_contract_address" "$address" "$(to_18 "111")" "$destination_address" "ibc/0/px"
    done <<<"$addresses"

  done < <(jq -c '.[]' "$bridge_contract_file")
}

function request_batch() {
  local chain_name=("$@")
  for chain in "${chain_name[@]}"; do
    length=$(cosmos_query "$chain" batch-fees | jq '.batch_fees | length')
    if [ "$length" -eq 0 ]; then
      continue
    fi

    while read -r token_contract; do
      denom=$(cosmos_query "$chain" denom "$token_contract" | jq -r '.denom')
      cosmos_tx "$chain" build-batch "$denom" "1" "1" "$(show_address "$FROM" -e)" --from "$FROM"
    done < <(jq -r '.[] | "\(.token_contract)"' "$(cosmos_query "$chain" batch-fees)")
  done
}

function run_test() {
  "$PROJECT_DIR"/tests/scripts/fxcore.sh init

  "$PROJECT_DIR"/tests/scripts/fxcore.sh start

  "$PROJECT_DIR"/tests/scripts/contract.sh stop >/dev/null

  LOCAL_PORT=8535 "$PROJECT_DIR"/tests/scripts/contract.sh start

  LOCAL_PORT=8535 "$PROJECT_DIR"/tests/scripts/contract.sh deploy_bridge_contract

  create_oracles "eth" "bsc"

  update_crosschain_oracles "eth" "bsc"

  create_oracle_bridger "eth" "bsc"

  LOCAL_PORT=8535 "$PROJECT_DIR"/tests/scripts/contract.sh init_bridge_contract

  LOCAL_PORT=8535 setup_bridge_server "eth" "bsc"

  docker-compose -f "$OUT_DIR/bridge-docker-compose.yml" up -d

  send_to_fx

  return
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
