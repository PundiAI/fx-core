// nolint:staticcheck
package keeper

import (
	"context"

	fxtypes "github.com/functionx/fx-core/v7/types"
	crosschainkeeper "github.com/functionx/fx-core/v7/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
	"github.com/functionx/fx-core/v7/x/gravity/types"
)

type queryServer struct {
	crosschainkeeper.Keeper
}

func NewQueryServerImpl(keeper crosschainkeeper.Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

var _ types.QueryServer = queryServer{}

func (k queryServer) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	response, err := k.Keeper.Params(c, &crosschaintypes.QueryParamsRequest{ChainName: ethtypes.ModuleName})
	if err != nil {
		return nil, err
	}
	return &types.QueryParamsResponse{Params: types.Params{
		GravityId:                      response.Params.GravityId,
		BridgeChainId:                  1,
		SignedValsetsWindow:            response.Params.SignedWindow,
		SignedBatchesWindow:            response.Params.SignedWindow,
		SignedClaimsWindow:             response.Params.SignedWindow,
		TargetBatchTimeout:             response.Params.ExternalBatchTimeout,
		AverageBlockTime:               response.Params.AverageBlockTime,
		AverageEthBlockTime:            response.Params.AverageExternalBlockTime,
		SlashFractionValset:            response.Params.SlashFraction,
		SlashFractionBatch:             response.Params.SlashFraction,
		SlashFractionClaim:             response.Params.SlashFraction,
		SlashFractionConflictingClaim:  response.Params.SlashFraction,
		UnbondSlashingValsetsWindow:    response.Params.SignedWindow,
		IbcTransferTimeoutHeight:       response.Params.IbcTransferTimeoutHeight,
		ValsetUpdatePowerChangePercent: response.Params.OracleSetUpdatePowerChangePercent,
	}}, nil
}

func (k queryServer) CurrentValset(c context.Context, _ *types.QueryCurrentValsetRequest) (*types.QueryCurrentValsetResponse, error) {
	response, err := k.Keeper.CurrentOracleSet(c, &crosschaintypes.QueryCurrentOracleSetRequest{ChainName: ethtypes.ModuleName})
	if err != nil {
		return nil, err
	}
	valset := &types.Valset{
		Nonce:   response.OracleSet.Nonce,
		Members: make([]*types.BridgeValidator, len(response.OracleSet.Members)),
		Height:  response.OracleSet.Height,
	}
	for i := 0; i < len(response.OracleSet.Members); i++ {
		valset.Members[i] = &types.BridgeValidator{
			Power:      response.OracleSet.Members[i].Power,
			EthAddress: response.OracleSet.Members[i].ExternalAddress,
		}
	}
	return &types.QueryCurrentValsetResponse{Valset: valset}, nil
}

func (k queryServer) ValsetRequest(c context.Context, req *types.QueryValsetRequestRequest) (*types.QueryValsetRequestResponse, error) {
	response, err := k.Keeper.OracleSetRequest(c, &crosschaintypes.QueryOracleSetRequestRequest{
		ChainName: ethtypes.ModuleName,
		Nonce:     req.Nonce,
	})
	if err != nil {
		return nil, err
	}
	valset := &types.Valset{
		Nonce:   response.OracleSet.Nonce,
		Members: make([]*types.BridgeValidator, len(response.OracleSet.Members)),
		Height:  response.OracleSet.Height,
	}
	for i := 0; i < len(response.OracleSet.Members); i++ {
		valset.Members[i] = &types.BridgeValidator{
			Power:      response.OracleSet.Members[i].Power,
			EthAddress: response.OracleSet.Members[i].ExternalAddress,
		}
	}
	return &types.QueryValsetRequestResponse{Valset: valset}, nil
}

