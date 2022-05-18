# Change log

## [Unreleased]

### CLI Breaking Changes

* `fxcored unsafe-reset-all` command has been moved to the `fxcored tendermint` sub-command.
* `fxcored tendermint update-validator` command has been rename to the `fxcored tendermint unsafe-reset-priv-validator`
* `fxcored tendermint update-node-key` command has been rename to the `fxcored tendermint unsafe-reset-node-key`
* Remove bech32 PubKey support, Use pubkey in JSON format
* `fxcored debug addr` command has been moved and rename to the `fxcored keys prase`.
* `fxcored keys add` command flags `--algo` the default is eth_secp256k1; `--coin-type` the default is 60
* `fxcored keys add` command output add the EIP55 address
* Remove Cli flags `--gas-prices` default value
* Change Cli flags `--gas` default value with `80000`
* Change the `fxcored config` command output to lowercase

### Improvements

* Updated Tendermint to v0.34.19;
* Updated Cosmos-sdk to v0.42.x;
* `MsgRequestBatch` add the field BaseFee
      
### API Breaking Changes

* Update FX metadata, delete `fx` denom
* Refactor `gravity` and `crosschain` module reset api routes
* The `gravity` and `crosschain` module add `ProjectedBatchTimeoutHeight` and `BridgeTokens` query api

### Features

* Support evm, enable ethereum compatibility
* Support EIP1559, the initial gas price is 500Gwei
* Account migrate, migrate fx prefix address to 0x prefix address, validator is not supported
* Add gRPC swagger-ui
* The `gravity/crosschain` module support targetIbc `0x` prefix
* Add `fxcored config update` command, only missing parts are added

### Bug Fixes

* Fix --node flag parsing. [issues#22](https://github.com/FunctionX/fx-core/issues/22)
* Fix --output flag parsing. [issues#34](https://github.com/FunctionX/fx-core/issues/34)

### Deprecated

