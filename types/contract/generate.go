package contract

// abigen version 1.10.2-stable
//go:generate abigen --abi ./artifacts/WFX.abi --bin ./artifacts/WFX.bin --type WFX --pkg contract --out ./wfx.go
//go:generate abigen --abi ./artifacts/FIP20.abi --bin ./artifacts/FIP20.bin --type FIP20 --pkg contract --out ./fip20.go
//go:generate abigen --abi ./artifacts/ERC1967Proxy.abi --bin ./artifacts/ERC1967Proxy.bin --type ERC1967Proxy --pkg contract --out ./erc1967_proxy.go
//go:generate abigen --abi ./artifacts/LPToken.abi --bin ./artifacts/LPToken.bin --type LPToken --pkg contract --out ./lptoken.go
