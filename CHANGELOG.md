# Change log

## [Unreleased]

### Improvements

* Bump tendermint to v0.34.20.
* Bump cosmos-sdk to v0.45.5.
* The IBC version was upgraded from Cosmos-SDK/x/ibc to IBC-Go v3.1.0
* Added modules: feegrant、authz、feemarket、evm、erc20、migrate
* Migrate modules: auth、bank、distribution、gov、slashing、ibc、crosschain(bsc、polygon、tron)
* The previous Oracle deposit will be automatically delegated to the validator with the highest power value after the upgrade.  Oracle can modify the validator itself, requiring a manual delegate payment
* `MsgRequestBatch` add the field BaseFee
* Delete gravity and crosschain module ibc sequence key 

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

### API Breaking Changes

* Update FX metadata, delete `fx` denom
* Refactor `gravity` and `crosschain` module reset api routes
* The `gravity` and `crosschain` module add `ProjectedBatchTimeoutHeight` and `BridgeTokens` query api
* The `gravity`、`crosschain` and `other` reset API route add `/fx` prefix

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

