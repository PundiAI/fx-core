#!/usr/bin/env bash

set -eo pipefail

readonly docker_images="ghcr.io/pundix/pundix:0.2.3"
readonly rpc_port="16657"

export DAEMON="pundixd"
export CHAIN_ID="PUNDIX"
export CHAIN_NAME="pundix"
export NODE_HOME="$OUT_DIR/.$CHAIN_NAME"
export LOCAL_MINT_DENOM="bsc$PURSE_ADDRESS"
LOCAL_STAKING_BOND_DENOM="$($DAEMON query fx-ibc-transfer denom-convert "transfer/$IBC_CHANNEL/eth$PUNDIX_ADDRESS" --home "$NODE_HOME" | jq -r .IBCDenom)"
export LOCAL_STAKING_BOND_DENOM
export STAKING_DENOM="$LOCAL_STAKING_BOND_DENOM"

export NODE_RPC="http://127.0.0.1:$rpc_port"
GAS_PRICES="$(echo "2*10^12" | bc)$STAKING_DENOM"
export GAS_PRICES
export BECH32_PREFIX="px"

function start() {
  [[ -d "$NODE_HOME" ]] && rm -r "$NODE_HOME"
  gen_cosmos_genesis

  $DAEMON add-genesis-account px1a53udazy8ayufvy0s434pfwjcedzqv34hargq6 "$(to_18 "10^5")${MINT_DENOM}" --home "$NODE_HOME"

  if docker stats --no-stream; then
    docker run -d --name "$CHAIN_NAME" --network bridge -v "${NODE_HOME}:/root/.$CHAIN_NAME" \
      -e LOCAL_MINT_DENOM="$LOCAL_MINT_DENOM" \
      -e LOCAL_STAKING_BOND_DENOM="$LOCAL_STAKING_BOND_DENOM" \
      -p "0.0.0.0:8090:9090" -p "0.0.0.0:$rpc_port:26657" -p "0.0.0.0:2317:1317" \
      "$docker_images" start
  else
    $DAEMON config config.toml rpc.laddr "tcp://0.0.0.0:$rpc_port" --home "$NODE_HOME"
    $DAEMON config app.toml api.address "tcp://0.0.0.0:2317" --home "$NODE_HOME"
    $DAEMON config app.toml grpc.address "0.0.0.0:8090" --home "$NODE_HOME"
    $DAEMON config config.toml p2p.laddr "tcp://0.0.0.0:36656" --home "$NODE_HOME"
    nohup $DAEMON start --home "$NODE_HOME" >"$NODE_HOME/$CHAIN_NAME.log" 2>&1 &
  fi
  node_catching_up "$NODE_RPC"

  cat >"$OUT_DIR/$CHAIN_NAME.json" <<EOF
{
  "chain_id": "$CHAIN_ID",
  "rpc_port": "$rpc_port",
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
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"
