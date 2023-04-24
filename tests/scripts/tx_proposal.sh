#!/usr/bin/env bash

set -eo pipefail

# check dependencies commands are installed
commands=(jq fxcored)
for cmd in "${commands[@]}"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "$cmd command not found, please install $cmd first" && exit 1
  fi
done

JSON_RPC="http://localhost:26657"
GAS_PRICES="4000000000000FX"
GAS_ADJUSTMENT="1.3"
CHAIN_ID="fxcore"
DATA_DIR=$(dirname "$(dirname "$(realpath "$0")")")/data/proposal

submit_proposal() {
  local data_file="$1"
  local metadata="$2"

  fxcored tx gov submit-proposal "$data_file" --gas-prices="$GAS_PRICES" --gas=auto --gas-adjustment="$GAS_ADJUSTMENT" --from fx1 -y --node "$JSON_RPC" --chain-id "$CHAIN_ID"
  for proposal_id in $(fxcored query gov proposals --node "$JSON_RPC" -o json | jq -r '.proposals[].id'); do
    messages=$(fxcored query gov proposal "${proposal_id}" --node "$JSON_RPC" -o json)
    if [[ "$(echo "${messages}" | jq -r '.messages[].metadata')" == "$metadata" && "$(echo "${messages}" | jq -r '.status')" == "PROPOSAL_STATUS_VOTING_PERIOD" ]]; then
      break
    fi
  done

  fxcored tx gov vote "${proposal_id}" yes --gas-prices="$GAS_PRICES" --gas=auto --gas-adjustment=1.4 --from fx1 -y --node "$JSON_RPC"
}

## register coin proposal
function msg_register_coin_proposal() {
  if [ -z "$1" ]; then
    title="register coin"
  else
    title="$1"
  fi

  if [ -z "$2" ]; then
    summary="erc20 register coin proposal"
  else
    summary="$2"
  fi

  metadata="$(echo '{"title": "'"$title"'","summary": "'"$summary"'","metadata":""}' | base64)"
  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/register_coin.json >register_coin_tmp.json && mv register_coin_tmp.json "${DATA_DIR}"/register_coin.json

  submit_proposal "${DATA_DIR}"/register_coin.json "$metadata"
}

## register erc20 proposal
function msg_register_erc20_proposal() {
  if [ -z "$1" ]; then
    title="register fip20 contract token"
  else
    title="$1"
  fi
  if [ -z "$2" ]; then
    summary="erc20 register fip20 contract token"
  else
    summary="$2"
  fi

  metadata="$(echo '{"title": "'"$title"'","summary": "'"$summary"'","metadata":""}' | base64)"
  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/register_erc20.json >register_erc20_tmp.json && mv register_erc20_tmp.json "${DATA_DIR}"/register_erc20.json

  for erc20address in $(jq -r '.messages[].erc20address' "${DATA_DIR}"/register_erc20.json); do
    if [[ $(fxcored q evm code "${erc20address}" --node "$JSON_RPC" -o json | jq -r '.code') == null ]]; then
      echo "contract $erc20address not deployed, cannot initiate MsgRegisterERC20 proposal" && exit 1
    fi
  done

  submit_proposal "${DATA_DIR}"/register_erc20.json "$metadata"
}

# update denom alias proposal
function msg_update_denom_alias_proposal() {
  if [ -z "$1" ]; then
    title="update denom alias"
  else
    title="$1"
  fi
  if [ -z "$2" ]; then
    summary="erc20 update denom alias proposal"
  else
    summary="$2"
  fi

  metadata="$(echo '{"title": "'"$title"'","summary": "'"$summary"'","metadata":""}' | base64)"
  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/update_denom_alias.json >update_denom_alias_tmp.json && mv update_denom_alias_tmp.json "${DATA_DIR}"/update_denom_alias.json

  for denom in $(jq -r '.messages[].denom' "${DATA_DIR}"/update_denom_alias.json); do
    if ! fxcored q erc20 token-pair "${denom}" --node "$JSON_RPC" >/dev/null 2>&1; then
      echo "token $denom unregistered, cannot initiate MsgUpdateDenomAlias proposal" && exit 1
    fi
  done

  submit_proposal "${DATA_DIR}"/update_denom_alias.json "$metadata"
}

