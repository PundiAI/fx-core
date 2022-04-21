# Change log

## [Unreleased]

### CLI Breaking Changes

* `fxcored unsafe-reset-all` command has been moved to the `fxcored tendermint` sub-command.
* `fxcored debug addr` command has been moved and rename to the `fxcored keys prase`.
* `fxcored tendermint update-validator` command has been rename to the `fxcored tendermint unsafe-reset-priv-validator`
* `fxcored tendermint update-node-key` command has been rename to the `fxcored tendermint unsafe-reset-node-key`
* Remove bech32 PubKey support, Use pubkey in JSON format
* `fxcored keys add` command flags `--algo` the default is eth_secp256k1; `--coin-type` the default is 60

### API Breaking Changes

* update FX metadata, delete `fx` denom

### Features

* add gRPC swagger-ui

### Bug Fixes

### Deprecated

