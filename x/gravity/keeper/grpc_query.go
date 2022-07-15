package keeper

import (
	"context"
	"sort"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/x/gravity/types"
)

var _ types.QueryServer = Keeper{}

// Params queries the params of the gravity module
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	var params types.Params
	k.paramSpace.GetParamSet(sdk.UnwrapSDKContext(c), &params)
	return &types.QueryParamsResponse{Params: params}, nil

}

// CurrentValset queries the CurrentValset of the gravity module
func (k Keeper) CurrentValset(c context.Context, _ *types.QueryCurrentValsetRequest) (*types.QueryCurrentValsetResponse, error) {
	return &types.QueryCurrentValsetResponse{Valset: k.GetCurrentValset(sdk.UnwrapSDKContext(c))}, nil
}

// ValsetRequest queries the ValsetRequest of the gravity module
func (k Keeper) ValsetRequest(c context.Context, req *types.QueryValsetRequestRequest) (*types.QueryValsetRequestResponse, error) {
	return &types.QueryValsetRequestResponse{Valset: k.GetValset(sdk.UnwrapSDKContext(c), req.Nonce)}, nil
}

// ValsetConfirm queries the ValsetConfirm of the gravity module
func (k Keeper) ValsetConfirm(c context.Context, req *types.QueryValsetConfirmRequest) (*types.QueryValsetConfirmResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "address invalid")
	}
	return &types.QueryValsetConfirmResponse{Confirm: k.GetValsetConfirm(sdk.UnwrapSDKContext(c), req.Nonce, addr)}, nil
}

// ValsetConfirmsByNonce queries the ValsetConfirmsByNonce of the gravity module
func (k Keeper) ValsetConfirmsByNonce(c context.Context, req *types.QueryValsetConfirmsByNonceRequest) (*types.QueryValsetConfirmsByNonceResponse, error) {
	var confirms []*types.MsgValsetConfirm
	k.IterateValsetConfirmByNonce(sdk.UnwrapSDKContext(c), req.Nonce, func(_ []byte, c types.MsgValsetConfirm) bool {
		confirms = append(confirms, &c)
		return false
	})
	return &types.QueryValsetConfirmsByNonceResponse{Confirms: confirms}, nil
}

// LastValsetRequests queries the LastValsetRequests of the gravity module
func (k Keeper) LastValsetRequests(c context.Context, _ *types.QueryLastValsetRequestsRequest) (*types.QueryLastValsetRequestsResponse, error) {
	valReq := k.GetValsets(sdk.UnwrapSDKContext(c))
	valReqLen := len(valReq)
	retLen := 0
	if valReqLen < maxValsetRequestsReturned {
		retLen = valReqLen
	} else {
		retLen = maxValsetRequestsReturned
	}
	return &types.QueryLastValsetRequestsResponse{Valsets: valReq[0:retLen]}, nil
}

// LastPendingValsetRequestByAddr queries the LastPendingValsetRequestByAddr of the gravity module
func (k Keeper) LastPendingValsetRequestByAddr(c context.Context, req *types.QueryLastPendingValsetRequestByAddrRequest) (*types.QueryLastPendingValsetRequestByAddrResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "address invalid")
	}

	var pendingValsetReq []*types.Valset
	k.IterateValsets(sdk.UnwrapSDKContext(c), func(_ []byte, val *types.Valset) bool {
		// foundConfirm is true if the operatorAddr has signed the valset we are currently looking at
		foundConfirm := k.GetValsetConfirm(sdk.UnwrapSDKContext(c), val.Nonce, addr) != nil
		// if this valset has NOT been signed by operatorAddr, store it in pendingValsetReq
		// and exit the loop
		if !foundConfirm {
			pendingValsetReq = append(pendingValsetReq, val)
		}
		// if we have more than 100 unconfirmed requests in
		// our array we should exit, pagination
		if len(pendingValsetReq) > 100 {
			return true
		}
		// return false to continue the loop
		return false
	})
	return &types.QueryLastPendingValsetRequestByAddrResponse{Valsets: pendingValsetReq}, nil
}

