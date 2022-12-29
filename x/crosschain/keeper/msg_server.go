package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

var _ types.MsgServer = MsgServer{}

type MsgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &MsgServer{Keeper: keeper}
}

func (s MsgServer) BondedOracle(c context.Context, msg *types.MsgBondedOracle) (*types.MsgBondedOracleResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle address")
	}
	bridgerAddr, err := sdk.AccAddressFromBech32(msg.BridgerAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "bridger address")
	}
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "validator address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	if !s.IsProposalOracle(ctx, msg.OracleAddress) {
		return nil, types.ErrNoFoundOracle
	}
	// check oracle has set bridger address
	if _, found := s.GetOracle(ctx, oracleAddr); found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle existed bridger address")
	}
	// check bridger address is bound to oracle
	if _, found := s.GetOracleAddressByBridgerKey(ctx, bridgerAddr); found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "bridger address is bound to oracle")
	}
	// check external address is bound to oracle
	if _, found := s.GetOracleByExternalAddress(ctx, msg.ExternalAddress); found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "external address is bound to oracle")
	}
	threshold := s.GetOracleDelegateThreshold(ctx)
	oracle := types.Oracle{
		OracleAddress:     oracleAddr.String(),
		BridgerAddress:    bridgerAddr.String(),
		ExternalAddress:   msg.ExternalAddress,
		DelegateAmount:    msg.DelegateAmount.Amount,
		StartHeight:       ctx.BlockHeight(),
		Online:            true,
		DelegateValidator: msg.ValidatorAddress,
		SlashTimes:        0,
	}
	if threshold.Denom != msg.DelegateAmount.Denom {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "delegate denom got %s, expected %s", msg.DelegateAmount.Denom, threshold.Denom)
	}
	if msg.DelegateAmount.IsLT(threshold) {
		return nil, types.ErrDelegateAmountBelowMinimum
	}
	if msg.DelegateAmount.Amount.GT(threshold.Amount.Mul(sdk.NewInt(s.GetOracleDelegateMultiple(ctx)))) {
		return nil, types.ErrDelegateAmountAboveMaximum
	}
	validator, found := s.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	delegateAddr := oracle.GetDelegateAddress(s.moduleName)
	if err := s.bankKeeper.SendCoins(ctx, oracleAddr, delegateAddr, sdk.NewCoins(msg.DelegateAmount)); err != nil {
		return nil, err
	}
	newShares, err := s.stakingKeeper.Delegate(ctx, delegateAddr, msg.DelegateAmount.Amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		stakingtypes.EventTypeDelegate,
		sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValidatorAddress),
		sdk.NewAttribute(sdk.AttributeKeyAmount, msg.DelegateAmount.String()),
		sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
	))

	s.SetOracle(ctx, oracle)
	s.SetOracleByBridger(ctx, bridgerAddr, oracleAddr)
	s.SetOracleByExternalAddress(ctx, msg.ExternalAddress, oracleAddr)
	s.CommonSetOracleTotalPower(ctx)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OracleAddress),
		),
	})

	return &types.MsgBondedOracleResponse{}, nil
}

func (s MsgServer) AddDelegate(c context.Context, msg *types.MsgAddDelegate) (*types.MsgAddDelegateResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	if !s.IsProposalOracle(ctx, msg.OracleAddress) {
		return nil, types.ErrNoFoundOracle
	}
	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, types.ErrNoFoundOracle
	}
	validator, found := s.stakingKeeper.GetValidator(ctx, oracle.GetValidator())
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	threshold := s.GetOracleDelegateThreshold(ctx)

	if threshold.Denom != msg.Amount.Denom {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "delegate denom got %s, expected %s", msg.Amount.Denom, threshold.Denom)
	}

	slashAmount := sdk.NewCoin(fxtypes.DefaultDenom, oracle.GetSlashAmount(s.GetSlashFraction(ctx)))
	if slashAmount.IsPositive() && msg.Amount.Amount.LTE(slashAmount.Amount) {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "not sufficient slash amount")
	}

	delegateCoin := sdk.NewCoin(fxtypes.DefaultDenom, msg.Amount.Amount.Sub(slashAmount.Amount))

	oracle.DelegateAmount = oracle.DelegateAmount.Add(delegateCoin.Amount)
	if oracle.DelegateAmount.Sub(threshold.Amount).IsNegative() {
		return nil, types.ErrDelegateAmountBelowMinimum
	}
	if oracle.DelegateAmount.GT(threshold.Amount.Mul(sdk.NewInt(s.GetOracleDelegateMultiple(ctx)))) {
		return nil, types.ErrDelegateAmountAboveMaximum
	}

	if slashAmount.IsPositive() {
		if err := s.bankKeeper.SendCoinsFromAccountToModule(ctx, oracleAddr, s.moduleName, sdk.NewCoins(slashAmount)); err != nil {
			return nil, err
		}
		if err := s.bankKeeper.BurnCoins(ctx, s.moduleName, sdk.NewCoins(slashAmount)); err != nil {
			return nil, err
		}
	}

	delegateAddr := oracle.GetDelegateAddress(s.moduleName)
	if err := s.bankKeeper.SendCoins(ctx, oracleAddr, delegateAddr, sdk.NewCoins(delegateCoin)); err != nil {
		return nil, err
	}
	newShares, err := s.stakingKeeper.Delegate(ctx, delegateAddr, delegateCoin.Amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}

	if !oracle.Online {
		oracle.Online = true
		oracle.StartHeight = ctx.BlockHeight()
	}

	s.SetOracle(ctx, oracle)
	s.CommonSetOracleTotalPower(ctx)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeDelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, oracle.DelegateValidator),
			sdk.NewAttribute(sdk.AttributeKeyAmount, delegateCoin.Amount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OracleAddress),
		),
	})

	return &types.MsgAddDelegateResponse{}, nil
}

