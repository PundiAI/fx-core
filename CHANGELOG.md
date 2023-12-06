<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.

Ref: https://keepachangelog.com/en/1.0.0/
-->

# Change log

## [v6.0.0-rc0]

### Features

* Supports precompile staking redelegate
* Supports Convert one-to-one token to many-to-one
* Supports layer2 cross-chain protocol integration

## [v5.0.0]

### Features

* Transfer validator permissions
* Edit consensus public key

### Bug Fixes

* Repair testnet slash period

## [v4.2.1]

### Bug Fixes

* Apply ClawbackVestingAccount Barberry patch & Bump SDK to v0.46.13

## [v4.2.0]

### Bug Fixes

* Fix: IBC transfer to evm
* Fix: fail to withdraw reward after validator is slashed
* Fix(ibc): Properly handle ordered channels in UnreceivedPackets query

## [v4.1.0]

### Bug Fixes

* Fix: WFX Token contract code
* Fix: can not transfer shares when redelegate
* Fix: add gravity module to ibc router

## [v4.0.0]

### Features

* Supports optimism and arbitrum cross-chain protocol integration
* Supports staking through contract calls, with the staking pre-compiled contract address being: 0x0000000000000000000000000000000000001003
* Supports cross-chain transactions through contract calls, with the gravity and IBC pre-compiled contract address being: 0x0000000000000000000000000000000000001004
* Supports adding multiple cross-chain tokens through a single proposal
* Supports contract calls through a proposal
* Supports IBC cross-chain transactions for contract tokens
* Supports cross-chain transactions of contract tokens to other chains
* Added doctor command line tool to check the fxcored working environment 

### Improvements

* Bump go-ethereum version to v1.10.26
* Bump cosmos-sdk to v0.46.12
* Bump tendermint to v0.34.27
* Bump ibc-go to v6.1.0

---

## [v3.1.0] - 2023-01-31

### Bug Fixes

* Fixed emit block bloom event in evm module

## [v3.0.0] - 2023-01-14

### Features

* Support Avalanche C-chain cross-chain
* Support cross-chain tokens: AVAX, SAVAX, QI, BAVA and WBTC (erc20)
* Support IBC standard transfer transaction
* Support the From address as 0x address when IBC cross-chain
* Migrate the gravity module to the eth module (unify all cross-chain logic)
* The fee must be empty when calling the contract Transfercrosschain method for ibc cross-chain
* When the contract self-destructs, the contract code cannot be deleted
* The EthereumTx transaction gas limit must be greater than 0
* The EthereumTx transaction From must be empty
* Upgrade WFX contract to support cross-chain transfer of FX Token contract

### Bug Fixes 

* Fix bridge oracle address delegation invalid
* Fix the bug that the alias field of metadata is set to "null"
* Fix keys command parse address

### Improvements

* Bump go-ethereum version to v1.10.19
* Bump cosmos-sdk to v0.45.11
* Bump tendermint to v0.34.23

## [v2.4.1-2] - 2022-10-14

### Bug Fixes

* Fix ibc app-transfer v3 grpc-gateway path
* Fix Rest API query tx
* Fix `make install`
* Fix v2.1 upgrade migrate event
* Fix parse ed25519 pubkey command
* Import deprecated msg
* Dragonberry Patch

### Improvements

* Bump cosmos-sdk to v0.45.10
* Bump tendermint to v0.34.22

## [v2.4.0] - 2022-10-14

### Bug Fixes

* Fix Dragonberry Patch
* Fix docker images tag

### Features

## [v2.3.1] - 2022-09-13

### Features

* `RegisterERC20Proposal`, `RegisterCoinProposal`, `ToggleTokenConversionProposal`, `UpdateDenomAliasProposal` proposal quorum changed from 40% to 25%
* configurable bypass-min-fee maxGas

## [v2.3.0] - 2022-08-22

### Bug Fixes

* Fix `gravity` module cancel out batch panic

### Features

* (fx/base) Add `GetGasPrice` gRPC query node gas price

### Deprecated

* (fx/other) Deprecate `GasPrice` gRPC query since `other` module will be deleted

## [v2.2.1] - 2022-07-28

### Bug Fixes

* Fix transaction msg `MsgConvertCoin` `MsgConvertERC20` too much gas
* Fix crosschain to ethereum
* Fix tendermint subcommand

## [v2.2.0] - 2022-07-22

### Features

* Add query oracle reward in the crosschain module
* Check fxcored version when synchronizing blocks from scratch
* Add denom many to one support
* Update RegisterCoinProposal support denom many to one
* Add UpdateDenomAliasProposal and MsgConvertDenom

### Improvements

* Add ibc transfer route event
* Add gravity and crosschain attention event claimHash

## [v2.1.1] - 2022-07-11

### Bug Fixes

* Add support for the `x-crisis-skip-assert-invariants` CLI flag to the `start` command 
* CLI parse legacy proposal `InitCrossChainParamsProposal` failed
* Deleted Polygon(USDT) and Tron(USDT) contracts and metadata initialization during migration and upgrade

### Improvements

* Refactor gravity handle FxOriginatedTokenClaim

## [v2.1.0] - 2022-06-29

### Improvements

* Bump tendermint to v0.34.20.
* Bump cosmos-sdk to v0.45.5.
* The IBC version was upgraded from Cosmos-SDK/x/ibc to IBC-Go v3.1.0
* Added modules: feegrant、authz、feemarket、evm、erc20、migrate
* Migrate modules: auth、bank、distribution、gov、slashing、ibc、crosschain(bsc、polygon、tron)
* The previous Oracle deposit will be automatically delegated to the validator with the highest power value after the upgrade.  Oracle can modify the validator itself, requiring a manual delegate payment
* `MsgRequestBatch` add the field BaseFee
* Delete gravity and crosschain module ibc sequence key 
* Update crosschain params AverageBlockTime
* Bump ethermint to v0.16.1-fxcore-v2.0.0-rc3.
* Update block max gas to 30_000_000

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
* Remove `network` command

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
* Fix ibc router is not empty, receive address parse error

### Deprecated

* Remove `network` command
