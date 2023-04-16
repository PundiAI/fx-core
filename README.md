# fxcore

**fxcore** is a blockchain built using Cosmos SDK and Tendermint and created with [Starport](https://github.com/tendermint/starport).

[![Version](https://img.shields.io/github/v/release/functionx/fx-core.svg)](https://github.com/functionx/fx-core/releases/latest)
[![API Reference](https://pkg.go.dev/badge/github.com/functionx/fx-core.svg)](https://pkg.go.dev/github.com/functionx/fx-core/v4)
[![License](https://img.shields.io/github/license/functionx/fx-core.svg)](https://github.com/functionx/fx-core/blob/main/LICENSE)
[![Tests](https://github.com/functionx/fx-core/actions/workflows/test.yml/badge.svg)](https://github.com/functionx/fx-core/actions/workflows/test.yml)
[![Lint](https://github.com/functionx/fx-core/actions/workflows/lint.yml/badge.svg)](https://github.com/functionx/fx-core/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/functionx/fx-core/v4)](https://goreportcard.com/report/github.com/functionx/fx-core/v4)

**Note**: Requires [Go 1.18+](https://go.dev/dl)

## Releases

Please do not use the main branch to install `fxcored` or run node. Use releases instead.

## Quick start

### Install

```
make install
```

### Usage

```
FunctionX Core Chain App

Usage:
  fxcored [command]

Available Commands:
  add-genesis-account Add a genesis account to genesis.json
  collect-gentxs      Collect genesis txs and output a genesis.json file
  config              Update or query an application configuration file
  debug               Tool for helping with debugging your application
  export              Export state to JSON
  gentx               Generate a genesis tx carrying a self delegation
  help                Help about any command
  init                Initialize private validator, p2p, genesis, and application configuration files
  keys                Manage your application's keys
  migrate             Migrate genesis to a specified target version
  network             Show fxcored network and upgrade info
  query               Querying subcommands
  start               Run the full node
  status              Query remote node for status
  tendermint          Tendermint subcommands
  testnet             Initialize files for a fxchain testnet
  tx                  Transactions subcommands
  unsafe-reset-all    Resets the blockchain database, removes address book files, and resets data/priv_validator_state.json to the genesis state
  validate-genesis    validates the genesis file at the default location or at the location passed as an arg
  version             Print the application binary version information

Flags:
  -h, --help                 help for fxcored
      --home string          directory for config and data (default "/root/.fxcore")
      --log_filter strings   The logging filter can discard custom log type (ABCIQuery)
      --log_format string    The logging format (json|plain) (default "plain")
      --log_level string     The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --trace                print out full stack trace on errors

Use "fxcored [command] --help" for more information about a command.
```

## Learn more

- [Function X Docs](https://functionx.gitbook.io)
- [Tendermint Starport](https://github.com/tendermint/starport)
- [Cosmos SDK Documentation](https://docs.cosmos.network)
