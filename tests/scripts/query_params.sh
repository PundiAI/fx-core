#!/usr/bin/env bash

set -eo pipefail

BASE_DIR=${PROJECT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]%/*}")" && pwd)}

# shellcheck source=/dev/null
source "$BASE_DIR"/scripts/setup_env.sh

## define an array to store the name of each module
cross_chain_module=(eth bsc polygon tron avalanche arbitrum optimism)

help="Command:
    queryParams                             query module params.(erc20 gov eth bsc polygon tron avalanche arbitrum optimism)
"

function queryParams() {
  if [ -z "$1" ]; then
    echo "What modules do you need to query?" && exit 1
  else
    module="$1"
  fi
  echo "${module} params: "
  # shellcheck disable=SC2199
  if [[ "${cross_chain_module[@]}" =~ ${module} ]]; then
    fxcored q crosschain "${module}" params --node "${JSON_RPC}" | jq .
  elif [[ "${module}" == "gov" ]]; then
    fxcored q "${module}" params --node "${JSON_RPC}" --msg-type="$2" | jq .
  else
    fxcored q "${module}" params --node "${JSON_RPC}" | jq .
  fi
}

if [ "$1" == "queryParams" ]; then
  "$@" || (echo "failed: $0" "$@" && exit 1)
else
  echo "$help"
fi
