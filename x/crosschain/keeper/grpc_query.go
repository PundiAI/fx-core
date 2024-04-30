package keeper

import (
	"context"
	"sort"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/functionx/fx-core/v7/contract"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := k.GetParams(sdk.UnwrapSDKContext(c))
	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Keeper) CurrentOracleSet(c context.Context, _ *types.QueryCurrentOracleSetRequest) (*types.QueryCurrentOracleSetResponse, error) {
	return &types.QueryCurrentOracleSetResponse{OracleSet: k.GetCurrentOracleSet(sdk.UnwrapSDKContext(c))}, nil
}

func (k Keeper) OracleSetRequest(c context.Context, req *types.QueryOracleSetRequestRequest) (*types.QueryOracleSetRequestResponse, error) {
	return &types.QueryOracleSetRequestResponse{OracleSet: k.GetOracleSet(sdk.UnwrapSDKContext(c), req.Nonce)}, nil
}

func (k Keeper) OracleSetConfirm(c context.Context, req *types.QueryOracleSetConfirmRequest) (*types.QueryOracleSetConfirmResponse, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(req.GetBridgerAddress())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bridger address")
	}
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	ctx := sdk.UnwrapSDKContext(c)
	oracleAddr, found := k.GetOracleAddrByBridgerAddr(ctx, bridgerAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	return &types.QueryOracleSetConfirmResponse{Confirm: k.GetOracleSetConfirm(ctx, req.Nonce, oracleAddr)}, nil
}

func (k Keeper) OracleSetConfirmsByNonce(c context.Context, req *types.QueryOracleSetConfirmsByNonceRequest) (*types.QueryOracleSetConfirmsByNonceResponse, error) {
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	var confirms []*types.MsgOracleSetConfirm
	k.IterateOracleSetConfirmByNonce(sdk.UnwrapSDKContext(c), req.Nonce, func(confirm *types.MsgOracleSetConfirm) bool {
		confirms = append(confirms, confirm)
		return false
	})
	return &types.QueryOracleSetConfirmsByNonceResponse{Confirms: confirms}, nil
}

func (k Keeper) LastOracleSetRequests(c context.Context, _ *types.QueryLastOracleSetRequestsRequest) (*types.QueryLastOracleSetRequestsResponse, error) {
	var oraclesSets []*types.OracleSet
	k.IterateOracleSets(sdk.UnwrapSDKContext(c), true, func(oracleSet *types.OracleSet) bool {
		if len(oraclesSets) >= types.MaxOracleSetRequestsResults {
			return true
		}
		oraclesSets = append(oraclesSets, oracleSet)
		return false
	})
	return &types.QueryLastOracleSetRequestsResponse{OracleSets: oraclesSets}, nil
}

