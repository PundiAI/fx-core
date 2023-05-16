#!/usr/bin/env bash

set -euo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

readonly script_dir="${PROJECT_DIR}/tests/scripts"
readonly contract_address="0x0000000000000000000000000000000000001003"

export NODE_HOME="$OUT_DIR/.fxcore"
export GAS_ADJUSTMENT=1.4

# Import external ethereum.sh files for contract interaction
function evm_interact() {
  "${script_dir}/ethereum.sh" "$@"
}

# Handle input validation and provide clear error messages
function validate_input() {
  if [[ -z "$1" ]]; then
    echo "Error: Missing required input. Exiting..." >&2
    exit 1
  fi
}

# Interact with EVM
function evm_command() {
  local action="$1"
  local contract="$2"
  local method="$3"
  local arguments=("${@:4}")

  validate_input "$action"
  validate_input "$contract"
  validate_input "$method"

  # Construct the command
  local cmd=(evm_interact "$action" "$contract" "$method")

  # Append additional arguments if provided
  if ((${#arguments[@]} > 0)); then
    cmd+=("${arguments[@]}")
  fi

  # Execute the command and handle potential errors
  if ! output=$("${cmd[@]}" 2>&1); then
    echo "Error: Failed to execute evm interact command. Exiting..." >&2
    echo "Command: ${cmd[*]}" >&2
    echo "Output: $output" >&2
    exit 1
  fi

  echo "$output"
}

# ARGS: [<contract_address>] [<validator_address>] [<delegation_value>] [<path-index>]
# DESC: Delegate to a validator
function send_delegate() {
  local contract="${1:-$contract_address}"
  local validator_address="$2"
  local delegation_value="$3"
  local path_index="${4:-0}"

  validate_input "$validator_address"
  validate_input "$delegation_value"

  evm_command send "$contract" "delegate(string)" "$validator_address" --value "$delegation_value" --index "$path_index" --disable-confirm true
}

# ARGS: [<contract_address>] [<validator_address>] [<path-index>]
# DESC: Withdraw rewards from a validator
function send_withdraw() {
  local contract="${1:-$contract_address}"
  local validator_address="$2"
  local path_index="${3:-0}"

  validate_input "$validator_address"

  evm_command send "$contract" "withdraw(string)" "$validator_address" --index "$path_index" --disable-confirm true
}

# ARGS: [<contract_address>] [<validator_address>] [<shares>] <path-index>
# DESC: Undelegate from a validator
function send_undelegate() {
  local contract="${1:-$contract_address}"
  local validator_address="$2"
  local shares="$3"
  local path_index="${4:-0}"

  validate_input "$validator_address"
  validate_input "$shares"

  evm_command send "$contract" "undelegate(string,uint256)" "$validator_address" "$shares" --index "$path_index" --disable-confirm true
}

# ARGS: [<contract_address>] [<validator_address>] [<receipt_address>] [<shares>] <path-index>
# DESC: Transfer shares to a receipt
function send_transfer_shares() {
  local contract="${1:-$contract_address}"
  local validator_address="$2"
  local receipt_address="$3"
  local shares="$4"
  local path_index="${5:-0}"

  validate_input "$validator_address"
  validate_input "$receipt_address"
  validate_input "$shares"

  evm_command send "$contract" "transferShares(string,address,uint256)" "$validator_address" "$receipt_address" "$shares" --index "$path_index" --disable-confirm true
}

# ARGS: [<contract_address>] [<validator_address>] [<owner>] [<spender>] [<shares>] <path-index>
# DESC: Approve shares to a spender
function send_approve_shares() {
  local contract="${1:-$contract_address}"
  local validator_address="$2"
  local spender="$3"
  local shares="$4"
  local path_index="${5:-0}"

  validate_input "$validator_address"
  validate_input "$spender"
  validate_input "$shares"

  evm_command send "$contract" "approveShares(string,address,uint256)" "$validator_address" "$spender" "$shares" --index "$path_index" --disable-confirm true
}

# ARGS: [<contract_address>] [<validator_address>] [<owner>] [<spender>] [<shares>] <path-index>
# DESC: Transfer from shares
function send_transfer_from_shares() {
  local contract="${1:-$contract_address}"
  local validator_address="$2"
  local owner="$3"
  local spender="$4"
  local shares="$5"
  local path_index="${6:-0}"

  validate_input "$validator_address"
  validate_input "$owner"
  validate_input "$spender"
  validate_input "$shares"

  evm_command send "$contract" "transferFromShares(string,address,address,uint256)" "$validator_address" "$owner" "$spender" "$shares" --index "$path_index" --disable-confirm true
}

# ARGS: [<contract_address>] [<validator_address>] [<delegate_address>]
# DESC: Query delegation information
function get_delegation_info() {
  local contract="${1:-$contract_address}"
  local validator_address="$2"
  local delegate_address="$3"

  validate_input "$validator_address"
  validate_input "$delegate_address"

  evm_command call "$contract" "delegation(string,address)(uint256,uint256)" "$validator_address" "$delegate_address"
}

# ARGS: [<contract_address>] [<validator_address>] [<delegate_address>]
# DESC: Query delegation rewards
function get_delegation_rewards() {
  local contract="${1:-$contract_address}"
  local validator_address="$2"
  local delegate_address="$3"

  validate_input "$validator_address"
  validate_input "$delegate_address"

  evm_command call "$contract" "delegationRewards(string,address)(uint256)" "$validator_address" "$delegate_address"
}

# ARGS: [<contract_address>] [<validator_address>] [<owner>] [<spender>]
# DESC: Get allowance shares
function get_allowance_shares() {
  local contract="${1:-$contract_address}"
  local validator_address="$2"
  local owner="$3"
  local spender="$4"

  validate_input "$validator_address"
  validate_input "$owner"
  validate_input "$spender"

  evm_command call "$contract" "allowanceShares(string,address,address)(uint256)" "$validator_address" "$owner" "$spender"
}

# use contract <true|false>
# DESC: address or contract for delegation
function staking_delegate_test() {
  local use_contract="$1"

  validate_input "$use_contract"

  local validators_list
  validators_list=$(validators_list)
  local validator_address
  validator_address=$(echo "${validators_list}" | jq -r '.[0]')
  local del_address
  del_address=$(show_address "${FROM}" "-e")
  local contract
  if [[ "$use_contract" == "true" ]]; then
    echo "use contract delegate"
    contract=$(evm_interact deploy_staking_contract)
    del_address="$contract"
  else
    contract="$contract_address"
  fi

  send_delegate "$contract" "$validator_address" "$(to_18 "10^2")"

  get_delegation_info "$contract" "$validator_address" "$del_address"

  get_delegation_rewards "$contract_address" "$validator_address" "$del_address"

  send_withdraw "$contract" "$validator_address"

  send_undelegate "$contract" "$validator_address" "$(to_18 "10^2")"

  get_delegation_info "$contract" "$validator_address" "$del_address"
}

# use contract <true|false>
# DESC: Address or contract for delegated transfer shares
function staking_shares_test() {
  local use_contract="$1"

  validate_input "$use_contract"

  local path_index=1
  local validators_list
  validators_list=$(validators_list)
  local validator_address
  validator_address=$(echo "${validators_list}" | jq -r '.[0]')
  local del_address
  del_address=$(show_address "$FROM" "-e")
  local contract
  local receipt_address
  if [[ "$use_contract" == "true" ]]; then
    echo "use contract delegate"
    contract=$(evm_interact deploy_staking_contract)
    receipt_address="$del_address"
    del_address="$contract"
    path_index=0
  else
    contract="$contract_address"
    local receipt_name="receipt"
    add_key "$receipt_name" "1" | jq -r ".address"
    receipt_address=$(show_address "$receipt_name" "-e")
    cosmos_transfer "$receipt_name" 10
  fi
  send_delegate "$contract" "$validator_address" "$(to_18 "200")"
  get_delegation_info "$contract" "$validator_address" "$del_address"
  send_transfer_shares "$contract" "$validator_address" "$receipt_address" "$(to_18 "200")"
  get_delegation_info "$contract" "$validator_address" "$del_address"
  get_delegation_info "$contract" "$validator_address" "$receipt_address"
  send_approve_shares "$contract_address" "$validator_address" "$del_address" "$(to_18 "200")" "$path_index"
  get_allowance_shares "$contract_address" "$validator_address" "$receipt_address" "$del_address"
  send_transfer_from_shares "$contract" "$validator_address" "$receipt_address" "$del_address" "$(to_18 "200")"
  get_delegation_info "$contract" "$validator_address" "$del_address"
  get_delegation_info "$contract" "$validator_address" "$receipt_address"
  get_allowance_shares "$contract_address" "$validator_address" "$receipt_address" "$del_address"
}

# DESC: Start node
function start_node() {
  "${script_dir}/fxcore.sh" init
  "${script_dir}/fxcore.sh" start

  sleep 10
}

# DESC: Perform staking operations
function staking_operations() {
  ## start node
  start_node

  ## start test
  staking_delegate_test false
  staking_delegate_test true
  staking_shares_test false
  staking_shares_test true
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
