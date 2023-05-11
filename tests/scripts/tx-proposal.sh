#!/usr/bin/env bash

set -eo pipefail

PROJECT_DIR="${PROJECT_DIR:-"$(git rev-parse --show-toplevel)"}"
export PROJECT_DIR
export OUT_DIR="${PROJECT_DIR}/out"

readonly proposals_file="${PROJECT_DIR}/tests/data/proposals.json"

## ARGS: <title> <summary>
## DESC: base64 encode metadata
function base64_metadata() {
  echo '{"title": "'"$1"'","summary": "'"$2"'","metadata":""}' | base64
}

## ARGS: <msg_type>
## DESC: get proposal template
function get_proposal_template() {
  local msg_type=$1
  jq -r --arg msg_type "$msg_type" '.[]|select(.msg_type == $msg_type)' "$proposals_file" >"$OUT_DIR/${msg_type##*.}.json"
}

## ARGS: <msg_type> [<amount>]
## DESC: query min deposit
function query_min_deposit() {
  local msg_type=$1 amount=$2
  if [[ -z "$msg_type" ]]; then
    echo "$(cosmos_query gov params | jq -r '.params.min_deposit|select(.denom=="'"$STAKING_DENOM"'")|.amount')$STAKING_DENOM" && return
  fi

  base_deposit="$(cosmos_query gov params --msg-type="$msg_type" | jq -r '.params.min_deposit|select(.denom=="'"$STAKING_DENOM"'")|.amount')"
  if [[ "$msg_type" == "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal" && -n "$amount" ]]; then
    deposit_threshold=$(cosmos_query gov egf-params | jq -r '.params.egf_deposit_threshold.amount')
    claim_ratio=$(cosmos_query gov egf-params | jq -r '.params.claim_ratio')

    amount_without=${amount%"$STAKING_DENOM"}
    if [[ $(echo "$amount_without - $deposit_threshold" | bc) -gt 0 ]]; then
      echo "$(echo "$amount_without * $claim_ratio" | bc)""$STAKING_DENOM"
    fi
  fi
  echo "${base_deposit}${STAKING_DENOM}"
}

## ARGS:  <option> [<proposal_id>]
## DESC: vote proposal
function vote() {
  local option=$1 proposal_id=${2:-""}

  for proposal_id in $(cosmos_query gov proposals | jq -r '.proposals[].id'); do
    messages=$(cosmos_query gov proposal "${proposal_id}")
    if [[ "$(echo "${messages}" | jq -r '.status')" == "PROPOSAL_STATUS_VOTING_PERIOD" ]]; then
      break
    fi
  done

  cosmos_tx gov vote "${proposal_id}" "$option" --from "$FROM"

  #  voting_period=$(cosmos_query gov params | jq -r '.params.voting_period')
  while true; do
    if [ "$($DAEMON query gov proposal "$proposal_id" | jq -r '.status')" != "PROPOSAL_STATUS_VOTING_PERIOD" ]; then
      break
    fi
    echo "wait for voting period"
    sleep 1
    #    sleep "${voting_period%?}"
  done
}

## ARGS: <proposal_file>
## DESC: submit proposal
function submit_proposal() {
  local proposal_file=$1
  msg_type=$(jq -r '.msg_type' "$proposal_file")

  if [ "$msg_type" == "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal" ]; then
    deposit=$(query_min_deposit "$msg_type" "$(jq -r '.amount' "$proposal_file")")
    json_processor "$proposal_file" '.deposit = "'"$deposit"'"'

    cosmos_tx gov submit-legacy-proposal community-pool-spend "$proposal_file" --from "$FROM"
  else
    title=$(jq -r '.title' "$proposal_file")
    summary=$(jq -r '.summary' "$proposal_file")
    metadata=$(base64_metadata "$title" "$summary")
    json_processor "$proposal_file" '.metadata = "'"$metadata"'"'

    deposit=$(query_min_deposit "$msg_type")
    json_processor "$proposal_file" '.deposit = "'"$deposit"'"'

    cosmos_tx gov submit-proposal "$proposal_file" --from "$FROM"
  fi
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"