func (s MsgServer) ReDelegate(c context.Context, msg *types.MsgReDelegate) (*types.MsgReDelegateResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, types.ErrNoFoundOracle
	}
	if !oracle.Online {
		return nil, types.ErrOracleNotOnLine
	}
	if oracle.DelegateValidator == msg.ValidatorAddress {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "validator address is not changed")
	}

	delegateAddr := oracle.GetDelegateAddress(s.moduleName)
	valSrcAddress := oracle.GetValidator()

	valDestAddress, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "validator address")
	}
	delegation, found := s.stakingKeeper.GetDelegation(ctx, delegateAddr, valSrcAddress)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "no delegation for (address, validator) tuple")
	}
	completionTime, err := s.stakingKeeper.BeginRedelegation(ctx, delegateAddr, valSrcAddress, valDestAddress, delegation.Shares)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		stakingtypes.EventTypeRedelegate,
		sdk.NewAttribute(stakingtypes.AttributeKeySrcValidator, oracle.DelegateValidator),
		sdk.NewAttribute(stakingtypes.AttributeKeyDstValidator, msg.ValidatorAddress),
		sdk.NewAttribute(sdk.AttributeKeyAmount, oracle.DelegateAmount.String()),
		sdk.NewAttribute(stakingtypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
	))
	return &types.MsgReDelegateResponse{}, err
}

func (s MsgServer) EditBridger(c context.Context, msg *types.MsgEditBridger) (*types.MsgEditBridgerResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle address")
	}
	bridgerAddr, err := sdk.AccAddressFromBech32(msg.BridgerAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, types.ErrNoFoundOracle
	}
	if !oracle.Online {
		return nil, types.ErrOracleNotOnLine
	}
	if oracle.BridgerAddress == msg.BridgerAddress {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "bridger address is not changed")
	}
	oracle.BridgerAddress = msg.BridgerAddress
	s.Keeper.SetOracle(ctx, oracle)
	s.Keeper.SetOracleByBridger(ctx, bridgerAddr, oracleAddr)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.OracleAddress),
	))
	return &types.MsgEditBridgerResponse{}, nil
}

func (s MsgServer) WithdrawReward(c context.Context, msg *types.MsgWithdrawReward) (*types.MsgWithdrawRewardResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, types.ErrNoFoundOracle
	}
	if !oracle.Online {
		return nil, types.ErrOracleNotOnLine
	}

	delegateAddr := oracle.GetDelegateAddress(s.moduleName)
	if _, err := s.distributionKeeper.WithdrawDelegationRewards(ctx, delegateAddr, oracle.GetValidator()); err != nil {
		return nil, err
	}
	balances := s.bankKeeper.GetAllBalances(ctx, delegateAddr)
	if !balances.IsAllPositive() {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "rewards")
	}
	if err = s.bankKeeper.SendCoins(ctx, delegateAddr, oracleAddr, balances); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OracleAddress),
		),
	)
	return &types.MsgWithdrawRewardResponse{}, nil
}

