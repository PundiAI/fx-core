package keeper

import (
	"context"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
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
	bridgerAddr, err := sdk.AccAddressFromBech32(req.GetBridgerAddress())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bridger address")
	}
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	sdkCtx := sdk.UnwrapSDKContext(c)
	oracleAddr, found := k.GetOracleAddressByBridgerKey(sdkCtx, bridgerAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	return &types.QueryOracleSetConfirmResponse{Confirm: k.GetOracleSetConfirm(sdkCtx, req.Nonce, oracleAddr)}, nil
}

// OracleSetConfirmsByNonce queries the OracleSetConfirmsByNonce of the bsc module
func (k Keeper) OracleSetConfirmsByNonce(c context.Context, req *types.QueryOracleSetConfirmsByNonceRequest) (*types.QueryOracleSetConfirmsByNonceResponse, error) {
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	var confirms []*types.MsgOracleSetConfirm
	k.IterateOracleSetConfirmByNonce(sdk.UnwrapSDKContext(c), req.Nonce, func(_ []byte, c types.MsgOracleSetConfirm) bool {
		confirms = append(confirms, &c)
		return false
	})
	return &types.QueryOracleSetConfirmsByNonceResponse{Confirms: confirms}, nil
}

// LastOracleSetRequests queries the LastOracleSetRequests of the bsc module
func (k Keeper) LastOracleSetRequests(c context.Context, _ *types.QueryLastOracleSetRequestsRequest) (*types.QueryLastOracleSetRequestsResponse, error) {
	oracleSets := k.GetOracleSets(sdk.UnwrapSDKContext(c))
	retLen := 0
	if len(oracleSets) < maxOracleSetRequestsReturned {
		retLen = len(oracleSets)
	} else {
		retLen = maxOracleSetRequestsReturned
	}
	return &types.QueryLastOracleSetRequestsResponse{OracleSets: oracleSets[0:retLen]}, nil
}

