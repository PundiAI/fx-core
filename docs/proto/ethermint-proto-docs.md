<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [ethermint/crypto/v1/ethsecp256k1/keys.proto](#ethermint/crypto/v1/ethsecp256k1/keys.proto)
    - [PrivKey](#ethermint.crypto.v1.ethsecp256k1.PrivKey)
    - [PubKey](#ethermint.crypto.v1.ethsecp256k1.PubKey)
  
- [ethermint/evm/v1/events.proto](#ethermint/evm/v1/events.proto)
    - [EventBlockBloom](#ethermint.evm.v1.EventBlockBloom)
    - [EventEthereumTx](#ethermint.evm.v1.EventEthereumTx)
    - [EventMessage](#ethermint.evm.v1.EventMessage)
    - [EventTxLog](#ethermint.evm.v1.EventTxLog)
  
- [ethermint/evm/v1/evm.proto](#ethermint/evm/v1/evm.proto)
    - [AccessTuple](#ethermint.evm.v1.AccessTuple)
    - [ChainConfig](#ethermint.evm.v1.ChainConfig)
    - [Log](#ethermint.evm.v1.Log)
    - [Params](#ethermint.evm.v1.Params)
    - [State](#ethermint.evm.v1.State)
    - [TraceConfig](#ethermint.evm.v1.TraceConfig)
    - [TransactionLogs](#ethermint.evm.v1.TransactionLogs)
    - [TxResult](#ethermint.evm.v1.TxResult)
  
- [ethermint/evm/v1/genesis.proto](#ethermint/evm/v1/genesis.proto)
    - [GenesisAccount](#ethermint.evm.v1.GenesisAccount)
    - [GenesisState](#ethermint.evm.v1.GenesisState)
  
- [ethermint/evm/v1/tx.proto](#ethermint/evm/v1/tx.proto)
    - [AccessListTx](#ethermint.evm.v1.AccessListTx)
    - [DynamicFeeTx](#ethermint.evm.v1.DynamicFeeTx)
    - [ExtensionOptionsEthereumTx](#ethermint.evm.v1.ExtensionOptionsEthereumTx)
    - [LegacyTx](#ethermint.evm.v1.LegacyTx)
    - [MsgEthereumTx](#ethermint.evm.v1.MsgEthereumTx)
    - [MsgEthereumTxResponse](#ethermint.evm.v1.MsgEthereumTxResponse)
    - [MsgUpdateParams](#ethermint.evm.v1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#ethermint.evm.v1.MsgUpdateParamsResponse)
  
    - [Msg](#ethermint.evm.v1.Msg)
  
- [ethermint/evm/v1/query.proto](#ethermint/evm/v1/query.proto)
    - [EstimateGasResponse](#ethermint.evm.v1.EstimateGasResponse)
    - [EthCallRequest](#ethermint.evm.v1.EthCallRequest)
    - [QueryAccountRequest](#ethermint.evm.v1.QueryAccountRequest)
    - [QueryAccountResponse](#ethermint.evm.v1.QueryAccountResponse)
    - [QueryBalanceRequest](#ethermint.evm.v1.QueryBalanceRequest)
    - [QueryBalanceResponse](#ethermint.evm.v1.QueryBalanceResponse)
    - [QueryBaseFeeRequest](#ethermint.evm.v1.QueryBaseFeeRequest)
    - [QueryBaseFeeResponse](#ethermint.evm.v1.QueryBaseFeeResponse)
    - [QueryCodeRequest](#ethermint.evm.v1.QueryCodeRequest)
    - [QueryCodeResponse](#ethermint.evm.v1.QueryCodeResponse)
    - [QueryCosmosAccountRequest](#ethermint.evm.v1.QueryCosmosAccountRequest)
    - [QueryCosmosAccountResponse](#ethermint.evm.v1.QueryCosmosAccountResponse)
    - [QueryParamsRequest](#ethermint.evm.v1.QueryParamsRequest)
    - [QueryParamsResponse](#ethermint.evm.v1.QueryParamsResponse)
    - [QueryStorageRequest](#ethermint.evm.v1.QueryStorageRequest)
    - [QueryStorageResponse](#ethermint.evm.v1.QueryStorageResponse)
    - [QueryTraceBlockRequest](#ethermint.evm.v1.QueryTraceBlockRequest)
    - [QueryTraceBlockResponse](#ethermint.evm.v1.QueryTraceBlockResponse)
    - [QueryTraceTxRequest](#ethermint.evm.v1.QueryTraceTxRequest)
    - [QueryTraceTxResponse](#ethermint.evm.v1.QueryTraceTxResponse)
    - [QueryTxLogsRequest](#ethermint.evm.v1.QueryTxLogsRequest)
    - [QueryTxLogsResponse](#ethermint.evm.v1.QueryTxLogsResponse)
    - [QueryValidatorAccountRequest](#ethermint.evm.v1.QueryValidatorAccountRequest)
    - [QueryValidatorAccountResponse](#ethermint.evm.v1.QueryValidatorAccountResponse)
  
    - [Query](#ethermint.evm.v1.Query)
  
- [ethermint/feemarket/v1/events.proto](#ethermint/feemarket/v1/events.proto)
    - [EventBlockGas](#ethermint.feemarket.v1.EventBlockGas)
    - [EventFeeMarket](#ethermint.feemarket.v1.EventFeeMarket)
  
- [ethermint/feemarket/v1/feemarket.proto](#ethermint/feemarket/v1/feemarket.proto)
    - [Params](#ethermint.feemarket.v1.Params)
  
- [ethermint/feemarket/v1/genesis.proto](#ethermint/feemarket/v1/genesis.proto)
    - [GenesisState](#ethermint.feemarket.v1.GenesisState)
  
- [ethermint/feemarket/v1/query.proto](#ethermint/feemarket/v1/query.proto)
    - [QueryBaseFeeRequest](#ethermint.feemarket.v1.QueryBaseFeeRequest)
    - [QueryBaseFeeResponse](#ethermint.feemarket.v1.QueryBaseFeeResponse)
    - [QueryBlockGasRequest](#ethermint.feemarket.v1.QueryBlockGasRequest)
    - [QueryBlockGasResponse](#ethermint.feemarket.v1.QueryBlockGasResponse)
    - [QueryParamsRequest](#ethermint.feemarket.v1.QueryParamsRequest)
    - [QueryParamsResponse](#ethermint.feemarket.v1.QueryParamsResponse)
  
    - [Query](#ethermint.feemarket.v1.Query)
  
- [ethermint/feemarket/v1/tx.proto](#ethermint/feemarket/v1/tx.proto)
    - [MsgUpdateParams](#ethermint.feemarket.v1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#ethermint.feemarket.v1.MsgUpdateParamsResponse)
  
    - [Msg](#ethermint.feemarket.v1.Msg)
  
- [ethermint/types/v1/account.proto](#ethermint/types/v1/account.proto)
    - [EthAccount](#ethermint.types.v1.EthAccount)
  
- [ethermint/types/v1/dynamic_fee.proto](#ethermint/types/v1/dynamic_fee.proto)
    - [ExtensionOptionDynamicFeeTx](#ethermint.types.v1.ExtensionOptionDynamicFeeTx)
  
- [ethermint/types/v1/indexer.proto](#ethermint/types/v1/indexer.proto)
    - [TxResult](#ethermint.types.v1.TxResult)
  
- [ethermint/types/v1/web3.proto](#ethermint/types/v1/web3.proto)
    - [ExtensionOptionsWeb3Tx](#ethermint.types.v1.ExtensionOptionsWeb3Tx)
  
- [Scalar Value Types](#scalar-value-types)



<a name="ethermint/crypto/v1/ethsecp256k1/keys.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/crypto/v1/ethsecp256k1/keys.proto



<a name="ethermint.crypto.v1.ethsecp256k1.PrivKey"></a>

### PrivKey
PrivKey defines a type alias for an ecdsa.PrivateKey that implements
Tendermint's PrivateKey interface.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  | key is the private key in byte form |






<a name="ethermint.crypto.v1.ethsecp256k1.PubKey"></a>

### PubKey
PubKey defines a type alias for an ecdsa.PublicKey that implements
Tendermint's PubKey interface. It represents the 33-byte compressed public
key format.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  | key is the public key in byte form |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/evm/v1/events.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/evm/v1/events.proto



<a name="ethermint.evm.v1.EventBlockBloom"></a>

### EventBlockBloom
EventBlockBloom defines an Ethereum block bloom filter event


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bloom` | [string](#string) |  | bloom is the bloom filter of the block |






<a name="ethermint.evm.v1.EventEthereumTx"></a>

### EventEthereumTx
EventEthereumTx defines the event for an Ethereum transaction


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `amount` | [string](#string) |  | amount |
| `eth_hash` | [string](#string) |  | eth_hash is the Ethereum hash of the transaction |
| `index` | [string](#string) |  | index of the transaction in the block |
| `gas_used` | [string](#string) |  | gas_used is the amount of gas used by the transaction |
| `hash` | [string](#string) |  | hash is the Tendermint hash of the transaction |
| `recipient` | [string](#string) |  | recipient of the transaction |
| `eth_tx_failed` | [string](#string) |  | eth_tx_failed contains a VM error should it occur |






<a name="ethermint.evm.v1.EventMessage"></a>

### EventMessage
EventMessage


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `module` | [string](#string) |  | module which emits the event |
| `sender` | [string](#string) |  | sender of the message |
| `tx_type` | [string](#string) |  | tx_type is the type of the message |






<a name="ethermint.evm.v1.EventTxLog"></a>

### EventTxLog
EventTxLog defines the event for an Ethereum transaction log


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tx_logs` | [string](#string) | repeated | tx_logs is an array of transaction logs |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/evm/v1/evm.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/evm/v1/evm.proto



<a name="ethermint.evm.v1.AccessTuple"></a>

### AccessTuple
AccessTuple is the element type of an access list.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is a hex formatted ethereum address |
| `storage_keys` | [string](#string) | repeated | storage_keys are hex formatted hashes of the storage keys |






<a name="ethermint.evm.v1.ChainConfig"></a>

### ChainConfig
ChainConfig defines the Ethereum ChainConfig parameters using *sdk.Int values
instead of *big.Int.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `homestead_block` | [string](#string) |  | homestead_block switch (nil no fork, 0 = already homestead) |
| `dao_fork_block` | [string](#string) |  | dao_fork_block corresponds to TheDAO hard-fork switch block (nil no fork) |
| `dao_fork_support` | [bool](#bool) |  | dao_fork_support defines whether the nodes supports or opposes the DAO hard-fork |
| `eip150_block` | [string](#string) |  | eip150_block: EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150) EIP150 HF block (nil no fork) |
| `eip150_hash` | [string](#string) |  | eip150_hash: EIP150 HF hash (needed for header only clients as only gas pricing changed) |
| `eip155_block` | [string](#string) |  | eip155_block: EIP155Block HF block |
| `eip158_block` | [string](#string) |  | eip158_block: EIP158 HF block |
| `byzantium_block` | [string](#string) |  | byzantium_block: Byzantium switch block (nil no fork, 0 = already on byzantium) |
| `constantinople_block` | [string](#string) |  | constantinople_block: Constantinople switch block (nil no fork, 0 = already activated) |
| `petersburg_block` | [string](#string) |  | petersburg_block: Petersburg switch block (nil same as Constantinople) |
| `istanbul_block` | [string](#string) |  | istanbul_block: Istanbul switch block (nil no fork, 0 = already on istanbul) |
| `muir_glacier_block` | [string](#string) |  | muir_glacier_block: Eip-2384 (bomb delay) switch block (nil no fork, 0 = already activated) |
| `berlin_block` | [string](#string) |  | berlin_block: Berlin switch block (nil = no fork, 0 = already on berlin) |
| `london_block` | [string](#string) |  | london_block: London switch block (nil = no fork, 0 = already on london) |
| `arrow_glacier_block` | [string](#string) |  | arrow_glacier_block: Eip-4345 (bomb delay) switch block (nil = no fork, 0 = already activated) |
| `gray_glacier_block` | [string](#string) |  | gray_glacier_block: EIP-5133 (bomb delay) switch block (nil = no fork, 0 = already activated) |
| `merge_netsplit_block` | [string](#string) |  | merge_netsplit_block: Virtual fork after The Merge to use as a network splitter |
| `shanghai_block` | [string](#string) |  | shanghai_block switch block (nil = no fork, 0 = already on shanghai) |
| `cancun_block` | [string](#string) |  | cancun_block switch block (nil = no fork, 0 = already on cancun) |






<a name="ethermint.evm.v1.Log"></a>

### Log
Log represents an protobuf compatible Ethereum Log that defines a contract
log event. These events are generated by the LOG opcode and stored/indexed by
the node.

NOTE: address, topics and data are consensus fields. The rest of the fields
are derived, i.e. filled in by the nodes, but not secured by consensus.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address of the contract that generated the event |
| `topics` | [string](#string) | repeated | topics is a list of topics provided by the contract. |
| `data` | [bytes](#bytes) |  | data which is supplied by the contract, usually ABI-encoded |
| `block_number` | [uint64](#uint64) |  | block_number of the block in which the transaction was included |
| `tx_hash` | [string](#string) |  | tx_hash is the transaction hash |
| `tx_index` | [uint64](#uint64) |  | tx_index of the transaction in the block |
| `block_hash` | [string](#string) |  | block_hash of the block in which the transaction was included |
| `index` | [uint64](#uint64) |  | index of the log in the block |
| `removed` | [bool](#bool) |  | removed is true if this log was reverted due to a chain reorganisation. You must pay attention to this field if you receive logs through a filter query. |






<a name="ethermint.evm.v1.Params"></a>

### Params
Params defines the EVM module parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `evm_denom` | [string](#string) |  | evm_denom represents the token denomination used to run the EVM state transitions. |
| `enable_create` | [bool](#bool) |  | enable_create toggles state transitions that use the vm.Create function |
| `enable_call` | [bool](#bool) |  | enable_call toggles state transitions that use the vm.Call function |
| `extra_eips` | [int64](#int64) | repeated | extra_eips defines the additional EIPs for the vm.Config |
| `chain_config` | [ChainConfig](#ethermint.evm.v1.ChainConfig) |  | chain_config defines the EVM chain configuration parameters |
| `allow_unprotected_txs` | [bool](#bool) |  | allow_unprotected_txs defines if replay-protected (i.e non EIP155 signed) transactions can be executed on the state machine. |






<a name="ethermint.evm.v1.State"></a>

### State
State represents a single Storage key value pair item.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  | key is the stored key |
| `value` | [string](#string) |  | value is the stored value for the given key |






<a name="ethermint.evm.v1.TraceConfig"></a>

### TraceConfig
TraceConfig holds extra parameters to trace functions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tracer` | [string](#string) |  | tracer is a custom javascript tracer |
| `timeout` | [string](#string) |  | timeout overrides the default timeout of 5 seconds for JavaScript-based tracing calls |
| `reexec` | [uint64](#uint64) |  | reexec defines the number of blocks the tracer is willing to go back |
| `disable_stack` | [bool](#bool) |  | disable_stack switches stack capture |
| `disable_storage` | [bool](#bool) |  | disable_storage switches storage capture |
| `debug` | [bool](#bool) |  | debug can be used to print output during capture end |
| `limit` | [int32](#int32) |  | limit defines the maximum length of output, but zero means unlimited |
| `overrides` | [ChainConfig](#ethermint.evm.v1.ChainConfig) |  | overrides can be used to execute a trace using future fork rules |
| `enable_memory` | [bool](#bool) |  | enable_memory switches memory capture |
| `enable_return_data` | [bool](#bool) |  | enable_return_data switches the capture of return data |
| `tracer_json_config` | [string](#string) |  | tracer_json_config configures the tracer using a JSON string |






<a name="ethermint.evm.v1.TransactionLogs"></a>

### TransactionLogs
TransactionLogs define the logs generated from a transaction execution
with a given hash. It it used for import/export data as transactions are not
persisted on blockchain state after an upgrade.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  | hash of the transaction |
| `logs` | [Log](#ethermint.evm.v1.Log) | repeated | logs is an array of Logs for the given transaction hash |






<a name="ethermint.evm.v1.TxResult"></a>

### TxResult
TxResult stores results of Tx execution.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  | contract_address contains the ethereum address of the created contract (if any). If the state transition is an evm.Call, the contract address will be empty. |
| `bloom` | [bytes](#bytes) |  | bloom represents the bloom filter bytes |
| `tx_logs` | [TransactionLogs](#ethermint.evm.v1.TransactionLogs) |  | tx_logs contains the transaction hash and the proto-compatible ethereum logs. |
| `ret` | [bytes](#bytes) |  | ret defines the bytes from the execution. |
| `reverted` | [bool](#bool) |  | reverted flag is set to true when the call has been reverted |
| `gas_used` | [uint64](#uint64) |  | gas_used notes the amount of gas consumed while execution |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/evm/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/evm/v1/genesis.proto



<a name="ethermint.evm.v1.GenesisAccount"></a>

### GenesisAccount
GenesisAccount defines an account to be initialized in the genesis state.
Its main difference between with Geth's GenesisAccount is that it uses a
custom storage type and that it doesn't contain the private key field.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address defines an ethereum hex formated address of an account |
| `code` | [string](#string) |  | code defines the hex bytes of the account code. |
| `storage` | [State](#ethermint.evm.v1.State) | repeated | storage defines the set of state key values for the account. |






<a name="ethermint.evm.v1.GenesisState"></a>

### GenesisState
GenesisState defines the evm module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `accounts` | [GenesisAccount](#ethermint.evm.v1.GenesisAccount) | repeated | accounts is an array containing the ethereum genesis accounts. |
| `params` | [Params](#ethermint.evm.v1.Params) |  | params defines all the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/evm/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/evm/v1/tx.proto



<a name="ethermint.evm.v1.AccessListTx"></a>

### AccessListTx
AccessListTx is the data of EIP-2930 access list transactions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  | chain_id of the destination EVM chain |
| `nonce` | [uint64](#uint64) |  | nonce corresponds to the account nonce (transaction sequence). |
| `gas_price` | [string](#string) |  | gas_price defines the value for each gas unit |
| `gas` | [uint64](#uint64) |  | gas defines the gas limit defined for the transaction. |
| `to` | [string](#string) |  | to is the recipient address in hex format |
| `value` | [string](#string) |  | value defines the unsigned integer value of the transaction amount. |
| `data` | [bytes](#bytes) |  | data is the data payload bytes of the transaction. |
| `accesses` | [AccessTuple](#ethermint.evm.v1.AccessTuple) | repeated | accesses is an array of access tuples |
| `v` | [bytes](#bytes) |  | v defines the signature value |
| `r` | [bytes](#bytes) |  | r defines the signature value |
| `s` | [bytes](#bytes) |  | s define the signature value |






<a name="ethermint.evm.v1.DynamicFeeTx"></a>

### DynamicFeeTx
DynamicFeeTx is the data of EIP-1559 dinamic fee transactions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  | chain_id of the destination EVM chain |
| `nonce` | [uint64](#uint64) |  | nonce corresponds to the account nonce (transaction sequence). |
| `gas_tip_cap` | [string](#string) |  | gas_tip_cap defines the max value for the gas tip |
| `gas_fee_cap` | [string](#string) |  | gas_fee_cap defines the max value for the gas fee |
| `gas` | [uint64](#uint64) |  | gas defines the gas limit defined for the transaction. |
| `to` | [string](#string) |  | to is the hex formatted address of the recipient |
| `value` | [string](#string) |  | value defines the the transaction amount. |
| `data` | [bytes](#bytes) |  | data is the data payload bytes of the transaction. |
| `accesses` | [AccessTuple](#ethermint.evm.v1.AccessTuple) | repeated | accesses is an array of access tuples |
| `v` | [bytes](#bytes) |  | v defines the signature value |
| `r` | [bytes](#bytes) |  | r defines the signature value |
| `s` | [bytes](#bytes) |  | s define the signature value |






<a name="ethermint.evm.v1.ExtensionOptionsEthereumTx"></a>

### ExtensionOptionsEthereumTx
ExtensionOptionsEthereumTx is an extension option for ethereum transactions






<a name="ethermint.evm.v1.LegacyTx"></a>

### LegacyTx
LegacyTx is the transaction data of regular Ethereum transactions.
NOTE: All non-protected transactions (i.e non EIP155 signed) will fail if the
AllowUnprotectedTxs parameter is disabled.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  | nonce corresponds to the account nonce (transaction sequence). |
| `gas_price` | [string](#string) |  | gas_price defines the value for each gas unit |
| `gas` | [uint64](#uint64) |  | gas defines the gas limit defined for the transaction. |
| `to` | [string](#string) |  | to is the hex formatted address of the recipient |
| `value` | [string](#string) |  | value defines the unsigned integer value of the transaction amount. |
| `data` | [bytes](#bytes) |  | data is the data payload bytes of the transaction. |
| `v` | [bytes](#bytes) |  | v defines the signature value |
| `r` | [bytes](#bytes) |  | r defines the signature value |
| `s` | [bytes](#bytes) |  | s define the signature value |






<a name="ethermint.evm.v1.MsgEthereumTx"></a>

### MsgEthereumTx
MsgEthereumTx encapsulates an Ethereum transaction as an SDK message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [google.protobuf.Any](#google.protobuf.Any) |  | data is inner transaction data of the Ethereum transaction |
| `size` | [double](#double) |  | size is the encoded storage size of the transaction (DEPRECATED) |
| `hash` | [string](#string) |  | hash of the transaction in hex format |
| `from` | [string](#string) |  | from is the ethereum signer address in hex format. This address value is checked against the address derived from the signature (V, R, S) using the secp256k1 elliptic curve |






<a name="ethermint.evm.v1.MsgEthereumTxResponse"></a>

### MsgEthereumTxResponse
MsgEthereumTxResponse defines the Msg/EthereumTx response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  | hash of the ethereum transaction in hex format. This hash differs from the Tendermint sha256 hash of the transaction bytes. See https://github.com/tendermint/tendermint/issues/6539 for reference |
| `logs` | [Log](#ethermint.evm.v1.Log) | repeated | logs contains the transaction hash and the proto-compatible ethereum logs. |
| `ret` | [bytes](#bytes) |  | ret is the returned data from evm function (result or data supplied with revert opcode) |
| `vm_error` | [string](#string) |  | vm_error is the error returned by vm execution |
| `gas_used` | [uint64](#uint64) |  | gas_used specifies how much gas was consumed by the transaction |






<a name="ethermint.evm.v1.MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams defines a Msg for updating the x/evm module parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address of the governance account. |
| `params` | [Params](#ethermint.evm.v1.Params) |  | params defines the x/evm parameters to update. NOTE: All parameters must be supplied. |






<a name="ethermint.evm.v1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse
MsgUpdateParamsResponse defines the response structure for executing a
MsgUpdateParams message.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ethermint.evm.v1.Msg"></a>

### Msg
Msg defines the evm Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `EthereumTx` | [MsgEthereumTx](#ethermint.evm.v1.MsgEthereumTx) | [MsgEthereumTxResponse](#ethermint.evm.v1.MsgEthereumTxResponse) | EthereumTx defines a method submitting Ethereum transactions. | POST|/ethermint/evm/v1/ethereum_tx|
| `UpdateParams` | [MsgUpdateParams](#ethermint.evm.v1.MsgUpdateParams) | [MsgUpdateParamsResponse](#ethermint.evm.v1.MsgUpdateParamsResponse) | UpdateParams defined a governance operation for updating the x/evm module parameters. The authority is hard-coded to the Cosmos SDK x/gov module account | |

 <!-- end services -->



<a name="ethermint/evm/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/evm/v1/query.proto



<a name="ethermint.evm.v1.EstimateGasResponse"></a>

### EstimateGasResponse
EstimateGasResponse defines EstimateGas response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas` | [uint64](#uint64) |  | gas returns the estimated gas |






<a name="ethermint.evm.v1.EthCallRequest"></a>

### EthCallRequest
EthCallRequest defines EthCall request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `args` | [bytes](#bytes) |  | args uses the same json format as the json rpc api. |
| `gas_cap` | [uint64](#uint64) |  | gas_cap defines the default gas cap to be used |
| `proposer_address` | [bytes](#bytes) |  | proposer_address of the requested block in hex format |
| `chain_id` | [int64](#int64) |  | chain_id is the eip155 chain id parsed from the requested block header |






<a name="ethermint.evm.v1.QueryAccountRequest"></a>

### QueryAccountRequest
QueryAccountRequest is the request type for the Query/Account RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the ethereum hex address to query the account for. |






<a name="ethermint.evm.v1.QueryAccountResponse"></a>

### QueryAccountResponse
QueryAccountResponse is the response type for the Query/Account RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balance` | [string](#string) |  | balance is the balance of the EVM denomination. |
| `code_hash` | [string](#string) |  | code_hash is the hex-formatted code bytes from the EOA. |
| `nonce` | [uint64](#uint64) |  | nonce is the account's sequence number. |






<a name="ethermint.evm.v1.QueryBalanceRequest"></a>

### QueryBalanceRequest
QueryBalanceRequest is the request type for the Query/Balance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the ethereum hex address to query the balance for. |






<a name="ethermint.evm.v1.QueryBalanceResponse"></a>

### QueryBalanceResponse
QueryBalanceResponse is the response type for the Query/Balance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balance` | [string](#string) |  | balance is the balance of the EVM denomination. |






<a name="ethermint.evm.v1.QueryBaseFeeRequest"></a>

### QueryBaseFeeRequest
QueryBaseFeeRequest defines the request type for querying the EIP1559 base
fee.






<a name="ethermint.evm.v1.QueryBaseFeeResponse"></a>

### QueryBaseFeeResponse
QueryBaseFeeResponse returns the EIP1559 base fee.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_fee` | [string](#string) |  | base_fee is the EIP1559 base fee |






<a name="ethermint.evm.v1.QueryCodeRequest"></a>

### QueryCodeRequest
QueryCodeRequest is the request type for the Query/Code RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the ethereum hex address to query the code for. |






<a name="ethermint.evm.v1.QueryCodeResponse"></a>

### QueryCodeResponse
QueryCodeResponse is the response type for the Query/Code RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code` | [bytes](#bytes) |  | code represents the code bytes from an ethereum address. |






<a name="ethermint.evm.v1.QueryCosmosAccountRequest"></a>

### QueryCosmosAccountRequest
QueryCosmosAccountRequest is the request type for the Query/CosmosAccount RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the ethereum hex address to query the account for. |






<a name="ethermint.evm.v1.QueryCosmosAccountResponse"></a>

### QueryCosmosAccountResponse
QueryCosmosAccountResponse is the response type for the Query/CosmosAccount
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cosmos_address` | [string](#string) |  | cosmos_address is the cosmos address of the account. |
| `sequence` | [uint64](#uint64) |  | sequence is the account's sequence number. |
| `account_number` | [uint64](#uint64) |  | account_number is the account number |






<a name="ethermint.evm.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest defines the request type for querying x/evm parameters.






<a name="ethermint.evm.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse defines the response type for querying x/evm parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#ethermint.evm.v1.Params) |  | params define the evm module parameters. |






<a name="ethermint.evm.v1.QueryStorageRequest"></a>

### QueryStorageRequest
QueryStorageRequest is the request type for the Query/Storage RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the ethereum hex address to query the storage state for. |
| `key` | [string](#string) |  | key defines the key of the storage state |






<a name="ethermint.evm.v1.QueryStorageResponse"></a>

### QueryStorageResponse
QueryStorageResponse is the response type for the Query/Storage RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [string](#string) |  | value defines the storage state value hash associated with the given key. |






<a name="ethermint.evm.v1.QueryTraceBlockRequest"></a>

### QueryTraceBlockRequest
QueryTraceBlockRequest defines TraceTx request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `txs` | [MsgEthereumTx](#ethermint.evm.v1.MsgEthereumTx) | repeated | txs is an array of messages in the block |
| `trace_config` | [TraceConfig](#ethermint.evm.v1.TraceConfig) |  | trace_config holds extra parameters to trace functions. |
| `block_number` | [int64](#int64) |  | block_number of the traced block |
| `block_hash` | [string](#string) |  | block_hash (hex) of the traced block |
| `block_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | block_time of the traced block |
| `proposer_address` | [bytes](#bytes) |  | proposer_address is the address of the requested block |
| `chain_id` | [int64](#int64) |  | chain_id is the eip155 chain id parsed from the requested block header |






<a name="ethermint.evm.v1.QueryTraceBlockResponse"></a>

### QueryTraceBlockResponse
QueryTraceBlockResponse defines TraceBlock response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | data is the response serialized in bytes |






<a name="ethermint.evm.v1.QueryTraceTxRequest"></a>

### QueryTraceTxRequest
QueryTraceTxRequest defines TraceTx request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg` | [MsgEthereumTx](#ethermint.evm.v1.MsgEthereumTx) |  | msg is the MsgEthereumTx for the requested transaction |
| `trace_config` | [TraceConfig](#ethermint.evm.v1.TraceConfig) |  | trace_config holds extra parameters to trace functions. |
| `predecessors` | [MsgEthereumTx](#ethermint.evm.v1.MsgEthereumTx) | repeated | predecessors is an array of transactions included in the same block need to be replayed first to get correct context for tracing. |
| `block_number` | [int64](#int64) |  | block_number of requested transaction |
| `block_hash` | [string](#string) |  | block_hash of requested transaction |
| `block_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | block_time of requested transaction |
| `proposer_address` | [bytes](#bytes) |  | proposer_address is the proposer of the requested block |
| `chain_id` | [int64](#int64) |  | chain_id is the the eip155 chain id parsed from the requested block header |






<a name="ethermint.evm.v1.QueryTraceTxResponse"></a>

### QueryTraceTxResponse
QueryTraceTxResponse defines TraceTx response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | data is the response serialized in bytes |






<a name="ethermint.evm.v1.QueryTxLogsRequest"></a>

### QueryTxLogsRequest
QueryTxLogsRequest is the request type for the Query/TxLogs RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  | hash is the ethereum transaction hex hash to query the logs for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="ethermint.evm.v1.QueryTxLogsResponse"></a>

### QueryTxLogsResponse
QueryTxLogsResponse is the response type for the Query/TxLogs RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `logs` | [Log](#ethermint.evm.v1.Log) | repeated | logs represents the ethereum logs generated from the given transaction. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="ethermint.evm.v1.QueryValidatorAccountRequest"></a>

### QueryValidatorAccountRequest
QueryValidatorAccountRequest is the request type for the
Query/ValidatorAccount RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cons_address` | [string](#string) |  | cons_address is the validator cons address to query the account for. |






<a name="ethermint.evm.v1.QueryValidatorAccountResponse"></a>

### QueryValidatorAccountResponse
QueryValidatorAccountResponse is the response type for the
Query/ValidatorAccount RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account_address` | [string](#string) |  | account_address is the cosmos address of the account in bech32 format. |
| `sequence` | [uint64](#uint64) |  | sequence is the account's sequence number. |
| `account_number` | [uint64](#uint64) |  | account_number is the account number |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ethermint.evm.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Account` | [QueryAccountRequest](#ethermint.evm.v1.QueryAccountRequest) | [QueryAccountResponse](#ethermint.evm.v1.QueryAccountResponse) | Account queries an Ethereum account. | GET|/ethermint/evm/v1/account/{address}|
| `CosmosAccount` | [QueryCosmosAccountRequest](#ethermint.evm.v1.QueryCosmosAccountRequest) | [QueryCosmosAccountResponse](#ethermint.evm.v1.QueryCosmosAccountResponse) | CosmosAccount queries an Ethereum account's Cosmos Address. | GET|/ethermint/evm/v1/cosmos_account/{address}|
| `ValidatorAccount` | [QueryValidatorAccountRequest](#ethermint.evm.v1.QueryValidatorAccountRequest) | [QueryValidatorAccountResponse](#ethermint.evm.v1.QueryValidatorAccountResponse) | ValidatorAccount queries an Ethereum account's from a validator consensus Address. | GET|/ethermint/evm/v1/validator_account/{cons_address}|
| `Balance` | [QueryBalanceRequest](#ethermint.evm.v1.QueryBalanceRequest) | [QueryBalanceResponse](#ethermint.evm.v1.QueryBalanceResponse) | Balance queries the balance of a the EVM denomination for a single EthAccount. | GET|/ethermint/evm/v1/balances/{address}|
| `Storage` | [QueryStorageRequest](#ethermint.evm.v1.QueryStorageRequest) | [QueryStorageResponse](#ethermint.evm.v1.QueryStorageResponse) | Storage queries the balance of all coins for a single account. | GET|/ethermint/evm/v1/storage/{address}/{key}|
| `Code` | [QueryCodeRequest](#ethermint.evm.v1.QueryCodeRequest) | [QueryCodeResponse](#ethermint.evm.v1.QueryCodeResponse) | Code queries the balance of all coins for a single account. | GET|/ethermint/evm/v1/codes/{address}|
| `Params` | [QueryParamsRequest](#ethermint.evm.v1.QueryParamsRequest) | [QueryParamsResponse](#ethermint.evm.v1.QueryParamsResponse) | Params queries the parameters of x/evm module. | GET|/ethermint/evm/v1/params|
| `EthCall` | [EthCallRequest](#ethermint.evm.v1.EthCallRequest) | [MsgEthereumTxResponse](#ethermint.evm.v1.MsgEthereumTxResponse) | EthCall implements the `eth_call` rpc api | GET|/ethermint/evm/v1/eth_call|
| `EstimateGas` | [EthCallRequest](#ethermint.evm.v1.EthCallRequest) | [EstimateGasResponse](#ethermint.evm.v1.EstimateGasResponse) | EstimateGas implements the `eth_estimateGas` rpc api | GET|/ethermint/evm/v1/estimate_gas|
| `TraceTx` | [QueryTraceTxRequest](#ethermint.evm.v1.QueryTraceTxRequest) | [QueryTraceTxResponse](#ethermint.evm.v1.QueryTraceTxResponse) | TraceTx implements the `debug_traceTransaction` rpc api | GET|/ethermint/evm/v1/trace_tx|
| `TraceBlock` | [QueryTraceBlockRequest](#ethermint.evm.v1.QueryTraceBlockRequest) | [QueryTraceBlockResponse](#ethermint.evm.v1.QueryTraceBlockResponse) | TraceBlock implements the `debug_traceBlockByNumber` and `debug_traceBlockByHash` rpc api | GET|/ethermint/evm/v1/trace_block|
| `BaseFee` | [QueryBaseFeeRequest](#ethermint.evm.v1.QueryBaseFeeRequest) | [QueryBaseFeeResponse](#ethermint.evm.v1.QueryBaseFeeResponse) | BaseFee queries the base fee of the parent block of the current block, it's similar to feemarket module's method, but also checks london hardfork status. | GET|/ethermint/evm/v1/base_fee|

 <!-- end services -->



<a name="ethermint/feemarket/v1/events.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/feemarket/v1/events.proto



<a name="ethermint.feemarket.v1.EventBlockGas"></a>

### EventBlockGas
EventBlockGas defines an Ethereum block gas event


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [string](#string) |  | height of the block |
| `amount` | [string](#string) |  | amount of gas wanted by the block |






<a name="ethermint.feemarket.v1.EventFeeMarket"></a>

### EventFeeMarket
EventFeeMarket is the event type for the fee market module


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_fee` | [string](#string) |  | base_fee for EIP-1559 blocks |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/feemarket/v1/feemarket.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/feemarket/v1/feemarket.proto



<a name="ethermint.feemarket.v1.Params"></a>

### Params
Params defines the EVM module parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `no_base_fee` | [bool](#bool) |  | no_base_fee forces the EIP-1559 base fee to 0 (needed for 0 price calls) |
| `base_fee_change_denominator` | [uint32](#uint32) |  | base_fee_change_denominator bounds the amount the base fee can change between blocks. |
| `elasticity_multiplier` | [uint32](#uint32) |  | elasticity_multiplier bounds the maximum gas limit an EIP-1559 block may have. |
| `enable_height` | [int64](#int64) |  | enable_height defines at which block height the base fee calculation is enabled. |
| `base_fee` | [string](#string) |  | base_fee for EIP-1559 blocks. |
| `min_gas_price` | [string](#string) |  | min_gas_price defines the minimum gas price value for cosmos and eth transactions |
| `min_gas_multiplier` | [string](#string) |  | min_gas_multiplier bounds the minimum gas used to be charged to senders based on gas limit |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/feemarket/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/feemarket/v1/genesis.proto



<a name="ethermint.feemarket.v1.GenesisState"></a>

### GenesisState
GenesisState defines the feemarket module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#ethermint.feemarket.v1.Params) |  | params defines all the parameters of the feemarket module. |
| `block_gas` | [uint64](#uint64) |  | block_gas is the amount of gas wanted on the last block before the upgrade. Zero by default. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/feemarket/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/feemarket/v1/query.proto



<a name="ethermint.feemarket.v1.QueryBaseFeeRequest"></a>

### QueryBaseFeeRequest
QueryBaseFeeRequest defines the request type for querying the EIP1559 base
fee.






<a name="ethermint.feemarket.v1.QueryBaseFeeResponse"></a>

### QueryBaseFeeResponse
QueryBaseFeeResponse returns the EIP1559 base fee.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_fee` | [string](#string) |  | base_fee is the EIP1559 base fee |






<a name="ethermint.feemarket.v1.QueryBlockGasRequest"></a>

### QueryBlockGasRequest
QueryBlockGasRequest defines the request type for querying the EIP1559 base
fee.






<a name="ethermint.feemarket.v1.QueryBlockGasResponse"></a>

### QueryBlockGasResponse
QueryBlockGasResponse returns block gas used for a given height.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas` | [int64](#int64) |  | gas is the returned block gas |






<a name="ethermint.feemarket.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest defines the request type for querying x/evm parameters.






<a name="ethermint.feemarket.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse defines the response type for querying x/evm parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#ethermint.feemarket.v1.Params) |  | params define the evm module parameters. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ethermint.feemarket.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#ethermint.feemarket.v1.QueryParamsRequest) | [QueryParamsResponse](#ethermint.feemarket.v1.QueryParamsResponse) | Params queries the parameters of x/feemarket module. | GET|/ethermint/feemarket/v1/params|
| `BaseFee` | [QueryBaseFeeRequest](#ethermint.feemarket.v1.QueryBaseFeeRequest) | [QueryBaseFeeResponse](#ethermint.feemarket.v1.QueryBaseFeeResponse) | BaseFee queries the base fee of the parent block of the current block. | GET|/ethermint/feemarket/v1/base_fee|
| `BlockGas` | [QueryBlockGasRequest](#ethermint.feemarket.v1.QueryBlockGasRequest) | [QueryBlockGasResponse](#ethermint.feemarket.v1.QueryBlockGasResponse) | BlockGas queries the gas used at a given block height | GET|/ethermint/feemarket/v1/block_gas|

 <!-- end services -->



<a name="ethermint/feemarket/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/feemarket/v1/tx.proto



<a name="ethermint.feemarket.v1.MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams defines a Msg for updating the x/feemarket module parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address of the governance account. |
| `params` | [Params](#ethermint.feemarket.v1.Params) |  | params defines the x/feemarket parameters to update. NOTE: All parameters must be supplied. |






<a name="ethermint.feemarket.v1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse
MsgUpdateParamsResponse defines the response structure for executing a
MsgUpdateParams message.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="ethermint.feemarket.v1.Msg"></a>

### Msg
Msg defines the erc20 Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `UpdateParams` | [MsgUpdateParams](#ethermint.feemarket.v1.MsgUpdateParams) | [MsgUpdateParamsResponse](#ethermint.feemarket.v1.MsgUpdateParamsResponse) | UpdateParams defined a governance operation for updating the x/feemarket module parameters. The authority is hard-coded to the Cosmos SDK x/gov module account | |

 <!-- end services -->



<a name="ethermint/types/v1/account.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/types/v1/account.proto



<a name="ethermint.types.v1.EthAccount"></a>

### EthAccount
EthAccount implements the authtypes.AccountI interface and embeds an
authtypes.BaseAccount type. It is compatible with the auth AccountKeeper.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_account` | [cosmos.auth.v1beta1.BaseAccount](#cosmos.auth.v1beta1.BaseAccount) |  | base_account is an authtypes.BaseAccount |
| `code_hash` | [string](#string) |  | code_hash is the hash calculated from the code contents |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/types/v1/dynamic_fee.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/types/v1/dynamic_fee.proto



<a name="ethermint.types.v1.ExtensionOptionDynamicFeeTx"></a>

### ExtensionOptionDynamicFeeTx
ExtensionOptionDynamicFeeTx is an extension option that specifies the maxPrioPrice for cosmos tx


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `max_priority_price` | [string](#string) |  | max_priority_price is the same as `max_priority_fee_per_gas` in eip-1559 spec |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/types/v1/indexer.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/types/v1/indexer.proto



<a name="ethermint.types.v1.TxResult"></a>

### TxResult
TxResult is the value stored in eth tx indexer


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `height` | [int64](#int64) |  | height of the blockchain |
| `tx_index` | [uint32](#uint32) |  | tx_index of the cosmos transaction |
| `msg_index` | [uint32](#uint32) |  | msg_index in a batch transaction |
| `eth_tx_index` | [int32](#int32) |  | eth_tx_index is the index in the list of valid eth tx in the block, aka. the transaction list returned by eth_getBlock api. |
| `failed` | [bool](#bool) |  | failed is true if the eth transaction did not go succeed |
| `gas_used` | [uint64](#uint64) |  | gas_used by the transaction. If it exceeds the block gas limit, it's set to gas limit, which is what's actually deducted by ante handler. |
| `cumulative_gas_used` | [uint64](#uint64) |  | cumulative_gas_used specifies the cumulated amount of gas used for all processed messages within the current batch transaction. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethermint/types/v1/web3.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethermint/types/v1/web3.proto



<a name="ethermint.types.v1.ExtensionOptionsWeb3Tx"></a>

### ExtensionOptionsWeb3Tx
ExtensionOptionsWeb3Tx is an extension option that specifies the typed chain id,
the fee payer as well as its signature data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `typed_data_chain_id` | [uint64](#uint64) |  | typed_data_chain_id is used only in EIP712 Domain and should match Ethereum network ID in a Web3 provider (e.g. Metamask). |
| `fee_payer` | [string](#string) |  | fee_payer is an account address for the fee payer. It will be validated during EIP712 signature checking. |
| `fee_payer_sig` | [bytes](#bytes) |  | fee_payer_sig is a signature data from the fee paying account, allows to perform fee delegation when using EIP712 Domain. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

