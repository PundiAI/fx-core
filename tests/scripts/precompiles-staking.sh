#!/usr/bin/env bash

set -eo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

readonly script_dir="${PROJECT_DIR}/tests/scripts"
readonly contract_address="0x0000000000000000000000000000000000001003"

export NODE_HOME="$OUT_DIR/.fxcore"
export GAS_ADJUSTMENT=1.4

## ARGS: <path-index> <contract_address> <input:[validator_address]> <value>
## DESC: delegate
function delegate() {
  "${script_dir}/contract.sh" send "$1" "$2" "delegate(string)" "$3" --value "$4" --disable-confirm true
}

## ARGS: <path-index> <contract_address> <input:[validator_address]>
## DESC: withdraw rewards
function withdraw() {
  "${script_dir}/contract.sh" send "$1" "$2" "withdraw(string)" "$3" --disable-confirm true
}

## ARGS: <path-index> <contract_address> <input:[validator_address]> <input:[shares]>
## DESC: undelegate
function undelegate() {
  "${script_dir}/contract.sh" send "$1" "$2" "undelegate(string,uint256)" "$3" "$4" --disable-confirm true
}

## ARGS: <contract_address> <input:[validator_address]> <input:[delegate_address]>
## DESC: query delegation
function delegation() {
  "${script_dir}/contract.sh" call "$1" "delegation(string,address)(uint256,uint256)" "$2" "$3"
}

## ARGS: <contract_address> <input:[validator_address]> <input:[delegate_address]>
## DESC: query delegation rewards
function delegation_rewards() {
  "${script_dir}/contract.sh" call "$1" "delegationRewards(string,address)(uint256)" "$2" "$3"
}

## ARGS: <path-index> <contract_address> <input:[validator_address]> <input:[address]> <input:[shares]>
## DESC: contract transfer shares
function transfer_shares() {
  "${script_dir}/contract.sh" send "$1" "$2" "transferShares(string,address,uint256)" "$3" "$4" "$5" --disable-confirm true
}

## ARGS: <path-index> <contract_address> <input:[validator_address]> <input:[address]> <input:[address]> <input:[shares]>
## DESC: contract transfer from shares
function transfer_from_shares() {
  "${script_dir}/contract.sh" send "$1" "$2" "transferFromShares(string,address,address,uint256)" "$3" "$4" "$5" "$6" --disable-confirm true
}

## ARGS: <path-index> <contract_address> <input:[validator_address]> <input:[address]> <input:[shares]>
## DESC: contract approve shares
function approve_shares() {
  "${script_dir}/contract.sh" send "$1" "$2" "approveShares(string,address,uint256)" "$3" "$4" "$5" --disable-confirm true
}

## ARGS: <contract_address> <input:[validator_address]> <input:[address]> <input:[address]>
## DESC: contract allowance shares
function allowance_shares() {
  "${script_dir}/contract.sh" call "$1" "allowanceShares(string,address,address)(uint256)" "$2" "$3" "$4"
}

# DESC: precompiled contract method calls
function start() {
  "${script_dir}/fxcore.sh" init
  "${script_dir}/fxcore.sh" start

  sleep 10

  staking_delegate fale

  staking_delegate true

  staking_shares

  ## todo add contract transfer shares
}

## ARGS: <true/false>
## DESC: delegated validator by address or contract
function staking_delegate() {
  ## set HDPath index
  index=0
  ## get validators list
  validators_list=$(validators_list)
  ## get first validator
  validator_address_0=$(echo "${validators_list}" | jq -r '.[0]')
  # set delegator address
  del_address=$(show_address "$FROM" "-e")
  # set withdraw address
  withdraw_address=$(show_address "$FROM" "-a")
  contract=${contract_address}

  ## contract delegate
  if [ "$1" = true ]; then
    ## deploy contract
    contract=$("${script_dir}/contract.sh" deploy_staking_contract)
    del_address=${contract}
    withdraw_address=$($DAEMON debug addr "${contract}" | jq -r ".bech32")
  fi

  txHash=$(delegate "${index}" "${contract}" "${validator_address_0}" "$(to_18 "10^2")")
  echo "==> ${del_address} delegate to ${validator_address_0} success. tx hash: ${txHash}"
  ## call the contract query delegate
  output=$(delegation "${contract}" "${validator_address_0}" "${del_address}")
  amount="$(awk 'NR==1' <(echo "$output"))"
  shares="$(awk 'NR==2' <(echo "$output"))"

  echo "==> ${del_address} delegation amount: ${amount}  shares: ${shares}"

  ## call the contract query delegation rewards
  rewards=$(delegation_rewards "${contract_address}" "${validator_address_0}" "${del_address}")
  echo "==> ${del_address} delegate rewards: ${rewards}"
  echo "==> Before withdraw rewards withdraw addr: ${withdraw_address} balance: $(cosmos_query bank balances "${withdraw_address}" --denom="${STAKING_DENOM}" | jq -r ".amount")"
  ## withdraw rewards
  txHash=$(withdraw "${index}" "${contract}" "${validator_address_0}")
  echo "==> ${del_address} withdraw rewards success. tx hash: ${txHash}"
  echo "==> After withdraw rewards withdraw addr: ${withdraw_address} balance: $(cosmos_query bank balances "${withdraw_address}" --denom="${STAKING_DENOM}" | jq -r ".amount")"
  ## undelegate
  txHash=$(undelegate "${index}" "${contract}" "${validator_address_0}" "$(to_18 "10^2")")
  echo "==> ${del_address} undelegate success. tx hash: ${txHash}"
  ## call the contract query delegate after undelegate
  output=$(delegation "${contract}" "${validator_address_0}" "${del_address}")
  amount="$(awk 'NR==1' <(echo "$output"))"
  shares="$(awk 'NR==2' <(echo "$output"))"
  echo "==> ${del_address} delegation after undelegation amount: ${amount}  shares: ${shares}"
}

