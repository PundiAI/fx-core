#!/usr/bin/env bash

set -eo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

## ARGS: <chain_name> <dest> <amount> <bridge-fee> [options...]
function send_to_external() {
  local chain_name=$1 && shift
  cosmos_tx "$chain_name" send-to-external "$@"
}

## ARGS: <src_port> <src_channel> <receiver> <amount> [options...]
function ibc_transfer() {
  cosmos_tx ibc-transfer transfer "$@"
}

## ARGS: <src_port> <src_channel> <receiver> <amount> [options...]
function fx_ibc_transfer() {
  cosmos_tx fx-ibc-transfer transfer "$@"
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