func (k Keeper) LastPendingOracleSetRequestByAddr(c context.Context, req *types.QueryLastPendingOracleSetRequestByAddrRequest) (*types.QueryLastPendingOracleSetRequestByAddrResponse, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(req.GetBridgerAddress())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bridger address")
	}

	ctx := sdk.UnwrapSDKContext(c)
	oracleAddr, found := k.GetOracleAddrByBridgerAddr(ctx, bridgerAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	var pendingOracleSetReq []*types.OracleSet
	k.IterateOracleSets(ctx, false, func(oracleSet *types.OracleSet) bool {
		if oracle.StartHeight > int64(oracleSet.Height) {
			return false
		}
		// found is true if the operatorAddr has signed the oracle set we are currently looking at
		// if this oracle set has NOT been signed by oracleAddr, store it in pendingOracleSetReq and exit the loop
		if found = k.GetOracleSetConfirm(ctx, oracleSet.Nonce, oracleAddr) != nil; !found {
			pendingOracleSetReq = append(pendingOracleSetReq, oracleSet)
		}
		// if we have more than 100 unconfirmed requests in
		// our array we should exit, pagination
		return len(pendingOracleSetReq) == types.MaxResults
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
		if err := contract.ValidateEthereumAddress(fee.TokenContract); err != nil {
			return nil, status.Error(codes.InvalidArgument, "token contract")
		}
	}
	allBatchFees := k.GetAllBatchFees(sdk.UnwrapSDKContext(c), types.MaxResults, req.MinBatchFees)
	return &types.QueryBatchFeeResponse{BatchFees: allBatchFees}, nil
}

func (k Keeper) LastPendingBatchRequestByAddr(c context.Context, req *types.QueryLastPendingBatchRequestByAddrRequest) (*types.QueryLastPendingBatchRequestByAddrResponse, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(req.GetBridgerAddress())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bridger address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	oracleAddr, found := k.GetOracleAddrByBridgerAddr(ctx, bridgerAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	var pendingBatchReq *types.OutgoingTxBatch
	k.IterateOutgoingTxBatches(ctx, func(batch *types.OutgoingTxBatch) bool {
		// filter startHeight before confirm
		if oracle.StartHeight > int64(batch.Block) {
			return false
		}
		foundConfirm := k.GetBatchConfirm(ctx, batch.TokenContract, batch.BatchNonce, oracleAddr) != nil
		if !foundConfirm {
			pendingBatchReq = batch
			return true
		}
		return false
	})
	return &types.QueryLastPendingBatchRequestByAddrResponse{Batch: pendingBatchReq}, nil
}

func (k Keeper) OutgoingTxBatches(c context.Context, _ *types.QueryOutgoingTxBatchesRequest) (*types.QueryOutgoingTxBatchesResponse, error) {
	var batches []*types.OutgoingTxBatch
	k.IterateOutgoingTxBatches(sdk.UnwrapSDKContext(c), func(batch *types.OutgoingTxBatch) bool {
		batches = append(batches, batch)
		return len(batches) == types.MaxResults
	})
	sort.Slice(batches, func(i, j int) bool {
		return batches[i].BatchTimeout < batches[j].BatchTimeout
	})
	return &types.QueryOutgoingTxBatchesResponse{Batches: batches}, nil
}

func (k Keeper) BatchRequestByNonce(c context.Context, req *types.QueryBatchRequestByNonceRequest) (*types.QueryBatchRequestByNonceResponse, error) {
	if err := contract.ValidateEthereumAddress(req.GetTokenContract()); err != nil {
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
	ctx := sdk.UnwrapSDKContext(c)
	oracleAddr, found := k.GetOracleAddrByBridgerAddr(ctx, bridgerAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	confirm := k.GetBatchConfirm(ctx, req.TokenContract, req.Nonce, oracleAddr)
	return &types.QueryBatchConfirmResponse{Confirm: confirm}, nil
}

// BatchConfirms returns the batch confirmations by nonce and token contract
func (k Keeper) BatchConfirms(c context.Context, req *types.QueryBatchConfirmsRequest) (*types.QueryBatchConfirmsResponse, error) {
	if err := contract.ValidateEthereumAddress(req.GetTokenContract()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "token contract address")
	}
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	var confirms []*types.MsgConfirmBatch
	k.IterateBatchConfirmByNonceAndTokenContract(sdk.UnwrapSDKContext(c), req.Nonce, req.TokenContract, func(confirm *types.MsgConfirmBatch) bool {
		confirms = append(confirms, confirm)
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

	oracle, found := k.GetOracleAddrByBridgerAddr(ctx, bridgerAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle")
	}
	lastEventNonce := k.GetLastEventNonceByOracle(ctx, oracle)
	return &types.QueryLastEventNonceByAddrResponse{EventNonce: lastEventNonce}, nil
}

func (k Keeper) DenomToToken(c context.Context, req *types.QueryDenomToTokenRequest) (*types.QueryDenomToTokenResponse, error) {
	if len(req.GetDenom()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "denom")
	}

	bridgeToken := k.GetDenomBridgeToken(sdk.UnwrapSDKContext(c), req.Denom)
	if bridgeToken == nil {
		return nil, status.Error(codes.NotFound, "bridge token")
	}
	return &types.QueryDenomToTokenResponse{
		Token: bridgeToken.Token,
	}, nil
}

func (k Keeper) TokenToDenom(c context.Context, req *types.QueryTokenToDenomRequest) (*types.QueryTokenToDenomResponse, error) {
	if err := contract.ValidateEthereumAddress(req.Token); err != nil {
		return nil, status.Error(codes.InvalidArgument, "token address")
	}
	bridgeToken := k.GetBridgeTokenDenom(sdk.UnwrapSDKContext(c), req.Token)
	if bridgeToken == nil {
		return nil, status.Error(codes.NotFound, "bridge token")
	}
	return &types.QueryTokenToDenomResponse{
		Denom: bridgeToken.Denom,
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

	oracleAddr, found := k.GetOracleAddrByBridgerAddr(ctx, bridgerAddr)
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
	if err := contract.ValidateEthereumAddress(req.GetExternalAddress()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "external address")
	}

	ctx := sdk.UnwrapSDKContext(c)
	oracleAddr, found := k.GetOracleAddrByExternalAddr(ctx, req.ExternalAddress)
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
	var batches []*types.OutgoingTxBatch
	k.IterateOutgoingTxBatches(ctx, func(batch *types.OutgoingTxBatch) bool {
		batches = append(batches, batch)
		return false
	})
	res := &types.QueryPendingSendToExternalResponse{
		TransfersInBatches: make([]*types.OutgoingTransferTx, 0),
		UnbatchedTransfers: make([]*types.OutgoingTransferTx, 0),
	}
	for _, batch := range batches {
		for _, tx := range batch.Transactions {
			if tx.Sender == req.SenderAddress {
				res.TransfersInBatches = append(res.TransfersInBatches, tx)
			}
		}
	}
	k.IterateUnbatchedTransactions(ctx, "", func(tx *types.OutgoingTransferTx) bool {
		if tx.Sender == req.SenderAddress {
			res.UnbatchedTransfers = append(res.UnbatchedTransfers, tx)
		}
		return false
	})
	return res, nil
}

func (k Keeper) GetPendingPoolSendToExternal(c context.Context, req *types.QueryPendingPoolSendToExternalRequest) (*types.QueryPendingPoolSendToExternalResponse, error) {
	if _, err := sdk.AccAddressFromBech32(req.GetSenderAddress()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "sender address")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	pendingOutgoingStore := prefix.NewStore(store, types.PendingOutgoingTxPoolKey)

	txs, pageRes, err := query.GenericFilteredPaginate(k.cdc, pendingOutgoingStore, req.Pagination, func(key []byte, tx *types.PendingOutgoingTransferTx) (*types.PendingOutgoingTransferTx, error) {
		if tx.Sender != req.SenderAddress {
			return nil, nil
		}
		return tx, nil
	}, func() *types.PendingOutgoingTransferTx {
		return &types.PendingOutgoingTransferTx{}
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPendingPoolSendToExternalResponse{Txs: txs, Pagination: pageRes}, nil
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

	oracle, found := k.GetOracleAddrByBridgerAddr(ctx, bridgerAddr)
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
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)
	timeout := k.CalExternalTimeoutHeight(ctx, params, params.ExternalBatchTimeout)
	return &types.QueryProjectedBatchTimeoutHeightResponse{TimeoutHeight: timeout}, nil
}

func (k Keeper) BridgeTokens(c context.Context, _ *types.QueryBridgeTokensRequest) (*types.QueryBridgeTokensResponse, error) {
	bridgeTokens := make([]*types.BridgeToken, 0)
	k.IterateBridgeTokenToDenom(sdk.UnwrapSDKContext(c), func(token *types.BridgeToken) bool {
		bridgeTokens = append(bridgeTokens, token)
		return false
	})
	return &types.QueryBridgeTokensResponse{BridgeTokens: bridgeTokens}, nil
}

func (k Keeper) BridgeCoinByDenom(c context.Context, req *types.QueryBridgeCoinByDenomRequest) (*types.QueryBridgeCoinByDenomResponse, error) {
	if len(req.GetDenom()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "denom")
	}
	ctx := sdk.UnwrapSDKContext(c)

	var bridgeCoinMetaData banktypes.Metadata
	k.bankKeeper.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
		if metadata.GetBase() == req.GetDenom() {
			bridgeCoinMetaData = metadata
			return true
		}
		if len(metadata.GetDenomUnits()) == 0 {
			return false
		}
		for _, alias := range metadata.GetDenomUnits()[0].GetAliases() {
			if alias == req.GetDenom() {
				bridgeCoinMetaData = metadata
				return true
			}
		}
		return false
	})
	if len(bridgeCoinMetaData.GetBase()) == 0 {
		return nil, status.Error(codes.NotFound, "denom")
	}

	bridgeCoinDenom := k.erc20Keeper.ToTargetDenom(
		ctx,
		req.GetDenom(),
		bridgeCoinMetaData.GetBase(),
		bridgeCoinMetaData.GetDenomUnits()[0].GetAliases(),
		fxtypes.ParseFxTarget(req.GetChainName()),
	)

	token := k.GetDenomBridgeToken(ctx, bridgeCoinDenom)
	if token == nil {
		return nil, status.Error(codes.NotFound, "denom")
	}

	supply := k.bankKeeper.GetSupply(ctx, bridgeCoinDenom)
	return &types.QueryBridgeCoinByDenomResponse{Coin: supply}, nil
}

func (k Keeper) BridgeChainList(_ context.Context, _ *types.QueryBridgeChainListRequest) (*types.QueryBridgeChainListResponse, error) {
	return &types.QueryBridgeChainListResponse{ChainNames: types.GetValidateChains()}, nil
}

func (k Keeper) BridgeCallConfirmByNonce(c context.Context, req *types.QueryBridgeCallConfirmByNonceRequest) (*types.QueryBridgeCallConfirmByNonceResponse, error) {
	if req.GetEventNonce() == 0 {
		return nil, status.Error(codes.InvalidArgument, "event nonce")
	}

	ctx := sdk.UnwrapSDKContext(c)
	currentOracleSet := k.GetCurrentOracleSet(ctx)
	confirmPowers := uint64(0)
	bridgeCallConfirms := make([]*types.MsgBridgeCallConfirm, 0)
	k.IterBridgeCallConfirmByNonce(ctx, req.GetEventNonce(), func(msg *types.MsgBridgeCallConfirm) bool {
		power, found := currentOracleSet.GetBridgePower(msg.ExternalAddress)
		if !found {
			return false
		}
		confirmPowers += power
		bridgeCallConfirms = append(bridgeCallConfirms, msg)
		return false
	})
	totalPower := currentOracleSet.GetTotalPower()
	requiredPower := types.AttestationVotesPowerThreshold.Mul(sdkmath.NewIntFromUint64(totalPower)).Quo(sdkmath.NewInt(100))
	enoughPower := requiredPower.GTE(sdkmath.NewIntFromUint64(confirmPowers))
	return &types.QueryBridgeCallConfirmByNonceResponse{Confirms: bridgeCallConfirms, EnoughPower: enoughPower}, nil
}

func (k Keeper) BridgeCallByNonce(c context.Context, req *types.QueryBridgeCallByNonceRequest) (*types.QueryBridgeCallByNonceResponse, error) {
	outgoingBridgeCall, found := k.GetOutgoingBridgeCallByNonce(sdk.UnwrapSDKContext(c), req.GetEventNonce())
	if !found {
		return nil, status.Error(codes.NotFound, "outgoing bridge call not found")
	}
	return &types.QueryBridgeCallByNonceResponse{BridgeCall: outgoingBridgeCall}, nil
}

func (k Keeper) BridgeCallBySender(c context.Context, req *types.QueryBridgeCallBySenderRequest) (*types.QueryBridgeCallBySenderResponse, error) {
	var outgoingBridgeCalls []*types.OutgoingBridgeCall
	k.IterateOutgoingBridgeCallsByAddress(sdk.UnwrapSDKContext(c), req.GetSenderAddress(), func(outgoingBridgeCall *types.OutgoingBridgeCall) bool {
		outgoingBridgeCalls = append(outgoingBridgeCalls, outgoingBridgeCall)
		return false
	})
	return &types.QueryBridgeCallBySenderResponse{BridgeCalls: outgoingBridgeCalls}, nil
}

func (k Keeper) LastPendingBridgeCallByAddr(c context.Context, req *types.QueryLastPendingBridgeCallByAddrRequest) (*types.QueryLastPendingBridgeCallByAddrResponse, error) {
	if len(req.GetBridgerAddress()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty external address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	bridgerAddr := sdk.MustAccAddressFromBech32(req.GetBridgerAddress())
	oracleAddr, found := k.GetOracleAddrByBridgerAddr(ctx, bridgerAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "bridger address not found")
	}
	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle not found")
	}

	unsignedOutgoingBridgeCall := make([]*types.OutgoingBridgeCall, 0)
	k.IterateOutgoingBridgeCalls(ctx, func(outgoingBridgeCall *types.OutgoingBridgeCall) bool {
		if oracle.StartHeight > int64(outgoingBridgeCall.BlockHeight) {
			return false
		}
		if k.HasBridgeCallConfirm(ctx, outgoingBridgeCall.Nonce, oracleAddr) {
			return false
		}
		unsignedOutgoingBridgeCall = append(unsignedOutgoingBridgeCall, outgoingBridgeCall)
		return len(unsignedOutgoingBridgeCall) == types.MaxResults
	})
	return &types.QueryLastPendingBridgeCallByAddrResponse{BridgeCalls: unsignedOutgoingBridgeCall}, nil
}
