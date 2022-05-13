# Change log

## [Unreleased]

### CLI Breaking Changes

* `fxcored unsafe-reset-all` command has been moved to the `fxcored tendermint` sub-command.
* `fxcored debug addr` command has been moved and rename to the `fxcored keys prase`.
* `fxcored tendermint update-validator` command has been rename to the `fxcored tendermint unsafe-reset-priv-validator`
* `fxcored tendermint update-node-key` command has been rename to the `fxcored tendermint unsafe-reset-node-key`
* Remove bech32 PubKey support, Use pubkey in JSON format
* `fxcored keys add` command flags `--algo` the default is eth_secp256k1; `--coin-type` the default is 60
* `fxcored keys add` command output add the EIP55 address

### API Breaking Changes

* update FX metadata, delete `fx` denom

### Features

* support evm, enable ethereum compatibility
* support EIP1559, the initial gas price is 500Gwei
* account migrate, migrate fx address to 0x address, validator is not supported
* add gRPC swagger-ui
* gravity/crosschain module support targetIbc `0x` prefix

### Bug Fixes

### Deprecated

