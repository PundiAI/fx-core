#!/usr/bin/env bash

set -eo pipefail

export REST_RPC=${REST_RPC:-"http://localhost:1317"}
export BECH32_PREFIX=${BECH32_PREFIX:-"fx"}
export MINT_DENOM=${MINT_DENOM:-"FX"}

if [ -z "$BECH32_PREFIX" ]; then
  BECH32_PREFIX=$(curl -s "$REST_RPC/cosmos/auth/v1beta1/bech32" | jq -r '.bech32_prefix') || echo "failed to get bech32_prefix"
fi

if [ -z "$MINT_DENOM" ]; then
  MINT_DENOM=$(curl -s "$REST_RPC/cosmos/mint/v1beta1/params" | jq -r '.params.mint_denom') || echo "failed to get mint_denom"
fi

## DESC: show validator reward
function show_validator_reward() {
  local decimals=18
  {
    echo "moniker#operator_address#jailed#status#commission_rate#self_delegated#3rd_party_delegated#block_reward#tx_fee_reward"
    while read -r operator_address jailed status tokens commission_rate moniker; do
      acc_address=$(curl -s "$REST_RPC/fx/auth/v1/bech32/$operator_address?prefix=${BECH32_PREFIX}" | jq -r '.address')
      self_delegated=$(curl -s "$REST_RPC/cosmos/staking/v1beta1/delegations/$acc_address" | jq -r '.delegation_responses[]|select(.delegation.validator_address == "'"$operator_address"'")|.balance.amount')
      local block_reward=0
      local tx_fee_reward=0
      while read -r denom amount; do
        if [ "$denom" == "$MINT_DENOM" ]; then
          block_reward=$amount
        else
          tx_fee_reward=$amount
        fi
      done < <(curl -s "$REST_RPC/cosmos/distribution/v1beta1/delegators/$acc_address/rewards" | jq -rc '.total[]|"\(.denom) \(.amount)"')
      while read -r denom amount; do
        if [ "$denom" == "$MINT_DENOM" ]; then
          block_reward=$(echo "$block_reward+$amount" | bc)
        else
          tx_fee_reward=$(echo "$tx_fee_reward+$amount" | bc)
        fi
      done < <(curl -s "$REST_RPC/cosmos/distribution/v1beta1/validators/$operator_address/commission" | jq -rc '.commission.commission[]|"\(.denom) \(.amount)"')
      self_delegated=${self_delegated:-0}
      party_delegated=$(echo "$tokens-$self_delegated" | bc)

      commission_rate=$(echo "scale=2;($commission_rate * 100)/1" | bc -l)
      self_delegated=$(echo "$self_delegated / 10^$decimals" | bc)
      party_delegated=$(echo "$party_delegated / 10^$decimals" | bc)
      block_reward=$(echo "$block_reward / 10^$decimals" | bc)
      echo "$moniker#$operator_address#$jailed#$status#$commission_rate%#$self_delegated#$party_delegated#$block_reward#$tx_fee_reward"
    done < <(curl -s "$REST_RPC/cosmos/staking/v1beta1/validators" | jq -r '.validators[]|"\(.operator_address) \(.jailed) \(.status) \(.tokens) \(.commission.commission_rates.rate) \(.description.moniker)"')
  } | column -t -s"#"
}

## DESC: show validator vote
function show_validator_vote() {
  if [ -z "$PROPOSAL_ID" ]; then
    {
      echo "proposal_id#status#title"
      curl -s "$REST_RPC/cosmos/gov/v1beta1/proposals?pagination.reverse=true&pagination.limit=10" | jq -r '.proposals[]|"\(.proposal_id)#\(.status)#\(.content.title)"'
    } | column -t -s"#"
    read -r -p "Please select a proposal id: " PROPOSAL_ID
  fi

  {
    echo "moniker#operator_address#acc_address#proposal_id#vote_option"
    while read -r operator_address moniker; do
      acc_address=$(curl -s "$REST_RPC/fx/auth/v1/bech32/$operator_address?prefix=${BECH32_PREFIX}" | jq -r '.address')
      option=$(curl -s "$REST_RPC/cosmos/tx/v1beta1/txs?events=message.sender='$acc_address'&proposal_vote.proposal_id='$PROPOSAL_ID'" | jq -r '.txs[].tx.body.messages[].option')
      option=${option:-"null"}
      echo "$moniker#$operator_address#$acc_address#$PROPOSAL_ID#$option"
    done < <(curl -s "$REST_RPC/cosmos/staking/v1beta1/validators" | jq -r '.validators[]|"\(.operator_address) \(.description.moniker)"')
  } | column -t -s"#"
}

## DESC: show module accounts
function show_module_account() {
  local denom=${1:-"$MINT_DENOM"}
  local decimals=18

  printf "%-25s %-45s %-20s %s\n" "name" "address" "permissions" "balance"
  while read -r name address permissions; do
    balance=$(curl -s "$REST_RPC/cosmos/bank/v1beta1/balances/$address" | jq -r ".balances[]|select(.denom == \"$denom\")|.amount")
    balance=$(echo "${balance:-0} / 10^$decimals" | bc)
    printf "%-25s %-45s %-20s %s\n" "$name" "$address" "$permissions" "$balance$denom"
  done < <(curl -s "$REST_RPC/cosmos/auth/v1beta1/module_accounts" | jq -r '.accounts[]|"\(.name) \(.base_account.address) \(.permissions)"')
}

[[ "$#" -gt 0 && "$(type -t "$1")" != "function" ]] && echo "invalid command: $1" && exit 1
"$@" || (echo "failed: $0" "$@" && exit 1)
