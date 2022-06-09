package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

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

func (s EthereumMsgServer) CreateOracleBridger(c context.Context, msg *types.MsgCreateOracleBridger) (*types.MsgCreateOracleBridgerResponse, error) {
	var err error
	var oracleAddr, bridgerAddr sdk.AccAddress
	if oracleAddr, err = sdk.AccAddressFromBech32(msg.OracleAddress); err != nil {
		return nil, sdkerrors.Wrap(types.ErrOracleAddress, msg.OracleAddress)
	}
	if bridgerAddr, err = sdk.AccAddressFromBech32(msg.BridgerAddress); err != nil {
		return nil, sdkerrors.Wrap(types.ErrBridgerAddress, msg.BridgerAddress)
	}
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(c)
	if !s.IsOracle(ctx, msg.OracleAddress) {
		return nil, types.ErrNotOracle
	}
	// check oracle has set orchestrator address
	if _, found := s.GetOracle(ctx, oracleAddr); found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "oracle existed orchestrator address")
	}
	validator, found := s.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	// check orchestrator address is bound to oracle
	if _, found := s.GetOracleAddressByBridgerKey(ctx, bridgerAddr); found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "orchestrator address is bound to oracle")
	}
	// check external address is bound to oracle
	if _, found := s.GetOracleByExternalAddress(ctx, msg.ExternalAddress); found {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "external address is bound to oracle")
	}

	threshold := s.GetOracleStakeThreshold(ctx)
	if threshold.Denom != msg.DelegateAmount.Denom {
		return nil, sdkerrors.Wrapf(types.ErrBadStakeDenom, "got %s, expected %s", msg.DelegateAmount.Denom, threshold.Denom)
	}
	if msg.DelegateAmount.IsLT(threshold) {
		return nil, types.ErrStakeAmountBelowMinimum
	}
	if msg.DelegateAmount.Amount.GT(threshold.Amount.Mul(sdk.NewInt(depositMultiple))) {
		return nil, types.ErrStakeAmountBelowMaximum
	}

	newShares, err := s.stakingKeeper.Delegate(ctx, oracleAddr, msg.DelegateAmount.Amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}

	oracle := types.Oracle{
		OracleAddress:   oracleAddr.String(),
		BridgerAddress:  bridgerAddr.String(),
		ExternalAddress: msg.ExternalAddress,
		DelegateAmount:  msg.DelegateAmount,
		StartHeight:     ctx.BlockHeight(),
		Jailed:          false,
		JailedHeight:    0,
	}
	// save oracle
	s.SetOracle(ctx, oracle)
	// set the orchestrator address
	s.SetOracleByBridger(ctx, oracleAddr, bridgerAddr)
	// set the ethereum address
	s.SetExternalAddressForOracle(ctx, oracleAddr, msg.ExternalAddress)
	// save total stake amount
	totalStake := s.GetTotalStake(ctx)
	s.SetTotalStake(ctx, totalStake.Add(msg.DelegateAmount))

	s.CommonSetOracleTotalPower(ctx)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeDelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.DelegateAmount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, s.moduleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OracleAddress),
		),
	})

	return &types.MsgCreateOracleBridgerResponse{}, nil
}

