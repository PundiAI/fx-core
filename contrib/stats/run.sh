#!/usr/bin/env bash

set -eo pipefail

export REST_RPC=${REST_RPC:-"https://fx-rest.functionx.io"}
export BECH32_PREFIX=${BECH32_PREFIX:-"fx"}
export MINT_DENOM=${MINT_DENOM:-"FX"}

function check_command() {
  commands=("$@")
  for cmd in "${commands[@]}"; do
    if ! command -v "$cmd" &>/dev/null; then
      echo "$cmd command not found, please install $cmd first" && exit 1
    fi
  done
}
check_command jq bc curl

if [ -z "$BECH32_PREFIX" ]; then
  BECH32_PREFIX=$(curl -s "$REST_RPC/cosmos/auth/v1beta1/bech32" | jq -r '.bech32_prefix') || echo "failed to get bech32_prefix"
fi

if [ -z "$MINT_DENOM" ]; then
  MINT_DENOM=$(curl -s "$REST_RPC/cosmos/mint/v1beta1/params" | jq -r '.params.mint_denom') || echo "failed to get mint_denom"
fi

function datetime_since() {
  python3 -c "import time; print((time.mktime(time.strptime('${1%.*}', '%Y-%m-%dT%H:%M:%S'))-time.mktime(time.strptime('${2%.*}', '%Y-%m-%dT%H:%M:%S'))))"
}

function datetime_add() {
  second=$(echo "${2:-0}" | bc)
  python3 -c "import time; print(time.strftime('%Y-%m-%dT%H:%M:%S', time.localtime(time.mktime(time.strptime('${1%.*}', '%Y-%m-%dT%H:%M:%S'))+$second)))"
}

function get_latest_block_and_time() {
  # get latest block header
  latest_block_header=$(curl -s "$REST_RPC/cosmos/base/tendermint/v1beta1/blocks/latest" | jq -r '.block.header')
  # get latest block time
  latest_block_time=$(echo "$latest_block_header" | jq -r '.time')
  # get latest block height
  latest_block_height=$(echo "$latest_block_header" | jq -r '.height')
  echo "$latest_block_height $latest_block_time"
}

function avg_block_time_interval() {
  local block_interval=20000
  read -r latest_block_height latest_block_time < <(get_latest_block_and_time)
  # get block time of latest_block - block_interval
  block_time=$(curl -s "$REST_RPC/cosmos/base/tendermint/v1beta1/blocks/$((latest_block_height - block_interval))" | jq -r '.block.header.time')
  # calculate avg block time interval
  block_time_interval=$(datetime_since "$latest_block_time" "$block_time")
  python3 -c "print($block_time_interval/$block_interval)"
}

function calc_upgrade_height() {
  local upgrade_time=${1:-$(datetime_add "$(date -u +%FT%T)" "14*3600")}
  real_block_time_interval=$(avg_block_time_interval)
  echo "real_block_time_interval: $real_block_time_interval"
  read -r latest_block_height latest_block_time < <(get_latest_block_and_time)
  echo "latest_block_height: $latest_block_height, latest_block_time: $latest_block_time"
  block_time_interval=$(datetime_since "$upgrade_time" "$latest_block_time")
  echo "upgrade height: $(echo "$latest_block_height+($block_time_interval)/$real_block_time_interval" | bc)"
}

function calc_upgrade_time() {
  local upgrade_height=${1:-$(curl -s "$REST_RPC/cosmos/base/tendermint/v1beta1/status" | jq -r '.sync_info.latest_block_height')}
  real_block_time_interval=$(avg_block_time_interval)
  echo "real_block_time_interval: $real_block_time_interval"
  read -r latest_block_height latest_block_time < <(get_latest_block_and_time)
  echo "latest_block_height: $latest_block_height, latest_block_time: $latest_block_time"
  block_time_interval=$(echo "$upgrade_height-$latest_block_height" | bc)
  echo "upgrade time: $(datetime_add "$latest_block_time" "$block_time_interval*$real_block_time_interval")"
}

