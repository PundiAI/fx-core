package keeper

import (
	"context"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/crosschain/types"
)

var _ types.QueryServer = Keeper{}

// Params queries the params of the bsc module
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	var params types.Params
	k.paramSpace.GetParamSet(sdk.UnwrapSDKContext(c), &params)
	return &types.QueryParamsResponse{Params: params}, nil

}

// CurrentOracleSet queries the CurrentOracleSet of the bsc module
func (k Keeper) CurrentOracleSet(c context.Context, _ *types.QueryCurrentOracleSetRequest) (*types.QueryCurrentOracleSetResponse, error) {
	return &types.QueryCurrentOracleSetResponse{OracleSet: k.GetCurrentOracleSet(sdk.UnwrapSDKContext(c))}, nil
}

// OracleSetRequest queries the OracleSetRequest of the bsc module
func (k Keeper) OracleSetRequest(c context.Context, req *types.QueryOracleSetRequestRequest) (*types.QueryOracleSetRequestResponse, error) {
	return &types.QueryOracleSetRequestResponse{OracleSet: k.GetOracleSet(sdk.UnwrapSDKContext(c), req.Nonce)}, nil
}

// OracleSetConfirm queries the OracleSetConfirm of the bsc module
func (k Keeper) OracleSetConfirm(c context.Context, req *types.QueryOracleSetConfirmRequest) (*types.QueryOracleSetConfirmResponse, error) {
	orcAddr, err := sdk.AccAddressFromBech32(req.OrchestratorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}
	sdkCtx := sdk.UnwrapSDKContext(c)
	oracleAddr, found := k.GetOracleAddressByOrchestratorKey(sdkCtx, orcAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNotOracle, orcAddr.String())
	}
	return &types.QueryOracleSetConfirmResponse{Confirm: k.GetOracleSetConfirm(sdkCtx, req.Nonce, oracleAddr)}, nil
}

// OracleSetConfirmsByNonce queries the OracleSetConfirmsByNonce of the bsc module
func (k Keeper) OracleSetConfirmsByNonce(c context.Context, req *types.QueryOracleSetConfirmsByNonceRequest) (*types.QueryOracleSetConfirmsByNonceResponse, error) {
	var confirms []*types.MsgOracleSetConfirm
	k.IterateOracleSetConfirmByNonce(sdk.UnwrapSDKContext(c), req.Nonce, func(_ []byte, c types.MsgOracleSetConfirm) bool {
		confirms = append(confirms, &c)
		return false
	})
	return &types.QueryOracleSetConfirmsByNonceResponse{Confirms: confirms}, nil
}

// LastOracleSetRequests queries the LastOracleSetRequests of the bsc module
func (k Keeper) LastOracleSetRequests(c context.Context, _ *types.QueryLastOracleSetRequestsRequest) (*types.QueryLastOracleSetRequestsResponse, error) {
	valReq := k.GetOracleSets(sdk.UnwrapSDKContext(c))
	valReqLen := len(valReq)
	retLen := 0
	if valReqLen < maxOracleSetRequestsReturned {
		retLen = valReqLen
	} else {
		retLen = maxOracleSetRequestsReturned
	}
	return &types.QueryLastOracleSetRequestsResponse{OracleSets: valReq[0:retLen]}, nil
}

// LastPendingOracleSetRequestByAddr queries the LastPendingOracleSetRequestByAddr of the bsc module
func (k Keeper) LastPendingOracleSetRequestByAddr(c context.Context, req *types.QueryLastPendingOracleSetRequestByAddrRequest) (*types.QueryLastPendingOracleSetRequestByAddrResponse, error) {
	orcAddr, err := sdk.AccAddressFromBech32(req.OrchestratorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}

	sdkCtx := sdk.UnwrapSDKContext(c)
	oracleAddr, found := k.GetOracleAddressByOrchestratorKey(sdkCtx, orcAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNotOracle, orcAddr.String())
	}
	oracle, found := k.GetOracle(sdkCtx, oracleAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNoOracleFound, oracleAddr.String())
	}
	var pendingOracleSetReq []*types.OracleSet
	k.IterateOracleSets(sdkCtx, func(_ []byte, oracleSet *types.OracleSet) bool {
		if oracle.StartHeight > int64(oracleSet.Height) {
			return false
		}
		// found is true if the operatorAddr has signed the valset we are currently looking at
		// if this valset has NOT been signed by oracleAddr, store it in pendingOracleSetReq and exit the loop
		if found = k.GetOracleSetConfirm(sdkCtx, oracleSet.Nonce, oracleAddr) != nil; !found {
			pendingOracleSetReq = append(pendingOracleSetReq, oracleSet)
		}
		// if we have more than 100 unconfirmed requests in
		// our array we should exit, pagination
		if len(pendingOracleSetReq) > MaxResults {
			return true
		}
		// return false to continue the loop
		return false
	})
	return &types.QueryLastPendingOracleSetRequestByAddrResponse{OracleSets: pendingOracleSetReq}, nil
}