func (s MsgServer) UnbondedOracle(c context.Context, msg *types.MsgUnbondedOracle) (*types.MsgUnbondedOracleResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	if s.IsProposalOracle(ctx, msg.OracleAddress) {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "need to pass a proposal to unbind")
	}
	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, types.ErrNoFoundOracle
	}
	if oracle.Online {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle on line")
	}
	delegateAddr := oracle.GetDelegateAddress(s.moduleName)
	validatorAddr := oracle.GetValidator()
	if _, found := s.stakingKeeper.GetUnbondingDelegation(ctx, delegateAddr, validatorAddr); found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "exist unbonding delegation")
	}
	balances := s.bankKeeper.GetAllBalances(ctx, delegateAddr)
	slashAmount := sdk.NewCoin(fxtypes.DefaultDenom, oracle.GetSlashAmount(s.GetSlashFraction(ctx)))
	if slashAmount.IsPositive() {
		if balances.AmountOf(fxtypes.DefaultDenom).LT(slashAmount.Amount) {
			return nil, sdkerrors.Wrap(types.ErrInvalid, "not sufficient slash amount")
		}
		if err := s.bankKeeper.SendCoinsFromAccountToModule(ctx, delegateAddr, s.moduleName, sdk.NewCoins(slashAmount)); err != nil {
			return nil, err
		}
		if err := s.bankKeeper.BurnCoins(ctx, s.moduleName, sdk.NewCoins(slashAmount)); err != nil {
			return nil, err
		}
	}
	sendCoins := balances.Sub(sdk.NewCoins(slashAmount))
	for i := 0; i < len(sendCoins); i++ {
		if !sendCoins[i].IsPositive() {
			sendCoins = append(sendCoins[:i], sendCoins[i+1:]...)
			i--
		}
	}
	if sendCoins.IsAllPositive() {
		if err := s.bankKeeper.SendCoins(ctx, delegateAddr, oracleAddr, sendCoins); err != nil {
			return nil, err
		}
	}

	s.DelOracleByExternalAddress(ctx, oracle.ExternalAddress)
	s.DelOracleByBridger(ctx, oracle.GetBridger())
	s.DelOracle(ctx, oracle.GetOracle())
	s.DelLastEventNonceByOracle(ctx, oracleAddr)

	return &types.MsgUnbondedOracleResponse{}, nil
}

func (s MsgServer) SendToExternal(c context.Context, msg *types.MsgSendToExternal) (*types.MsgSendToExternalResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "sender address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	txID, err := s.AddToOutgoingPool(ctx, sender, msg.Dest, msg.Amount, msg.BridgeFee)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

	return &types.MsgSendToExternalResponse{
		OutgoingTxId: txID,
	}, nil
}

func (s MsgServer) CancelSendToExternal(c context.Context, msg *types.MsgCancelSendToExternal) (*types.MsgCancelSendToExternalResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "sender address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err = s.RemoveFromOutgoingPoolAndRefund(ctx, msg.TransactionId, sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

	return &types.MsgCancelSendToExternalResponse{}, nil
}

func (s MsgServer) RequestBatch(c context.Context, msg *types.MsgRequestBatch) (*types.MsgRequestBatchResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "sender address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	bridgeToken := s.GetDenomByBridgeToken(ctx, msg.Denom)
	if bridgeToken == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "bridge token is not exist")
	}

	_, found := s.GetOracleAddressByBridgerKey(ctx, sender)
	if !found {
		if !s.IsProposalOracle(ctx, msg.Sender) {
			return nil, sdkerrors.Wrap(types.ErrEmpty, "sender must be oracle or bridger")
		}
	}

	batch, err := s.BuildOutgoingTxBatch(ctx, bridgeToken.Token, msg.FeeReceive, OutgoingTxBatchSize, msg.MinimumFee, msg.BaseFee)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

	return &types.MsgRequestBatchResponse{
		BatchNonce: batch.BatchNonce,
	}, nil
}

func (s MsgServer) ConfirmBatch(c context.Context, msg *types.MsgConfirmBatch) (*types.MsgConfirmBatchResponse, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(msg.BridgerAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "bridger address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// fetch the outgoing batch given the nonce
	batch := s.GetOutgoingTxBatch(ctx, msg.TokenContract, msg.Nonce)
	if batch == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find batch")
	}

	checkpoint, err := batch.GetCheckpoint(s.GetGravityID(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, err.Error())
	}

	oracleAddr, err := s.confirmHandlerCommon(ctx, bridgerAddr, msg.ExternalAddress, msg.Signature, checkpoint)
	if err != nil {
		return nil, err
	}
	// check if we already have this confirm
	if s.GetBatchConfirm(ctx, msg.Nonce, msg.TokenContract, oracleAddr) != nil {
		return nil, sdkerrors.Wrap(types.ErrDuplicate, "signature")
	}
	s.SetBatchConfirm(ctx, oracleAddr, msg)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.BridgerAddress),
	))

	return nil, nil
}

