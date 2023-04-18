#!/usr/bin/env bash

set -eo pipefail

commands=(jq curl)
for cmd in "${commands[@]}"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "$cmd command not found, please install $cmd first" && exit 1
  fi
done

export REST_RPC=${REST_RPC:-"https://fx-rest.functionx.io"}

function format_balance() {
  local balances=$1
  local denom=$2
  amount=$(echo "$balances" | jq -r '.[]|select(.denom == "'"$denom"'")|.amount')
  echo | awk "{printf(\"%.2f\",${amount:-0}/1e${DECIMALS:-18})}"
}

printf "%-25s %-45s %-20s %s\n" "name" "address" "permissions" "balance"
while read -r name address permissions; do
  balances=$(curl -s "$REST_RPC/cosmos/bank/v1beta1/balances/$address" | jq -c '.balances')
  printf "%-25s %-45s %-20s %s\n" "$name" "$address" "$permissions" "$(format_balance "$balances" FX)FX"
done < <(curl -s "$REST_RPC/cosmos/auth/v1beta1/module_accounts" | jq -r '.accounts[]|"\(.name) \(.base_account.address) \(.permissions)"')
