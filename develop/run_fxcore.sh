#!/bin/bash

set -e

if [[ "$1" == "init" ]]; then
  if [ "$(uname)" == 'Darwin' ]; then
    rm -rf ~/.fxcore
  else
    sudo rm -rf ~/.fxcore
  fi

  # Initialize private validator, p2p, genesis, and application configuration files
  fxcored init --chain-id=fxcore --denom=FX local

  # fxcored config config.toml rpc.pprof_laddr ""
  # open prometheus
  fxcored config config.toml instrumentation.prometheus true

  fxcored keys add fx1
  fxcored add-genesis-account fx1 3000000000000000000000000FX
  fxcored gentx fx1 100000000000000000000FX --chain-id=fxcore \
    --moniker="fx-validator" \
    --commission-max-change-rate=0.01 \
    --commission-max-rate=0.2 \
    --commission-rate=0.01 \
    --details="Details A Function X foundation self-hosted validator." \
    --security-contact="contact@functionx.io" \
    --website="functionx.io"
  fxcored collect-gentxs
  sed -i '' 's/378604525462891000000000000/438600525462891000000000000/g' ~/.fxcore/config/genesis.json
fi

fxcored start --api.enable --rpc.laddr 'tcp://0.0.0.0:26657' \
  --minimum-gas-prices '6000000000000FX' \
  --log_filter='ABCIQuery'
