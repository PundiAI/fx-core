#!/usr/bin/env bash

set -eo pipefail

export FX_HOME=${FX_HOME:-"/tmp/fxcore"}

if [[ "$1" == "init" ]]; then
  if [ -d "$FX_HOME" ]; then
    echo "node home '$FX_HOME' already exists"
    read -rp "Are you sure you want to delete all the data and start over? [y/N] " input
    [[ "$input" != "y" && "$input" != "Y" ]] && exit 1
    rm -r "$FX_HOME"
  fi

  # Initialize private validator, p2p, genesis, and application configuration files
  fxcored init local --chain-id fxcore --default-denom FX

  fxcored config config.toml rpc.cors_allowed_origins "*"
  # open prometheus
  fxcored config config.toml instrumentation.prometheus true
  # consensus
  fxcored config config.toml consensus.timeout_commit 1s
  # open rest and swagger
  fxcored config app.toml api.enable true
  fxcored config app.toml api.enabled-unsafe-cors true
  fxcored config app.toml api.swagger true
  # open telemetry
  fxcored config app.toml telemetry.enabled true
  fxcored config app.toml telemetry.enable-service-label true
  fxcored config app.toml telemetry.prometheus-retention-time 60
  # web3 api
  fxcored config app.toml json-rpc.api "eth,txpool,personal,net,debug,web3"

  # update fxcore client config
  fxcored config set client chain-id "fxcore"
  fxcored config set client keyring-backend "test"
  fxcored config set client output "json"
  fxcored config set client broadcast-mode "sync"

  echo "test test test test test test test test test test test junk" | fxcored keys add fx1 --recover
  if [ -n "${2:-""}" ]; then
    fxcored genesis add-genesis-account fx1 10004000000000000000000000FX
  else
    fxcored genesis add-genesis-account fx1 4000000000000000000000FX
  fi
  fxcored genesis gentx fx1 100000000000000000000FX --chain-id=fxcore \
    --gas="200000" \
    --moniker="fx-validator" \
    --commission-max-change-rate="0.01" \
    --commission-max-rate="0.2" \
    --commission-rate="0.01" \
    --details="Details A Function X foundation self-hosted validator." \
    --security-contact="contact@functionx.io" \
    --website="functionx.io"
  fxcored genesis collect-gentxs
fi

fxcored start
