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

## define an array to store the name of each module
cross_chain_module=(eth bsc polygon tron avalanche)

function QueryParams() {
  if [ -z "$1" ]; then
    echo "What modules do you need to query?" && exit 1
  else
    module="$1"
  fi
  echo "Querying ${module} params..."
  # shellcheck disable=SC2199
  if [[ "${cross_chain_module[@]}" =~ ${module} ]]; then
    fxcored q crosschain "${module}" params --node "${JSON_RPC}" | jq .
  else
    fxcored q "${module}" params --node "${JSON_RPC}" | jq .
  fi
}
