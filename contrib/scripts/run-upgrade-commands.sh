#!/usr/bin/env bash

set -o errexit -o nounset

# check dependencies commands are installed
commands=(jq curl)
for cmd in "${commands[@]}"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "$cmd command not found, please install $cmd first" && exit 1
  fi
done

UPGRADE_HEIGHT=${1:-10}
proposal_id=${2:-""}

DEFAULT_NODE_HOME=${HOME}/.fxcore
NODE_HOME=${FX_RUN_HOME:-$DEFAULT_NODE_HOME}
echo "NODE_HOME = ${NODE_HOME}"

DEFAULT_BINARY=$GOPATH/bin/fxcored
BINARY=${FX_RUN_BINARY:-$DEFAULT_BINARY}
echo "BINARY = ${BINARY}"

USER_MNEMONIC="test test test test test test test test test test test junk"
CHAINID=fxcore

if test -f "$BINARY"; then

  echo "wait 3 seconds for blockchain to start"
  sleep 3

	$BINARY config chain-id $CHAINID --home "$NODE_HOME"
	$BINARY config output json --home "$NODE_HOME"
	$BINARY config keyring-backend test --home "$NODE_HOME"
  $BINARY config node tcp://localhost:26657 --home "$NODE_HOME"
  $BINARY config broadcast-mode block --home "$NODE_HOME"
  $BINARY config --home "$NODE_HOME"

  key=$($BINARY keys show fx1 --home "$NODE_HOME")

  if [ -z "$key" ]; then
    echo "$USER_MNEMONIC" | $BINARY --home "$NODE_HOME" keys add fx1 --recover --keyring-backend=test
  fi

  # $BINARY keys list --home $NODE_HOME

  upgrade_height=`expr $($BINARY status -o json | jq -r '.SyncInfo.latest_block_height|tonumber') + ${UPGRADE_HEIGHT}`
  printf "\n"
  echo "Upgrade Height = ${upgrade_height}"
  printf "Submitting proposal... \n"
  $BINARY tx gov submit-proposal software-upgrade fxv4 \
  --title fxv4 \
  --deposit "$($BINARY q gov params | jq -r '.deposit_params.min_deposit[0].amount')FX" \
  --upgrade-height "${upgrade_height}" \
  --upgrade-info "upgrade to fxv4" \
  --description "upgrade to fxv4" \
  --gas auto \
  --gas-prices 4000000000000FX \
  --gas-adjustment=1.3 \
  --from fx1 \
  --home "${NODE_HOME}" \
  --yes
  printf "Done \n"

  printf "Casting vote... \n"


  if [ -z "$proposal_id" ]; then
    proposal_id=$($BINARY q gov proposals --status=voting_period | jq -r '.proposals[0].proposal_id')
  fi

  echo "Vote ProposalID  =  ${proposal_id}"

  $BINARY tx gov vote "${proposal_id}" yes \
  --gas auto \
  --gas-prices 4000000000000FX\
  --gas-adjustment=1.3 \
  --from fx1 \
  --home "${NODE_HOME}" \
  --yes

  printf "Done \n"

else
  echo "Please build fxcored v3 and move to $GOPATH/bin/fxcored"
fi
