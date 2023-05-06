#!/usr/bin/env bash

set -eo pipefail

readonly project_dir="$(git rev-parse --show-toplevel)"
readonly out_dir="${project_dir}/out"
readonly data_dir="${project_dir}/tests/data"
readonly proposal_dir="${data_dir}/proposal"

## ARGS: <data_file> <metadata>
## DESC: submit a proposal
function submit_proposal() {
  local data_file="$1"
  local metadata="$2"

  cosmos_tx gov submit-proposal "$data_file" --from fx1 -y

  for proposal_id in $(cosmos_query gov proposals | jq -r '.proposals[].id'); do
    messages=$(cosmos_query gov proposal "${proposal_id}")
    if [[ "$(echo "${messages}" | jq -r '.messages[].metadata')" == "$metadata" &&
    "$(echo "${messages}" | jq -r '.status')" == "PROPOSAL_STATUS_VOTING_PERIOD" ]]; then
      break
    fi
  done

  cosmos_tx gov vote "${proposal_id}" yes --from fx1 -y
}

## ARGS: <title> <summary>
## DESC: base64 encode metadata
function base64_metadata() {
  echo '{"title": "'"$1"'","summary": "'"$2"'","metadata":""}' | base64
}

## ARGS: <msg_type> <amount>
## DESC: query min deposit
function query_min_deposit() {
  local msg_type=$1
  local amount=$2
  if [[ -z "$msg_type" ]]; then
    echo "$(cosmos_query gov params | jq -r '.params.min_deposit[0].amount')FX"
    return
  fi
  baseDeposit="$(cosmos_query gov params --msg-type="$msg_type" | jq -r '.params.min_deposit[0].amount')"
  if [[ "$msg_type" == "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal" && -n "$amount" ]]; then
    depositThreshold=$(cosmos_query gov egf-params | jq -r '.params.egf_deposit_threshold.amount')
    claimRatio=$(cosmos_query gov egf-params | jq -r '.params.claim_ratio')

    amount_without=${amount/FX/}
    flag=$(echo "$amount_without" - "$depositThreshold" | bc)
    if [[ flag -gt 0 ]]; then
      baseDeposit=$(echo "$amount_without" \* "$claimRatio" | bc)
    fi
  fi
  echo "$baseDeposit"FX
}

## DESC: submit register coin proposal
function register_coin_proposal() {

  metadata=$(base64_metadata "register coin" "This proposal creates and registers an ERC20 representation for the coin")

  deposit=$(query_min_deposit "/fx.erc20.v1.MsgRegisterCoin")

  cp "${proposal_dir}/register_coin.json" "${out_dir}/register_coin.json"
  json_processor "${out_dir}/register_coin.json" '.metadata = "'"$metadata"'"'
  json_processor "${out_dir}/register_coin.json" '.deposit = "'"$deposit"'"'

  submit_proposal "${out_dir}"/register_coin.json "$metadata"
}

## DESC: submit register erc20 proposal
function register_erc20_proposal() {

  metadata=$(base64_metadata "register fip20" "This proposal registers and creates a corresponding native coin for fip20 tokens")

  deposit=$(query_min_deposit "/fx.erc20.v1.MsgRegisterERC20")

  cp "${proposal_dir}/register_erc20.json" "${out_dir}/register_erc20.json"
  json_processor "${out_dir}/register_erc20.json" '.metadata = "'"$metadata"'"'
  json_processor "${out_dir}/register_erc20.json" '.deposit = "'"$deposit"'"'

  for erc20address in $(jq -r '.messages[].erc20address' "${out_dir}/register_erc20.json"); do
    if [[ $(cosmos_query evm code "${erc20address}" | jq -r '.code') == null ]]; then
      echo "contract $erc20address not deployed, cannot initiate MsgRegisterERC20 proposal" && exit 1
    fi
  done

  submit_proposal "${out_dir}/register_erc20.json" "$metadata"
}

## DESC: submit update denom alias proposal
function update_denom_alias_proposal() {

  metadata=$(base64_metadata "update denom alias" "This proposal allows modifying the alias of the bridge corresponding to the coin")

  deposit=$(query_min_deposit "/fx.erc20.v1.MsgUpdateDenomAlias")

  cp "${proposal_dir}/update_denom_alias.json" "${out_dir}/update_denom_alias.json"
  json_processor "${out_dir}/update_denom_alias.json" '.metadata = "'"$metadata"'"'
  json_processor "${out_dir}/update_denom_alias.json" '.deposit = "'"$deposit"'"'

  for denom in $(jq -r '.messages[].denom' "${out_dir}/update_denom_alias.json"); do
    cosmos_query erc20 token-pair "${denom}"
  done

  submit_proposal "${out_dir}/update_denom_alias.json" "$metadata"
}

