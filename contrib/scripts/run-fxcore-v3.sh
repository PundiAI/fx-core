#!/usr/bin/env bash

set -o errexit -o nounset

NODE_HOME=$(realpath ./build/.fxcore)
echo "NODE_HOME = ${NODE_HOME}"
BINARY=$NODE_HOME/cosmovisor/genesis/bin/fxcored
echo "BINARY = ${BINARY}"
CHAINID=fxcore

USER_MNEMONIC="test test test test test test test test test test test junk"

if ! test -f "./build/bin/fxcored3"; then
  echo "fxcored v3 does not exist, please build it to ./build/bin/fxcored3 first"
  exit
fi


rm -rf ./build/.fxcore

mkdir -p "$NODE_HOME"/cosmovisor/genesis/bin
cp ./build/bin/fxcored3 "$NODE_HOME"/cosmovisor/genesis/bin/fxcored
$BINARY init upgrader --chain-id $CHAINID --home "$NODE_HOME"


if ! test -f "./build/bin/fxcored4"; then
  echo "fxcored v3 does not exist, please build it to ./build/bin/fxcored4 first"
  exit
fi

mkdir -p "$NODE_HOME"/cosmovisor/upgrades/fxv4/bin
cp ./build/bin/fxcored4 "$NODE_HOME"/cosmovisor/upgrades/fxv4/bin/fxcored

GOPATH=$(go env GOPATH)

export DAEMON_NAME=fxcored
export DAEMON_HOME="$NODE_HOME"
COSMOVISOR=$GOPATH/bin/cosmovisor

$BINARY config chain-id $CHAINID --home "$NODE_HOME"
$BINARY config keyring-backend test --home "$NODE_HOME"

tmp=$(mktemp)

# update genesis total supply
jq '.app_state.bank.supply[0].amount = "388604525462891000000000000"' "$NODE_HOME"/config/genesis.json > "$tmp" && mv "$tmp" "$NODE_HOME"/config/genesis.json
# update gov voting period
jq '.app_state.gov.voting_params.voting_period = "15s"' "$NODE_HOME"/config/genesis.json > "$tmp" && mv "$tmp" "$NODE_HOME"/config/genesis.json

echo "$USER_MNEMONIC" | $BINARY --home "$NODE_HOME" keys add fx1 --recover
$BINARY add-genesis-account fx1 10004000000000000000000000FX --home "$NODE_HOME"

$BINARY gentx fx1 100000000000000000000FX --home $NODE_HOME --chain-id "$CHAINID"
$BINARY collect-gentxs --home "$NODE_HOME"
$BINARY config app.toml minimum-gas-prices "4000000000000FX"

$COSMOVISOR run start --home "$NODE_HOME" --x-crisis-skip-assert-invariants

