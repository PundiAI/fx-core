#!/usr/bin/env bash

set -eo pipefail

BASE_DIR=${PROJECT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]%/*}")" && pwd)}

# shellcheck source=/dev/null
source "$BASE_DIR"/scripts/setup_env.sh

DATA_DIR="$BASE_DIR"/data/proposal

help="Command:
    registerCoinProposal                                    submit register coin proposal
    registerERC20Proposal                                   submit register ERC20 token proposal   
    updateDenomAliasProposal                                submit update denom alias proposal
    toggleTokenConversionProposal                           submit toggle token conversion proposal
    updateCrossChainParamsProposal                          submit update crosschain params proposal 
    updateERC20ParamsProposal                               submit update ERC20 params proposal 
    updateGovParamsProposal                                 submit update GOV params proposal 
    callContractProposal                                    submit evm call contract proposal
    distributionCommunityPoolSpendProposal                  submit community pool spend proposal
"

submitProposal() {
  local data_file="$1"
  local metadata="$2"

  if ! fxcored tx gov submit-proposal "$data_file" --gas-prices="$GAS_PRICES" --gas=auto --gas-adjustment="$GAS_ADJUSTMENT" --from fx1 -y --node "$JSON_RPC" --chain-id "$CHAIN_ID"; then
    exit 1
  fi

  for proposal_id in $(fxcored query gov proposals --node "$JSON_RPC" -o json | jq -r '.proposals[].id'); do
    messages=$(fxcored query gov proposal "${proposal_id}" --node "$JSON_RPC" -o json)
    if [[ "$(echo "${messages}" | jq -r '.messages[].metadata')" == "$metadata" && "$(echo "${messages}" | jq -r '.status')" == "PROPOSAL_STATUS_VOTING_PERIOD" ]]; then
      break
    fi
  done

  fxcored tx gov vote "${proposal_id}" yes --gas-prices="$GAS_PRICES" --gas=auto --gas-adjustment=1.4 --from fx1 -y --node "$JSON_RPC"
}

function base64metadata() {
  echo '{"title": "'"$1"'","summary": "'"$2"'","metadata":""}' | base64
}

function queryMinDeposit() {
  if [[ -z $1 ]]; then
    echo "$(fxcored q gov params | jq -r '.params.min_deposit[0].amount')FX"
  else
    baseDeposit="$(fxcored q gov params --msg-type="$1" | jq -r '.params.min_deposit[0].amount')"
    if [[ "$1" == "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal" && -n $2 ]]; then
      depositThreshold=$(fxcored q gov egf-params | jq -r '.params.egf_deposit_threshold.amount')
      claimRatio=$(fxcored q gov egf-params | jq -r '.params.claim_ratio')

      amount="$2"
      amount_without_fx=${amount/FX/}

      flag=$(echo "$amount_without_fx" - "$depositThreshold" | bc)
      if [[ flag -gt 0 ]]; then
        baseDeposit=$(echo "$amount_without_fx" \* "$claimRatio" | bc)
      fi
    fi
    echo "$baseDeposit"FX
  fi
}

## register coin proposal
function registerCoinProposal() {

  metadata=$(base64metadata "register coin" "This proposal creates and registers an ERC20 representation for the coin")

  deposit=$(queryMinDeposit "/fx.erc20.v1.MsgRegisterCoin")

  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/register_coin.json >register_coin_tmp.json && mv register_coin_tmp.json "${DATA_DIR}"/register_coin.json
  jq '.deposit = "'"$deposit"'" ' "${DATA_DIR}"/register_coin.json >register_coin_tmp.json && mv register_coin_tmp.json "${DATA_DIR}"/register_coin.json

  submitProposal "${DATA_DIR}"/register_coin.json "$metadata"
}

## register erc20 proposal
function registerERC20Proposal() {

  metadata=$(base64metadata "register fip20" "This proposal registers and creates a corresponding native coin for fip20 tokens")

  deposit=$(queryMinDeposit "/fx.erc20.v1.MsgRegisterERC20")

  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/register_erc20.json >register_erc20_tmp.json && mv register_erc20_tmp.json "${DATA_DIR}"/register_erc20.json
  jq '.deposit = "'"$deposit"'" ' "${DATA_DIR}"/register_erc20.json >register_erc20_tmp.json && mv register_erc20_tmp.json "${DATA_DIR}"/register_erc20.json

  for erc20address in $(jq -r '.messages[].erc20address' "${DATA_DIR}"/register_erc20.json); do
    if [[ $(fxcored q evm code "${erc20address}" --node "$JSON_RPC" -o json | jq -r '.code') == null ]]; then
      echo "contract $erc20address not deployed, cannot initiate MsgRegisterERC20 proposal" && exit 1
    fi
  done

  submitProposal "${DATA_DIR}"/register_erc20.json "$metadata"
}

# update denom alias proposal
function updateDenomAliasProposal() {

  metadata=$(base64metadata "update denom alias" "This proposal allows modifying the alias of the bridge corresponding to the coin")

  deposit=$(queryMinDeposit "/fx.erc20.v1.MsgUpdateDenomAlias")

  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/update_denom_alias.json >update_denom_alias_tmp.json && mv update_denom_alias_tmp.json "${DATA_DIR}"/update_denom_alias.json
  jq '.deposit = "'"$deposit"'" ' "${DATA_DIR}"/update_denom_alias.json >update_denom_alias_tmp.json && mv update_denom_alias_tmp.json "${DATA_DIR}"/update_denom_alias.json

  for denom in $(jq -r '.messages[].denom' "${DATA_DIR}"/update_denom_alias.json); do
    if ! fxcored q erc20 token-pair "${denom}" --node "$JSON_RPC" >/dev/null 2>&1; then
      echo "token $denom unregistered, cannot initiate MsgUpdateDenomAlias proposal" && exit 1
    fi
  done

  submitProposal "${DATA_DIR}"/update_denom_alias.json "$metadata"
}

## toggle token conversion proposal
function toggleTokenConversionProposal() {

  metadata=$(base64metadata "toggle token conversion" "This proposal is used to enable or disable the conversion of coins and tokens")

  deposit=$(queryMinDeposit "/fx.erc20.v1.MsgToggleTokenConversion")

  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/toggle_token_conversion.json >toggle_token_conversion_tmp.json && mv toggle_token_conversion_tmp.json "${DATA_DIR}"/toggle_token_conversion.json
  jq '.deposit = "'"$deposit"'" ' "${DATA_DIR}"/toggle_token_conversion.json >toggle_token_conversion_tmp.json && mv toggle_token_conversion_tmp.json "${DATA_DIR}"/toggle_token_conversion.json

  for token in $(jq -r '.messages[].token' "${DATA_DIR}"/toggle_token_conversion.json); do
    if ! fxcored q erc20 token-pair "${token}" --node "$JSON_RPC" >/dev/null 2>&1; then
      echo "$token unregistered, cannot initiate MsgToggleTokenConversion proposal" && exit 1
    fi
  done

  submitProposal "${DATA_DIR}"/toggle_token_conversion.json "$metadata"
}

## update crosschain params
function updateCrossChainParamsProposal() {

  metadata=$(base64metadata "update crosschain Params" "This proposal is used to modify the variable parameters of the cross-chain")

  deposit=$(queryMinDeposit "/fx.gravity.crosschain.v1.MsgUpdateParams")

  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/update_crosschain_params.json >update_crosschain_params_tmp.json && mv update_crosschain_params_tmp.json "${DATA_DIR}"/update_crosschain_params.json
  jq '.deposit = "'"$deposit"'" ' "${DATA_DIR}"/update_crosschain_params.json >update_crosschain_params_tmp.json && mv update_crosschain_params_tmp.json "${DATA_DIR}"/update_crosschain_params.json

  submitProposal "${DATA_DIR}"/update_crosschain_params.json "$metadata"
}

## update erc20 params proposal
function updateERC20ParamsProposal() {

  metadata=$(base64metadata "update erc20 Params" "This proposal is used to modify the variable parameters of the erc20 module")

  deposit=$(queryMinDeposit "/fx.erc20.v1.MsgUpdateParams")

  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/update_erc20_params.json >update_erc20_params_tmp.json && mv update_erc20_params_tmp.json "${DATA_DIR}"/update_erc20_params.json
  jq '.deposit = "'"$deposit"'" ' "${DATA_DIR}"/update_erc20_params.json >update_erc20_params_tmp.json && mv update_erc20_params_tmp.json "${DATA_DIR}"/update_erc20_params.json

  submitProposal "${DATA_DIR}"/update_erc20_params.json "$metadata"
}

## update gov params proposal
function updateGovParamsProposal() {

  metadata=$(base64metadata "update gov Params" "This proposal is used to modify the variable parameters of the gov module")

  deposit=$(queryMinDeposit "/fx.gov.v1.MsgUpdateParams")

  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/update_gov_params.json >update_gov_params_tmp.json && mv update_gov_params_tmp.json "${DATA_DIR}"/update_gov_params.json
  jq '.deposit = "'"$deposit"'" ' "${DATA_DIR}"/update_gov_params.json >update_gov_params_tmp.json && mv update_gov_params_tmp.json "${DATA_DIR}"/update_gov_params.json

  submitProposal "${DATA_DIR}"/update_gov_params.json "$metadata"
}

## call contract
function callContractProposal() {

  metadata=$(base64metadata "evm call contract" "This proposal is used to call the method of the evm contract")

  deposit=$(queryMinDeposit "/fx.evm.v1.MsgCallContract")

  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/call_contract.json >call_contract_tmp.json && mv call_contract_tmp.json "${DATA_DIR}"/call_contract.json
  jq '.deposit = "'"$deposit"'" ' "${DATA_DIR}"/call_contract.json >call_contract_tmp.json && mv call_contract_tmp.json "${DATA_DIR}"/call_contract.json

  submitProposal "${DATA_DIR}"/call_contract.json "$metadata"
}

## apply for Genesis Ecological Fund
function distributionCommunityPoolSpendProposal() {

  title=$(jq -r '.title' "${DATA_DIR}"/leagcy_community_pool_spend.json)
  summary=$(jq -r '.description' "${DATA_DIR}"/leagcy_community_pool_spend.json)

  metadata=$(base64metadata "$title" "$summary")

  amount=$(jq -r '.amount' "${DATA_DIR}"/leagcy_community_pool_spend.json)

  deposit=$(queryMinDeposit "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal" "$amount")

  jq '.deposit = "'"$deposit"'" ' "${DATA_DIR}"/leagcy_community_pool_spend.json >leagcy_community_pool_spend_tmp.json && mv leagcy_community_pool_spend_tmp.json "${DATA_DIR}"/leagcy_community_pool_spend.json

  if ! fxcored tx gov submit-legacy-proposal community-pool-spend "${DATA_DIR}"/leagcy_community_pool_spend.json --gas-prices="$GAS_PRICES" --gas=auto --gas-adjustment="$GAS_ADJUSTMENT" --from fx1 -y --node "$JSON_RPC" --chain-id "$CHAIN_ID"; then
    exit 1
  fi

  for proposal_id in $(fxcored query gov proposals --node "$JSON_RPC" -o json | jq -r '.proposals[].id'); do
    messages=$(fxcored query gov proposal "${proposal_id}" --node "$JSON_RPC" -o json)
    if [[ "$(echo "${messages}" | jq -r '.messages[].metadata')" == "$metadata" && "$(echo "${messages}" | jq -r '.status')" == "PROPOSAL_STATUS_VOTING_PERIOD" ]]; then
      break
    fi
  done

  fxcored tx gov vote "${proposal_id}" yes --gas-prices="$GAS_PRICES" --gas=auto --gas-adjustment=1.4 --from fx1 -y --node "$JSON_RPC"
}

if [ "$1" == "registerCoinProposal" ] || [ "$1" == "registerERC20Proposal" ] || [ "$1" == "updateDenomAliasProposal" ] || [ "$1" == "toggleTokenConversionProposal" ] ||
  [ "$1" == "updateCrossChainParamsProposal" ] || [ "$1" == "updateERC20ParamsProposal" ] || [ "$1" == "updateGovParamsProposal" ] || [ "$1" == "callContractProposal" ] ||
  [ "$1" == "distributionCommunityPoolSpendProposal" ]; then
  "$@" || (echo "failed: $0" "$@" && exit 1)
else
  echo "$help"
fi