## toggle token conversion proposal
function msg_toggle_token_conversion_proposal() {
  if [ -z "$1" ]; then
    title="toggle token conversion"
  else
    title="$1"
  fi
  if [ -z "$2" ]; then
    summary="toggle token conversion proposal"
  else
    summary="$2"
  fi

  metadata="$(echo '{"title": "'"$title"'","summary": "'"$summary"'","metadata":""}' | base64)"
  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/toggle_token_conversion.json >toggle_token_conversion_tmp.json && mv toggle_token_conversion_tmp.json "${DATA_DIR}"/toggle_token_conversion.json

  for token in $(jq -r '.messages[].token' "${DATA_DIR}"/toggle_token_conversion.json); do
    if ! fxcored q erc20 token-pair "${token}" --node "$JSON_RPC" >/dev/null 2>&1; then
      echo "$token unregistered, cannot initiate MsgToggleTokenConversion proposal" && exit 1
    fi
  done

  submit_proposal "${DATA_DIR}"/toggle_token_conversion.json "$metadata"
}

## update crosschain params
function update_crosschain_params_proposal() {
  if [ -z "$1" ]; then
    title="update crosschain Params"
  else
    title="$1"
  fi
  if [ -z "$2" ]; then
    summary="update crosschain Params proposal"
  else
    summary="$2"
  fi

  metadata="$(echo '{"title": "'"$title"'","summary": "'"$summary"'","metadata":""}' | base64)"
  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/update_crosschain_params.json >update_crosschain_params_tmp.json && mv update_crosschain_params_tmp.json "${DATA_DIR}"/update_crosschain_params.json

  submit_proposal "${DATA_DIR}"/update_crosschain_params.json "$metadata"
}

## update erc20 params proposal
function update_erc20_params_proposal() {
  if [ -z "$1" ]; then
    title="update erc20 Params"
  else
    title="$1"
  fi
  if [ -z "$2" ]; then
    summary="update erc20 Params proposal"
  else
    summary="$2"
  fi

  metadata="$(echo '{"title": "'"$title"'","summary": "'"$summary"'","metadata":""}' | base64)"
  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/update_erc20_params.json >update_erc20_params_tmp.json && mv update_erc20_params_tmp.json "${DATA_DIR}"/update_erc20_params.json

  submit_proposal "${DATA_DIR}"/update_erc20_params.json "$metadata"
}

## update gov params proposal
function update_gov_params_proposal() {
  if [ -z "$1" ]; then
    title="update gov Params"
  else
    title="$1"
  fi
  if [ -z "$2" ]; then
    summary="update gov Params proposal"
  else
    summary="$2"
  fi

  metadata="$(echo '{"title": "'"$title"'","summary": "'"$summary"'","metadata":""}' | base64)"
  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/update_gov_params.json >update_gov_params_tmp.json && mv update_gov_params_tmp.json "${DATA_DIR}"/update_gov_params.json

  submit_proposal "${DATA_DIR}"/update_gov_params.json "$metadata"
}

## call contract
function call_contract_proposal() {
  if [ -z "$1" ]; then
    title="call contract"
  else
    title="$1"
  fi
  if [ -z "$2" ]; then
    summary="call contract proposal"
  else
    summary="$2"
  fi

  metadata="$(echo '{"title": "'"$title"'","summary": "'"$summary"'","metadata":""}' | base64)"
  jq '.metadata = "'"$metadata"'" ' "${DATA_DIR}"/call_contract.json >call_contract_tmp.json && mv call_contract_tmp.json "${DATA_DIR}"/call_contract.json

  submit_proposal "${DATA_DIR}"/call_contract.json "$metadata"
}

## apply for Genesis Ecological Fund
function distribution_community_pool_spend_proposal() {
  title=$(jq -r '.title' "${DATA_DIR}"/leagcy_community_pool_spend.json)
  summary=$(jq -r '.description' "${DATA_DIR}"/leagcy_community_pool_spend.json)
  metadata="$(echo '{"title": "'"$title"'","summary": "'"$summary"'","metadata":""}' | base64)"

  fxcored tx gov submit-legacy-proposal community-pool-spend "${DATA_DIR}"/leagcy_community_pool_spend.json --gas-prices="$GAS_PRICES" --gas=auto --gas-adjustment="$GAS_ADJUSTMENT" --from fx1 -y --node "$JSON_RPC" --chain-id "$CHAIN_ID"

  for proposal_id in $(fxcored query gov proposals --node "$JSON_RPC" -o json | jq -r '.proposals[].id'); do
    messages=$(fxcored query gov proposal "${proposal_id}" --node "$JSON_RPC" -o json)
    if [[ "$(echo "${messages}" | jq -r '.messages[].metadata')" == "$metadata" && "$(echo "${messages}" | jq -r '.status')" == "PROPOSAL_STATUS_VOTING_PERIOD" ]]; then
      break
    fi
  done

  fxcored tx gov vote "${proposal_id}" yes --gas-prices="$GAS_PRICES" --gas=auto --gas-adjustment=1.4 --from fx1 -y --node "$JSON_RPC"
}
