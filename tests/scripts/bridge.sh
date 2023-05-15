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
  "${PROJECT_DIR}/tests/scripts/tx-proposal.sh" "$@"
}

function proposal_register_coin() {
  local proposal_file="$1"
  min_deposit=$(proposal query_min_deposit)
  cosmos_tx gov submit-proposal register-coin "$proposal_file" --title="Register Coin" --description="Register Coin" --deposit="$min_deposit" --from "$FROM"
  proposal vote yes
}

function create_oracles() {
  local chain_name=("$@")
  local index=${bridger_start_index}

  for chain in "${chain_name[@]}"; do
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
    [ ! -f "$oracle_file" ] && continue
    oracles=()
    while read -r oracle_address; do
      oracles+=("$oracle_address")
    done < <(jq -r '.[] | "\(.oracle_address)"' "$oracle_file")

    oracles_list=$(
      IFS=,
      echo "${oracles[*]}"
    )
    min_deposit=$(proposal query_min_deposit)

    if [[ "$(cosmos_version | grep "v3")" != "" ]]; then
      cosmos_tx crosschain update-crosschain-oracles "$chain" "$min_deposit" --oracles "$oracles_list" --title="Update $chain chain oracles" --desc="oracles description" --from "$FROM"
    else
      cosmos_tx "$chain" update-crosschain-oracles "$oracles_list" --deposit="$min_deposit" --title="Update $chain chain oracles" --description="oracles description" --from "$FROM"
    fi
    proposal vote yes
  done
}

function create_oracle_bridger() {
  local chain_name=("$@")
  for chain in "${chain_name[@]}"; do
    local oracle_file="$OUT_DIR/$chain-bridge-oracle.json"
    [ ! -f "$oracle_file" ] && continue

    min_deposit=$(proposal query_min_deposit)
    validator_address=$(show_address "$FROM" -a --bech val)

    while read -r oracle_name oracle_address oracle_index bridge_name bridge_address bridge_index external_address; do
      add_key "$oracle_name" "$oracle_index"
      add_key "$bridge_name" "$bridge_index"

      cosmos_tx bank send "$FROM" "$oracle_address" "$min_deposit" --from "$FROM"
      cosmos_transfer "$oracle_name" 100
      cosmos_transfer "$bridge_name" 500

      if [[ "$(cosmos_version | grep "v3")" != "" ]]; then
        cosmos_tx crosschain create-oracle-bridger "$chain" "$validator_address" "$bridge_address" "$external_address" "$min_deposit" --from "$oracle_name"
      else
        cosmos_tx "$chain" create-oracle-bridger "$validator_address" "$bridge_address" "$external_address" "$min_deposit" --from "$oracle_name"
      fi
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

function register_coin() {
  while read -r chain_name address decimals symbol target_ibc name; do
    if [[ "$symbol" == "FX" ]]; then
      cat >"$OUT_DIR/coin.json" <<EOF
{
      "description": "The native staking token of the Function X",
      "denom_units": [
        {
          "denom": "FX",
          "exponent": 0,
          "aliases": []
        }
      ],
      "base": "FX",
      "display": "FX",
      "name": "Function X",
      "symbol": "FX"
}
EOF
    else
      DENOM="${chain_name}${address}"
      if [ "$target_ibc" != "null" ]; then
        DENOM="$(convert_ibc_denom "${target_ibc}/${chain_name}${address}")"
      fi
      cat >"$OUT_DIR/coin.json" <<EOF
{
      "description": "The cross chain token of the Function X",
      "denom_units": [
        {
          "denom": "$DENOM",
          "exponent": 0,
          "aliases": []
        },
        {
          "denom": "$symbol",
          "exponent": $decimals,
          "aliases": []
        }
      ],
      "base": "$DENOM",
      "display": "$DENOM",
      "name": "$name",
      "symbol": "$symbol"
}
EOF
    fi
    proposal_register_coin "$OUT_DIR/coin.json"
  done < <(jq -r '.bridge_token_list.one_to_one[] | "\(.chain_name) \(.address) \(.decimals) \(.symbol) \(.target_ibc) \(.name)"' "$bridge_contract_file")

  while read -r chain_list base_denom symbol decimals name; do
    aliases=()

    while read -r chain_name address target_ibc; do
      denom="${chain_name}${address}"
      denom="${chain_name}${address}"
      if [ "$target_ibc" != "null" ]; then
        denom="$(convert_ibc_denom "${target_ibc}/${chain_name}${address}")"
      fi
      aliases+=("\"$denom\"")
    done < <(echo "$chain_list" | jq -r '.[] | "\(.chain_name) \(.address) \(.target_ibc)"')

    IFS=,
    alias_str="${aliases[*]}"

    cat >"$OUT_DIR/coin.json" <<EOF
 {
      "description": "The cross chain token of the Function X",
      "denom_units": [
        {
          "denom": "$base_denom",
          "exponent": 0,
          "aliases": [$alias_str]
        },
        {
          "denom": "$symbol",
          "exponent": "$decimals",
          "aliases": []
        }
      ],
      "base": "$base_denom",
      "display": "$base_denom",
      "name": "$name",
      "symbol": "$symbol"
    },
EOF
    proposal_register_coin "$OUT_DIR/coin.json"
  done < <(jq -r '.bridge_token_list.one_to_many[] | "\(.chain_list) \(.base_denom) \(.symbol) \(.decimals) \(.name)"' "$bridge_contract_file")
}

function request_batch() {
  local chain_name=("$@")
  for chain in "${chain_name[@]}"; do
    if [[ "$(cosmos_version | grep "v3")" != "" ]]; then
      length=$(cosmos_query crosschain batch-fees "$chain" | jq '.batch_fees | length')
    else
      length=$(cosmos_query "$chain" batch-fees | jq '.batch_fees | length')
    fi

    if [ "$length" -eq 0 ]; then
      continue
    fi

    while read -r token_contract; do

      if [[ "$(cosmos_version | grep "v3")" != "" ]]; then
        denom=$(cosmos_query crosschain denom "$chain" "$token_contract" | jq -r '.denom')
        cosmos_tx crosschain build-batch "$chain" "$denom" "1" "$(show_address "$FROM" -e)" --from "eth-bridger-0"
      else
        denom=$(cosmos_query "$chain" denom "$token_contract" | jq -r '.denom')
        cosmos_tx "$chain" build-batch "$denom" "1" "1" "$(show_address "$FROM" -e)" --from "$FROM"
      fi

    done < <(cosmos_query crosschain batch-fees "$chain" | jq -r '.batch_fees[] | "\(.token_contract)"')
  done
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