## DESC: submit toggle token conversion proposal
function toggle_token_conversion_proposal() {

  metadata=$(base64_metadata "toggle token conversion" "This proposal is used to enable or disable the conversion of coins and tokens")

  deposit=$(query_min_deposit "/fx.erc20.v1.MsgToggleTokenConversion")

  cp "${proposal_dir}/toggle_token_conversion.json" "${out_dir}/toggle_token_conversion.json"
  json_processor "${out_dir}/toggle_token_conversion.json" '.metadata = "'"$metadata"'"'
  json_processor "${out_dir}/toggle_token_conversion.json" '.deposit = "'"$deposit"'"'

  for token in $(jq -r '.messages[].token' "${out_dir}/toggle_token_conversion.json"); do
    cosmos_query erc20 token-pair "${token}"
  done

  submit_proposal "${out_dir}/toggle_token_conversion.json" "$metadata"
}

## DESC: submit update crosschain params proposal
function update_crosschain_params_proposal() {

  metadata=$(base64_metadata "update crosschain Params" "This proposal is used to modify the variable parameters of the cross-chain")

  deposit=$(query_min_deposit "/fx.gravity.crosschain.v1.MsgUpdateParams")

  cp "${proposal_dir}/update_crosschain_params.json" "${out_dir}/update_crosschain_params.json"
  json_processor "${out_dir}/update_crosschain_params.json" '.metadata = "'"$metadata"'"'
  json_processor "${out_dir}/update_crosschain_params.json" '.deposit = "'"$deposit"'"'

  submit_proposal "${out_dir}/update_crosschain_params.json" "$metadata"
}

## DESC: submit update ecc20 params proposal
function update_ecc20_params_proposal() {

  metadata=$(base64_metadata "update erc20 Params" "This proposal is used to modify the variable parameters of the erc20 module")

  deposit=$(query_min_deposit "/fx.erc20.v1.MsgUpdateParams")

  cp "${proposal_dir}/update_erc20_params.json" "${out_dir}/update_erc20_params.json"
  json_processor "${out_dir}/update_erc20_params.json" '.metadata = "'"$metadata"'"'
  json_processor "${out_dir}/update_erc20_params.json" '.deposit = "'"$deposit"'"'

  submit_proposal "${out_dir}/update_erc20_params.json" "$metadata"
}

## DESC: submit update gov params proposal
function update_gov_params_proposal() {

  metadata=$(base64_metadata "update gov Params" "This proposal is used to modify the variable parameters of the gov module")

  deposit=$(query_min_deposit "/fx.gov.v1.MsgUpdateParams")

  cp "${proposal_dir}/update_gov_params.json" "${out_dir}/update_gov_params.json"
  json_processor "${out_dir}/update_gov_params.json" '.metadata = "'"$metadata"'"'
  json_processor "${out_dir}/update_gov_params.json" '.deposit = "'"$deposit"'"'

  submit_proposal "${out_dir}/update_gov_params.json" "$metadata"
}

## DESC: submit call contract proposal
function call_contract_proposal() {

  metadata=$(base64_metadata "evm call contract" "This proposal is used to call the method of the evm contract")

  deposit=$(query_min_deposit "/fx.evm.v1.MsgCallContract")

  cp "${proposal_dir}/call_contract.json" "${out_dir}/call_contract.json"
  json_processor "${out_dir}/call_contract.json" '.metadata = "'"$metadata"'"'
  json_processor "${out_dir}/call_contract.json" '.deposit = "'"$deposit"'"'

  submit_proposal "${out_dir}/call_contract.json" "$metadata"
}

## DESC: submit distribution community pool spend proposal
function distribution_community_pool_spend_proposal() {

  title=$(jq -r '.title' "${proposal_dir}/leagcy_community_pool_spend.json")
  summary=$(jq -r '.description' "${proposal_dir}/leagcy_community_pool_spend.json")

  metadata=$(base64_metadata "$title" "$summary")

  amount=$(jq -r '.amount' "${proposal_dir}/leagcy_community_pool_spend.json")

  deposit=$(query_min_deposit "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal" "$amount")

  cp "${proposal_dir}/leagcy_community_pool_spend.json" "${out_dir}/leagcy_community_pool_spend.json"
  json_processor "${out_dir}/leagcy_community_pool_spend.json" '.metadata = "'"$metadata"'"'
  json_processor "${out_dir}/leagcy_community_pool_spend.json" '.deposit = "'"$deposit"'"'

  cosmos_tx gov submit-legacy-proposal community-pool-spend "${out_dir}/leagcy_community_pool_spend.json" --from fx1 -y

  for proposal_id in $(cosmos_query gov proposals | jq -r '.proposals[].id'); do
    messages=$(cosmos_query gov proposal "${proposal_id}")
    if [[ "$(echo "${messages}" | jq -r '.messages[].metadata')" == "$metadata" &&
    "$(echo "${messages}" | jq -r '.status')" == "PROPOSAL_STATUS_VOTING_PERIOD" ]]; then
      break
    fi
  done

  cosmos_tx gov vote "${proposal_id}" yes --from fx1 -y
}

# shellcheck source=/dev/null
. "${project_dir}/tests/scripts/setup-env.sh"

check_command jq "$DAEMON"
