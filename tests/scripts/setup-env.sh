#!/usr/bin/env bash

set -eo pipefail

export DEBUG=${DEBUG:-"false"}
[[ "$DEBUG" == "true" ]] && set -x

PROJECT_DIR="${PROJECT_DIR:-"$(git rev-parse --show-toplevel)"}"
export PROJECT_DIR
export OUT_DIR="${PROJECT_DIR}/out"

export DAEMON=${DAEMON:-"fxcored"}
export CHAIN_ID=${CHAIN_ID:-"fxcore"}
export CHAIN_NAME=${CHAIN_NAME:-"fxcore"}
export NODE_RPC=${NODE_RPC:-"http://127.0.0.1:26657"}
export REST_RPC=${REST_RPC:-"http://127.0.0.1:1317"}
export NODE_GRPC=${NODE_GRPC:-"127.0.0.1:9090"}
export NODE_HOME=${NODE_HOME:-"$HOME/.$CHAIN_NAME"}
export BECH32_PREFIX="fx"

export DENOM=${DENOM:-"FX"}
export DECIMALS=${DECIMALS:-"18"}
export STAKING_DENOM=${STAKING_DENOM:-"$DENOM"}
export MINT_DENOM=${MINT_DENOM:-"$STAKING_DENOM"}

export OUTPUT=${OUTPUT:-"json"}
export KEYRING_BACKEND=${KEYRING_BACKEND:-"test"}
export BROADCAST_MODE=${BROADCAST_MODE:-"block"}
export GAS_ADJUSTMENT=${GAS_ADJUSTMENT:-1.3}
GAS_PRICES="$(echo "4*10^12" | bc)$STAKING_DENOM"
export GAS_PRICES

export TEST_MNEMONIC=${TEST_MNEMONIC:-"test test test test test test test test test test test junk"}
export FROM=${FROM:-"test1"}

export DOCKER_NETWORK=${DOCKER_NETWORK:-"test-net"}

if [ -f "$PROJECT_DIR/tests/scripts/.env" ]; then
  # shellcheck source=/dev/null
  source "$PROJECT_DIR/tests/scripts/.env"
fi

mkdir -p "${OUT_DIR}"

function check_command() {
  commands=("$@")
  for cmd in "${commands[@]}"; do
    if ! command -v "$cmd" &>/dev/null; then
      echo "$cmd command not found, please install $cmd first" && exit 1
    fi
  done
}