// BatchFees queries the batch fees from unbatched pool
func (k Keeper) BatchFees(c context.Context, req *types.QueryBatchFeeRequest) (*types.QueryBatchFeeResponse, error) {
	if req.MinBatchFees == nil {
		req.MinBatchFees = make([]types.MinBatchFee, 0)
	}
	allBatchFees := k.GetAllBatchFees(sdk.UnwrapSDKContext(c), MaxResults, req.MinBatchFees)
	return &types.QueryBatchFeeResponse{BatchFees: allBatchFees}, nil
}

// LastPendingBatchRequestByAddr queries the LastPendingBatchRequestByAddr of the bsc module
func (k Keeper) LastPendingBatchRequestByAddr(c context.Context, req *types.QueryLastPendingBatchRequestByAddrRequest) (*types.QueryLastPendingBatchRequestByAddrResponse, error) {
	orcAddr, err := sdk.AccAddressFromBech32(req.OrchestratorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "address invalid")
	}
	sdkCtx := sdk.UnwrapSDKContext(c)
	oracleAddr, found := k.GetOracleAddressByOrchestratorKey(sdkCtx, orcAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNotOracle, orcAddr.String())
	}
	oracle, found := k.GetOracle(sdkCtx, oracleAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNoOracleFound, oracleAddr.String())
	}
	var pendingBatchReq *types.OutgoingTxBatch
	k.IterateOutgoingTXBatches(sdkCtx, func(_ []byte, batch *types.OutgoingTxBatch) bool {
		// filter startHeight before confirm
		if oracle.StartHeight > int64(batch.Block) {
			return false
		}
		foundConfirm := k.GetBatchConfirm(sdkCtx, batch.BatchNonce, batch.TokenContract, oracleAddr) != nil
		if !foundConfirm {
			pendingBatchReq = batch
			return true
		}
		return false
	})
	return &types.QueryLastPendingBatchRequestByAddrResponse{Batch: pendingBatchReq}, nil
}

// OutgoingTxBatches queries the OutgoingTxBatches of the bsc module
func (k Keeper) OutgoingTxBatches(c context.Context, _ *types.QueryOutgoingTxBatchesRequest) (*types.QueryOutgoingTxBatchesResponse, error) {
	var batches []*types.OutgoingTxBatch
	k.IterateOutgoingTXBatches(sdk.UnwrapSDKContext(c), func(_ []byte, batch *types.OutgoingTxBatch) bool {
		batches = append(batches, batch)
		return len(batches) == MaxResults
	})
	sort.Slice(batches, func(i, j int) bool {
		return batches[i].BatchTimeout < batches[j].BatchTimeout
	})
	return &types.QueryOutgoingTxBatchesResponse{Batches: batches}, nil
}

// BatchRequestByNonce queries the BatchRequestByNonce of the bsc module
func (k Keeper) BatchRequestByNonce(c context.Context, req *types.QueryBatchRequestByNonceRequest) (*types.QueryBatchRequestByNonceResponse, error) {
	if err := types.ValidateExternalAddress(req.TokenContract); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "token contract address")
	}
	foundBatch := k.GetOutgoingTXBatch(sdk.UnwrapSDKContext(c), req.TokenContract, req.Nonce)
	if foundBatch == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Can not find tx batch")
	}
	return &types.QueryBatchRequestByNonceResponse{Batch: foundBatch}, nil
}

