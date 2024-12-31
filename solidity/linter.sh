#!/usr/bin/env bash

set -eo pipefail

npx solhint --config ./.solhint.json --max-warnings 0 "contracts/**/*.sol"

# Verify both files exist
for file in ./contracts/bridge/{FxBridgeLogic,FxBridgeLogicETH}.sol; do
  if [[ ! -f "$file" ]]; then
    echo "Error: $file not found"
    exit 1
  fi
done

# Find the last occurrence of INIT to handle multiple matches
init_line=$(grep -n "INIT" ./contracts/bridge/FxBridgeLogic.sol | tail -1 | cut -d: -f1)
if [[ -z "$init_line" ]]; then
  echo "Error: INIT marker not found in FxBridgeLogic.sol"
  exit 1
fi

# Calculate lines to compare
total_lines=$(wc -l <./contracts/bridge/FxBridgeLogic.sol)
lines=$((total_lines - init_line))

# Compare files - exit with error if they are identical
if ! diff <(tail -n "$lines" ./contracts/bridge/FxBridgeLogic.sol) <(tail -n "$lines" ./contracts/bridge/FxBridgeLogicETH.sol); then
  echo "Error: FxBridgeLogic.sol and FxBridgeLogicETH.sol have identical implementations"
  exit 1
fi