function help() {
  printf "%s\n" "$0"
  printf "Usage: %s [<command>]\n\n" "$0"
  printf "Commands:\n\n"
  {
    local label
    items=$(echo "${BASH_SOURCE[*]:1}" | xargs grep "^##\ \(DESC\|ARGS\)\|^function\ ")
    IFS=$'\n'
    for item in $items; do
      [[ "$item" == *\#\#\ ARGS:* ]] && label=${item##*\#\#\ ARGS:\ } && continue
      [[ "$item" == *\#\#\ DESC:* ]] && label="${label:-"-"}#${item##*\#\#\ DESC:\ }" && continue
      [[ "$item" == *function* ]] && printf "    %s#%s\n" "$(echo "${item##*function\ }" | cut -d \( -f 1)" "$label"
      label=""
    done
  } | column -t -s '#'
  printf "\nEnvironment variables:\n\n"
  {
    items=$(echo "${BASH_SOURCE[*]:1}" | xargs grep "^export ")
    IFS=$'\n'
    for item in $items; do
      printf "    %s\n" "${item##*export }"
    done
  } | column -t -s '#'
}

## ARGS: <json_file> <jq_opt...>
function json_processor() {
  local json_file=$1 && shift
  local jq_opt=("$@")

  jq "${jq_opt[@]}" "$json_file" >"$json_file.tmp" &&
    mv "$json_file.tmp" "$json_file"
}

## ARGS: <Elapsed Percentage (0-100)> <Total length in chars>
function bar() {
  ((elapsed = $1 * $2 / 100))

  # Create the bar with spaces.
  printf -v prog "%${elapsed}s"
  printf -v total "%$(($2 - elapsed))s"

  printf '%s\r' "[${prog// /-}${total}]"
}

function echo_error() {
  echo -e "\\033[31m$*\\033[m"
}

## ARGS: <args...>
function cosmos_tx() {
  $DAEMON tx "$@" --gas-prices="$GAS_PRICES" --gas=auto --gas-adjustment="$GAS_ADJUSTMENT" --node="$NODE_RPC" --home="$NODE_HOME" -y
}

## ARGS: <to> <amount> [<denom>]
function cosmos_transfer() {
  local to=$1 amount=$2 denom=${3:-$STAKING_DENOM}
  to_address=$($DAEMON keys show "$to" --home "$NODE_HOME" -a)
  cosmos_tx bank send "$FROM" "$to_address" "$(to_18 "$amount")$denom" --from "$FROM"
}

## ARGS: <args...>
function cosmos_query() {
  $DAEMON query "$@" --node="$NODE_RPC" --home="$NODE_HOME"
}

function cosmos_version() {
  $DAEMON version
}

function to_18() {
  echo "$1 * 10^18" | bc
}

function from_18() {
  echo "$1 / 10^18" | bc
}

function gen_cosmos_genesis() {
  $DAEMON init --chain-id="$CHAIN_ID" local-test --home "$NODE_HOME"

  $DAEMON config app.toml grpc-web.enable false --home "$NODE_HOME"
  $DAEMON config app.toml api.enable true --home "$NODE_HOME"

  $DAEMON config app.toml json-rpc.address "0.0.0.0:8545" --home "$NODE_HOME"
  $DAEMON config app.toml json-rpc.ws-address "0.0.0.0:8546" --home "$NODE_HOME"

  $DAEMON config config.toml consensus.timeout_commit 1s --home "$NODE_HOME"
  $DAEMON config config.toml rpc.pprof_laddr "" --home "$NODE_HOME"
  $DAEMON config config.toml rpc.laddr "tcp://0.0.0.0:26657" --home "$NODE_HOME"

  $DAEMON config chain-id "$CHAIN_ID" --home "$NODE_HOME"
  $DAEMON config keyring-backend "$KEYRING_BACKEND" --home "$NODE_HOME"
  $DAEMON config output "$OUTPUT" --home "$NODE_HOME"
  $DAEMON config broadcast-mode "$BROADCAST_MODE" --home "$NODE_HOME"
  $DAEMON config node "$NODE_RPC" --home "$NODE_HOME"

  echo "$TEST_MNEMONIC" | $DAEMON keys add "$FROM" --recover --home "$NODE_HOME"
  genesis_amount="$(to_18 "10^6")${STAKING_DENOM}"
  [[ -n "$MINT_DENOM" && "$STAKING_DENOM" != "$MINT_DENOM" ]] && genesis_amount="$genesis_amount,$(to_18 "10^6")${MINT_DENOM}"
  $DAEMON add-genesis-account "$FROM" "$genesis_amount" --home "$NODE_HOME"

  set +e && supply="$($DAEMON validate-genesis --home "$NODE_HOME" 2>&1 | grep "expected .*$STAKING_DENOM" | cut -d " " -f 14)" && set -e
  if [[ -n "$supply" ]]; then
    json_processor "$NODE_HOME/config/genesis.json" ".app_state.bank.supply[0].amount = \"${supply%%"$STAKING_DENOM"}\""
  fi

  json_processor "$NODE_HOME/config/genesis.json" '.app_state.gov.voting_params.voting_period = "5s"'

  $DAEMON gentx "$FROM" "$(to_18 100)${STAKING_DENOM}" --chain-id="${CHAIN_ID}" \
    --moniker="test-validator" \
    --commission-max-change-rate="0.01" \
    --commission-max-rate="0.2" \
    --commission-rate="0.03" \
    --gas="20000000" \
    --gas-prices="" \
    --home "$NODE_HOME"
  $DAEMON collect-gentxs --home "$NODE_HOME"
}

function node_catching_up() {
  local node_url=${1:-"$REST_RPC"}
  local timeout=${2:-"10"}
  for i in $(seq "$timeout"); do
    sync_state=$(curl -s "$node_url/status" | jq -r '.result.sync_info.catching_up')
    if [ "$sync_state" != "false" ]; then
      sleep 1 && echo "Node is syncing... $i" && continue
    fi
    return 0
  done
  echo "Timeout: Node is not catching up"
  return 1
}

function show_address() {
  local from=${1:-"$FROM"} && shift
  local flags=("$@")
  $DAEMON keys show "$from" "${flags[@]}" --home "$NODE_HOME"
}

function validators_list() {
  cosmos_query staking validators --home "$NODE_HOME" --output json | jq -r ['.validators[] | select(.status == "BOND_STATUS_BONDED") | .operator_address']
}

function add_key() {
  local name=$1 index=$2
  $DAEMON keys delete "$name" --home "$NODE_HOME" -y >/dev/null 2>&1
  echo "$TEST_MNEMONIC" | $DAEMON keys add "$name" --index "$index" --home "$NODE_HOME" --recover
}

function export_key() {
  local name=$1 out=$2
  shift 2
  $DAEMON keys export "$name" "$@" --home "$NODE_HOME" >"$out"
}

function cosmos_grpc() {
  grpcurl -plaintext "$NODE_GRPC" "$@"
}

function cosmos_reset() {
  curl -s "$REST_RPC/$*" | jq -r '.result'
}

function sha256sum() {
  echo -n "$1" | shasum -a 256 | awk '{print $1}'
}

function convert_ibc_denom() {
  echo "ibc/$(sha256sum "$1" | tr '[:lower:]' '[:upper:]')"
}

function docker_run() {
  local opts=$1 && shift
  local args=("$@")
  local name=$CHAIN_NAME
  [[ "$opts" == "--rm" ]] && name="$name-tmp"
  IFS=' ' read -r -a opts_ary <<<"$opts"
  docker run "${opts_ary[@]}" --name "$name" --network "$DOCKER_NETWORK" -v "$NODE_HOME:$NODE_HOME" \
    "$DOCKER_IMAGE" "${args[@]}" --home "$NODE_HOME"
}

function docker_stop() {
  local container=${1:-"$CHAIN_NAME"}
  if docker ps -a | grep "$container" >/dev/null; then
    docker stop "$container" && docker rm "$container" && sleep 1
  fi
}

function create_docker_network() {
  local network=${1:-"$DOCKER_NETWORK"}
  [[ "$network" != "$DOCKER_NETWORK" ]] && export DOCKER_NETWORK=$network

  # check docker is running
  if docker stats --no-stream >/dev/null; then
    # check docker network exists
    if ! docker network ls -f "name=$network" | grep "$network" >/dev/null; then
      docker network create "$network"
    fi
  fi
}
