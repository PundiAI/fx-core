#!/usr/bin/env bash

set -eo pipefail

export DEBUG=${DEBUG:-"false"}
[[ "$DEBUG" == "true" ]] && set -x

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
  local json_file=$1
  shift
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

function cosmos_tx() {
  $DAEMON tx "$@" "${tx_flags_ary[@]}" || (echo "failed: $DAEMON tx $*" && exit 1)
}

## ARGS: <to> <amount> [<denom>]
function cosmos_transfer() {
  local to=$1 amount=$2 denom=${3:-$STAKING_DENOM}
  node_catching_up "$NODE_RPC"
  cosmos_tx bank send "$FROM" "$($DAEMON keys show "$to" --home "$NODE_HOME" -a)" "$(to_18 "$amount")$denom" --from "$FROM"
}

function cosmos_query() {
  $DAEMON query "$@" "${query_flags_ary[@]}" || (echo "failed: $DAEMON query $*" && exit 1)
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

  $DAEMON config config.toml consensus.timeout_commit 1s --home "$NODE_HOME"
  $DAEMON config config.toml rpc.pprof_laddr "" --home "$NODE_HOME"
  $DAEMON config config.toml rpc.laddr "tcp://0.0.0.0:26657" --home "$NODE_HOME"

  $DAEMON config chain-id "$CHAIN_ID" --home "$NODE_HOME"
  $DAEMON config keyring-backend "$KEYRING_BACKEND" --home "$NODE_HOME"
  $DAEMON config output "$OUTPUT" --home "$NODE_HOME"
  $DAEMON config broadcast-mode "$BROADCAST_MODE" --home "$NODE_HOME"
  $DAEMON config node "$NODE_RPC" --home "$NODE_HOME"

  echo "$TEST_MNEMONIC" | $DAEMON keys add "$FROM" --recover --home "$NODE_HOME"
  genesis_amount="$(to_18 "10^5")${STAKING_DENOM}"
  [[ -n "$MINT_DENOM" && "$STAKING_DENOM" != "$MINT_DENOM" ]] && genesis_amount="$genesis_amount,$(to_18 "10^5")${MINT_DENOM}"
  $DAEMON add-genesis-account "$FROM" "$genesis_amount" --home "$NODE_HOME"

  set +e && supply="$($DAEMON validate-genesis --home "$NODE_HOME" 2>&1 | grep "expected .*$STAKING_DENOM" | cut -d " " -f 14)" && set -e
  if [ -n "$supply" ]; then
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
  while true; do
    sync_state=$(curl -s "$node_url/status" | jq -r '.result.sync_info.catching_up')
    if [ "$sync_state" != "false" ]; then
      sleep 1
      echo "Node is syncing..." && continue
    fi
    break
  done
}

function show_address() {
  local from=${1:-"$FROM"} opt=${2:-""}
  $DAEMON keys show "$from" "$opt" --home "$NODE_HOME"
}

function add_key() {
  local name=$1 index=$2
  echo "$TEST_MNEMONIC" | $DAEMON keys add "$name" --index "$index" --home "$NODE_HOME" --recover
}

function cosmos_grpc() {
  grpcurl -plaintext "$NODE_GRPC" "$@"
}

function cosmos_reset() {
  curl -s "$REST_RPC/$*" | jq -r '.result'
}

export DAEMON=${DAEMON:-"fxcored"}
export CHAIN_ID=${CHAIN_ID:-"fxcore"}
export CHAIN_NAME=${CHAIN_NAME:-"fxcore"}
export NODE_RPC=${NODE_RPC:-"http://127.0.0.1:26657"}
export REST_RPC=${REST_RPC:-"http://127.0.0.1:1317"}
export NODE_GRPC=${NODE_GRPC:-"127.0.0.1:9090"}
export NODE_HOME=${NODE_HOME:-"$HOME/.$CHAIN_NAME"}

export DENOM=${DENOM:-"FX"}
export DECIMALS=${DECIMALS:-"18"}
export STAKING_DENOM=${STAKING_DENOM:-"$DENOM"}
export MINT_DENOM=${MINT_DENOM:-"$STAKING_DENOM"}

export OUTPUT=${OUTPUT:-"json"}
export KEYRING_BACKEND=${KEYRING_BACKEND:-"test"}
export BROADCAST_MODE=${BROADCAST_MODE:-"block"}
export GAS_ADJUSTMENT=${GAS_ADJUSTMENT:-1.3}
export GAS_PRICES=${GAS_PRICES:-"$(echo "4*10^12" | bc)$STAKING_DENOM"}

export TEST_MNEMONIC=${TEST_MNEMONIC:-"test test test test test test test test test test test junk"}
export FROM=${FROM:-"test1"}

export QUERY_FLAGS=${QUERY_FLAGS:-"--node=$NODE_RPC --output=$OUTPUT"}
IFS=' ' read -r -a query_flags_ary <<<"$QUERY_FLAGS"

export TX_FLAGS=${TX_FLAGS:-"--keyring-backend=$KEYRING_BACKEND --gas-prices=$GAS_PRICES --gas=auto --gas-adjustment=$GAS_ADJUSTMENT --broadcast-mode=$BROADCAST_MODE --output=$OUTPUT --node=$NODE_RPC --chain-id=$CHAIN_ID --home=$NODE_HOME -y"}
IFS=' ' read -r -a tx_flags_ary <<<"$TX_FLAGS"

if [[ "$1" == "help" || "$#" -eq 0 ]]; then
  help && exit 0
fi

if [[ "$#" -gt 0 && "$(type -t "$1")" != "function" ]]; then
  echo "invalid command: $1" && help && exit 1
fi

if ! "$@"; then
  echo "failed: $0" "$@" && exit 1
fi
