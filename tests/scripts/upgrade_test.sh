#!/usr/bin/env bash

set -eo pipefail

JSON_RPC=${JSON_RPC:-"http://127.0.0.1:26657"}
out_dir=$(mktemp -d)
function clean() {
  rm -rf "$out_dir"
}
trap clean EXIT

cat <<EOF >"$out_dir"/register_coin.json
{
  "description": "Cross chain token of Function X",
  "denom_units": [
    {
      "denom": "test",
      "aliases": [
        "eth0x0000000000000000000000000000000000000000"
      ]
    },
    {
      "denom": "TEST",
      "exponent": 18
    }
  ],
  "base": "test",
  "display": "test",
  "name": "upgrade test token",
  "symbol": "TEST"
}
EOF

## register-coin
fxcored tx gov submit-legacy-proposal register-coin "$out_dir"/register_coin.json --title "Register test" \
  --description "This proposal creates and registers an ERC20 representation of test that can be bridged cross chains" \
  --deposit 10000000000000000000000FX --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.3 --from fx1 -y --node "${JSON_RPC}"

for proposal_id in $(fxcored query gov proposals -o json | jq -r '.proposals[].id'); do
  messages=$(fxcored query gov proposal "${proposal_id}" -o json)
  if [[ "$(echo "${messages}" | jq -r '.messages[].metadata')" == "$(echo '{"title":"Register test","summary":"This proposal creates and registers an ERC20 representation of test that can be bridged cross chains","metadata":""}' | base64)" && "$(echo "${messages}" | jq -r '.status')" == "PROPOSAL_STATUS_VOTING_PERIOD" ]]; then
    break
  fi
done

fxcored tx gov vote "${proposal_id}" yes --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.4 --from fx1 -y --node "${JSON_RPC}"

sleep 20

## update-denom-alias
fxcored tx gov submit-legacy-proposal update-denom-alias test bsc0x0000000000000000000000000000000000000000 --title "update denom alias" \
  --description "This proposal update denom alias" --deposit 10000000000000000000000FX --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.3 --from fx1 -y --node "${JSON_RPC}"

for proposal_id in $(fxcored query gov proposals -o json | jq -r '.proposals[].id'); do
  messages=$(fxcored query gov proposal "${proposal_id}" -o json)
  if [[ "$(echo "${messages}" | jq -r '.messages[].metadata')" == "$(echo '{"title":"update denom alias","summary":"This proposal update denom alias","metadata":""}' | base64)" && "$(echo "${messages}" | jq -r '.status')" == "PROPOSAL_STATUS_VOTING_PERIOD" ]]; then
    break
  fi
done

fxcored tx gov vote "${proposal_id}" yes --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.4 --from fx1 -y --node "${JSON_RPC}"

sleep 1

## toggle token conversion
fxcored tx gov submit-legacy-proposal toggle-token-conversion test --title "toggle token conversion" --description "This proposal toggle token conversion" \
  --deposit 10000000000000000000000FX --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.3 --from fx1 -y --node "tcp://127.0.0.1:26657"

for proposal_id in $(fxcored query gov proposals -o json | jq -r '.proposals[].id'); do
  messages=$(fxcored query gov proposal "${proposal_id}" -o json)
  if [[ "$(echo "${messages}" | jq -r '.messages[].metadata')" == "$(echo '{"title":"toggle token conversion","summary":"This proposal toggle token conversion","metadata":""}' | base64)" && "$(echo "${messages}" | jq -r '.status')" == "PROPOSAL_STATUS_VOTING_PERIOD" ]]; then
    break
  fi
done

fxcored tx gov vote "${proposal_id}" yes --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.4 --from fx1 -y --node "${JSON_RPC}"

sleep 3

cat <<EOF >"$out_dir"/msg_register_coin.json
{
  "messages": [
    {
      "@type": "/fx.erc20.v1.MsgRegisterCoin",
      "authority": "fx10d07y265gmmuvt4z0w9aw880jnsr700jqjzsmz",
      "metadata": {"description":"Cross chain token of Function X","denom_units":[{"denom":"test2","aliases":["eth0x0000000000000000000000000000000000000001"]},{"denom":"TEST2","exponent":18}],"base":"test2","display":"test2","name":"upgrade test2 token","symbol":"TEST2"}
    }
  ],
  "metadata": "eyJ0aXRsZSI6IlJlZ2lzdGVyIHRlc3QyIiwic3VtbWFyeSI6IlJlZ2lzdGVyIHRlc3QyIiwibWV0YWRhdGEiOiIifQo=",
  "deposit": "1000000000000000000000000FX"
}
EOF

