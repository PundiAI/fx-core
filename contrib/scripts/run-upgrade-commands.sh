#!/bin/sh

set -o errexit -o nounset

# check dependencies commands are installed
commands=(jq curl)
for cmd in "${commands[@]}"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "$cmd command not found, please install $cmd first" && exit 1
  fi
done

UPGRADE_HEIGHT=${1:-6}
proposal_id=${2:-""}

#if [ -z "$1" ]; then
#  echo "Need to set an upgrade height as the first argument: ex. ./run-upgrade-commands.sh 10"
#  exit 1
#fi

#NODE_HOME=$(realpath ./build/.fxcore)
NODE_HOME=${HOME}/.fxcore
echo "NODE_HOME = ${NODE_HOME}"

#BINARY=$NODE_HOME/cosmovisor/genesis/bin/fxcored
BINARY=$GOPATH/bin/fxcored
echo "BINARY = ${BINARY}"

USER_MNEMONIC="test test test test test test test test test test test junk"
CHAINID=fxcore

if test -f "$BINARY"; then

  echo "wait 3 seconds for blockchain to start"
  sleep 3

	$BINARY config chain-id $CHAINID --home "$NODE_HOME"
	$BINARY config output json --home "$NODE_HOME"
	$BINARY config keyring-backend test --home "$NODE_HOME"
  $BINARY config --home "$NODE_HOME"


  key=$($BINARY keys show fx1 --home "$NODE_HOME")

  if [ -z "$key" ]; then
    echo "$USER_MNEMONIC" | $BINARY --home "$NODE_HOME" keys add fx1 --recover --keyring-backend=test
  fi

  # $BINARY keys list --home $NODE_HOME

  upgrade_height=`expr $($BINARY status -o json | jq -r '.SyncInfo.latest_block_height|tonumber') + ${UPGRADE_HEIGHT}`
  printf "\n"
  printf "Submitting proposal... \n"
  $BINARY tx gov submit-proposal software-upgrade fxv4 \
  --title fxv4 \
  --deposit "$($BINARY q gov params | jq -r '.deposit_params.min_deposit[0].amount')FX" \
  --upgrade-height ${upgrade_height} \
  --upgrade-info "upgrade to fxv4" \
  --description "upgrade to fxv4" \
  --gas auto \
  --gas-prices 4000000000000FX \
  --gas-adjustment=1.3 \
  --from fx1 \
  --keyring-backend test \
  --chain-id $CHAINID \
  --home "${NODE_HOME}" \
  --node tcp://localhost:26657 \
  --broadcast-mode block \
  --yes
  echo "Done \n"

  echo "Casting vote... \n"


  if [ -z "$proposal_id" ]; then
    proposal_id=$($BINARY q gov proposals --status=voting_period | jq -r '.proposals[0].proposal_id')
  fi


  $BINARY tx gov vote "${proposal_id}" yes \
  --gas auto \
  --gas-prices 4000000000000FX\
  --gas-adjustment=1.3 \
  --from fx1 \
  --keyring-backend test \
  --chain-id $CHAINID \
  --home "${NODE_HOME}" \
  --node tcp://localhost:26657 \
  --broadcast-mode block \
  --yes

  printf "Done \n"

else
  echo "Please build fxcored v3 and move to $GOPATH/bin/fxcored"
fi
