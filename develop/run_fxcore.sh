#!/bin/bash

set -o errexit -o pipefail

if [[ "$1" == "init" ]]; then
  if [ -d ~/.fxcore ]; then
    read -rp "Are you sure you want to delete all the data and start over? [y/N] " input
    if [[ "$input" != "y" && "$input" != "Y" ]]; then
      exit 1
    fi
    rm -r ~/.fxcore
  fi

  # Initialize private validator, p2p, genesis, and application configuration files
  fxcored init local --chain-id fxcore

  fxcored config config.toml rpc.cors_allowed_origins "*"
  # open prometheus
  fxcored config config.toml instrumentation.prometheus true
  # consensus
  # fxcored config config.toml consensus.timeout_commit 1s
  # open rest and swagger
  fxcored config app.toml api.enable true
  fxcored config app.toml api.enabled-unsafe-cors
  fxcored config app.toml api.swagger true
  # open telemetry
  fxcored config app.toml telemetry.enabled true
  fxcored config app.toml telemetry.enable-service-label true
  fxcored config app.toml telemetry.prometheus-retention-time 60
  # web3 api
  fxcored config app.toml json-rpc.api "eth,txpool,personal,net,debug,web3"

  # update fxcore client config
  fxcored config chain-id fxcore
  fxcored config keyring-backend test
  fxcored config output json
  fxcored config broadcast-mode "block"

  fxcored keys add fx1 --algo secp256k1 --coin-type 118
  fxcored add-genesis-account fx1 4000000000000000000000FX
  fxcored gentx fx1 100000000000000000000FX --chain-id=fxcore \
    --gas="200000" \
    --moniker="fx-validator" \
    --commission-max-change-rate=0.01 \
    --commission-max-rate=0.2 \
    --commission-rate=0.01 \
    --details="Details A Function X foundation self-hosted validator." \
    --security-contact="contact@functionx.io" \
    --website="functionx.io"
  fxcored collect-gentxs
fi

fxcored start
