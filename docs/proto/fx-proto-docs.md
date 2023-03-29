<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [fx/base/v1/query.proto](#fx/base/v1/query.proto)
    - [GetGasPriceRequest](#fx.base.v1.GetGasPriceRequest)
    - [GetGasPriceResponse](#fx.base.v1.GetGasPriceResponse)
  
    - [Query](#fx.base.v1.Query)
  
- [fx/crosschain/v1/types.proto](#fx/crosschain/v1/types.proto)
    - [Attestation](#fx.gravity.crosschain.v1.Attestation)
    - [BatchFees](#fx.gravity.crosschain.v1.BatchFees)
    - [BridgeToken](#fx.gravity.crosschain.v1.BridgeToken)
    - [BridgeValidator](#fx.gravity.crosschain.v1.BridgeValidator)
    - [ERC20Token](#fx.gravity.crosschain.v1.ERC20Token)
    - [IDSet](#fx.gravity.crosschain.v1.IDSet)
    - [InitCrossChainParamsProposal](#fx.gravity.crosschain.v1.InitCrossChainParamsProposal)
    - [LastObservedBlockHeight](#fx.gravity.crosschain.v1.LastObservedBlockHeight)
    - [MinBatchFee](#fx.gravity.crosschain.v1.MinBatchFee)
    - [Oracle](#fx.gravity.crosschain.v1.Oracle)
    - [OracleSet](#fx.gravity.crosschain.v1.OracleSet)
    - [OutgoingTransferTx](#fx.gravity.crosschain.v1.OutgoingTransferTx)
    - [OutgoingTxBatch](#fx.gravity.crosschain.v1.OutgoingTxBatch)
    - [Params](#fx.gravity.crosschain.v1.Params)
    - [ProposalOracle](#fx.gravity.crosschain.v1.ProposalOracle)
    - [UpdateChainOraclesProposal](#fx.gravity.crosschain.v1.UpdateChainOraclesProposal)
  
    - [ClaimType](#fx.gravity.crosschain.v1.ClaimType)
  
- [fx/crosschain/v1/tx.proto](#fx/crosschain/v1/tx.proto)
    - [MsgAddDelegate](#fx.gravity.crosschain.v1.MsgAddDelegate)
    - [MsgAddDelegateResponse](#fx.gravity.crosschain.v1.MsgAddDelegateResponse)
    - [MsgAddOracleDeposit](#fx.gravity.crosschain.v1.MsgAddOracleDeposit)
    - [MsgBondedOracle](#fx.gravity.crosschain.v1.MsgBondedOracle)
    - [MsgBondedOracleResponse](#fx.gravity.crosschain.v1.MsgBondedOracleResponse)
    - [MsgBridgeTokenClaim](#fx.gravity.crosschain.v1.MsgBridgeTokenClaim)
    - [MsgBridgeTokenClaimResponse](#fx.gravity.crosschain.v1.MsgBridgeTokenClaimResponse)
    - [MsgCancelSendToExternal](#fx.gravity.crosschain.v1.MsgCancelSendToExternal)
    - [MsgCancelSendToExternalResponse](#fx.gravity.crosschain.v1.MsgCancelSendToExternalResponse)
    - [MsgConfirmBatch](#fx.gravity.crosschain.v1.MsgConfirmBatch)
    - [MsgConfirmBatchResponse](#fx.gravity.crosschain.v1.MsgConfirmBatchResponse)
    - [MsgEditBridger](#fx.gravity.crosschain.v1.MsgEditBridger)
    - [MsgEditBridgerResponse](#fx.gravity.crosschain.v1.MsgEditBridgerResponse)
    - [MsgIncreaseBridgeFee](#fx.gravity.crosschain.v1.MsgIncreaseBridgeFee)
    - [MsgIncreaseBridgeFeeResponse](#fx.gravity.crosschain.v1.MsgIncreaseBridgeFeeResponse)
    - [MsgOracleSetConfirm](#fx.gravity.crosschain.v1.MsgOracleSetConfirm)
    - [MsgOracleSetConfirmResponse](#fx.gravity.crosschain.v1.MsgOracleSetConfirmResponse)
    - [MsgOracleSetUpdatedClaim](#fx.gravity.crosschain.v1.MsgOracleSetUpdatedClaim)
    - [MsgOracleSetUpdatedClaimResponse](#fx.gravity.crosschain.v1.MsgOracleSetUpdatedClaimResponse)
    - [MsgReDelegate](#fx.gravity.crosschain.v1.MsgReDelegate)
    - [MsgReDelegateResponse](#fx.gravity.crosschain.v1.MsgReDelegateResponse)
    - [MsgRequestBatch](#fx.gravity.crosschain.v1.MsgRequestBatch)
    - [MsgRequestBatchResponse](#fx.gravity.crosschain.v1.MsgRequestBatchResponse)
    - [MsgSendToExternal](#fx.gravity.crosschain.v1.MsgSendToExternal)
    - [MsgSendToExternalClaim](#fx.gravity.crosschain.v1.MsgSendToExternalClaim)
    - [MsgSendToExternalClaimResponse](#fx.gravity.crosschain.v1.MsgSendToExternalClaimResponse)
    - [MsgSendToExternalResponse](#fx.gravity.crosschain.v1.MsgSendToExternalResponse)
    - [MsgSendToFxClaim](#fx.gravity.crosschain.v1.MsgSendToFxClaim)
    - [MsgSendToFxClaimResponse](#fx.gravity.crosschain.v1.MsgSendToFxClaimResponse)
    - [MsgSetOrchestratorAddress](#fx.gravity.crosschain.v1.MsgSetOrchestratorAddress)
    - [MsgUnbondedOracle](#fx.gravity.crosschain.v1.MsgUnbondedOracle)
    - [MsgUnbondedOracleResponse](#fx.gravity.crosschain.v1.MsgUnbondedOracleResponse)
    - [MsgUpdateParams](#fx.gravity.crosschain.v1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#fx.gravity.crosschain.v1.MsgUpdateParamsResponse)
    - [MsgWithdrawReward](#fx.gravity.crosschain.v1.MsgWithdrawReward)
    - [MsgWithdrawRewardResponse](#fx.gravity.crosschain.v1.MsgWithdrawRewardResponse)
  
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
    - [QueryBridgeTokensRequest](#fx.gravity.crosschain.v1.QueryBridgeTokensRequest)
    - [QueryBridgeTokensResponse](#fx.gravity.crosschain.v1.QueryBridgeTokensResponse)
    - [QueryCurrentOracleSetRequest](#fx.gravity.crosschain.v1.QueryCurrentOracleSetRequest)
    - [QueryCurrentOracleSetResponse](#fx.gravity.crosschain.v1.QueryCurrentOracleSetResponse)
    - [QueryDenomToTokenRequest](#fx.gravity.crosschain.v1.QueryDenomToTokenRequest)
    - [QueryDenomToTokenResponse](#fx.gravity.crosschain.v1.QueryDenomToTokenResponse)
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
    - [QueryOracleByBridgerAddrRequest](#fx.gravity.crosschain.v1.QueryOracleByBridgerAddrRequest)
    - [QueryOracleByExternalAddrRequest](#fx.gravity.crosschain.v1.QueryOracleByExternalAddrRequest)
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
    - [QueryProjectedBatchTimeoutHeightRequest](#fx.gravity.crosschain.v1.QueryProjectedBatchTimeoutHeightRequest)
    - [QueryProjectedBatchTimeoutHeightResponse](#fx.gravity.crosschain.v1.QueryProjectedBatchTimeoutHeightResponse)
    - [QueryTokenToDenomRequest](#fx.gravity.crosschain.v1.QueryTokenToDenomRequest)
    - [QueryTokenToDenomResponse](#fx.gravity.crosschain.v1.QueryTokenToDenomResponse)
  
    - [Query](#fx.gravity.crosschain.v1.Query)
  
- [fx/erc20/v1/erc20.proto](#fx/erc20/v1/erc20.proto)
    - [RegisterCoinProposal](#fx.erc20.v1.RegisterCoinProposal)
    - [RegisterERC20Proposal](#fx.erc20.v1.RegisterERC20Proposal)
    - [ToggleTokenConversionProposal](#fx.erc20.v1.ToggleTokenConversionProposal)
    - [TokenPair](#fx.erc20.v1.TokenPair)
    - [UpdateDenomAliasProposal](#fx.erc20.v1.UpdateDenomAliasProposal)
  
    - [Owner](#fx.erc20.v1.Owner)
  
- [fx/erc20/v1/genesis.proto](#fx/erc20/v1/genesis.proto)
    - [GenesisState](#fx.erc20.v1.GenesisState)
    - [Params](#fx.erc20.v1.Params)
  
- [fx/erc20/v1/query.proto](#fx/erc20/v1/query.proto)
    - [QueryAliasDenomRequest](#fx.erc20.v1.QueryAliasDenomRequest)
    - [QueryAliasDenomResponse](#fx.erc20.v1.QueryAliasDenomResponse)
    - [QueryDenomAliasesRequest](#fx.erc20.v1.QueryDenomAliasesRequest)
    - [QueryDenomAliasesResponse](#fx.erc20.v1.QueryDenomAliasesResponse)
    - [QueryParamsRequest](#fx.erc20.v1.QueryParamsRequest)
    - [QueryParamsResponse](#fx.erc20.v1.QueryParamsResponse)
    - [QueryTokenPairRequest](#fx.erc20.v1.QueryTokenPairRequest)
    - [QueryTokenPairResponse](#fx.erc20.v1.QueryTokenPairResponse)
    - [QueryTokenPairsRequest](#fx.erc20.v1.QueryTokenPairsRequest)
    - [QueryTokenPairsResponse](#fx.erc20.v1.QueryTokenPairsResponse)
  
    - [Query](#fx.erc20.v1.Query)
  
- [fx/erc20/v1/tx.proto](#fx/erc20/v1/tx.proto)
    - [MsgConvertCoin](#fx.erc20.v1.MsgConvertCoin)
    - [MsgConvertCoinResponse](#fx.erc20.v1.MsgConvertCoinResponse)
    - [MsgConvertDenom](#fx.erc20.v1.MsgConvertDenom)
    - [MsgConvertDenomResponse](#fx.erc20.v1.MsgConvertDenomResponse)
    - [MsgConvertERC20](#fx.erc20.v1.MsgConvertERC20)
    - [MsgConvertERC20Response](#fx.erc20.v1.MsgConvertERC20Response)
    - [MsgUpdateParams](#fx.erc20.v1.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#fx.erc20.v1.MsgUpdateParamsResponse)
  
    - [Msg](#fx.erc20.v1.Msg)
  
- [fx/ibc/applications/transfer/v1/query.proto](#fx/ibc/applications/transfer/v1/query.proto)
    - [Query](#fx.ibc.applications.transfer.v1.Query)
  
- [fx/ibc/applications/transfer/v1/transfer.proto](#fx/ibc/applications/transfer/v1/transfer.proto)
    - [FungibleTokenPacketData](#fx.ibc.applications.transfer.v1.FungibleTokenPacketData)
  
- [fx/ibc/applications/transfer/v1/tx.proto](#fx/ibc/applications/transfer/v1/tx.proto)
    - [MsgTransfer](#fx.ibc.applications.transfer.v1.MsgTransfer)
  
    - [Msg](#fx.ibc.applications.transfer.v1.Msg)
  
- [fx/legacy/gravity/v1/types.proto](#fx/legacy/gravity/v1/types.proto)
    - [Attestation](#fx.gravity.v1.Attestation)
    - [BatchFees](#fx.gravity.v1.BatchFees)
    - [BridgeValidator](#fx.gravity.v1.BridgeValidator)
    - [ERC20ToDenom](#fx.gravity.v1.ERC20ToDenom)
    - [ERC20Token](#fx.gravity.v1.ERC20Token)
    - [LastObservedEthereumBlockHeight](#fx.gravity.v1.LastObservedEthereumBlockHeight)
    - [MinBatchFee](#fx.gravity.v1.MinBatchFee)
    - [OutgoingTransferTx](#fx.gravity.v1.OutgoingTransferTx)
    - [OutgoingTxBatch](#fx.gravity.v1.OutgoingTxBatch)
    - [Valset](#fx.gravity.v1.Valset)
  
    - [ClaimType](#fx.gravity.v1.ClaimType)
  
- [fx/legacy/gravity/v1/tx.proto](#fx/legacy/gravity/v1/tx.proto)
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
  
- [fx/legacy/gravity/v1/genesis.proto](#fx/legacy/gravity/v1/genesis.proto)
    - [GenesisState](#fx.gravity.v1.GenesisState)
    - [Params](#fx.gravity.v1.Params)
  
- [fx/legacy/gravity/v1/query.proto](#fx/legacy/gravity/v1/query.proto)
    - [QueryBatchConfirmRequest](#fx.gravity.v1.QueryBatchConfirmRequest)
    - [QueryBatchConfirmResponse](#fx.gravity.v1.QueryBatchConfirmResponse)
    - [QueryBatchConfirmsRequest](#fx.gravity.v1.QueryBatchConfirmsRequest)
    - [QueryBatchConfirmsResponse](#fx.gravity.v1.QueryBatchConfirmsResponse)
    - [QueryBatchFeeRequest](#fx.gravity.v1.QueryBatchFeeRequest)
    - [QueryBatchFeeResponse](#fx.gravity.v1.QueryBatchFeeResponse)
    - [QueryBatchRequestByNonceRequest](#fx.gravity.v1.QueryBatchRequestByNonceRequest)
    - [QueryBatchRequestByNonceResponse](#fx.gravity.v1.QueryBatchRequestByNonceResponse)
    - [QueryBridgeTokensRequest](#fx.gravity.v1.QueryBridgeTokensRequest)
    - [QueryBridgeTokensResponse](#fx.gravity.v1.QueryBridgeTokensResponse)
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
    - [QueryLastEventBlockHeightByAddrRequest](#fx.gravity.v1.QueryLastEventBlockHeightByAddrRequest)
    - [QueryLastEventBlockHeightByAddrResponse](#fx.gravity.v1.QueryLastEventBlockHeightByAddrResponse)
    - [QueryLastEventNonceByAddrRequest](#fx.gravity.v1.QueryLastEventNonceByAddrRequest)
    - [QueryLastEventNonceByAddrResponse](#fx.gravity.v1.QueryLastEventNonceByAddrResponse)
    - [QueryLastObservedBlockHeightRequest](#fx.gravity.v1.QueryLastObservedBlockHeightRequest)
    - [QueryLastObservedBlockHeightResponse](#fx.gravity.v1.QueryLastObservedBlockHeightResponse)
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
    - [QueryProjectedBatchTimeoutHeightRequest](#fx.gravity.v1.QueryProjectedBatchTimeoutHeightRequest)
    - [QueryProjectedBatchTimeoutHeightResponse](#fx.gravity.v1.QueryProjectedBatchTimeoutHeightResponse)
    - [QueryValsetConfirmRequest](#fx.gravity.v1.QueryValsetConfirmRequest)
    - [QueryValsetConfirmResponse](#fx.gravity.v1.QueryValsetConfirmResponse)
    - [QueryValsetConfirmsByNonceRequest](#fx.gravity.v1.QueryValsetConfirmsByNonceRequest)
    - [QueryValsetConfirmsByNonceResponse](#fx.gravity.v1.QueryValsetConfirmsByNonceResponse)
    - [QueryValsetRequestRequest](#fx.gravity.v1.QueryValsetRequestRequest)
    - [QueryValsetRequestResponse](#fx.gravity.v1.QueryValsetRequestResponse)
  
    - [Query](#fx.gravity.v1.Query)
  
- [fx/legacy/other/query.proto](#fx/legacy/other/query.proto)
    - [GasPriceRequest](#fx.other.GasPriceRequest)
    - [GasPriceResponse](#fx.other.GasPriceResponse)
  
    - [Query](#fx.other.Query)
  
- [fx/migrate/v1/migrate.proto](#fx/migrate/v1/migrate.proto)
    - [MigrateRecord](#fx.migrate.v1.MigrateRecord)
  
- [fx/migrate/v1/genesis.proto](#fx/migrate/v1/genesis.proto)
    - [GenesisState](#fx.migrate.v1.GenesisState)
  
- [fx/migrate/v1/query.proto](#fx/migrate/v1/query.proto)
    - [QueryMigrateCheckAccountRequest](#fx.migrate.v1.QueryMigrateCheckAccountRequest)
    - [QueryMigrateCheckAccountResponse](#fx.migrate.v1.QueryMigrateCheckAccountResponse)
    - [QueryMigrateRecordRequest](#fx.migrate.v1.QueryMigrateRecordRequest)
    - [QueryMigrateRecordResponse](#fx.migrate.v1.QueryMigrateRecordResponse)
  
    - [Query](#fx.migrate.v1.Query)
  
- [fx/migrate/v1/tx.proto](#fx/migrate/v1/tx.proto)
    - [MsgMigrateAccount](#fx.migrate.v1.MsgMigrateAccount)
    - [MsgMigrateAccountResponse](#fx.migrate.v1.MsgMigrateAccountResponse)
  
    - [Msg](#fx.migrate.v1.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="fx/base/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/base/v1/query.proto



<a name="fx.base.v1.GetGasPriceRequest"></a>

### GetGasPriceRequest







<a name="fx.base.v1.GetGasPriceResponse"></a>

### GetGasPriceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas_prices` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.base.v1.Query"></a>

### Query


| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `GetGasPrice` | [GetGasPriceRequest](#fx.base.v1.GetGasPriceRequest) | [GetGasPriceResponse](#fx.base.v1.GetGasPriceResponse) |  | GET|/fx/base/v1/gas_price|

 <!-- end services -->



<a name="fx/crosschain/v1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/crosschain/v1/types.proto



<a name="fx.gravity.crosschain.v1.Attestation"></a>

### Attestation
Attestation is an aggregate of `claims` that eventually becomes `observed` by
all bridger set
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
| `total_amount` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.BridgeToken"></a>

### BridgeToken
BridgeToken


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token` | [string](#string) |  |  |
| `denom` | [string](#string) |  |  |
| `channel_ibc` | [string](#string) |  | Deprecated: after upgrade v3 |






<a name="fx.gravity.crosschain.v1.BridgeValidator"></a>

### BridgeValidator
BridgeValidator represents a validator's external address and its power


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `power` | [uint64](#uint64) |  |  |
| `external_address` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.ERC20Token"></a>

### ERC20Token
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
Deprecated: after block 5713000


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  |  |
| `description` | [string](#string) |  |  |
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






<a name="fx.gravity.crosschain.v1.MinBatchFee"></a>

### MinBatchFee



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_contract` | [string](#string) |  |  |
| `baseFee` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.Oracle"></a>

### Oracle



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_address` | [string](#string) |  |  |
| `bridger_address` | [string](#string) |  |  |
| `external_address` | [string](#string) |  |  |
| `delegate_amount` | [string](#string) |  |  |
| `start_height` | [int64](#int64) |  | start oracle height |
| `online` | [bool](#bool) |  |  |
| `delegate_validator` | [string](#string) |  |  |
| `slash_times` | [int64](#int64) |  |  |






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
| `token` | [ERC20Token](#fx.gravity.crosschain.v1.ERC20Token) |  |  |
| `fee` | [ERC20Token](#fx.gravity.crosschain.v1.ERC20Token) |  |  |






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
| `oracles` | [string](#string) | repeated | Deprecated: after block 5713000 |
| `delegate_threshold` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `delegate_multiple` | [int64](#int64) |  |  |






<a name="fx.gravity.crosschain.v1.ProposalOracle"></a>

### ProposalOracle



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracles` | [string](#string) | repeated |  |






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


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="fx/crosschain/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/crosschain/v1/tx.proto



<a name="fx.gravity.crosschain.v1.MsgAddDelegate"></a>

### MsgAddDelegate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `oracle_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="fx.gravity.crosschain.v1.MsgAddDelegateResponse"></a>

### MsgAddDelegateResponse







<a name="fx.gravity.crosschain.v1.MsgAddOracleDeposit"></a>

### MsgAddOracleDeposit
Deprecated: after block 5713000


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_address` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgBondedOracle"></a>

### MsgBondedOracle



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `oracle_address` | [string](#string) |  |  |
| `bridger_address` | [string](#string) |  |  |
| `external_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |
| `delegate_amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="fx.gravity.crosschain.v1.MsgBondedOracleResponse"></a>

### MsgBondedOracleResponse







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
| `bridger_address` | [string](#string) |  |  |
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
| `bridger_address` | [string](#string) |  |  |
| `external_address` | [string](#string) |  |  |
| `signature` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgConfirmBatchResponse"></a>

### MsgConfirmBatchResponse







<a name="fx.gravity.crosschain.v1.MsgEditBridger"></a>

### MsgEditBridger



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `oracle_address` | [string](#string) |  |  |
| `bridger_address` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgEditBridgerResponse"></a>

### MsgEditBridgerResponse







<a name="fx.gravity.crosschain.v1.MsgIncreaseBridgeFee"></a>

### MsgIncreaseBridgeFee



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `transaction_id` | [uint64](#uint64) |  |  |
| `sender` | [string](#string) |  |  |
| `add_bridge_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="fx.gravity.crosschain.v1.MsgIncreaseBridgeFeeResponse"></a>

### MsgIncreaseBridgeFeeResponse







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
| `bridger_address` | [string](#string) |  |  |
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
| `bridger_address` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgOracleSetUpdatedClaimResponse"></a>

### MsgOracleSetUpdatedClaimResponse







<a name="fx.gravity.crosschain.v1.MsgReDelegate"></a>

### MsgReDelegate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `oracle_address` | [string](#string) |  |  |
| `validator_address` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgReDelegateResponse"></a>

### MsgReDelegateResponse







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
| `fee_receive` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |
| `base_fee` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgRequestBatchResponse"></a>

### MsgRequestBatchResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch_nonce` | [uint64](#uint64) |  |  |






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
| `bridger_address` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgSendToExternalClaimResponse"></a>

### MsgSendToExternalClaimResponse







<a name="fx.gravity.crosschain.v1.MsgSendToExternalResponse"></a>

### MsgSendToExternalResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `outgoing_tx_id` | [uint64](#uint64) |  |  |






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
| `bridger_address` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgSendToFxClaimResponse"></a>

### MsgSendToFxClaimResponse







<a name="fx.gravity.crosschain.v1.MsgSetOrchestratorAddress"></a>

### MsgSetOrchestratorAddress
Deprecated: after block 5713000


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_address` | [string](#string) |  |  |
| `bridger_address` | [string](#string) |  |  |
| `external_address` | [string](#string) |  |  |
| `deposit` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgUnbondedOracle"></a>

### MsgUnbondedOracle



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `oracle_address` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgUnbondedOracleResponse"></a>

### MsgUnbondedOracleResponse







<a name="fx.gravity.crosschain.v1.MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams is the Msg/UpdateParams request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `authority` | [string](#string) |  | authority is the address of the governance account. |
| `params` | [Params](#fx.gravity.crosschain.v1.Params) |  | params defines the x/crosschain parameters to update.

NOTE: All parameters must be supplied. |






<a name="fx.gravity.crosschain.v1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse
MsgUpdateParamsResponse defines the response structure for executing a
MsgUpdateParams message.






<a name="fx.gravity.crosschain.v1.MsgWithdrawReward"></a>

### MsgWithdrawReward



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `oracle_address` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.MsgWithdrawRewardResponse"></a>

### MsgWithdrawRewardResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.gravity.crosschain.v1.Msg"></a>

### Msg
Msg defines the state transitions possible within gravity

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `BondedOracle` | [MsgBondedOracle](#fx.gravity.crosschain.v1.MsgBondedOracle) | [MsgBondedOracleResponse](#fx.gravity.crosschain.v1.MsgBondedOracleResponse) |  | |
| `AddDelegate` | [MsgAddDelegate](#fx.gravity.crosschain.v1.MsgAddDelegate) | [MsgAddDelegateResponse](#fx.gravity.crosschain.v1.MsgAddDelegateResponse) |  | |
| `ReDelegate` | [MsgReDelegate](#fx.gravity.crosschain.v1.MsgReDelegate) | [MsgReDelegateResponse](#fx.gravity.crosschain.v1.MsgReDelegateResponse) |  | |
| `EditBridger` | [MsgEditBridger](#fx.gravity.crosschain.v1.MsgEditBridger) | [MsgEditBridgerResponse](#fx.gravity.crosschain.v1.MsgEditBridgerResponse) |  | |
| `WithdrawReward` | [MsgWithdrawReward](#fx.gravity.crosschain.v1.MsgWithdrawReward) | [MsgWithdrawRewardResponse](#fx.gravity.crosschain.v1.MsgWithdrawRewardResponse) |  | |
| `UnbondedOracle` | [MsgUnbondedOracle](#fx.gravity.crosschain.v1.MsgUnbondedOracle) | [MsgUnbondedOracleResponse](#fx.gravity.crosschain.v1.MsgUnbondedOracleResponse) |  | |
| `OracleSetConfirm` | [MsgOracleSetConfirm](#fx.gravity.crosschain.v1.MsgOracleSetConfirm) | [MsgOracleSetConfirmResponse](#fx.gravity.crosschain.v1.MsgOracleSetConfirmResponse) |  | |
| `OracleSetUpdateClaim` | [MsgOracleSetUpdatedClaim](#fx.gravity.crosschain.v1.MsgOracleSetUpdatedClaim) | [MsgOracleSetUpdatedClaimResponse](#fx.gravity.crosschain.v1.MsgOracleSetUpdatedClaimResponse) |  | |
| `BridgeTokenClaim` | [MsgBridgeTokenClaim](#fx.gravity.crosschain.v1.MsgBridgeTokenClaim) | [MsgBridgeTokenClaimResponse](#fx.gravity.crosschain.v1.MsgBridgeTokenClaimResponse) |  | |
| `SendToFxClaim` | [MsgSendToFxClaim](#fx.gravity.crosschain.v1.MsgSendToFxClaim) | [MsgSendToFxClaimResponse](#fx.gravity.crosschain.v1.MsgSendToFxClaimResponse) |  | |
| `SendToExternal` | [MsgSendToExternal](#fx.gravity.crosschain.v1.MsgSendToExternal) | [MsgSendToExternalResponse](#fx.gravity.crosschain.v1.MsgSendToExternalResponse) |  | |
| `CancelSendToExternal` | [MsgCancelSendToExternal](#fx.gravity.crosschain.v1.MsgCancelSendToExternal) | [MsgCancelSendToExternalResponse](#fx.gravity.crosschain.v1.MsgCancelSendToExternalResponse) |  | |
| `SendToExternalClaim` | [MsgSendToExternalClaim](#fx.gravity.crosschain.v1.MsgSendToExternalClaim) | [MsgSendToExternalClaimResponse](#fx.gravity.crosschain.v1.MsgSendToExternalClaimResponse) |  | |
| `RequestBatch` | [MsgRequestBatch](#fx.gravity.crosschain.v1.MsgRequestBatch) | [MsgRequestBatchResponse](#fx.gravity.crosschain.v1.MsgRequestBatchResponse) |  | |
| `ConfirmBatch` | [MsgConfirmBatch](#fx.gravity.crosschain.v1.MsgConfirmBatch) | [MsgConfirmBatchResponse](#fx.gravity.crosschain.v1.MsgConfirmBatchResponse) |  | |
| `UpdateParams` | [MsgUpdateParams](#fx.gravity.crosschain.v1.MsgUpdateParams) | [MsgUpdateParamsResponse](#fx.gravity.crosschain.v1.MsgUpdateParamsResponse) | UpdateParams defines a governance operation for updating the x/crosschain module parameters. The authority is hard-coded to the x/gov module account. | |
| `IncreaseBridgeFee` | [MsgIncreaseBridgeFee](#fx.gravity.crosschain.v1.MsgIncreaseBridgeFee) | [MsgIncreaseBridgeFeeResponse](#fx.gravity.crosschain.v1.MsgIncreaseBridgeFeeResponse) |  | |

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
| `proposal_oracle` | [ProposalOracle](#fx.gravity.crosschain.v1.ProposalOracle) |  |  |
| `last_observed_oracle_set` | [OracleSet](#fx.gravity.crosschain.v1.OracleSet) |  |  |
| `last_slashed_batch_block` | [uint64](#uint64) |  |  |
| `last_slashed_oracle_set_nonce` | [uint64](#uint64) |  |  |





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
| `chain_name` | [string](#string) |  |  |
| `token_contract` | [string](#string) |  |  |
| `bridger_address` | [string](#string) |  |  |
| `nonce` | [uint64](#uint64) |  |  |






<a name="fx.gravity.crosschain.v1.QueryBatchConfirmResponse"></a>

### QueryBatchConfirmResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirm` | [MsgConfirmBatch](#fx.gravity.crosschain.v1.MsgConfirmBatch) |  |  |






<a name="fx.gravity.crosschain.v1.QueryBatchConfirmsRequest"></a>

### QueryBatchConfirmsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `token_contract` | [string](#string) |  |  |
| `nonce` | [uint64](#uint64) |  |  |






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
| `minBatchFees` | [MinBatchFee](#fx.gravity.crosschain.v1.MinBatchFee) | repeated |  |






<a name="fx.gravity.crosschain.v1.QueryBatchFeeResponse"></a>

### QueryBatchFeeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch_fees` | [BatchFees](#fx.gravity.crosschain.v1.BatchFees) | repeated |  |






<a name="fx.gravity.crosschain.v1.QueryBatchRequestByNonceRequest"></a>

### QueryBatchRequestByNonceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `token_contract` | [string](#string) |  |  |
| `nonce` | [uint64](#uint64) |  |  |






<a name="fx.gravity.crosschain.v1.QueryBatchRequestByNonceResponse"></a>

### QueryBatchRequestByNonceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch` | [OutgoingTxBatch](#fx.gravity.crosschain.v1.OutgoingTxBatch) |  |  |






<a name="fx.gravity.crosschain.v1.QueryBridgeTokensRequest"></a>

### QueryBridgeTokensRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryBridgeTokensResponse"></a>

### QueryBridgeTokensResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bridge_tokens` | [BridgeToken](#fx.gravity.crosschain.v1.BridgeToken) | repeated |  |






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
| `chain_name` | [string](#string) |  |  |
| `denom` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryDenomToTokenResponse"></a>

### QueryDenomToTokenResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastEventBlockHeightByAddrRequest"></a>

### QueryLastEventBlockHeightByAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `bridger_address` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastEventBlockHeightByAddrResponse"></a>

### QueryLastEventBlockHeightByAddrResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_height` | [uint64](#uint64) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastEventNonceByAddrRequest"></a>

### QueryLastEventNonceByAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `bridger_address` | [string](#string) |  |  |






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
| `chain_name` | [string](#string) |  |  |
| `bridger_address` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastPendingBatchRequestByAddrResponse"></a>

### QueryLastPendingBatchRequestByAddrResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch` | [OutgoingTxBatch](#fx.gravity.crosschain.v1.OutgoingTxBatch) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastPendingOracleSetRequestByAddrRequest"></a>

### QueryLastPendingOracleSetRequestByAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `bridger_address` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryLastPendingOracleSetRequestByAddrResponse"></a>

### QueryLastPendingOracleSetRequestByAddrResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle_sets` | [OracleSet](#fx.gravity.crosschain.v1.OracleSet) | repeated |  |






<a name="fx.gravity.crosschain.v1.QueryOracleByAddrRequest"></a>

### QueryOracleByAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `oracle_address` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOracleByBridgerAddrRequest"></a>

### QueryOracleByBridgerAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `bridger_address` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOracleByExternalAddrRequest"></a>

### QueryOracleByExternalAddrRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `external_address` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOracleResponse"></a>

### QueryOracleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `oracle` | [Oracle](#fx.gravity.crosschain.v1.Oracle) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOracleSetConfirmRequest"></a>

### QueryOracleSetConfirmRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `bridger_address` | [string](#string) |  |  |
| `nonce` | [uint64](#uint64) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOracleSetConfirmResponse"></a>

### QueryOracleSetConfirmResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirm` | [MsgOracleSetConfirm](#fx.gravity.crosschain.v1.MsgOracleSetConfirm) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOracleSetConfirmsByNonceRequest"></a>

### QueryOracleSetConfirmsByNonceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `nonce` | [uint64](#uint64) |  |  |






<a name="fx.gravity.crosschain.v1.QueryOracleSetConfirmsByNonceResponse"></a>

### QueryOracleSetConfirmsByNonceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirms` | [MsgOracleSetConfirm](#fx.gravity.crosschain.v1.MsgOracleSetConfirm) | repeated |  |






<a name="fx.gravity.crosschain.v1.QueryOracleSetRequestRequest"></a>

### QueryOracleSetRequestRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `nonce` | [uint64](#uint64) |  |  |






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
| `chain_name` | [string](#string) |  |  |
| `sender_address` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryPendingSendToExternalResponse"></a>

### QueryPendingSendToExternalResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `transfers_in_batches` | [OutgoingTransferTx](#fx.gravity.crosschain.v1.OutgoingTransferTx) | repeated |  |
| `unbatched_transfers` | [OutgoingTransferTx](#fx.gravity.crosschain.v1.OutgoingTransferTx) | repeated |  |






<a name="fx.gravity.crosschain.v1.QueryProjectedBatchTimeoutHeightRequest"></a>

### QueryProjectedBatchTimeoutHeightRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryProjectedBatchTimeoutHeightResponse"></a>

### QueryProjectedBatchTimeoutHeightResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `timeout_height` | [uint64](#uint64) |  |  |






<a name="fx.gravity.crosschain.v1.QueryTokenToDenomRequest"></a>

### QueryTokenToDenomRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_name` | [string](#string) |  |  |
| `token` | [string](#string) |  |  |






<a name="fx.gravity.crosschain.v1.QueryTokenToDenomResponse"></a>

### QueryTokenToDenomResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.gravity.crosschain.v1.Query"></a>

### Query
Query defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#fx.gravity.crosschain.v1.QueryParamsRequest) | [QueryParamsResponse](#fx.gravity.crosschain.v1.QueryParamsResponse) | Deployments queries deployments | GET|/fx/crosschain/v1/params|
| `CurrentOracleSet` | [QueryCurrentOracleSetRequest](#fx.gravity.crosschain.v1.QueryCurrentOracleSetRequest) | [QueryCurrentOracleSetResponse](#fx.gravity.crosschain.v1.QueryCurrentOracleSetResponse) |  | GET|/fx/crosschain/v1/oracle_set/current|
| `OracleSetRequest` | [QueryOracleSetRequestRequest](#fx.gravity.crosschain.v1.QueryOracleSetRequestRequest) | [QueryOracleSetRequestResponse](#fx.gravity.crosschain.v1.QueryOracleSetRequestResponse) |  | GET|/fx/crosschain/v1/oracle_set/request|
| `OracleSetConfirm` | [QueryOracleSetConfirmRequest](#fx.gravity.crosschain.v1.QueryOracleSetConfirmRequest) | [QueryOracleSetConfirmResponse](#fx.gravity.crosschain.v1.QueryOracleSetConfirmResponse) |  | GET|/fx/crosschain/v1/oracle_set/confirm|
| `OracleSetConfirmsByNonce` | [QueryOracleSetConfirmsByNonceRequest](#fx.gravity.crosschain.v1.QueryOracleSetConfirmsByNonceRequest) | [QueryOracleSetConfirmsByNonceResponse](#fx.gravity.crosschain.v1.QueryOracleSetConfirmsByNonceResponse) |  | GET|/fx/crosschain/v1/oracle_set/confirms|
| `LastOracleSetRequests` | [QueryLastOracleSetRequestsRequest](#fx.gravity.crosschain.v1.QueryLastOracleSetRequestsRequest) | [QueryLastOracleSetRequestsResponse](#fx.gravity.crosschain.v1.QueryLastOracleSetRequestsResponse) |  | GET|/fx/crosschain/v1/oracle_set/requests|
| `LastPendingOracleSetRequestByAddr` | [QueryLastPendingOracleSetRequestByAddrRequest](#fx.gravity.crosschain.v1.QueryLastPendingOracleSetRequestByAddrRequest) | [QueryLastPendingOracleSetRequestByAddrResponse](#fx.gravity.crosschain.v1.QueryLastPendingOracleSetRequestByAddrResponse) |  | GET|/fx/crosschain/v1/oracle_set/last|
| `LastPendingBatchRequestByAddr` | [QueryLastPendingBatchRequestByAddrRequest](#fx.gravity.crosschain.v1.QueryLastPendingBatchRequestByAddrRequest) | [QueryLastPendingBatchRequestByAddrResponse](#fx.gravity.crosschain.v1.QueryLastPendingBatchRequestByAddrResponse) |  | GET|/fx/crosschain/v1/batch/last|
| `LastEventNonceByAddr` | [QueryLastEventNonceByAddrRequest](#fx.gravity.crosschain.v1.QueryLastEventNonceByAddrRequest) | [QueryLastEventNonceByAddrResponse](#fx.gravity.crosschain.v1.QueryLastEventNonceByAddrResponse) |  | GET|/fx/crosschain/v1/oracle/event_nonce|
| `LastEventBlockHeightByAddr` | [QueryLastEventBlockHeightByAddrRequest](#fx.gravity.crosschain.v1.QueryLastEventBlockHeightByAddrRequest) | [QueryLastEventBlockHeightByAddrResponse](#fx.gravity.crosschain.v1.QueryLastEventBlockHeightByAddrResponse) |  | GET|/fx/crosschain/v1/oracle/event/block_height|
| `BatchFees` | [QueryBatchFeeRequest](#fx.gravity.crosschain.v1.QueryBatchFeeRequest) | [QueryBatchFeeResponse](#fx.gravity.crosschain.v1.QueryBatchFeeResponse) |  | GET|/fx/crosschain/v1/batch_fees|
| `LastObservedBlockHeight` | [QueryLastObservedBlockHeightRequest](#fx.gravity.crosschain.v1.QueryLastObservedBlockHeightRequest) | [QueryLastObservedBlockHeightResponse](#fx.gravity.crosschain.v1.QueryLastObservedBlockHeightResponse) |  | GET|/fx/crosschain/v1/observed/block_height|
| `OutgoingTxBatches` | [QueryOutgoingTxBatchesRequest](#fx.gravity.crosschain.v1.QueryOutgoingTxBatchesRequest) | [QueryOutgoingTxBatchesResponse](#fx.gravity.crosschain.v1.QueryOutgoingTxBatchesResponse) |  | GET|/fx/crosschain/v1/batch/outgoing_tx|
| `BatchRequestByNonce` | [QueryBatchRequestByNonceRequest](#fx.gravity.crosschain.v1.QueryBatchRequestByNonceRequest) | [QueryBatchRequestByNonceResponse](#fx.gravity.crosschain.v1.QueryBatchRequestByNonceResponse) |  | GET|/fx/crosschain/v1/batch/request|
| `BatchConfirm` | [QueryBatchConfirmRequest](#fx.gravity.crosschain.v1.QueryBatchConfirmRequest) | [QueryBatchConfirmResponse](#fx.gravity.crosschain.v1.QueryBatchConfirmResponse) |  | GET|/fx/crosschain/v1/batch/confirm|
| `BatchConfirms` | [QueryBatchConfirmsRequest](#fx.gravity.crosschain.v1.QueryBatchConfirmsRequest) | [QueryBatchConfirmsResponse](#fx.gravity.crosschain.v1.QueryBatchConfirmsResponse) |  | GET|/fx/crosschain/v1/batch/confirms|
| `TokenToDenom` | [QueryTokenToDenomRequest](#fx.gravity.crosschain.v1.QueryTokenToDenomRequest) | [QueryTokenToDenomResponse](#fx.gravity.crosschain.v1.QueryTokenToDenomResponse) |  | GET|/fx/crosschain/v1/denom|
| `DenomToToken` | [QueryDenomToTokenRequest](#fx.gravity.crosschain.v1.QueryDenomToTokenRequest) | [QueryDenomToTokenResponse](#fx.gravity.crosschain.v1.QueryDenomToTokenResponse) |  | GET|/fx/crosschain/v1/token|
| `GetOracleByAddr` | [QueryOracleByAddrRequest](#fx.gravity.crosschain.v1.QueryOracleByAddrRequest) | [QueryOracleResponse](#fx.gravity.crosschain.v1.QueryOracleResponse) |  | GET|/fx/crosschain/v1/oracle_by_addr|
| `GetOracleByExternalAddr` | [QueryOracleByExternalAddrRequest](#fx.gravity.crosschain.v1.QueryOracleByExternalAddrRequest) | [QueryOracleResponse](#fx.gravity.crosschain.v1.QueryOracleResponse) |  | GET|/fx/crosschain/v1/oracle_by_external_addr|
| `GetOracleByBridgerAddr` | [QueryOracleByBridgerAddrRequest](#fx.gravity.crosschain.v1.QueryOracleByBridgerAddrRequest) | [QueryOracleResponse](#fx.gravity.crosschain.v1.QueryOracleResponse) |  | GET|/fx/crosschain/v1/oracle_by_bridger_addr|
| `GetPendingSendToExternal` | [QueryPendingSendToExternalRequest](#fx.gravity.crosschain.v1.QueryPendingSendToExternalRequest) | [QueryPendingSendToExternalResponse](#fx.gravity.crosschain.v1.QueryPendingSendToExternalResponse) |  | GET|/fx/crosschain/v1/pending_send_to_external|
| `Oracles` | [QueryOraclesRequest](#fx.gravity.crosschain.v1.QueryOraclesRequest) | [QueryOraclesResponse](#fx.gravity.crosschain.v1.QueryOraclesResponse) | Validators queries all oracle that match the given status. | GET|/fx/crosschain/v1/oracles|
| `ProjectedBatchTimeoutHeight` | [QueryProjectedBatchTimeoutHeightRequest](#fx.gravity.crosschain.v1.QueryProjectedBatchTimeoutHeightRequest) | [QueryProjectedBatchTimeoutHeightResponse](#fx.gravity.crosschain.v1.QueryProjectedBatchTimeoutHeightResponse) |  | GET|/fx/crosschain/v1/projected_batch_timeout|
| `BridgeTokens` | [QueryBridgeTokensRequest](#fx.gravity.crosschain.v1.QueryBridgeTokensRequest) | [QueryBridgeTokensResponse](#fx.gravity.crosschain.v1.QueryBridgeTokensResponse) |  | GET|/fx/crosschain/v1/bridge_tokens|

 <!-- end services -->



<a name="fx/erc20/v1/erc20.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/erc20/v1/erc20.proto



<a name="fx.erc20.v1.RegisterCoinProposal"></a>

### RegisterCoinProposal
RegisterCoinProposal is a gov Content type to register a token pair


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | title of the proposal |
| `description` | [string](#string) |  | proposal description |
| `metadata` | [cosmos.bank.v1beta1.Metadata](#cosmos.bank.v1beta1.Metadata) |  | token pair of Cosmos native denom and ERC20 token address |






<a name="fx.erc20.v1.RegisterERC20Proposal"></a>

### RegisterERC20Proposal
RegisterCoinProposal is a gov Content type to register a token pair


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | title of the proposal |
| `description` | [string](#string) |  | proposal description |
| `erc20address` | [string](#string) |  | contract address of ERC20 token |
| `aliases` | [string](#string) | repeated | aliases is a list of string aliases for the given denom |






<a name="fx.erc20.v1.ToggleTokenConversionProposal"></a>

### ToggleTokenConversionProposal
ToggleTokenConversionProposal is a gov Content type to toggle the conversion
of a token pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | title of the proposal |
| `description` | [string](#string) |  | proposal description |
| `token` | [string](#string) |  | token identifier can be either the hex contract address of the ERC20 or the Cosmos base denomination |






<a name="fx.erc20.v1.TokenPair"></a>

### TokenPair
TokenPair defines an instance that records pairing consisting of a Cosmos
native Coin and an ERC20 token address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `erc20_address` | [string](#string) |  | address of ERC20 contract token |
| `denom` | [string](#string) |  | cosmos base denomination to be mapped to |
| `enabled` | [bool](#bool) |  | shows token mapping enable status |
| `contract_owner` | [Owner](#fx.erc20.v1.Owner) |  | ERC20 owner address ENUM (0 invalid, 1 ModuleAccount, 2 external address) |






<a name="fx.erc20.v1.UpdateDenomAliasProposal"></a>

### UpdateDenomAliasProposal
UpdateDenomAliasProposal is a gov Content type to update denom alias


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | title of the proposal |
| `description` | [string](#string) |  | proposal description |
| `denom` | [string](#string) |  | denom |
| `alias` | [string](#string) |  | alias update |





 <!-- end messages -->


<a name="fx.erc20.v1.Owner"></a>

### Owner
Owner enumerates the ownership of a ERC20 contract.

| Name | Number | Description |
| ---- | ------ | ----------- |
| OWNER_UNSPECIFIED | 0 | OWNER_UNSPECIFIED defines an invalid/undefined owner. |
| OWNER_MODULE | 1 | OWNER_MODULE erc20 is owned by the erc20 module account. |
| OWNER_EXTERNAL | 2 | EXTERNAL erc20 is owned by an external account. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="fx/erc20/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/erc20/v1/genesis.proto



<a name="fx.erc20.v1.GenesisState"></a>

### GenesisState
GenesisState defines the module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#fx.erc20.v1.Params) |  | module parameters |
| `token_pairs` | [TokenPair](#fx.erc20.v1.TokenPair) | repeated | registered token pairs |






<a name="fx.erc20.v1.Params"></a>

### Params
Params defines the erc20 module params


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enable_erc20` | [bool](#bool) |  | parameter to enable the intrarelaying of Cosmos coins <--> ERC20 tokens. |
| `enable_evm_hook` | [bool](#bool) |  | parameter to enable the EVM hook to convert an ERC20 token to a Cosmos Coin by transferring the Tokens through a MsgEthereumTx to the ModuleAddress Ethereum address. |
| `ibc_timeout` | [google.protobuf.Duration](#google.protobuf.Duration) |  | parameter to set ibc timeout |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="fx/erc20/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/erc20/v1/query.proto



<a name="fx.erc20.v1.QueryAliasDenomRequest"></a>

### QueryAliasDenomRequest
QueryAliasDenomRequest is the request type for the Query/AliasDenom RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `alias` | [string](#string) |  |  |






<a name="fx.erc20.v1.QueryAliasDenomResponse"></a>

### QueryAliasDenomResponse
QueryAliasDenomResponse is the response type for the Query/AliasDenom RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |






<a name="fx.erc20.v1.QueryDenomAliasesRequest"></a>

### QueryDenomAliasesRequest
QueryDenomAliasesRequest is the request type for the Query/DenomAliases RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |






<a name="fx.erc20.v1.QueryDenomAliasesResponse"></a>

### QueryDenomAliasesResponse
QueryDenomAliasesResponse is the response type for the Query/DenomAliases RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `aliases` | [string](#string) | repeated |  |






<a name="fx.erc20.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="fx.erc20.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#fx.erc20.v1.Params) |  |  |






<a name="fx.erc20.v1.QueryTokenPairRequest"></a>

### QueryTokenPairRequest
QueryTokenPairRequest is the request type for the Query/TokenPair RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token` | [string](#string) |  | token identifier can be either the hex contract address of the ERC20 or the Cosmos base denomination |






<a name="fx.erc20.v1.QueryTokenPairResponse"></a>

### QueryTokenPairResponse
QueryTokenPairResponse is the response type for the Query/TokenPair RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_pair` | [TokenPair](#fx.erc20.v1.TokenPair) |  |  |






<a name="fx.erc20.v1.QueryTokenPairsRequest"></a>

### QueryTokenPairsRequest
QueryTokenPairsRequest is the request type for the Query/TokenPairs RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="fx.erc20.v1.QueryTokenPairsResponse"></a>

### QueryTokenPairsResponse
QueryTokenPairsResponse is the response type for the Query/TokenPairs RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_pairs` | [TokenPair](#fx.erc20.v1.TokenPair) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.erc20.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `TokenPairs` | [QueryTokenPairsRequest](#fx.erc20.v1.QueryTokenPairsRequest) | [QueryTokenPairsResponse](#fx.erc20.v1.QueryTokenPairsResponse) | Retrieves registered token pairs | GET|/fx/erc20/v1/token_pairs|
| `TokenPair` | [QueryTokenPairRequest](#fx.erc20.v1.QueryTokenPairRequest) | [QueryTokenPairResponse](#fx.erc20.v1.QueryTokenPairResponse) | Retrieves a registered token pair | GET|/fx/erc20/v1/token_pairs/{token}|
| `Params` | [QueryParamsRequest](#fx.erc20.v1.QueryParamsRequest) | [QueryParamsResponse](#fx.erc20.v1.QueryParamsResponse) | Params retrieves the erc20 module params | GET|/fx/erc20/v1/params|
| `DenomAliases` | [QueryDenomAliasesRequest](#fx.erc20.v1.QueryDenomAliasesRequest) | [QueryDenomAliasesResponse](#fx.erc20.v1.QueryDenomAliasesResponse) | Retrieves registered denom aliases | GET|/fx/erc20/v1/denom/aliases|
| `AliasDenom` | [QueryAliasDenomRequest](#fx.erc20.v1.QueryAliasDenomRequest) | [QueryAliasDenomResponse](#fx.erc20.v1.QueryAliasDenomResponse) | Retrieves registered alias denom | GET|/fx/erc20/v1/alias/denom|

 <!-- end services -->



<a name="fx/erc20/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/erc20/v1/tx.proto



<a name="fx.erc20.v1.MsgConvertCoin"></a>

### MsgConvertCoin
MsgConvertCoin defines a Msg to convert a Cosmos Coin to a ERC20 token


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | Cosmos coin which denomination is registered on erc20 bridge. The coin amount defines the total ERC20 tokens to convert. |
| `receiver` | [string](#string) |  | recipient hex address to receive ERC20 token |
| `sender` | [string](#string) |  | cosmos bech32 address from the owner of the given ERC20 tokens |






<a name="fx.erc20.v1.MsgConvertCoinResponse"></a>

### MsgConvertCoinResponse
MsgConvertCoinResponse returns no fields






<a name="fx.erc20.v1.MsgConvertDenom"></a>

### MsgConvertDenom
MsgConvertDenom defines a Msg to convert an denom to other denom


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | msg sender |
| `receiver` | [string](#string) |  | receiver address |
| `coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | coin to convert |
| `target` | [string](#string) |  | target denom |






<a name="fx.erc20.v1.MsgConvertDenomResponse"></a>

### MsgConvertDenomResponse
MsgConvertDenomResponse returns no fields






<a name="fx.erc20.v1.MsgConvertERC20"></a>

### MsgConvertERC20
MsgConvertERC20 defines a Msg to convert an ERC20 token to a Cosmos SDK coin.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  | ERC20 token contract address registered on erc20 bridge |
| `amount` | [string](#string) |  | amount of ERC20 tokens to mint |
| `receiver` | [string](#string) |  | bech32 address to receive SDK coins. |
| `sender` | [string](#string) |  | sender hex address from the owner of the given ERC20 tokens |






<a name="fx.erc20.v1.MsgConvertERC20Response"></a>

### MsgConvertERC20Response
MsgConvertERC20Response returns no fields






<a name="fx.erc20.v1.MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams is the Msg/UpdateParams request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address of the governance account. |
| `params` | [Params](#fx.erc20.v1.Params) |  | params defines the x/erc20 parameters to update.

NOTE: All parameters must be supplied. |






<a name="fx.erc20.v1.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse
MsgUpdateParamsResponse defines the response structure for executing a
MsgUpdateParams message.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.erc20.v1.Msg"></a>

### Msg
Msg defines the erc20 Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ConvertCoin` | [MsgConvertCoin](#fx.erc20.v1.MsgConvertCoin) | [MsgConvertCoinResponse](#fx.erc20.v1.MsgConvertCoinResponse) | ConvertCoin mints a ERC20 representation of the SDK Coin denom that is registered on the token mapping. | |
| `ConvertERC20` | [MsgConvertERC20](#fx.erc20.v1.MsgConvertERC20) | [MsgConvertERC20Response](#fx.erc20.v1.MsgConvertERC20Response) | ConvertERC20 mints a Cosmos coin representation of the ERC20 token contract that is registered on the token mapping. | |
| `ConvertDenom` | [MsgConvertDenom](#fx.erc20.v1.MsgConvertDenom) | [MsgConvertDenomResponse](#fx.erc20.v1.MsgConvertDenomResponse) | ConvertDenom convert denom to other denom | |
| `UpdateParams` | [MsgUpdateParams](#fx.erc20.v1.MsgUpdateParams) | [MsgUpdateParamsResponse](#fx.erc20.v1.MsgUpdateParamsResponse) | UpdateParams defines a governance operation for updating the x/erc20 module parameters. The authority is hard-coded to the x/gov module account. | |

 <!-- end services -->



<a name="fx/ibc/applications/transfer/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/ibc/applications/transfer/v1/query.proto


 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.ibc.applications.transfer.v1.Query"></a>

### Query
Query provides defines the gRPC querier service.
Deprecated: This service is deprecated. It may be removed in the next
version. Replace ibc.applications.transfer.v1.Query

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `DenomTrace` | [.ibc.applications.transfer.v1.QueryDenomTraceRequest](#ibc.applications.transfer.v1.QueryDenomTraceRequest) | [.ibc.applications.transfer.v1.QueryDenomTraceResponse](#ibc.applications.transfer.v1.QueryDenomTraceResponse) | DenomTrace queries a denomination trace information. | |
| `DenomTraces` | [.ibc.applications.transfer.v1.QueryDenomTracesRequest](#ibc.applications.transfer.v1.QueryDenomTracesRequest) | [.ibc.applications.transfer.v1.QueryDenomTracesResponse](#ibc.applications.transfer.v1.QueryDenomTracesResponse) | DenomTraces queries all denomination traces. | |
| `Params` | [.ibc.applications.transfer.v1.QueryParamsRequest](#ibc.applications.transfer.v1.QueryParamsRequest) | [.ibc.applications.transfer.v1.QueryParamsResponse](#ibc.applications.transfer.v1.QueryParamsResponse) | Params queries all parameters of the ibc-transfer module. | |
| `DenomHash` | [.ibc.applications.transfer.v1.QueryDenomHashRequest](#ibc.applications.transfer.v1.QueryDenomHashRequest) | [.ibc.applications.transfer.v1.QueryDenomHashResponse](#ibc.applications.transfer.v1.QueryDenomHashResponse) | DenomHash queries a denomination hash information. | |
| `EscrowAddress` | [.ibc.applications.transfer.v1.QueryEscrowAddressRequest](#ibc.applications.transfer.v1.QueryEscrowAddressRequest) | [.ibc.applications.transfer.v1.QueryEscrowAddressResponse](#ibc.applications.transfer.v1.QueryEscrowAddressResponse) | EscrowAddress returns the escrow address for a particular port and channel id. | |

 <!-- end services -->



<a name="fx/ibc/applications/transfer/v1/transfer.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/ibc/applications/transfer/v1/transfer.proto



<a name="fx.ibc.applications.transfer.v1.FungibleTokenPacketData"></a>

### FungibleTokenPacketData
FungibleTokenPacketData defines a struct for the packet payload
See FungibleTokenPacketData spec:
https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer#data-structures


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  | the token denomination to be transferred |
| `amount` | [string](#string) |  | the token amount to be transferred |
| `sender` | [string](#string) |  | the sender address |
| `receiver` | [string](#string) |  | the recipient address on the destination chain |
| `router` | [string](#string) |  | the router is hook destination chain |
| `fee` | [string](#string) |  | the fee is destination fee |
| `memo` | [string](#string) |  | optional memo |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="fx/ibc/applications/transfer/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/ibc/applications/transfer/v1/tx.proto



<a name="fx.ibc.applications.transfer.v1.MsgTransfer"></a>

### MsgTransfer
MsgTransfer defines a msg to transfer fungible tokens (i.e Coins) between
ICS20 enabled chains. See ICS Spec here:
https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer#data-structures


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `source_port` | [string](#string) |  | the port on which the packet will be sent |
| `source_channel` | [string](#string) |  | the channel by which the packet will be sent |
| `token` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | the tokens to be transferred |
| `sender` | [string](#string) |  | the sender address |
| `receiver` | [string](#string) |  | the recipient address on the destination chain |
| `timeout_height` | [ibc.core.client.v1.Height](#ibc.core.client.v1.Height) |  | Timeout height relative to the current block height. The timeout is disabled when set to 0. |
| `timeout_timestamp` | [uint64](#uint64) |  | Timeout timestamp (in nanoseconds) relative to the current block timestamp. The timeout is disabled when set to 0. |
| `router` | [string](#string) |  | the router is hook destination chain |
| `fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | the tokens to be destination fee |
| `memo` | [string](#string) |  | optional memo |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.ibc.applications.transfer.v1.Msg"></a>

### Msg
Msg defines the ibc/transfer Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Transfer` | [MsgTransfer](#fx.ibc.applications.transfer.v1.MsgTransfer) | [.ibc.applications.transfer.v1.MsgTransferResponse](#ibc.applications.transfer.v1.MsgTransferResponse) | Transfer defines a rpc handler method for MsgTransfer. | |

 <!-- end services -->



<a name="fx/legacy/gravity/v1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/legacy/gravity/v1/types.proto



<a name="fx.gravity.v1.Attestation"></a>

### Attestation
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `observed` | [bool](#bool) |  |  |
| `votes` | [string](#string) | repeated |  |
| `height` | [uint64](#uint64) |  |  |
| `claim` | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="fx.gravity.v1.BatchFees"></a>

### BatchFees
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_contract` | [string](#string) |  |  |
| `total_fees` | [string](#string) |  |  |
| `total_txs` | [uint64](#uint64) |  |  |
| `total_amount` | [string](#string) |  |  |






<a name="fx.gravity.v1.BridgeValidator"></a>

### BridgeValidator
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `power` | [uint64](#uint64) |  |  |
| `eth_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.ERC20ToDenom"></a>

### ERC20ToDenom
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `erc20` | [string](#string) |  |  |
| `denom` | [string](#string) |  |  |






<a name="fx.gravity.v1.ERC20Token"></a>

### ERC20Token
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract` | [string](#string) |  |  |
| `amount` | [string](#string) |  |  |






<a name="fx.gravity.v1.LastObservedEthereumBlockHeight"></a>

### LastObservedEthereumBlockHeight
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fx_block_height` | [uint64](#uint64) |  |  |
| `eth_block_height` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.MinBatchFee"></a>

### MinBatchFee
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_contract` | [string](#string) |  |  |
| `baseFee` | [string](#string) |  |  |






<a name="fx.gravity.v1.OutgoingTransferTx"></a>

### OutgoingTransferTx
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `sender` | [string](#string) |  |  |
| `dest_address` | [string](#string) |  |  |
| `erc20_token` | [ERC20Token](#fx.gravity.v1.ERC20Token) |  |  |
| `erc20_fee` | [ERC20Token](#fx.gravity.v1.ERC20Token) |  |  |






<a name="fx.gravity.v1.OutgoingTxBatch"></a>

### OutgoingTxBatch
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch_nonce` | [uint64](#uint64) |  |  |
| `batch_timeout` | [uint64](#uint64) |  |  |
| `transactions` | [OutgoingTransferTx](#fx.gravity.v1.OutgoingTransferTx) | repeated |  |
| `token_contract` | [string](#string) |  |  |
| `block` | [uint64](#uint64) |  |  |
| `feeReceive` | [string](#string) |  |  |






<a name="fx.gravity.v1.Valset"></a>

### Valset
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `members` | [BridgeValidator](#fx.gravity.v1.BridgeValidator) | repeated |  |
| `height` | [uint64](#uint64) |  |  |





 <!-- end messages -->


<a name="fx.gravity.v1.ClaimType"></a>

### ClaimType
Deprecated: after upgrade v3

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



<a name="fx/legacy/gravity/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/legacy/gravity/v1/tx.proto



<a name="fx.gravity.v1.MsgCancelSendToEth"></a>

### MsgCancelSendToEth
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `transaction_id` | [uint64](#uint64) |  |  |
| `sender` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgCancelSendToEthResponse"></a>

### MsgCancelSendToEthResponse
Deprecated: after upgrade v3






<a name="fx.gravity.v1.MsgConfirmBatch"></a>

### MsgConfirmBatch
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
| `eth_signer` | [string](#string) |  |  |
| `orchestrator` | [string](#string) |  |  |
| `signature` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgConfirmBatchResponse"></a>

### MsgConfirmBatchResponse
Deprecated: after upgrade v3






<a name="fx.gravity.v1.MsgDepositClaim"></a>

### MsgDepositClaim
Deprecated: after upgrade v3


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
Deprecated: after upgrade v3






<a name="fx.gravity.v1.MsgFxOriginatedTokenClaim"></a>

### MsgFxOriginatedTokenClaim
Deprecated: after upgrade v3


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
Deprecated: after upgrade v3






<a name="fx.gravity.v1.MsgRequestBatch"></a>

### MsgRequestBatch
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `denom` | [string](#string) |  |  |
| `minimum_fee` | [string](#string) |  |  |
| `feeReceive` | [string](#string) |  |  |
| `base_fee` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgRequestBatchResponse"></a>

### MsgRequestBatchResponse
Deprecated: after upgrade v3






<a name="fx.gravity.v1.MsgSendToEth"></a>

### MsgSendToEth
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  |  |
| `eth_dest` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |
| `bridge_fee` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  |  |






<a name="fx.gravity.v1.MsgSendToEthResponse"></a>

### MsgSendToEthResponse
Deprecated: after upgrade v3






<a name="fx.gravity.v1.MsgSetOrchestratorAddress"></a>

### MsgSetOrchestratorAddress
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator` | [string](#string) |  |  |
| `orchestrator` | [string](#string) |  |  |
| `eth_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgSetOrchestratorAddressResponse"></a>

### MsgSetOrchestratorAddressResponse
Deprecated: after upgrade v3






<a name="fx.gravity.v1.MsgValsetConfirm"></a>

### MsgValsetConfirm
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `orchestrator` | [string](#string) |  |  |
| `eth_address` | [string](#string) |  |  |
| `signature` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgValsetConfirmResponse"></a>

### MsgValsetConfirmResponse
Deprecated: after upgrade v3






<a name="fx.gravity.v1.MsgValsetUpdatedClaim"></a>

### MsgValsetUpdatedClaim
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `event_nonce` | [uint64](#uint64) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `valset_nonce` | [uint64](#uint64) |  |  |
| `members` | [BridgeValidator](#fx.gravity.v1.BridgeValidator) | repeated |  |
| `orchestrator` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgValsetUpdatedClaimResponse"></a>

### MsgValsetUpdatedClaimResponse
Deprecated: after upgrade v3






<a name="fx.gravity.v1.MsgWithdrawClaim"></a>

### MsgWithdrawClaim
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `event_nonce` | [uint64](#uint64) |  |  |
| `block_height` | [uint64](#uint64) |  |  |
| `batch_nonce` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
| `orchestrator` | [string](#string) |  |  |






<a name="fx.gravity.v1.MsgWithdrawClaimResponse"></a>

### MsgWithdrawClaimResponse
Deprecated: after upgrade v3





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.gravity.v1.Msg"></a>

### Msg
Deprecated: after upgrade v3

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ValsetConfirm` | [MsgValsetConfirm](#fx.gravity.v1.MsgValsetConfirm) | [MsgValsetConfirmResponse](#fx.gravity.v1.MsgValsetConfirmResponse) | Deprecated: Please use crosschain Msg.OracleSetConfirm | |
| `SendToEth` | [MsgSendToEth](#fx.gravity.v1.MsgSendToEth) | [MsgSendToEthResponse](#fx.gravity.v1.MsgSendToEthResponse) | Deprecated: Please use crosschain Msg.SendToExternal | |
| `RequestBatch` | [MsgRequestBatch](#fx.gravity.v1.MsgRequestBatch) | [MsgRequestBatchResponse](#fx.gravity.v1.MsgRequestBatchResponse) | Deprecated: Please use crosschain Msg.RequestBatch | |
| `ConfirmBatch` | [MsgConfirmBatch](#fx.gravity.v1.MsgConfirmBatch) | [MsgConfirmBatchResponse](#fx.gravity.v1.MsgConfirmBatchResponse) | Deprecated: Please use crosschain Msg.ConfirmBatch | |
| `DepositClaim` | [MsgDepositClaim](#fx.gravity.v1.MsgDepositClaim) | [MsgDepositClaimResponse](#fx.gravity.v1.MsgDepositClaimResponse) | Deprecated: Please use crosschain Msg.SendToFxClaim | |
| `WithdrawClaim` | [MsgWithdrawClaim](#fx.gravity.v1.MsgWithdrawClaim) | [MsgWithdrawClaimResponse](#fx.gravity.v1.MsgWithdrawClaimResponse) | Deprecated: Please use crosschain Msg.SendToExternalClaim | |
| `CancelSendToEth` | [MsgCancelSendToEth](#fx.gravity.v1.MsgCancelSendToEth) | [MsgCancelSendToEthResponse](#fx.gravity.v1.MsgCancelSendToEthResponse) | Deprecated: Please use crosschain Msg.CancelSendToExternal | |
| `ValsetUpdateClaim` | [MsgValsetUpdatedClaim](#fx.gravity.v1.MsgValsetUpdatedClaim) | [MsgValsetUpdatedClaimResponse](#fx.gravity.v1.MsgValsetUpdatedClaimResponse) | Deprecated: Please use crosschain Msg.OracleSetUpdateClaim | |

 <!-- end services -->



<a name="fx/legacy/gravity/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/legacy/gravity/v1/genesis.proto



<a name="fx.gravity.v1.GenesisState"></a>

### GenesisState
Deprecated: after upgrade v3


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
| `last_observed_valset` | [Valset](#fx.gravity.v1.Valset) |  |  |
| `last_slashed_batch_block` | [uint64](#uint64) |  |  |
| `last_slashed_valset_nonce` | [uint64](#uint64) |  |  |
| `last_tx_pool_id` | [uint64](#uint64) |  |  |
| `last_batch_id` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.Params"></a>

### Params
Deprecated: after upgrade v3


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



<a name="fx/legacy/gravity/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/legacy/gravity/v1/query.proto



<a name="fx.gravity.v1.QueryBatchConfirmRequest"></a>

### QueryBatchConfirmRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
| `address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryBatchConfirmResponse"></a>

### QueryBatchConfirmResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirm` | [MsgConfirmBatch](#fx.gravity.v1.MsgConfirmBatch) |  |  |






<a name="fx.gravity.v1.QueryBatchConfirmsRequest"></a>

### QueryBatchConfirmsRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryBatchConfirmsResponse"></a>

### QueryBatchConfirmsResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirms` | [MsgConfirmBatch](#fx.gravity.v1.MsgConfirmBatch) | repeated |  |






<a name="fx.gravity.v1.QueryBatchFeeRequest"></a>

### QueryBatchFeeRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `minBatchFees` | [MinBatchFee](#fx.gravity.v1.MinBatchFee) | repeated |  |






<a name="fx.gravity.v1.QueryBatchFeeResponse"></a>

### QueryBatchFeeResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch_fees` | [BatchFees](#fx.gravity.v1.BatchFees) | repeated |  |






<a name="fx.gravity.v1.QueryBatchRequestByNonceRequest"></a>

### QueryBatchRequestByNonceRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryBatchRequestByNonceResponse"></a>

### QueryBatchRequestByNonceResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch` | [OutgoingTxBatch](#fx.gravity.v1.OutgoingTxBatch) |  |  |






<a name="fx.gravity.v1.QueryBridgeTokensRequest"></a>

### QueryBridgeTokensRequest
Deprecated: after upgrade v3






<a name="fx.gravity.v1.QueryBridgeTokensResponse"></a>

### QueryBridgeTokensResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bridge_tokens` | [ERC20ToDenom](#fx.gravity.v1.ERC20ToDenom) | repeated |  |






<a name="fx.gravity.v1.QueryCurrentValsetRequest"></a>

### QueryCurrentValsetRequest
Deprecated: after upgrade v3






<a name="fx.gravity.v1.QueryCurrentValsetResponse"></a>

### QueryCurrentValsetResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `valset` | [Valset](#fx.gravity.v1.Valset) |  |  |






<a name="fx.gravity.v1.QueryDelegateKeyByEthRequest"></a>

### QueryDelegateKeyByEthRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `eth_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryDelegateKeyByEthResponse"></a>

### QueryDelegateKeyByEthResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |
| `orchestrator_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryDelegateKeyByOrchestratorRequest"></a>

### QueryDelegateKeyByOrchestratorRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `orchestrator_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryDelegateKeyByOrchestratorResponse"></a>

### QueryDelegateKeyByOrchestratorResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |
| `eth_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryDelegateKeyByValidatorRequest"></a>

### QueryDelegateKeyByValidatorRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `validator_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryDelegateKeyByValidatorResponse"></a>

### QueryDelegateKeyByValidatorResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `eth_address` | [string](#string) |  |  |
| `orchestrator_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryDenomToERC20Request"></a>

### QueryDenomToERC20Request
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryDenomToERC20Response"></a>

### QueryDenomToERC20Response
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `erc20` | [string](#string) |  |  |
| `fx_originated` | [bool](#bool) |  |  |






<a name="fx.gravity.v1.QueryERC20ToDenomRequest"></a>

### QueryERC20ToDenomRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `erc20` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryERC20ToDenomResponse"></a>

### QueryERC20ToDenomResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom` | [string](#string) |  |  |
| `fx_originated` | [bool](#bool) |  |  |






<a name="fx.gravity.v1.QueryLastEventBlockHeightByAddrRequest"></a>

### QueryLastEventBlockHeightByAddrRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryLastEventBlockHeightByAddrResponse"></a>

### QueryLastEventBlockHeightByAddrResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_height` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.QueryLastEventNonceByAddrRequest"></a>

### QueryLastEventNonceByAddrRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryLastEventNonceByAddrResponse"></a>

### QueryLastEventNonceByAddrResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `event_nonce` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.QueryLastObservedBlockHeightRequest"></a>

### QueryLastObservedBlockHeightRequest
Deprecated: after upgrade v3






<a name="fx.gravity.v1.QueryLastObservedBlockHeightResponse"></a>

### QueryLastObservedBlockHeightResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `eth_block_height` | [uint64](#uint64) |  |  |
| `block_height` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.QueryLastPendingBatchRequestByAddrRequest"></a>

### QueryLastPendingBatchRequestByAddrRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryLastPendingBatchRequestByAddrResponse"></a>

### QueryLastPendingBatchRequestByAddrResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch` | [OutgoingTxBatch](#fx.gravity.v1.OutgoingTxBatch) |  |  |






<a name="fx.gravity.v1.QueryLastPendingValsetRequestByAddrRequest"></a>

### QueryLastPendingValsetRequestByAddrRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryLastPendingValsetRequestByAddrResponse"></a>

### QueryLastPendingValsetRequestByAddrResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `valsets` | [Valset](#fx.gravity.v1.Valset) | repeated |  |






<a name="fx.gravity.v1.QueryLastValsetRequestsRequest"></a>

### QueryLastValsetRequestsRequest
Deprecated: after upgrade v3






<a name="fx.gravity.v1.QueryLastValsetRequestsResponse"></a>

### QueryLastValsetRequestsResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `valsets` | [Valset](#fx.gravity.v1.Valset) | repeated |  |






<a name="fx.gravity.v1.QueryOutgoingTxBatchesRequest"></a>

### QueryOutgoingTxBatchesRequest
Deprecated: after upgrade v3






<a name="fx.gravity.v1.QueryOutgoingTxBatchesResponse"></a>

### QueryOutgoingTxBatchesResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batches` | [OutgoingTxBatch](#fx.gravity.v1.OutgoingTxBatch) | repeated |  |






<a name="fx.gravity.v1.QueryParamsRequest"></a>

### QueryParamsRequest
Deprecated: after upgrade v3






<a name="fx.gravity.v1.QueryParamsResponse"></a>

### QueryParamsResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#fx.gravity.v1.Params) |  |  |






<a name="fx.gravity.v1.QueryPendingSendToEthRequest"></a>

### QueryPendingSendToEthRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender_address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryPendingSendToEthResponse"></a>

### QueryPendingSendToEthResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `transfers_in_batches` | [OutgoingTransferTx](#fx.gravity.v1.OutgoingTransferTx) | repeated |  |
| `unbatched_transfers` | [OutgoingTransferTx](#fx.gravity.v1.OutgoingTransferTx) | repeated |  |






<a name="fx.gravity.v1.QueryProjectedBatchTimeoutHeightRequest"></a>

### QueryProjectedBatchTimeoutHeightRequest
Deprecated: after upgrade v3






<a name="fx.gravity.v1.QueryProjectedBatchTimeoutHeightResponse"></a>

### QueryProjectedBatchTimeoutHeightResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `timeout_height` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.QueryValsetConfirmRequest"></a>

### QueryValsetConfirmRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `address` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryValsetConfirmResponse"></a>

### QueryValsetConfirmResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirm` | [MsgValsetConfirm](#fx.gravity.v1.MsgValsetConfirm) |  |  |






<a name="fx.gravity.v1.QueryValsetConfirmsByNonceRequest"></a>

### QueryValsetConfirmsByNonceRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.QueryValsetConfirmsByNonceResponse"></a>

### QueryValsetConfirmsByNonceResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirms` | [MsgValsetConfirm](#fx.gravity.v1.MsgValsetConfirm) | repeated |  |






<a name="fx.gravity.v1.QueryValsetRequestRequest"></a>

### QueryValsetRequestRequest
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |






<a name="fx.gravity.v1.QueryValsetRequestResponse"></a>

### QueryValsetRequestResponse
Deprecated: after upgrade v3


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `valset` | [Valset](#fx.gravity.v1.Valset) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.gravity.v1.Query"></a>

### Query
Deprecated: after upgrade v3

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#fx.gravity.v1.QueryParamsRequest) | [QueryParamsResponse](#fx.gravity.v1.QueryParamsResponse) | Deprecated: Please use crosschain query.Params | GET|/fx/gravity/v1/params|
| `CurrentValset` | [QueryCurrentValsetRequest](#fx.gravity.v1.QueryCurrentValsetRequest) | [QueryCurrentValsetResponse](#fx.gravity.v1.QueryCurrentValsetResponse) | Deprecated: Please use crosschain query.CurrentOracleSet | GET|/fx/gravity/v1/valset/current|
| `ValsetRequest` | [QueryValsetRequestRequest](#fx.gravity.v1.QueryValsetRequestRequest) | [QueryValsetRequestResponse](#fx.gravity.v1.QueryValsetRequestResponse) | Deprecated: Please use crosschain query.OracleSetRequest | GET|/fx/gravity/v1/valset/request|
| `ValsetConfirm` | [QueryValsetConfirmRequest](#fx.gravity.v1.QueryValsetConfirmRequest) | [QueryValsetConfirmResponse](#fx.gravity.v1.QueryValsetConfirmResponse) | Deprecated: Please use crosschain query.OracleSetConfirm | GET|/fx/gravity/v1/valset/confirm|
| `ValsetConfirmsByNonce` | [QueryValsetConfirmsByNonceRequest](#fx.gravity.v1.QueryValsetConfirmsByNonceRequest) | [QueryValsetConfirmsByNonceResponse](#fx.gravity.v1.QueryValsetConfirmsByNonceResponse) | Deprecated: Please use crosschain query.OracleSetConfirmsByNonce | GET|/fx/gravity/v1/valset/confirms|
| `LastValsetRequests` | [QueryLastValsetRequestsRequest](#fx.gravity.v1.QueryLastValsetRequestsRequest) | [QueryLastValsetRequestsResponse](#fx.gravity.v1.QueryLastValsetRequestsResponse) | Deprecated: Please use crosschain query.LastOracleSetRequests | GET|/fx/gravity/v1/valset/requests|
| `LastPendingValsetRequestByAddr` | [QueryLastPendingValsetRequestByAddrRequest](#fx.gravity.v1.QueryLastPendingValsetRequestByAddrRequest) | [QueryLastPendingValsetRequestByAddrResponse](#fx.gravity.v1.QueryLastPendingValsetRequestByAddrResponse) | Deprecated: Please use crosschain query.LastPendingOracleSetRequestByAddr | GET|/fx/gravity/v1/valset/last|
| `LastPendingBatchRequestByAddr` | [QueryLastPendingBatchRequestByAddrRequest](#fx.gravity.v1.QueryLastPendingBatchRequestByAddrRequest) | [QueryLastPendingBatchRequestByAddrResponse](#fx.gravity.v1.QueryLastPendingBatchRequestByAddrResponse) | Deprecated: Please use crosschain query.LastPendingBatchRequestByAddr | GET|/fx/gravity/v1/batch/last|
| `LastEventNonceByAddr` | [QueryLastEventNonceByAddrRequest](#fx.gravity.v1.QueryLastEventNonceByAddrRequest) | [QueryLastEventNonceByAddrResponse](#fx.gravity.v1.QueryLastEventNonceByAddrResponse) | Deprecated: Please use crosschain query.LastEventNonceByAddr | GET|/fx/gravity/v1/oracle/event_nonce/{address}|
| `LastEventBlockHeightByAddr` | [QueryLastEventBlockHeightByAddrRequest](#fx.gravity.v1.QueryLastEventBlockHeightByAddrRequest) | [QueryLastEventBlockHeightByAddrResponse](#fx.gravity.v1.QueryLastEventBlockHeightByAddrResponse) | Deprecated: Please use crosschain query.LastEventBlockHeightByAddr | GET|/fx/gravity/v1/oracle/event/block_height/{address}|
| `BatchFees` | [QueryBatchFeeRequest](#fx.gravity.v1.QueryBatchFeeRequest) | [QueryBatchFeeResponse](#fx.gravity.v1.QueryBatchFeeResponse) | Deprecated: Please use crosschain query.BatchFees | GET|/fx/gravity/v1/batch_fees|
| `LastObservedBlockHeight` | [QueryLastObservedBlockHeightRequest](#fx.gravity.v1.QueryLastObservedBlockHeightRequest) | [QueryLastObservedBlockHeightResponse](#fx.gravity.v1.QueryLastObservedBlockHeightResponse) | Deprecated: Please use crosschain query.LastObservedBlockHeight | GET|/fx/gravity/v1/observed/block_height|
| `OutgoingTxBatches` | [QueryOutgoingTxBatchesRequest](#fx.gravity.v1.QueryOutgoingTxBatchesRequest) | [QueryOutgoingTxBatchesResponse](#fx.gravity.v1.QueryOutgoingTxBatchesResponse) | Deprecated: Please use crosschain query.OutgoingTxBatches | GET|/fx/gravity/v1/batch/outgoing_tx|
| `BatchRequestByNonce` | [QueryBatchRequestByNonceRequest](#fx.gravity.v1.QueryBatchRequestByNonceRequest) | [QueryBatchRequestByNonceResponse](#fx.gravity.v1.QueryBatchRequestByNonceResponse) | Deprecated: Please use crosschain query.BatchRequestByNonce | GET|/fx/gravity/v1/batch/request|
| `BatchConfirm` | [QueryBatchConfirmRequest](#fx.gravity.v1.QueryBatchConfirmRequest) | [QueryBatchConfirmResponse](#fx.gravity.v1.QueryBatchConfirmResponse) | Deprecated: Please use crosschain query.BatchConfirm | GET|/fx/gravity/v1/batch/confirm|
| `BatchConfirms` | [QueryBatchConfirmsRequest](#fx.gravity.v1.QueryBatchConfirmsRequest) | [QueryBatchConfirmsResponse](#fx.gravity.v1.QueryBatchConfirmsResponse) | Deprecated: Please use crosschain query.BatchConfirms | GET|/fx/gravity/v1/batch/confirms|
| `ERC20ToDenom` | [QueryERC20ToDenomRequest](#fx.gravity.v1.QueryERC20ToDenomRequest) | [QueryERC20ToDenomResponse](#fx.gravity.v1.QueryERC20ToDenomResponse) | Deprecated: Please use crosschain query.TokenToDenom | GET|/fx/gravity/v1/denom|
| `DenomToERC20` | [QueryDenomToERC20Request](#fx.gravity.v1.QueryDenomToERC20Request) | [QueryDenomToERC20Response](#fx.gravity.v1.QueryDenomToERC20Response) | Deprecated: Please use crosschain query.DenomToToken | GET|/fx/gravity/v1/erc20|
| `GetDelegateKeyByValidator` | [QueryDelegateKeyByValidatorRequest](#fx.gravity.v1.QueryDelegateKeyByValidatorRequest) | [QueryDelegateKeyByValidatorResponse](#fx.gravity.v1.QueryDelegateKeyByValidatorResponse) | Deprecated: Please use crosschain query.GetOracleByAddr | GET|/fx/gravity/v1/delegate_key_by_validator|
| `GetDelegateKeyByEth` | [QueryDelegateKeyByEthRequest](#fx.gravity.v1.QueryDelegateKeyByEthRequest) | [QueryDelegateKeyByEthResponse](#fx.gravity.v1.QueryDelegateKeyByEthResponse) | Deprecated: Please use crosschain query.GetOracleByExternalAddr | GET|/fx/gravity/v1/delegate_key_by_eth|
| `GetDelegateKeyByOrchestrator` | [QueryDelegateKeyByOrchestratorRequest](#fx.gravity.v1.QueryDelegateKeyByOrchestratorRequest) | [QueryDelegateKeyByOrchestratorResponse](#fx.gravity.v1.QueryDelegateKeyByOrchestratorResponse) | Deprecated: Please use crosschain query.GetOracleByBridgerAddr | GET|/fx/gravity/v1/delegate_key_by_orchestrator|
| `GetPendingSendToEth` | [QueryPendingSendToEthRequest](#fx.gravity.v1.QueryPendingSendToEthRequest) | [QueryPendingSendToEthResponse](#fx.gravity.v1.QueryPendingSendToEthResponse) | Deprecated: Please use crosschain query.GetPendingSendToExternal | GET|/fx/gravity/v1/pending_send_to_eth|
| `ProjectedBatchTimeoutHeight` | [QueryProjectedBatchTimeoutHeightRequest](#fx.gravity.v1.QueryProjectedBatchTimeoutHeightRequest) | [QueryProjectedBatchTimeoutHeightResponse](#fx.gravity.v1.QueryProjectedBatchTimeoutHeightResponse) | Deprecated: Please use crosschain query.ProjectedBatchTimeoutHeight | GET|/fx/gravity/v1/projected_batch_timeout|
| `BridgeTokens` | [QueryBridgeTokensRequest](#fx.gravity.v1.QueryBridgeTokensRequest) | [QueryBridgeTokensResponse](#fx.gravity.v1.QueryBridgeTokensResponse) | Deprecated: Please use crosschain query.BridgeTokens | GET|/fx/gravity/v1/bridge_tokens|

 <!-- end services -->



<a name="fx/legacy/other/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/legacy/other/query.proto



<a name="fx.other.GasPriceRequest"></a>

### GasPriceRequest
Deprecated: GasPriceRequest






<a name="fx.other.GasPriceResponse"></a>

### GasPriceResponse
Deprecated: GasPriceResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas_prices` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.other.Query"></a>

### Query
Deprecated: Query

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `FxGasPrice` | [GasPriceRequest](#fx.other.GasPriceRequest) | [GasPriceResponse](#fx.other.GasPriceResponse) | Deprecated: Please use base query.GetGasPrice | GET|/fx/other/gas_price|
| `GasPrice` | [GasPriceRequest](#fx.other.GasPriceRequest) | [GasPriceResponse](#fx.other.GasPriceResponse) | Deprecated: Please use base query.GetGasPrice | GET|/other/v1/gas_price|

 <!-- end services -->



<a name="fx/migrate/v1/migrate.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/migrate/v1/migrate.proto



<a name="fx.migrate.v1.MigrateRecord"></a>

### MigrateRecord



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from` | [string](#string) |  |  |
| `to` | [string](#string) |  |  |
| `height` | [int64](#int64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="fx/migrate/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/migrate/v1/genesis.proto



<a name="fx.migrate.v1.GenesisState"></a>

### GenesisState
GenesisState defines the module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `migrate_records` | [MigrateRecord](#fx.migrate.v1.MigrateRecord) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="fx/migrate/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/migrate/v1/query.proto



<a name="fx.migrate.v1.QueryMigrateCheckAccountRequest"></a>

### QueryMigrateCheckAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from` | [string](#string) |  | migrate from address |
| `to` | [string](#string) |  | migrate to address |






<a name="fx.migrate.v1.QueryMigrateCheckAccountResponse"></a>

### QueryMigrateCheckAccountResponse







<a name="fx.migrate.v1.QueryMigrateRecordRequest"></a>

### QueryMigrateRecordRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  |  |






<a name="fx.migrate.v1.QueryMigrateRecordResponse"></a>

### QueryMigrateRecordResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `found` | [bool](#bool) |  | has migrate true-> migrated, false-> not migrated. |
| `migrateRecord` | [MigrateRecord](#fx.migrate.v1.MigrateRecord) |  | migrateRecord defines the the migrate record. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.migrate.v1.Query"></a>

### Query
Query provides defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `MigrateRecord` | [QueryMigrateRecordRequest](#fx.migrate.v1.QueryMigrateRecordRequest) | [QueryMigrateRecordResponse](#fx.migrate.v1.QueryMigrateRecordResponse) | DenomTrace queries a denomination trace information. | GET|/fx/migrate/v1/record/{address}|
| `MigrateCheckAccount` | [QueryMigrateCheckAccountRequest](#fx.migrate.v1.QueryMigrateCheckAccountRequest) | [QueryMigrateCheckAccountResponse](#fx.migrate.v1.QueryMigrateCheckAccountResponse) |  | GET|/fx/migrate/v1/check/account|

 <!-- end services -->



<a name="fx/migrate/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## fx/migrate/v1/tx.proto



<a name="fx.migrate.v1.MsgMigrateAccount"></a>

### MsgMigrateAccount



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `from` | [string](#string) |  |  |
| `to` | [string](#string) |  |  |
| `signature` | [string](#string) |  |  |






<a name="fx.migrate.v1.MsgMigrateAccountResponse"></a>

### MsgMigrateAccountResponse






 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.migrate.v1.Msg"></a>

### Msg
Msg defines the state transitions possible within gravity

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `MigrateAccount` | [MsgMigrateAccount](#fx.migrate.v1.MsgMigrateAccount) | [MsgMigrateAccountResponse](#fx.migrate.v1.MsgMigrateAccountResponse) |  | |

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
