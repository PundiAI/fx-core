package keeper

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/functionx/fx-core/x/crosschain/types"
)

const MaxResults = 100
const maxOracleSetRequestsReturned = 5

const (

	// OracleSets

	// QueryOracleSetRequest This retrieves a specific validator set by it's nonce
	// used to compare what's on Ethereum with what's in Cosmos
	// to perform slashing / validation of system consistency
	QueryOracleSetRequest = "oracleSetRequest"
	// QueryOracleSetConfirmsByNonce Gets all the confirmation signatures for a given validator
	// set, used by the relayer to package the validator set and
	// it's signatures into an Ethereum transaction
	QueryOracleSetConfirmsByNonce = "oracleSetConfirms"
	// QueryLastOracleSetRequests Gets the last N (where N is currently 5) validator sets that
	// have been produced by the chain. Useful to see if any recently
	// signed requests can be submitted.
	QueryLastOracleSetRequests = "lastOracleSetRequests"
	// QueryLastPendingOracleSetRequestByAddr Gets a list of unsigned oracleSets for a given validators delegate
	// orchestrator address. Up to 100 are sent at a time
	QueryLastPendingOracleSetRequestByAddr = "lastPendingOracleSetRequest"

	QueryCurrentOracleSet = "currentOracleSet"
	QueryOracleSetConfirm = "oracleSetConfirm"

	// QueryGravityID used by the contract deployer script. GravityID is set in the Genesis
	// file, then read by the contract deployer and deployed to Ethereum
	// a unique GravityID ensures that even if the same validator set with
	// the same keys is running on two chains these chains can have independent
	// bridges
	QueryGravityID = "gravityID"

	// Batches
	// note the current logic here constrains batch throughput to one
	// batch (of any type) per Cosmos block.

	// QueryBatch This retrieves a specific batch by it's nonce and token contract
	// or in the case of a Cosmos originated address it's denom
	QueryBatch = "batch"
	// QueryLastPendingBatchRequestByAddr Get the last unsigned batch (of any denom) for the validators
	// orchestrator to sign
	QueryLastPendingBatchRequestByAddr = "lastPendingBatchRequest"
	// QueryOutgoingTxBatches gets the last 100 outgoing batches, regardless of denom, useful
	// for a relayed to see what is available to relay
	QueryOutgoingTxBatches = "lastBatches"
	// QueryBatchConfirms Used by the relayer to package a batch with signatures required
	// to submit to Ethereum
	QueryBatchConfirms = "batchConfirms"
	// QueryBatchFees Used to query all pending SendToEth transactions and fees available for each
	// token type, a relayer can then estimate their potential profit when requesting
	// a batch
	QueryBatchFees = "batchFees"

	// QueryTokenToDenom Token mapping
	// This retrieves the denom which is represented by a given ERC20 contract
	QueryTokenToDenom = "TokenToDenom"
	// QueryDenomToToken This retrieves the ERC20 contract which represents a given denom
	QueryDenomToToken = "DenomToToken"

	// QueryPendingSendToExternal Query pending transactions
	QueryPendingSendToExternal = "PendingSendToExternal"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		if len(path) <= 0 {
			return nil, sdkerrors.ErrInvalidRequest
		}
		switch path[0] {
		// OracleSets
		case QueryCurrentOracleSet:
			return queryCurrentOracleSet(ctx, keeper)
		case QueryOracleSetRequest:
			if len(path) != 2 {
				return nil, sdkerrors.ErrInvalidRequest
			}
			return queryOracleSetRequest(ctx, path[1], keeper)
		case QueryOracleSetConfirm:
			if len(path) != 3 {
				return nil, sdkerrors.ErrInvalidRequest
			}
			return queryOracleSetConfirm(ctx, path[1], path[2], keeper)
		case QueryOracleSetConfirmsByNonce:
			if len(path) != 2 {
				return nil, sdkerrors.ErrInvalidRequest
			}
			return queryAllOracleSetConfirms(ctx, path[1], keeper)
		case QueryLastOracleSetRequests:
			return lastOracleSetRequests(ctx, keeper)
		case QueryLastPendingOracleSetRequestByAddr:
			if len(path) != 2 {
				return nil, sdkerrors.ErrInvalidRequest
			}
			return lastPendingOracleSetRequest(ctx, path[1], keeper)

		// Batches
		case QueryBatch:
			if len(path) != 3 {
				return nil, sdkerrors.ErrInvalidRequest
			}
			return queryBatch(ctx, path[1], path[2], keeper)
		case QueryBatchConfirms:
			if len(path) != 3 {
				return nil, sdkerrors.ErrInvalidRequest
			}
			return queryAllBatchConfirms(ctx, path[1], path[2], keeper)
		case QueryLastPendingBatchRequestByAddr:
			if len(path) != 2 {
				return nil, sdkerrors.ErrInvalidRequest
			}
			return lastPendingBatchRequest(ctx, path[1], keeper)
		case QueryOutgoingTxBatches:
			return lastBatchesRequest(ctx, keeper)
		case QueryBatchFees:
			return queryBatchFees(ctx, keeper, req, legacyQuerierCdc)

		case QueryGravityID:
			return queryGravityID(ctx, keeper)

		// Token mappings
		case QueryDenomToToken:
			if len(path) != 2 {
				return nil, sdkerrors.ErrInvalidRequest
			}
			return queryDenomToToken(ctx, path[1], keeper)
		case QueryTokenToDenom:
			if len(path) != 2 {
				return nil, sdkerrors.ErrInvalidRequest
			}
			return queryTokenToDenom(ctx, path[1], keeper)

		// Pending transactions
		case QueryPendingSendToExternal:
			if len(path) != 2 {
				return nil, sdkerrors.ErrInvalidRequest
			}
			return queryPendingSendToExternal(ctx, path[1], keeper)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint", keeper.moduleName)
		}
	}
}

