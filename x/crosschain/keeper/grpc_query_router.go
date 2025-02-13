package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

var _ types.QueryServer = RouterKeeper{}

func (k RouterKeeper) getQueryServerByChainName(chainName string) (types.QueryServer, error) {
	if !k.router.HasRoute(chainName) {
		return nil, status.Error(codes.InvalidArgument, "chain name not found:"+chainName)
	}
	return k.router.GetRoute(chainName).QueryServer, nil
}

func (k RouterKeeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.Params(c, req)
	}
}

func (k RouterKeeper) CurrentOracleSet(c context.Context, req *types.QueryCurrentOracleSetRequest) (*types.QueryCurrentOracleSetResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.CurrentOracleSet(c, req)
	}
}

func (k RouterKeeper) OracleSetRequest(c context.Context, req *types.QueryOracleSetRequestRequest) (*types.QueryOracleSetRequestResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.OracleSetRequest(c, req)
	}
}

func (k RouterKeeper) OracleSetConfirm(c context.Context, req *types.QueryOracleSetConfirmRequest) (*types.QueryOracleSetConfirmResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.OracleSetConfirm(c, req)
	}
}

func (k RouterKeeper) OracleSetConfirmsByNonce(c context.Context, req *types.QueryOracleSetConfirmsByNonceRequest) (*types.QueryOracleSetConfirmsByNonceResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.OracleSetConfirmsByNonce(c, req)
	}
}

func (k RouterKeeper) LastOracleSetRequests(c context.Context, req *types.QueryLastOracleSetRequestsRequest) (*types.QueryLastOracleSetRequestsResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.LastOracleSetRequests(c, req)
	}
}

func (k RouterKeeper) LastPendingOracleSetRequestByAddr(c context.Context, req *types.QueryLastPendingOracleSetRequestByAddrRequest) (*types.QueryLastPendingOracleSetRequestByAddrResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.LastPendingOracleSetRequestByAddr(c, req)
	}
}

func (k RouterKeeper) LastPendingBatchRequestByAddr(c context.Context, req *types.QueryLastPendingBatchRequestByAddrRequest) (*types.QueryLastPendingBatchRequestByAddrResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.LastPendingBatchRequestByAddr(c, req)
	}
}

func (k RouterKeeper) OutgoingTxBatches(c context.Context, req *types.QueryOutgoingTxBatchesRequest) (*types.QueryOutgoingTxBatchesResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.OutgoingTxBatches(c, req)
	}
}

func (k RouterKeeper) OutgoingTxBatch(c context.Context, req *types.QueryOutgoingTxBatchRequest) (*types.QueryOutgoingTxBatchResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.OutgoingTxBatch(c, req)
	}
}

func (k RouterKeeper) BatchConfirm(c context.Context, req *types.QueryBatchConfirmRequest) (*types.QueryBatchConfirmResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BatchConfirm(c, req)
	}
}

func (k RouterKeeper) BatchConfirms(c context.Context, req *types.QueryBatchConfirmsRequest) (*types.QueryBatchConfirmsResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BatchConfirms(c, req)
	}
}

func (k RouterKeeper) LastEventNonceByAddr(c context.Context, req *types.QueryLastEventNonceByAddrRequest) (*types.QueryLastEventNonceByAddrResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.LastEventNonceByAddr(c, req)
	}
}

func (k RouterKeeper) DenomToToken(c context.Context, req *types.QueryDenomToTokenRequest) (*types.QueryDenomToTokenResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.DenomToToken(c, req)
	}
}

func (k RouterKeeper) TokenToDenom(c context.Context, req *types.QueryTokenToDenomRequest) (*types.QueryTokenToDenomResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.TokenToDenom(c, req)
	}
}

func (k RouterKeeper) GetOracleByAddr(c context.Context, req *types.QueryOracleByAddrRequest) (*types.QueryOracleResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.GetOracleByAddr(c, req)
	}
}

func (k RouterKeeper) GetOracleByBridgerAddr(c context.Context, req *types.QueryOracleByBridgerAddrRequest) (*types.QueryOracleResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.GetOracleByBridgerAddr(c, req)
	}
}

func (k RouterKeeper) GetOracleByExternalAddr(c context.Context, req *types.QueryOracleByExternalAddrRequest) (*types.QueryOracleResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.GetOracleByExternalAddr(c, req)
	}
}

