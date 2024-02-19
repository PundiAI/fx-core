#!/usr/bin/env bash

set -eo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

export CHAIN_NAME="fxcore"
export NODE_HOME="$OUT_DIR/.$CHAIN_NAME"

export DOCKER_IMAGE=${DOCKER_IMAGE:-"ghcr.io/functionx/fx-core:latest"}
export DAEMON="docker run --rm -i --network $DOCKER_NETWORK -v $NODE_HOME:$NODE_HOME $DOCKER_IMAGE"

export NODE_RPC="http://$CHAIN_NAME:26657"
export NODE_GRPC="$CHAIN_NAME:9090"
export REST_RPC="http://$CHAIN_NAME:1317"

GAS_PRICES="$(echo "4*10^12" | bc)$STAKING_DENOM"
export GAS_PRICES

function init() {
  [[ -d "$NODE_HOME" ]] && docker_stop && rm -r "$NODE_HOME"
  gen_cosmos_genesis

  cat >"$OUT_DIR/$CHAIN_NAME.json" <<EOF
{
  "chain_id": "$CHAIN_ID",
  "node_rpc": "$NODE_RPC",
  "node_grpc": "http://${NODE_GRPC}",
  "rest_rpc": "$REST_RPC",
  "node_home": "$NODE_HOME",
  "mint_denom": "$MINT_DENOM",
  "staking_denom": "$STAKING_DENOM",
  "gas_prices": "$GAS_PRICES",
  "bech32_prefix": "$BECH32_PREFIX"
}
EOF
  sleep 2
}

function start() {
  docker_run "-d -p 0.0.0.0:9090:9090 -p 0.0.0.0:26657:26657 -p 0.0.0.0:1317:1317 -p 0.0.0.0:8545:8545" start
  node_catching_up "http://127.0.0.1:26657"
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
