#!/usr/bin/env bash

set -eo pipefail

readonly bridger_start_index=2
readonly bridger_oracle_number=3

project_dir="$(git rev-parse --show-toplevel)"
readonly project_dir
readonly out_dir="${project_dir}/out"

export NODE_HOME="$out_dir/.fxcore"

function create_oracles() {
  local chain_name=("$@")
  index=${bridger_start_index}

  for chain in "${chain_name[@]}"; do
    local oracle_file="$out_dir/$chain-bridge-oracle.json"

    oracles=$($DAEMON query "$chain" oracles | jq -r '.oracles[]')
    if [ ${#oracles} -gt 0 ]; then
      continue
    fi
    echo "[]" >"$oracle_file"

    for ((i = 0; i < "$bridger_oracle_number"; i++)); do
      oracle_address=$(add_key "$chain-oracle-$i" "$index" | jq -r ".address")
      bridger_address=$(add_key "$chain-bridger-$i" "$((index + 1))" | jq -r ".address")
      external_address=$(show_address "$chain-bridger-$i" -e)

      jq -cs add "$oracle_file" <(echo "[{\"oracle_address\":\"$oracle_address\",\"bridge_address\":\"$bridger_address\",\"external_address\":\"$external_address\",\"oracle_index\":\"$index\",\"bridge_index\":\"$((index + 1))\"}]") >"$oracle_file.tmp" &&
        mv "$oracle_file.tmp" "$oracle_file"

      index=$((index + 2))
    done
  done
}

. "${project_dir}/tests/scripts/setup-env.sh"
