#!/bin/bash

set -e

if [[ "$1" == "init" ]]; then
  rm -r ~/.fxcore

  # Initialize private validator, p2p, genesis, and application configuration files
  fxcored init --chain-id=fxcore --denom=FX local

  # fxcored config config.toml rpc.pprof_laddr ""
  # open prometheus
  fxcored config config.toml instrumentation.prometheus true

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

fxcored start --minimum-gas-prices '4000000000000FX' --log_filter='ABCIQuery'
