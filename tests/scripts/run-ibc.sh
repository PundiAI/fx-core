#!/usr/bin/env bash

project_dir="$(git rev-parse --show-toplevel)"
readonly project_dir
readonly out_dir="${project_dir}/out"

readonly ibc_from="ibc-$FROM"
readonly ibc_home_dir="$out_dir/.ibcrelayer"
[[ -d "${ibc_home_dir}" ]] && rm -rf "${ibc_home_dir}"
mkdir -p "${ibc_home_dir}"
readonly a_chain_name="fxcore"
readonly b_chain_name="pundix"
readonly docker_image="functionx/ibc-relayer:latest"

for chain_name in $a_chain_name $b_chain_name; do
  (
    # shellcheck source=/dev/null
    . "${project_dir}/tests/scripts/$chain_name.sh"

    echo "$TEST_MNEMONIC" | $DAEMON keys add "$ibc_from" --index=1 --recover --home "$NODE_HOME"

    node_catching_up "$NODE_RPC"
    cosmos_tx bank send "$FROM" "$($DAEMON keys show "$ibc_from" --home "$NODE_HOME" -a)" "$(to_18 "200")$STAKING_DENOM" \
      --from "$FROM" --home "$NODE_HOME"
  )
done

cat >"${ibc_home_dir}"/config/config.yaml <<EOF
global:
  api-listen-addr: :5183
  timeout: 10s
  light-cache-size: 20
chains:
-
  chain-id: "$(jq -r '.chain_id' "$out_dir/$a_chain_name.json")"
  key-type: "mnemonic"
  key-value: "$TEST_MNEMONIC"
  pub-key-type: "ethermint/PubKeyEthSecp256k1"
  hd-path: "m/44'/60'/0'/0/1"

  rpc-addr: "$(jq -r '.node_rpc' "$out_dir/$a_chain_name.json")"
  account-prefix: $(jq -r '.bech32_prefix' "$out_dir/$a_chain_name.json")
  gas-adjustment: 1.01
  gas-prices: "$(jq -r '.gas_price' "$out_dir/$a_chain_name.json")"
  trusting-period: 29m
  skip-un-relay-sequences: []
  iterator-block-config:
    delay-second: 3
    handler-block-count: 100
    batch-handle-block-count: 15
-
  chain-id: "$(jq -r '.chain_id' "$out_dir/$b_chain_name.json")"
  key-type: "mnemonic"
  key-value: "$TEST_MNEMONIC"
  pub-key-type: "tendermint/PubKeySecp256k1"
  hd-path: "m/44'/118'/0'/0/1"

  rpc-addr: "$(jq -r '.node_rpc' "$out_dir/$b_chain_name.json")"
  account-prefix: $(jq -r '.bech32_prefix' "$out_dir/$b_chain_name.json")
  gas-adjustment: 1.01
  gas-prices: "$(jq -r '.gas_price' "$out_dir/$b_chain_name.json")"
  trusting-period: 29m
  skip-un-relay-sequences: []
  iterator-block-config:
    delay-second: 3
    handler-block-count: 100
    batch-handle-block-count: 15

paths: {}
EOF

docker run --rm --name ibc-relay-temp -v "${ibc_home_dir}":/root/.relayer --network bridge "$docker_image" \
  config init --home=/root/.relayer

docker run --rm --name ibc-relay-temp -v "${ibc_home_dir}":/root/.relayer --network bridge "$docker_image" \
  paths generate \
  "$(jq -r '.chain_id' "$out_dir/$a_chain_name.json")" "$(jq -r '.chain_id' "$out_dir/$b_chain_name.json")" \
  transfer --port=transfer --home=/root/.relayer

docker run --rm --name ibc-relay-temp -v "${ibc_home_dir}":/root/.relayer --network bridge "$docker_image" \
  light init "$(jq -r '.chain_id' "$out_dir/$a_chain_name.json")" -f --home=/root/.relayer

docker run --rm --name ibc-relay-temp -v "${ibc_home_dir}":/root/.relayer --network bridge "$docker_image" \
  light init "$(jq -r '.chain_id' "$out_dir/$b_chain_name.json")" -f --home=/root/.relayer

docker run --rm --name ibc-relay-temp -v "${ibc_home_dir}":/root/.relayer --network bridge "$docker_image" \
  tx link transfer -d --home=/root/.relayer

docker run -itd --name ibc-relay -v "${ibc_home_dir}":/root/.relayer --network bridge "$docker_image" \
  start transfer --home=/root/.relayer --time-threshold=19m --notify.enable=false --debug
