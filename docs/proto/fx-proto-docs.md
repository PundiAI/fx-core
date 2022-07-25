<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [fx/crosschain/v1/crosschain.proto](#fx/crosschain/v1/crosschain.proto)
    - [Attestation](#fx.gravity.crosschain.v1.Attestation)
    - [BatchFees](#fx.gravity.crosschain.v1.BatchFees)
    - [BridgeToken](#fx.gravity.crosschain.v1.BridgeToken)
    - [BridgeValidator](#fx.gravity.crosschain.v1.BridgeValidator)
    - [ChainOracle](#fx.gravity.crosschain.v1.ChainOracle)
    - [ExternalToken](#fx.gravity.crosschain.v1.ExternalToken)
    - [IDSet](#fx.gravity.crosschain.v1.IDSet)
    - [InitCrossChainParamsProposal](#fx.gravity.crosschain.v1.InitCrossChainParamsProposal)
    - [LastObservedBlockHeight](#fx.gravity.crosschain.v1.LastObservedBlockHeight)
    - [Oracle](#fx.gravity.crosschain.v1.Oracle)
    - [OracleSet](#fx.gravity.crosschain.v1.OracleSet)
    - [OutgoingTransferTx](#fx.gravity.crosschain.v1.OutgoingTransferTx)
    - [OutgoingTxBatch](#fx.gravity.crosschain.v1.OutgoingTxBatch)
    - [Params](#fx.gravity.crosschain.v1.Params)
    - [UpdateChainOraclesProposal](#fx.gravity.crosschain.v1.UpdateChainOraclesProposal)
  
    - [ClaimType](#fx.gravity.crosschain.v1.ClaimType)
    - [SignType](#fx.gravity.crosschain.v1.SignType)
  
- [fx/crosschain/v1/tx.proto](#fx/crosschain/v1/tx.proto)
    - [MsgAddOracleDeposit](#fx.gravity.crosschain.v1.MsgAddOracleDeposit)
    - [MsgAddOracleDepositResponse](#fx.gravity.crosschain.v1.MsgAddOracleDepositResponse)
    - [MsgBridgeTokenClaim](#fx.gravity.crosschain.v1.MsgBridgeTokenClaim)
    - [MsgBridgeTokenClaimResponse](#fx.gravity.crosschain.v1.MsgBridgeTokenClaimResponse)
    - [MsgCancelSendToExternal](#fx.gravity.crosschain.v1.MsgCancelSendToExternal)
    - [MsgCancelSendToExternalResponse](#fx.gravity.crosschain.v1.MsgCancelSendToExternalResponse)
    - [MsgConfirmBatch](#fx.gravity.crosschain.v1.MsgConfirmBatch)
    - [MsgConfirmBatchResponse](#fx.gravity.crosschain.v1.MsgConfirmBatchResponse)
    - [MsgOracleSetConfirm](#fx.gravity.crosschain.v1.MsgOracleSetConfirm)
    - [MsgOracleSetConfirmResponse](#fx.gravity.crosschain.v1.MsgOracleSetConfirmResponse)
    - [MsgOracleSetUpdatedClaim](#fx.gravity.crosschain.v1.MsgOracleSetUpdatedClaim)
    - [MsgOracleSetUpdatedClaimResponse](#fx.gravity.crosschain.v1.MsgOracleSetUpdatedClaimResponse)
    - [MsgRequestBatch](#fx.gravity.crosschain.v1.MsgRequestBatch)
    - [MsgRequestBatchResponse](#fx.gravity.crosschain.v1.MsgRequestBatchResponse)
    - [MsgSendToExternal](#fx.gravity.crosschain.v1.MsgSendToExternal)
    - [MsgSendToExternalClaim](#fx.gravity.crosschain.v1.MsgSendToExternalClaim)
    - [MsgSendToExternalClaimResponse](#fx.gravity.crosschain.v1.MsgSendToExternalClaimResponse)
    - [MsgSendToExternalResponse](#fx.gravity.crosschain.v1.MsgSendToExternalResponse)
    - [MsgSendToFxClaim](#fx.gravity.crosschain.v1.MsgSendToFxClaim)
    - [MsgSendToFxClaimResponse](#fx.gravity.crosschain.v1.MsgSendToFxClaimResponse)
    - [MsgSetOrchestratorAddress](#fx.gravity.crosschain.v1.MsgSetOrchestratorAddress)
    - [MsgSetOrchestratorAddressResponse](#fx.gravity.crosschain.v1.MsgSetOrchestratorAddressResponse)
  
    - [Msg](#fx.gravity.crosschain.v1.Msg)
  
- [fx/crosschain/v1/genesis.proto](#fx/crosschain/v1/genesis.proto)
    - [GenesisState](#fx.gravity.crosschain.v1.GenesisState)
  
- [fx/crosschain/v1/query.proto](#fx/crosschain/v1/query.proto)
    - [QueryBatchConfirmRequest](#fx.gravity.crosschain.v1.QueryBatchConfirmRequest)
    - [QueryBatchConfirmResponse](#fx.gravity.crosschain.v1.QueryBatchConfirmResponse)
    - [QueryBatchConfirmsRequest](#fx.gravity.crosschain.v1.QueryBatchConfirmsRequest)
    - [QueryBatchConfirmsResponse](#fx.gravity.crosschain.v1.QueryBatchConfirmsResponse)
    - [QueryBatchFeeRequest](#fx.gravity.crosschain.v1.QueryBatchFeeRequest)
    - [QueryBatchFeeResponse](#fx.gravity.crosschain.v1.QueryBatchFeeResponse)
    - [QueryBatchRequestByNonceRequest](#fx.gravity.crosschain.v1.QueryBatchRequestByNonceRequest)
    - [QueryBatchRequestByNonceResponse](#fx.gravity.crosschain.v1.QueryBatchRequestByNonceResponse)
    - [QueryCurrentOracleSetRequest](#fx.gravity.crosschain.v1.QueryCurrentOracleSetRequest)
    - [QueryCurrentOracleSetResponse](#fx.gravity.crosschain.v1.QueryCurrentOracleSetResponse)
    - [QueryDenomToTokenRequest](#fx.gravity.crosschain.v1.QueryDenomToTokenRequest)
    - [QueryDenomToTokenResponse](#fx.gravity.crosschain.v1.QueryDenomToTokenResponse)
    - [QueryIbcSequenceHeightRequest](#fx.gravity.crosschain.v1.QueryIbcSequenceHeightRequest)
    - [QueryIbcSequenceHeightResponse](#fx.gravity.crosschain.v1.QueryIbcSequenceHeightResponse)
    - [QueryLastEventBlockHeightByAddrRequest](#fx.gravity.crosschain.v1.QueryLastEventBlockHeightByAddrRequest)
    - [QueryLastEventBlockHeightByAddrResponse](#fx.gravity.crosschain.v1.QueryLastEventBlockHeightByAddrResponse)
    - [QueryLastEventNonceByAddrRequest](#fx.gravity.crosschain.v1.QueryLastEventNonceByAddrRequest)
    - [QueryLastEventNonceByAddrResponse](#fx.gravity.crosschain.v1.QueryLastEventNonceByAddrResponse)
    - [QueryLastObservedBlockHeightRequest](#fx.gravity.crosschain.v1.QueryLastObservedBlockHeightRequest)
    - [QueryLastObservedBlockHeightResponse](#fx.gravity.crosschain.v1.QueryLastObservedBlockHeightResponse)
    - [QueryLastOracleSetRequestsRequest](#fx.gravity.crosschain.v1.QueryLastOracleSetRequestsRequest)
    - [QueryLastOracleSetRequestsResponse](#fx.gravity.crosschain.v1.QueryLastOracleSetRequestsResponse)
    - [QueryLastPendingBatchRequestByAddrRequest](#fx.gravity.crosschain.v1.QueryLastPendingBatchRequestByAddrRequest)
    - [QueryLastPendingBatchRequestByAddrResponse](#fx.gravity.crosschain.v1.QueryLastPendingBatchRequestByAddrResponse)
    - [QueryLastPendingOracleSetRequestByAddrRequest](#fx.gravity.crosschain.v1.QueryLastPendingOracleSetRequestByAddrRequest)
    - [QueryLastPendingOracleSetRequestByAddrResponse](#fx.gravity.crosschain.v1.QueryLastPendingOracleSetRequestByAddrResponse)
    - [QueryOracleByAddrRequest](#fx.gravity.crosschain.v1.QueryOracleByAddrRequest)
    - [QueryOracleByExternalAddrRequest](#fx.gravity.crosschain.v1.QueryOracleByExternalAddrRequest)
    - [QueryOracleByOrchestratorRequest](#fx.gravity.crosschain.v1.QueryOracleByOrchestratorRequest)
    - [QueryOracleResponse](#fx.gravity.crosschain.v1.QueryOracleResponse)
    - [QueryOracleSetConfirmRequest](#fx.gravity.crosschain.v1.QueryOracleSetConfirmRequest)
    - [QueryOracleSetConfirmResponse](#fx.gravity.crosschain.v1.QueryOracleSetConfirmResponse)
    - [QueryOracleSetConfirmsByNonceRequest](#fx.gravity.crosschain.v1.QueryOracleSetConfirmsByNonceRequest)
    - [QueryOracleSetConfirmsByNonceResponse](#fx.gravity.crosschain.v1.QueryOracleSetConfirmsByNonceResponse)
    - [QueryOracleSetRequestRequest](#fx.gravity.crosschain.v1.QueryOracleSetRequestRequest)
    - [QueryOracleSetRequestResponse](#fx.gravity.crosschain.v1.QueryOracleSetRequestResponse)
    - [QueryOraclesRequest](#fx.gravity.crosschain.v1.QueryOraclesRequest)
    - [QueryOraclesResponse](#fx.gravity.crosschain.v1.QueryOraclesResponse)
    - [QueryOutgoingTxBatchesRequest](#fx.gravity.crosschain.v1.QueryOutgoingTxBatchesRequest)
    - [QueryOutgoingTxBatchesResponse](#fx.gravity.crosschain.v1.QueryOutgoingTxBatchesResponse)
    - [QueryParamsRequest](#fx.gravity.crosschain.v1.QueryParamsRequest)
    - [QueryParamsResponse](#fx.gravity.crosschain.v1.QueryParamsResponse)
    - [QueryPendingSendToExternalRequest](#fx.gravity.crosschain.v1.QueryPendingSendToExternalRequest)
    - [QueryPendingSendToExternalResponse](#fx.gravity.crosschain.v1.QueryPendingSendToExternalResponse)
    - [QueryTokenToDenomRequest](#fx.gravity.crosschain.v1.QueryTokenToDenomRequest)
    - [QueryTokenToDenomResponse](#fx.gravity.crosschain.v1.QueryTokenToDenomResponse)
  
    - [Query](#fx.gravity.crosschain.v1.Query)
  
- [fx/gravity/v1/attestation.proto](#fx/gravity/v1/attestation.proto)
    - [Attestation](#fx.gravity.v1.Attestation)
    - [ERC20Token](#fx.gravity.v1.ERC20Token)
  
    - [ClaimType](#fx.gravity.v1.ClaimType)
  
- [fx/gravity/v1/batch.proto](#fx/gravity/v1/batch.proto)
    - [OutgoingTransferTx](#fx.gravity.v1.OutgoingTransferTx)
    - [OutgoingTxBatch](#fx.gravity.v1.OutgoingTxBatch)
  
- [fx/gravity/v1/ethereum_signer.proto](#fx/gravity/v1/ethereum_signer.proto)
    - [SignType](#fx.gravity.v1.SignType)
  
- [fx/gravity/v1/types.proto](#fx/gravity/v1/types.proto)
    - [BridgeValidator](#fx.gravity.v1.BridgeValidator)
    - [ERC20ToDenom](#fx.gravity.v1.ERC20ToDenom)
    - [LastObservedEthereumBlockHeight](#fx.gravity.v1.LastObservedEthereumBlockHeight)
    - [Valset](#fx.gravity.v1.Valset)
  
- [fx/gravity/v1/msgs.proto](#fx/gravity/v1/msgs.proto)
    - [MsgCancelSendToEth](#fx.gravity.v1.MsgCancelSendToEth)
    - [MsgCancelSendToEthResponse](#fx.gravity.v1.MsgCancelSendToEthResponse)
    - [MsgConfirmBatch](#fx.gravity.v1.MsgConfirmBatch)
    - [MsgConfirmBatchResponse](#fx.gravity.v1.MsgConfirmBatchResponse)
    - [MsgDepositClaim](#fx.gravity.v1.MsgDepositClaim)
    - [MsgDepositClaimResponse](#fx.gravity.v1.MsgDepositClaimResponse)
    - [MsgFxOriginatedTokenClaim](#fx.gravity.v1.MsgFxOriginatedTokenClaim)
    - [MsgFxOriginatedTokenClaimResponse](#fx.gravity.v1.MsgFxOriginatedTokenClaimResponse)
    - [MsgRequestBatch](#fx.gravity.v1.MsgRequestBatch)
    - [MsgRequestBatchResponse](#fx.gravity.v1.MsgRequestBatchResponse)
    - [MsgSendToEth](#fx.gravity.v1.MsgSendToEth)
    - [MsgSendToEthResponse](#fx.gravity.v1.MsgSendToEthResponse)
    - [MsgSetOrchestratorAddress](#fx.gravity.v1.MsgSetOrchestratorAddress)
    - [MsgSetOrchestratorAddressResponse](#fx.gravity.v1.MsgSetOrchestratorAddressResponse)
    - [MsgValsetConfirm](#fx.gravity.v1.MsgValsetConfirm)
    - [MsgValsetConfirmResponse](#fx.gravity.v1.MsgValsetConfirmResponse)
    - [MsgValsetUpdatedClaim](#fx.gravity.v1.MsgValsetUpdatedClaim)
    - [MsgValsetUpdatedClaimResponse](#fx.gravity.v1.MsgValsetUpdatedClaimResponse)
    - [MsgWithdrawClaim](#fx.gravity.v1.MsgWithdrawClaim)
    - [MsgWithdrawClaimResponse](#fx.gravity.v1.MsgWithdrawClaimResponse)
  
    - [Msg](#fx.gravity.v1.Msg)
  
- [fx/gravity/v1/genesis.proto](#fx/gravity/v1/genesis.proto)
    - [GenesisState](#fx.gravity.v1.GenesisState)
    - [Params](#fx.gravity.v1.Params)
  
- [fx/gravity/v1/pool.proto](#fx/gravity/v1/pool.proto)
    - [BatchFees](#fx.gravity.v1.BatchFees)
    - [IDSet](#fx.gravity.v1.IDSet)
  
- [fx/gravity/v1/query.proto](#fx/gravity/v1/query.proto)
    - [QueryBatchConfirmRequest](#fx.gravity.v1.QueryBatchConfirmRequest)
    - [QueryBatchConfirmResponse](#fx.gravity.v1.QueryBatchConfirmResponse)
    - [QueryBatchConfirmsRequest](#fx.gravity.v1.QueryBatchConfirmsRequest)
    - [QueryBatchConfirmsResponse](#fx.gravity.v1.QueryBatchConfirmsResponse)
    - [QueryBatchFeeRequest](#fx.gravity.v1.QueryBatchFeeRequest)
    - [QueryBatchFeeResponse](#fx.gravity.v1.QueryBatchFeeResponse)
    - [QueryBatchRequestByNonceRequest](#fx.gravity.v1.QueryBatchRequestByNonceRequest)
    - [QueryBatchRequestByNonceResponse](#fx.gravity.v1.QueryBatchRequestByNonceResponse)
    - [QueryCurrentValsetRequest](#fx.gravity.v1.QueryCurrentValsetRequest)
    - [QueryCurrentValsetResponse](#fx.gravity.v1.QueryCurrentValsetResponse)
    - [QueryDelegateKeyByEthRequest](#fx.gravity.v1.QueryDelegateKeyByEthRequest)
    - [QueryDelegateKeyByEthResponse](#fx.gravity.v1.QueryDelegateKeyByEthResponse)
    - [QueryDelegateKeyByOrchestratorRequest](#fx.gravity.v1.QueryDelegateKeyByOrchestratorRequest)
    - [QueryDelegateKeyByOrchestratorResponse](#fx.gravity.v1.QueryDelegateKeyByOrchestratorResponse)
    - [QueryDelegateKeyByValidatorRequest](#fx.gravity.v1.QueryDelegateKeyByValidatorRequest)
    - [QueryDelegateKeyByValidatorResponse](#fx.gravity.v1.QueryDelegateKeyByValidatorResponse)
    - [QueryDenomToERC20Request](#fx.gravity.v1.QueryDenomToERC20Request)
    - [QueryDenomToERC20Response](#fx.gravity.v1.QueryDenomToERC20Response)
    - [QueryERC20ToDenomRequest](#fx.gravity.v1.QueryERC20ToDenomRequest)
    - [QueryERC20ToDenomResponse](#fx.gravity.v1.QueryERC20ToDenomResponse)
    - [QueryIbcSequenceHeightRequest](#fx.gravity.v1.QueryIbcSequenceHeightRequest)
    - [QueryIbcSequenceHeightResponse](#fx.gravity.v1.QueryIbcSequenceHeightResponse)
    - [QueryLastEventBlockHeightByAddrRequest](#fx.gravity.v1.QueryLastEventBlockHeightByAddrRequest)
    - [QueryLastEventBlockHeightByAddrResponse](#fx.gravity.v1.QueryLastEventBlockHeightByAddrResponse)
    - [QueryLastEventNonceByAddrRequest](#fx.gravity.v1.QueryLastEventNonceByAddrRequest)
    - [QueryLastEventNonceByAddrResponse](#fx.gravity.v1.QueryLastEventNonceByAddrResponse)
    - [QueryLastObservedEthBlockHeightRequest](#fx.gravity.v1.QueryLastObservedEthBlockHeightRequest)
    - [QueryLastObservedEthBlockHeightResponse](#fx.gravity.v1.QueryLastObservedEthBlockHeightResponse)
    - [QueryLastPendingBatchRequestByAddrRequest](#fx.gravity.v1.QueryLastPendingBatchRequestByAddrRequest)
    - [QueryLastPendingBatchRequestByAddrResponse](#fx.gravity.v1.QueryLastPendingBatchRequestByAddrResponse)
    - [QueryLastPendingValsetRequestByAddrRequest](#fx.gravity.v1.QueryLastPendingValsetRequestByAddrRequest)
    - [QueryLastPendingValsetRequestByAddrResponse](#fx.gravity.v1.QueryLastPendingValsetRequestByAddrResponse)
    - [QueryLastValsetRequestsRequest](#fx.gravity.v1.QueryLastValsetRequestsRequest)
    - [QueryLastValsetRequestsResponse](#fx.gravity.v1.QueryLastValsetRequestsResponse)
    - [QueryOutgoingTxBatchesRequest](#fx.gravity.v1.QueryOutgoingTxBatchesRequest)
    - [QueryOutgoingTxBatchesResponse](#fx.gravity.v1.QueryOutgoingTxBatchesResponse)
    - [QueryParamsRequest](#fx.gravity.v1.QueryParamsRequest)
    - [QueryParamsResponse](#fx.gravity.v1.QueryParamsResponse)
    - [QueryPendingSendToEthRequest](#fx.gravity.v1.QueryPendingSendToEthRequest)
    - [QueryPendingSendToEthResponse](#fx.gravity.v1.QueryPendingSendToEthResponse)
    - [QueryValsetConfirmRequest](#fx.gravity.v1.QueryValsetConfirmRequest)
    - [QueryValsetConfirmResponse](#fx.gravity.v1.QueryValsetConfirmResponse)
    - [QueryValsetConfirmsByNonceRequest](#fx.gravity.v1.QueryValsetConfirmsByNonceRequest)
    - [QueryValsetConfirmsByNonceResponse](#fx.gravity.v1.QueryValsetConfirmsByNonceResponse)
    - [QueryValsetRequestRequest](#fx.gravity.v1.QueryValsetRequestRequest)
    - [QueryValsetRequestResponse](#fx.gravity.v1.QueryValsetRequestResponse)
  
    - [Query](#fx.gravity.v1.Query)
  
- [fx/other/query.proto](#fx/other/query.proto)
    - [GasPriceRequest](#fx.other.GasPriceRequest)
    - [GasPriceResponse](#fx.other.GasPriceResponse)
  
    - [Query](#fx.other.Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="fx/crosschain/v1/crosschain.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/crosschain/v1/crosschain.proto



<a name="fx.gravity.crosschain.v1.Attestation"></a>

### Attestation
Attestation is an aggregate of `claims` that eventually becomes `observed` by
all orchestrators
EVENT_NONCE:
EventNonce a nonce provided by the gravity contract that is unique per event
fired These event nonces must be relayed in order. This is a correctness
issue, if relaying out of order transaction replay attacks become possible
OBSERVED:
Observed indicates that >67% of validators have attested to the event,
and that the event should be executed by the gravity state machine

The actual content of the claims is passed in with the transaction making the
claim and then passed through the call stack alongside the attestation while
it is processed the key in which the attestation is stored is keyed on the
exact details of the claim but there is no reason to store those exact
details becuause the next message sender will kindly provide you with them.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `observed` | [bool](#bool) |  |  |
| `votes` | [string](#string) | repeated |  |
| `height` | [uint64](#uint64) |  |  |
| `claim` | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="fx.gravity.crosschain.v1.BatchFees"></a>

### BatchFees



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_contract` | [string](#string) |  |  |
| `total_fees` | [string](#string) |  |  |
| `total_txs` | [uint64](#uint64) |  |  |






<a name="fx.gravity.crosschain.v1.BridgeToken"></a>

### BridgeToken
BridgeToken


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token` | [string](#string) |  |  |
| `denom` | [string](#string) |  |  |
| `channel_ibc` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.BridgeValidator"></a>

### BridgeValidator
BridgeValidator represents a validator's external address and its power


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `power` | [uint64](#uint64) |  |  |
| `external_address` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.ChainOracle"></a>

### ChainOracle
module oracles


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracles` | [string](#string) | repeated |  |






<a name="fx.gravity.crosschain.v1.ExternalToken"></a>

### ExternalToken
ERC20Token unique identifier for an Ethereum ERC20 token.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.IDSet"></a>

### IDSet
IDSet represents a set of IDs


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `ids` | [uint64](#uint64) | repeated |  |






<a name="fx.gravity.crosschain.v1.InitCrossChainParamsProposal"></a>

### InitCrossChainParamsProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | the title of the update proposal |
| `description` | [string](#string) |  | the description of the proposal |
| `params` | [Params](#fx.gravity.crosschain.v1.Params) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.LastObservedBlockHeight"></a>

### LastObservedBlockHeight
LastObservedBlockHeight stores the last observed
external block height along with the our block height that
it was observed at. These two numbers can be used to project
outward and always produce batches with timeouts in the future
even if no Ethereum block height has been relayed for a long time


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `external_block_height` | [uint64](#uint64) |  |  |
| `block_height` | [uint64](#uint64) |  |  |






<a name="fx.gravity.crosschain.v1.Oracle"></a>

### Oracle



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_address` | [string](#string) |  |  |
| `orchestrator_address` | [string](#string) |  |  |
| `external_address` | [string](#string) |  |  |
| `deposit_amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `start_height` | [int64](#int64) |  | start oracle height |
| `jailed` | [bool](#bool) |  |  |
| `jailed_height` | [int64](#int64) |  |  |






<a name="fx.gravity.crosschain.v1.OracleSet"></a>

### OracleSet
OracleSet is the external Chain Bridge Multsig Set, each gravity validator
also maintains an external key to sign messages, these are used to check
signatures on external because of the significant gas savings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `members` | [BridgeValidator](#fx.gravity.crosschain.v1.BridgeValidator) | repeated |  |
| `height` | [uint64](#uint64) |  |  |






<a name="fx.gravity.crosschain.v1.OutgoingTransferTx"></a>

### OutgoingTransferTx
OutgoingTransferTx represents an individual send from gravity to ETH


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `sender` | [string](#string) |  |  |
| `dest_address` | [string](#string) |  |  |
| `token` | [ExternalToken](#fx.gravity.crosschain.v1.ExternalToken) |  |  |
| `fee` | [ExternalToken](#fx.gravity.crosschain.v1.ExternalToken) |  |  |






<a name="fx.gravity.crosschain.v1.OutgoingTxBatch"></a>

### OutgoingTxBatch
OutgoingTxBatch represents a batch of transactions going from gravity to ETH


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch_nonce` | [uint64](#uint64) |  |  |
| `batch_timeout` | [uint64](#uint64) |  |  |
| `transactions` | [OutgoingTransferTx](#fx.gravity.crosschain.v1.OutgoingTransferTx) | repeated |  |
| `token_contract` | [string](#string) |  |  |
| `block` | [uint64](#uint64) |  |  |
| `feeReceive` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.Params"></a>

### Params
oracle_set_update_power_change_percent

If power change between validators of CurrentOracleSet and latest oracle set
request is > 10%


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gravity_id` | [string](#string) |  |  |
| `average_block_time` | [uint64](#uint64) |  |  |
| `external_batch_timeout` | [uint64](#uint64) |  |  |
| `average_external_block_time` | [uint64](#uint64) |  |  |
| `signed_window` | [uint64](#uint64) |  |  |
| `slash_fraction` | [bytes](#bytes) |  |  |
| `oracle_set_update_power_change_percent` | [bytes](#bytes) |  |  |
| `ibc_transfer_timeout_height` | [uint64](#uint64) |  |  |
| `oracles` | [string](#string) | repeated |  |
| `deposit_threshold` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="fx.gravity.crosschain.v1.UpdateChainOraclesProposal"></a>

### UpdateChainOraclesProposal



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | the title of the update proposal |
| `description` | [string](#string) |  | the description of the proposal |
| `oracles` | [string](#string) | repeated |  |
| `chain_name` | [string](#string) |  |  |





 <!-- end messages -->


<a name="fx.gravity.crosschain.v1.ClaimType"></a>

### ClaimType
ClaimType is the cosmos type of an event from the counterpart chain that can
be handled

| Name | Number | Description |
| ---- | ------ | ----------- |
| CLAIM_TYPE_UNSPECIFIED | 0 |  |
| CLAIM_TYPE_SEND_TO_FX | 1 |  |
| CLAIM_TYPE_SEND_TO_EXTERNAL | 2 |  |
| CLAIM_TYPE_BRIDGE_TOKEN | 3 |  |
| CLAIM_TYPE_ORACLE_SET_UPDATED | 4 |  |



<a name="fx.gravity.crosschain.v1.SignType"></a>

### SignType
SignType defines messages that have been signed by an orchestrator

| Name | Number | Description |
| ---- | ------ | ----------- |
| SIGN_TYPE_UNSPECIFIED | 0 |  |
| SIGN_TYPE_ORCHESTRATOR_SIGNED_MULTI_SIG_UPDATE | 1 |  |
| SIGN_TYPE_ORCHESTRATOR_SIGNED_WITHDRAW_BATCH | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="fx/crosschain/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/crosschain/v1/tx.proto



<a name="fx.gravity.crosschain.v1.MsgAddOracleDeposit"></a>

### MsgAddOracleDeposit



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgAddOracleDepositResponse"></a>

### MsgAddOracleDepositResponse







<a name="fx.gravity.crosschain.v1.MsgBridgeTokenClaim"></a>

### MsgBridgeTokenClaim



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `event_nonce` | [uint64](#uint64) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `symbol` | [string](#string) |  |  |
| `decimals` | [uint64](#uint64) |  |  |
| `orchestrator` | [string](#string) |  |  |
| `channel_ibc` | [string](#string) |  | Bridge Token channel IBC |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgBridgeTokenClaimResponse"></a>

### MsgBridgeTokenClaimResponse







<a name="fx.gravity.crosschain.v1.MsgCancelSendToExternal"></a>

### MsgCancelSendToExternal
This call allows the sender (and only the sender)
to cancel a given MsgSendToExternal and recieve a refund
of the tokens


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `transaction_id` | [uint64](#uint64) |  |  |
| `sender` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgCancelSendToExternalResponse"></a>

### MsgCancelSendToExternalResponse







<a name="fx.gravity.crosschain.v1.MsgConfirmBatch"></a>

### MsgConfirmBatch
MsgConfirmBatch
When validators observe a MsgRequestBatch they form a batch by ordering
transactions currently in the txqueue in order of highest to lowest fee,
cutting off when the batch either reaches a hardcoded maximum size (to be
decided, probably around 100) or when transactions stop being profitable
(determine this without nondeterminism) This message includes the batch
as well as an Bsc signature over this batch by the validator
-------------


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
| `orchestrator_address` | [string](#string) |  |  |
| `external_address` | [string](#string) |  |  |
| `signature` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgConfirmBatchResponse"></a>

### MsgConfirmBatchResponse







<a name="fx.gravity.crosschain.v1.MsgOracleSetConfirm"></a>

### MsgOracleSetConfirm
MsgOracleSetConfirm
this is the message sent by the validators when they wish to submit their
signatures over the validator set at a given block height. A validator must
first call MsgSetEthAddress to set their Ethereum address to be used for
signing. Then someone (anyone) must make a OracleSetRequest, the request is
essentially a messaging mechanism to determine which block all validators
should submit signatures over. Finally validators sign the validator set,
powers, and Ethereum addresses of the entire validator set at the height of a
OracleSetRequest and submit that signature with this message.

If a sufficient number of validators (66% of voting power) (A) have set
Ethereum addresses and (B) submit OracleSetConfirm messages with their
signatures it is then possible for anyone to view these signatures in the
chain store and submit them to Ethereum to update the validator set
-------------


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `orchestrator_address` | [string](#string) |  |  |
| `external_address` | [string](#string) |  |  |
| `signature` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgOracleSetConfirmResponse"></a>

### MsgOracleSetConfirmResponse







<a name="fx.gravity.crosschain.v1.MsgOracleSetUpdatedClaim"></a>

### MsgOracleSetUpdatedClaim
This informs the Cosmos module that a validator
set has been updated.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `event_nonce` | [uint64](#uint64) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `oracle_set_nonce` | [uint64](#uint64) |  |  |
| `members` | [BridgeValidator](#fx.gravity.crosschain.v1.BridgeValidator) | repeated |  |
| `orchestrator` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgOracleSetUpdatedClaimResponse"></a>

### MsgOracleSetUpdatedClaimResponse







<a name="fx.gravity.crosschain.v1.MsgRequestBatch"></a>

### MsgRequestBatch
MsgRequestBatch
this is a message anyone can send that requests a batch of transactions to
send across the bridge be created for whatever block height this message is
included in. This acts as a coordination point, the handler for this message
looks at the AddToOutgoingPool tx's in the store and generates a batch, also
available in the store tied to this message. The validators then grab this
batch, sign it, submit the signatures with a MsgConfirmBatch before a relayer
can finally submit the batch
-------------
feeReceive:


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `denom` | [string](#string) |  |  |
| `minimum_fee` | [string](#string) |  |  |
| `feeReceive` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgRequestBatchResponse"></a>

### MsgRequestBatchResponse







<a name="fx.gravity.crosschain.v1.MsgSendToExternal"></a>

### MsgSendToExternal
MsgSendToExternal
This is the message that a user calls when they want to bridge an asset
it will later be removed when it is included in a batch and successfully
submitted tokens are removed from the users balance immediately
-------------
AMOUNT:
the coin to send across the bridge, note the restriction that this is a
single coin not a set of coins that is normal in other Payment messages
FEE:
the fee paid for the bridge, distinct from the fee paid to the chain to
actually send this message in the first place. So a successful send has
two layers of fees for the user


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `dest` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `bridge_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgSendToExternalClaim"></a>

### MsgSendToExternalClaim
SendToExternalClaim claims that a batch of withdrawal
operations on the bridge contract was executed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `event_nonce` | [uint64](#uint64) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `batch_nonce` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
| `orchestrator` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgSendToExternalClaimResponse"></a>

### MsgSendToExternalClaimResponse







<a name="fx.gravity.crosschain.v1.MsgSendToExternalResponse"></a>

### MsgSendToExternalResponse







<a name="fx.gravity.crosschain.v1.MsgSendToFxClaim"></a>

### MsgSendToFxClaim
MsgSendToFxClaim
When more than 66% of the active validator set has
claimed to have seen the deposit enter the bsc blockchain coins are
issued to the Payment address in question
-------------


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `event_nonce` | [uint64](#uint64) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |
| `sender` | [string](#string) |  |  |
| `receiver` | [string](#string) |  |  |
| `target_ibc` | [string](#string) |  |  |
| `orchestrator` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgSendToFxClaimResponse"></a>

### MsgSendToFxClaimResponse







<a name="fx.gravity.crosschain.v1.MsgSetOrchestratorAddress"></a>

### MsgSetOrchestratorAddress



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle` | [string](#string) |  |  |
| `orchestrator` | [string](#string) |  |  |
| `external_address` | [string](#string) |  |  |
| `deposit` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgSetOrchestratorAddressResponse"></a>

### MsgSetOrchestratorAddressResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.gravity.crosschain.v1.Msg"></a>

### Msg
Msg defines the state transitions possible within gravity

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `SetOrchestratorAddress` | [MsgSetOrchestratorAddress](#fx.gravity.crosschain.v1.MsgSetOrchestratorAddress) | [MsgSetOrchestratorAddressResponse](#fx.gravity.crosschain.v1.MsgSetOrchestratorAddressResponse) |  | |
| `AddOracleDeposit` | [MsgAddOracleDeposit](#fx.gravity.crosschain.v1.MsgAddOracleDeposit) | [MsgAddOracleDepositResponse](#fx.gravity.crosschain.v1.MsgAddOracleDepositResponse) |  | |
| `OracleSetConfirm` | [MsgOracleSetConfirm](#fx.gravity.crosschain.v1.MsgOracleSetConfirm) | [MsgOracleSetConfirmResponse](#fx.gravity.crosschain.v1.MsgOracleSetConfirmResponse) |  | |
| `OracleSetUpdateClaim` | [MsgOracleSetUpdatedClaim](#fx.gravity.crosschain.v1.MsgOracleSetUpdatedClaim) | [MsgOracleSetUpdatedClaimResponse](#fx.gravity.crosschain.v1.MsgOracleSetUpdatedClaimResponse) |  | |
| `BridgeTokenClaim` | [MsgBridgeTokenClaim](#fx.gravity.crosschain.v1.MsgBridgeTokenClaim) | [MsgBridgeTokenClaimResponse](#fx.gravity.crosschain.v1.MsgBridgeTokenClaimResponse) |  | |
| `SendToFxClaim` | [MsgSendToFxClaim](#fx.gravity.crosschain.v1.MsgSendToFxClaim) | [MsgSendToFxClaimResponse](#fx.gravity.crosschain.v1.MsgSendToFxClaimResponse) |  | |
| `SendToExternal` | [MsgSendToExternal](#fx.gravity.crosschain.v1.MsgSendToExternal) | [MsgSendToExternalResponse](#fx.gravity.crosschain.v1.MsgSendToExternalResponse) |  | |
| `CancelSendToExternal` | [MsgCancelSendToExternal](#fx.gravity.crosschain.v1.MsgCancelSendToExternal) | [MsgCancelSendToExternalResponse](#fx.gravity.crosschain.v1.MsgCancelSendToExternalResponse) |  | |
| `SendToExternalClaim` | [MsgSendToExternalClaim](#fx.gravity.crosschain.v1.MsgSendToExternalClaim) | [MsgSendToExternalClaimResponse](#fx.gravity.crosschain.v1.MsgSendToExternalClaimResponse) |  | |
| `RequestBatch` | [MsgRequestBatch](#fx.gravity.crosschain.v1.MsgRequestBatch) | [MsgRequestBatchResponse](#fx.gravity.crosschain.v1.MsgRequestBatchResponse) |  | |
| `ConfirmBatch` | [MsgConfirmBatch](#fx.gravity.crosschain.v1.MsgConfirmBatch) | [MsgConfirmBatchResponse](#fx.gravity.crosschain.v1.MsgConfirmBatchResponse) |  | |

 <!-- end services -->



<a name="fx/crosschain/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/crosschain/v1/genesis.proto



<a name="fx.gravity.crosschain.v1.GenesisState"></a>

### GenesisState
GenesisState struct


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#fx.gravity.crosschain.v1.Params) |  |  |
| `last_observed_event_nonce` | [uint64](#uint64) |  |  |
| `last_observed_block_height` | [LastObservedBlockHeight](#fx.gravity.crosschain.v1.LastObservedBlockHeight) |  |  |
| `oracles` | [Oracle](#fx.gravity.crosschain.v1.Oracle) | repeated |  |
| `oracle_sets` | [OracleSet](#fx.gravity.crosschain.v1.OracleSet) | repeated |  |
| `bridge_tokens` | [BridgeToken](#fx.gravity.crosschain.v1.BridgeToken) | repeated |  |
| `unbatched_transfers` | [OutgoingTransferTx](#fx.gravity.crosschain.v1.OutgoingTransferTx) | repeated |  |
| `batches` | [OutgoingTxBatch](#fx.gravity.crosschain.v1.OutgoingTxBatch) | repeated |  |
| `oracle_set_confirms` | [MsgOracleSetConfirm](#fx.gravity.crosschain.v1.MsgOracleSetConfirm) | repeated |  |
| `batch_confirms` | [MsgConfirmBatch](#fx.gravity.crosschain.v1.MsgConfirmBatch) | repeated |  |
| `attestations` | [Attestation](#fx.gravity.crosschain.v1.Attestation) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="fx/crosschain/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/crosschain/v1/query.proto



<a name="fx.gravity.crosschain.v1.QueryBatchConfirmRequest"></a>

### QueryBatchConfirmRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
| `orchestrator_address` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryBatchConfirmResponse"></a>

### QueryBatchConfirmResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirm` | [MsgConfirmBatch](#fx.gravity.crosschain.v1.MsgConfirmBatch) |  |  |






<a name="fx.gravity.crosschain.v1.QueryBatchConfirmsRequest"></a>

### QueryBatchConfirmsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryBatchConfirmsResponse"></a>

### QueryBatchConfirmsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirms` | [MsgConfirmBatch](#fx.gravity.crosschain.v1.MsgConfirmBatch) | repeated |  |






<a name="fx.gravity.crosschain.v1.QueryBatchFeeRequest"></a>

### QueryBatchFeeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryBatchFeeResponse"></a>

### QueryBatchFeeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch_fees` | [BatchFees](#fx.gravity.crosschain.v1.BatchFees) | repeated |  |






<a name="fx.gravity.crosschain.v1.QueryBatchRequestByNonceRequest"></a>

### QueryBatchRequestByNonceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryBatchRequestByNonceResponse"></a>

### QueryBatchRequestByNonceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch` | [OutgoingTxBatch](#fx.gravity.crosschain.v1.OutgoingTxBatch) |  |  |






<a name="fx.gravity.crosschain.v1.QueryCurrentOracleSetRequest"></a>

### QueryCurrentOracleSetRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryCurrentOracleSetResponse"></a>

### QueryCurrentOracleSetResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_set` | [OracleSet](#fx.gravity.crosschain.v1.OracleSet) |  |  |






<a name="fx.gravity.crosschain.v1.QueryDenomToTokenRequest"></a>

### QueryDenomToTokenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryDenomToTokenResponse"></a>

### QueryDenomToTokenResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token` | [string](#string) |  |  |
| `channel_ibc` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryIbcSequenceHeightRequest"></a>

### QueryIbcSequenceHeightRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source_port` | [string](#string) |  |  |
| `source_channel` | [string](#string) |  |  |
| `sequence` | [uint64](#uint64) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryIbcSequenceHeightResponse"></a>

### QueryIbcSequenceHeightResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `found` | [bool](#bool) |  |  |
| `block_height` | [uint64](#uint64) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastEventBlockHeightByAddrRequest"></a>

### QueryLastEventBlockHeightByAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `orchestrator_address` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastEventBlockHeightByAddrResponse"></a>

### QueryLastEventBlockHeightByAddrResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_height` | [uint64](#uint64) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastEventNonceByAddrRequest"></a>

### QueryLastEventNonceByAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `orchestrator_address` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastEventNonceByAddrResponse"></a>

### QueryLastEventNonceByAddrResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `event_nonce` | [uint64](#uint64) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastObservedBlockHeightRequest"></a>

### QueryLastObservedBlockHeightRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastObservedBlockHeightResponse"></a>

### QueryLastObservedBlockHeightResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `external_block_height` | [uint64](#uint64) |  |  |
| `block_height` | [uint64](#uint64) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastOracleSetRequestsRequest"></a>

### QueryLastOracleSetRequestsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastOracleSetRequestsResponse"></a>

### QueryLastOracleSetRequestsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_sets` | [OracleSet](#fx.gravity.crosschain.v1.OracleSet) | repeated |  |






<a name="fx.gravity.crosschain.v1.QueryLastPendingBatchRequestByAddrRequest"></a>

### QueryLastPendingBatchRequestByAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `orchestrator_address` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastPendingBatchRequestByAddrResponse"></a>

### QueryLastPendingBatchRequestByAddrResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch` | [OutgoingTxBatch](#fx.gravity.crosschain.v1.OutgoingTxBatch) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastPendingOracleSetRequestByAddrRequest"></a>

### QueryLastPendingOracleSetRequestByAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `orchestrator_address` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastPendingOracleSetRequestByAddrResponse"></a>

### QueryLastPendingOracleSetRequestByAddrResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_sets` | [OracleSet](#fx.gravity.crosschain.v1.OracleSet) | repeated |  |






<a name="fx.gravity.crosschain.v1.QueryOracleByAddrRequest"></a>

### QueryOracleByAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_address` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOracleByExternalAddrRequest"></a>

### QueryOracleByExternalAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `external_address` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOracleByOrchestratorRequest"></a>

### QueryOracleByOrchestratorRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `orchestrator_address` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOracleResponse"></a>

### QueryOracleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle` | [Oracle](#fx.gravity.crosschain.v1.Oracle) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOracleSetConfirmRequest"></a>

### QueryOracleSetConfirmRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `orchestrator_address` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOracleSetConfirmResponse"></a>

### QueryOracleSetConfirmResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirm` | [MsgOracleSetConfirm](#fx.gravity.crosschain.v1.MsgOracleSetConfirm) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOracleSetConfirmsByNonceRequest"></a>

### QueryOracleSetConfirmsByNonceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOracleSetConfirmsByNonceResponse"></a>

### QueryOracleSetConfirmsByNonceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirms` | [MsgOracleSetConfirm](#fx.gravity.crosschain.v1.MsgOracleSetConfirm) | repeated |  |






<a name="fx.gravity.crosschain.v1.QueryOracleSetRequestRequest"></a>

### QueryOracleSetRequestRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOracleSetRequestResponse"></a>

### QueryOracleSetRequestResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_set` | [OracleSet](#fx.gravity.crosschain.v1.OracleSet) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOraclesRequest"></a>

### QueryOraclesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOraclesResponse"></a>

### QueryOraclesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracles` | [Oracle](#fx.gravity.crosschain.v1.Oracle) | repeated | oracles contains all the queried oracles. |






<a name="fx.gravity.crosschain.v1.QueryOutgoingTxBatchesRequest"></a>

### QueryOutgoingTxBatchesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOutgoingTxBatchesResponse"></a>

### QueryOutgoingTxBatchesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batches` | [OutgoingTxBatch](#fx.gravity.crosschain.v1.OutgoingTxBatch) | repeated |  |






<a name="fx.gravity.crosschain.v1.QueryParamsRequest"></a>

### QueryParamsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryParamsResponse"></a>

### QueryParamsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#fx.gravity.crosschain.v1.Params) |  |  |






<a name="fx.gravity.crosschain.v1.QueryPendingSendToExternalRequest"></a>

### QueryPendingSendToExternalRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender_address` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryPendingSendToExternalResponse"></a>

### QueryPendingSendToExternalResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `transfers_in_batches` | [OutgoingTransferTx](#fx.gravity.crosschain.v1.OutgoingTransferTx) | repeated |  |
| `unbatched_transfers` | [OutgoingTransferTx](#fx.gravity.crosschain.v1.OutgoingTransferTx) | repeated |  |






<a name="fx.gravity.crosschain.v1.QueryTokenToDenomRequest"></a>

### QueryTokenToDenomRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryTokenToDenomResponse"></a>

### QueryTokenToDenomResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `channel_ibc` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.gravity.crosschain.v1.Query"></a>

### Query
Query defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#fx.gravity.crosschain.v1.QueryParamsRequest) | [QueryParamsResponse](#fx.gravity.crosschain.v1.QueryParamsResponse) | Deployments queries deployments | GET|/crosschain/v1beta/params|
| `CurrentOracleSet` | [QueryCurrentOracleSetRequest](#fx.gravity.crosschain.v1.QueryCurrentOracleSetRequest) | [QueryCurrentOracleSetResponse](#fx.gravity.crosschain.v1.QueryCurrentOracleSetResponse) |  | GET|/crosschain/v1beta/oracle_set/current|
| `OracleSetRequest` | [QueryOracleSetRequestRequest](#fx.gravity.crosschain.v1.QueryOracleSetRequestRequest) | [QueryOracleSetRequestResponse](#fx.gravity.crosschain.v1.QueryOracleSetRequestResponse) |  | GET|/crosschain/v1beta/oracle_set/request|
| `OracleSetConfirm` | [QueryOracleSetConfirmRequest](#fx.gravity.crosschain.v1.QueryOracleSetConfirmRequest) | [QueryOracleSetConfirmResponse](#fx.gravity.crosschain.v1.QueryOracleSetConfirmResponse) |  | GET|/crosschain/v1beta/oracle_set/confirm|
| `OracleSetConfirmsByNonce` | [QueryOracleSetConfirmsByNonceRequest](#fx.gravity.crosschain.v1.QueryOracleSetConfirmsByNonceRequest) | [QueryOracleSetConfirmsByNonceResponse](#fx.gravity.crosschain.v1.QueryOracleSetConfirmsByNonceResponse) |  | GET|/crosschain/v1beta/oracle_set/confirms|
| `LastOracleSetRequests` | [QueryLastOracleSetRequestsRequest](#fx.gravity.crosschain.v1.QueryLastOracleSetRequestsRequest) | [QueryLastOracleSetRequestsResponse](#fx.gravity.crosschain.v1.QueryLastOracleSetRequestsResponse) |  | GET|/crosschain/v1beta/oracle_set/requests|
| `LastPendingOracleSetRequestByAddr` | [QueryLastPendingOracleSetRequestByAddrRequest](#fx.gravity.crosschain.v1.QueryLastPendingOracleSetRequestByAddrRequest) | [QueryLastPendingOracleSetRequestByAddrResponse](#fx.gravity.crosschain.v1.QueryLastPendingOracleSetRequestByAddrResponse) |  | GET|/crosschain/v1beta/oracle_set/last|
| `LastPendingBatchRequestByAddr` | [QueryLastPendingBatchRequestByAddrRequest](#fx.gravity.crosschain.v1.QueryLastPendingBatchRequestByAddrRequest) | [QueryLastPendingBatchRequestByAddrResponse](#fx.gravity.crosschain.v1.QueryLastPendingBatchRequestByAddrResponse) |  | GET|/crosschain/v1beta/batch/last|
| `LastEventNonceByAddr` | [QueryLastEventNonceByAddrRequest](#fx.gravity.crosschain.v1.QueryLastEventNonceByAddrRequest) | [QueryLastEventNonceByAddrResponse](#fx.gravity.crosschain.v1.QueryLastEventNonceByAddrResponse) |  | GET|/crosschain/v1beta/oracle/event_nonce|
| `LastEventBlockHeightByAddr` | [QueryLastEventBlockHeightByAddrRequest](#fx.gravity.crosschain.v1.QueryLastEventBlockHeightByAddrRequest) | [QueryLastEventBlockHeightByAddrResponse](#fx.gravity.crosschain.v1.QueryLastEventBlockHeightByAddrResponse) |  | GET|/crosschain/v1beta/oracle/event/block_height|
| `BatchFees` | [QueryBatchFeeRequest](#fx.gravity.crosschain.v1.QueryBatchFeeRequest) | [QueryBatchFeeResponse](#fx.gravity.crosschain.v1.QueryBatchFeeResponse) |  | GET|/crosschain/v1beta/batch_fees|
| `LastObservedBlockHeight` | [QueryLastObservedBlockHeightRequest](#fx.gravity.crosschain.v1.QueryLastObservedBlockHeightRequest) | [QueryLastObservedBlockHeightResponse](#fx.gravity.crosschain.v1.QueryLastObservedBlockHeightResponse) |  | GET|/crosschain/v1beta/block_height|
| `OutgoingTxBatches` | [QueryOutgoingTxBatchesRequest](#fx.gravity.crosschain.v1.QueryOutgoingTxBatchesRequest) | [QueryOutgoingTxBatchesResponse](#fx.gravity.crosschain.v1.QueryOutgoingTxBatchesResponse) |  | GET|/crosschain/v1beta/batch/outgoing_tx|
| `BatchRequestByNonce` | [QueryBatchRequestByNonceRequest](#fx.gravity.crosschain.v1.QueryBatchRequestByNonceRequest) | [QueryBatchRequestByNonceResponse](#fx.gravity.crosschain.v1.QueryBatchRequestByNonceResponse) |  | GET|/crosschain/v1beta/batch/request|
| `BatchConfirm` | [QueryBatchConfirmRequest](#fx.gravity.crosschain.v1.QueryBatchConfirmRequest) | [QueryBatchConfirmResponse](#fx.gravity.crosschain.v1.QueryBatchConfirmResponse) |  | GET|/crosschain/v1beta/batch/confirm|
| `BatchConfirms` | [QueryBatchConfirmsRequest](#fx.gravity.crosschain.v1.QueryBatchConfirmsRequest) | [QueryBatchConfirmsResponse](#fx.gravity.crosschain.v1.QueryBatchConfirmsResponse) |  | GET|/crosschain/v1beta/batch/confirms|
| `TokenToDenom` | [QueryTokenToDenomRequest](#fx.gravity.crosschain.v1.QueryTokenToDenomRequest) | [QueryTokenToDenomResponse](#fx.gravity.crosschain.v1.QueryTokenToDenomResponse) |  | GET|/crosschain/v1beta/bridge_denom|
| `DenomToToken` | [QueryDenomToTokenRequest](#fx.gravity.crosschain.v1.QueryDenomToTokenRequest) | [QueryDenomToTokenResponse](#fx.gravity.crosschain.v1.QueryDenomToTokenResponse) |  | GET|/crosschain/v1beta/bridge_token|
| `GetOracleByAddr` | [QueryOracleByAddrRequest](#fx.gravity.crosschain.v1.QueryOracleByAddrRequest) | [QueryOracleResponse](#fx.gravity.crosschain.v1.QueryOracleResponse) |  | GET|/crosschain/v1beta/query_oracle_by_addr|
| `GetOracleByExternalAddr` | [QueryOracleByExternalAddrRequest](#fx.gravity.crosschain.v1.QueryOracleByExternalAddrRequest) | [QueryOracleResponse](#fx.gravity.crosschain.v1.QueryOracleResponse) |  | GET|/crosschain/v1beta/query_oracle_by_external_addr|
| `GetOracleByOrchestrator` | [QueryOracleByOrchestratorRequest](#fx.gravity.crosschain.v1.QueryOracleByOrchestratorRequest) | [QueryOracleResponse](#fx.gravity.crosschain.v1.QueryOracleResponse) |  | GET|/crosschain/v1beta/query_oracle_by_orchestrator|
| `GetPendingSendToExternal` | [QueryPendingSendToExternalRequest](#fx.gravity.crosschain.v1.QueryPendingSendToExternalRequest) | [QueryPendingSendToExternalResponse](#fx.gravity.crosschain.v1.QueryPendingSendToExternalResponse) |  | GET|/crosschain/v1beta/query_pending_send_to_external|
| `GetIbcSequenceHeightByChannel` | [QueryIbcSequenceHeightRequest](#fx.gravity.crosschain.v1.QueryIbcSequenceHeightRequest) | [QueryIbcSequenceHeightResponse](#fx.gravity.crosschain.v1.QueryIbcSequenceHeightResponse) |  | GET|/crosschain/v1beta/query_ibc_sequence_height|
| `Oracles` | [QueryOraclesRequest](#fx.gravity.crosschain.v1.QueryOraclesRequest) | [QueryOraclesResponse](#fx.gravity.crosschain.v1.QueryOraclesResponse) | Validators queries all oracle that match the given status. | GET|/crosschain/v1beta1/oracles|

 <!-- end services -->



<a name="fx/gravity/v1/attestation.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/gravity/v1/attestation.proto



<a name="fx.gravity.v1.Attestation"></a>

### Attestation
Attestation is an aggregate of `claims` that eventually becomes `observed` by
all orchestrators
EVENT_NONCE:
EventNonce a nonce provided by the gravity contract that is unique per event
fired These event nonces must be relayed in order. This is a correctness
issue, if relaying out of order transaction replay attacks become possible
OBSERVED:
Observed indicates that >67% of validators have attested to the event,
and that the event should be executed by the gravity state machine

The actual content of the claims is passed in with the transaction making the
claim and then passed through the call stack alongside the attestation while
it is processed the key in which the attestation is stored is keyed on the
exact details of the claim but there is no reason to store those exact
details becuause the next message sender will kindly provide you with them.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `observed` | [bool](#bool) |  |  |
| `votes` | [string](#string) | repeated |  |
| `height` | [uint64](#uint64) |  |  |
| `claim` | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="fx.gravity.v1.ERC20Token"></a>

### ERC20Token
ERC20Token unique identifier for an Ethereum ERC20 token.
CONTRACT:
The contract address on ETH of the token, this could be a Cosmos
originated token, if so it will be the ERC20 address of the representation
(note: developers should look up the token symbol using the address on ETH to
display for UI)


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |





 <!-- end messages -->


<a name="fx.gravity.v1.ClaimType"></a>

### ClaimType
ClaimType is the cosmos type of an event from the counterpart chain that can
be handled

| Name | Number | Description |
| ---- | ------ | ----------- |
| CLAIM_TYPE_UNSPECIFIED | 0 |  |
| CLAIM_TYPE_DEPOSIT | 1 |  |
| CLAIM_TYPE_WITHDRAW | 2 |  |
| CLAIM_TYPE_ORIGINATED_TOKEN | 3 |  |
| CLAIM_TYPE_VALSET_UPDATED | 4 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="fx/gravity/v1/batch.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/gravity/v1/batch.proto



<a name="fx.gravity.v1.OutgoingTransferTx"></a>

### OutgoingTransferTx
OutgoingTransferTx represents an individual send from gravity to ETH


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `sender` | [string](#string) |  |  |
| `dest_address` | [string](#string) |  |  |
| `erc20_token` | [ERC20Token](#fx.gravity.v1.ERC20Token) |  |  |
| `erc20_fee` | [ERC20Token](#fx.gravity.v1.ERC20Token) |  |  |






<a name="fx.gravity.v1.OutgoingTxBatch"></a>

### OutgoingTxBatch
OutgoingTxBatch represents a batch of transactions going from gravity to ETH


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch_nonce` | [uint64](#uint64) |  |  |
| `batch_timeout` | [uint64](#uint64) |  |  |
| `transactions` | [OutgoingTransferTx](#fx.gravity.v1.OutgoingTransferTx) | repeated |  |
| `token_contract` | [string](#string) |  |  |
| `block` | [uint64](#uint64) |  |  |
| `feeReceive` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="fx/gravity/v1/ethereum_signer.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/gravity/v1/ethereum_signer.proto


 <!-- end messages -->


<a name="fx.gravity.v1.SignType"></a>

### SignType
SignType defines messages that have been signed by an orchestrator

| Name | Number | Description |
| ---- | ------ | ----------- |
| SIGN_TYPE_UNSPECIFIED | 0 |  |
| SIGN_TYPE_ORCHESTRATOR_SIGNED_MULTI_SIG_UPDATE | 1 |  |
| SIGN_TYPE_ORCHESTRATOR_SIGNED_WITHDRAW_BATCH | 2 |  |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="fx/gravity/v1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/gravity/v1/types.proto



<a name="fx.gravity.v1.BridgeValidator"></a>

### BridgeValidator
BridgeValidator represents a validator's ETH address and its power


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `power` | [uint64](#uint64) |  |  |
| `eth_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.ERC20ToDenom"></a>

### ERC20ToDenom
This records the relationship between an ERC20 token and the denom
of the corresponding fx originated asset


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `erc20` | [string](#string) |  |  |
| `denom` | [string](#string) |  |  |






<a name="fx.gravity.v1.LastObservedEthereumBlockHeight"></a>

### LastObservedEthereumBlockHeight
LastObservedEthereumBlockHeight stores the last observed
Ethereum block height along with the fx block height that
it was observed at. These two numbers can be used to project
outward and always produce batches with timeouts in the future
even if no Ethereum block height has been relayed for a long time


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fx_block_height` | [uint64](#uint64) |  |  |
| `eth_block_height` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.Valset"></a>

### Valset
Valset is the Ethereum Bridge Multsig Set, each gravity validator also
maintains an ETH key to sign messages, these are used to check signatures on
ETH because of the significant gas savings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `members` | [BridgeValidator](#fx.gravity.v1.BridgeValidator) | repeated |  |
| `height` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="fx/gravity/v1/msgs.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/gravity/v1/msgs.proto



<a name="fx.gravity.v1.MsgCancelSendToEth"></a>

### MsgCancelSendToEth
This call allows the sender (and only the sender)
to cancel a given MsgSendToEth and recieve a refund
of the tokens


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `transaction_id` | [uint64](#uint64) |  |  |
| `sender` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgCancelSendToEthResponse"></a>

### MsgCancelSendToEthResponse







<a name="fx.gravity.v1.MsgConfirmBatch"></a>

### MsgConfirmBatch
MsgConfirmBatch
When validators observe a MsgRequestBatch they form a batch by ordering
transactions currently in the txqueue in order of highest to lowest fee,
cutting off when the batch either reaches a hardcoded maximum size (to be
decided, probably around 100) or when transactions stop being profitable
(determine this without nondeterminism) This message includes the batch
as well as an Ethereum signature over this batch by the validator
-------------


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
| `eth_signer` | [string](#string) |  |  |
| `orchestrator` | [string](#string) |  |  |
| `signature` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgConfirmBatchResponse"></a>

### MsgConfirmBatchResponse







<a name="fx.gravity.v1.MsgDepositClaim"></a>

### MsgDepositClaim
EthereumBridgeDepositClaim
When more than 66% of the active validator set has
claimed to have seen the deposit enter the ethereum blockchain coins are
issued to the Cosmos address in question
-------------


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `event_nonce` | [uint64](#uint64) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |
| `eth_sender` | [string](#string) |  |  |
| `fx_receiver` | [string](#string) |  |  |
| `target_ibc` | [string](#string) |  |  |
| `orchestrator` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgDepositClaimResponse"></a>

### MsgDepositClaimResponse







<a name="fx.gravity.v1.MsgFxOriginatedTokenClaim"></a>

### MsgFxOriginatedTokenClaim



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `event_nonce` | [uint64](#uint64) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `symbol` | [string](#string) |  |  |
| `decimals` | [uint64](#uint64) |  |  |
| `orchestrator` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgFxOriginatedTokenClaimResponse"></a>

### MsgFxOriginatedTokenClaimResponse







<a name="fx.gravity.v1.MsgRequestBatch"></a>

### MsgRequestBatch



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `denom` | [string](#string) |  |  |
| `minimum_fee` | [string](#string) |  |  |
| `feeReceive` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgRequestBatchResponse"></a>

### MsgRequestBatchResponse







<a name="fx.gravity.v1.MsgSendToEth"></a>

### MsgSendToEth
MsgSendToEth
This is the message that a user calls when they want to bridge an asset
it will later be removed when it is included in a batch and successfully
submitted tokens are removed from the users balance immediately
-------------
AMOUNT:
the coin to send across the bridge, note the restriction that this is a
single coin not a set of coins that is normal in other Cosmos messages
FEE:
the fee paid for the bridge, distinct from the fee paid to the chain to
actually send this message in the first place. So a successful send has
two layers of fees for the user


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `eth_dest` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `bridge_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="fx.gravity.v1.MsgSendToEthResponse"></a>

### MsgSendToEthResponse







<a name="fx.gravity.v1.MsgSetOrchestratorAddress"></a>

### MsgSetOrchestratorAddress
MsgSetOrchestratorAddress
this message allows validators to delegate their voting responsibilities
to a given key. This key is then used as an optional authentication method
for sigining oracle claims
VALIDATOR
The validator field is a cosmosvaloper1... string (i.e. sdk.ValAddress)
that references a validator in the active set
ORCHESTRATOR
The orchestrator field is a cosmos1... string  (i.e. sdk.AccAddress) that
references the key that is being delegated to
ETH_ADDRESS
This is a hex encoded 0x Ethereum public key that will be used by this
validator on Ethereum


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |
| `orchestrator` | [string](#string) |  |  |
| `eth_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgSetOrchestratorAddressResponse"></a>

### MsgSetOrchestratorAddressResponse







<a name="fx.gravity.v1.MsgValsetConfirm"></a>

### MsgValsetConfirm
MsgValsetConfirm
this is the message sent by the validators when they wish to submit their
signatures over the validator set at a given block height. A validator must
first call MsgSetEthAddress to set their Ethereum address to be used for
signing. Then someone (anyone) must make a ValsetRequest, the request is
essentially a messaging mechanism to determine which block all validators
should submit signatures over. Finally validators sign the validator set,
powers, and Ethereum addresses of the entire validator set at the height of a
ValsetRequest and submit that signature with this message.

If a sufficient number of validators (66% of voting power) (A) have set
Ethereum addresses and (B) submit ValsetConfirm messages with their
signatures it is then possible for anyone to view these signatures in the
chain store and submit them to Ethereum to update the validator set
-------------


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `orchestrator` | [string](#string) |  |  |
| `eth_address` | [string](#string) |  |  |
| `signature` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgValsetConfirmResponse"></a>

### MsgValsetConfirmResponse







<a name="fx.gravity.v1.MsgValsetUpdatedClaim"></a>

### MsgValsetUpdatedClaim
This informs the Cosmos module that a validator
set has been updated.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `event_nonce` | [uint64](#uint64) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `valset_nonce` | [uint64](#uint64) |  |  |
| `members` | [BridgeValidator](#fx.gravity.v1.BridgeValidator) | repeated |  |
| `orchestrator` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgValsetUpdatedClaimResponse"></a>

### MsgValsetUpdatedClaimResponse







<a name="fx.gravity.v1.MsgWithdrawClaim"></a>

### MsgWithdrawClaim
WithdrawClaim claims that a batch of withdrawal
operations on the bridge contract was executed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `event_nonce` | [uint64](#uint64) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `batch_nonce` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
| `orchestrator` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgWithdrawClaimResponse"></a>

### MsgWithdrawClaimResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.gravity.v1.Msg"></a>

### Msg
Msg defines the state transitions possible within gravity

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ValsetConfirm` | [MsgValsetConfirm](#fx.gravity.v1.MsgValsetConfirm) | [MsgValsetConfirmResponse](#fx.gravity.v1.MsgValsetConfirmResponse) |  | |
| `SendToEth` | [MsgSendToEth](#fx.gravity.v1.MsgSendToEth) | [MsgSendToEthResponse](#fx.gravity.v1.MsgSendToEthResponse) |  | |
| `RequestBatch` | [MsgRequestBatch](#fx.gravity.v1.MsgRequestBatch) | [MsgRequestBatchResponse](#fx.gravity.v1.MsgRequestBatchResponse) |  | |
| `ConfirmBatch` | [MsgConfirmBatch](#fx.gravity.v1.MsgConfirmBatch) | [MsgConfirmBatchResponse](#fx.gravity.v1.MsgConfirmBatchResponse) |  | |
| `DepositClaim` | [MsgDepositClaim](#fx.gravity.v1.MsgDepositClaim) | [MsgDepositClaimResponse](#fx.gravity.v1.MsgDepositClaimResponse) |  | |
| `WithdrawClaim` | [MsgWithdrawClaim](#fx.gravity.v1.MsgWithdrawClaim) | [MsgWithdrawClaimResponse](#fx.gravity.v1.MsgWithdrawClaimResponse) |  | |
| `ValsetUpdateClaim` | [MsgValsetUpdatedClaim](#fx.gravity.v1.MsgValsetUpdatedClaim) | [MsgValsetUpdatedClaimResponse](#fx.gravity.v1.MsgValsetUpdatedClaimResponse) |  | |
| `SetOrchestratorAddress` | [MsgSetOrchestratorAddress](#fx.gravity.v1.MsgSetOrchestratorAddress) | [MsgSetOrchestratorAddressResponse](#fx.gravity.v1.MsgSetOrchestratorAddressResponse) |  | |
| `CancelSendToEth` | [MsgCancelSendToEth](#fx.gravity.v1.MsgCancelSendToEth) | [MsgCancelSendToEthResponse](#fx.gravity.v1.MsgCancelSendToEthResponse) |  | |
| `FxOriginatedTokenClaim` | [MsgFxOriginatedTokenClaim](#fx.gravity.v1.MsgFxOriginatedTokenClaim) | [MsgFxOriginatedTokenClaimResponse](#fx.gravity.v1.MsgFxOriginatedTokenClaimResponse) |  | |

 <!-- end services -->



<a name="fx/gravity/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/gravity/v1/genesis.proto



<a name="fx.gravity.v1.GenesisState"></a>

### GenesisState
GenesisState struct


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#fx.gravity.v1.Params) |  |  |
| `last_observed_nonce` | [uint64](#uint64) |  |  |
| `last_observed_block_height` | [LastObservedEthereumBlockHeight](#fx.gravity.v1.LastObservedEthereumBlockHeight) |  |  |
| `delegate_keys` | [MsgSetOrchestratorAddress](#fx.gravity.v1.MsgSetOrchestratorAddress) | repeated |  |
| `valsets` | [Valset](#fx.gravity.v1.Valset) | repeated |  |
| `erc20_to_denoms` | [ERC20ToDenom](#fx.gravity.v1.ERC20ToDenom) | repeated |  |
| `unbatched_transfers` | [OutgoingTransferTx](#fx.gravity.v1.OutgoingTransferTx) | repeated |  |
| `batches` | [OutgoingTxBatch](#fx.gravity.v1.OutgoingTxBatch) | repeated |  |
| `batch_confirms` | [MsgConfirmBatch](#fx.gravity.v1.MsgConfirmBatch) | repeated |  |
| `valset_confirms` | [MsgValsetConfirm](#fx.gravity.v1.MsgValsetConfirm) | repeated |  |
| `attestations` | [Attestation](#fx.gravity.v1.Attestation) | repeated |  |






<a name="fx.gravity.v1.Params"></a>

### Params
valset_update_power_change_percent

If power change between validators of CurrentValset and latest valset request
is > 10%


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gravity_id` | [string](#string) |  |  |
| `contract_source_hash` | [string](#string) |  |  |
| `bridge_eth_address` | [string](#string) |  |  |
| `bridge_chain_id` | [uint64](#uint64) |  |  |
| `signed_valsets_window` | [uint64](#uint64) |  |  |
| `signed_batches_window` | [uint64](#uint64) |  |  |
| `signed_claims_window` | [uint64](#uint64) |  |  |
| `target_batch_timeout` | [uint64](#uint64) |  |  |
| `average_block_time` | [uint64](#uint64) |  |  |
| `average_eth_block_time` | [uint64](#uint64) |  |  |
| `slash_fraction_valset` | [bytes](#bytes) |  |  |
| `slash_fraction_batch` | [bytes](#bytes) |  |  |
| `slash_fraction_claim` | [bytes](#bytes) |  |  |
| `slash_fraction_conflicting_claim` | [bytes](#bytes) |  |  |
| `unbond_slashing_valsets_window` | [uint64](#uint64) |  |  |
| `ibc_transfer_timeout_height` | [uint64](#uint64) |  |  |
| `valset_update_power_change_percent` | [bytes](#bytes) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="fx/gravity/v1/pool.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/gravity/v1/pool.proto



<a name="fx.gravity.v1.BatchFees"></a>

### BatchFees



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_contract` | [string](#string) |  |  |
| `total_fees` | [string](#string) |  |  |
| `total_txs` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.IDSet"></a>

### IDSet
IDSet represents a set of IDs


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `ids` | [uint64](#uint64) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="fx/gravity/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/gravity/v1/query.proto



<a name="fx.gravity.v1.QueryBatchConfirmRequest"></a>

### QueryBatchConfirmRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `contract_address` | [string](#string) |  |  |
| `address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryBatchConfirmResponse"></a>

### QueryBatchConfirmResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirm` | [MsgConfirmBatch](#fx.gravity.v1.MsgConfirmBatch) |  |  |






<a name="fx.gravity.v1.QueryBatchConfirmsRequest"></a>

### QueryBatchConfirmsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `contract_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryBatchConfirmsResponse"></a>

### QueryBatchConfirmsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirms` | [MsgConfirmBatch](#fx.gravity.v1.MsgConfirmBatch) | repeated |  |






<a name="fx.gravity.v1.QueryBatchFeeRequest"></a>

### QueryBatchFeeRequest







<a name="fx.gravity.v1.QueryBatchFeeResponse"></a>

### QueryBatchFeeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch_fees` | [BatchFees](#fx.gravity.v1.BatchFees) | repeated |  |






<a name="fx.gravity.v1.QueryBatchRequestByNonceRequest"></a>

### QueryBatchRequestByNonceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `contract_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryBatchRequestByNonceResponse"></a>

### QueryBatchRequestByNonceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch` | [OutgoingTxBatch](#fx.gravity.v1.OutgoingTxBatch) |  |  |






<a name="fx.gravity.v1.QueryCurrentValsetRequest"></a>

### QueryCurrentValsetRequest







<a name="fx.gravity.v1.QueryCurrentValsetResponse"></a>

### QueryCurrentValsetResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `valset` | [Valset](#fx.gravity.v1.Valset) |  |  |






<a name="fx.gravity.v1.QueryDelegateKeyByEthRequest"></a>

### QueryDelegateKeyByEthRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `eth_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryDelegateKeyByEthResponse"></a>

### QueryDelegateKeyByEthResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |
| `orchestrator_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryDelegateKeyByOrchestratorRequest"></a>

### QueryDelegateKeyByOrchestratorRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `orchestrator_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryDelegateKeyByOrchestratorResponse"></a>

### QueryDelegateKeyByOrchestratorResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |
| `eth_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryDelegateKeyByValidatorRequest"></a>

### QueryDelegateKeyByValidatorRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryDelegateKeyByValidatorResponse"></a>

### QueryDelegateKeyByValidatorResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `eth_address` | [string](#string) |  |  |
| `orchestrator_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryDenomToERC20Request"></a>

### QueryDenomToERC20Request



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryDenomToERC20Response"></a>

### QueryDenomToERC20Response



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `erc20` | [string](#string) |  |  |
| `fx_originated` | [bool](#bool) |  |  |






<a name="fx.gravity.v1.QueryERC20ToDenomRequest"></a>

### QueryERC20ToDenomRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `erc20` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryERC20ToDenomResponse"></a>

### QueryERC20ToDenomResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `fx_originated` | [bool](#bool) |  |  |






<a name="fx.gravity.v1.QueryIbcSequenceHeightRequest"></a>

### QueryIbcSequenceHeightRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sourcePort` | [string](#string) |  |  |
| `sourceChannel` | [string](#string) |  |  |
| `sequence` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.QueryIbcSequenceHeightResponse"></a>

### QueryIbcSequenceHeightResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `found` | [bool](#bool) |  |  |
| `height` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.QueryLastEventBlockHeightByAddrRequest"></a>

### QueryLastEventBlockHeightByAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryLastEventBlockHeightByAddrResponse"></a>

### QueryLastEventBlockHeightByAddrResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_height` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.QueryLastEventNonceByAddrRequest"></a>

### QueryLastEventNonceByAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryLastEventNonceByAddrResponse"></a>

### QueryLastEventNonceByAddrResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `event_nonce` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.QueryLastObservedEthBlockHeightRequest"></a>

### QueryLastObservedEthBlockHeightRequest







<a name="fx.gravity.v1.QueryLastObservedEthBlockHeightResponse"></a>

### QueryLastObservedEthBlockHeightResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `blockHeight` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.QueryLastPendingBatchRequestByAddrRequest"></a>

### QueryLastPendingBatchRequestByAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryLastPendingBatchRequestByAddrResponse"></a>

### QueryLastPendingBatchRequestByAddrResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch` | [OutgoingTxBatch](#fx.gravity.v1.OutgoingTxBatch) |  |  |






<a name="fx.gravity.v1.QueryLastPendingValsetRequestByAddrRequest"></a>

### QueryLastPendingValsetRequestByAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryLastPendingValsetRequestByAddrResponse"></a>

### QueryLastPendingValsetRequestByAddrResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `valsets` | [Valset](#fx.gravity.v1.Valset) | repeated |  |






<a name="fx.gravity.v1.QueryLastValsetRequestsRequest"></a>

### QueryLastValsetRequestsRequest







<a name="fx.gravity.v1.QueryLastValsetRequestsResponse"></a>

### QueryLastValsetRequestsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `valsets` | [Valset](#fx.gravity.v1.Valset) | repeated |  |






<a name="fx.gravity.v1.QueryOutgoingTxBatchesRequest"></a>

### QueryOutgoingTxBatchesRequest







<a name="fx.gravity.v1.QueryOutgoingTxBatchesResponse"></a>

### QueryOutgoingTxBatchesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batches` | [OutgoingTxBatch](#fx.gravity.v1.OutgoingTxBatch) | repeated |  |






<a name="fx.gravity.v1.QueryParamsRequest"></a>

### QueryParamsRequest







<a name="fx.gravity.v1.QueryParamsResponse"></a>

### QueryParamsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#fx.gravity.v1.Params) |  |  |






<a name="fx.gravity.v1.QueryPendingSendToEthRequest"></a>

### QueryPendingSendToEthRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryPendingSendToEthResponse"></a>

### QueryPendingSendToEthResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `transfers_in_batches` | [OutgoingTransferTx](#fx.gravity.v1.OutgoingTransferTx) | repeated |  |
| `unbatched_transfers` | [OutgoingTransferTx](#fx.gravity.v1.OutgoingTransferTx) | repeated |  |






<a name="fx.gravity.v1.QueryValsetConfirmRequest"></a>

### QueryValsetConfirmRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryValsetConfirmResponse"></a>

### QueryValsetConfirmResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirm` | [MsgValsetConfirm](#fx.gravity.v1.MsgValsetConfirm) |  |  |






<a name="fx.gravity.v1.QueryValsetConfirmsByNonceRequest"></a>

### QueryValsetConfirmsByNonceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.QueryValsetConfirmsByNonceResponse"></a>

### QueryValsetConfirmsByNonceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirms` | [MsgValsetConfirm](#fx.gravity.v1.MsgValsetConfirm) | repeated |  |






<a name="fx.gravity.v1.QueryValsetRequestRequest"></a>

### QueryValsetRequestRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.QueryValsetRequestResponse"></a>

### QueryValsetRequestResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `valset` | [Valset](#fx.gravity.v1.Valset) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.gravity.v1.Query"></a>

### Query
Query defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#fx.gravity.v1.QueryParamsRequest) | [QueryParamsResponse](#fx.gravity.v1.QueryParamsResponse) | Deployments queries deployments | GET|/gravity/v1beta/params|
| `CurrentValset` | [QueryCurrentValsetRequest](#fx.gravity.v1.QueryCurrentValsetRequest) | [QueryCurrentValsetResponse](#fx.gravity.v1.QueryCurrentValsetResponse) |  | GET|/gravity/v1beta/valset/current|
| `ValsetRequest` | [QueryValsetRequestRequest](#fx.gravity.v1.QueryValsetRequestRequest) | [QueryValsetRequestResponse](#fx.gravity.v1.QueryValsetRequestResponse) |  | GET|/gravity/v1beta/valset|
| `ValsetConfirm` | [QueryValsetConfirmRequest](#fx.gravity.v1.QueryValsetConfirmRequest) | [QueryValsetConfirmResponse](#fx.gravity.v1.QueryValsetConfirmResponse) |  | GET|/gravity/v1beta/valset/confirm|
| `ValsetConfirmsByNonce` | [QueryValsetConfirmsByNonceRequest](#fx.gravity.v1.QueryValsetConfirmsByNonceRequest) | [QueryValsetConfirmsByNonceResponse](#fx.gravity.v1.QueryValsetConfirmsByNonceResponse) |  | GET|/gravity/v1beta/confirms/{nonce}|
| `LastValsetRequests` | [QueryLastValsetRequestsRequest](#fx.gravity.v1.QueryLastValsetRequestsRequest) | [QueryLastValsetRequestsResponse](#fx.gravity.v1.QueryLastValsetRequestsResponse) |  | GET|/gravity/v1beta/valset/requests|
| `LastPendingValsetRequestByAddr` | [QueryLastPendingValsetRequestByAddrRequest](#fx.gravity.v1.QueryLastPendingValsetRequestByAddrRequest) | [QueryLastPendingValsetRequestByAddrResponse](#fx.gravity.v1.QueryLastPendingValsetRequestByAddrResponse) |  | GET|/gravity/v1beta/valset/last|
| `LastPendingBatchRequestByAddr` | [QueryLastPendingBatchRequestByAddrRequest](#fx.gravity.v1.QueryLastPendingBatchRequestByAddrRequest) | [QueryLastPendingBatchRequestByAddrResponse](#fx.gravity.v1.QueryLastPendingBatchRequestByAddrResponse) |  | GET|/gravity/v1beta/batch/{address}|
| `LastEventNonceByAddr` | [QueryLastEventNonceByAddrRequest](#fx.gravity.v1.QueryLastEventNonceByAddrRequest) | [QueryLastEventNonceByAddrResponse](#fx.gravity.v1.QueryLastEventNonceByAddrResponse) |  | GET|/gravity/v1beta/oracle/event_nonce/{address}|
| `LastEventBlockHeightByAddr` | [QueryLastEventBlockHeightByAddrRequest](#fx.gravity.v1.QueryLastEventBlockHeightByAddrRequest) | [QueryLastEventBlockHeightByAddrResponse](#fx.gravity.v1.QueryLastEventBlockHeightByAddrResponse) |  | GET|/gravity/v1beta/oracle/event/block_height/{address}|
| `BatchFees` | [QueryBatchFeeRequest](#fx.gravity.v1.QueryBatchFeeRequest) | [QueryBatchFeeResponse](#fx.gravity.v1.QueryBatchFeeResponse) |  | GET|/gravity/v1beta/batch_fees|
| `LastObservedEthBlockHeight` | [QueryLastObservedEthBlockHeightRequest](#fx.gravity.v1.QueryLastObservedEthBlockHeightRequest) | [QueryLastObservedEthBlockHeightResponse](#fx.gravity.v1.QueryLastObservedEthBlockHeightResponse) |  | GET|/gravity/v1beta/eth/block_height|
| `OutgoingTxBatches` | [QueryOutgoingTxBatchesRequest](#fx.gravity.v1.QueryOutgoingTxBatchesRequest) | [QueryOutgoingTxBatchesResponse](#fx.gravity.v1.QueryOutgoingTxBatchesResponse) |  | GET|/gravity/v1beta/batch/outgoing_tx|
| `BatchRequestByNonce` | [QueryBatchRequestByNonceRequest](#fx.gravity.v1.QueryBatchRequestByNonceRequest) | [QueryBatchRequestByNonceResponse](#fx.gravity.v1.QueryBatchRequestByNonceResponse) |  | GET|/gravity/v1beta/batch/{nonce}|
| `BatchConfirm` | [QueryBatchConfirmRequest](#fx.gravity.v1.QueryBatchConfirmRequest) | [QueryBatchConfirmResponse](#fx.gravity.v1.QueryBatchConfirmResponse) |  | GET|/gravity/v1beta/batch/confirm|
| `BatchConfirms` | [QueryBatchConfirmsRequest](#fx.gravity.v1.QueryBatchConfirmsRequest) | [QueryBatchConfirmsResponse](#fx.gravity.v1.QueryBatchConfirmsResponse) |  | GET|/gravity/v1beta/batch/confirms|
| `ERC20ToDenom` | [QueryERC20ToDenomRequest](#fx.gravity.v1.QueryERC20ToDenomRequest) | [QueryERC20ToDenomResponse](#fx.gravity.v1.QueryERC20ToDenomResponse) |  | GET|/gravity/v1beta/fx_originated/erc20_to_denom|
| `DenomToERC20` | [QueryDenomToERC20Request](#fx.gravity.v1.QueryDenomToERC20Request) | [QueryDenomToERC20Response](#fx.gravity.v1.QueryDenomToERC20Response) |  | GET|/gravity/v1beta/fx_originated/denom_to_erc20|
| `GetDelegateKeyByValidator` | [QueryDelegateKeyByValidatorRequest](#fx.gravity.v1.QueryDelegateKeyByValidatorRequest) | [QueryDelegateKeyByValidatorResponse](#fx.gravity.v1.QueryDelegateKeyByValidatorResponse) |  | GET|/gravity/v1beta/query_delegate_key_by_validator|
| `GetDelegateKeyByEth` | [QueryDelegateKeyByEthRequest](#fx.gravity.v1.QueryDelegateKeyByEthRequest) | [QueryDelegateKeyByEthResponse](#fx.gravity.v1.QueryDelegateKeyByEthResponse) |  | GET|/gravity/v1beta/query_delegate_key_by_eth|
| `GetDelegateKeyByOrchestrator` | [QueryDelegateKeyByOrchestratorRequest](#fx.gravity.v1.QueryDelegateKeyByOrchestratorRequest) | [QueryDelegateKeyByOrchestratorResponse](#fx.gravity.v1.QueryDelegateKeyByOrchestratorResponse) |  | GET|/gravity/v1beta/query_delegate_key_by_orchestrator|
| `GetPendingSendToEth` | [QueryPendingSendToEthRequest](#fx.gravity.v1.QueryPendingSendToEthRequest) | [QueryPendingSendToEthResponse](#fx.gravity.v1.QueryPendingSendToEthResponse) |  | GET|/gravity/v1beta/query_pending_send_to_eth|
| `GetIbcSequenceHeightByChannel` | [QueryIbcSequenceHeightRequest](#fx.gravity.v1.QueryIbcSequenceHeightRequest) | [QueryIbcSequenceHeightResponse](#fx.gravity.v1.QueryIbcSequenceHeightResponse) |  | GET|/gravity/v1beta/query_ibc_sequence_height|

 <!-- end services -->



<a name="fx/other/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/other/query.proto



<a name="fx.other.GasPriceRequest"></a>

### GasPriceRequest







<a name="fx.other.GasPriceResponse"></a>

### GasPriceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas_prices` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.other.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `GasPrice` | [GasPriceRequest](#fx.other.GasPriceRequest) | [GasPriceResponse](#fx.other.GasPriceResponse) |  | GET|/other/v1/gas_price|

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
