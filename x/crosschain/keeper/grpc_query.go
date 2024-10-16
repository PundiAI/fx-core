package keeper

import (
	"context"
	"sort"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

var _ types.QueryServer = QueryServer{}

type QueryServer struct {
	Keeper
}

func NewQueryServerImpl(keeper Keeper) types.QueryServer {
	return &QueryServer{Keeper: keeper}
}

func (k QueryServer) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := k.GetParams(sdk.UnwrapSDKContext(c))
	return &types.QueryParamsResponse{Params: params}, nil
}

func (k QueryServer) CurrentOracleSet(c context.Context, _ *types.QueryCurrentOracleSetRequest) (*types.QueryCurrentOracleSetResponse, error) {
	return &types.QueryCurrentOracleSetResponse{OracleSet: k.GetCurrentOracleSet(sdk.UnwrapSDKContext(c))}, nil
}

func (k QueryServer) OracleSetRequest(c context.Context, req *types.QueryOracleSetRequestRequest) (*types.QueryOracleSetRequestResponse, error) {
	return &types.QueryOracleSetRequestResponse{OracleSet: k.GetOracleSet(sdk.UnwrapSDKContext(c), req.Nonce)}, nil
}

func (k QueryServer) OracleSetConfirm(c context.Context, req *types.QueryOracleSetConfirmRequest) (*types.QueryOracleSetConfirmResponse, error) {
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	ctx := sdk.UnwrapSDKContext(c)
	oracleAddr, err := k.BridgeAddrToOracleAddr(ctx, req.GetBridgerAddress())
	if err != nil {
		return nil, err
	}
	return &types.QueryOracleSetConfirmResponse{Confirm: k.GetOracleSetConfirm(ctx, req.Nonce, oracleAddr)}, nil
}

