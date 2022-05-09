<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [crosschain/v1/crosschain.proto](#crosschain/v1/crosschain.proto)
    - [Attestation](#fx.gravity.crosschain.v1.Attestation)
    - [BatchFees](#fx.gravity.crosschain.v1.BatchFees)
    - [BridgeToken](#fx.gravity.crosschain.v1.BridgeToken)
    - [BridgeValidator](#fx.gravity.crosschain.v1.BridgeValidator)
    - [ChainOracle](#fx.gravity.crosschain.v1.ChainOracle)
    - [ExternalToken](#fx.gravity.crosschain.v1.ExternalToken)
    - [IDSet](#fx.gravity.crosschain.v1.IDSet)
    - [InitCrossChainParamsProposal](#fx.gravity.crosschain.v1.InitCrossChainParamsProposal)
    - [LastObservedBlockHeight](#fx.gravity.crosschain.v1.LastObservedBlockHeight)
    - [MinBatchFee](#fx.gravity.crosschain.v1.MinBatchFee)
    - [Oracle](#fx.gravity.crosschain.v1.Oracle)
    - [OracleSet](#fx.gravity.crosschain.v1.OracleSet)
    - [OutgoingTransferTx](#fx.gravity.crosschain.v1.OutgoingTransferTx)
    - [OutgoingTxBatch](#fx.gravity.crosschain.v1.OutgoingTxBatch)
    - [Params](#fx.gravity.crosschain.v1.Params)
    - [UpdateChainOraclesProposal](#fx.gravity.crosschain.v1.UpdateChainOraclesProposal)
  
    - [ClaimType](#fx.gravity.crosschain.v1.ClaimType)
    - [SignType](#fx.gravity.crosschain.v1.SignType)
  
- [crosschain/v1/genesis.proto](#crosschain/v1/genesis.proto)
    - [GenesisState](#fx.gravity.crosschain.v1.GenesisState)
  
- [crosschain/v1/tx.proto](#crosschain/v1/tx.proto)
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
  
- [crosschain/v1/query.proto](#crosschain/v1/query.proto)
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
    - [QueryProjectedBatchTimeoutHeightRequest](#fx.gravity.crosschain.v1.QueryProjectedBatchTimeoutHeightRequest)
    - [QueryProjectedBatchTimeoutHeightResponse](#fx.gravity.crosschain.v1.QueryProjectedBatchTimeoutHeightResponse)
    - [QueryTokenToDenomRequest](#fx.gravity.crosschain.v1.QueryTokenToDenomRequest)
    - [QueryTokenToDenomResponse](#fx.gravity.crosschain.v1.QueryTokenToDenomResponse)
  
    - [Query](#fx.gravity.crosschain.v1.Query)
  
- [ethereum/erc20/v1/erc20.proto](#ethereum/erc20/v1/erc20.proto)
    - [RegisterCoinProposal](#fx.ethereum.erc20.v1.RegisterCoinProposal)
    - [RegisterERC20Proposal](#fx.ethereum.erc20.v1.RegisterERC20Proposal)
    - [ToggleTokenRelayProposal](#fx.ethereum.erc20.v1.ToggleTokenRelayProposal)
    - [TokenPair](#fx.ethereum.erc20.v1.TokenPair)
  
    - [Owner](#fx.ethereum.erc20.v1.Owner)
  
- [ethereum/feemarket/v1/feemarket.proto](#ethereum/feemarket/v1/feemarket.proto)
    - [Params](#fx.ethereum.feemarket.v1.Params)
  
- [ethereum/evm/v1/evm.proto](#ethereum/evm/v1/evm.proto)
    - [AccessTuple](#fx.ethereum.evm.v1.AccessTuple)
    - [ChainConfig](#fx.ethereum.evm.v1.ChainConfig)
    - [Log](#fx.ethereum.evm.v1.Log)
    - [Params](#fx.ethereum.evm.v1.Params)
    - [State](#fx.ethereum.evm.v1.State)
    - [TraceConfig](#fx.ethereum.evm.v1.TraceConfig)
    - [TransactionLogs](#fx.ethereum.evm.v1.TransactionLogs)
    - [TxResult](#fx.ethereum.evm.v1.TxResult)
  
- [ethereum/erc20/v1/genesis.proto](#ethereum/erc20/v1/genesis.proto)
    - [GenesisState](#fx.ethereum.erc20.v1.GenesisState)
    - [Params](#fx.ethereum.erc20.v1.Params)
  
- [ethereum/erc20/v1/query.proto](#ethereum/erc20/v1/query.proto)
    - [QueryParamsRequest](#fx.ethereum.erc20.v1.QueryParamsRequest)
    - [QueryParamsResponse](#fx.ethereum.erc20.v1.QueryParamsResponse)
    - [QueryTokenPairRequest](#fx.ethereum.erc20.v1.QueryTokenPairRequest)
    - [QueryTokenPairResponse](#fx.ethereum.erc20.v1.QueryTokenPairResponse)
    - [QueryTokenPairsRequest](#fx.ethereum.erc20.v1.QueryTokenPairsRequest)
    - [QueryTokenPairsResponse](#fx.ethereum.erc20.v1.QueryTokenPairsResponse)
  
    - [Query](#fx.ethereum.erc20.v1.Query)
  
- [ethereum/erc20/v1/tx.proto](#ethereum/erc20/v1/tx.proto)
    - [MsgConvertCoin](#fx.ethereum.erc20.v1.MsgConvertCoin)
    - [MsgConvertCoinResponse](#fx.ethereum.erc20.v1.MsgConvertCoinResponse)
    - [MsgConvertERC20](#fx.ethereum.erc20.v1.MsgConvertERC20)
    - [MsgConvertERC20Response](#fx.ethereum.erc20.v1.MsgConvertERC20Response)
  
    - [Msg](#fx.ethereum.erc20.v1.Msg)
  
- [ethereum/evm/v1/genesis.proto](#ethereum/evm/v1/genesis.proto)
    - [GenesisAccount](#fx.ethereum.evm.v1.GenesisAccount)
    - [GenesisState](#fx.ethereum.evm.v1.GenesisState)
  
- [ethereum/evm/v1/tx.proto](#ethereum/evm/v1/tx.proto)
    - [AccessListTx](#fx.ethereum.evm.v1.AccessListTx)
    - [DynamicFeeTx](#fx.ethereum.evm.v1.DynamicFeeTx)
    - [ExtensionOptionsEthereumTx](#fx.ethereum.evm.v1.ExtensionOptionsEthereumTx)
    - [LegacyTx](#fx.ethereum.evm.v1.LegacyTx)
    - [MsgEthereumTx](#fx.ethereum.evm.v1.MsgEthereumTx)
    - [MsgEthereumTxResponse](#fx.ethereum.evm.v1.MsgEthereumTxResponse)
  
    - [Msg](#fx.ethereum.evm.v1.Msg)
  
- [ethereum/evm/v1/query.proto](#ethereum/evm/v1/query.proto)
    - [EstimateGasResponse](#fx.ethereum.evm.v1.EstimateGasResponse)
    - [EthCallRequest](#fx.ethereum.evm.v1.EthCallRequest)
    - [QueryAccountRequest](#fx.ethereum.evm.v1.QueryAccountRequest)
    - [QueryAccountResponse](#fx.ethereum.evm.v1.QueryAccountResponse)
    - [QueryBalanceRequest](#fx.ethereum.evm.v1.QueryBalanceRequest)
    - [QueryBalanceResponse](#fx.ethereum.evm.v1.QueryBalanceResponse)
    - [QueryCodeRequest](#fx.ethereum.evm.v1.QueryCodeRequest)
    - [QueryCodeResponse](#fx.ethereum.evm.v1.QueryCodeResponse)
    - [QueryCosmosAccountRequest](#fx.ethereum.evm.v1.QueryCosmosAccountRequest)
    - [QueryCosmosAccountResponse](#fx.ethereum.evm.v1.QueryCosmosAccountResponse)
    - [QueryModuleEnableRequest](#fx.ethereum.evm.v1.QueryModuleEnableRequest)
    - [QueryModuleEnableResponse](#fx.ethereum.evm.v1.QueryModuleEnableResponse)
    - [QueryParamsRequest](#fx.ethereum.evm.v1.QueryParamsRequest)
    - [QueryParamsResponse](#fx.ethereum.evm.v1.QueryParamsResponse)
    - [QueryStorageRequest](#fx.ethereum.evm.v1.QueryStorageRequest)
    - [QueryStorageResponse](#fx.ethereum.evm.v1.QueryStorageResponse)
    - [QueryTraceBlockRequest](#fx.ethereum.evm.v1.QueryTraceBlockRequest)
    - [QueryTraceBlockResponse](#fx.ethereum.evm.v1.QueryTraceBlockResponse)
    - [QueryTraceTxRequest](#fx.ethereum.evm.v1.QueryTraceTxRequest)
    - [QueryTraceTxResponse](#fx.ethereum.evm.v1.QueryTraceTxResponse)
    - [QueryTxLogsRequest](#fx.ethereum.evm.v1.QueryTxLogsRequest)
    - [QueryTxLogsResponse](#fx.ethereum.evm.v1.QueryTxLogsResponse)
    - [QueryValidatorAccountRequest](#fx.ethereum.evm.v1.QueryValidatorAccountRequest)
    - [QueryValidatorAccountResponse](#fx.ethereum.evm.v1.QueryValidatorAccountResponse)
  
    - [Query](#fx.ethereum.evm.v1.Query)
  
- [ethereum/feemarket/v1/genesis.proto](#ethereum/feemarket/v1/genesis.proto)
    - [GenesisState](#fx.ethereum.feemarket.v1.GenesisState)
  
- [ethereum/feemarket/v1/query.proto](#ethereum/feemarket/v1/query.proto)
    - [QueryBaseFeeRequest](#fx.ethereum.feemarket.v1.QueryBaseFeeRequest)
    - [QueryBaseFeeResponse](#fx.ethereum.feemarket.v1.QueryBaseFeeResponse)
    - [QueryBlockGasRequest](#fx.ethereum.feemarket.v1.QueryBlockGasRequest)
    - [QueryBlockGasResponse](#fx.ethereum.feemarket.v1.QueryBlockGasResponse)
    - [QueryModuleEnableRequest](#fx.ethereum.feemarket.v1.QueryModuleEnableRequest)
    - [QueryModuleEnableResponse](#fx.ethereum.feemarket.v1.QueryModuleEnableResponse)
    - [QueryParamsRequest](#fx.ethereum.feemarket.v1.QueryParamsRequest)
    - [QueryParamsResponse](#fx.ethereum.feemarket.v1.QueryParamsResponse)
  
    - [Query](#fx.ethereum.feemarket.v1.Query)
  
- [ethereum/types/v1/web3.proto](#ethereum/types/v1/web3.proto)
    - [ExtensionOptionsWeb3Tx](#fx.ethereum.types.v1.ExtensionOptionsWeb3Tx)
  
- [gravity/v1/attestation.proto](#gravity/v1/attestation.proto)
    - [Attestation](#fx.gravity.v1.Attestation)
    - [ERC20Token](#fx.gravity.v1.ERC20Token)
  
    - [ClaimType](#fx.gravity.v1.ClaimType)
  
- [gravity/v1/batch.proto](#gravity/v1/batch.proto)
    - [OutgoingTransferTx](#fx.gravity.v1.OutgoingTransferTx)
    - [OutgoingTxBatch](#fx.gravity.v1.OutgoingTxBatch)
  
- [gravity/v1/ethereum_signer.proto](#gravity/v1/ethereum_signer.proto)
    - [SignType](#fx.gravity.v1.SignType)
  
- [gravity/v1/types.proto](#gravity/v1/types.proto)
    - [BridgeValidator](#fx.gravity.v1.BridgeValidator)
    - [ERC20ToDenom](#fx.gravity.v1.ERC20ToDenom)
    - [LastObservedEthereumBlockHeight](#fx.gravity.v1.LastObservedEthereumBlockHeight)
    - [Valset](#fx.gravity.v1.Valset)
  
- [gravity/v1/tx.proto](#gravity/v1/tx.proto)
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
  
- [gravity/v1/genesis.proto](#gravity/v1/genesis.proto)
    - [GenesisState](#fx.gravity.v1.GenesisState)
    - [Params](#fx.gravity.v1.Params)
  
- [gravity/v1/pool.proto](#gravity/v1/pool.proto)
    - [BatchFees](#fx.gravity.v1.BatchFees)
    - [IDSet](#fx.gravity.v1.IDSet)
    - [MinBatchFee](#fx.gravity.v1.MinBatchFee)
  
- [gravity/v1/query.proto](#gravity/v1/query.proto)
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
    - [QueryIbcSequenceHeightRequest](#fx.gravity.v1.QueryIbcSequenceHeightRequest)
    - [QueryIbcSequenceHeightResponse](#fx.gravity.v1.QueryIbcSequenceHeightResponse)
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
  
- [ibc/applications/transfer/v1/transfer.proto](#ibc/applications/transfer/v1/transfer.proto)
    - [DenomTrace](#fx.ibc.applications.transfer.v1.DenomTrace)
    - [FungibleTokenPacketData](#fx.ibc.applications.transfer.v1.FungibleTokenPacketData)
    - [Params](#fx.ibc.applications.transfer.v1.Params)
  
- [ibc/applications/transfer/v1/genesis.proto](#ibc/applications/transfer/v1/genesis.proto)
    - [GenesisState](#fx.ibc.applications.transfer.v1.GenesisState)
  
- [ibc/applications/transfer/v1/query.proto](#ibc/applications/transfer/v1/query.proto)
    - [QueryDenomTraceRequest](#fx.ibc.applications.transfer.v1.QueryDenomTraceRequest)
    - [QueryDenomTraceResponse](#fx.ibc.applications.transfer.v1.QueryDenomTraceResponse)
    - [QueryDenomTracesRequest](#fx.ibc.applications.transfer.v1.QueryDenomTracesRequest)
    - [QueryDenomTracesResponse](#fx.ibc.applications.transfer.v1.QueryDenomTracesResponse)
    - [QueryParamsRequest](#fx.ibc.applications.transfer.v1.QueryParamsRequest)
    - [QueryParamsResponse](#fx.ibc.applications.transfer.v1.QueryParamsResponse)
  
    - [Query](#fx.ibc.applications.transfer.v1.Query)
  
- [ibc/applications/transfer/v1/tx.proto](#ibc/applications/transfer/v1/tx.proto)
    - [MsgTransfer](#fx.ibc.applications.transfer.v1.MsgTransfer)
    - [MsgTransferResponse](#fx.ibc.applications.transfer.v1.MsgTransferResponse)
  
    - [Msg](#fx.ibc.applications.transfer.v1.Msg)
  
- [migrate/v1/migrate.proto](#migrate/v1/migrate.proto)
    - [MigrateRecord](#fx.migrate.v1.MigrateRecord)
  
- [migrate/v1/query.proto](#migrate/v1/query.proto)
    - [QueryMigrateRecordRequest](#fx.migrate.v1.QueryMigrateRecordRequest)
    - [QueryMigrateRecordResponse](#fx.migrate.v1.QueryMigrateRecordResponse)
  
    - [Query](#fx.migrate.v1.Query)
  
- [migrate/v1/tx.proto](#migrate/v1/tx.proto)
    - [MsgMigrateAccount](#fx.migrate.v1.MsgMigrateAccount)
    - [MsgMigrateAccountResponse](#fx.migrate.v1.MsgMigrateAccountResponse)
  
    - [Msg](#fx.migrate.v1.Msg)
  
- [other/query.proto](#other/query.proto)
    - [GasPriceRequest](#fx.other.GasPriceRequest)
    - [GasPriceResponse](#fx.other.GasPriceResponse)
  
    - [Query](#fx.other.Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="crosschain/v1/crosschain.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## crosschain/v1/crosschain.proto



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
| `total_amount` | [string](#string) |  |  |






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



<a name="crosschain/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## crosschain/v1/genesis.proto



<a name="fx.gravity.crosschain.v1.GenesisState"></a>

### GenesisState
GenesisState struct





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="crosschain/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## crosschain/v1/tx.proto



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
| `base_fee` | [string](#string) |  |  |






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



<a name="crosschain/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## crosschain/v1/query.proto



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
| `nonce` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
| `chain_name` | [string](#string) |  |  |






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
| `LastObservedBlockHeight` | [QueryLastObservedBlockHeightRequest](#fx.gravity.crosschain.v1.QueryLastObservedBlockHeightRequest) | [QueryLastObservedBlockHeightResponse](#fx.gravity.crosschain.v1.QueryLastObservedBlockHeightResponse) |  | GET|/crosschain/v1beta/observed/block_height|
| `OutgoingTxBatches` | [QueryOutgoingTxBatchesRequest](#fx.gravity.crosschain.v1.QueryOutgoingTxBatchesRequest) | [QueryOutgoingTxBatchesResponse](#fx.gravity.crosschain.v1.QueryOutgoingTxBatchesResponse) |  | GET|/crosschain/v1beta/batch/outgoing_tx|
| `BatchRequestByNonce` | [QueryBatchRequestByNonceRequest](#fx.gravity.crosschain.v1.QueryBatchRequestByNonceRequest) | [QueryBatchRequestByNonceResponse](#fx.gravity.crosschain.v1.QueryBatchRequestByNonceResponse) |  | GET|/crosschain/v1beta/batch/request|
| `BatchConfirm` | [QueryBatchConfirmRequest](#fx.gravity.crosschain.v1.QueryBatchConfirmRequest) | [QueryBatchConfirmResponse](#fx.gravity.crosschain.v1.QueryBatchConfirmResponse) |  | GET|/crosschain/v1beta/batch/confirm|
| `BatchConfirms` | [QueryBatchConfirmsRequest](#fx.gravity.crosschain.v1.QueryBatchConfirmsRequest) | [QueryBatchConfirmsResponse](#fx.gravity.crosschain.v1.QueryBatchConfirmsResponse) |  | GET|/crosschain/v1beta/batch/confirms|
| `TokenToDenom` | [QueryTokenToDenomRequest](#fx.gravity.crosschain.v1.QueryTokenToDenomRequest) | [QueryTokenToDenomResponse](#fx.gravity.crosschain.v1.QueryTokenToDenomResponse) |  | GET|/crosschain/v1beta/denom|
| `DenomToToken` | [QueryDenomToTokenRequest](#fx.gravity.crosschain.v1.QueryDenomToTokenRequest) | [QueryDenomToTokenResponse](#fx.gravity.crosschain.v1.QueryDenomToTokenResponse) |  | GET|/crosschain/v1beta/token|
| `GetOracleByAddr` | [QueryOracleByAddrRequest](#fx.gravity.crosschain.v1.QueryOracleByAddrRequest) | [QueryOracleResponse](#fx.gravity.crosschain.v1.QueryOracleResponse) |  | GET|/crosschain/v1beta/oracle_by_addr|
| `GetOracleByExternalAddr` | [QueryOracleByExternalAddrRequest](#fx.gravity.crosschain.v1.QueryOracleByExternalAddrRequest) | [QueryOracleResponse](#fx.gravity.crosschain.v1.QueryOracleResponse) |  | GET|/crosschain/v1beta/oracle_by_external_addr|
| `GetOracleByOrchestrator` | [QueryOracleByOrchestratorRequest](#fx.gravity.crosschain.v1.QueryOracleByOrchestratorRequest) | [QueryOracleResponse](#fx.gravity.crosschain.v1.QueryOracleResponse) |  | GET|/crosschain/v1beta/oracle_by_orchestrator|
| `GetPendingSendToExternal` | [QueryPendingSendToExternalRequest](#fx.gravity.crosschain.v1.QueryPendingSendToExternalRequest) | [QueryPendingSendToExternalResponse](#fx.gravity.crosschain.v1.QueryPendingSendToExternalResponse) |  | GET|/crosschain/v1beta/pending_send_to_external|
| `GetIbcSequenceHeightByChannel` | [QueryIbcSequenceHeightRequest](#fx.gravity.crosschain.v1.QueryIbcSequenceHeightRequest) | [QueryIbcSequenceHeightResponse](#fx.gravity.crosschain.v1.QueryIbcSequenceHeightResponse) |  | GET|/crosschain/v1beta/ibc_sequence_height|
| `Oracles` | [QueryOraclesRequest](#fx.gravity.crosschain.v1.QueryOraclesRequest) | [QueryOraclesResponse](#fx.gravity.crosschain.v1.QueryOraclesResponse) | Validators queries all oracle that match the given status. | GET|/crosschain/v1beta1/oracles|
| `ProjectedBatchTimeoutHeight` | [QueryProjectedBatchTimeoutHeightRequest](#fx.gravity.crosschain.v1.QueryProjectedBatchTimeoutHeightRequest) | [QueryProjectedBatchTimeoutHeightResponse](#fx.gravity.crosschain.v1.QueryProjectedBatchTimeoutHeightResponse) |  | GET|/crosschain/v1beta1/projected_batch_timeout|
| `BridgeTokens` | [QueryBridgeTokensRequest](#fx.gravity.crosschain.v1.QueryBridgeTokensRequest) | [QueryBridgeTokensResponse](#fx.gravity.crosschain.v1.QueryBridgeTokensResponse) |  | GET|/gravity/v1beta1/bridge_tokens|

 <!-- end services -->



<a name="ethereum/erc20/v1/erc20.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethereum/erc20/v1/erc20.proto



<a name="fx.ethereum.erc20.v1.RegisterCoinProposal"></a>

### RegisterCoinProposal
RegisterCoinProposal is a gov Content type to register a token pair


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | title of the proposal |
| `description` | [string](#string) |  | proposal description |
| `metadata` | [cosmos.bank.v1beta1.Metadata](#cosmos.bank.v1beta1.Metadata) |  | token pair of Cosmos native denom and ERC20 token address |






<a name="fx.ethereum.erc20.v1.RegisterERC20Proposal"></a>

### RegisterERC20Proposal
RegisterCoinProposal is a gov Content type to register a token pair


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | title of the proposal |
| `description` | [string](#string) |  | proposal description |
| `erc20address` | [string](#string) |  | contract address of ERC20 token |






<a name="fx.ethereum.erc20.v1.ToggleTokenRelayProposal"></a>

### ToggleTokenRelayProposal
ToggleTokenRelayProposal is a gov Content type to toggle
the internal relaying of a token pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | title of the proposal |
| `description` | [string](#string) |  | proposal description |
| `token` | [string](#string) |  | token identifier can be either the hex contract address of the ERC20 or the Cosmos base denomination |






<a name="fx.ethereum.erc20.v1.TokenPair"></a>

### TokenPair
TokenPair defines an instance that records pairing consisting of a Cosmos
native Coin and an ERC20 token address.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `erc20_address` | [string](#string) |  | address of ERC20 contract token |
| `denom` | [string](#string) |  | cosmos base denomination to be mapped to |
| `enabled` | [bool](#bool) |  | shows token mapping enable status |
| `contract_owner` | [Owner](#fx.ethereum.erc20.v1.Owner) |  | ERC20 owner address ENUM (0 invalid, 1 ModuleAccount, 2 external address) |





 <!-- end messages -->


<a name="fx.ethereum.erc20.v1.Owner"></a>

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



<a name="ethereum/feemarket/v1/feemarket.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethereum/feemarket/v1/feemarket.proto



<a name="fx.ethereum.feemarket.v1.Params"></a>

### Params
Params defines the EVM module parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_fee_change_denominator` | [uint32](#uint32) |  | base fee change denominator bounds the amount the base fee can change between blocks. |
| `elasticity_multiplier` | [uint32](#uint32) |  | elasticity multiplier bounds the maximum gas limit an EIP-1559 block may have. |
| `base_fee` | [string](#string) |  | base fee for EIP-1559 blocks. |
| `min_base_fee` | [string](#string) |  | min base fee for EIP-1559 blocks. |
| `max_base_fee` | [string](#string) |  | max base fee for EIP-1559 blocks. |
| `max_gas` | [string](#string) |  | replace block max gas, if > 0 |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethereum/evm/v1/evm.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethereum/evm/v1/evm.proto



<a name="fx.ethereum.evm.v1.AccessTuple"></a>

### AccessTuple
AccessTuple is the element type of an access list.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | hex formatted ethereum address |
| `storage_keys` | [string](#string) | repeated | hex formatted hashes of the storage keys |






<a name="fx.ethereum.evm.v1.ChainConfig"></a>

### ChainConfig
ChainConfig defines the Ethereum ChainConfig parameters using *sdk.Int values
instead of *big.Int.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `homestead_block` | [string](#string) |  | Homestead switch block (nil no fork, 0 = already homestead) |
| `dao_fork_block` | [string](#string) |  | TheDAO hard-fork switch block (nil no fork) |
| `dao_fork_support` | [bool](#bool) |  | Whether the nodes supports or opposes the DAO hard-fork |
| `eip150_block` | [string](#string) |  | EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150) EIP150 HF block (nil no fork) |
| `eip150_hash` | [string](#string) |  | EIP150 HF hash (needed for header only clients as only gas pricing changed) |
| `eip155_block` | [string](#string) |  | EIP155Block HF block |
| `eip158_block` | [string](#string) |  | EIP158 HF block |
| `byzantium_block` | [string](#string) |  | Byzantium switch block (nil no fork, 0 = already on byzantium) |
| `constantinople_block` | [string](#string) |  | Constantinople switch block (nil no fork, 0 = already activated) |
| `petersburg_block` | [string](#string) |  | Petersburg switch block (nil same as Constantinople) |
| `istanbul_block` | [string](#string) |  | Istanbul switch block (nil no fork, 0 = already on istanbul) |
| `muir_glacier_block` | [string](#string) |  | Eip-2384 (bomb delay) switch block (nil no fork, 0 = already activated) |
| `berlin_block` | [string](#string) |  | Berlin switch block (nil = no fork, 0 = already on berlin) |
| `london_block` | [string](#string) |  | London switch block (nil = no fork, 0 = already on london) |
| `arrow_glacier_block` | [string](#string) |  | Eip-4345 (bomb delay) switch block (nil = no fork, 0 = already activated) |
| `merge_fork_block` | [string](#string) |  | EIP-3675 (TheMerge) switch block (nil = no fork, 0 = already in merge proceedings) |






<a name="fx.ethereum.evm.v1.Log"></a>

### Log
Log represents an protobuf compatible Ethereum Log that defines a contract
log event. These events are generated by the LOG opcode and stored/indexed by
the node.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address of the contract that generated the event |
| `topics` | [string](#string) | repeated | list of topics provided by the contract. |
| `data` | [bytes](#bytes) |  | supplied by the contract, usually ABI-encoded |
| `block_number` | [uint64](#uint64) |  | block in which the transaction was included |
| `tx_hash` | [string](#string) |  | hash of the transaction |
| `tx_index` | [uint64](#uint64) |  | index of the transaction in the block |
| `block_hash` | [string](#string) |  | hash of the block in which the transaction was included |
| `index` | [uint64](#uint64) |  | index of the log in the block |
| `removed` | [bool](#bool) |  | The Removed field is true if this log was reverted due to a chain reorganisation. You must pay attention to this field if you receive logs through a filter query. |






<a name="fx.ethereum.evm.v1.Params"></a>

### Params
Params defines the EVM module parameters


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enable_create` | [bool](#bool) |  | enable create toggles state transitions that use the vm.Create function |
| `enable_call` | [bool](#bool) |  | enable call toggles state transitions that use the vm.Call function |
| `extra_eips` | [int64](#int64) | repeated | extra eips defines the additional EIPs for the vm.Config |






<a name="fx.ethereum.evm.v1.State"></a>

### State
State represents a single Storage key value pair item.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [string](#string) |  |  |
| `value` | [string](#string) |  |  |






<a name="fx.ethereum.evm.v1.TraceConfig"></a>

### TraceConfig
TraceConfig holds extra parameters to trace functions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `tracer` | [string](#string) |  | custom javascript tracer |
| `timeout` | [string](#string) |  | overrides the default timeout of 5 seconds for JavaScript-based tracing calls |
| `reexec` | [uint64](#uint64) |  | number of blocks the tracer is willing to go back |
| `disable_stack` | [bool](#bool) |  | disable stack capture |
| `disable_storage` | [bool](#bool) |  | disable storage capture |
| `debug` | [bool](#bool) |  | print output during capture end |
| `limit` | [int32](#int32) |  | maximum length of output, but zero means unlimited |
| `overrides` | [ChainConfig](#fx.ethereum.evm.v1.ChainConfig) |  | Chain overrides, can be used to execute a trace using future fork rules |
| `enable_memory` | [bool](#bool) |  | enable memory capture |
| `enable_return_data` | [bool](#bool) |  | enable return data capture |






<a name="fx.ethereum.evm.v1.TransactionLogs"></a>

### TransactionLogs
TransactionLogs define the logs generated from a transaction execution
with a given hash. It it used for import/export data as transactions are not
persisted on blockchain state after an upgrade.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  |  |
| `logs` | [Log](#fx.ethereum.evm.v1.Log) | repeated |  |






<a name="fx.ethereum.evm.v1.TxResult"></a>

### TxResult
TxResult stores results of Tx execution.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  | contract_address contains the ethereum address of the created contract (if any). If the state transition is an evm.Call, the contract address will be empty. |
| `bloom` | [bytes](#bytes) |  | bloom represents the bloom filter bytes |
| `tx_logs` | [TransactionLogs](#fx.ethereum.evm.v1.TransactionLogs) |  | tx_logs contains the transaction hash and the proto-compatible ethereum logs. |
| `ret` | [bytes](#bytes) |  | ret defines the bytes from the execution. |
| `reverted` | [bool](#bool) |  | reverted flag is set to true when the call has been reverted |
| `gas_used` | [uint64](#uint64) |  | gas_used notes the amount of gas consumed while execution |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethereum/erc20/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethereum/erc20/v1/genesis.proto



<a name="fx.ethereum.erc20.v1.GenesisState"></a>

### GenesisState
GenesisState defines the module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#fx.ethereum.erc20.v1.Params) |  | module parameters |
| `token_pairs` | [TokenPair](#fx.ethereum.erc20.v1.TokenPair) | repeated | registered token pairs |






<a name="fx.ethereum.erc20.v1.Params"></a>

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



<a name="ethereum/erc20/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethereum/erc20/v1/query.proto



<a name="fx.ethereum.erc20.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="fx.ethereum.erc20.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#fx.ethereum.erc20.v1.Params) |  |  |






<a name="fx.ethereum.erc20.v1.QueryTokenPairRequest"></a>

### QueryTokenPairRequest
QueryTokenPairRequest is the request type for the Query/TokenPair RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token` | [string](#string) |  | token identifier can be either the hex contract address of the ERC20 or the Cosmos base denomination |






<a name="fx.ethereum.erc20.v1.QueryTokenPairResponse"></a>

### QueryTokenPairResponse
QueryTokenPairResponse is the response type for the Query/TokenPair RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_pair` | [TokenPair](#fx.ethereum.erc20.v1.TokenPair) |  |  |






<a name="fx.ethereum.erc20.v1.QueryTokenPairsRequest"></a>

### QueryTokenPairsRequest
QueryTokenPairsRequest is the request type for the Query/TokenPairs RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="fx.ethereum.erc20.v1.QueryTokenPairsResponse"></a>

### QueryTokenPairsResponse
QueryTokenPairsResponse is the response type for the Query/TokenPairs RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_pairs` | [TokenPair](#fx.ethereum.erc20.v1.TokenPair) | repeated |  |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.ethereum.erc20.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `TokenPairs` | [QueryTokenPairsRequest](#fx.ethereum.erc20.v1.QueryTokenPairsRequest) | [QueryTokenPairsResponse](#fx.ethereum.erc20.v1.QueryTokenPairsResponse) | Retrieves registered token pairs | GET|/ethereum/erc20/v1/token_pairs|
| `TokenPair` | [QueryTokenPairRequest](#fx.ethereum.erc20.v1.QueryTokenPairRequest) | [QueryTokenPairResponse](#fx.ethereum.erc20.v1.QueryTokenPairResponse) | Retrieves a registered token pair | GET|/ethereum/erc20/v1/token_pairs/{token}|
| `Params` | [QueryParamsRequest](#fx.ethereum.erc20.v1.QueryParamsRequest) | [QueryParamsResponse](#fx.ethereum.erc20.v1.QueryParamsResponse) | Params retrieves the erc20 module params | GET|/ethereum/erc20/v1/params|

 <!-- end services -->



<a name="ethereum/erc20/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethereum/erc20/v1/tx.proto



<a name="fx.ethereum.erc20.v1.MsgConvertCoin"></a>

### MsgConvertCoin
MsgConvertCoin defines a Msg to convert a Cosmos Coin to a ERC20 token


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `coin` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | Cosmos coin which denomination is registered on erc20 bridge. The coin amount defines the total ERC20 tokens to convert. |
| `receiver` | [string](#string) |  | recipient hex address to receive ERC20 token |
| `sender` | [string](#string) |  | cosmos bech32 address from the owner of the given ERC20 tokens |






<a name="fx.ethereum.erc20.v1.MsgConvertCoinResponse"></a>

### MsgConvertCoinResponse
MsgConvertCoinResponse returns no fields






<a name="fx.ethereum.erc20.v1.MsgConvertERC20"></a>

### MsgConvertERC20
MsgConvertERC20 defines a Msg to convert an ERC20 token to a Cosmos SDK coin.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  | ERC20 token contract address registered on erc20 bridge |
| `amount` | [string](#string) |  | amount of ERC20 tokens to mint |
| `receiver` | [string](#string) |  | bech32 address to receive SDK coins. |
| `sender` | [string](#string) |  | sender hex address from the owner of the given ERC20 tokens |






<a name="fx.ethereum.erc20.v1.MsgConvertERC20Response"></a>

### MsgConvertERC20Response
MsgConvertERC20Response returns no fields





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.ethereum.erc20.v1.Msg"></a>

### Msg
Msg defines the erc20 Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `ConvertCoin` | [MsgConvertCoin](#fx.ethereum.erc20.v1.MsgConvertCoin) | [MsgConvertCoinResponse](#fx.ethereum.erc20.v1.MsgConvertCoinResponse) | ConvertCoin mints a ERC20 representation of the SDK Coin denom that is registered on the token mapping. | GET|/erc20/v1/tx/convert_coin|
| `ConvertERC20` | [MsgConvertERC20](#fx.ethereum.erc20.v1.MsgConvertERC20) | [MsgConvertERC20Response](#fx.ethereum.erc20.v1.MsgConvertERC20Response) | ConvertERC20 mints a Cosmos coin representation of the ERC20 token contract that is registered on the token mapping. | GET|/erc20/v1/tx/convert_erc20|

 <!-- end services -->



<a name="ethereum/evm/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethereum/evm/v1/genesis.proto



<a name="fx.ethereum.evm.v1.GenesisAccount"></a>

### GenesisAccount
GenesisAccount defines an account to be initialized in the genesis state.
Its main difference between with Geth's GenesisAccount is that it uses a
custom storage type and that it doesn't contain the private key field.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address defines an ethereum hex formated address of an account |
| `code` | [string](#string) |  | code defines the hex bytes of the account code. |
| `storage` | [State](#fx.ethereum.evm.v1.State) | repeated | storage defines the set of state key values for the account. |






<a name="fx.ethereum.evm.v1.GenesisState"></a>

### GenesisState
GenesisState defines the evm module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `accounts` | [GenesisAccount](#fx.ethereum.evm.v1.GenesisAccount) | repeated | accounts is an array containing the ethereum genesis accounts. |
| `params` | [Params](#fx.ethereum.evm.v1.Params) |  | params defines all the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethereum/evm/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethereum/evm/v1/tx.proto



<a name="fx.ethereum.evm.v1.AccessListTx"></a>

### AccessListTx
AccessListTx is the data of EIP-2930 access list transactions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  | destination EVM chain ID |
| `nonce` | [uint64](#uint64) |  | nonce corresponds to the account nonce (transaction sequence). |
| `gas_price` | [string](#string) |  | gas price defines the value for each gas unit |
| `gas` | [uint64](#uint64) |  | gas defines the gas limit defined for the transaction. |
| `to` | [string](#string) |  | hex formatted address of the recipient |
| `value` | [string](#string) |  | value defines the unsigned integer value of the transaction amount. |
| `data` | [bytes](#bytes) |  | input defines the data payload bytes of the transaction. |
| `accesses` | [AccessTuple](#fx.ethereum.evm.v1.AccessTuple) | repeated |  |
| `v` | [bytes](#bytes) |  | v defines the signature value |
| `r` | [bytes](#bytes) |  | r defines the signature value |
| `s` | [bytes](#bytes) |  | s define the signature value |






<a name="fx.ethereum.evm.v1.DynamicFeeTx"></a>

### DynamicFeeTx
DynamicFeeTx is the data of EIP-1559 dinamic fee transactions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `chain_id` | [string](#string) |  | destination EVM chain ID |
| `nonce` | [uint64](#uint64) |  | nonce corresponds to the account nonce (transaction sequence). |
| `gas_tip_cap` | [string](#string) |  | gas tip cap defines the max value for the gas tip |
| `gas_fee_cap` | [string](#string) |  | gas fee cap defines the max value for the gas fee |
| `gas` | [uint64](#uint64) |  | gas defines the gas limit defined for the transaction. |
| `to` | [string](#string) |  | hex formatted address of the recipient |
| `value` | [string](#string) |  | value defines the the transaction amount. |
| `data` | [bytes](#bytes) |  | input defines the data payload bytes of the transaction. |
| `accesses` | [AccessTuple](#fx.ethereum.evm.v1.AccessTuple) | repeated |  |
| `v` | [bytes](#bytes) |  | v defines the signature value |
| `r` | [bytes](#bytes) |  | r defines the signature value |
| `s` | [bytes](#bytes) |  | s define the signature value |






<a name="fx.ethereum.evm.v1.ExtensionOptionsEthereumTx"></a>

### ExtensionOptionsEthereumTx







<a name="fx.ethereum.evm.v1.LegacyTx"></a>

### LegacyTx
LegacyTx is the transaction data of regular Ethereum transactions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  | nonce corresponds to the account nonce (transaction sequence). |
| `gas_price` | [string](#string) |  | gas price defines the value for each gas unit |
| `gas` | [uint64](#uint64) |  | gas defines the gas limit defined for the transaction. |
| `to` | [string](#string) |  | hex formatted address of the recipient |
| `value` | [string](#string) |  | value defines the unsigned integer value of the transaction amount. |
| `data` | [bytes](#bytes) |  | input defines the data payload bytes of the transaction. |
| `v` | [bytes](#bytes) |  | v defines the signature value |
| `r` | [bytes](#bytes) |  | r defines the signature value |
| `s` | [bytes](#bytes) |  | s define the signature value |






<a name="fx.ethereum.evm.v1.MsgEthereumTx"></a>

### MsgEthereumTx
MsgEthereumTx encapsulates an Ethereum transaction as an SDK message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [google.protobuf.Any](#google.protobuf.Any) |  | inner transaction data

caches |
| `size` | [double](#double) |  | encoded storage size of the transaction |
| `hash` | [string](#string) |  | transaction hash in hex format |
| `from` | [string](#string) |  | ethereum signer address in hex format. This address value is checked against the address derived from the signature (V, R, S) using the secp256k1 elliptic curve |






<a name="fx.ethereum.evm.v1.MsgEthereumTxResponse"></a>

### MsgEthereumTxResponse
MsgEthereumTxResponse defines the Msg/EthereumTx response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  | ethereum transaction hash in hex format. This hash differs from the Tendermint sha256 hash of the transaction bytes. See https://github.com/tendermint/tendermint/issues/6539 for reference |
| `logs` | [Log](#fx.ethereum.evm.v1.Log) | repeated | logs contains the transaction hash and the proto-compatible ethereum logs. |
| `ret` | [bytes](#bytes) |  | returned data from evm function (result or data supplied with revert opcode) |
| `vm_error` | [string](#string) |  | vm error is the error returned by vm execution |
| `gas_used` | [uint64](#uint64) |  | gas consumed by the transaction |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.ethereum.evm.v1.Msg"></a>

### Msg
Msg defines the evm Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `EthereumTx` | [MsgEthereumTx](#fx.ethereum.evm.v1.MsgEthereumTx) | [MsgEthereumTxResponse](#fx.ethereum.evm.v1.MsgEthereumTxResponse) | EthereumTx defines a method submitting Ethereum transactions. | POST|/evm/v1/ethereum_tx|

 <!-- end services -->



<a name="ethereum/evm/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethereum/evm/v1/query.proto



<a name="fx.ethereum.evm.v1.EstimateGasResponse"></a>

### EstimateGasResponse
EstimateGasResponse defines EstimateGas response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas` | [uint64](#uint64) |  | the estimated gas |






<a name="fx.ethereum.evm.v1.EthCallRequest"></a>

### EthCallRequest
EthCallRequest defines EthCall request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `args` | [bytes](#bytes) |  | same json format as the json rpc api. |
| `gas_cap` | [uint64](#uint64) |  | the default gas cap to be used |






<a name="fx.ethereum.evm.v1.QueryAccountRequest"></a>

### QueryAccountRequest
QueryAccountRequest is the request type for the Query/Account RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the ethereum hex address to query the account for. |






<a name="fx.ethereum.evm.v1.QueryAccountResponse"></a>

### QueryAccountResponse
QueryAccountResponse is the response type for the Query/Account RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balance` | [string](#string) |  | balance is the balance of the EVM denomination. |
| `code_hash` | [string](#string) |  | code hash is the hex-formatted code bytes from the EOA. |
| `nonce` | [uint64](#uint64) |  | nonce is the account's sequence number. |






<a name="fx.ethereum.evm.v1.QueryBalanceRequest"></a>

### QueryBalanceRequest
QueryBalanceRequest is the request type for the Query/Balance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the ethereum hex address to query the balance for. |






<a name="fx.ethereum.evm.v1.QueryBalanceResponse"></a>

### QueryBalanceResponse
QueryBalanceResponse is the response type for the Query/Balance RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `balance` | [string](#string) |  | balance is the balance of the EVM denomination. |






<a name="fx.ethereum.evm.v1.QueryCodeRequest"></a>

### QueryCodeRequest
QueryCodeRequest is the request type for the Query/Code RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the ethereum hex address to query the code for. |






<a name="fx.ethereum.evm.v1.QueryCodeResponse"></a>

### QueryCodeResponse
QueryCodeResponse is the response type for the Query/Code RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code` | [bytes](#bytes) |  | code represents the code bytes from an ethereum address. |






<a name="fx.ethereum.evm.v1.QueryCosmosAccountRequest"></a>

### QueryCosmosAccountRequest
QueryCosmosAccountRequest is the request type for the Query/CosmosAccount RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the ethereum hex address to query the account for. |






<a name="fx.ethereum.evm.v1.QueryCosmosAccountResponse"></a>

### QueryCosmosAccountResponse
QueryCosmosAccountResponse is the response type for the Query/CosmosAccount
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cosmos_address` | [string](#string) |  | cosmos_address is the cosmos address of the account. |
| `sequence` | [uint64](#uint64) |  | sequence is the account's sequence number. |
| `account_number` | [uint64](#uint64) |  | account_number is the account numbert |






<a name="fx.ethereum.evm.v1.QueryModuleEnableRequest"></a>

### QueryModuleEnableRequest
QueryModuleEnableRequest defines the request type for querying the module is
enable.






<a name="fx.ethereum.evm.v1.QueryModuleEnableResponse"></a>

### QueryModuleEnableResponse
QueryModuleEnableResponse returns module is enable.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enable` | [bool](#bool) |  |  |






<a name="fx.ethereum.evm.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest defines the request type for querying x/evm parameters.






<a name="fx.ethereum.evm.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse defines the response type for querying x/evm parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#fx.ethereum.evm.v1.Params) |  | params define the evm module parameters. |






<a name="fx.ethereum.evm.v1.QueryStorageRequest"></a>

### QueryStorageRequest
QueryStorageRequest is the request type for the Query/Storage RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the ethereum hex address to query the storage state for. |
| `key` | [string](#string) |  | key defines the key of the storage state |






<a name="fx.ethereum.evm.v1.QueryStorageResponse"></a>

### QueryStorageResponse
QueryStorageResponse is the response type for the Query/Storage RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [string](#string) |  | key defines the storage state value hash associated with the given key. |






<a name="fx.ethereum.evm.v1.QueryTraceBlockRequest"></a>

### QueryTraceBlockRequest
QueryTraceBlockRequest defines TraceTx request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `txs` | [MsgEthereumTx](#fx.ethereum.evm.v1.MsgEthereumTx) | repeated | txs messages in the block |
| `trace_config` | [TraceConfig](#fx.ethereum.evm.v1.TraceConfig) |  | TraceConfig holds extra parameters to trace functions. |
| `block_number` | [int64](#int64) |  | block number |
| `block_hash` | [string](#string) |  | block hex hash |
| `block_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | block time |






<a name="fx.ethereum.evm.v1.QueryTraceBlockResponse"></a>

### QueryTraceBlockResponse
QueryTraceBlockResponse defines TraceBlock response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  |  |






<a name="fx.ethereum.evm.v1.QueryTraceTxRequest"></a>

### QueryTraceTxRequest
QueryTraceTxRequest defines TraceTx request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msg` | [MsgEthereumTx](#fx.ethereum.evm.v1.MsgEthereumTx) |  | msgEthereumTx for the requested transaction |
| `trace_config` | [TraceConfig](#fx.ethereum.evm.v1.TraceConfig) |  | TraceConfig holds extra parameters to trace functions. |
| `predecessors` | [MsgEthereumTx](#fx.ethereum.evm.v1.MsgEthereumTx) | repeated | the predecessor transactions included in the same block need to be replayed first to get correct context for tracing. |
| `block_number` | [int64](#int64) |  | block number of requested transaction |
| `block_hash` | [string](#string) |  | block hex hash of requested transaction |
| `block_time` | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | block time of requested transaction |






<a name="fx.ethereum.evm.v1.QueryTraceTxResponse"></a>

### QueryTraceTxResponse
QueryTraceTxResponse defines TraceTx response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | response serialized in bytes |






<a name="fx.ethereum.evm.v1.QueryTxLogsRequest"></a>

### QueryTxLogsRequest
QueryTxLogsRequest is the request type for the Query/TxLogs RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  | hash is the ethereum transaction hex hash to query the logs for. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="fx.ethereum.evm.v1.QueryTxLogsResponse"></a>

### QueryTxLogsResponse
QueryTxLogs is the response type for the Query/TxLogs RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `logs` | [Log](#fx.ethereum.evm.v1.Log) | repeated | logs represents the ethereum logs generated from the given transaction. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="fx.ethereum.evm.v1.QueryValidatorAccountRequest"></a>

### QueryValidatorAccountRequest
QueryValidatorAccountRequest is the request type for the
Query/ValidatorAccount RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `cons_address` | [string](#string) |  | cons_address is the validator cons address to query the account for. |






<a name="fx.ethereum.evm.v1.QueryValidatorAccountResponse"></a>

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


<a name="fx.ethereum.evm.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Account` | [QueryAccountRequest](#fx.ethereum.evm.v1.QueryAccountRequest) | [QueryAccountResponse](#fx.ethereum.evm.v1.QueryAccountResponse) | Account queries an Ethereum account. | GET|/ethereum/evm/v1/account/{address}|
| `CosmosAccount` | [QueryCosmosAccountRequest](#fx.ethereum.evm.v1.QueryCosmosAccountRequest) | [QueryCosmosAccountResponse](#fx.ethereum.evm.v1.QueryCosmosAccountResponse) | CosmosAccount queries an Ethereum account's Cosmos Address. | GET|/ethereum/evm/v1/cosmos_account/{address}|
| `ValidatorAccount` | [QueryValidatorAccountRequest](#fx.ethereum.evm.v1.QueryValidatorAccountRequest) | [QueryValidatorAccountResponse](#fx.ethereum.evm.v1.QueryValidatorAccountResponse) | ValidatorAccount queries an Ethereum account's from a validator consensus Address. | GET|/ethereum/evm/v1/validator_account/{cons_address}|
| `Balance` | [QueryBalanceRequest](#fx.ethereum.evm.v1.QueryBalanceRequest) | [QueryBalanceResponse](#fx.ethereum.evm.v1.QueryBalanceResponse) | Balance queries the balance of a the EVM denomination for a single EthAccount. | GET|/ethereum/evm/v1/balances/{address}|
| `Storage` | [QueryStorageRequest](#fx.ethereum.evm.v1.QueryStorageRequest) | [QueryStorageResponse](#fx.ethereum.evm.v1.QueryStorageResponse) | Storage queries the balance of all coins for a single account. | GET|/ethereum/evm/v1/storage/{address}/{key}|
| `Code` | [QueryCodeRequest](#fx.ethereum.evm.v1.QueryCodeRequest) | [QueryCodeResponse](#fx.ethereum.evm.v1.QueryCodeResponse) | Code queries the balance of all coins for a single account. | GET|/ethereum/evm/v1/codes/{address}|
| `Params` | [QueryParamsRequest](#fx.ethereum.evm.v1.QueryParamsRequest) | [QueryParamsResponse](#fx.ethereum.evm.v1.QueryParamsResponse) | Params queries the parameters of x/evm module. | GET|/ethereum/evm/v1/params|
| `EthCall` | [EthCallRequest](#fx.ethereum.evm.v1.EthCallRequest) | [MsgEthereumTxResponse](#fx.ethereum.evm.v1.MsgEthereumTxResponse) | EthCall implements the `eth_call` rpc api | GET|/ethereum/evm/v1/eth_call|
| `EstimateGas` | [EthCallRequest](#fx.ethereum.evm.v1.EthCallRequest) | [EstimateGasResponse](#fx.ethereum.evm.v1.EstimateGasResponse) | EstimateGas implements the `eth_estimateGas` rpc api | GET|/ethereum/evm/v1/estimate_gas|
| `TraceTx` | [QueryTraceTxRequest](#fx.ethereum.evm.v1.QueryTraceTxRequest) | [QueryTraceTxResponse](#fx.ethereum.evm.v1.QueryTraceTxResponse) | TraceTx implements the `debug_traceTransaction` rpc api | GET|/ethereum/evm/v1/trace_tx|
| `TraceBlock` | [QueryTraceBlockRequest](#fx.ethereum.evm.v1.QueryTraceBlockRequest) | [QueryTraceBlockResponse](#fx.ethereum.evm.v1.QueryTraceBlockResponse) | TraceBlock implements the `debug_traceBlockByNumber` and `debug_traceBlockByHash` rpc api | GET|/ethereum/evm/v1/trace_block|

 <!-- end services -->



<a name="ethereum/feemarket/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethereum/feemarket/v1/genesis.proto



<a name="fx.ethereum.feemarket.v1.GenesisState"></a>

### GenesisState
GenesisState defines the feemarket module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#fx.ethereum.feemarket.v1.Params) |  | params defines all the paramaters of the module. |
| `block_gas` | [uint64](#uint64) |  | block gas is the amount of gas used on the last block before the upgrade. Zero by default. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ethereum/feemarket/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethereum/feemarket/v1/query.proto



<a name="fx.ethereum.feemarket.v1.QueryBaseFeeRequest"></a>

### QueryBaseFeeRequest
QueryBaseFeeRequest defines the request type for querying the EIP1559 base
fee.






<a name="fx.ethereum.feemarket.v1.QueryBaseFeeResponse"></a>

### QueryBaseFeeResponse
BaseFeeResponse returns the EIP1559 base fee.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `base_fee` | [string](#string) |  |  |






<a name="fx.ethereum.feemarket.v1.QueryBlockGasRequest"></a>

### QueryBlockGasRequest
QueryBlockGasRequest defines the request type for querying the EIP1559 base
fee.






<a name="fx.ethereum.feemarket.v1.QueryBlockGasResponse"></a>

### QueryBlockGasResponse
QueryBlockGasResponse returns block gas used for a given height.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `gas` | [int64](#int64) |  |  |






<a name="fx.ethereum.feemarket.v1.QueryModuleEnableRequest"></a>

### QueryModuleEnableRequest
QueryModuleEnableRequest defines the request type for querying the module is
enable.






<a name="fx.ethereum.feemarket.v1.QueryModuleEnableResponse"></a>

### QueryModuleEnableResponse
QueryModuleEnableResponse returns module is enable.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `enable` | [bool](#bool) |  |  |






<a name="fx.ethereum.feemarket.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest defines the request type for querying x/evm parameters.






<a name="fx.ethereum.feemarket.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse defines the response type for querying x/evm parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#fx.ethereum.feemarket.v1.Params) |  | params define the evm module parameters. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.ethereum.feemarket.v1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Params` | [QueryParamsRequest](#fx.ethereum.feemarket.v1.QueryParamsRequest) | [QueryParamsResponse](#fx.ethereum.feemarket.v1.QueryParamsResponse) | Params queries the parameters of x/feemarket module. | GET|/ethereum/feemarket/evm/v1/params|
| `BaseFee` | [QueryBaseFeeRequest](#fx.ethereum.feemarket.v1.QueryBaseFeeRequest) | [QueryBaseFeeResponse](#fx.ethereum.feemarket.v1.QueryBaseFeeResponse) | BaseFee queries the base fee of the parent block of the current block. | GET|/ethereum/feemarket/evm/v1/base_fee|
| `BlockGas` | [QueryBlockGasRequest](#fx.ethereum.feemarket.v1.QueryBlockGasRequest) | [QueryBlockGasResponse](#fx.ethereum.feemarket.v1.QueryBlockGasResponse) | BlockGas queries the gas used at a given block height | GET|/ethereum/feemarket/evm/v1/block_gas|

 <!-- end services -->



<a name="ethereum/types/v1/web3.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ethereum/types/v1/web3.proto



<a name="fx.ethereum.types.v1.ExtensionOptionsWeb3Tx"></a>

### ExtensionOptionsWeb3Tx



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `typed_data_chain_id` | [uint64](#uint64) |  | typed data chain id used only in EIP712 Domain and should match Ethereum network ID in a Web3 provider (e.g. Metamask). |
| `fee_payer` | [string](#string) |  | fee payer is an account address for the fee payer. It will be validated during EIP712 signature checking. |
| `fee_payer_sig` | [bytes](#bytes) |  | fee payer sig is a signature data from the fee paying account, allows to perform fee delegation when using EIP712 Domain. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="gravity/v1/attestation.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## gravity/v1/attestation.proto



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



<a name="gravity/v1/batch.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## gravity/v1/batch.proto



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



<a name="gravity/v1/ethereum_signer.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## gravity/v1/ethereum_signer.proto


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



<a name="gravity/v1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## gravity/v1/types.proto



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



<a name="gravity/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## gravity/v1/tx.proto



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
| `fee_receive` | [string](#string) |  |  |
| `base_fee` | [string](#string) |  |  |






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



<a name="gravity/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## gravity/v1/genesis.proto



<a name="fx.gravity.v1.GenesisState"></a>

### GenesisState
GenesisState struct


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#fx.gravity.v1.Params) |  |  |
| `last_observed_nonce` | [uint64](#uint64) |  |  |
| `valsets` | [Valset](#fx.gravity.v1.Valset) | repeated |  |
| `valset_confirms` | [MsgValsetConfirm](#fx.gravity.v1.MsgValsetConfirm) | repeated |  |
| `batches` | [OutgoingTxBatch](#fx.gravity.v1.OutgoingTxBatch) | repeated |  |
| `batch_confirms` | [MsgConfirmBatch](#fx.gravity.v1.MsgConfirmBatch) | repeated |  |
| `attestations` | [Attestation](#fx.gravity.v1.Attestation) | repeated |  |
| `delegate_keys` | [MsgSetOrchestratorAddress](#fx.gravity.v1.MsgSetOrchestratorAddress) | repeated |  |
| `erc20_to_denoms` | [ERC20ToDenom](#fx.gravity.v1.ERC20ToDenom) | repeated |  |
| `unbatched_transfers` | [OutgoingTransferTx](#fx.gravity.v1.OutgoingTransferTx) | repeated |  |
| `module_coins` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






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



<a name="gravity/v1/pool.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## gravity/v1/pool.proto



<a name="fx.gravity.v1.BatchFees"></a>

### BatchFees



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_contract` | [string](#string) |  |  |
| `total_fees` | [string](#string) |  |  |
| `total_txs` | [uint64](#uint64) |  |  |
| `total_amount` | [string](#string) |  |  |






<a name="fx.gravity.v1.IDSet"></a>

### IDSet
IDSet represents a set of IDs


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `ids` | [uint64](#uint64) | repeated |  |






<a name="fx.gravity.v1.MinBatchFee"></a>

### MinBatchFee



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `token_contract` | [string](#string) |  |  |
| `baseFee` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="gravity/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## gravity/v1/query.proto



<a name="fx.gravity.v1.QueryBatchConfirmRequest"></a>

### QueryBatchConfirmRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `nonce` | [uint64](#uint64) |  |  |
| `token_contract` | [string](#string) |  |  |
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
| `token_contract` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryBatchConfirmsResponse"></a>

### QueryBatchConfirmsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `confirms` | [MsgConfirmBatch](#fx.gravity.v1.MsgConfirmBatch) | repeated |  |






<a name="fx.gravity.v1.QueryBatchFeeRequest"></a>

### QueryBatchFeeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `minBatchFees` | [MinBatchFee](#fx.gravity.v1.MinBatchFee) | repeated |  |






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
| `token_contract` | [string](#string) |  |  |






<a name="fx.gravity.v1.QueryBatchRequestByNonceResponse"></a>

### QueryBatchRequestByNonceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `batch` | [OutgoingTxBatch](#fx.gravity.v1.OutgoingTxBatch) |  |  |






<a name="fx.gravity.v1.QueryBridgeTokensRequest"></a>

### QueryBridgeTokensRequest







<a name="fx.gravity.v1.QueryBridgeTokensResponse"></a>

### QueryBridgeTokensResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `bridge_tokens` | [ERC20ToDenom](#fx.gravity.v1.ERC20ToDenom) | repeated |  |






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






<a name="fx.gravity.v1.QueryLastObservedBlockHeightRequest"></a>

### QueryLastObservedBlockHeightRequest







<a name="fx.gravity.v1.QueryLastObservedBlockHeightResponse"></a>

### QueryLastObservedBlockHeightResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `eth_block_height` | [uint64](#uint64) |  |  |
| `block_height` | [uint64](#uint64) |  |  |






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






<a name="fx.gravity.v1.QueryProjectedBatchTimeoutHeightRequest"></a>

### QueryProjectedBatchTimeoutHeightRequest







<a name="fx.gravity.v1.QueryProjectedBatchTimeoutHeightResponse"></a>

### QueryProjectedBatchTimeoutHeightResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `timeout_height` | [uint64](#uint64) |  |  |






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
| `ValsetRequest` | [QueryValsetRequestRequest](#fx.gravity.v1.QueryValsetRequestRequest) | [QueryValsetRequestResponse](#fx.gravity.v1.QueryValsetRequestResponse) |  | GET|/gravity/v1beta/valset/request|
| `ValsetConfirm` | [QueryValsetConfirmRequest](#fx.gravity.v1.QueryValsetConfirmRequest) | [QueryValsetConfirmResponse](#fx.gravity.v1.QueryValsetConfirmResponse) |  | GET|/gravity/v1beta/valset/confirm|
| `ValsetConfirmsByNonce` | [QueryValsetConfirmsByNonceRequest](#fx.gravity.v1.QueryValsetConfirmsByNonceRequest) | [QueryValsetConfirmsByNonceResponse](#fx.gravity.v1.QueryValsetConfirmsByNonceResponse) |  | GET|/gravity/v1beta/valset/confirms|
| `LastValsetRequests` | [QueryLastValsetRequestsRequest](#fx.gravity.v1.QueryLastValsetRequestsRequest) | [QueryLastValsetRequestsResponse](#fx.gravity.v1.QueryLastValsetRequestsResponse) |  | GET|/gravity/v1beta/valset/requests|
| `LastPendingValsetRequestByAddr` | [QueryLastPendingValsetRequestByAddrRequest](#fx.gravity.v1.QueryLastPendingValsetRequestByAddrRequest) | [QueryLastPendingValsetRequestByAddrResponse](#fx.gravity.v1.QueryLastPendingValsetRequestByAddrResponse) |  | GET|/gravity/v1beta/valset/last|
| `LastPendingBatchRequestByAddr` | [QueryLastPendingBatchRequestByAddrRequest](#fx.gravity.v1.QueryLastPendingBatchRequestByAddrRequest) | [QueryLastPendingBatchRequestByAddrResponse](#fx.gravity.v1.QueryLastPendingBatchRequestByAddrResponse) |  | GET|/gravity/v1beta/batch/last|
| `LastEventNonceByAddr` | [QueryLastEventNonceByAddrRequest](#fx.gravity.v1.QueryLastEventNonceByAddrRequest) | [QueryLastEventNonceByAddrResponse](#fx.gravity.v1.QueryLastEventNonceByAddrResponse) |  | GET|/gravity/v1beta/oracle/event_nonce/{address}|
| `LastEventBlockHeightByAddr` | [QueryLastEventBlockHeightByAddrRequest](#fx.gravity.v1.QueryLastEventBlockHeightByAddrRequest) | [QueryLastEventBlockHeightByAddrResponse](#fx.gravity.v1.QueryLastEventBlockHeightByAddrResponse) |  | GET|/gravity/v1beta/oracle/event/block_height/{address}|
| `BatchFees` | [QueryBatchFeeRequest](#fx.gravity.v1.QueryBatchFeeRequest) | [QueryBatchFeeResponse](#fx.gravity.v1.QueryBatchFeeResponse) |  | GET|/gravity/v1beta/batch_fees|
| `LastObservedBlockHeight` | [QueryLastObservedBlockHeightRequest](#fx.gravity.v1.QueryLastObservedBlockHeightRequest) | [QueryLastObservedBlockHeightResponse](#fx.gravity.v1.QueryLastObservedBlockHeightResponse) |  | GET|/gravity/v1beta/observed/block_height|
| `OutgoingTxBatches` | [QueryOutgoingTxBatchesRequest](#fx.gravity.v1.QueryOutgoingTxBatchesRequest) | [QueryOutgoingTxBatchesResponse](#fx.gravity.v1.QueryOutgoingTxBatchesResponse) |  | GET|/gravity/v1beta/batch/outgoing_tx|
| `BatchRequestByNonce` | [QueryBatchRequestByNonceRequest](#fx.gravity.v1.QueryBatchRequestByNonceRequest) | [QueryBatchRequestByNonceResponse](#fx.gravity.v1.QueryBatchRequestByNonceResponse) |  | GET|/gravity/v1beta/batch/request|
| `BatchConfirm` | [QueryBatchConfirmRequest](#fx.gravity.v1.QueryBatchConfirmRequest) | [QueryBatchConfirmResponse](#fx.gravity.v1.QueryBatchConfirmResponse) |  | GET|/gravity/v1beta/batch/confirm|
| `BatchConfirms` | [QueryBatchConfirmsRequest](#fx.gravity.v1.QueryBatchConfirmsRequest) | [QueryBatchConfirmsResponse](#fx.gravity.v1.QueryBatchConfirmsResponse) |  | GET|/gravity/v1beta/batch/confirms|
| `ERC20ToDenom` | [QueryERC20ToDenomRequest](#fx.gravity.v1.QueryERC20ToDenomRequest) | [QueryERC20ToDenomResponse](#fx.gravity.v1.QueryERC20ToDenomResponse) |  | GET|/gravity/v1beta/denom|
| `DenomToERC20` | [QueryDenomToERC20Request](#fx.gravity.v1.QueryDenomToERC20Request) | [QueryDenomToERC20Response](#fx.gravity.v1.QueryDenomToERC20Response) |  | GET|/gravity/v1beta/erc20|
| `GetDelegateKeyByValidator` | [QueryDelegateKeyByValidatorRequest](#fx.gravity.v1.QueryDelegateKeyByValidatorRequest) | [QueryDelegateKeyByValidatorResponse](#fx.gravity.v1.QueryDelegateKeyByValidatorResponse) |  | GET|/gravity/v1beta/delegate_key_by_validator|
| `GetDelegateKeyByEth` | [QueryDelegateKeyByEthRequest](#fx.gravity.v1.QueryDelegateKeyByEthRequest) | [QueryDelegateKeyByEthResponse](#fx.gravity.v1.QueryDelegateKeyByEthResponse) |  | GET|/gravity/v1beta/delegate_key_by_eth|
| `GetDelegateKeyByOrchestrator` | [QueryDelegateKeyByOrchestratorRequest](#fx.gravity.v1.QueryDelegateKeyByOrchestratorRequest) | [QueryDelegateKeyByOrchestratorResponse](#fx.gravity.v1.QueryDelegateKeyByOrchestratorResponse) |  | GET|/gravity/v1beta/delegate_key_by_orchestrator|
| `GetPendingSendToEth` | [QueryPendingSendToEthRequest](#fx.gravity.v1.QueryPendingSendToEthRequest) | [QueryPendingSendToEthResponse](#fx.gravity.v1.QueryPendingSendToEthResponse) |  | GET|/gravity/v1beta/pending_send_to_eth|
| `GetIbcSequenceHeightByChannel` | [QueryIbcSequenceHeightRequest](#fx.gravity.v1.QueryIbcSequenceHeightRequest) | [QueryIbcSequenceHeightResponse](#fx.gravity.v1.QueryIbcSequenceHeightResponse) |  | GET|/gravity/v1beta/ibc_sequence_height|
| `ProjectedBatchTimeoutHeight` | [QueryProjectedBatchTimeoutHeightRequest](#fx.gravity.v1.QueryProjectedBatchTimeoutHeightRequest) | [QueryProjectedBatchTimeoutHeightResponse](#fx.gravity.v1.QueryProjectedBatchTimeoutHeightResponse) |  | GET|/gravity/v1beta1/projected_batch_timeout|
| `BridgeTokens` | [QueryBridgeTokensRequest](#fx.gravity.v1.QueryBridgeTokensRequest) | [QueryBridgeTokensResponse](#fx.gravity.v1.QueryBridgeTokensResponse) |  | GET|/gravity/v1beta1/bridge_tokens|

 <!-- end services -->



<a name="ibc/applications/transfer/v1/transfer.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/transfer/v1/transfer.proto



<a name="fx.ibc.applications.transfer.v1.DenomTrace"></a>

### DenomTrace
DenomTrace contains the base denomination for ICS20 fungible tokens and the
source tracing information path.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `path` | [string](#string) |  | path defines the chain of port/channel identifiers used for tracing the source of the fungible token. |
| `base_denom` | [string](#string) |  | base denomination of the relayed fungible token. |






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






<a name="fx.ibc.applications.transfer.v1.Params"></a>

### Params
Params defines the set of IBC transfer parameters.
NOTE: To prevent a single token from being transferred, set the
TransfersEnabled parameter to true and then set the bank module's SendEnabled
parameter for the denomination to false.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `send_enabled` | [bool](#bool) |  | send_enabled enables or disables all cross-chain token transfers from this chain. |
| `receive_enabled` | [bool](#bool) |  | receive_enabled enables or disables all cross-chain token transfers to this chain. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/transfer/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/transfer/v1/genesis.proto



<a name="fx.ibc.applications.transfer.v1.GenesisState"></a>

### GenesisState
GenesisState defines the ibc-transfer genesis state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `port_id` | [string](#string) |  |  |
| `denom_traces` | [DenomTrace](#fx.ibc.applications.transfer.v1.DenomTrace) | repeated |  |
| `params` | [Params](#fx.ibc.applications.transfer.v1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="ibc/applications/transfer/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/transfer/v1/query.proto



<a name="fx.ibc.applications.transfer.v1.QueryDenomTraceRequest"></a>

### QueryDenomTraceRequest
QueryDenomTraceRequest is the request type for the Query/DenomTrace RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `hash` | [string](#string) |  | hash (in hex format) of the denomination trace information. |






<a name="fx.ibc.applications.transfer.v1.QueryDenomTraceResponse"></a>

### QueryDenomTraceResponse
QueryDenomTraceResponse is the response type for the Query/DenomTrace RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom_trace` | [DenomTrace](#fx.ibc.applications.transfer.v1.DenomTrace) |  | denom_trace returns the requested denomination trace information. |






<a name="fx.ibc.applications.transfer.v1.QueryDenomTracesRequest"></a>

### QueryDenomTracesRequest
QueryConnectionsRequest is the request type for the Query/DenomTraces RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="fx.ibc.applications.transfer.v1.QueryDenomTracesResponse"></a>

### QueryDenomTracesResponse
QueryConnectionsResponse is the response type for the Query/DenomTraces RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `denom_traces` | [DenomTrace](#fx.ibc.applications.transfer.v1.DenomTrace) | repeated | denom_traces returns all denominations trace information. |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<a name="fx.ibc.applications.transfer.v1.QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="fx.ibc.applications.transfer.v1.QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#fx.ibc.applications.transfer.v1.Params) |  | params defines the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.ibc.applications.transfer.v1.Query"></a>

### Query
Query provides defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `DenomTrace` | [QueryDenomTraceRequest](#fx.ibc.applications.transfer.v1.QueryDenomTraceRequest) | [QueryDenomTraceResponse](#fx.ibc.applications.transfer.v1.QueryDenomTraceResponse) | DenomTrace queries a denomination trace information. | GET|/ibc/applications/transfer/v1beta1/denom_traces/{hash}|
| `DenomTraces` | [QueryDenomTracesRequest](#fx.ibc.applications.transfer.v1.QueryDenomTracesRequest) | [QueryDenomTracesResponse](#fx.ibc.applications.transfer.v1.QueryDenomTracesResponse) | DenomTraces queries all denomination traces. | GET|/ibc/applications/transfer/v1beta1/denom_traces|
| `Params` | [QueryParamsRequest](#fx.ibc.applications.transfer.v1.QueryParamsRequest) | [QueryParamsResponse](#fx.ibc.applications.transfer.v1.QueryParamsResponse) | Params queries all parameters of the ibc-transfer module. | GET|/ibc/applications/transfer/v1beta1/params|

 <!-- end services -->



<a name="ibc/applications/transfer/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## ibc/applications/transfer/v1/tx.proto



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






<a name="fx.ibc.applications.transfer.v1.MsgTransferResponse"></a>

### MsgTransferResponse
MsgTransferResponse defines the Msg/Transfer response type.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="fx.ibc.applications.transfer.v1.Msg"></a>

### Msg
Msg defines the ibc/transfer Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Transfer` | [MsgTransfer](#fx.ibc.applications.transfer.v1.MsgTransfer) | [MsgTransferResponse](#fx.ibc.applications.transfer.v1.MsgTransferResponse) | Transfer defines a rpc handler method for MsgTransfer. | |

 <!-- end services -->



<a name="migrate/v1/migrate.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## migrate/v1/migrate.proto



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



<a name="migrate/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## migrate/v1/query.proto



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
| `MigrateRecord` | [QueryMigrateRecordRequest](#fx.migrate.v1.QueryMigrateRecordRequest) | [QueryMigrateRecordResponse](#fx.migrate.v1.QueryMigrateRecordResponse) | DenomTrace queries a denomination trace information. | GET|/migrate/v1/record/{address}|

 <!-- end services -->



<a name="migrate/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## migrate/v1/tx.proto



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



<a name="other/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## other/query.proto



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
