package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/crosschain/types"
)

const depositMultiple = 10

var _ types.MsgServer = EthereumMsgServer{}

type EthereumMsgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) ProposalMsgServer {
	return &EthereumMsgServer{Keeper: keeper}
}

func (s EthereumMsgServer) SetOrchestratorAddress(c context.Context, msg *types.MsgSetOrchestratorAddress) (*types.MsgSetOrchestratorAddressResponse, error) {
	var err error
	var oracleAddress, orchestratorAddr sdk.AccAddress
	if oracleAddress, err = sdk.AccAddressFromBech32(msg.Oracle); err != nil {
		return nil, sdkerrors.Wrap(types.ErrOracleAddress, msg.Oracle)
	}
	if orchestratorAddr, err = sdk.AccAddressFromBech32(msg.Orchestrator); err != nil {
		return nil, sdkerrors.Wrap(types.ErrOrchestratorAddress, msg.Orchestrator)
	}
	ctx := sdk.UnwrapSDKContext(c)
	if !s.IsOracle(ctx, msg.Oracle) {
		return nil, types.ErrNotOracle
	}
	// check oracle has set orchestrator address
	if _, found := s.GetOracle(ctx, oracleAddress); found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle existed orchestrator address")
	}
	// check orchestrator address is bound to oracle
	if _, found := s.GetOracleAddressByOrchestratorKey(ctx, orchestratorAddr); found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "orchestrator address is bound to oracle")
	}
	// check external address is bound to oracle
	if _, found := s.GetOracleByExternalAddress(ctx, msg.ExternalAddress); found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "external address is bound to oracle")
	}

	depositThreshold := s.GetOracleDepositThreshold(ctx)
	if depositThreshold.Denom != msg.Deposit.Denom {
		return nil, sdkerrors.Wrapf(types.ErrBadDepositDenom, "got %s, expected %s", msg.Deposit.Denom, depositThreshold.Denom)
	}
	if msg.Deposit.IsLT(depositThreshold) {
		return nil, types.ErrDepositAmountBelowMinimum
	}
	if msg.Deposit.Amount.GT(depositThreshold.Amount.Mul(sdk.NewInt(depositMultiple))) {
		return nil, types.ErrDepositAmountBelowMaximum
	}

	if err = s.bankKeeper.SendCoinsFromAccountToModule(ctx, oracleAddress, s.moduleName, sdk.NewCoins(msg.Deposit)); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	oracle := types.Oracle{
		OracleAddress:       oracleAddress.String(),
		OrchestratorAddress: orchestratorAddr.String(),
		ExternalAddress:     msg.ExternalAddress,
		DepositAmount:       msg.Deposit,
		StartHeight:         ctx.BlockHeight(),
		Jailed:              false,
		JailedHeight:        0,
	}
	// save oracle
	s.SetOracle(ctx, oracle)
	// set the orchestrator address
	s.SetOracleByOrchestrator(ctx, oracleAddress, orchestratorAddr)
	// set the ethereum address
	s.SetExternalAddressForOracle(ctx, oracleAddress, msg.ExternalAddress)
	// save total deposit amount
	totalDeposit := s.GetTotalDeposit(ctx)
	s.SetTotalDeposit(ctx, totalDeposit.Add(msg.Deposit))

	s.CommonSetOracleTotalPower(ctx)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, s.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Oracle),
	))

	return &types.MsgSetOrchestratorAddressResponse{}, nil
}

func (s EthereumMsgServer) AddOracleDeposit(c context.Context, msg *types.MsgAddOracleDeposit) (*types.MsgAddOracleDepositResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.Oracle)
	if err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(c)
	// ensure that the oracle exists
	if !s.IsOracle(ctx, oracleAddr.String()) {
		return nil, types.ErrNotOracle
	}
	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, types.ErrNoOracleFound
	}
	depositThreshold := s.GetOracleDepositThreshold(ctx)
	// check deposit denom
	if depositThreshold.Denom != msg.Amount.Denom {
		return nil, sdkerrors.Wrapf(types.ErrBadDepositDenom, "got %s, expected %s", msg.Amount.Denom, depositThreshold.Denom)
	}
	// check oracle total deposit grate then minimum deposit amount
	deposit := oracle.DepositAmount.Add(msg.Amount)
	if deposit.Amount.Sub(depositThreshold.Amount).IsNegative() {
		return nil, types.ErrDepositAmountBelowMinimum
	}
	if deposit.Amount.GT(depositThreshold.Amount.Mul(sdk.NewInt(depositMultiple))) {
		return nil, types.ErrDepositAmountBelowMaximum
	}

	totalDeposit := s.GetTotalDeposit(ctx)
	totalDeposit = totalDeposit.Add(msg.Amount)
	if err = s.bankKeeper.SendCoinsFromAccountToModule(ctx, oracleAddr, s.moduleName, sdk.NewCoins(msg.Amount)); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, err.Error())
	}
	// save new total deposit
	s.SetTotalDeposit(ctx, totalDeposit)
	if oracle.Jailed {
		oracle.Jailed = false
		oracle.StartHeight = ctx.BlockHeight()
	}
	// save oracle new deposit
	oracle.DepositAmount = deposit
	s.SetOracle(ctx, oracle)

	s.CommonSetOracleTotalPower(ctx)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, s.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Oracle),
	))

	return &types.MsgAddOracleDepositResponse{}, nil
}

