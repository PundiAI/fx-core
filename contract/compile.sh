#!/usr/bin/env bash

set -eo pipefail

commands=(git yarn jq abigen)
for cmd in "${commands[@]}"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "$cmd command not found, please install $cmd first" && exit 1
  fi
done

abigen_version=$(abigen --version | awk '{print $3}')
if ! [[ "$abigen_version" =~ ^1.12.0-stable.* ]]; then
  echo "expected abigen version 1.12.0, but got $abigen_version, please upgrade abigen first" && exit 1
fi

project_dir="$(git rev-parse --show-toplevel)"
if [ ! -d "$project_dir/solidity/contracts/node_modules" ]; then
  echo "===> Installing node modules"
  (cd "$project_dir/solidity" && yarn install)
fi

if [ -d "$project_dir/solidity/artifacts" ]; then
  echo "===> Cleaning artifacts"
  (cd "$project_dir/solidity" && yarn clean)
fi

echo "===> Compiling contracts"
(cd "$project_dir/solidity" && yarn compile)

[[ ! -d "$project_dir/contract/artifacts" ]] && mkdir -p "$project_dir/contract/artifacts"

# add core contracts
contracts=(WFXUpgradable FIP20Upgradable ICrossChain IStaking IFxBridgeLogic)
contracts_test=(CrossChainTest StakingTest)
# add 3rd party contracts
contracts+=(ERC1967Proxy)
contracts_test+=(ERC721TokenTest)

for contract in "${contracts[@]}"; do
  echo "===> Ethereum ABI wrapper code generator: $contract"
  file_path=$(find "$project_dir/solidity/artifacts" -name "${contract}.json" -type f)
  jq -c '.abi' "$file_path" >"$project_dir/contract/artifacts/${contract}.abi"
  jq -r '.bytecode' "$file_path" >"$project_dir/contract/artifacts/${contract}.bin"
  abigen --abi "$project_dir/contract/artifacts/${contract}.abi" \
    --bin "$project_dir/contract/artifacts/${contract}.bin" \
    --type "${contract}" --pkg contract \
    --out "$project_dir/contract/${contract}.go"
done

# test contracts
for contract_test in "${contracts_test[@]}"; do
  echo "===> Ethereum ABI wrapper code generator: $contract_test"
  file_path=$(find "$project_dir/solidity/artifacts" -name "${contract_test}.json" -type f)
  jq -c '.abi' "$file_path" >"$project_dir/contract/artifacts/${contract_test}.abi"
  jq -r '.bytecode' "$file_path" >"$project_dir/contract/artifacts/${contract_test}.bin"
  abigen --abi "$project_dir/contract/artifacts/${contract_test}.abi" \
    --bin "$project_dir/contract/artifacts/${contract_test}.bin" \
    --type "${contract_test}" --pkg contract \
    --out "$project_dir/tests/contract/${contract_test}.go"
done

rm -rf "$project_dir/contract/artifacts"
