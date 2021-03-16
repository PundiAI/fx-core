# fxcore

**fxcore** is a blockchain built using Cosmos SDK and Tendermint and created with [Starport](https://github.com/tendermint/starport).
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
      --home string          directory for config and data (default "/Users/pundix055/.fxcore")
      --log_filter strings   The logging filter can discard custom log type (ABCIQuery) (default "")
      --log_format string    The logging format (json|plain) (default "plain")
      --log_level string     The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --trace                print out full stack trace on errors

Use "fxcored [command] --help" for more information about a command.
```

### Example

Local node startup

```
fxcored init --chain-id=fxcore --denom=FX local
fxcored keys add fx1
fxcored add-genesis-account fx1 100000000000000000000000000FX
fxcored gentx fx1 1000000000000000000000000FX --chain-id=fxcore \
    --moniker="fx1-validator" \
    --commission-max-change-rate=0.01 \
    --commission-max-rate=1.0 \
    --commission-rate=0.07 \
    --details="..." \
    --security-contact="..." \
    --website="..."
fxcored collect-gentxs
fxcored start -rpc.laddr tcp://0.0.0.0:26657
```

## Learn more

- [Starport](https://github.com/tendermint/starport)
- [Cosmos SDK documentation](https://docs.cosmos.network)
- [Cosmos SDK Tutorials](https://tutorials.cosmos.network)
- [Discord](https://discord.gg/W8trcGV)
