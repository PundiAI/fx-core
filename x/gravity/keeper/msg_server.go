// nolint:staticcheck
package keeper

import (
	"context"

	crosschainkeeper "github.com/functionx/fx-core/v7/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v7/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
	"github.com/functionx/fx-core/v7/x/gravity/types"
)

type msgServer struct {
	crosschainkeeper.MsgServer
}

func NewMsgServerImpl(keeper crosschainkeeper.Keeper) types.MsgServer {
	return &msgServer{MsgServer: crosschainkeeper.MsgServer{Keeper: keeper}}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) ValsetConfirm(c context.Context, msg *types.MsgValsetConfirm) (*types.MsgValsetConfirmResponse, error) {
	_, err := k.MsgServer.OracleSetConfirm(c, &crosschaintypes.MsgOracleSetConfirm{
		Nonce:           msg.Nonce,
		BridgerAddress:  msg.Orchestrator,
		ExternalAddress: msg.EthAddress,
		Signature:       msg.Signature,
		ChainName:       ethtypes.ModuleName,
	})
	return &types.MsgValsetConfirmResponse{}, err
}

func (k msgServer) SendToEth(c context.Context, msg *types.MsgSendToEth) (*types.MsgSendToEthResponse, error) {
	_, err := k.MsgServer.SendToExternal(c, &crosschaintypes.MsgSendToExternal{
		Sender:    msg.Sender,
		Dest:      msg.EthDest,
		Amount:    msg.Amount,
		BridgeFee: msg.BridgeFee,
		ChainName: ethtypes.ModuleName,
	})
	return &types.MsgSendToEthResponse{}, err
}

func (k msgServer) RequestBatch(c context.Context, msg *types.MsgRequestBatch) (*types.MsgRequestBatchResponse, error) {
	_, err := k.MsgServer.RequestBatch(c, &crosschaintypes.MsgRequestBatch{
		Sender:     msg.Sender,
		Denom:      msg.Denom,
		MinimumFee: msg.MinimumFee,
		FeeReceive: msg.FeeReceive,
		ChainName:  ethtypes.ModuleName,
		BaseFee:    msg.BaseFee,
	})
	return &types.MsgRequestBatchResponse{}, err
}

func (k msgServer) ConfirmBatch(c context.Context, msg *types.MsgConfirmBatch) (*types.MsgConfirmBatchResponse, error) {
	_, err := k.MsgServer.ConfirmBatch(c, &crosschaintypes.MsgConfirmBatch{
		Nonce:           msg.Nonce,
		TokenContract:   msg.TokenContract,
		BridgerAddress:  msg.Orchestrator,
		ExternalAddress: msg.EthSigner,
		Signature:       msg.Signature,
		ChainName:       ethtypes.ModuleName,
	})
	return &types.MsgConfirmBatchResponse{}, err
}

func (k msgServer) DepositClaim(c context.Context, msg *types.MsgDepositClaim) (*types.MsgDepositClaimResponse, error) {
	_, err := k.MsgServer.SendToFxClaim(c, &crosschaintypes.MsgSendToFxClaim{
		EventNonce:     msg.EventNonce,
		BlockHeight:    msg.BlockHeight,
		TokenContract:  msg.TokenContract,
		Amount:         msg.Amount,
		Sender:         msg.EthSender,
		Receiver:       msg.FxReceiver,
		TargetIbc:      msg.TargetIbc,
		BridgerAddress: msg.Orchestrator,
		ChainName:      ethtypes.ModuleName,
	})
	return &types.MsgDepositClaimResponse{}, err
}

func (k msgServer) WithdrawClaim(c context.Context, msg *types.MsgWithdrawClaim) (*types.MsgWithdrawClaimResponse, error) {
	_, err := k.MsgServer.SendToExternalClaim(c, &crosschaintypes.MsgSendToExternalClaim{
		EventNonce:     msg.EventNonce,
		BlockHeight:    msg.BlockHeight,
		BatchNonce:     msg.BatchNonce,
		TokenContract:  msg.TokenContract,
		BridgerAddress: msg.Orchestrator,
		ChainName:      ethtypes.ModuleName,
	})
	return &types.MsgWithdrawClaimResponse{}, err
}

func (k msgServer) CancelSendToEth(c context.Context, msg *types.MsgCancelSendToEth) (*types.MsgCancelSendToEthResponse, error) {
	_, err := k.MsgServer.CancelSendToExternal(c, &crosschaintypes.MsgCancelSendToExternal{
		TransactionId: msg.TransactionId,
		Sender:        msg.Sender,
		ChainName:     ethtypes.ModuleName,
	})
	return &types.MsgCancelSendToEthResponse{}, err
}

func (k msgServer) ValsetUpdateClaim(c context.Context, msg *types.MsgValsetUpdatedClaim) (*types.MsgValsetUpdatedClaimResponse, error) {
	msg2 := &crosschaintypes.MsgOracleSetUpdatedClaim{
		EventNonce:     msg.EventNonce,
		BlockHeight:    msg.BlockHeight,
		OracleSetNonce: msg.ValsetNonce,
		Members:        make([]crosschaintypes.BridgeValidator, len(msg.Members)),
		BridgerAddress: msg.Orchestrator,
		ChainName:      ethtypes.ModuleName,
	}
	for i := 0; i < len(msg.Members); i++ {
		msg2.Members[i] = crosschaintypes.BridgeValidator{
			Power:           msg.Members[i].Power,
			ExternalAddress: msg.Members[i].EthAddress,
		}
	}
	_, err := k.MsgServer.OracleSetUpdateClaim(c, msg2)
	return &types.MsgValsetUpdatedClaimResponse{}, err
}
