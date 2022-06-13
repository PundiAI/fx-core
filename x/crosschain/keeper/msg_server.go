package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/crosschain/types"
)

var _ types.MsgServer = EthereumMsgServer{}

type EthereumMsgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &EthereumMsgServer{Keeper: keeper}
}

func (s EthereumMsgServer) CreateOracleBridger(c context.Context, msg *types.MsgCreateOracleBridger) (*types.MsgCreateOracleBridgerResponse, error) {
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
		return nil, sdkerrors.Wrap(types.ErrNoFoundOracle, msg.OracleAddress)
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
		Jailed:            false,
		JailedHeight:      0,
		DelegateValidator: msg.ValidatorAddress,
		IsValidator:       false,
	}
	oravleVal, found := s.stakingKeeper.GetValidator(ctx, oracleAddr.Bytes())
	if found {
		if msg.OracleAddress != msg.ValidatorAddress {
			return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle is a validator but validator address is not itself")
		}
		if msg.DelegateAmount.IsPositive() {
			return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle is a validator, cannot delegate here")
		}
		if msg.DelegateAmount.Denom != "" && msg.DelegateAmount.Denom != threshold.Denom {
			return nil, sdkerrors.Wrapf(types.ErrInvalid, "delegate denom, got %s, expected %s", msg.DelegateAmount.Denom, threshold.Denom)
		}
		oracle.IsValidator = true
		oracle.Jailed = oravleVal.Jailed
		oracle.JailedHeight = ctx.BlockHeight()
		oracle.DelegateAmount = sdk.NewInt(oravleVal.ConsensusPower(sdk.DefaultPowerReduction))
	} else {
		validator, found := s.stakingKeeper.GetValidator(ctx, valAddr)
		if !found {
			return nil, stakingtypes.ErrNoValidatorFound
		}
		if threshold.Denom != msg.DelegateAmount.Denom {
			return nil, sdkerrors.Wrapf(types.ErrInvalid, "delegate denom, got %s, expected %s", msg.DelegateAmount.Denom, threshold.Denom)
		}
		if msg.DelegateAmount.IsLT(threshold) {
			return nil, types.ErrDelegateAmountBelowMinimum
		}
		if msg.DelegateAmount.Amount.GT(threshold.Amount.Mul(sdk.NewInt(s.GetOracleDelegateMultiple(ctx)))) {
			return nil, types.ErrDelegateAmountBelowMaximum
		}

		deleteAddr := types.GetOracleDelegateAddress(msg.ChainName, oracleAddr)
		newShares, err := s.stakingKeeper.Delegate(ctx, deleteAddr, msg.DelegateAmount.Amount, stakingtypes.Unbonded, validator, true)
		if err != nil {
			return nil, err
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			stakingtypes.EventTypeDelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.DelegateAmount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
		))
	}

	// save oracle
	s.SetOracle(ctx, oracle)
	// set the bridger address
	s.SetOracleByBridger(ctx, oracleAddr, bridgerAddr)
	// set the external address
	s.SetExternalAddressForOracle(ctx, oracleAddr, msg.ExternalAddress)
	// update oracle total power
	s.CommonSetOracleTotalPower(ctx)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OracleAddress),
		),
	})

	return &types.MsgCreateOracleBridgerResponse{}, nil
}

func (s EthereumMsgServer) AddOracleDelegate(c context.Context, msg *types.MsgAddOracleDelegate) (*types.MsgAddOracleDelegateResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNoFoundOracle, msg.OracleAddress)
	}
	if oracle.IsValidator {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle is a validator, cannot delegate here")
	}
	valAddr, err := sdk.ValAddressFromBech32(oracle.DelegateValidator)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "validator address")
	}
	validator, found := s.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	threshold := s.GetOracleDelegateThreshold(ctx)
	// check delegate denom
	if threshold.Denom != msg.Amount.Denom {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "delegate denom, got %s, expected %s", msg.Amount.Denom, threshold.Denom)
	}
	// check oracle total delegateAmount grate then minimum delegateAmount amount
	delegateAmount := oracle.DelegateAmount.Add(msg.Amount.Amount)
	if delegateAmount.Sub(threshold.Amount).IsNegative() {
		return nil, types.ErrDelegateAmountBelowMinimum
	}
	if delegateAmount.GT(threshold.Amount.Mul(sdk.NewInt(s.GetOracleDelegateMultiple(ctx)))) {
		return nil, types.ErrDelegateAmountBelowMaximum
	}

	deleteAddr := types.GetOracleDelegateAddress(msg.ChainName, oracleAddr)
	newShares, err := s.stakingKeeper.Delegate(ctx, deleteAddr, msg.Amount.Amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}

	if oracle.Jailed {
		oracle.Jailed = false
		oracle.StartHeight = ctx.BlockHeight()
	}
	// save oracle new delegateAmount
	oracle.DelegateAmount = delegateAmount
	s.SetOracle(ctx, oracle)

	s.CommonSetOracleTotalPower(ctx)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeDelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, oracle.DelegateValidator),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OracleAddress),
		),
	})

	return &types.MsgAddOracleDelegateResponse{}, nil
}