// SendToExternal handles MsgSendToExternal
func (s EthereumMsgServer) SendToExternal(c context.Context, msg *types.MsgSendToExternal) (*types.MsgSendToExternalResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	txID, err := s.AddToOutgoingPool(ctx, sender, msg.Dest, msg.Amount, msg.BridgeFee)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, err.Error())
	}

	_ = txID
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, s.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

	return &types.MsgSendToExternalResponse{}, nil
}

func (s EthereumMsgServer) CancelSendToExternal(c context.Context, msg *types.MsgCancelSendToExternal) (*types.MsgCancelSendToExternalResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	err = s.RemoveFromOutgoingPoolAndRefund(ctx, msg.TransactionId, sender)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, s.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

	return &types.MsgCancelSendToExternalResponse{}, nil
}

// RequestBatch handles MsgRequestBatch
func (s EthereumMsgServer) RequestBatch(c context.Context, msg *types.MsgRequestBatch) (*types.MsgRequestBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	bridgeToken := s.GetDenomByBridgeToken(ctx, msg.Denom)
	if bridgeToken == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "bridge token is not exist")
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	_, found := s.GetOracleAddressByOrchestratorKey(ctx, sender)
	if !found {
		if !s.IsOracle(ctx, msg.Sender) {
			return nil, sdkerrors.Wrap(types.ErrEmpty, "oracle or orchestrator")
		}
	}

	batch, err := s.BuildOutgoingTXBatch(ctx, bridgeToken.Token, msg.FeeReceive, OutgoingTxBatchSize, msg.MinimumFee, *msg.BaseFee)
	if err != nil {
		return nil, err
	}

	_ = batch
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, s.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

	return &types.MsgRequestBatchResponse{}, nil
}

// ConfirmBatch handles MsgConfirmBatch
func (s EthereumMsgServer) ConfirmBatch(c context.Context, msg *types.MsgConfirmBatch) (*types.MsgConfirmBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// fetch the outgoing batch given the nonce
	batch := s.GetOutgoingTXBatch(ctx, msg.TokenContract, msg.Nonce)
	if batch == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find batch")
	}
	orchestratorAddr, err := sdk.AccAddressFromBech32(msg.OrchestratorAddress)
	if err != nil {
		return nil, types.ErrOrchestratorAddress
	}
	checkpoint, err := batch.GetCheckpoint(s.GetGravityID(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "checkpoint generation")
	}

	oracleAddr, err := s.confirmHandlerCommon(ctx, orchestratorAddr, msg.ExternalAddress, msg.Signature, checkpoint)
	if err != nil {
		return nil, err
	}
	// check if we already have this confirm
	if s.GetBatchConfirm(ctx, msg.Nonce, msg.TokenContract, oracleAddr) != nil {
		return nil, sdkerrors.Wrap(types.ErrDuplicate, "duplicate signature")
	}
	key := s.SetBatchConfirm(ctx, oracleAddr, msg)

	_ = key
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, s.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.OrchestratorAddress),
	))

	return nil, nil
}

// OracleSetConfirm handles MsgOracleSetConfirm
func (s EthereumMsgServer) OracleSetConfirm(c context.Context, msg *types.MsgOracleSetConfirm) (*types.MsgOracleSetConfirmResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	oracleSet := s.GetOracleSet(ctx, msg.Nonce)
	if oracleSet == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find oracleSet")
	}
	orchestratorAddr, err := sdk.AccAddressFromBech32(msg.OrchestratorAddress)
	if err != nil {
		return nil, types.ErrOrchestratorAddress
	}
	checkpoint := oracleSet.GetCheckpoint(s.GetGravityID(ctx))
	oracleAddr, err := s.confirmHandlerCommon(ctx, orchestratorAddr, msg.ExternalAddress, msg.Signature, checkpoint)
	if err != nil {
		return nil, err
	}
	// check if we already have this confirm
	if s.GetOracleSetConfirm(ctx, msg.Nonce, oracleAddr) != nil {
		return nil, sdkerrors.Wrap(types.ErrDuplicate, "duplicate signature")
	}
	key := s.SetOracleSetConfirm(ctx, oracleAddr, *msg)

	_ = key
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, s.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.OrchestratorAddress),
	))

	return &types.MsgOracleSetConfirmResponse{}, nil
}