func (k QueryServer) OracleSetConfirmsByNonce(c context.Context, req *types.QueryOracleSetConfirmsByNonceRequest) (*types.QueryOracleSetConfirmsByNonceResponse, error) {
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

func (k QueryServer) LastOracleSetRequests(c context.Context, _ *types.QueryLastOracleSetRequestsRequest) (*types.QueryLastOracleSetRequestsResponse, error) {
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

func (k QueryServer) LastPendingOracleSetRequestByAddr(c context.Context, req *types.QueryLastPendingOracleSetRequestByAddrRequest) (*types.QueryLastPendingOracleSetRequestByAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	oracle, err := k.BridgeAddrToOracle(ctx, req.GetBridgerAddress())
	if err != nil {
		return nil, err
	}
	oracleAddr := oracle.GetOracle()
	var pendingOracleSetReq []*types.OracleSet
	k.IterateOracleSets(ctx, false, func(oracleSet *types.OracleSet) bool {
		if oracle.StartHeight > int64(oracleSet.Height) {
			return false
		}
		// found is true if the operatorAddr has signed the oracle set we are currently looking at
		// if this oracle set has NOT been signed by oracleAddr, store it in pendingOracleSetReq and exit the loop
		if found := k.GetOracleSetConfirm(ctx, oracleSet.Nonce, oracleAddr) != nil; !found {
			pendingOracleSetReq = append(pendingOracleSetReq, oracleSet)
		}
		// if we have more than 100 unconfirmed requests in
		// our array we should exit, pagination
		return len(pendingOracleSetReq) == types.MaxResults
	})
	return &types.QueryLastPendingOracleSetRequestByAddrResponse{OracleSets: pendingOracleSetReq}, nil
}

func (k QueryServer) LastPendingBatchRequestByAddr(c context.Context, req *types.QueryLastPendingBatchRequestByAddrRequest) (*types.QueryLastPendingBatchRequestByAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	oracle, err := k.BridgeAddrToOracle(ctx, req.GetBridgerAddress())
	if err != nil {
		return nil, err
	}
	oracleAddr := oracle.GetOracle()
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

func (k QueryServer) OutgoingTxBatches(c context.Context, _ *types.QueryOutgoingTxBatchesRequest) (*types.QueryOutgoingTxBatchesResponse, error) {
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

func (k QueryServer) OutgoingTxBatch(c context.Context, req *types.QueryOutgoingTxBatchRequest) (*types.QueryOutgoingTxBatchResponse, error) {
	if err := types.ValidateExternalAddr(req.ChainName, req.GetTokenContract()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "token contract address")
	}
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	foundBatch := k.GetOutgoingTxBatch(sdk.UnwrapSDKContext(c), req.TokenContract, req.Nonce)
	if foundBatch == nil {
		return nil, status.Error(codes.NotFound, "tx batch")
	}
	return &types.QueryOutgoingTxBatchResponse{Batch: foundBatch}, nil
}

func (k QueryServer) BatchConfirm(c context.Context, req *types.QueryBatchConfirmRequest) (*types.QueryBatchConfirmResponse, error) {
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	ctx := sdk.UnwrapSDKContext(c)
	oracleAddr, err := k.BridgeAddrToOracleAddr(ctx, req.GetBridgerAddress())
	if err != nil {
		return nil, err
	}
	confirm := k.GetBatchConfirm(ctx, req.TokenContract, req.Nonce, oracleAddr)
	return &types.QueryBatchConfirmResponse{Confirm: confirm}, nil
}

// BatchConfirms returns the batch confirmations by nonce and token contract
func (k QueryServer) BatchConfirms(c context.Context, req *types.QueryBatchConfirmsRequest) (*types.QueryBatchConfirmsResponse, error) {
	if err := types.ValidateExternalAddr(req.ChainName, req.GetTokenContract()); err != nil {
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
func (k QueryServer) LastEventNonceByAddr(c context.Context, req *types.QueryLastEventNonceByAddrRequest) (*types.QueryLastEventNonceByAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	oracleAddr, err := k.BridgeAddrToOracleAddr(ctx, req.GetBridgerAddress())
	if err != nil {
		return nil, err
	}
	lastEventNonce := k.GetLastEventNonceByOracle(ctx, oracleAddr)
	return &types.QueryLastEventNonceByAddrResponse{EventNonce: lastEventNonce}, nil
}

func (k QueryServer) DenomToToken(c context.Context, req *types.QueryDenomToTokenRequest) (*types.QueryDenomToTokenResponse, error) {
	if len(req.GetDenom()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "denom")
	}

	ctx := sdk.UnwrapSDKContext(c)
	bridgeToken, err := k.erc20Keeper.GetBridgeToken(ctx, req.Denom, req.ChainName)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &types.QueryDenomToTokenResponse{Token: bridgeToken.Contract}, nil
}

func (k QueryServer) TokenToDenom(c context.Context, req *types.QueryTokenToDenomRequest) (*types.QueryTokenToDenomResponse, error) {
	if err := types.ValidateExternalAddr(req.ChainName, req.GetToken()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "token address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	baseDenom, err := k.erc20Keeper.GetBaseDenom(ctx, req.Token)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	_, err = k.erc20Keeper.GetBridgeToken(ctx, baseDenom, req.ChainName)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &types.QueryTokenToDenomResponse{Denom: baseDenom}, nil
}

func (k QueryServer) GetOracleByAddr(c context.Context, req *types.QueryOracleByAddrRequest) (*types.QueryOracleResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(req.OracleAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "oracle address")
	}
	oracle, found := k.GetOracle(sdk.UnwrapSDKContext(c), oracleAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "oracle not found")
	}
	return &types.QueryOracleResponse{Oracle: &oracle}, nil
}

func (k QueryServer) GetOracleByBridgerAddr(c context.Context, req *types.QueryOracleByBridgerAddrRequest) (*types.QueryOracleResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	oracle, err := k.BridgeAddrToOracle(ctx, req.GetBridgerAddress())
	if err != nil {
		return nil, err
	}
	return &types.QueryOracleResponse{Oracle: &oracle}, nil
}

func (k QueryServer) GetOracleByExternalAddr(c context.Context, req *types.QueryOracleByExternalAddrRequest) (*types.QueryOracleResponse, error) {
	if err := types.ValidateExternalAddr(req.ChainName, req.GetExternalAddress()); err != nil {
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

func (k QueryServer) LastObservedBlockHeight(c context.Context, _ *types.QueryLastObservedBlockHeightRequest) (*types.QueryLastObservedBlockHeightResponse, error) {
	blockHeight := k.GetLastObservedBlockHeight(sdk.UnwrapSDKContext(c))
	return &types.QueryLastObservedBlockHeightResponse{
		ExternalBlockHeight: blockHeight.ExternalBlockHeight,
		BlockHeight:         blockHeight.BlockHeight,
	}, nil
}

func (k QueryServer) LastEventBlockHeightByAddr(c context.Context, req *types.QueryLastEventBlockHeightByAddrRequest) (*types.QueryLastEventBlockHeightByAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	oracleAddr, err := k.BridgeAddrToOracleAddr(ctx, req.GetBridgerAddress())
	if err != nil {
		return nil, err
	}

	lastEventBlockHeight := k.GetLastEventBlockHeightByOracle(ctx, oracleAddr)
	return &types.QueryLastEventBlockHeightByAddrResponse{BlockHeight: lastEventBlockHeight}, nil
}

func (k QueryServer) Oracles(c context.Context, _ *types.QueryOraclesRequest) (*types.QueryOraclesResponse, error) {
	oracles := k.GetAllOracles(sdk.UnwrapSDKContext(c), false)
	return &types.QueryOraclesResponse{Oracles: oracles}, nil
}

func (k QueryServer) ProjectedBatchTimeoutHeight(c context.Context, _ *types.QueryProjectedBatchTimeoutHeightRequest) (*types.QueryProjectedBatchTimeoutHeightResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	timeout := k.CalExternalTimeoutHeight(ctx, GetExternalBatchTimeout)
	return &types.QueryProjectedBatchTimeoutHeightResponse{TimeoutHeight: timeout}, nil
}

func (k QueryServer) BridgeTokens(c context.Context, _ *types.QueryBridgeTokensRequest) (*types.QueryBridgeTokensResponse, error) {
	bridgeTokens := make([]*types.BridgeToken, 0)
	// todo: need implement
	return &types.QueryBridgeTokensResponse{BridgeTokens: bridgeTokens}, nil
}

func (k QueryServer) BridgeCoinByDenom(c context.Context, req *types.QueryBridgeCoinByDenomRequest) (*types.QueryBridgeCoinByDenomResponse, error) {
	if len(req.GetDenom()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "denom")
	}
	supply, err := k.BridgeCoinSupply(sdk.UnwrapSDKContext(c), req.GetDenom(), req.GetChainName())
	if err != nil {
		return nil, status.Error(codes.NotFound, "denom")
	}
	return &types.QueryBridgeCoinByDenomResponse{Coin: supply}, nil
}

func (k QueryServer) BridgeChainList(_ context.Context, _ *types.QueryBridgeChainListRequest) (*types.QueryBridgeChainListResponse, error) {
	return &types.QueryBridgeChainListResponse{ChainNames: types.GetSupportChains()}, nil
}

func (k QueryServer) BridgeCalls(c context.Context, req *types.QueryBridgeCallsRequest) (*types.QueryBridgeCallsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	pendingStore := prefix.NewStore(store, types.OutgoingBridgeCallNonceKey)

	datas, pageRes, err := query.GenericFilteredPaginate(k.cdc, pendingStore, req.Pagination, func(key []byte, data *types.OutgoingBridgeCall) (*types.OutgoingBridgeCall, error) {
		return data, nil
	}, func() *types.OutgoingBridgeCall {
		return &types.OutgoingBridgeCall{}
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryBridgeCallsResponse{BridgeCalls: datas, Pagination: pageRes}, nil
}

func (k QueryServer) BridgeCallConfirmByNonce(c context.Context, req *types.QueryBridgeCallConfirmByNonceRequest) (*types.QueryBridgeCallConfirmByNonceResponse, error) {
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
	enoughPower := sdkmath.NewIntFromUint64(confirmPowers).GTE(requiredPower)
	return &types.QueryBridgeCallConfirmByNonceResponse{Confirms: bridgeCallConfirms, EnoughPower: enoughPower}, nil
}

func (k QueryServer) BridgeCallByNonce(c context.Context, req *types.QueryBridgeCallByNonceRequest) (*types.QueryBridgeCallByNonceResponse, error) {
	outgoingBridgeCall, found := k.GetOutgoingBridgeCallByNonce(sdk.UnwrapSDKContext(c), req.GetNonce())
	if !found {
		return nil, status.Error(codes.NotFound, "outgoing bridge call not found")
	}
	return &types.QueryBridgeCallByNonceResponse{BridgeCall: outgoingBridgeCall}, nil
}

func (k QueryServer) BridgeCallBySender(c context.Context, req *types.QueryBridgeCallBySenderRequest) (*types.QueryBridgeCallBySenderResponse, error) {
	var outgoingBridgeCalls []*types.OutgoingBridgeCall
	k.IterateOutgoingBridgeCallsByAddress(sdk.UnwrapSDKContext(c), req.GetSenderAddress(), func(outgoingBridgeCall *types.OutgoingBridgeCall) bool {
		outgoingBridgeCalls = append(outgoingBridgeCalls, outgoingBridgeCall)
		return false
	})
	return &types.QueryBridgeCallBySenderResponse{BridgeCalls: outgoingBridgeCalls}, nil
}

func (k QueryServer) LastPendingBridgeCallByAddr(c context.Context, req *types.QueryLastPendingBridgeCallByAddrRequest) (*types.QueryLastPendingBridgeCallByAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	oracle, err := k.BridgeAddrToOracle(ctx, req.GetBridgerAddress())
	if err != nil {
		return nil, err
	}

	oracleAddr := oracle.GetOracle()
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

func (k QueryServer) BridgeAddrToOracleAddr(ctx sdk.Context, bridgeAddr string) (sdk.AccAddress, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(bridgeAddr)
	if err != nil {
		return sdk.AccAddress{}, status.Error(codes.InvalidArgument, "bridger address")
	}
	oracleAddr, found := k.GetOracleAddrByBridgerAddr(ctx, bridgerAddr)
	if !found {
		return sdk.AccAddress{}, status.Error(codes.NotFound, "oracle not found by bridger address")
	}
	return oracleAddr, nil
}

func (k QueryServer) BridgeAddrToOracle(ctx sdk.Context, bridgeAddr string) (types.Oracle, error) {
	oracleAddr, err := k.BridgeAddrToOracleAddr(ctx, bridgeAddr)
	if err != nil {
		return types.Oracle{}, err
	}
	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		return types.Oracle{}, status.Error(codes.NotFound, "oracle not found")
	}
	return oracle, nil
}

func (k QueryServer) PendingExecuteClaim(c context.Context, req *types.QueryPendingExecuteClaimRequest) (*types.QueryPendingExecuteClaimResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := prefix.NewStore(ctx.KVStore(k.Keeper.storeKey), types.PendingExecuteClaimKey)
	var claims []*codectypes.Any
	pageRes, err := query.Paginate(store, req.Pagination, func(key, value []byte) error {
		var claim types.ExternalClaim
		if err := k.cdc.UnmarshalInterface(value, &claim); err != nil {
			return err
		}
		anyClaim, err := codectypes.NewAnyWithValue(claim)
		if err != nil {
			return err
		}
		claims = append(claims, anyClaim)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "paginate: %v", err)
	}
	return &types.QueryPendingExecuteClaimResponse{Claims: claims, Pagination: pageRes}, nil
}
