package crosschain

// solc version 0.8.19 https://github.com/ethereum/solidity/releases
//go:generate solc --abi ./solidity/CrossChain.sol -o ./artifacts --overwrite
//go:generate solc --bin ./solidity/CrossChain.sol -o ./artifacts --overwrite
// abigen version 1.11.5-stable https://github.com/ethereum/go-ethereum/releases
//go:generate abigen --abi ./artifacts/CrossChain.abi --bin ./artifacts/CrossChain.bin --type crosschain --pkg crosschain --out ./crosschain_contract.go

// solc version 0.8.19 https://github.com/ethereum/solidity/releases
//go:generate solc --abi ./solidity/crosschain_test.sol -o ./artifacts --overwrite
//go:generate solc --bin ./solidity/crosschain_test.sol -o ./artifacts --overwrite
// abigen version 1.11.5-stable https://github.com/ethereum/go-ethereum/releases
//go:generate abigen --abi ./artifacts/crosschain_test.abi --bin ./artifacts/crosschain_test.bin --type crosschain_test --pkg crosschain_test --out ./crosschain_contract_test.go
