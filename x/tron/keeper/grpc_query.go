package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	"github.com/functionx/fx-core/v7/x/tron/types"
)

var _ crosschaintypes.QueryServer = Keeper{}

// BatchFees queries the batch fees from unbatched pool
func (k Keeper) BatchFees(c context.Context, req *crosschaintypes.QueryBatchFeeRequest) (*crosschaintypes.QueryBatchFeeResponse, error) {
	if req.GetMinBatchFees() == nil {
		req.MinBatchFees = make([]crosschaintypes.MinBatchFee, 0)
	}
	for _, fee := range req.MinBatchFees {
		if fee.BaseFee.IsNil() || fee.BaseFee.IsNegative() {
			return nil, status.Error(codes.InvalidArgument, "base fee")
		}
		if err := types.ValidateTronAddress(fee.TokenContract); err != nil {
			return nil, status.Error(codes.InvalidArgument, "token contract")
		}
	}
	allBatchFees := k.GetAllBatchFees(sdk.UnwrapSDKContext(c), crosschaintypes.MaxResults, req.MinBatchFees)
	return &crosschaintypes.QueryBatchFeeResponse{BatchFees: allBatchFees}, nil
}

// BatchRequestByNonce queries the BatchRequestByNonce of the bsc module
func (k Keeper) BatchRequestByNonce(c context.Context, req *crosschaintypes.QueryBatchRequestByNonceRequest) (*crosschaintypes.QueryBatchRequestByNonceResponse, error) {
	if err := types.ValidateTronAddress(req.GetTokenContract()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "token contract address")
	}
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	foundBatch := k.GetOutgoingTxBatch(sdk.UnwrapSDKContext(c), req.TokenContract, req.Nonce)
	if foundBatch == nil {
		return nil, status.Error(codes.NotFound, "tx batch")
	}
	return &crosschaintypes.QueryBatchRequestByNonceResponse{Batch: foundBatch}, nil
}

// BatchConfirms returns the batch confirmations by nonce and token contract
func (k Keeper) BatchConfirms(c context.Context, req *crosschaintypes.QueryBatchConfirmsRequest) (*crosschaintypes.QueryBatchConfirmsResponse, error) {
	if err := types.ValidateTronAddress(req.GetTokenContract()); err != nil {
		return nil, status.Error(codes.InvalidArgument, "token contract address")
	}
	if req.GetNonce() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "nonce")
	}
	var confirms []*crosschaintypes.MsgConfirmBatch
	k.IterateBatchConfirmByNonceAndTokenContract(sdk.UnwrapSDKContext(c), req.Nonce, req.TokenContract, func(confirm *crosschaintypes.MsgConfirmBatch) bool {
		confirms = append(confirms, confirm)
		return false
	})
	return &crosschaintypes.QueryBatchConfirmsResponse{Confirms: confirms}, nil
}

func (k Keeper) TokenToDenom(c context.Context, req *crosschaintypes.QueryTokenToDenomRequest) (*crosschaintypes.QueryTokenToDenomResponse, error) {
	if err := types.ValidateTronAddress(req.Token); err != nil {
		return nil, status.Error(codes.InvalidArgument, "token address")
	}
	bridgeToken := k.GetBridgeTokenDenom(sdk.UnwrapSDKContext(c), req.Token)
	if bridgeToken == nil {
		return nil, status.Error(codes.NotFound, "bridge token")
	}
	return &crosschaintypes.QueryTokenToDenomResponse{
		Denom: bridgeToken.Denom,
	}, nil
}

func (k Keeper) GetOracleByExternalAddr(c context.Context, req *crosschaintypes.QueryOracleByExternalAddrRequest) (*crosschaintypes.QueryOracleResponse, error) {
	if err := types.ValidateTronAddress(req.GetExternalAddress()); err != nil {
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
	return &crosschaintypes.QueryOracleResponse{Oracle: &oracle}, nil
}
