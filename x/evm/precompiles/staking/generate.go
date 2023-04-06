package staking

// solc version 0.8.19 https://github.com/ethereum/solidity/releases
//go:generate solc --abi ./solidity/Staking.sol -o ./artifacts --overwrite
//go:generate solc --bin ./solidity/Staking.sol -o ./artifacts --overwrite
// abigen version 1.11.5-stable https://github.com/ethereum/go-ethereum/releases
//go:generate abigen --abi ./artifacts/Staking.abi --bin ./artifacts/Staking.bin --type staking --pkg staking --out ./staking_contract.go

// solc version 0.8.19 https://github.com/ethereum/solidity/releases
//go:generate solc --abi ./solidity/staking_test.sol -o ./artifacts --overwrite
//go:generate solc --bin ./solidity/staking_test.sol -o ./artifacts --overwrite
// abigen version 1.11.5-stable https://github.com/ethereum/go-ethereum/releases
//go:generate abigen --abi ./artifacts/staking_test.abi --bin ./artifacts/staking_test.bin --type staking_test --pkg staking_test --out ./staking_contract_test.go
