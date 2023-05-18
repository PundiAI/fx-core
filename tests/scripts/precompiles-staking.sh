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

  evm_command send "$contract" "delegate(string)" "$validator_address" --value "$delegation_value" --index "$path_index"
}

# ARGS: [<contract_address>] [<validator_address>] [<path-index>]
# DESC: Withdraw rewards from a validator
function send_withdraw() {
  local contract="${1:-$contract_address}"
  local validator_address="$2"
  local path_index="${3:-0}"

  validate_input "$validator_address"

  evm_command send "$contract" "withdraw(string)" "$validator_address" --index "$path_index"
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

  evm_command send "$contract" "undelegate(string,uint256)" "$validator_address" "$shares" --index "$path_index"
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

  evm_command send "$contract" "transferShares(string,address,uint256)" "$validator_address" "$receipt_address" "$shares" --index "$path_index"
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

  evm_command send "$contract" "approveShares(string,address,uint256)" "$validator_address" "$spender" "$shares" --index "$path_index"
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

  evm_command send "$contract" "transferFromShares(string,address,address,uint256)" "$validator_address" "$owner" "$spender" "$shares" --index "$path_index"
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

# Import assertion file
function assert() {
  "${script_dir}/assert.sh" "$@"
}

# use contract <true|false>
# DESC: address or contract for delegation
function staking_delegate_test() {

  assert log_header "Test assert : staking_delegate_test"

  local use_contract="$1"

  validate_input "$use_contract"

  local validators_list
  validators_list=$(validators_list)
  local validator_address
  validator_address=$(echo "${validators_list}" | jq -r '.[0]')
  local del_address
  del_address=$(show_address "${FROM}" "-e")
  local delegate_amount
  delegate_amount=$(to_18 "10^2")

  local contract
  if [[ "$use_contract" == "true" ]]; then
    echo "use contract delegate"
    contract=$(evm_interact deploy_contract StakingTest)
    del_address="$contract"
  else
    contract="$contract_address"
  fi

  assert assert_not_empty "$(send_delegate "$contract" "$validator_address" "$delegate_amount")" "send delegate not empty"

  delegation_info=$(get_delegation_info "$contract" "$validator_address" "$del_address")
  assert assert_ge "$(awk 'NR==1' <(echo "$delegation_info"))" "$delegate_amount" "delegation amount greater than or equal to delegate amount"

  rewards_amount="$(get_delegation_rewards "$contract_address" "$validator_address" "$del_address")"
  assert assert_gt "${rewards_amount}" 0 "delegation reward is greater than 0"

  assert assert_not_empty "$(send_withdraw "$contract" "$validator_address")" "send withdraw not empty"

  assert assert_not_empty "$(send_undelegate "$contract" "$validator_address" "$delegate_amount")" "send undelegate not empty"

  undelegate_delegation_info=$(get_delegation_info "$contract" "$validator_address" "$del_address")
  undelegate_amount=$(echo "$(awk 'NR==1' <(echo "$delegation_info"))" - "$(awk 'NR==1' <(echo "$undelegate_delegation_info"))" | bc)

  assert assert_eq "$delegate_amount" "$undelegate_amount" "unbound amount is equal"
}

# use contract <true|false>
# DESC: Address or contract for delegated transfer shares
function staking_shares_test() {

  assert log_header "Test assert : staking_shares_test"

  local use_contract="$1"

  validate_input "$use_contract"

  local path_index=1
  local validators_list
  validators_list=$(validators_list)
  local validator_address
  validator_address=$(echo "${validators_list}" | jq -r '.[0]')
  local del_address
  del_address=$(show_address "$FROM" "-e")

  local delegate_amount
  delegate_amount=$(to_18 "10^2")

  local contract
  local receipt_address
  if [[ "$use_contract" == "true" ]]; then
    echo "use contract delegate"
    contract=$(evm_interact deploy_contract StakingTest)
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

  assert assert_not_empty "$(send_delegate "$contract" "$validator_address" "$delegate_amount")" "send delegate not empty"

  delegate_info_by_del_address=$(get_delegation_info "$contract" "$validator_address" "$del_address")
  assert assert_ge "$(awk 'NR==1' <(echo "$delegate_info_by_del_address"))" "$delegate_amount" "$del_address delegation amount greater than or equal to $delegate_amount"

  assert assert_not_empty "$(send_transfer_shares "$contract" "$validator_address" "$receipt_address" "$(awk 'NR==2' <(echo "$delegate_info_by_del_address"))")" "send transfer shares not empty"

  delegate_info_by_del_address_after_transfer_shares=$(get_delegation_info "$contract" "$validator_address" "$del_address")
  assert assert_eq "$(awk 'NR==1' <(echo "$delegate_info_by_del_address_after_transfer_shares"))" 0 "$del_address transfer shares after amount is zero"

  delegate_info_by_receipt_address_after_transfer_shares=$(get_delegation_info "$contract" "$validator_address" "$receipt_address")
  assert assert_ge "$(awk 'NR==1' <(echo "$delegate_info_by_receipt_address_after_transfer_shares"))" "$(awk 'NR==1' <(echo "$delegate_info_by_del_address"))" "$receipt_address delegation amount greater than or equal to $delegate_amount"

  assert assert_not_empty "$(send_approve_shares "$contract_address" "$validator_address" "$del_address" "$delegate_amount" "$path_index")" "$receipt_address approve shares not empty"

  allowance_shares=$(get_allowance_shares "$contract_address" "$validator_address" "$receipt_address" "$del_address")
  assert assert_gt "$allowance_shares" 0 "allowance shares greater than 0"

  assert assert_not_empty "$(send_transfer_from_shares "$contract" "$validator_address" "$receipt_address" "$del_address" "$delegate_amount")" "send transfer from shares not empty"

  allowance_shares=$(get_allowance_shares "$contract_address" "$validator_address" "$receipt_address" "$del_address")
  assert assert_eq "$allowance_shares" 0 "allowance shares equal to 0"

  delegate_info_by_del_address_after_transfer_from_shares=$(get_delegation_info "$contract" "$validator_address" "$del_address")
  delegate_info_by_del_address=$(echo "$(awk 'NR==1' <(echo "$delegate_info_by_del_address_after_transfer_from_shares"))" - "$(awk 'NR==1' <(echo "$delegate_info_by_del_address_after_transfer_shares"))" | bc)

  delegate_info_by_receipt_address_after_transfer_from_shares=$(get_delegation_info "$contract" "$validator_address" "$receipt_address")
  delegate_info_by_receipt_address=$(echo "$(awk 'NR==1' <(echo "$delegate_info_by_receipt_address_after_transfer_shares"))" - "$(awk 'NR==1' <(echo "$delegate_info_by_receipt_address_after_transfer_from_shares"))" | bc)

  assert assert_eq "$delegate_info_by_del_address" "$delegate_info_by_receipt_address" "transfer from shares amount successfully"
  assert assert_eq "$delegate_info_by_del_address" "$delegate_amount" "transfer from shares amount is equal to delegate amount"
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