func (s EthereumMsgServer) AddOracleDelegate(c context.Context, msg *types.MsgAddOracleDelegate) (*types.MsgAddOracleDelegateResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
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
	valAddr, err := sdk.ValAddressFromBech32(oracle.DelegateValidator)
	if err != nil {
		return nil, err
	}
	validator, found := s.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	threshold := s.GetOracleStakeThreshold(ctx)
	// check stake denom
	if threshold.Denom != msg.Amount.Denom {
		return nil, sdkerrors.Wrapf(types.ErrBadStakeDenom, "got %s, expected %s", msg.Amount.Denom, threshold.Denom)
	}
	// check oracle total stakeAmount grate then minimum stakeAmount amount
	stakeAmount := oracle.DelegateAmount.Add(msg.Amount)
	if stakeAmount.Amount.Sub(threshold.Amount).IsNegative() {
		return nil, types.ErrStakeAmountBelowMinimum
	}
	if stakeAmount.Amount.GT(threshold.Amount.Mul(sdk.NewInt(depositMultiple))) {
		return nil, types.ErrStakeAmountBelowMaximum
	}

	totalStake := s.GetTotalStake(ctx)
	totalStake = totalStake.Add(msg.Amount)

	newShares, err := s.stakingKeeper.Delegate(ctx, oracleAddr, msg.Amount.Amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}
	// save new total stakeAmount
	s.SetTotalStake(ctx, totalStake)
	if oracle.Jailed {
		oracle.Jailed = false
		oracle.StartHeight = ctx.BlockHeight()
	}
	// save oracle new stakeAmount
	oracle.DelegateAmount = stakeAmount
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
			sdk.NewAttribute(sdk.AttributeKeyModule, s.moduleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OracleAddress),
		),
	})

	return &types.MsgAddOracleDelegateResponse{}, nil
}

func (s EthereumMsgServer) EditOracle(ctx context.Context, oracle *types.MsgEditOracle) (*types.MsgEditOracleResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s EthereumMsgServer) WithdrawReward(ctx context.Context, reward *types.MsgWithdrawReward) (*types.MsgWithdrawRewardResponse, error) {
	//TODO implement me
	panic("implement me")
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
	_, found := s.GetOracleAddressByBridgerKey(ctx, sender)
	if !found {
		if !s.IsOracle(ctx, msg.Sender) {
			return nil, sdkerrors.Wrap(types.ErrEmpty, "oracle or orchestrator")
		}
	}

	batch, err := s.BuildOutgoingTxBatch(ctx, bridgeToken.Token, msg.FeeReceive, OutgoingTxBatchSize, msg.MinimumFee, *msg.BaseFee)
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
	bridgerAddr, err := sdk.AccAddressFromBech32(msg.BridgerAddress)
	if err != nil {
		return nil, types.ErrBridgerAddress
	}
	checkpoint, err := batch.GetCheckpoint(s.GetGravityID(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "checkpoint generation")
	}

	oracleAddr, err := s.confirmHandlerCommon(ctx, bridgerAddr, msg.ExternalAddress, msg.Signature, checkpoint)
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
		sdk.NewAttribute(sdk.AttributeKeySender, msg.BridgerAddress),
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
	bridgerAddr, err := sdk.AccAddressFromBech32(msg.BridgerAddress)
	if err != nil {
		return nil, types.ErrBridgerAddress
	}
	checkpoint := oracleSet.GetCheckpoint(s.GetGravityID(ctx))
	oracleAddr, err := s.confirmHandlerCommon(ctx, bridgerAddr, msg.ExternalAddress, msg.Signature, checkpoint)
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

func checkOrchestratorIsOracle(ctx sdk.Context, keeper Keeper, orchestrator string) error {
	orcAddr, err := sdk.AccAddressFromBech32(orchestrator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "orchestrator")
	}
	oracleAddr, found := keeper.GetOracleAddressByBridgerKey(ctx, orcAddr)
	if !found {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "oracle")
	}

	oracle, found := keeper.GetOracle(ctx, oracleAddr)
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

func (s EthereumMsgServer) confirmHandlerCommon(ctx sdk.Context, bridgerAddr sdk.AccAddress, signatureAddr, signature string, checkpoint []byte) (oracleAddr sdk.AccAddress, err error) {
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
	if oracle.BridgerAddress != bridgerAddr.String() {
		return nil, sdkerrors.Wrapf(types.ErrInvalid, "got %s, expected %s", bridgerAddr, oracle.BridgerAddress)
	}
	if err = types.ValidateEthereumSignature(checkpoint, sigBytes, oracle.ExternalAddress); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s found %s", oracle.ExternalAddress, hex.EncodeToString(checkpoint), signature))
	}
	return oracle.GetOracle(), nil
}