func (k queryServer) ValsetConfirm(c context.Context, req *types.QueryValsetConfirmRequest) (*types.QueryValsetConfirmResponse, error) {
	response, err := k.Keeper.OracleSetConfirm(c, &crosschaintypes.QueryOracleSetConfirmRequest{
		ChainName:      ethtypes.ModuleName,
		BridgerAddress: req.Address,
		Nonce:          req.Nonce,
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryValsetConfirmResponse{Confirm: &types.MsgValsetConfirm{
		Nonce:        response.Confirm.Nonce,
		Orchestrator: response.Confirm.BridgerAddress,
		EthAddress:   response.Confirm.ExternalAddress,
		Signature:    response.Confirm.Signature,
	}}, nil
}

func (k queryServer) ValsetConfirmsByNonce(c context.Context, req *types.QueryValsetConfirmsByNonceRequest) (*types.QueryValsetConfirmsByNonceResponse, error) {
	response, err := k.Keeper.OracleSetConfirmsByNonce(c, &crosschaintypes.QueryOracleSetConfirmsByNonceRequest{
		ChainName: ethtypes.ModuleName,
		Nonce:     req.Nonce,
	})
	if err != nil {
		return nil, err
	}
	confirms := make([]*types.MsgValsetConfirm, len(response.Confirms))
	for i := 0; i < len(response.Confirms); i++ {
		confirms[i] = &types.MsgValsetConfirm{
			Nonce:        response.Confirms[i].Nonce,
			Orchestrator: response.Confirms[i].BridgerAddress,
			EthAddress:   response.Confirms[i].ExternalAddress,
			Signature:    response.Confirms[i].Signature,
		}
	}
	return &types.QueryValsetConfirmsByNonceResponse{Confirms: confirms}, nil
}

func (k queryServer) LastValsetRequests(c context.Context, _ *types.QueryLastValsetRequestsRequest) (*types.QueryLastValsetRequestsResponse, error) {
	response, err := k.Keeper.LastOracleSetRequests(c, &crosschaintypes.QueryLastOracleSetRequestsRequest{ChainName: ethtypes.ModuleName})
	if err != nil {
		return nil, err
	}
	valsets := make([]*types.Valset, len(response.OracleSets))
	for i := 0; i < len(response.OracleSets); i++ {
		valsets[i] = &types.Valset{
			Nonce:   response.OracleSets[i].Nonce,
			Members: make([]*types.BridgeValidator, len(response.OracleSets[i].Members)),
			Height:  response.OracleSets[i].Height,
		}
		for j := 0; j < len(response.OracleSets[i].Members); j++ {
			valsets[i].Members[j] = &types.BridgeValidator{
				Power:      response.OracleSets[i].Members[j].Power,
				EthAddress: response.OracleSets[i].Members[j].ExternalAddress,
			}
		}
	}
	return &types.QueryLastValsetRequestsResponse{Valsets: valsets}, nil
}

func (k queryServer) LastPendingValsetRequestByAddr(c context.Context, req *types.QueryLastPendingValsetRequestByAddrRequest) (*types.QueryLastPendingValsetRequestByAddrResponse, error) {
	response, err := k.Keeper.LastPendingOracleSetRequestByAddr(c, &crosschaintypes.QueryLastPendingOracleSetRequestByAddrRequest{
		ChainName:      ethtypes.ModuleName,
		BridgerAddress: req.Address,
	})
	if err != nil {
		return nil, err
	}
	valsets := make([]*types.Valset, len(response.OracleSets))
	for i := 0; i < len(response.OracleSets); i++ {
		valsets[i] = &types.Valset{
			Nonce:   response.OracleSets[i].Nonce,
			Members: make([]*types.BridgeValidator, len(response.OracleSets[i].Members)),
			Height:  response.OracleSets[i].Height,
		}
		for j := 0; j < len(response.OracleSets[i].Members); j++ {
			valsets[i].Members[j] = &types.BridgeValidator{
				Power:      response.OracleSets[i].Members[j].Power,
				EthAddress: response.OracleSets[i].Members[j].ExternalAddress,
			}
		}
	}

	return &types.QueryLastPendingValsetRequestByAddrResponse{Valsets: valsets}, nil
}

func (k queryServer) BatchFees(c context.Context, req *types.QueryBatchFeeRequest) (*types.QueryBatchFeeResponse, error) {
	minBatchFees := make([]crosschaintypes.MinBatchFee, len(req.MinBatchFees))
	for i := 0; i < len(req.MinBatchFees); i++ {
		minBatchFees[i] = crosschaintypes.MinBatchFee{
			TokenContract: req.MinBatchFees[i].TokenContract,
			BaseFee:       req.MinBatchFees[i].BaseFee,
		}
	}
	response, err := k.Keeper.BatchFees(c, &crosschaintypes.QueryBatchFeeRequest{
		ChainName:    ethtypes.ModuleName,
		MinBatchFees: minBatchFees,
	})
	if err != nil {
		return nil, err
	}
	batchFees := make([]*types.BatchFees, len(response.BatchFees))
	for i := 0; i < len(response.BatchFees); i++ {
		batchFees[i] = &types.BatchFees{
			TokenContract: response.BatchFees[i].TokenContract,
			TotalFees:     response.BatchFees[i].TotalFees,
			TotalTxs:      response.BatchFees[i].TotalTxs,
			TotalAmount:   response.BatchFees[i].TotalAmount,
		}
	}
	return &types.QueryBatchFeeResponse{BatchFees: batchFees}, nil
}

func (k queryServer) LastPendingBatchRequestByAddr(c context.Context, req *types.QueryLastPendingBatchRequestByAddrRequest) (*types.QueryLastPendingBatchRequestByAddrResponse, error) {
	response, err := k.Keeper.LastPendingBatchRequestByAddr(c, &crosschaintypes.QueryLastPendingBatchRequestByAddrRequest{
		ChainName:      ethtypes.ModuleName,
		BridgerAddress: req.Address,
	})
	if err != nil {
		return nil, err
	}

	outgoingTxBatch := &types.OutgoingTxBatch{
		BatchNonce:    response.Batch.BatchNonce,
		BatchTimeout:  response.Batch.BatchTimeout,
		Transactions:  make([]*types.OutgoingTransferTx, len(response.Batch.Transactions)),
		TokenContract: response.Batch.TokenContract,
		Block:         response.Batch.Block,
		FeeReceive:    response.Batch.FeeReceive,
	}
	for i := 0; i < len(response.Batch.Transactions); i++ {
		outgoingTxBatch.Transactions[i] = &types.OutgoingTransferTx{
			Id:          response.Batch.Transactions[i].Id,
			Sender:      response.Batch.Transactions[i].Sender,
			DestAddress: response.Batch.Transactions[i].DestAddress,
			Erc20Token: &types.ERC20Token{
				Contract: response.Batch.Transactions[i].Token.Contract,
				Amount:   response.Batch.Transactions[i].Token.Amount,
			},
			Erc20Fee: &types.ERC20Token{
				Contract: response.Batch.Transactions[i].Fee.Contract,
				Amount:   response.Batch.Transactions[i].Fee.Amount,
			},
		}
	}
	return &types.QueryLastPendingBatchRequestByAddrResponse{Batch: outgoingTxBatch}, nil
}

func (k queryServer) OutgoingTxBatches(c context.Context, _ *types.QueryOutgoingTxBatchesRequest) (*types.QueryOutgoingTxBatchesResponse, error) {
	response, err := k.Keeper.OutgoingTxBatches(c, &crosschaintypes.QueryOutgoingTxBatchesRequest{ChainName: ethtypes.ModuleName})
	if err != nil {
		return nil, err
	}
	batches := make([]*types.OutgoingTxBatch, len(response.Batches))
	for i := 0; i < len(response.Batches); i++ {
		batches[i] = &types.OutgoingTxBatch{
			BatchNonce:    response.Batches[i].BatchNonce,
			BatchTimeout:  response.Batches[i].BatchTimeout,
			Transactions:  make([]*types.OutgoingTransferTx, len(response.Batches[i].Transactions)),
			TokenContract: response.Batches[i].TokenContract,
			Block:         response.Batches[i].Block,
			FeeReceive:    response.Batches[i].FeeReceive,
		}
		for j := 0; j < len(response.Batches[i].Transactions); j++ {
			batches[i].Transactions[j] = &types.OutgoingTransferTx{
				Id:          response.Batches[i].Transactions[j].Id,
				Sender:      response.Batches[i].Transactions[j].Sender,
				DestAddress: response.Batches[i].Transactions[j].DestAddress,
				Erc20Token: &types.ERC20Token{
					Contract: response.Batches[i].Transactions[j].Token.Contract,
					Amount:   response.Batches[i].Transactions[j].Token.Amount,
				},
				Erc20Fee: &types.ERC20Token{
					Contract: response.Batches[i].Transactions[j].Fee.Contract,
					Amount:   response.Batches[i].Transactions[j].Fee.Amount,
				},
			}
		}
	}
	return &types.QueryOutgoingTxBatchesResponse{Batches: batches}, nil
}

func (k queryServer) BatchRequestByNonce(c context.Context, req *types.QueryBatchRequestByNonceRequest) (*types.QueryBatchRequestByNonceResponse, error) {
	response, err := k.Keeper.BatchRequestByNonce(c, &crosschaintypes.QueryBatchRequestByNonceRequest{
		ChainName:     ethtypes.ModuleName,
		TokenContract: req.TokenContract,
		Nonce:         req.Nonce,
	})
	if err != nil {
		return nil, err
	}
	outgoingTxBatch := &types.OutgoingTxBatch{
		BatchNonce:    response.Batch.BatchNonce,
		BatchTimeout:  response.Batch.BatchTimeout,
		Transactions:  make([]*types.OutgoingTransferTx, len(response.Batch.Transactions)),
		TokenContract: response.Batch.TokenContract,
		Block:         response.Batch.Block,
		FeeReceive:    response.Batch.FeeReceive,
	}
	for i := 0; i < len(response.Batch.Transactions); i++ {
		outgoingTxBatch.Transactions[i] = &types.OutgoingTransferTx{
			Id:          response.Batch.Transactions[i].Id,
			Sender:      response.Batch.Transactions[i].Sender,
			DestAddress: response.Batch.Transactions[i].DestAddress,
			Erc20Token: &types.ERC20Token{
				Contract: response.Batch.Transactions[i].Token.Contract,
				Amount:   response.Batch.Transactions[i].Token.Amount,
			},
			Erc20Fee: &types.ERC20Token{
				Contract: response.Batch.Transactions[i].Fee.Contract,
				Amount:   response.Batch.Transactions[i].Fee.Amount,
			},
		}
	}
	return &types.QueryBatchRequestByNonceResponse{Batch: outgoingTxBatch}, nil
}

func (k queryServer) BatchConfirm(c context.Context, req *types.QueryBatchConfirmRequest) (*types.QueryBatchConfirmResponse, error) {
	response, err := k.Keeper.BatchConfirm(c, &crosschaintypes.QueryBatchConfirmRequest{
		ChainName:      ethtypes.ModuleName,
		TokenContract:  req.TokenContract,
		BridgerAddress: req.Address,
		Nonce:          req.Nonce,
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryBatchConfirmResponse{Confirm: &types.MsgConfirmBatch{
		Nonce:         response.Confirm.Nonce,
		TokenContract: response.Confirm.TokenContract,
		EthSigner:     response.Confirm.ExternalAddress,
		Orchestrator:  response.Confirm.BridgerAddress,
		Signature:     response.Confirm.Signature,
	}}, nil
}

func (k queryServer) BatchConfirms(c context.Context, req *types.QueryBatchConfirmsRequest) (*types.QueryBatchConfirmsResponse, error) {
	response, err := k.Keeper.BatchConfirms(c, &crosschaintypes.QueryBatchConfirmsRequest{
		ChainName:     ethtypes.ModuleName,
		TokenContract: req.TokenContract,
		Nonce:         req.Nonce,
	})
	if err != nil {
		return nil, err
	}
	confirms := make([]*types.MsgConfirmBatch, len(response.Confirms))
	for i := 0; i < len(response.Confirms); i++ {
		confirms[i] = &types.MsgConfirmBatch{
			Nonce:         response.Confirms[i].Nonce,
			TokenContract: response.Confirms[i].TokenContract,
			EthSigner:     response.Confirms[i].ExternalAddress,
			Orchestrator:  response.Confirms[i].BridgerAddress,
			Signature:     response.Confirms[i].Signature,
		}
	}
	return &types.QueryBatchConfirmsResponse{Confirms: confirms}, nil
}

func (k queryServer) LastEventNonceByAddr(c context.Context, req *types.QueryLastEventNonceByAddrRequest) (*types.QueryLastEventNonceByAddrResponse, error) {
	response, err := k.Keeper.LastEventNonceByAddr(c, &crosschaintypes.QueryLastEventNonceByAddrRequest{
		ChainName:      ethtypes.ModuleName,
		BridgerAddress: req.Address,
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryLastEventNonceByAddrResponse{EventNonce: response.EventNonce}, nil
}

func (k queryServer) DenomToERC20(c context.Context, req *types.QueryDenomToERC20Request) (*types.QueryDenomToERC20Response, error) {
	response, err := k.Keeper.DenomToToken(c, &crosschaintypes.QueryDenomToTokenRequest{
		ChainName: ethtypes.ModuleName,
		Denom:     req.Denom,
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryDenomToERC20Response{Erc20: response.Token, FxOriginated: req.Denom == fxtypes.DefaultDenom}, err
}

func (k queryServer) ERC20ToDenom(c context.Context, req *types.QueryERC20ToDenomRequest) (*types.QueryERC20ToDenomResponse, error) {
	response, err := k.Keeper.TokenToDenom(c, &crosschaintypes.QueryTokenToDenomRequest{
		ChainName: ethtypes.ModuleName,
		Token:     req.Erc20,
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryERC20ToDenomResponse{Denom: response.Denom, FxOriginated: response.Denom == fxtypes.DefaultDenom}, nil
}

func (k queryServer) GetDelegateKeyByValidator(c context.Context, req *types.QueryDelegateKeyByValidatorRequest) (*types.QueryDelegateKeyByValidatorResponse, error) {
	response, err := k.Keeper.GetOracleByAddr(c, &crosschaintypes.QueryOracleByAddrRequest{
		ChainName:     ethtypes.ModuleName,
		OracleAddress: req.ValidatorAddress,
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryDelegateKeyByValidatorResponse{
		EthAddress:          response.Oracle.ExternalAddress,
		OrchestratorAddress: response.Oracle.BridgerAddress,
	}, nil
}

func (k queryServer) GetDelegateKeyByOrchestrator(c context.Context, req *types.QueryDelegateKeyByOrchestratorRequest) (*types.QueryDelegateKeyByOrchestratorResponse, error) {
	response, err := k.Keeper.GetOracleByBridgerAddr(c, &crosschaintypes.QueryOracleByBridgerAddrRequest{
		ChainName:      ethtypes.ModuleName,
		BridgerAddress: req.OrchestratorAddress,
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryDelegateKeyByOrchestratorResponse{
		ValidatorAddress: response.Oracle.OracleAddress,
		EthAddress:       response.Oracle.ExternalAddress,
	}, nil
}

func (k queryServer) GetDelegateKeyByEth(c context.Context, req *types.QueryDelegateKeyByEthRequest) (*types.QueryDelegateKeyByEthResponse, error) {
	response, err := k.Keeper.GetOracleByExternalAddr(c, &crosschaintypes.QueryOracleByExternalAddrRequest{
		ChainName:       ethtypes.ModuleName,
		ExternalAddress: req.EthAddress,
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryDelegateKeyByEthResponse{
		ValidatorAddress:    response.Oracle.OracleAddress,
		OrchestratorAddress: response.Oracle.BridgerAddress,
	}, nil
}

func (k queryServer) GetPendingSendToEth(c context.Context, req *types.QueryPendingSendToEthRequest) (*types.QueryPendingSendToEthResponse, error) {
	response, err := k.Keeper.GetPendingSendToExternal(c, &crosschaintypes.QueryPendingSendToExternalRequest{
		ChainName:     ethtypes.ModuleName,
		SenderAddress: req.SenderAddress,
	})
	if err != nil {
		return nil, err
	}

	res := &types.QueryPendingSendToEthResponse{
		TransfersInBatches: make([]*types.OutgoingTransferTx, len(response.TransfersInBatches)),
		UnbatchedTransfers: make([]*types.OutgoingTransferTx, len(response.UnbatchedTransfers)),
	}
	for i := 0; i < len(response.TransfersInBatches); i++ {
		res.TransfersInBatches[i] = &types.OutgoingTransferTx{
			Id:          response.TransfersInBatches[i].Id,
			Sender:      response.TransfersInBatches[i].Sender,
			DestAddress: response.TransfersInBatches[i].DestAddress,
			Erc20Token: &types.ERC20Token{
				Contract: response.TransfersInBatches[i].Token.Contract,
				Amount:   response.TransfersInBatches[i].Token.Amount,
			},
			Erc20Fee: &types.ERC20Token{
				Contract: response.TransfersInBatches[i].Fee.Contract,
				Amount:   response.TransfersInBatches[i].Fee.Amount,
			},
		}
	}
	for i := 0; i < len(response.UnbatchedTransfers); i++ {
		res.UnbatchedTransfers[i] = &types.OutgoingTransferTx{
			Id:          response.UnbatchedTransfers[i].Id,
			Sender:      response.UnbatchedTransfers[i].Sender,
			DestAddress: response.UnbatchedTransfers[i].DestAddress,
			Erc20Token: &types.ERC20Token{
				Contract: response.UnbatchedTransfers[i].Token.Contract,
				Amount:   response.UnbatchedTransfers[i].Token.Amount,
			},
			Erc20Fee: &types.ERC20Token{
				Contract: response.UnbatchedTransfers[i].Fee.Contract,
				Amount:   response.UnbatchedTransfers[i].Fee.Amount,
			},
		}
	}

	return res, nil
}

func (k queryServer) LastObservedBlockHeight(c context.Context, _ *types.QueryLastObservedBlockHeightRequest) (*types.QueryLastObservedBlockHeightResponse, error) {
	response, err := k.Keeper.LastObservedBlockHeight(c, &crosschaintypes.QueryLastObservedBlockHeightRequest{ChainName: ethtypes.ModuleName})
	if err != nil {
		return nil, err
	}
	return &types.QueryLastObservedBlockHeightResponse{BlockHeight: response.BlockHeight, EthBlockHeight: response.ExternalBlockHeight}, nil
}

func (k queryServer) LastEventBlockHeightByAddr(c context.Context, req *types.QueryLastEventBlockHeightByAddrRequest) (*types.QueryLastEventBlockHeightByAddrResponse, error) {
	response, err := k.Keeper.LastEventBlockHeightByAddr(c, &crosschaintypes.QueryLastEventBlockHeightByAddrRequest{
		ChainName:      ethtypes.ModuleName,
		BridgerAddress: req.Address,
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryLastEventBlockHeightByAddrResponse{BlockHeight: response.BlockHeight}, nil
}

func (k queryServer) ProjectedBatchTimeoutHeight(c context.Context, _ *types.QueryProjectedBatchTimeoutHeightRequest) (*types.QueryProjectedBatchTimeoutHeightResponse, error) {
	response, err := k.Keeper.ProjectedBatchTimeoutHeight(c, &crosschaintypes.QueryProjectedBatchTimeoutHeightRequest{ChainName: ethtypes.ModuleName})
	if err != nil {
		return nil, err
	}
	return &types.QueryProjectedBatchTimeoutHeightResponse{TimeoutHeight: response.TimeoutHeight}, nil
}

func (k queryServer) BridgeTokens(c context.Context, _ *types.QueryBridgeTokensRequest) (*types.QueryBridgeTokensResponse, error) {
	response, err := k.Keeper.BridgeTokens(c, &crosschaintypes.QueryBridgeTokensRequest{ChainName: ethtypes.ModuleName})
	if err != nil {
		return nil, err
	}
	bridgeTokens := make([]*types.ERC20ToDenom, len(response.BridgeTokens))
	for i := 0; i < len(response.BridgeTokens); i++ {
		bridgeTokens[i] = &types.ERC20ToDenom{
			Erc20: response.BridgeTokens[i].Token,
			Denom: response.BridgeTokens[i].Denom,
		}
	}
	return &types.QueryBridgeTokensResponse{BridgeTokens: bridgeTokens}, nil
}
