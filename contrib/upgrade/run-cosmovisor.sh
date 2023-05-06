#!/usr/bin/env bash

set -o errexit -o nounset

export CUR_VERSION=${CUR_VERSION:-"v3"}
export NEXT_VERSION=${NEXT_VERSION:-"v4"}

export NODE_HOME=${NODE_HOME:-"./out/.fxcore"}
echo "NODE_HOME = ${NODE_HOME}"
export BINARY=${BINARY:-"$NODE_HOME/cosmovisor/genesis/bin/fxcored"}
echo "BINARY = ${BINARY}"

export CHAIN_ID=${CHAIN_ID:-"fxcore"}

if [ -d "$NODE_HOME" ]; then
  echo "node home $NODE_HOME already exists"
  read -rp "Are you sure you want to delete all the data and start over? [y/N] " input
  [[ "$input" != "y" && "$input" != "Y" ]] && exit 1
  rm -r "$NODE_HOME"
fi

if ! test -f "fxcored-${CUR_VERSION}"; then
  echo "fxcored ${CUR_VERSION} does not exist, please build it to fxcored-${CUR_VERSION} first" && exit 1
fi

mkdir -p "$NODE_HOME"/cosmovisor/genesis/bin
cp "fxcored-${CUR_VERSION}" "$NODE_HOME"/cosmovisor/genesis/bin/fxcored

$BINARY init upgrader --chain-id "$CHAIN_ID" --home "$NODE_HOME"
$BINARY config chain-id "$CHAIN_ID" --home "$NODE_HOME"
$BINARY config keyring-backend test --home "$NODE_HOME"

if ! test -f "fxcored-${NEXT_VERSION}"; then
  echo "fxcored ${NEXT_VERSION} does not exist, please build it to fxcored-${NEXT_VERSION} first" && exit 1
fi

mkdir -p "$NODE_HOME/cosmovisor/upgrades/fx${NEXT_VERSION}/bin"
cp "fxcored-${NEXT_VERSION}" "$NODE_HOME/cosmovisor/upgrades/fx${NEXT_VERSION}/bin/fxcored"

readonly genesis_tmp="$NODE_HOME"/config/genesis.json.tmp
# update genesis total supply
jq '.app_state.bank.supply[0].amount = "388604525462891000000000000"' "$NODE_HOME"/config/genesis.json >"$genesis_tmp" &&
  mv "$genesis_tmp" "$NODE_HOME"/config/genesis.json
# update gov voting period
jq '.app_state.gov.voting_params.voting_period = "15s"' "$NODE_HOME"/config/genesis.json >"$genesis_tmp" &&
  mv "$genesis_tmp" "$NODE_HOME"/config/genesis.json

echo "$TEST_MNEMONIC" | $BINARY --home "$NODE_HOME" keys add fx1 --recover
$BINARY add-genesis-account fx1 10004000000000000000000000FX --home "$NODE_HOME"

$BINARY gentx fx1 100000000000000000000FX --home "$NODE_HOME" --chain-id "$CHAIN_ID"
$BINARY collect-gentxs --home "$NODE_HOME"
$BINARY config app.toml minimum-gas-prices "4000000000000FX"

export DAEMON_NAME=fxcored
export DAEMON_HOME="$NODE_HOME"
cosmovisor run start --home "$NODE_HOME" --x-crisis-skip-assert-invariants