function show_validator_reward() {
  local decimals=18
  staking_token_denom=$(curl -s "$REST_RPC/cosmos/staking/v1beta1/params" | jq -r '.params.bond_denom')
  {
    if [ "$staking_token_denom" == "$MINT_DENOM" ]; then
      echo "operator_address#jailed#status#commission_rate#self_delegated#3rd_party_delegated#block_reward#moniker"
    else
      echo "operator_address#jailed#status#commission_rate#self_delegated#3rd_party_delegated#block_reward#tx_fee_reward#moniker"
    fi
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

      commission_rate=$(printf "%.2f" "$(echo "scale=2;($commission_rate * 100)/1" | bc -l)")
      self_delegated=$(echo "scale=2;$self_delegated / 10^$decimals" | bc)
      party_delegated=$(echo "scale=2;$party_delegated / 10^$decimals" | bc)
      block_reward=$(echo "scale=2;$block_reward / 10^$decimals" | bc)
      if [ "$staking_token_denom" == "$MINT_DENOM" ]; then
        echo "$operator_address#$jailed#$status#$commission_rate%#$self_delegated#$party_delegated#$block_reward#$moniker"
      else
        echo "$operator_address#$jailed#$status#$commission_rate%#$self_delegated#$party_delegated#$block_reward#$tx_fee_reward#$moniker"
      fi
    done < <(curl -s "$REST_RPC/cosmos/staking/v1beta1/validators" | jq -r '.validators[]|"\(.operator_address) \(.jailed) \(.status) \(.tokens) \(.commission.commission_rates.rate) \(.description.moniker)"')
  } | column -t -s"#"
}

