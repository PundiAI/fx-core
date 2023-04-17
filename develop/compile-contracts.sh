#!/usr/bin/env bash

set -eo pipefail

commands=(jq abigen)
for cmd in "${commands[@]}"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "$cmd command not found, please install $cmd first" && exit 1
  fi
done

if [ ! -d "./solidity/contracts/node_modules" ]; then
  echo "===> Installing node modules"
  (cd ./solidity && yarn install)
fi

if [ -d "./solidity/artifacts" ]; then
  echo "===> Cleaning artifacts"
  (cd ./solidity && yarn clean)
fi

echo "===> Compiling contracts"
(cd ./solidity && yarn compile)

[[ ! -d "./contract/artifacts" ]] && mkdir -p ./contract/artifacts

# add core contracts
contracts=(WFXUpgradable FIP20Upgradable ICrossChain CrossChainTest IStaking StakingTest)
# add 3rd party contracts
contracts+=(ERC1967Proxy ERC721)

for contract in "${contracts[@]}"; do
  echo "===> Ethereum ABI wrapper code generator: $contract"
  file_path=$(find ./solidity/artifacts -name "${contract}.json" -type f)
  jq -c '.abi' "$file_path" > "./contract/artifacts/${contract}.abi"
  jq -r '.bytecode' "$file_path" > "./contract/artifacts/${contract}.bin"
  abigen --abi "./contract/artifacts/${contract}.abi" --bin "./contract/artifacts/${contract}.bin" --type "${contract}" --pkg contract --out "./contract/${contract}.go"
done