## register-coin
fxcored tx gov submit-proposal "$out_dir"/msg_register_coin.json --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.3 --from fx1 -y --node "${JSON_RPC}" --chain-id fxcore

for proposal_id in $(fxcored query gov proposals -o json | jq -r '.proposals[].id'); do
  messages=$(fxcored query gov proposal "${proposal_id}" -o json)
  if [[ "$(echo "${messages}" | jq -r '.messages[].metadata')" == "$(echo '{"title":"Register test2","summary":"Register test2","metadata":""}' | base64)" && "$(echo "${messages}" | jq -r '.status')" == "PROPOSAL_STATUS_VOTING_PERIOD" ]]; then
    break
  fi
done

fxcored tx gov vote "${proposal_id}" yes --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.4 --from fx1 -y --node "${JSON_RPC}"

sleep 20

cat <<EOF >"$out_dir"/many_msg.json
{
  "messages": [
    {
      "@type": "/fx.erc20.v1.MsgUpdateDenomAlias",
      "authority": "fx10d07y265gmmuvt4z0w9aw880jnsr700jqjzsmz",
      "denom": "test2",
      "alias": "tron0x0000000000000000000000000000000000000001"
    },
    {
      "@type": "/fx.erc20.v1.MsgUpdateDenomAlias",
      "authority": "fx10d07y265gmmuvt4z0w9aw880jnsr700jqjzsmz",
      "denom": "test2",
      "alias": "bsc0x0000000000000000000000000000000000000001"
    },
    {
      "@type": "/fx.erc20.v1.MsgUpdateDenomAlias",
      "authority": "fx10d07y265gmmuvt4z0w9aw880jnsr700jqjzsmz",
      "denom": "test2",
      "alias": "polygon0x0000000000000000000000000000000000000001"
    }
  ],
  "metadata": "eyJ0aXRsZSI6Im1hbnkgbXNnIiwic3VtbWFyeSI6Im1hbnkgbXNnIiwibWV0YWRhdGEiOiIifQo=",
  "deposit": "1000000000000000000000000FX"
}
EOF

fxcored tx gov submit-proposal "$out_dir"/many_msg.json --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.3 --from fx1 -y --node "${JSON_RPC}" --chain-id fxcore

for proposal_id in $(fxcored query gov proposals -o json | jq -r '.proposals[].id'); do
  messages=$(fxcored query gov proposal "${proposal_id}" -o json)
  if [[ "$(echo "${messages}" | jq -r '.messages[].metadata')" == "$(echo '{"title":"many msg","summary":"many msg","metadata":""}' | base64)" && "$(echo "${messages}" | jq -r '.status')" == "PROPOSAL_STATUS_VOTING_PERIOD" ]]; then
    break
  fi
done

fxcored tx gov vote "${proposal_id}" yes --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.4 --from fx1 -y --node "${JSON_RPC}"

sleep 2

