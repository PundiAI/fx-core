#!/usr/bin/env bash

set -eo pipefail

export UPGRADE_HEIGHT_INTERVAL=${1:-10}
export UPGRADE_NAME=${2:-"$UPGRADE_NAME"}
[[ -z "$UPGRADE_NAME" ]] && echo "UPGRADE_NAME is required" && exit 1

export CHAIN_ID=${CHAIN_ID:-"fxcore"}
export NODE_HOME=${NODE_HOME:-"$HOME/.fxcore"}
export BINARY=${BINARY:-"fxcored"}
! test -f "$BINARY" && echo "Binary not found at $BINARY"

$BINARY config chain-id "$CHAIN_ID" --home "$NODE_HOME"
$BINARY config output json --home "$NODE_HOME"
$BINARY config keyring-backend test --home "$NODE_HOME"
$BINARY config node tcp://localhost:26657 --home "$NODE_HOME"
$BINARY config broadcast-mode block --home "$NODE_HOME"
$BINARY config --home "$NODE_HOME"

while true; do
  sync_state=$("$BINARY" status --home "$NODE_HOME" | jq -r '.SyncInfo.catching_up')
  if [ "$sync_state" != "false" ]; then
    echo "Node is syncing..." && sleep 0.1 && continue
  fi
  break
done

if [ -z "$($BINARY keys show fx1 --home "$NODE_HOME")" ]; then
  echo "test test test test test test test test test test test junk" | $BINARY --home "$NODE_HOME" keys add fx1 --recover
fi

upgrade_height=$($BINARY status --home "$NODE_HOME" | jq -r '.SyncInfo.latest_block_height|tonumber + '"${UPGRADE_HEIGHT_INTERVAL}"'')
readonly upgrade_height
printf "\n"
echo "Upgrade Height = ${upgrade_height}"
printf "Submitting proposal... \n"
$BINARY tx gov submit-proposal software-upgrade "$UPGRADE_NAME" \
  --title "$UPGRADE_NAME" \
  --deposit "$($BINARY query gov params --home "$NODE_HOME" | jq -r '.deposit_params.min_deposit[0].amount')FX" \
  --upgrade-height "${upgrade_height}" \
  --upgrade-info "upgrade to $UPGRADE_NAME" \
  --description "upgrade to $UPGRADE_NAME" \
  --gas auto --gas-prices 4000000000000FX --gas-adjustment=1.3 \
  --from fx1 \
  --home "${NODE_HOME}" \
  --yes
printf "Done \n"

printf "Casting vote... \n"
PROPOSAL_ID=$($BINARY query gov proposals --status=voting_period --home "$NODE_HOME" | jq -r '.proposals[0].PROPOSAL_ID')

echo "Vote ProposalID  =  ${PROPOSAL_ID}"

$BINARY tx gov vote "${PROPOSAL_ID}" yes \
  --gas auto --gas-prices 4000000000000FX --gas-adjustment=1.3 \
  --from fx1 \
  --home "${NODE_HOME}" \
  --yes

printf "Done \n"