func (s EthereumMsgServer) EditOracle(c context.Context, msg *types.MsgEditOracle) (*types.MsgEditOracleResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNoFoundOracle, msg.OracleAddress)
	}
	if oracle.Jailed {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle jailed")
	}
	if oracle.IsValidator {
		if msg.ValidatorAddress != "" && msg.ValidatorAddress != msg.OracleAddress {
			return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle is a validator, cannot edit validator address")
		}
	} else {
		if msg.ValidatorAddress != "" && msg.ValidatorAddress != msg.OracleAddress {
			delegateAddress := types.GetOracleDelegateAddress(msg.ChainName, oracleAddr)
			valSrcAddress, err := sdk.ValAddressFromBech32(oracle.DelegateValidator)
			if err != nil {
				return nil, sdkerrors.Wrap(types.ErrInvalid, "validator address")
			}
			valDestAddress, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
			if err != nil {
				return nil, sdkerrors.Wrap(types.ErrInvalid, "validator address")
			}
			sharesAmount, err := s.stakingKeeper.ValidateUnbondAmount(ctx, delegateAddress, valSrcAddress, oracle.DelegateAmount)
			if err != nil {
				return nil, err
			}
			completionTime, err := s.stakingKeeper.BeginRedelegation(ctx, delegateAddress, valSrcAddress, valDestAddress, sharesAmount)
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
		}
	}
	//TODO implement me edit bridger and external address
	return &types.MsgEditOracleResponse{}, err
}

func (s EthereumMsgServer) WithdrawReward(c context.Context, msg *types.MsgWithdrawReward) (*types.MsgWithdrawRewardResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNoFoundOracle, msg.OracleAddress)
	}
	if oracle.Jailed {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle jailed")
	}
	if oracle.IsValidator {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle is a validator, cannot withdraw reward here")
	}
	validatorAddr, err := sdk.ValAddressFromBech32(oracle.DelegateValidator)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "validator address")
	}

	deleteAddr := types.GetOracleDelegateAddress(msg.ChainName, oracleAddr)
	rewards, err := s.distributionKeeper.WithdrawDelegationRewards(ctx, deleteAddr, validatorAddr)
	if err != nil {
		return nil, err
	}
	if err = s.bankKeeper.SendCoins(ctx, deleteAddr, oracleAddr, rewards); err != nil {
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

// SendToExternal handles MsgSendToExternal
func (s EthereumMsgServer) SendToExternal(c context.Context, msg *types.MsgSendToExternal) (*types.MsgSendToExternalResponse, error) {
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

func (s EthereumMsgServer) CancelSendToExternal(c context.Context, msg *types.MsgCancelSendToExternal) (*types.MsgCancelSendToExternalResponse, error) {
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

// RequestBatch handles MsgRequestBatch
func (s EthereumMsgServer) RequestBatch(c context.Context, msg *types.MsgRequestBatch) (*types.MsgRequestBatchResponse, error) {
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

	batch, err := s.BuildOutgoingTxBatch(ctx, bridgeToken.Token, msg.FeeReceive, OutgoingTxBatchSize, msg.MinimumFee, *msg.BaseFee)
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

// ConfirmBatch handles MsgConfirmBatch
func (s EthereumMsgServer) ConfirmBatch(c context.Context, msg *types.MsgConfirmBatch) (*types.MsgConfirmBatchResponse, error) {
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
	key := s.SetBatchConfirm(ctx, oracleAddr, msg)

	_ = key
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.BridgerAddress),
	))

	return nil, nil
}

// OracleSetConfirm handles MsgOracleSetConfirm
func (s EthereumMsgServer) OracleSetConfirm(c context.Context, msg *types.MsgOracleSetConfirm) (*types.MsgOracleSetConfirmResponse, error) {
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
	key := s.SetOracleSetConfirm(ctx, oracleAddr, msg)

	_ = key
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.BridgerAddress),
	))

	return &types.MsgOracleSetConfirmResponse{}, nil
}

