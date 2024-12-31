#!/usr/bin/env bash

set -eo pipefail

npx solhint --config ./.solhint.json --max-warnings 0 "contracts/**/*.sol"

lines=$(($(wc -l ./contracts/bridge/FxBridgeLogic.sol | awk '{print $1}') - $(grep -n "INIT" ./contracts/bridge/FxBridgeLogic.sol | cut -d: -f1)))
if ! diff <(tail -n $lines ./contracts/bridge/FxBridgeLogic.sol) <(tail -n $lines ./contracts/bridge/FxBridgeLogicETH.sol); then
  echo "FxBridgeLogic.sol and FxBridgeLogicETH.sol are not different"
  exit 1
fi