// LastPendingOracleSetRequestByAddr queries the LastPendingOracleSetRequestByAddr of the bsc module
func (k Keeper) LastPendingOracleSetRequestByAddr(c context.Context, req *types.QueryLastPendingOracleSetRequestByAddrRequest) (*types.QueryLastPendingOracleSetRequestByAddrResponse, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(req.GetBridgerAddress())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bridger address")
	}

	sdkCtx := sdk.UnwrapSDKContext(c)
	oracleAddr, found := k.GetOracleAddressByBridgerKey(sdkCtx, bridgerAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	oracle, found := k.GetOracle(sdkCtx, oracleAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	var pendingOracleSetReq []*types.OracleSet
	k.IterateOracleSets(sdkCtx, func(_ []byte, oracleSet *types.OracleSet) bool {
		if oracle.StartHeight > int64(oracleSet.Height) {
			return false
		}
		// found is true if the operatorAddr has signed the oracle set we are currently looking at
		// if this oracle set has NOT been signed by oracleAddr, store it in pendingOracleSetReq and exit the loop
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
	if req.GetMinBatchFees() == nil {
		req.MinBatchFees = make([]types.MinBatchFee, 0)
	}
	for _, fee := range req.MinBatchFees {
		if fee.BaseFee.IsNil() || fee.BaseFee.IsNegative() {
			return nil, status.Error(codes.InvalidArgument, "base fee")
		}
		if err := fxtypes.ValidateEthereumAddress(fee.TokenContract); err != nil {
			return nil, status.Error(codes.InvalidArgument, "token contract")
		}
	}
	allBatchFees := k.GetAllBatchFees(sdk.UnwrapSDKContext(c), MaxResults, req.MinBatchFees)
	return &types.QueryBatchFeeResponse{BatchFees: allBatchFees}, nil
}

// LastPendingBatchRequestByAddr queries the LastPendingBatchRequestByAddr of the bsc module
func (k Keeper) LastPendingBatchRequestByAddr(c context.Context, req *types.QueryLastPendingBatchRequestByAddrRequest) (*types.QueryLastPendingBatchRequestByAddrResponse, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(req.GetBridgerAddress())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bridger address")
	}
	sdkCtx := sdk.UnwrapSDKContext(c)
	oracleAddr, found := k.GetOracleAddressByBridgerKey(sdkCtx, bridgerAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	oracle, found := k.GetOracle(sdkCtx, oracleAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	var pendingBatchReq *types.OutgoingTxBatch
	k.IterateOutgoingTxBatches(sdkCtx, func(_ []byte, batch *types.OutgoingTxBatch) bool {
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
	k.IterateOutgoingTxBatches(sdk.UnwrapSDKContext(c), func(_ []byte, batch *types.OutgoingTxBatch) bool {
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
	if err := fxtypes.ValidateEthereumAddress(req.GetTokenContract()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "token contract address")
	}
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	foundBatch := k.GetOutgoingTxBatch(sdk.UnwrapSDKContext(c), req.TokenContract, req.Nonce)
	if foundBatch == nil {
		return nil, status.Error(codes.NotFound, "tx batch")
	}
	return &types.QueryBatchRequestByNonceResponse{Batch: foundBatch}, nil
}

func (k Keeper) BatchConfirm(c context.Context, req *types.QueryBatchConfirmRequest) (*types.QueryBatchConfirmResponse, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(req.BridgerAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bridger address")
	}
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	sdkCtx := sdk.UnwrapSDKContext(c)
	oracleAddr, found := k.GetOracleAddressByBridgerKey(sdkCtx, bridgerAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	confirm := k.GetBatchConfirm(sdkCtx, req.Nonce, req.TokenContract, oracleAddr)
	return &types.QueryBatchConfirmResponse{Confirm: confirm}, nil
}

// BatchConfirms returns the batch confirmations by nonce and token contract
func (k Keeper) BatchConfirms(c context.Context, req *types.QueryBatchConfirmsRequest) (*types.QueryBatchConfirmsResponse, error) {
	if err := fxtypes.ValidateEthereumAddress(req.GetTokenContract()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "token contract address")
	}
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	var confirms []*types.MsgConfirmBatch
	k.IterateBatchConfirmByNonceAndTokenContract(sdk.UnwrapSDKContext(c), req.Nonce, req.TokenContract, func(_ []byte, c types.MsgConfirmBatch) bool {
		confirms = append(confirms, &c)
		return false
	})
	return &types.QueryBatchConfirmsResponse{Confirms: confirms}, nil
}

// LastEventNonceByAddr returns the last event nonce for the given validator address, this allows eth oracles to figure out where they left off
func (k Keeper) LastEventNonceByAddr(c context.Context, req *types.QueryLastEventNonceByAddrRequest) (*types.QueryLastEventNonceByAddrResponse, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(req.BridgerAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bridger address")
	}
	ctx := sdk.UnwrapSDKContext(c)

	oracle, found := k.GetOracleAddressByBridgerKey(ctx, bridgerAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	lastEventNonce := k.GetLastEventNonceByOracle(ctx, oracle)
	return &types.QueryLastEventNonceByAddrResponse{EventNonce: lastEventNonce}, nil
}

func (k Keeper) DenomToToken(c context.Context, req *types.QueryDenomToTokenRequest) (*types.QueryDenomToTokenResponse, error) {
	if len(req.GetDenom()) <= 0 {
		return nil, status.Error(codes.InvalidArgument, "denom")
	}

	bridgeToken := k.GetDenomByBridgeToken(sdk.UnwrapSDKContext(c), req.Denom)
	if bridgeToken == nil {
		return nil, status.Error(codes.NotFound, "bridge token")
	}
	return &types.QueryDenomToTokenResponse{
		Token:      bridgeToken.Token,
		ChannelIbc: bridgeToken.ChannelIbc,
	}, nil
}

func (k Keeper) TokenToDenom(c context.Context, req *types.QueryTokenToDenomRequest) (*types.QueryTokenToDenomResponse, error) {
	if err := fxtypes.ValidateEthereumAddress(req.GetToken()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "token address")
	}
	bridgeToken := k.GetBridgeTokenDenom(sdk.UnwrapSDKContext(c), req.Token)
	if bridgeToken == nil {
		return nil, status.Error(codes.NotFound, "bridge token")
	}
	return &types.QueryTokenToDenomResponse{
		Denom:      bridgeToken.Denom,
		ChannelIbc: bridgeToken.ChannelIbc,
	}, nil
}

func (k Keeper) GetOracleByAddr(c context.Context, req *types.QueryOracleByAddrRequest) (*types.QueryOracleResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(req.OracleAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "oracle address")
	}
	oracle, found := k.GetOracle(sdk.UnwrapSDKContext(c), oracleAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	return &types.QueryOracleResponse{Oracle: &oracle}, nil
}

func (k Keeper) GetOracleByBridgerAddr(c context.Context, req *types.QueryOracleByBridgerAddrRequest) (*types.QueryOracleResponse, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(req.GetBridgerAddress())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bridger address")
	}
	ctx := sdk.UnwrapSDKContext(c)

	oracleAddr, found := k.GetOracleAddressByBridgerKey(ctx, bridgerAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	return &types.QueryOracleResponse{Oracle: &oracle}, nil
}

func (k Keeper) GetOracleByExternalAddr(c context.Context, req *types.QueryOracleByExternalAddrRequest) (*types.QueryOracleResponse, error) {
	if err := fxtypes.ValidateEthereumAddress(req.GetExternalAddress()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "external address")
	}

	ctx := sdk.UnwrapSDKContext(c)
	oracleAddr, found := k.GetOracleByExternalAddress(ctx, req.ExternalAddress)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	return &types.QueryOracleResponse{Oracle: &oracle}, nil
}

func (k Keeper) GetPendingSendToExternal(c context.Context, req *types.QueryPendingSendToExternalRequest) (*types.QueryPendingSendToExternalResponse, error) {
	if _, err := sdk.AccAddressFromBech32(req.GetSenderAddress()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "sender address")
	}

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

func (k Keeper) LastObservedBlockHeight(c context.Context, _ *types.QueryLastObservedBlockHeightRequest) (*types.QueryLastObservedBlockHeightResponse, error) {
	blockHeight := k.GetLastObservedBlockHeight(sdk.UnwrapSDKContext(c))
	return &types.QueryLastObservedBlockHeightResponse{
		ExternalBlockHeight: blockHeight.ExternalBlockHeight,
		BlockHeight:         blockHeight.BlockHeight,
	}, nil
}

func (k Keeper) LastEventBlockHeightByAddr(c context.Context, req *types.QueryLastEventBlockHeightByAddrRequest) (*types.QueryLastEventBlockHeightByAddrResponse, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(req.GetBridgerAddress())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bridger address")
	}
	ctx := sdk.UnwrapSDKContext(c)

	oracle, found := k.GetOracleAddressByBridgerKey(ctx, bridgerAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}

	lastEventBlockHeight := k.GetLastEventBlockHeightByOracle(ctx, oracle)
	return &types.QueryLastEventBlockHeightByAddrResponse{BlockHeight: lastEventBlockHeight}, nil
}

func (k Keeper) Oracles(c context.Context, _ *types.QueryOraclesRequest) (*types.QueryOraclesResponse, error) {
	oracles := k.GetAllOracles(sdk.UnwrapSDKContext(c), false)
	return &types.QueryOraclesResponse{Oracles: oracles}, nil
}

func (k Keeper) ProjectedBatchTimeoutHeight(c context.Context, _ *types.QueryProjectedBatchTimeoutHeightRequest) (*types.QueryProjectedBatchTimeoutHeightResponse, error) {
	timeout := k.GetBatchTimeoutHeight(sdk.UnwrapSDKContext(c))
	return &types.QueryProjectedBatchTimeoutHeightResponse{TimeoutHeight: timeout}, nil
}

func (k Keeper) BridgeTokens(c context.Context, _ *types.QueryBridgeTokensRequest) (*types.QueryBridgeTokensResponse, error) {
	var bridgeTokens = make([]*types.BridgeToken, 0)
	k.IterateBridgeTokenToDenom(sdk.UnwrapSDKContext(c), func(bytes []byte, token *types.BridgeToken) bool {
		bridgeTokens = append(bridgeTokens, token)
		return false
	})
	return &types.QueryBridgeTokensResponse{BridgeTokens: bridgeTokens}, nil
}