// SendToExternalClaim handles MsgSendToExternalClaim
// executed aka 'observed' and had its slashing window expire) that will never be cleaned up in the end block. This
// should not be a security risk as 'old' events can never execute but it does store spam in the chain.
func (s EthereumMsgServer) SendToExternalClaim(c context.Context, msg *types.MsgSendToExternalClaim) (*types.MsgSendToExternalClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := checkOrchestratorIsOracle(ctx, s.Keeper, msg.BridgerAddress); err != nil {
		return nil, err
	}

	anyMsg, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendToExternalClaimResponse{}, s.claimHandlerCommon(ctx, anyMsg, msg, msg.Type())
}

// SendToFxClaim handles MsgSendToFxClaim
// executed aka 'observed' and had it's slashing window expire) that will never be cleaned up in the endblocker. This
// should not be a security risk as 'old' events can never execute but it does store spam in the chain.
func (s EthereumMsgServer) SendToFxClaim(c context.Context, msg *types.MsgSendToFxClaim) (*types.MsgSendToFxClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := checkOrchestratorIsOracle(ctx, s.Keeper, msg.BridgerAddress); err != nil {
		return nil, err
	}

	anyMsg, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendToFxClaimResponse{}, s.claimHandlerCommon(ctx, anyMsg, msg, msg.Type())
}

func (s EthereumMsgServer) BridgeTokenClaim(c context.Context, msg *types.MsgBridgeTokenClaim) (*types.MsgBridgeTokenClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := checkOrchestratorIsOracle(ctx, s.Keeper, msg.BridgerAddress); err != nil {
		return nil, err
	}

	anyMsg, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgBridgeTokenClaimResponse{}, s.claimHandlerCommon(ctx, anyMsg, msg, msg.Type())
}

// OracleSetUpdateClaim handles claims for executing a oracle set update on Ethereum
func (s EthereumMsgServer) OracleSetUpdateClaim(c context.Context, msg *types.MsgOracleSetUpdatedClaim) (*types.MsgOracleSetUpdatedClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := checkOrchestratorIsOracle(ctx, s.Keeper, msg.BridgerAddress); err != nil {
		return nil, err
	}

	for _, member := range msg.Members {
		if _, found := s.GetOracleByExternalAddress(ctx, member.ExternalAddress); !found {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "member oracle")
		}
	}

	anyMsg, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgOracleSetUpdatedClaimResponse{}, s.claimHandlerCommon(ctx, anyMsg, msg, msg.Type())
}

func checkOrchestratorIsOracle(ctx sdk.Context, keeper Keeper, bridgerAddress string) error {
	bridgerAddr, err := sdk.AccAddressFromBech32(bridgerAddress)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalid, "bridger address")
	}
	oracleAddr, found := keeper.GetOracleAddressByBridgerKey(ctx, bridgerAddr)
	if !found {
		return sdkerrors.Wrapf(types.ErrNoFoundOracle, "by bridger address: %s", bridgerAddress)
	}
	oracle, found := keeper.GetOracle(ctx, oracleAddr)
	if !found {
		return sdkerrors.Wrap(types.ErrNoFoundOracle, oracleAddr.String())
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

func (s EthereumMsgServer) confirmHandlerCommon(ctx sdk.Context, bridgerAddr sdk.AccAddress, signatureAddr, signature string, checkpoint []byte) (oracleAddr sdk.AccAddress, err error) {
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature decoding")
	}

	oracleAddr, found := s.GetOracleByExternalAddress(ctx, signatureAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNoFoundOracle, "by bridger address")
	}

	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrNoFoundOracle, oracleAddr.String())
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
	return oracle.GetOracle(), nil
}