cat <<EOF >"$out_dir"/updateCrosschainParams.json
{
    "messages":[
        {
            "@type":"/fx.gravity.crosschain.v1.MsgUpdateParams",
            "authority":"fx10d07y265gmmuvt4z0w9aw880jnsr700jqjzsmz",
            "chain_name":"bsc",
            "params":{"gravityId":"fx-bsc-bridge","average_block_time":7000,"external_batch_timeout":43200000,"average_external_block_time":5000,"signed_window":30000,"slash_fraction":"0.810000000000000000","oracle_set_update_power_change_percent":"0.100000000000000000","ibc_transfer_timeout_height":20000,"delegate_threshold":{"denom":"FX","amount":"10000000000000000000000"},"delegate_multiple":10}
        },
        {
            "@type":"/fx.gravity.crosschain.v1.MsgUpdateParams",
            "authority":"fx10d07y265gmmuvt4z0w9aw880jnsr700jqjzsmz",
            "chain_name":"eth",
            "params":{"gravityId":"fx-eth-bridge","average_block_time":7000,"external_batch_timeout":43200000,"average_external_block_time":5000,"signed_window":30000,"slash_fraction":"0.810000000000000000","oracle_set_update_power_change_percent":"0.100000000000000000","ibc_transfer_timeout_height":20000,"delegate_threshold":{"denom":"FX","amount":"10000000000000000000000"},"delegate_multiple":10}
        },
        {
            "@type":"/fx.gravity.crosschain.v1.MsgUpdateParams",
            "authority":"fx10d07y265gmmuvt4z0w9aw880jnsr700jqjzsmz",
            "chain_name":"polygon",
            "params":{"gravityId":"fx-polygon-bridge","average_block_time":7000,"external_batch_timeout":43200000,"average_external_block_time":5000,"signed_window":30000,"slash_fraction":"0.810000000000000000","oracle_set_update_power_change_percent":"0.100000000000000000","ibc_transfer_timeout_height":20000,"delegate_threshold":{"denom":"FX","amount":"10000000000000000000000"},"delegate_multiple":10}
        },
        {
            "@type":"/fx.gravity.crosschain.v1.MsgUpdateParams",
            "authority":"fx10d07y265gmmuvt4z0w9aw880jnsr700jqjzsmz",
            "chain_name":"tron",
            "params":{"gravityId":"fx-tron-bridge","average_block_time":7000,"external_batch_timeout":43200000,"average_external_block_time":5000,"signed_window":30000,"slash_fraction":"0.810000000000000000","oracle_set_update_power_change_percent":"0.100000000000000000","ibc_transfer_timeout_height":20000,"delegate_threshold":{"denom":"FX","amount":"10000000000000000000000"},"delegate_multiple":10}
        },
        {
            "@type":"/fx.gravity.crosschain.v1.MsgUpdateParams",
            "authority":"fx10d07y265gmmuvt4z0w9aw880jnsr700jqjzsmz",
            "chain_name":"avalanche",
            "params":{"gravityId":"fx-avalanche-bridge","average_block_time":7000,"external_batch_timeout":43200000,"average_external_block_time":5000,"signed_window":30000,"slash_fraction":"0.810000000000000000","oracle_set_update_power_change_percent":"0.100000000000000000","ibc_transfer_timeout_height":20000,"delegate_threshold":{"denom":"FX","amount":"10000000000000000000000"},"delegate_multiple":10}
        }
    ],
    "metadata": "eyJ0aXRsZSI6InVwZGF0ZSBDcm9zc2NoYWluIFBhcmFtcyIsInN1bW1hcnkiOiJ1cGRhdGUgQ3Jvc3NjaGFpbiBQYXJhbXMiLCJtZXRhZGF0YSI6IiJ9Cg==",
    "deposit":"1000000000000000000000000FX"
}
EOF

fxcored tx gov submit-proposal "$out_dir"/updateCrosschainParams.json --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.3 --from fx1 -y --node "${JSON_RPC}" --chain-id fxcore
for proposal_id in $(fxcored query gov proposals -o json | jq -r '.proposals[].id'); do
  messages=$(fxcored query gov proposal "${proposal_id}" -o json)
  if [[ "$(echo "${messages}" | jq -r '.messages[].metadata')" == "$(echo '{"title":"update Crosschain Params","summary":"update Crosschain Params","metadata":""}' | base64)" && "$(echo "${messages}" | jq -r '.status')" == "PROPOSAL_STATUS_VOTING_PERIOD" ]]; then
    break
  fi
done

fxcored tx gov vote "${proposal_id}" yes --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.4 --from fx1 -y --node "${JSON_RPC}"

sleep 2

cat <<EOF >"$out_dir"/updateERC20Params.json
{
    "messages":[
        {
            "@type":"/fx.erc20.v1.MsgUpdateParams",
            "authority":"fx10d07y265gmmuvt4z0w9aw880jnsr700jqjzsmz",
            "params": {"enable_erc20":true,"enable_evm_hook":true,"ibc_timeout":"12h"}
        }
    ],
    "metadata": "eyJ0aXRsZSI6InVwZGF0ZSBlcmMyMCBQYXJhbXMiLCJzdW1tYXJ5IjoidXBkYXRlIGVyYzIwIFBhcmFtcyIsIm1ldGFkYXRhIjoiIn0K",
    "deposit":"1000000000000000000000000FX"
}
EOF