func (k Keeper) BatchConfirm(c context.Context, req *types.QueryBatchConfirmRequest) (*types.QueryBatchConfirmResponse, error) {
	orcAddr, err := sdk.AccAddressFromBech32(req.OrchestratorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "orchestrator")
	}
	sdkCtx := sdk.UnwrapSDKContext(c)
	oracleAddr, found := k.GetOracleAddressByOrchestratorKey(sdkCtx, orcAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNotOracle, orcAddr.String())
	}
	confirm := k.GetBatchConfirm(sdkCtx, req.GetNonce(), req.TokenContract, oracleAddr)
	return &types.QueryBatchConfirmResponse{Confirm: confirm}, nil
}

// BatchConfirms returns the batch confirmations by nonce and token contract
func (k Keeper) BatchConfirms(c context.Context, req *types.QueryBatchConfirmsRequest) (*types.QueryBatchConfirmsResponse, error) {
	var confirms []*types.MsgConfirmBatch
	k.IterateBatchConfirmByNonceAndTokenContract(sdk.UnwrapSDKContext(c), req.Nonce, req.TokenContract, func(_ []byte, c types.MsgConfirmBatch) bool {
		confirms = append(confirms, &c)
		return false
	})
	return &types.QueryBatchConfirmsResponse{Confirms: confirms}, nil
}

// LastEventNonceByAddr returns the last event nonce for the given validator address, this allows eth oracles to figure out where they left off
func (k Keeper) LastEventNonceByAddr(c context.Context, req *types.QueryLastEventNonceByAddrRequest) (*types.QueryLastEventNonceByAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	orcAddr, err := sdk.AccAddressFromBech32(req.OrchestratorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "orchestrator")
	}
	oracle, found := k.GetOracleAddressByOrchestratorKey(ctx, orcAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrEmpty, "address")
	}
	lastEventNonce := k.GetLastEventNonceByOracle(ctx, oracle)
	return &types.QueryLastEventNonceByAddrResponse{EventNonce: lastEventNonce}, nil
}

func (k Keeper) DenomToToken(c context.Context, req *types.QueryDenomToTokenRequest) (*types.QueryDenomToTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	bridgeToken := k.GetDenomByBridgeToken(ctx, req.Denom)
	if bridgeToken == nil {
		return nil, sdkerrors.Wrap(types.ErrEmpty, "bridge token is not exist")
	}
	return &types.QueryDenomToTokenResponse{
		Token:      bridgeToken.Token,
		ChannelIbc: bridgeToken.ChannelIbc,
	}, nil
}

func (k Keeper) TokenToDenom(c context.Context, req *types.QueryTokenToDenomRequest) (*types.QueryTokenToDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	bridgeToken := k.GetBridgeTokenDenom(ctx, req.Token)
	if bridgeToken == nil {
		return nil, sdkerrors.Wrap(types.ErrEmpty, "bridge token is not exist")
	}
	return &types.QueryTokenToDenomResponse{
		Denom:      bridgeToken.Denom,
		ChannelIbc: bridgeToken.ChannelIbc,
	}, nil
}

func (k Keeper) GetOracleByAddr(c context.Context, req *types.QueryOracleByAddrRequest) (*types.QueryOracleResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	oracleAddr, err := sdk.AccAddressFromBech32(req.OracleAddress)
	if err != nil {
		return nil, err
	}
	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "No oracleAddr")
	}
	return &types.QueryOracleResponse{Oracle: &oracle}, nil
}

func (k Keeper) GetOracleByOrchestrator(c context.Context, req *types.QueryOracleByOrchestratorRequest) (*types.QueryOracleResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	orcAddr, err := sdk.AccAddressFromBech32(req.OrchestratorAddress)
	if err != nil {
		return nil, err
	}
	oracleAddr, found := k.GetOracleAddressByOrchestratorKey(ctx, orcAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "No Orchestrator")
	}
	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "No oracleAddr")
	}
	return &types.QueryOracleResponse{Oracle: &oracle}, nil
}