// SendToExternalClaim handles MsgSendToExternalClaim
// executed aka 'observed' and had its slashing window expire) that will never be cleaned up in the end block. This
// should not be a security risk as 'old' events can never execute but it does store spam in the chain.
func (s EthereumMsgServer) SendToExternalClaim(c context.Context, msg *types.MsgSendToExternalClaim) (*types.MsgSendToExternalClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := checkOrchestratorIsOracle(ctx, s.Keeper, msg.Orchestrator); err != nil {
		return nil, err
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendToExternalClaimResponse{}, s.claimHandlerCommon(ctx, any, msg, msg.Type())
}

// SendToFxClaim handles MsgSendToFxClaim
// executed aka 'observed' and had it's slashing window expire) that will never be cleaned up in the endblocker. This
// should not be a security risk as 'old' events can never execute but it does store spam in the chain.
func (s EthereumMsgServer) SendToFxClaim(c context.Context, msg *types.MsgSendToFxClaim) (*types.MsgSendToFxClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := checkOrchestratorIsOracle(ctx, s.Keeper, msg.Orchestrator); err != nil {
		return nil, err
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendToFxClaimResponse{}, s.claimHandlerCommon(ctx, any, msg, msg.Type())
}

func (s EthereumMsgServer) BridgeTokenClaim(c context.Context, msg *types.MsgBridgeTokenClaim) (*types.MsgBridgeTokenClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := checkOrchestratorIsOracle(ctx, s.Keeper, msg.Orchestrator); err != nil {
		return nil, err
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgBridgeTokenClaimResponse{}, s.claimHandlerCommon(ctx, any, msg, msg.Type())
}

// OracleSetUpdateClaim handles claims for executing a oracle set update on Ethereum
func (s EthereumMsgServer) OracleSetUpdateClaim(c context.Context, msg *types.MsgOracleSetUpdatedClaim) (*types.MsgOracleSetUpdatedClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := checkOrchestratorIsOracle(ctx, s.Keeper, msg.Orchestrator); err != nil {
		return nil, err
	}

	for _, member := range msg.Members {
		if _, found := s.GetOracleByExternalAddress(ctx, member.ExternalAddress); !found {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "member oracle")
		}
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgOracleSetUpdatedClaimResponse{}, s.claimHandlerCommon(ctx, any, msg, msg.Type())
}

func checkOrchestratorIsOracle(ctx sdk.Context, keep Keeper, orchestrator string) error {
	orcAddr, err := sdk.AccAddressFromBech32(orchestrator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "orchestrator")
	}
	oracleAddr, found := keep.GetOracleAddressByOrchestratorKey(ctx, orcAddr)
	if !found {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "oracle")
	}

	oracle, found := keep.GetOracle(ctx, oracleAddr)
	if !found {
		return types.ErrNoOracleFound
	}
	if oracle.Jailed {
		return sdkerrors.Wrapf(types.ErrOracleJailed, oracle.OracleAddress)
	}
	return nil
}

// claimHandlerCommon is an internal function that provides common code for processing claims once they are
// translated from the message to the Ethereum claim interface
func (s EthereumMsgServer) claimHandlerCommon(ctx sdk.Context, msgAny *codectypes.Any, msg types.ExternalClaim, msgType string) error {
	// Add the claim to the store
	_, err := s.Attest(ctx, msg, msgAny)
	if err != nil {
		return sdkerrors.Wrap(err, "create attestation")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, s.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.GetClaimer().String()),
	))

	return nil
}

func (s EthereumMsgServer) confirmHandlerCommon(ctx sdk.Context, orchestratorAddr sdk.AccAddress, signatureAddr, signature string, checkpoint []byte) (oracleAddr sdk.AccAddress, err error) {
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}

	oracleAddr, found := s.GetOracleByExternalAddress(ctx, signatureAddr)
	if !found {
		return nil, types.ErrNotOracle
	}

	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, types.ErrNoOracleFound
	}

	if oracle.ExternalAddress != signatureAddr {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "got %s, expected %s", signatureAddr, oracle.ExternalAddress)
	}
	if oracle.OrchestratorAddress != orchestratorAddr.String() {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "got %s, expected %s", orchestratorAddr, oracle.OrchestratorAddress)
	}
	if err = types.ValidateEthereumSignature(checkpoint, sigBytes, oracle.ExternalAddress); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s found %s", oracle.ExternalAddress, hex.EncodeToString(checkpoint), signature))
	}
	return oracle.GetOracle(), nil
}
