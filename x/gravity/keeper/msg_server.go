package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	crosschainkeeper "github.com/functionx/fx-core/x/crosschain/keeper"

	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/x/eth/types"
	"github.com/functionx/fx-core/x/gravity/types"
)

type legacyMsgServer interface {
	SendToExternal(c context.Context, msg *crosschaintypes.MsgSendToExternal) (*crosschaintypes.MsgSendToExternalResponse, error)
}

var _ types.MsgServer = msgServer{}
var _ crosschaintypes.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

func NewLegacyMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) SendToEth(ctx context.Context, msg *types.MsgSendToEth) (*types.MsgSendToEthResponse, error) {
	_, err := k.legacyMsgServer.SendToExternal(ctx, &crosschaintypes.MsgSendToExternal{
		Sender:    msg.Sender,
		Dest:      msg.EthDest,
		Amount:    msg.Amount,
		BridgeFee: msg.BridgeFee,
		ChainName: ethtypes.ModuleName,
	})
	if err != nil {
		return nil, err
	}
	return &types.MsgSendToEthResponse{}, nil
}

func NewMsgServerImpl(keeper Keeper) crosschainkeeper.ProposalMsgServer {
	return &crosschainkeeper.EthereumMsgServer{Keeper: keeper.Keeper}
}

func (k msgServer) CreateOracleBridger(c context.Context, msg *crosschaintypes.MsgCreateOracleBridger) (*crosschaintypes.MsgCreateOracleBridgerResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, "oracle address")
	}
	bridgerAddr, err := sdk.AccAddressFromBech32(msg.BridgerAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, "bridger address")
	}
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, "validator address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	if !k.IsOracle(ctx, msg.OracleAddress) {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrNoFoundOracle, msg.OracleAddress)
	}
	// check oracle has set bridger address
	if _, found := k.GetOracle(ctx, oracleAddr); found {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, "oracle existed bridger address")
	}
	oracleIsValidator := false
	_, found := k.stakingKeeper.GetValidator(ctx, oracleAddr.Bytes())
	if found {
		oracleIsValidator = true
		if msg.OracleAddress != msg.ValidatorAddress {
			return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, "oracle is a validator but validator address is not itself")
		}
	}

	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	// check bridger address is bound to oracle
	if _, found := k.GetOracleAddressByBridgerKey(ctx, bridgerAddr); found {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, "bridger address is bound to oracle")
	}
	// check external address is bound to oracle
	if _, found := k.GetOracleByExternalAddress(ctx, msg.ExternalAddress); found {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, "external address is bound to oracle")
	}

	threshold := k.GetOracleDelegateThreshold(ctx)
	if threshold.Denom != msg.DelegateAmount.Denom {
		return nil, sdkerrors.Wrapf(crosschaintypes.ErrInvalid, "delegate denom, got %k, expected %k", msg.DelegateAmount.Denom, threshold.Denom)
	}
	if msg.DelegateAmount.IsLT(threshold) {
		return nil, crosschaintypes.ErrDelegateAmountBelowMinimum
	}
	if msg.DelegateAmount.Amount.GT(threshold.Amount.Mul(sdk.NewInt(k.GetOracleDelegateMultiple(ctx)))) {
		return nil, crosschaintypes.ErrDelegateAmountBelowMaximum
	}

	deleteAddr := crosschaintypes.GetOracleDelegateAddress(msg.ChainName, oracleAddr)
	newShares, err := k.stakingKeeper.Delegate(ctx, deleteAddr, msg.DelegateAmount.Amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}

	oracle := crosschaintypes.Oracle{
		OracleAddress:     oracleAddr.String(),
		BridgerAddress:    bridgerAddr.String(),
		ExternalAddress:   msg.ExternalAddress,
		DelegateAmount:    msg.DelegateAmount,
		StartHeight:       ctx.BlockHeight(),
		Jailed:            false,
		JailedHeight:      0,
		DelegateValidator: msg.ValidatorAddress,
		OracleIsValidator: oracleIsValidator,
	}
	// save oracle
	k.SetOracle(ctx, oracle)
	// set the bridger address
	k.SetOracleByBridger(ctx, oracleAddr, bridgerAddr)
	// set the ethereum address
	k.SetExternalAddressForOracle(ctx, oracleAddr, msg.ExternalAddress)
	// save total stake amount
	totalStake := k.GetTotalDelegate(ctx)
	k.SetTotalDelegate(ctx, totalStake.Add(msg.DelegateAmount))

	k.CommonSetOracleTotalPower(ctx)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			stakingtypes.EventTypeDelegate,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.DelegateAmount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OracleAddress),
		),
	})

	return &crosschaintypes.MsgCreateOracleBridgerResponse{}, nil
}

func (k msgServer) AddOracleDelegate(c context.Context, msg *crosschaintypes.MsgAddOracleDelegate) (*crosschaintypes.MsgAddOracleDelegateResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, "oracle address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	oracle, found := k.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrNoFoundOracle, msg.OracleAddress)
	}
	valAddr, err := sdk.ValAddressFromBech32(oracle.DelegateValidator)
	if err != nil {
		return nil, sdkerrors.Wrap(crosschaintypes.ErrInvalid, "validator address")
	}
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	threshold := k.GetOracleDelegateThreshold(ctx)
	// check stake denom
	if threshold.Denom != msg.Amount.Denom {
		return nil, sdkerrors.Wrapf(crosschaintypes.ErrInvalid, "delegate denom, got %k, expected %k", msg.Amount.Denom, threshold.Denom)
	}
	// check oracle total delegateAmount grate then minimum delegateAmount amount
	delegateAmount := oracle.DelegateAmount.Add(msg.Amount)
	if delegateAmount.Amount.Sub(threshold.Amount).IsNegative() {
		return nil, crosschaintypes.ErrDelegateAmountBelowMinimum
	}
	if delegateAmount.Amount.GT(threshold.Amount.Mul(sdk.NewInt(k.GetOracleDelegateMultiple(ctx)))) {
		return nil, crosschaintypes.ErrDelegateAmountBelowMaximum
	}

	totalDelegateAmount := k.GetTotalDelegate(ctx)
	totalDelegateAmount = totalDelegateAmount.Add(msg.Amount)

	deleteAddr := crosschaintypes.GetOracleDelegateAddress(msg.ChainName, oracleAddr)
	newShares, err := k.stakingKeeper.Delegate(ctx, deleteAddr, msg.Amount.Amount, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}
	// save new total delegateAmount
	k.SetTotalDelegate(ctx, totalDelegateAmount)
	if oracle.Jailed {
		oracle.Jailed = false
		oracle.StartHeight = ctx.BlockHeight()
	}
	// save oracle new delegateAmount
	oracle.DelegateAmount = delegateAmount
	k.SetOracle(ctx, oracle)

	k.CommonSetOracleTotalPower(ctx)

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

	return &crosschaintypes.MsgAddOracleDelegateResponse{}, nil
}

func (k msgServer) EditOracle(ctx context.Context, oracle *crosschaintypes.MsgEditOracle) (*crosschaintypes.MsgEditOracleResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k msgServer) WithdrawReward(ctx context.Context, reward *crosschaintypes.MsgWithdrawReward) (*crosschaintypes.MsgWithdrawRewardResponse, error) {
	//TODO implement me
	panic("implement me")
}