func (k Keeper) GetOracleByExternalAddr(c context.Context, req *types.QueryOracleByExternalAddrRequest) (*types.QueryOracleResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := types.ValidateExternalAddress(req.ExternalAddress); err != nil {
		return nil, sdkerrors.Wrap(err, "invalid bsc address")
	}
	oracleAddr, found := k.GetOracleByExternalAddress(ctx, req.ExternalAddress)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "No Orchestrator")
	}
	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "No oracleAddr")
	}
	return &types.QueryOracleResponse{Oracle: &oracle}, nil
}

func (k Keeper) GetPendingSendToExternal(c context.Context, req *types.QueryPendingSendToExternalRequest) (*types.QueryPendingSendToExternalResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	batches := k.GetOutgoingTxBatches(ctx)
	unbatchedTx := k.GetUnbatchedTransactions(ctx)
	senderAddress := req.SenderAddress
	var res = &types.QueryPendingSendToExternalResponse{
		TransfersInBatches: make([]*types.OutgoingTransferTx, 0),
		UnbatchedTransfers: make([]*types.OutgoingTransferTx, 0),
	}
	for _, batch := range batches {
		for _, tx := range batch.Transactions {
			if tx.Sender == senderAddress {
				res.TransfersInBatches = append(res.TransfersInBatches, tx)
			}
		}
	}
	for _, tx := range unbatchedTx {
		if tx.Sender == senderAddress {
			res.UnbatchedTransfers = append(res.UnbatchedTransfers, tx)
		}
	}
	return res, nil
}

func (k Keeper) GetIbcSequenceHeightByChannel(c context.Context, req *types.QueryIbcSequenceHeightRequest) (*types.QueryIbcSequenceHeightResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	height, found := k.GetIbcSequenceHeight(ctx, req.GetSourcePort(), req.GetSourceChannel(), req.GetSequence())
	return &types.QueryIbcSequenceHeightResponse{
		Found:       found,
		BlockHeight: height,
	}, nil
}

func (k Keeper) LastObservedBlockHeight(c context.Context, _ *types.QueryLastObservedBlockHeightRequest) (*types.QueryLastObservedBlockHeightResponse, error) {
	blockHeight := k.GetLastObservedBlockHeight(sdk.UnwrapSDKContext(c))
	return &types.QueryLastObservedBlockHeightResponse{
		ExternalBlockHeight: blockHeight.ExternalBlockHeight,
		BlockHeight:         blockHeight.BlockHeight,
	}, nil
}

func (k Keeper) LastEventBlockHeightByAddr(c context.Context, req *types.QueryLastEventBlockHeightByAddrRequest) (*types.QueryLastEventBlockHeightByAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	orcAddr, err := sdk.AccAddressFromBech32(req.OrchestratorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "orchestrator")
	}

	oracle, found := k.GetOracleAddressByOrchestratorKey(ctx, orcAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrEmpty, "address")
	}

	lastEventBlockHeight := k.getLastEventBlockHeightByOracle(ctx, oracle)
	return &types.QueryLastEventBlockHeightByAddrResponse{BlockHeight: lastEventBlockHeight}, nil
}

func (k Keeper) Oracles(c context.Context, req *types.QueryOraclesRequest) (*types.QueryOraclesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	oracles := k.GetAllOracles(ctx)
	return &types.QueryOraclesResponse{Oracles: oracles}, nil
}

func (k Keeper) ProjectedBatchTimeoutHeight(c context.Context, req *types.QueryProjectedBatchTimeoutHeightRequest) (*types.QueryProjectedBatchTimeoutHeightResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	timeout := k.GetBatchTimeoutHeight(ctx)
	return &types.QueryProjectedBatchTimeoutHeightResponse{TimeoutHeight: timeout}, nil
}

func (k Keeper) BridgeTokens(c context.Context, _ *types.QueryBridgeTokensRequest) (*types.QueryBridgeTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var bridgeTokens = make([]*types.BridgeToken, 0)
	k.IterateBridgeTokenToDenom(ctx, func(bytes []byte, token *types.BridgeToken) bool {
		bridgeTokens = append(bridgeTokens, token)
		return false
	})
	return &types.QueryBridgeTokensResponse{BridgeTokens: bridgeTokens}, nil
}