func (s MsgServer) OracleSetConfirm(c context.Context, msg *types.MsgOracleSetConfirm) (*types.MsgOracleSetConfirmResponse, error) {
	bridgerAddr, err := sdk.AccAddressFromBech32(msg.BridgerAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "bridger address")
	}

	ctx := sdk.UnwrapSDKContext(c)
	oracleSet := s.GetOracleSet(ctx, msg.Nonce)
	if oracleSet == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find oracleSet")
	}

	checkpoint, err := oracleSet.GetCheckpoint(s.GetGravityID(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, err.Error())
	}
	oracleAddr, err := s.confirmHandlerCommon(ctx, bridgerAddr, msg.ExternalAddress, msg.Signature, checkpoint)
	if err != nil {
		return nil, err
	}
	// check if we already have this confirm
	if s.GetOracleSetConfirm(ctx, msg.Nonce, oracleAddr) != nil {
		return nil, sdkerrors.Wrap(types.ErrDuplicate, "signature")
	}
	s.SetOracleSetConfirm(ctx, oracleAddr, msg)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.BridgerAddress),
	))

	return &types.MsgOracleSetConfirmResponse{}, nil
}

func (s MsgServer) SendToExternalClaim(c context.Context, msg *types.MsgSendToExternalClaim) (*types.MsgSendToExternalClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.claimHandlerCommon(ctx, msg); err != nil {
		return nil, err
	}
	return &types.MsgSendToExternalClaimResponse{}, nil
}

func (s MsgServer) SendToFxClaim(c context.Context, msg *types.MsgSendToFxClaim) (*types.MsgSendToFxClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.claimHandlerCommon(ctx, msg); err != nil {
		return nil, err
	}
	return &types.MsgSendToFxClaimResponse{}, nil
}

func (s MsgServer) BridgeTokenClaim(c context.Context, msg *types.MsgBridgeTokenClaim) (*types.MsgBridgeTokenClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.claimHandlerCommon(ctx, msg); err != nil {
		return nil, err
	}
	return &types.MsgBridgeTokenClaimResponse{}, nil
}

// OracleSetUpdateClaim handles claims for executing a oracle set update on Ethereum
func (s MsgServer) OracleSetUpdateClaim(c context.Context, msg *types.MsgOracleSetUpdatedClaim) (*types.MsgOracleSetUpdatedClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	bridgerAddr := msg.GetClaimer()
	oracleAddr, err := s.checkBridgerIsOracle(ctx, bridgerAddr)
	if err != nil {
		return nil, err
	}

	for _, member := range msg.Members {
		if _, found := s.GetOracleByExternalAddress(ctx, member.ExternalAddress); !found {
			return nil, sdkerrors.Wrap(types.ErrInvalid, "external address")
		}
	}
	// Add the claim to the store
	if _, err := s.Attest(ctx, oracleAddr, msg); err != nil {
		return nil, sdkerrors.Wrap(err, "create attestation")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, s.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, bridgerAddr.String()),
	))

	return &types.MsgOracleSetUpdatedClaimResponse{}, nil
}

func (s MsgServer) checkBridgerIsOracle(ctx sdk.Context, bridgerAddr sdk.AccAddress) (oracleAddr sdk.AccAddress, err error) {
	oracleAddr, found := s.GetOracleAddressByBridgerKey(ctx, bridgerAddr)
	if !found {
		return oracleAddr, types.ErrNoFoundOracle
	}
	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return oracleAddr, types.ErrNoFoundOracle
	}
	if !oracle.Online {
		return oracleAddr, types.ErrOracleNotOnLine
	}
	return oracleAddr, nil
}

// claimHandlerCommon is an internal function that provides common code for processing claims once they are
// translated from the message to the Ethereum claim interface
func (s MsgServer) claimHandlerCommon(ctx sdk.Context, msg types.ExternalClaim) (err error) {
	bridgerAddr := msg.GetClaimer()
	oracleAddr, err := s.checkBridgerIsOracle(ctx, bridgerAddr)
	if err != nil {
		return err
	}

	// Add the claim to the store
	if _, err := s.Attest(ctx, oracleAddr, msg); err != nil {
		return sdkerrors.Wrap(err, "create attestation")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, s.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, bridgerAddr.String()),
	))

	return nil
}

func (s MsgServer) confirmHandlerCommon(ctx sdk.Context, bridgerAddr sdk.AccAddress, signatureAddr, signature string, checkpoint []byte) (oracleAddr sdk.AccAddress, err error) {
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}

	oracleAddr, found := s.GetOracleByExternalAddress(ctx, signatureAddr)
	if !found {
		return nil, types.ErrNoFoundOracle
	}

	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, types.ErrNoFoundOracle
	}

	if oracle.ExternalAddress != signatureAddr {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "got %s, expected %s", signatureAddr, oracle.ExternalAddress)
	}
	if oracle.BridgerAddress != bridgerAddr.String() {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "got %s, expected %s", bridgerAddr, oracle.BridgerAddress)
	}
	if err = types.ValidateEthereumSignature(checkpoint, sigBytes, oracle.ExternalAddress); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s found %s", oracle.ExternalAddress, hex.EncodeToString(checkpoint), signature))
	}
	return oracleAddr, nil
}
