#!/usr/bin/env bash

set -eo pipefail

commands=(jq curl)
for cmd in "${commands[@]}"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "$cmd command not found, please install $cmd first" && exit 1
  fi
done

export REST_RPC=${REST_RPC:-"https://fx-rest.functionx.io"}

bech32_prefix=$(curl -s "$REST_RPC/cosmos/auth/v1beta1/bech32" | jq -r '.bech32_prefix')

mint_denom=$(curl -s "$REST_RPC/cosmos/mint/v1beta1/params" | jq -r '.params.mint_denom')

{
  echo "moniker#operator_address#jailed#status#commission_rate#self_delegated#3rd_party_delegated#block_reward#tx_fee_reward"
  while read -r operator_address jailed status tokens commission_rate moniker; do
    acc_address=$(curl -s "$REST_RPC/fx/auth/v1/bech32/$operator_address?prefix=${bech32_prefix}" | jq -r '.address')
    self_delegated=$(curl -s "$REST_RPC/cosmos/staking/v1beta1/delegations/$acc_address" | jq -r '.delegation_responses[]|select(.delegation.validator_address == "'"$operator_address"'")|.balance.amount')
    block_reward=0
    tx_fee_reward=0
    while read -r denom amount; do
      if [ "$denom" == "$mint_denom" ]; then
        block_reward=$amount
      else
        tx_fee_reward=$amount
      fi
    done < <(curl -s "$REST_RPC/cosmos/distribution/v1beta1/delegators/$acc_address/rewards" | jq -rc '.total[]|"\(.denom) \(.amount)"')
    while read -r denom amount; do
      if [ "$denom" == "$mint_denom" ]; then
        block_reward=$(echo "$block_reward+$amount" | bc)
      else
        tx_fee_reward=$(echo "$tx_fee_reward+$amount" | bc)
      fi
    done < <(curl -s "$REST_RPC/cosmos/distribution/v1beta1/validators/$operator_address/commission" | jq -rc '.commission.commission[]|"\(.denom) \(.amount)"')
    self_delegated=${self_delegated:-0}
    party_delegated=$(echo "$tokens-$self_delegated" | bc)
    echo "$moniker#$operator_address#$jailed#$status#$commission_rate#$self_delegated#$party_delegated#$block_reward#$tx_fee_reward"
  done < <(curl -s "$REST_RPC/cosmos/staking/v1beta1/validators" | jq -r '.validators[]|"\(.operator_address) \(.jailed) \(.status) \(.tokens) \(.commission.commission_rates.rate) \(.description.moniker)"')
} | column -t -s"#"
