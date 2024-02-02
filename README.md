# fxcore

**fxcore** is a blockchain built using Cosmos SDK and Tendermint and created with [Starport](https://github.com/tendermint/starport).

[![Version](https://img.shields.io/github/v/release/functionx/fx-core.svg)](https://github.com/functionx/fx-core/releases/latest)
[![API Reference](https://pkg.go.dev/badge/github.com/functionx/fx-core.svg)](https://pkg.go.dev/github.com/functionx/fx-core/v7)
[![License](https://img.shields.io/github/license/functionx/fx-core.svg)](https://github.com/functionx/fx-core/blob/main/LICENSE)
[![Tests](https://github.com/functionx/fx-core/actions/workflows/test.yml/badge.svg)](https://github.com/functionx/fx-core/actions/workflows/test.yml)
[![Lint](https://github.com/functionx/fx-core/actions/workflows/lint.yml/badge.svg)](https://github.com/functionx/fx-core/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/functionx/fx-core/v7)](https://goreportcard.com/report/github.com/functionx/fx-core/v7)

**Note**: Requires [Go 1.21+](https://go.dev/dl)

## Releases

Please do not use the main branch to install `fxcored` or run node. Use releases instead.

## Quick start

### Install

```
make install
```

### Usage

```
FunctionX Core BlockChain App

Usage:
  fxcored [command]

Available Commands:
  add-genesis-account Add a genesis account to genesis.json
  collect-gentxs      Collect genesis txs and output a genesis.json file
  config              Create or query an application CLI configuration file
  data                Modify data or query data in database
  debug               Tool for helping with debugging your application
  doctor              Check your system for potential problems
  export              Export state to JSON
  gentx               Generate a genesis tx carrying a self delegation
  help                Help about any command
  init                Initialize private validator, p2p, genesis, application and client configuration files
  keys                Manage your application's keys
  pre-upgrade         Called by cosmovisor, before migrations upgrade
  prune               Prune app history states by keeping the recent heights and deleting old heights
  query               Querying subcommands
  rollback            rollback cosmos-sdk and tendermint state by one height
  rosetta             spin up a rosetta server
  start               Run the full node
  status              Query remote node for status
  tendermint          Tendermint subcommands
  testnet             Initialize files for a fxcore local testnet
  tx                  Transactions subcommands
  validate-genesis    validates the genesis file at the default location or at the location passed as an arg
  version             Print the application binary version information

Flags:
  -h, --help                help for fxcored
      --home string         directory for config and data (default "/root/.fxcore")
      --log_format string   The logging format (json|plain) (default "plain")
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --trace               print out full stack trace on errors

Use "fxcored [command] --help" for more information about a command.
```

## Learn more

- [Function X Docs](https://functionx.gitbook.io)
- [Tendermint Starport](https://github.com/tendermint/starport)
- [Cosmos SDK Documentation](https://docs.cosmos.network)

## License

This project is licensed under the [Apache 2.0 License](LICENSE).