function show_oracle_reward() {
  local json_rpc_url=${1:-"https://fx-json.functionx.io:26657"}
  local DAEMON=${DAEMON:-"fxcored"}
  local decimals=18
  declare -a arr=("eth" "bsc" "polygon" "tron" "avalanche")
  {
    for chain_name in "${arr[@]}"; do
      echo "${chain_name}_oracle_address#delegate_amount#start_height#online#delegate_validator#reward"
      while read -r oracle_address delegate_amount start_height online delegate_validator; do
        reward_amount=$($DAEMON q crosschain "$chain_name" reward "$oracle_address" --node "$json_rpc_url" | jq -r '.total[0].amount')
        reward_amount=$(echo "scale=6; $reward_amount / 10^$decimals" | bc)
        delegate_amount=$(echo "scale=6; $delegate_amount / 10^$decimals" | bc)
        echo "$oracle_address#$delegate_amount#$start_height#$online#$delegate_validator#$reward_amount"
      done < <(curl -s "$REST_RPC/fx/crosschain/v1/oracles?chain_name=$chain_name" | jq -r '.oracles[]|"\(.oracle_address) \(.delegate_amount) \(.start_height) \(.online) \(.delegate_validator)"')
    done
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
  local after_days=${1:-"365"}

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
  printf "        |-> ⚠️ blocks_per_year here is calculated based on 5s block time, the actual block time is not fixed, so this value is inaccurate\n"
  printf "\n"

  # query total supply
  supply=$(curl -s "$REST_RPC/cosmos/bank/v1beta1/supply")

  # query staking params
  bond_denom=$(curl -s "$REST_RPC/cosmos/staking/v1beta1/params" | jq -r '.params.bond_denom')
  # staking total supply
  staking_total_supply=$(echo "$supply" | jq -r ".supply[]|select(.denom == \"$bond_denom\")|.amount")
  printf "current staking total supply:\t\t%s\n" "$(echo "scale=2;$staking_total_supply/10^18" | bc)"

  # query staking pool
  staking_pool=$(curl -s "$REST_RPC/cosmos/staking/v1beta1/pool" | jq -r '.pool')

  # bonded tokens
  bonded_tokens=$(echo "$staking_pool" | jq -r '.bonded_tokens')
  printf "current staking bonded tokens:\t\t%s\n" "$(echo "scale=2;$bonded_tokens/10^18" | bc)"

  # bonded ratio = bonded_tokens / staking_total_supply
  bonded_ratio=$(printf "%.18f" "$(echo "scale=18;$bonded_tokens/$staking_total_supply" | bc)")
  printf "current staking bonded ratio:\t\t%s\n" "$(printf "%.6f" "$(echo "scale=6;$bonded_ratio*100/1" | bc)")%"
  printf "    |-> bonded_ratio = bonded_tokens / staking_total_supply\n"
  printf "\n"

  total_supply=$staking_total_supply
  if [ "$mint_denom" != "$bond_denom" ]; then
    # mint denom total supply
    mint_total_supply=$(echo "$supply" | jq -r ".supply[]|select(.denom == \"$mint_denom\")|.amount")
    printf "current mint total supply:\t\t%s\n" "$(echo "scale=2;$mint_total_supply/10^18" | bc)"
    total_supply=$mint_total_supply
  fi

  # inflation_per_block = (1 - bonded_ratio/goal_bonded) * inflation_rate_change / blocks_per_year
  inflation_per_block=$(printf "%.18f" "$(echo "scale=18;(1 - $bonded_ratio/$goal_bonded) * $inflation_rate_change / $blocks_per_year" | bc)")
  printf "current mint inflation per block:\t%s\n" "$(printf "%.6f" "$(echo "scale=6;$inflation_per_block*100/1" | bc)")%"
  printf "    |-> inflation_per_block = (1 - bonded_ratio/goal_bonded) * inflation_rate_change / blocks_per_year\n"

  # query mint inflation
  inflation=$(curl -s "$REST_RPC/cosmos/mint/v1beta1/inflation" | jq -r '.inflation')
  printf "current mint inflation:\t\t\t%s\n" "$(echo "scale=6;$inflation*100/1" | bc)%"
  printf "    |-> inflation = inflation + inflation_per_block\n"

  # query annual provisions
  annual_provisions=$(curl -s "$REST_RPC/cosmos/mint/v1beta1/annual_provisions" | jq -r '.annual_provisions')
  printf "current mint annual provisions:\t\t%s\n" "$(echo "scale=2;$annual_provisions/10^18" | bc)"
  printf "    |-> annual_provisions = inflation * staking_total_supply\n"

  # average inflation per block = annual_provisions / blocks_per_year
  average_inflation_per_block=$(echo "scale=2;$annual_provisions/$blocks_per_year" | bc)
  printf "average inflation per block:\t\t%s\n" "$(echo "scale=2;$average_inflation_per_block/10^18" | bc)"
  printf "    |-> average_inflation_per_block = annual_provisions / blocks_per_year\n"
  printf "\n"

  check_command python3
  printf "expected after %s days: (only for reference, not accurate)\n" "$after_days"
  real_block_time_interval=$(avg_block_time_interval)
  python3 <<EOF
inflation = $inflation
inflation_per_block = $inflation_per_block
total_supply = $total_supply
real_avg_block_per_day = int(24*3600/$real_block_time_interval)
inflation_history = []
total_supply_history = []
per=1
if ($after_days >= 30) : per=10
if ($after_days >= 180) : per=30
print("\t\t\ttotal supply \t\tinflation")

for i in range($after_days*real_avg_block_per_day+1):
    inflation = inflation + inflation_per_block
    annual_provisions = inflation * $staking_total_supply
    total_supply = total_supply + (annual_provisions/(365*real_avg_block_per_day))
    if i != 0 and (i/real_avg_block_per_day) % per == 0:
        inflation_history.append(inflation)
        total_supply_history.append(total_supply)
        continue
    if $after_days*real_avg_block_per_day == i:
        inflation_history.append(inflation)
        total_supply_history.append(total_supply)

last = len(total_supply_history)-1
for i in range(last):
    print("    after %3d days\t%.6f\t%.2f%%" % ((i+1)*per, total_supply_history[i]/10**18, inflation_history[i]*100))

print("    after %3d days\t%.6f\t%.2f%%" % ($after_days, total_supply_history[last]/10**18, inflation_history[last]*100))
EOF
  printf "    |-> ⚠️ in blockchain, inflation will not decrease after inflation < min_inflation_rate\n"
  printf "    |-> ⚠️ in blockchain, inflation will not increase after inflation > max_inflation_rate\n"
}

function help() {
  printf "This script is used to stats cosmos chain info.\n"
  printf "Usage:\n"
  printf "    %s <command> [args]\n" "$0"
  printf "The commands are:\n"
  {
    printf "    show_validator_reward # # show validator reward\n"
    printf "    show_validator_vote # # show validator vote\n"
    printf "    show_module_account # # show all module account\n"
    printf "    show_mint_info # [<after_days>] # show mint module info\n"
    printf "    help # # show this help message\n"
    printf "\nVersion: Alpha\n"
  } | column -t -s "#"
}

[[ "$#" -eq 0 || "$1" == "help" ]] && help && exit 0
[[ "$#" -gt 0 && "$(type -t "$1")" != "function" ]] && echo "invalid args: $1" && exit 1
"$@" || (echo "failed: $0" "$@" && exit 1)
