#!/usr/bin/env bash

set -eo pipefail

export REST_RPC=${REST_RPC:-"https://fx-rest.functionx.io"}
export BECH32_PREFIX=${BECH32_PREFIX:-"fx"}
export MINT_DENOM=${MINT_DENOM:-"FX"}

if [ -z "$BECH32_PREFIX" ]; then
  BECH32_PREFIX=$(curl -s "$REST_RPC/cosmos/auth/v1beta1/bech32" | jq -r '.bech32_prefix') || echo "failed to get bech32_prefix"
fi

if [ -z "$MINT_DENOM" ]; then
  MINT_DENOM=$(curl -s "$REST_RPC/cosmos/mint/v1beta1/params" | jq -r '.params.mint_denom') || echo "failed to get mint_denom"
fi

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

function show_mint_info() {
  printf "mint module docs: https://docs.cosmos.network/main/modules/mint/\n\n"
  # query mint module params
  min_params=$(curl -s "$REST_RPC/cosmos/mint/v1beta1/params" | jq -r '.params')
  # mint denom
  mint_denom=$(echo "$min_params" | jq -r '.mint_denom')
  printf "current mint module params: (cant be modified by proposal)\n"

  # maximum annual change in inflation rate
  inflation_rate_change=$(echo "$min_params" | jq -r '.inflation_rate_change')
  printf "    minting inflation rate change:\t%s\n" "$(echo "scale=2;$inflation_rate_change*100/1" | bc)%"

  # maximum inflation rate
  max_inflation_rate=$(echo "$min_params" | jq -r '.inflation_max')
  printf "    minting max inflation rate:\t\t%s\n" "$(echo "scale=2;$max_inflation_rate*100/1" | bc)%"

  # minimum inflation rate
  min_inflation_rate=$(echo "$min_params" | jq -r '.inflation_min')
  printf "    minting min inflation rate:\t\t%s\n" "$(echo "scale=2;$min_inflation_rate*100/1" | bc)%"

  # goal of percent bonded atoms
  goal_bonded=$(echo "$min_params" | jq -r '.goal_bonded')
  printf "    minting goal bonded:\t\t%s\n" "$(echo "scale=2;$goal_bonded*100/1" | bc)%"

  # expected blocks per year
  blocks_per_year=$(echo "$min_params" | jq -r '.blocks_per_year')
  printf "    minting blocks per year:\t\t%s\n" "$blocks_per_year"
  printf "        |-> ⚠️blocks_per_year here is calculated based on 5s block time, the actual block time is not fixed, so this value is inaccurate\n"
  printf "\n"

  # query total supply by mint denom
  total_supply=$(curl -s "$REST_RPC/cosmos/bank/v1beta1/supply" | jq -r ".supply[]|select(.denom == \"$mint_denom\")|.amount")
  printf "current total supply:\t\t\t%s\n" "$(echo "scale=2;$total_supply/10^18" | bc)"

  # query staking pool
  staking_pool=$(curl -s "$REST_RPC/cosmos/staking/v1beta1/pool" | jq -r '.pool')

  # bonded tokens
  bonded_tokens=$(echo "$staking_pool" | jq -r '.bonded_tokens')
  printf "current staking bonded tokens:\t\t%s\n" "$(echo "scale=2;$bonded_tokens/10^18" | bc)"

  # bonded ratio = bonded_tokens / total_supply
  bonded_ratio=$(printf "%.6f" "$(echo "scale=6;$bonded_tokens*100/$total_supply" | bc)")
  printf "current staking bonded ratio:\t\t%s\n" "$bonded_ratio%"
  printf "    |-> ⚠️bonded_ratio = bonded_tokens / total_supply\n"
  printf "\n"

  # query minting inflation
  inflation=$(curl -s "$REST_RPC/cosmos/mint/v1beta1/inflation" | jq -r '.inflation')
  printf "current minting inflation:\t\t%s\n" "$(echo "scale=2;$inflation*100/1" | bc)%"
  printf "    |-> ⚠️inflation = latest_inflation + ((1 - bonded_ratio/goal_bonded) * inflation_rate_change) / blocks_per_year\n"

  # query annual provisions
  annual_provisions=$(curl -s "$REST_RPC/cosmos/mint/v1beta1/annual_provisions" | jq -r '.annual_provisions')
  printf "current minting annual provisions:\t%s\n" "$(echo "scale=2;$annual_provisions/10^18" | bc)"
  printf "    |-> ⚠️annual_provisions = inflation * total_supply\n"

  # average inflation per block = annual_provisions / blocks_per_year
  average_inflation_per_block=$(echo "scale=2;$annual_provisions/$blocks_per_year" | bc)
  printf "average inflation per block:\t\t%s\n" "$(echo "scale=2;$average_inflation_per_block/10^18" | bc)"
  printf "\n"
}

function help() {
  printf "This script is used to stats cosmos chain info.\n"
  printf "Usage:\n"
  printf "    %s <command> [args]\n" "$0"
  printf "The commands are:\n"
  printf "    show_validator_rewards \t show validator rewards\n"
  printf "    show_validator_votes \t show validator votes\n"
  printf "    show_module_account \t show all module accounts\n"
  printf "    show_mint_info \t\t show mint module info\n"
  printf "    help \t\t\t show this help message\n"
  printf "\nVersion: Alpha\n"
}

[[ "$#" -eq 0 || "$1" == "help" ]] && help && exit 0
[[ "$#" -gt 0 && "$(type -t "$1")" != "function" ]] && echo "invalid args: $1" && exit 1
"$@" || (echo "failed: $0" "$@" && exit 1)
