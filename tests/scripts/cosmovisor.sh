#!/usr/bin/env bash

set -eo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

export NODE_HOME="$OUT_DIR/.$CHAIN_NAME"
export DOCKER_IMAGE="ghcr.io/functionx/fxcorevisor:latest"
export DAEMON="docker run --rm -i --network $DOCKER_NETWORK --entrypoint /root/.fxcore/cosmovisor/genesis/bin/fxcored -v $NODE_HOME:$NODE_HOME $DOCKER_IMAGE"

function run_cosmovisor() {
  local init=${1:-""}

  docker rm -f fxcore

  if [[ "$init" == "init" ]]; then
    gen_cosmos_genesis
    json_processor "$NODE_HOME/config/genesis.json" '.app_state.gov.voting_params.voting_period = "15s"'
    json_processor "$NODE_HOME/config/genesis.json" '.initial_height = "2100000"'
  fi

  docker run -d --name fxcore \
    -p 127.0.0.1:26657:26657 -p 127.0.0.1:1317:1317 \
    -p 127.0.0.1:8545:8545 -p 127.0.0.1:8546:8546 \
    -v "$NODE_HOME"/data:/root/.fxcore/data \
    -v "$NODE_HOME"/config:/root/.fxcore/config \
    -v "$NODE_HOME"/keyring-test:/root/.fxcore/keyring-test \
    "$DOCKER_IMAGE" start --x-crisis-skip-assert-invariants
}

function send_upgrade_proposal() {
  upgrade_name=${1:-"$UPGRADE_NAME"}
  [[ -z "$upgrade_name" ]] && echo "upgrade name is required" && exit 1
  upgrade_height_interval=${2:-20}

  node_catching_up

  export NODE_HOME="/root/.fxcore"
  export DAEMON="docker exec -i fxcore /root/.fxcore/cosmovisor/current/bin/fxcored"
  upgrade_height=$($DAEMON status --home "$NODE_HOME" | jq -r '.SyncInfo.latest_block_height|tonumber + '"${upgrade_height_interval}"'')
  readonly upgrade_height
  echo "Upgrade Height = ${upgrade_height}"

  printf "Submitting proposal... \n"
  deposit=$(cosmos_query gov params | jq -r '.deposit_params.min_deposit[0].amount')$STAKING_DENOM
  proposal_cmd="submit-proposal"
  if [[ "${upgrade_name}" =~ v[0-9]+\.[0-9]+\.* ]]; then
    deposit=$(cosmos_query gov params | jq -r '.params.min_deposit[0].amount')$STAKING_DENOM
    proposal_cmd="submit-legacy-proposal"
    export GAS_ADJUSTMENT=1.5
  fi
  cosmos_tx gov "$proposal_cmd" software-upgrade "$upgrade_name" \
    --title "$upgrade_name" \
    --deposit "${deposit}" \
    --upgrade-height "${upgrade_height}" \
    --upgrade-info "upgrade to $upgrade_name" \
    --description "upgrade to $upgrade_name" \
    --no-validate=true \
    --from "${FROM}"

  printf "Casting vote... \n"
  PROPOSAL_ID=$(cosmos_query gov proposals --status=voting_period | jq -r '.proposals[0].proposal_id // .proposals[0].id')
  echo "Vote ProposalID  =  ${PROPOSAL_ID}"

  cosmos_tx gov vote "${PROPOSAL_ID}" yes --from "${FROM}"
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
