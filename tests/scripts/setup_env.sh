#!/usr/bin/env bash

# check dependencies commands are installed
commands=(jq fxcored)
for cmd in "${commands[@]}"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "$cmd command not found, please install $cmd first" && exit 1
  fi
done

export JSON_RPC="http://127.0.0.1:26657"
export GAS_PRICES="4000000000000FX"
export GAS_ADJUSTMENT="1.3"
export CHAIN_ID="fxcore"
