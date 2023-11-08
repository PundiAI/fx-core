# Fx Bridge Contract Upgrade Process

```
1. Deploy new bridge logic contract

npx hardhat deploy-contract --contract-name FxBridgeLogicETH --is-ledger --driver-path "m/44'/60'/0'/0/0" --network <network>

2. send upgradeTo

npx hardhat send 0x6f1D09Fed11115d65E1071CD2109eDb300D80A27 "upgradeTo(address)" <new bridge logic address> --driver-path "m/44'/60'/0'/0/0" --network <network>

3. send migrate

npx hardhat send 0x6f1D09Fed11115d65E1071CD2109eDb300D80A27 "migrate()" --driver-path "m/44'/60'/0'/0/0" --network <network>
```