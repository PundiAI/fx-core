#!/usr/bin/env bash

set -eo pipefail

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/setup-env.sh"

## ARGS: <self_delegation_amount>
function create_validator() {
  local self_delegation_amount=$1

  cosmos_tx staking create-validator \
    --amount "$self_delegation_amount$STAKING_DENOM" \
    --pubkey "$($DAEMON tendermint show-validator --home "$NODE_HOME")" \
    --commission-max-change-rate=0.01 \
    --commission-max-rate=0.2 \
    --commission-rate=0.03 \
    --min-self-delegation="$(to_18 100)" \
    --from "$FROM"
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
