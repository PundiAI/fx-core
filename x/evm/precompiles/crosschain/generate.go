package crosschain

// solc version 0.8.19 https://github.com/ethereum/solidity/releases
//go:generate solc --abi ./crosschain_test.sol -o ./artifacts --overwrite
//go:generate solc --bin ./crosschain_test.sol -o ./artifacts --overwrite
// abigen version 1.11.5-stable https://github.com/ethereum/go-ethereum/releases
//go:generate abigen --abi ./artifacts/crosschain_test.abi --bin ./artifacts/crosschain_test.bin --type crosschain_test --pkg crosschain_test --out ./crosschain_contract_test.go
