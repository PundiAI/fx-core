#!/bin/bash

set -e

if [[ "$1" == "init" ]]; then
  read -p "Are you sure you want to delete all the data and start over? [y/N] " input
  if [[ "$input" != "y" && "$input" != "Y" ]]; then
    exit 1
  fi
  rm -r ~/.fxcore

  # Initialize private validator, p2p, genesis, and application configuration files
  fxcored init --chain-id=fxcore --denom=FX local

  # fxcored config config.toml rpc.pprof_laddr ""
  # open prometheus
  fxcored config config.toml instrumentation.prometheus true

  # update fxcore client config
  fxcored config chain-id fxcore
  fxcored config keyring-backend test
  fxcored config output json
  fxcored config node "tcp://127.0.0.1:26657"
  fxcored config broadcast-mode "block"

  fxcored keys add fx1
  fxcored add-genesis-account fx1 4000000000000000000000FX
  fxcored gentx fx1 1000000000000000000000FX --chain-id=fxcore \
    --moniker="fx-validator" \
    --commission-max-change-rate=0.01 \
    --commission-max-rate=0.2 \
    --commission-rate=0.01 \
    --details="Details A Function X foundation self-hosted validator." \
    --security-contact="contact@functionx.io" \
    --website="functionx.io"
  fxcored collect-gentxs
fi

fxcored start --log_filter='ABCIQuery'