function staking_shares() {
  ## set HDPath index
  index_0=0
  index_1=1
  ## get validators list
  validators_list=$(validators_list)
  ## get first validator
  validator_address_0=$(echo "${validators_list}" | jq -r '.[0]')
  # set delegator address
  del_address=$(show_address "$FROM" "-e")

  contract=${contract_address}
  # receipt address
  receipt_name="receipt"
  add_key "${receipt_name}" "1" | jq -r ".address"
  receipt_address=$(show_address "$receipt_name" "-e")

  cosmos_transfer "${receipt_name}" 10

  ## delegate
  txHash=$(delegate "${index_0}" "${contract}" "${validator_address_0}" "$(to_18 "200")")
  echo "==> ${del_address} delegate to ${validator_address_0} success. tx hash: ${txHash}"

  ## call the contract query delegate
  output=$(delegation "${contract}" "${validator_address_0}" "${del_address}")
  amount="$(awk 'NR==1' <(echo "$output"))"
  shares="$(awk 'NR==2' <(echo "$output"))"
  echo "==> ${del_address} delegation amount: ${amount}  shares: ${shares}"

  ## transfer shares
  txHash=$(transfer_shares "${index_0}" "${contract}" "${validator_address_0}" "${receipt_address}" "${shares}")
  echo "==> ${del_address} transfer shares to ${receipt_address} success. tx hash: ${txHash}"

  ## query del_address delegate info
  output=$(delegation "${contract}" "${validator_address_0}" "${del_address}")
  amount="$(awk 'NR==1' <(echo "$output"))"
  shares="$(awk 'NR==2' <(echo "$output"))"
  echo "==> ${del_address} delegation amount: ${amount}  shares: ${shares}"

  ## query receipt_address delegate info
  output=$(delegation "${contract}" "${validator_address_0}" "${receipt_address}")
  amount="$(awk 'NR==1' <(echo "$output"))"
  shares="$(awk 'NR==2' <(echo "$output"))"
  echo "==> ${receipt_address} delegation amount: ${amount}  shares: ${shares}"

  ## approve receipt_address shares
  # shellcheck disable=SC2154
  txHash=$(approve_shares "${index_1}" "${contract}" "${validator_address_0}" "${del_address}" "${shares}")
  echo "==> ${receipt_address} approve shares to ${del_address} success. tx hash: ${txHash}"

  ## query del_address allowance shares
  allowanceShares=$(allowance_shares "${contract}" "${validator_address_0}" "${receipt_address}" "${del_address}")
  echo "==> ${receipt_address} allowance shares to ${del_address}: ${allowanceShares}"

  ## del_address transfer from shares
  txHash=$(transfer_from_shares "${index_0}" "${contract}" "${validator_address_0}" "${receipt_address}" "${del_address}" "${shares}")
  echo "==> ${del_address} transfer from shares to ${validator_address_0} success. tx hash: ${txHash}"

  ## query del_address delegate info
  output=$(delegation "${contract}" "${validator_address_0}" "${del_address}")
  amount="$(awk 'NR==1' <(echo "$output"))"
  shares="$(awk 'NR==2' <(echo "$output"))"
  echo "==> ${del_address} delegation amount: ${amount}  shares: ${shares}"

  ## query receipt_address delegate info
  output=$(delegation "${contract}" "${validator_address_0}" "${receipt_address}")
  amount="$(awk 'NR==1' <(echo "$output"))"
  shares="$(awk 'NR==2' <(echo "$output"))"
  echo "==> ${receipt_address} delegation amount: ${amount}  shares: ${shares}"

  ## query del_address allowance shares
  allowanceShares=$(allowance_shares "${contract}" "${validator_address_0}" "${receipt_address}" "${del_address}")
  echo "==> ${receipt_address} allowance shares to ${del_address}: ${allowanceShares}"
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
