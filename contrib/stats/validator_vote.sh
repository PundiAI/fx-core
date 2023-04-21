#!/usr/bin/env bash

set -eo pipefail

commands=(jq curl)
for cmd in "${commands[@]}"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "$cmd command not found, please install $cmd first" && exit 1
  fi
done

export REST_RPC=${REST_RPC:-"https://fx-rest.functionx.io"}

if [ -z "$PROPOSAL_ID" ]; then
  {
    echo "proposal_id#status#title"
    curl -s "$REST_RPC/cosmos/gov/v1beta1/proposals?pagination.reverse=true&pagination.limit=10" | jq -r '.proposals[]|"\(.proposal_id)#\(.status)#\(.content.title)"'
  } | column -t -s"#"
  read -r -p "Please select a proposal id: " PROPOSAL_ID
fi

bech32_prefix=$(curl -s "$REST_RPC/cosmos/auth/v1beta1/bech32" | jq -r '.bech32_prefix')

{
  echo "moniker#operator_address#acc_address#proposal_id#vote_option"
  while read -r operator_address moniker; do
    acc_address=$(curl -s "$REST_RPC/fx/auth/v1/bech32/$operator_address?prefix=${bech32_prefix}" | jq -r '.address')
    option=$(curl -s "$REST_RPC/cosmos/tx/v1beta1/txs?events=message.sender='$acc_address'&proposal_vote.proposal_id='$PROPOSAL_ID'" | jq -r '.txs[].tx.body.messages[].option')
    option=${option:-"null"}
    echo "$moniker#$operator_address#$acc_address#$PROPOSAL_ID#$option"
  done < <(curl -s "$REST_RPC/cosmos/staking/v1beta1/validators" | jq -r '.validators[]|"\(.operator_address) \(.description.moniker)"')
} | column -t -s"#"
