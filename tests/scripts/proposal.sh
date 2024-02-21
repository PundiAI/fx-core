#!/usr/bin/env bash

set -eo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

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

## ARGS: [<msg_type>] [<amount>]
## DESC: query min deposit
function query_min_deposit() {
  local msg_type=$1 amount=$2
  if [[ -z "$msg_type" ]]; then
    if [[ "$(cosmos_version | grep "v3")" != "" ]]; then
      echo "$(cosmos_query gov params | jq -r '.deposit_params.min_deposit[]|select(.denom=="'"$STAKING_DENOM"'")|.amount')$STAKING_DENOM" && return
    fi
    echo "$(cosmos_query gov params | jq -r '.params.min_deposit[]|select(.denom=="'"$STAKING_DENOM"'")|.amount')$STAKING_DENOM" && return
  fi

  base_deposit="$(cosmos_query gov params --msg-type="$msg_type" | jq -r '.params.min_deposit[]|select(.denom=="'"$STAKING_DENOM"'")|.amount')"
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

## ARGS: <proposal_id> <deposit_amount>
function deposit() {
  local proposal_id=$1 deposit_amount=$2
  cosmos_tx gov deposit "$proposal_id" "$(to_18 "$deposit_amount")$STAKING_DENOM" --from "$FROM"
}

## ARGS: <option> [<proposal_id>]
## DESC: vote proposal
function vote() {
  local option=$1 proposal_id=${2:-""}

  if [[ -z "$proposal_id" ]]; then
    proposal_id="$(cosmos_query gov proposals --reverse --limit 1 --status=voting_period | jq -r '.proposals[0].proposal_id // .proposals[0].id')"
  fi

  [[ "$(cosmos_query gov proposal "${proposal_id}" | jq -r '.status')" != "PROPOSAL_STATUS_VOTING_PERIOD" ]] &&
    echo "proposal is not in voting period" && return

  cosmos_tx gov vote "${proposal_id}" "$option" --from "$FROM"

  while true; do
    [[ "$($DAEMON query gov proposal "$proposal_id" | jq -r '.status')" != "PROPOSAL_STATUS_VOTING_PERIOD" ]] && break
    echo "wait for voting period"
    sleep 1
  done
}

## ARGS: <subspace> <key> <value>
function param_change() {
  local subspace=$1 key=$2 value=$3

  min_deposit=$(query_min_deposit)
  cosmos_tx gov submit-legacy-proposal param-change <(
    cat <<EOF
{
  "title":"Change Genesis Params",
  "description": "test",
  "changes": [
    {
      "subspace": "$subspace",
      "key": "$key",
      "value": "$value"
    }
  ],
  "deposit": "$min_deposit$STAKING_DENOM"
}
EOF
  ) --from "$FROM"
}

## ARGS: <receive_address> <spend_amount> [<denom>]
function spend_community_pool_proposal() {
  local receive_address=$1 spend_amount=$2 denom=${3:-$STAKING_DENOM}

  current_community_pool_balance=$(cosmos_query distribution community-pool | jq -r ".pool[]|select(.denom == \"$denom\")|.amount")
  echo "the community pool balance is: $(from_18 "$current_community_pool_balance")$denom"

  min_deposit=$(query_min_deposit)
  cosmos_tx gov submit-legacy-proposal community-pool-spend <(
    cat <<EOF
{
  "title":"Spend Community Pool",
  "description": "test",
  "recipient": "$receive_address",
  "amount": "$spend_amount$denom",
  "deposit": "$min_deposit$STAKING_DENOM"
}
EOF
  ) --from "$FROM"
}

## ARGS: <proposal_file>
## DESC: submit proposal
function submit_proposal() {
  local proposal_file=$1
  msg_type=$(jq -r '.msg_type' "$proposal_file")

  if [[ "$msg_type" == "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal" ]]; then
    deposit=$(query_min_deposit "$msg_type" "$(jq -r '.amount' "$proposal_file")")
    json_processor "$proposal_file" '.deposit = "'"$deposit"'"'

    cosmos_tx gov submit-legacy-proposal community-pool-spend "$proposal_file" --from "$FROM"
  else
    title=$(jq -r '.title' "$proposal_file")
    summary=$(jq -r '.summary' "$proposal_file")
    metadata=$(base64_metadata "$title" "$summary")
    json_processor "$proposal_file" '.proposal.metadata = "'"$metadata"'"'

    deposit=$(query_min_deposit "$msg_type")
    json_processor "$proposal_file" '.proposal.deposit = "'"$deposit"'"'
    json_processor "$proposal_file" -r '.proposal'
    cosmos_tx gov submit-proposal "$proposal_file" --from "$FROM"
  fi
}

## ARGS: <base_denom> <name> <symbol> <decimals> [aliases...]
## DESC: register coin
function register_coin() {
  local base_denom=$1 name=$2 symbol=$3 decimals=$4 && shift 4
  aliases=$(IFS=$',' && echo "$*")

  get_proposal_template "/fx.erc20.v1.MsgRegisterCoin"
  register_file="$OUT_DIR/MsgRegisterCoin.json"

  cat >"$OUT_DIR/coin.json" <<EOF
{
      "description": "The cross chain token of the Function X",
      "denom_units": [
        {
          "denom": "$base_denom",
          "exponent": 0,
          "aliases": [$aliases]
        },
        {
          "denom": "$symbol",
          "exponent": $decimals,
          "aliases": []
        }
      ],
      "base": "$base_denom",
      "display": "$base_denom",
      "name": "$name",
      "symbol": "$symbol"
}
EOF
  metadata=$(jq -r -c . "$OUT_DIR/coin.json")
  json_processor "$register_file" ".proposal.messages[0].metadata = $metadata"

  submit_proposal "$register_file"
  vote yes
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
