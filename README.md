# Pundi AIFX

**Pundi AIFX** is a blockchain built using Cosmos SDK and Cometbft.

## Maintenance status (read this first)

**Official maintenance of this repository and the reference `fxcored` binaries ended on March 1, 2026.** There will be no guaranteed security fixes, upgrade coordination, or release support from the core team after that date. A small number of validators may still be online; if you operate one, treat the network as **community-run / at your own risk**.

For a step-by-step **validator exit** walkthrough (unbonding, keys, and shutting down the node), see [docs/validator-exit-guide.md](docs/validator-exit-guide.md).

### If you still run a validator — quick checklist

1. **Decide whether to keep operating.** Without official maintenance, you are responsible for incident response, chain halts, and any emergency patches yourself (or via a community fork).
2. **Plan key and funds safety.** Ensure operator keys, consensus keys, and key backups are stored offline and access-controlled. If you stop validating, follow the normal Cosmos SDK flow for **unbonding** / **undelegation** and respect unbonding time on **mainnet `fxcore`** or **testnet `dhobyghaut`** as applicable.
3. **Stop the node cleanly when you exit.** Disable orchestration (systemd, Kubernetes, cosmovisor, etc.), then shut down `fxcored` so the validator does not keep signing while you intend to leave the active set.
4. **Archive what you need.** Export or retain any accounting, logs, or state snapshots required for your records before decommissioning hardware.
5. **Use tagged releases only — never `main`.** Historical builds and upgrade heights are described in [Releases](https://github.com/pundiai/fx-core/releases); `main` is not a supported production branch.
6. **Docs may be stale.** The former documentation site is [Pundi AIFX Docs](https://pundi.gitbook.io); verify any procedure against on-chain parameters and your own node version.

[![Version](https://img.shields.io/github/v/release/pundiai/fx-core.svg)](https://github.com/pundiai/fx-core/releases/latest)
[![License](https://img.shields.io/github/license/pundiai/fx-core.svg)](https://github.com/pundiai/fx-core/blob/main/LICENSE)
[![API Reference](https://pkg.go.dev/badge/github.com/pundiai/fx-core.svg)](https://pkg.go.dev/github.com/pundiai/fx-core/v8)
[![Go Report Card](https://goreportcard.com/badge/github.com/pundiai/fx-core/v8)](https://goreportcard.com/report/github.com/pundiai/fx-core/v8)
[![Go Tests/Lint](https://github.com/pundiai/fx-core/actions/workflows/golang.yml/badge.svg)](https://github.com/pundiai/fx-core/actions/workflows/golang.yml)
[![CodeQL](https://github.com/pundiai/fx-core/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/pundiai/fx-core/actions/workflows/codeql-analysis.yml)
[![Solidity](https://github.com/pundiai/fx-core/actions/workflows/solidity.yml/badge.svg)](https://github.com/pundiai/fx-core/actions/workflows/solidity.yml)

**Note**: Requires [Go 1.23+](https://go.dev/dl)

## Releases

Please do not use the main branch to install `fxcored` or run node. Use releases instead.

## Quick start

### Install

```
make install
```

### Usage

```
Pundi AIFX BlockChain App

Usage:
  fxcored [command]

Available Commands:
  add-genesis-account Add a genesis account to genesis.json
  collect-gentxs      Collect genesis txs and output a genesis.json file
  comet               CometBFT subcommands
  config              Create or query an application CLI configuration file
  data                Modify data or query data in database
  debug               Tool for helping with debugging your application
  doctor              Check your system for potential problems
  export              Export state to JSON
  export-delegates    Export all delegates and holders
  genesis             Application's genesis-related subcommands
  gentx               Generate a genesis tx carrying a self delegation
  help                Help about any command
  index-eth-tx        Index historical eth txs
  init                Initialize private validator, p2p, genesis, application and client configuration files
  keys                Manage your application's keys
  pre-upgrade         Called by cosmovisor, before migrations upgrade
  prune               Prune app history states by keeping the recent heights and deleting old heights
  query               Querying subcommands
  rollback            rollback cosmos-sdk and tendermint state by one height
  rosetta             spin up a rosetta server
  snapshots           Manage local snapshots
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
      --log_no_color        Disable colored logs
      --trace               print out full stack trace on errors

Use "fxcored [command] --help" for more information about a command.
```

## Learn more

- [Pundi AIFX Docs](https://pundi.gitbook.io)
- [Cosmos SDK Documentation](https://docs.cosmos.network)

## License

This project is licensed under the [Apache 2.0 License](LICENSE).