// BatchFees queries the batch fees from unbatched pool
func (k Keeper) BatchFees(c context.Context, req *types.QueryBatchFeeRequest) (*types.QueryBatchFeeResponse, error) {
	if req.MinBatchFees == nil {
		req.MinBatchFees = make([]types.MinBatchFee, 0)
	}
	return &types.QueryBatchFeeResponse{BatchFees: k.GetAllBatchFees(sdk.UnwrapSDKContext(c), req.MinBatchFees)}, nil
}

// LastPendingBatchRequestByAddr queries the LastPendingBatchRequestByAddr of the gravity module
func (k Keeper) LastPendingBatchRequestByAddr(c context.Context, req *types.QueryLastPendingBatchRequestByAddrRequest) (*types.QueryLastPendingBatchRequestByAddrResponse, error) {
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "address invalid")
	}

	var pendingBatchReq *types.OutgoingTxBatch
	k.IterateOutgoingTXBatches(sdk.UnwrapSDKContext(c), func(_ []byte, batch *types.OutgoingTxBatch) bool {
		foundConfirm := k.GetBatchConfirm(sdk.UnwrapSDKContext(c), batch.BatchNonce, batch.TokenContract, addr) != nil
		if !foundConfirm {
			pendingBatchReq = batch
			return true
		}
		return false
	})

	return &types.QueryLastPendingBatchRequestByAddrResponse{Batch: pendingBatchReq}, nil
}

// OutgoingTxBatches queries the OutgoingTxBatches of the gravity module
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

// BatchRequestByNonce queries the BatchRequestByNonce of the gravity module
func (k Keeper) BatchRequestByNonce(c context.Context, req *types.QueryBatchRequestByNonceRequest) (*types.QueryBatchRequestByNonceResponse, error) {
	if err := types.ValidateEthAddressAndValidateChecksum(req.TokenContract); err != nil {
		return nil, status.Error(codes.InvalidArgument, "token contract")
	}
	foundBatch := k.GetOutgoingTXBatch(sdk.UnwrapSDKContext(c), req.TokenContract, req.Nonce)
	if foundBatch == nil {
		return nil, status.Error(codes.NotFound, "tx batch")
	}
	return &types.QueryBatchRequestByNonceResponse{Batch: foundBatch}, nil
}

func (k Keeper) BatchConfirm(ctx context.Context, req *types.QueryBatchConfirmRequest) (*types.QueryBatchConfirmResponse, error) {
	orchestrator, err := sdk.AccAddressFromBech32(req.GetAddress())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "address")
	}
	confirm := k.GetBatchConfirm(sdk.UnwrapSDKContext(ctx), req.GetNonce(), req.GetTokenContract(), orchestrator)
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
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "address")
	}
	valAddr, found := k.GetOrchestratorValidator(ctx, addr)
	if !found {
		return nil, status.Error(codes.NotFound, "address")
	}
	lastEventNonce := k.GetLastEventNonceByValidator(ctx, valAddr)
	return &types.QueryLastEventNonceByAddrResponse{EventNonce: lastEventNonce}, nil
}

// DenomToERC20 queries the Cosmos Denom that maps to an Ethereum ERC20
func (k Keeper) DenomToERC20(c context.Context, req *types.QueryDenomToERC20Request) (*types.QueryDenomToERC20Response, error) {
	ctx := sdk.UnwrapSDKContext(c)
	fxOriginated, erc20, err := k.DenomToERC20Lookup(ctx, req.Denom)
	return &types.QueryDenomToERC20Response{Erc20: erc20, FxOriginated: fxOriginated}, err
}

// ERC20ToDenom queries the ERC20 contract that maps to an Ethereum ERC20 if any
func (k Keeper) ERC20ToDenom(c context.Context, req *types.QueryERC20ToDenomRequest) (*types.QueryERC20ToDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	fxOriginated, denom := k.ERC20ToDenomLookup(ctx, req.Erc20)
	return &types.QueryERC20ToDenomResponse{Denom: denom, FxOriginated: fxOriginated}, nil
}

func (k Keeper) GetDelegateKeyByValidator(c context.Context, req *types.QueryDelegateKeyByValidatorRequest) (*types.QueryDelegateKeyByValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	keys := k.GetDelegateKeys(ctx)
	reqValidator, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "address")
	}
	for _, key := range keys {
		keyValidator, err := sdk.ValAddressFromBech32(key.Validator)
		// this should be impossible due to the validate basic on the set orchestrator message
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if reqValidator.Equals(keyValidator) {
			return &types.QueryDelegateKeyByValidatorResponse{EthAddress: key.EthAddress, OrchestratorAddress: key.Orchestrator}, nil
		}

	}
	return nil, status.Error(codes.NotFound, "validator")
}