func queryOracleSetRequest(ctx sdk.Context, nonceStr string, keeper Keeper) ([]byte, error) {
	nonce, err := uint64FromString(nonceStr)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "nonce")
	}

	oracleSet := keeper.GetOracleSet(ctx, nonce)
	if oracleSet == nil {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, oracleSet)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// allOracleSetConfirmsByNonce returns all the confirm messages for a given nonce
// When nothing found an empty json array is returned. No pagination.
func queryAllOracleSetConfirms(ctx sdk.Context, nonceStr string, keeper Keeper) ([]byte, error) {
	nonce, err := uint64FromString(nonceStr)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "nonce")
	}

	var confirms []*types.MsgOracleSetConfirm
	keeper.IterateOracleSetConfirmByNonce(ctx, nonce, func(_ []byte, c types.MsgOracleSetConfirm) bool {
		confirms = append(confirms, &c)
		return false
	})
	if len(confirms) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, confirms)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// allBatchConfirms returns all the confirm messages for a given nonce
// When nothing found an empty json array is returned. No pagination.
func queryAllBatchConfirms(ctx sdk.Context, nonceStr string, tokenContract string, keeper Keeper) ([]byte, error) {
	nonce, err := uint64FromString(nonceStr)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "nonce")
	}
	if len(tokenContract) <= 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "token contract")
	}

	var confirms []types.MsgConfirmBatch
	keeper.IterateBatchConfirmByNonceAndTokenContract(ctx, nonce, tokenContract, func(_ []byte, c types.MsgConfirmBatch) bool {
		confirms = append(confirms, c)
		return false
	})
	if len(confirms) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, confirms)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// lastOracleSetRequests returns up to maxOracleSetRequestsReturned oracleSets from the store
func lastOracleSetRequests(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var valReq []*types.OracleSet
	keeper.IterateOracleSets(ctx, func(_ []byte, val *types.OracleSet) bool {
		valReq = append(valReq, val)
		return false
	})
	valReqLen := len(valReq)
	retLen := 0
	if valReqLen < maxOracleSetRequestsReturned {
		retLen = valReqLen
	} else {
		retLen = maxOracleSetRequestsReturned
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, valReq[0:retLen])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// lastPendingOracleSetRequest gets a list of validator sets that this validator has not signed
// limited by 100 sets per request.
func lastPendingOracleSetRequest(ctx sdk.Context, bridgerAddr string, keeper Keeper) ([]byte, error) {
	bridger, err := sdk.AccAddressFromBech32(bridgerAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "bridger address")
	}

	var pendingOracleSetReq []*types.OracleSet
	keeper.IterateOracleSets(ctx, func(_ []byte, val *types.OracleSet) bool {
		// foundConfirm is true if the operatorAddr has signed the oracleSet we are currently looking at
		foundConfirm := keeper.GetOracleSetConfirm(ctx, val.Nonce, bridger) != nil
		// if this oracleSet has NOT been signed by operatorAddr, store it in pendingOracleSetReq
		// and exit the loop
		if !foundConfirm {
			pendingOracleSetReq = append(pendingOracleSetReq, val)
		}
		// if we have more than 100 unconfirmed requests in
		// our array we should exit, pagination
		if len(pendingOracleSetReq) > MaxResults {
			return true
		}
		// return false to continue the loop
		return false
	})
	if len(pendingOracleSetReq) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pendingOracleSetReq)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryCurrentOracleSet(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	oracleSet := keeper.GetCurrentOracleSet(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, oracleSet)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// queryOracleSetConfirm returns the confirm msg for single orchestrator address and nonce
// When nothing found a nil value is returned
func queryOracleSetConfirm(ctx sdk.Context, nonceStr, bridgerAddr string, keeper Keeper) ([]byte, error) {
	nonce, err := uint64FromString(nonceStr)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "nonce")
	}

	bridger, err := sdk.AccAddressFromBech32(bridgerAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "bridger address")
	}

	oracleSet := keeper.GetOracleSetConfirm(ctx, nonce, bridger)
	if oracleSet == nil {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, *oracleSet)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

