package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	crosschainkeeper "github.com/functionx/fx-core/x/crosschain/keeper"

	"github.com/functionx/fx-core/x/crosschain/types"
)

var _ types.QueryServer = Keeper{}

// BatchFees queries the batch fees from unbatched pool
func (k Keeper) BatchFees(c context.Context, req *types.QueryBatchFeeRequest) (*types.QueryBatchFeeResponse, error) {
	if req.GetMinBatchFees() == nil {
		req.MinBatchFees = make([]types.MinBatchFee, 0)
	}
	for _, fee := range req.MinBatchFees {
		if fee.BaseFee.IsNil() || fee.BaseFee.IsNegative() {
			return nil, sdkerrors.Wrap(types.ErrInvalid, "base fee")
		}
		if err := types.ValidateEthereumAddress(fee.TokenContract); err != nil {
			return nil, sdkerrors.Wrap(types.ErrInvalid, "token contract")
		}
	}
	allBatchFees := k.GetAllBatchFees(sdk.UnwrapSDKContext(c), crosschainkeeper.MaxResults, req.MinBatchFees)
	return &types.QueryBatchFeeResponse{BatchFees: allBatchFees}, nil
}

// BatchRequestByNonce queries the BatchRequestByNonce of the bsc module
func (k Keeper) BatchRequestByNonce(c context.Context, req *types.QueryBatchRequestByNonceRequest) (*types.QueryBatchRequestByNonceResponse, error) {
	if err := types.ValidateEthereumAddress(req.GetTokenContract()); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "token contract address")
	}
	if req.GetNonce() <= 0 {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "nonce")
	}
	foundBatch := k.GetOutgoingTXBatch(sdk.UnwrapSDKContext(c), req.TokenContract, req.Nonce)
	if foundBatch == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "can not find tx batch")
	}
	return &types.QueryBatchRequestByNonceResponse{Batch: foundBatch}, nil
}

// BatchConfirms returns the batch confirmations by nonce and token contract
func (k Keeper) BatchConfirms(c context.Context, req *types.QueryBatchConfirmsRequest) (*types.QueryBatchConfirmsResponse, error) {
	if err := types.ValidateEthereumAddress(req.GetTokenContract()); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "token contract address")
	}
	if req.GetNonce() <= 0 {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "nonce")
	}
	var confirms []*types.MsgConfirmBatch
	k.IterateBatchConfirmByNonceAndTokenContract(sdk.UnwrapSDKContext(c), req.Nonce, req.TokenContract, func(_ []byte, c types.MsgConfirmBatch) bool {
		confirms = append(confirms, &c)
		return false
	})
	return &types.QueryBatchConfirmsResponse{Confirms: confirms}, nil
}

func (k Keeper) TokenToDenom(c context.Context, req *types.QueryTokenToDenomRequest) (*types.QueryTokenToDenomResponse, error) {
	if err := types.ValidateEthereumAddress(req.GetToken()); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "token address")
	}
	bridgeToken := k.GetBridgeTokenDenom(sdk.UnwrapSDKContext(c), req.Token)
	if bridgeToken == nil {
		return nil, sdkerrors.Wrap(types.ErrEmpty, "bridge token is not exist")
	}
	return &types.QueryTokenToDenomResponse{
		Denom:      bridgeToken.Denom,
		ChannelIbc: bridgeToken.ChannelIbc,
	}, nil
}

func (k Keeper) GetOracleByExternalAddr(c context.Context, req *types.QueryOracleByExternalAddrRequest) (*types.QueryOracleResponse, error) {
	if err := types.ValidateEthereumAddress(req.GetExternalAddress()); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "external address")
	}

	ctx := sdk.UnwrapSDKContext(c)
	oracleAddr, found := k.GetOracleByExternalAddress(ctx, req.ExternalAddress)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrNoFoundOracle, "by external address: %s", req.ExternalAddress)
	}
	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNoFoundOracle, oracleAddr.String())
	}
	return &types.QueryOracleResponse{Oracle: &oracle}, nil
}
