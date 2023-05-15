#!/usr/bin/env bash

set -eo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

readonly rpc_port="16657"
readonly grpc_port="8090"
readonly rest_port="2317"

export DAEMON="pundixd"
export CHAIN_ID="PUNDIX"
export CHAIN_NAME="pundix"
export BECH32_PREFIX="px"
export NODE_HOME="$OUT_DIR/.$CHAIN_NAME"

[[ -z "$PURSE_ADDRESS" ]] && echo "PURSE_ADDRESS is not set"
export LOCAL_MINT_DENOM="bsc$PURSE_ADDRESS"
[[ -z "$PUNDIX_ADDRESS" ]] && echo "PUNDIX_ADDRESS is not set"
LOCAL_STAKING_BOND_DENOM="$(convert_ibc_denom "transfer/$IBC_CHANNEL/eth$PUNDIX_ADDRESS")"
export LOCAL_STAKING_BOND_DENOM

export STAKING_DENOM="$LOCAL_STAKING_BOND_DENOM"
export MINT_DENOM="$LOCAL_MINT_DENOM"

export DOCKER_IMAGE=${DOCKER_IMAGE:-"ghcr.io/pundix/pundix:0.2.3"}
readonly docker_env="-e LOCAL_MINT_DENOM=$LOCAL_MINT_DENOM -e LOCAL_STAKING_BOND_DENOM=$LOCAL_STAKING_BOND_DENOM"
export DAEMON="docker run --rm -i --network $DOCKER_NETWORK $docker_env -v $NODE_HOME:$NODE_HOME $DOCKER_IMAGE"

export NODE_RPC="http://$CHAIN_NAME:$rpc_port"
export NODE_GRPC="$CHAIN_NAME:$grpc_port"
export REST_RPC="http://$CHAIN_NAME:$rest_port"

GAS_PRICES="$(echo "2*10^12" | bc)$STAKING_DENOM"
export GAS_PRICES

function init() {
  [[ -d "$NODE_HOME" ]] && docker_stop && rm -r "$NODE_HOME"
  gen_cosmos_genesis

  docker_run "--rm $docker_env" add-genesis-account px1a53udazy8ayufvy0s434pfwjcedzqv34hargq6 "$(to_18 "10^5")${MINT_DENOM}"
  docker_run "--rm" config config.toml rpc.laddr "tcp://0.0.0.0:$rpc_port"
  docker_run "--rm" config app.toml grpc.address "0.0.0.0:$grpc_port"
  docker_run "--rm" config app.toml api.address "tcp://0.0.0.0:$rest_port"

  cat >"$OUT_DIR/$CHAIN_NAME.json" <<EOF
{
  "chain_id": "$CHAIN_ID",
  "node_rpc": "$NODE_RPC",
  "node_grpc": "$NODE_GRPC",
  "rest_rpc": "$REST_RPC",
  "node_home": "$NODE_HOME",
  "mint_denom": "$MINT_DENOM",
  "staking_denom": "$STAKING_DENOM",
  "gas_prices": "$GAS_PRICES",
  "bech32_prefix": "$BECH32_PREFIX"
}
EOF
}

function start() {
  docker_run "-d $docker_env -p 0.0.0.0:$grpc_port:$grpc_port -p 0.0.0.0:$rpc_port:$rpc_port -p 0.0.0.0:$rest_port:$rest_port" start
  node_catching_up "http://127.0.0.1:$rpc_port"
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