type MultiSigUpdateResponse struct {
	OracleSet  types.OracleSet `json:"oracleSet"`
	Signatures [][]byte        `json:"signatures,omitempty"`
}

// lastPendingBatchRequest gets the latest batch that has NOT been signed by operatorAddr
func lastPendingBatchRequest(ctx sdk.Context, bridgerAddr string, keeper Keeper) ([]byte, error) {
	bridger, err := sdk.AccAddressFromBech32(bridgerAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "bridger address")
	}

	var pendingBatchReq *types.OutgoingTxBatch
	keeper.IterateOutgoingTxBatches(ctx, func(_ []byte, batch *types.OutgoingTxBatch) bool {
		foundConfirm := keeper.GetBatchConfirm(ctx, batch.BatchNonce, batch.TokenContract, bridger) != nil
		if !foundConfirm {
			pendingBatchReq = batch
			return true
		}
		return false
	})
	if pendingBatchReq == nil {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pendingBatchReq)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// Gets MaxResults batches from store. Does not select by token type or anything
func lastBatchesRequest(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var batches []*types.OutgoingTxBatch
	keeper.IterateOutgoingTxBatches(ctx, func(_ []byte, batch *types.OutgoingTxBatch) bool {
		batches = append(batches, batch)
		return len(batches) == MaxResults
	})
	if len(batches) == 0 {
		return nil, nil
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, batches)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryBatchFees(ctx sdk.Context, keeper Keeper, req abci.RequestQuery, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryBatchFeeRequest
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if params.MinBatchFees == nil {
		params.MinBatchFees = make([]types.MinBatchFee, 0)
	}
	val := types.QueryBatchFeeResponse{BatchFees: keeper.GetAllBatchFees(ctx, MaxResults, params.MinBatchFees)}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, val)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

// queryBatch gets a batch by tokenContract and nonce
func queryBatch(ctx sdk.Context, nonceStr string, tokenContract string, keeper Keeper) ([]byte, error) {
	nonce, err := uint64FromString(nonceStr)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "nonce")
	}
	if types.ValidateEthereumAddress(tokenContract) != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "token contract")
	}
	foundBatch := keeper.GetOutgoingTxBatch(ctx, tokenContract, nonce)
	if foundBatch == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Can not find tx batch")
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, foundBatch)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
	return res, nil
}

func queryGravityID(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	gravityID := keeper.GetGravityID(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, gravityID)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}

func queryDenomToToken(ctx sdk.Context, denom string, keeper Keeper) ([]byte, error) {
	if len(denom) <= 0 {
		return nil, sdkerrors.Wrap(types.ErrUnknown, "denom")
	}
	bridgeToken := keeper.GetDenomByBridgeToken(ctx, denom)
	if bridgeToken == nil {
		return nil, sdkerrors.Wrap(types.ErrEmpty, "bridge token is not exist")
	}
	bytes, err := codec.MarshalJSONIndent(types.ModuleCdc, types.QueryDenomToTokenResponse{
		Token:      bridgeToken.Token,
		ChannelIbc: bridgeToken.ChannelIbc,
	})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bytes, nil
}

func queryTokenToDenom(ctx sdk.Context, token string, keeper Keeper) ([]byte, error) {
	if types.ValidateEthereumAddress(token) != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "token")
	}
	bridgeToken := keeper.GetBridgeTokenDenom(ctx, token)
	if bridgeToken == nil {
		return nil, sdkerrors.Wrap(types.ErrEmpty, "bridge token is not exist")
	}
	bytes, err := codec.MarshalJSONIndent(types.ModuleCdc, types.QueryTokenToDenomResponse{
		Denom:      bridgeToken.Denom,
		ChannelIbc: bridgeToken.ChannelIbc,
	})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bytes, nil
}

func queryPendingSendToExternal(ctx sdk.Context, senderAddr string, k Keeper) ([]byte, error) {
	if _, err := sdk.AccAddressFromBech32(senderAddr); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "sender address")
	}
	batches := k.GetOutgoingTxBatches(ctx)
	unbatchedTx := k.GetUnbatchedTransactions(ctx)
	senderAddress := senderAddr
	res := types.QueryPendingSendToExternalResponse{}
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
	bytes, err := codec.MarshalJSONIndent(types.ModuleCdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bytes, nil
}

// uint64FromString to parse out a uint64 for a nonce
func uint64FromString(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
