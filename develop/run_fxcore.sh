#!/bin/bash

set -e

if [[ "$1" == "init" ]]; then
  if [ -d ~/.fxcore ]; then
    read -p "Are you sure you want to delete all the data and start over? [y/N] " input
    if [[ "$input" != "y" && "$input" != "Y" ]]; then
      exit 1
    fi
    rm -r ~/.fxcore
  fi

  # Initialize private validator, p2p, genesis, and application configuration files
  fxcored init local --chain-id fxcore

  # set mini gas price
  fxcored config app.toml minimum-gas-prices 1000000000FX
  # consensus
  fxcored config config.toml consensus.timeout_commit 1s
  # open rest and swagger
  fxcored config app.toml api.enable true
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

  fxcored keys add fx1
  perl -pi -e 's|378604525462891000000000000|778600525462891000000000000|g' ~/.fxcore/config/genesis.json
  perl -pi -e 's|"amount": "10000000000000000000000"|"amount": "100000000000000000000"|g' ~/.fxcore/config/genesis.json
  perl -pi -e 's|"voting_period": "1209600s"|"voting_period": "30s"|g' ~/.fxcore/config/genesis.json
  perl -pi -e 's|"unbonding_time": "1814400s"|"unbonding_time": "1800s"|g' ~/.fxcore/config/genesis.json
  fxcored add-genesis-account fx1 400000000000000000000000000FX
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