func (k RouterKeeper) LastObservedBlockHeight(c context.Context, req *types.QueryLastObservedBlockHeightRequest) (*types.QueryLastObservedBlockHeightResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.LastObservedBlockHeight(c, req)
	}
}

func (k RouterKeeper) LastEventBlockHeightByAddr(c context.Context, req *types.QueryLastEventBlockHeightByAddrRequest) (*types.QueryLastEventBlockHeightByAddrResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.LastEventBlockHeightByAddr(c, req)
	}
}

func (k RouterKeeper) Oracles(c context.Context, req *types.QueryOraclesRequest) (*types.QueryOraclesResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.Oracles(c, req)
	}
}

func (k RouterKeeper) ProjectedBatchTimeoutHeight(c context.Context, req *types.QueryProjectedBatchTimeoutHeightRequest) (*types.QueryProjectedBatchTimeoutHeightResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.ProjectedBatchTimeoutHeight(c, req)
	}
}

func (k RouterKeeper) BridgeTokens(c context.Context, req *types.QueryBridgeTokensRequest) (*types.QueryBridgeTokensResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeTokens(c, req)
	}
}

func (k RouterKeeper) BridgeTokensByChain(c context.Context, req *types.QueryBridgeTokensByChainRequest) (*types.QueryBridgeTokensByChainResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeTokensByChain(c, req)
	}
}

func (k RouterKeeper) BridgeTokensByDenom(c context.Context, req *types.QueryBridgeTokensByDenomRequest) (*types.QueryBridgeTokensByDenomResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeTokensByDenom(c, req)
	}
}

func (k RouterKeeper) BridgeTokensByERC20(c context.Context, req *types.QueryBridgeTokensByERC20Request) (*types.QueryBridgeTokensByERC20Response, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeTokensByERC20(c, req)
	}
}

func (k RouterKeeper) BridgeCoinByDenom(c context.Context, req *types.QueryBridgeCoinByDenomRequest) (*types.QueryBridgeCoinByDenomResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeCoinByDenom(c, req)
	}
}

func (k RouterKeeper) BridgeChainList(c context.Context, req *types.QueryBridgeChainListRequest) (*types.QueryBridgeChainListResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(ethtypes.ModuleName); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeChainList(c, req)
	}
}

func (k RouterKeeper) BridgeCalls(c context.Context, req *types.QueryBridgeCallsRequest) (*types.QueryBridgeCallsResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeCalls(c, req)
	}
}

func (k RouterKeeper) BridgeCallConfirmByNonce(c context.Context, req *types.QueryBridgeCallConfirmByNonceRequest) (*types.QueryBridgeCallConfirmByNonceResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeCallConfirmByNonce(c, req)
	}
}

func (k RouterKeeper) BridgeCallByNonce(c context.Context, req *types.QueryBridgeCallByNonceRequest) (*types.QueryBridgeCallByNonceResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeCallByNonce(c, req)
	}
}

func (k RouterKeeper) BridgeCallBySender(c context.Context, req *types.QueryBridgeCallBySenderRequest) (*types.QueryBridgeCallBySenderResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeCallBySender(c, req)
	}
}

func (k RouterKeeper) LastPendingBridgeCallByAddr(c context.Context, req *types.QueryLastPendingBridgeCallByAddrRequest) (*types.QueryLastPendingBridgeCallByAddrResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.LastPendingBridgeCallByAddr(c, req)
	}
}

func (k RouterKeeper) PendingExecuteClaim(ctx context.Context, req *types.QueryPendingExecuteClaimRequest) (*types.QueryPendingExecuteClaimResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.PendingExecuteClaim(ctx, req)
	}
}

func (k RouterKeeper) BridgeCallQuoteByNonce(ctx context.Context, req *types.QueryBridgeCallQuoteByNonceRequest) (*types.QueryBridgeCallQuoteByNonceResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(req.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeCallQuoteByNonce(ctx, req)
	}
}

func (k RouterKeeper) BridgeCallsByFeeReceiver(ctx context.Context, request *types.QueryBridgeCallsByFeeReceiverRequest) (*types.QueryBridgeCallsByFeeReceiverResponse, error) {
	if queryServer, err := k.getQueryServerByChainName(request.ChainName); err != nil {
		return nil, err
	} else {
		return queryServer.BridgeCallsByFeeReceiver(ctx, request)
	}
}
