#!/usr/bin/env bash

set -eo pipefail

PROJECT_DIR="${PROJECT_DIR:-"$(git rev-parse --show-toplevel)"}"
export PROJECT_DIR
export OUT_DIR="${PROJECT_DIR}/out"

readonly bridger_start_index=2
readonly bridger_oracle_number=3

export NODE_HOME="$OUT_DIR/.fxcore"

function create_oracles() {
  local chain_name=("$@")
  local index=${bridger_start_index}

  for chain in "${chain_name[@]}"; do
    oracles=$($DAEMON query "$chain" oracles | jq -r '.oracles[]')
    if [ ${#oracles} -gt 0 ]; then
      echo "oracles already exist"
      continue
    fi

    local oracle_file="$OUT_DIR/$chain-bridge-oracle.json"
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

. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"
