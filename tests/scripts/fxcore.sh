#!/usr/bin/env bash

set -eo pipefail

PROJECT_DIR="${PROJECT_DIR:-"$(git rev-parse --show-toplevel)"}"
export PROJECT_DIR
export OUT_DIR="${PROJECT_DIR}/out"

readonly docker_image="ghcr.io/functionx/fx-core:4.0.0-rc1"

export NODE_HOME="$OUT_DIR/.$CHAIN_NAME"

function init() {
  [[ -d "$NODE_HOME" ]] && rm -r "$NODE_HOME"
  gen_cosmos_genesis
  cat >"$OUT_DIR/$CHAIN_NAME.json" <<EOF
{
  "chain_id": "$CHAIN_ID",
  "node_rpc": "$NODE_RPC",
  "node_home": "$NODE_HOME",
  "mint_denom": "$MINT_DENOM",
  "staking_denom": "$STAKING_DENOM",
  "gas_prices": "$GAS_PRICES",
  "bech32_prefix": "$BECH32_PREFIX"
}
EOF
}

function start() {
  if docker stats --no-stream; then
    docker run -d --name "$CHAIN_NAME" --network bridge -v "${NODE_HOME}:/root/.$CHAIN_NAME" \
      -p "0.0.0.0:9090:9090" -p "0.0.0.0:26657:26657" -p "0.0.0.0:1317:1317" \
      "$docker_image" start
  else
    $DAEMON config config.toml rpc.laddr "tcp://0.0.0.0:26657" --home "$NODE_HOME"
    nohup "$DAEMON" start --home "$NODE_HOME" >"$NODE_HOME/$CHAIN_NAME.log" &
  fi
  node_catching_up "$NODE_RPC"

}

function stop() {
  if docker stats --no-stream; then
    docker stop "$CHAIN_NAME"
    docker rm "$CHAIN_NAME"
  else
    pkill -f "$DAEMON"
  fi
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"
