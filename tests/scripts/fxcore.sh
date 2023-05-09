#!/usr/bin/env bash

project_dir="$(git rev-parse --show-toplevel)"
readonly project_dir
readonly out_dir="${project_dir}/out"

readonly docker_image="ghcr.io/functionx/fx-core:4.0.0-rc1"
readonly rpc_port="26657"

export DAEMON="fxcored"
export CHAIN_ID="fxcore"
export CHAIN_NAME="fxcore"
export NODE_HOME="$out_dir/.$CHAIN_NAME"
export STAKING_DENOM="FX"
export MINT_DENOM="FX"
export NODE_RPC="http://127.0.0.1:$rpc_port"
GAS_PRICES="$(echo "4*10^12" | bc)$STAKING_DENOM"
export GAS_PRICES
export BECH32_PREFIX="fx"

function start() {
  [[ -d "$NODE_HOME" ]] && rm -r "$NODE_HOME"
  gen_cosmos_genesis
  
  $DAEMON config config.toml rpc.laddr "tcp://0.0.0.0:$rpc_port" --home "$NODE_HOME"
  
  if docker stats --no-stream; then
    docker run -d --name "$CHAIN_NAME" --network bridge -v "${NODE_HOME}:/root/.$CHAIN_NAME" \
      -p "0.0.0.0:9090:9090" -p "0.0.0.0:$rpc_port:26657" -p "0.0.0.0:1317:1317" "$docker_image" start
  else
    nohup $DAEMON start --home "$NODE_HOME" >"$NODE_HOME/$CHAIN_NAME.log" &
  fi
  node_catching_up "$NODE_RPC"

  cat >"$out_dir/$CHAIN_NAME.json" <<EOF
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

function stop() {
  if docker stats --no-stream; then
    docker stop "$CHAIN_NAME"
    docker rm "$CHAIN_NAME"
  else
    pkill -f "$DAEMON"
  fi
}

# shellcheck source=/dev/null
. "${project_dir}/tests/scripts/setup-env.sh"