fxcored tx gov submit-proposal "$out_dir"/updateERC20Params.json --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.3 --from fx1 -y --node "${JSON_RPC}" --chain-id fxcore
for proposal_id in $(fxcored query gov proposals -o json | jq -r '.proposals[].id'); do
  messages=$(fxcored query gov proposal "${proposal_id}" -o json)
  if [[ "$(echo "${messages}" | jq -r '.messages[].metadata')" == "$(echo '{"title":"update erc20 Params","summary":"update erc20 Params","metadata":""}' | base64)" && "$(echo "${messages}" | jq -r '.status')" == "PROPOSAL_STATUS_VOTING_PERIOD" ]]; then
    break
  fi
done

fxcored tx gov vote "${proposal_id}" yes --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.4 --from fx1 -y --node "${JSON_RPC}"

sleep 2

cat <<EOF >"$out_dir"/updateGovParams.json
{
    "messages":[
        {
            "@type":"/fx.gov.v1.MsgUpdateParams",
            "authority":"fx10d07y265gmmuvt4z0w9aw880jnsr700jqjzsmz",
            "params": {"msg_type":"","min_deposit":[{"denom":"FX","amount":"1000"}],"min_initial_deposit":{"denom":"FX","amount":"1000"},"voting_period":"12097000s","quorum":"0.3","max_deposit_period":"12097000s","threshold":"0.5","veto_threshold":"0.334"}
        },
        {
            "@type":"/fx.gov.v1.MsgUpdateParams",
            "authority":"fx10d07y265gmmuvt4z0w9aw880jnsr700jqjzsmz",
            "params": {"msg_type":"/fx.erc20.v1.MsgUpdateDenomAlias","min_deposit":[{"denom":"FX","amount":"1000"}],"min_initial_deposit":{"denom":"FX","amount":"1000"},"voting_period":"15s","quorum":"0.25","max_deposit_period":"12097000s","threshold":"0.5","veto_threshold":"0.334"}
        },
        {
            "@type":"/fx.gov.v1.MsgUpdateParams",
            "authority":"fx10d07y265gmmuvt4z0w9aw880jnsr700jqjzsmz",
            "params": {"msg_type":"/fx.evm.v1.MsgCallContract","min_deposit":[{"denom":"FX","amount":"1000"}],"min_initial_deposit":{"denom":"FX","amount":"1000"},"voting_period":"20s","quorum":"0.25","max_deposit_period":"12097000s","threshold":"0.5","veto_threshold":"0.334"}
        }   
    ],
    "metadata":"eyJ0aXRsZSI6InVwZGF0ZSBnb3YgUGFyYW1zIiwic3VtbWFyeSI6InVwZGF0ZSBnb3YgUGFyYW1zIiwibWV0YWRhdGEiOiIifQo=",
    "deposit":"1000000000000000000000000FX"
}
EOF

fxcored tx gov submit-proposal "$out_dir"/updateGovParams.json --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.3 --from fx1 -y --node "${JSON_RPC}" --chain-id fxcore
for proposal_id in $(fxcored query gov proposals -o json | jq -r '.proposals[].id'); do
  messages=$(fxcored query gov proposal "${proposal_id}" -o json)
  if [[ "$(echo "${messages}" | jq -r '.messages[].metadata')" == "$(echo '{"title":"update gov Params","summary":"update gov Params","metadata":""}' | base64)" && "$(echo "${messages}" | jq -r '.status')" == "PROPOSAL_STATUS_VOTING_PERIOD" ]]; then
    break
  fi
done

fxcored tx gov vote "${proposal_id}" yes --gas-prices=4000000000000FX --gas=auto --gas-adjustment=1.4 --from fx1 -y --node "${JSON_RPC}"

sleep 2

## migrated parameter query
echo ''

echo 'Query the parameter values after migration of each module...'
echo "erc20 params: $(fxcored q erc20 params --node "${JSON_RPC}" | jq .)"
echo ''
echo "eth params: $(fxcored q crosschain eth params --node "${JSON_RPC}" | jq .)"
echo ''
echo "bsc params: $(fxcored q crosschain bsc params --node "${JSON_RPC}" | jq .)"
echo ''
echo "polygon params: $(fxcored q crosschain polygon params --node "${JSON_RPC}" | jq .)"
echo ''
echo "tron params: $(fxcored q crosschain tron params --node "${JSON_RPC}" | jq .)"
echo ''
echo "avalanche params: $(fxcored q crosschain avalanche params --node "${JSON_RPC}" | jq .)"
echo ''
