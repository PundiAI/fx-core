# Fx Bridge Contract Upgrade Process

```shell
# setup env
cd solidity

yarn install

yarn typechain

# deploy bridge logic contract
export GOERLI_URL="https://goerli.infura.io/v3/xxxxxxx"

npx hardhat deploy-contract --contract-name FxBridgeLogicETH --is-ledger true --driver-path "m/44'/60'/0'/0/0" --network goerli

# verify bridge logic contract
export ETHERSCAN_API_KEY="xxxxxxx"

npx hardhat verify <new bridge logic address> --network goerli

# upgrade bridge logic contract
npx hardhat send 0xB1B68DFC4eE0A3123B897107aFbF43CEFEE9b0A2 "upgradeTo(address)" <new bridge logic address> --is-ledger true --driver-path "m/44'/60'/0'/0/1" --network goerli
```