func (k Keeper) GetDelegateKeyByOrchestrator(c context.Context, req *types.QueryDelegateKeyByOrchestratorRequest) (*types.QueryDelegateKeyByOrchestratorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	keys := k.GetDelegateKeys(ctx)
	reqOrchestrator, err := sdk.AccAddressFromBech32(req.OrchestratorAddress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "address")
	}
	for _, key := range keys {
		keyOrchestrator, err := sdk.AccAddressFromBech32(key.Orchestrator)
		// this should be impossible due to the validate basic on the set orchestrator message
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if reqOrchestrator.Equals(keyOrchestrator) {
			return &types.QueryDelegateKeyByOrchestratorResponse{ValidatorAddress: key.Validator, EthAddress: key.EthAddress}, nil
		}

	}
	return nil, status.Error(codes.NotFound, "validator")
}

func (k Keeper) GetDelegateKeyByEth(c context.Context, req *types.QueryDelegateKeyByEthRequest) (*types.QueryDelegateKeyByEthResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	keys := k.GetDelegateKeys(ctx)
	if err := types.ValidateEthAddressAndValidateChecksum(req.EthAddress); err != nil {
		return nil, status.Error(codes.InvalidArgument, "address")
	}
	for _, key := range keys {
		if req.EthAddress == key.EthAddress {
			return &types.QueryDelegateKeyByEthResponse{ValidatorAddress: key.Validator, OrchestratorAddress: key.Orchestrator}, nil
		}

	}
	return nil, status.Error(codes.NotFound, "validator")
}

func (k Keeper) GetPendingSendToEth(c context.Context, req *types.QueryPendingSendToEthRequest) (*types.QueryPendingSendToEthResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	batches := k.GetOutgoingTxBatches(ctx)
	unbatchedTx := k.GetPoolTransactions(ctx)
	senderAddress := req.SenderAddress
	var res = &types.QueryPendingSendToEthResponse{
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
	return &types.QueryIbcSequenceHeightResponse{Found: found, Height: height}, nil
}

func (k Keeper) LastObservedBlockHeight(c context.Context, _ *types.QueryLastObservedBlockHeightRequest) (*types.QueryLastObservedBlockHeightResponse, error) {
	blockHeight := k.GetLastObservedEthBlockHeight(sdk.UnwrapSDKContext(c))
	return &types.QueryLastObservedBlockHeightResponse{BlockHeight: blockHeight.FxBlockHeight, EthBlockHeight: blockHeight.EthBlockHeight}, nil
}

func (k Keeper) LastEventBlockHeightByAddr(c context.Context, req *types.QueryLastEventBlockHeightByAddrRequest) (*types.QueryLastEventBlockHeightByAddrResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "address")
	}
	valAddr, found := k.GetOrchestratorValidator(ctx, addr)
	if !found {
		return nil, status.Error(codes.NotFound, "address")
	}
	lastEventBlockHeight := k.getLastEventBlockHeightByValidator(ctx, valAddr)
	return &types.QueryLastEventBlockHeightByAddrResponse{BlockHeight: lastEventBlockHeight}, nil
}

func (k Keeper) ProjectedBatchTimeoutHeight(c context.Context, _ *types.QueryProjectedBatchTimeoutHeightRequest) (*types.QueryProjectedBatchTimeoutHeightResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	timeout := k.GetBatchTimeoutHeight(ctx)
	return &types.QueryProjectedBatchTimeoutHeightResponse{TimeoutHeight: timeout}, nil
}

func (k Keeper) BridgeTokens(c context.Context, _ *types.QueryBridgeTokensRequest) (*types.QueryBridgeTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var bridgeTokens = make([]*types.ERC20ToDenom, 0)
	k.IterateERC20ToDenom(ctx, func(bytes []byte, token *types.ERC20ToDenom) bool {
		bridgeTokens = append(bridgeTokens, token)
		return false
	})
	return &types.QueryBridgeTokensResponse{BridgeTokens: bridgeTokens}, nil
}
