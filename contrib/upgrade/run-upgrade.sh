#!/usr/bin/env bash

set -o errexit -o nounset

export UPGRADE_HEIGHT_INTERVAL=${1:-10}
export PROPOSAL_ID=${2:-""}

export NEXT_VERSION=${NEXT_VERSION:-"v4"}

export NODE_HOME=${NODE_HOME:-"./out/.fxcore"}
echo "NODE_HOME = ${NODE_HOME}"
export BINARY=${BINARY:-"$NODE_HOME/cosmovisor/genesis/bin/fxcored"}
echo "BINARY = ${BINARY}"

export CHAIN_ID=${CHAIN_ID:-"fxcore"}

if ! test -f "$BINARY"; then
  echo "Binary not found at $BINARY"
fi

$BINARY config chain-id "$CHAIN_ID" --home "$NODE_HOME"
$BINARY config output json --home "$NODE_HOME"
$BINARY config keyring-backend test --home "$NODE_HOME"
$BINARY config node tcp://localhost:26657 --home "$NODE_HOME"
$BINARY config broadcast-mode block --home "$NODE_HOME"
$BINARY config --home "$NODE_HOME"

while true; do
  sync_state=$("$BINARY" status --home "$NODE_HOME" | jq -r '.SyncInfo.catching_up')
  if [ "$sync_state" != "false" ]; then
    echo "Node is syncing..." && continue
  fi
  break
done

if [ -z "$($BINARY keys show fx1 --home "$NODE_HOME")" ]; then
  echo "$TEST_MNEMONIC" | $BINARY --home "$NODE_HOME" keys add fx1 --recover
fi

upgrade_height=$($BINARY status --home "$NODE_HOME" | jq -r '.SyncInfo.latest_block_height|tonumber + '"${UPGRADE_HEIGHT_INTERVAL}"'')
readonly upgrade_height
printf "\n"
echo "Upgrade Height = ${upgrade_height}"
printf "Submitting proposal... \n"
$BINARY tx gov submit-proposal software-upgrade "fx$NEXT_VERSION" \
  --title "fx$NEXT_VERSION" \
  --deposit "$($BINARY query gov params --home "$NODE_HOME" | jq -r '.deposit_params.min_deposit[0].amount')FX" \
  --upgrade-height "${upgrade_height}" \
  --upgrade-info "upgrade to fx$NEXT_VERSION" \
  --description "upgrade to fx$NEXT_VERSION" \
  --gas auto --gas-prices 4000000000000FX --gas-adjustment=1.3 \
  --from fx1 \
  --home "${NODE_HOME}" \
  --yes
printf "Done \n"

printf "Casting vote... \n"

if [ -z "$PROPOSAL_ID" ]; then
  PROPOSAL_ID=$($BINARY query gov proposals --status=voting_period --home "$NODE_HOME" | jq -r '.proposals[0].PROPOSAL_ID')
fi

echo "Vote ProposalID  =  ${PROPOSAL_ID}"

$BINARY tx gov vote "${PROPOSAL_ID}" yes \
  --gas auto --gas-prices 4000000000000FX --gas-adjustment=1.3 \
  --from fx1 \
  --home "${NODE_HOME}" \
  --yes

printf "Done \n"
