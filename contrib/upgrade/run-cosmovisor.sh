#!/usr/bin/env bash

set -eo pipefail

readonly init=${1:-"no"}
export CURRENT_VERSION=${2:-"$CURRENT_VERSION"}
[[ -z "$CURRENT_VERSION" ]] && echo "CURRENT_VERSION is required" && exit 1
export UPGRADE_NAME=${3:-"$UPGRADE_NAME"}

export CHAIN_ID=${CHAIN_ID:-"fxcore"}
export NODE_HOME=${NODE_HOME:-"$HOME/.fxcore"}
export BINARY=${BINARY:-"$NODE_HOME/cosmovisor/genesis/bin/fxcored"}

function build_binary() {
  local version=${1:-""}
  echo "build binary for version $version"
  [[ ! -d "/tmp/$CHAIN_ID-cache" ]] && mkdir -p "/tmp/$CHAIN_ID-cache"
  (
    cd "/tmp/$CHAIN_ID-cache" || exit 1
    if [[ ! -d "fx-core/.git" ]]; then
      git clone https://github.com/functionx/fx-core.git
    fi
    cd fx-core || exit 1
    git fetch --all
    git checkout "release/$version"
    make build
    mkdir -p "$NODE_HOME/cosmovisor/upgrades/${version}/bin"
    cp "./build/bin/fxcored" "$NODE_HOME/cosmovisor/upgrades/${version}/bin/fxcored"
  )
}

if [ -n "$UPGRADE_NAME" ]; then
  build_binary "$UPGRADE_NAME"
fi

build_binary "$CURRENT_VERSION"
mkdir -p "$(dirname "$BINARY")"
[[ -f "$BINARY" ]] && rm "$BINARY"
ln -s "$NODE_HOME/cosmovisor/upgrades/${CURRENT_VERSION}/bin/fxcored" "$BINARY"

if [ "$init" == "init" ]; then
  [[ -d "$NODE_HOME/data" ]] && rm -r "$NODE_HOME/data"
  [[ -d "$NODE_HOME/config" ]] && rm -r "$NODE_HOME/config"
  [[ -d "$NODE_HOME/keyring-test" ]] && rm -r "$NODE_HOME/keyring-test"

  $BINARY init upgrader --chain-id "$CHAIN_ID" --home "$NODE_HOME"
  $BINARY config chain-id "$CHAIN_ID" --home "$NODE_HOME"
  $BINARY config keyring-backend test --home "$NODE_HOME"

  readonly genesis_tmp="$NODE_HOME"/config/genesis.json.tmp
  # update genesis total supply
  jq '.app_state.bank.supply[0].amount = "388604525462891000000000000"' "$NODE_HOME"/config/genesis.json >"$genesis_tmp" &&
    mv "$genesis_tmp" "$NODE_HOME"/config/genesis.json
  # update gov voting period
  jq '.app_state.gov.voting_params.voting_period = "15s"' "$NODE_HOME"/config/genesis.json >"$genesis_tmp" &&
    mv "$genesis_tmp" "$NODE_HOME"/config/genesis.json

  echo "test test test test test test test test test test test junk" | $BINARY keys add fx1 --recover --home "$NODE_HOME"
  $BINARY add-genesis-account fx1 10004000000000000000000000FX --home "$NODE_HOME"

  $BINARY gentx fx1 100000000000000000000FX --chain-id "$CHAIN_ID" --home "$NODE_HOME"
  $BINARY collect-gentxs --home "$NODE_HOME"
fi

echo "start fxcore ..."
if docker stats --no-stream; then
  docker run -d --name fxcore \
    -p 0.0.0.0:26656:26656 -p 127.0.0.1:26657:26657 -p 127.0.0.1:1317:1317 -p 127.0.0.1:8545:8545 -p 127.0.0.1:8546:8546 \
    -v "$NODE_HOME":/root/.fxcore ghcr.io/functionx/fxcorevisor:latest run start --x-crisis-skip-assert-invariants
else
  export DAEMON_NAME=fxcored
  export DAEMON_HOME="$NODE_HOME"
  cosmovisor run start --x-crisis-skip-assert-invariants --home "$NODE_HOME"
fi
