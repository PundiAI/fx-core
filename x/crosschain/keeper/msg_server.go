package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/hashicorp/go-metrics"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
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
		return nil, types.ErrInvalid.Wrapf("oracle address")
	}
	bridgerAddr, err := sdk.AccAddressFromBech32(msg.BridgerAddress)
	if err != nil {
		return nil, types.ErrInvalid.Wrapf("bridger address")
	}
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, types.ErrInvalid.Wrapf("validator address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	if !s.IsProposalOracle(ctx, msg.OracleAddress) {
		return nil, types.ErrNoFoundOracle
	}
	// check oracle has set bridger address
	if s.HasOracle(ctx, oracleAddr) {
		return nil, types.ErrInvalid.Wrapf("oracle existed bridger address")
	}
	// check bridger address is bound to oracle
	if s.HasOracleAddrByBridgerAddr(ctx, bridgerAddr) {
		return nil, types.ErrInvalid.Wrapf("bridger address is bound to oracle")
	}
	// check external address is bound to oracle
	if s.HasOracleAddrByExternalAddr(ctx, msg.ExternalAddress) {
		return nil, types.ErrInvalid.Wrapf("external address is bound to oracle")
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
		return nil, types.ErrInvalid.Wrapf("delegate denom got %s, expected %s", msg.DelegateAmount.Denom, threshold.Denom)
	}
	if msg.DelegateAmount.IsLT(threshold) {
		return nil, types.ErrDelegateAmountBelowMinimum
	}
	if msg.DelegateAmount.Amount.GT(threshold.Amount.Mul(sdkmath.NewInt(s.GetOracleDelegateMultiple(ctx)))) {
		return nil, types.ErrDelegateAmountAboveMaximum
	}

	delegateAddr := oracle.GetDelegateAddress(s.moduleName)
	if err = s.bankKeeper.SendCoins(ctx, oracleAddr, delegateAddr, sdk.NewCoins(msg.DelegateAmount)); err != nil {
		return nil, err
	}
	msgDelegate := stakingtypes.NewMsgDelegate(delegateAddr.String(), valAddr.String(), msg.DelegateAmount)
	if _, err = s.stakingMsgServer.Delegate(ctx, msgDelegate); err != nil {
		return nil, err
	}

	s.SetOracle(ctx, oracle)
	s.SetOracleAddrByBridgerAddr(ctx, bridgerAddr, oracleAddr)
	s.SetOracleAddrByExternalAddr(ctx, msg.ExternalAddress, oracleAddr)
	s.SetLastTotalPower(ctx)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.OracleAddress),
	))

	return &types.MsgBondedOracleResponse{}, nil
}

//nolint:gocyclo // need to refactor
func (s MsgServer) AddDelegate(c context.Context, msg *types.MsgAddDelegate) (*types.MsgAddDelegateResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		return nil, types.ErrInvalid.Wrapf("oracle address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	if !s.IsProposalOracle(ctx, msg.OracleAddress) {
		return nil, types.ErrNoFoundOracle
	}
	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, types.ErrNoFoundOracle
	}

	threshold := s.GetOracleDelegateThreshold(ctx)

	if threshold.Denom != msg.Amount.Denom {
		return nil, types.ErrInvalid.Wrapf("delegate denom got %s, expected %s", msg.Amount.Denom, threshold.Denom)
	}

	slashAmount := types.NewDelegateAmount(oracle.GetSlashAmount(s.GetSlashFraction(ctx)))
	if slashAmount.IsPositive() && msg.Amount.Amount.LT(slashAmount.Amount) {
		return nil, types.ErrInvalid.Wrapf("not sufficient slash amount")
	}

	delegateCoin := types.NewDelegateAmount(msg.Amount.Amount.Sub(slashAmount.Amount))

	oracle.DelegateAmount = oracle.DelegateAmount.Add(delegateCoin.Amount)
	if oracle.DelegateAmount.Sub(threshold.Amount).IsNegative() {
		return nil, types.ErrDelegateAmountBelowMinimum
	}
	if oracle.DelegateAmount.GT(threshold.Amount.Mul(sdkmath.NewInt(s.GetOracleDelegateMultiple(ctx)))) {
		return nil, types.ErrDelegateAmountAboveMaximum
	}

	if slashAmount.IsPositive() {
		if err = s.bankKeeper.SendCoinsFromAccountToModule(ctx, oracleAddr, s.moduleName, sdk.NewCoins(slashAmount)); err != nil {
			return nil, err
		}
		if err = s.bankKeeper.BurnCoins(ctx, s.moduleName, sdk.NewCoins(slashAmount)); err != nil {
			return nil, err
		}
	}

	if delegateCoin.IsPositive() {
		delegateAddr := oracle.GetDelegateAddress(s.moduleName)
		if err = s.bankKeeper.SendCoins(ctx, oracleAddr, delegateAddr, sdk.NewCoins(delegateCoin)); err != nil {
			return nil, err
		}
		msgDelegate := stakingtypes.NewMsgDelegate(delegateAddr.String(), oracle.GetValidator().String(), delegateCoin)
		if _, err = s.stakingMsgServer.Delegate(c, msgDelegate); err != nil {
			return nil, err
		}
	}

	if !oracle.Online {
		oracle.Online = true
		oracle.StartHeight = ctx.BlockHeight()
		if !ctx.IsCheckTx() {
			telemetry.SetGaugeWithLabels(
				[]string{types.ModuleName, "oracle_status"},
				float32(0),
				[]metrics.Label{
					telemetry.NewLabel("module", s.moduleName),
					telemetry.NewLabel("address", oracle.OracleAddress),
				},
			)
		}
	}
	oracle.SlashTimes = 0

	s.SetOracle(ctx, oracle)
	s.SetLastTotalPower(ctx)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.OracleAddress),
		),
	)

	return &types.MsgAddDelegateResponse{}, nil
}

func (s MsgServer) ReDelegate(c context.Context, msg *types.MsgReDelegate) (*types.MsgReDelegateResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		return nil, types.ErrInvalid.Wrapf("oracle address")
	}
	valDstAddress, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, types.ErrInvalid.Wrapf("validator address")
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
		return nil, types.ErrInvalid.Wrapf("validator address is not changed")
	}
	delegateAddr := oracle.GetDelegateAddress(s.moduleName)
	valSrcAddress := oracle.GetValidator()
	delegateToken, err := s.GetOracleDelegateToken(ctx, delegateAddr, valSrcAddress)
	if err != nil {
		return nil, err
	}
	msgBeginRedelegate := stakingtypes.NewMsgBeginRedelegate(delegateAddr.String(), valSrcAddress.String(), valDstAddress.String(), types.NewDelegateAmount(delegateToken))
	if _, err = s.stakingMsgServer.BeginRedelegate(c, msgBeginRedelegate); err != nil {
		return nil, err
	}
	oracle.DelegateValidator = msg.ValidatorAddress
	s.SetOracle(ctx, oracle)
	return &types.MsgReDelegateResponse{}, err
}

func (s MsgServer) EditBridger(c context.Context, msg *types.MsgEditBridger) (*types.MsgEditBridgerResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		return nil, types.ErrInvalid.Wrapf("oracle address")
	}
	bridgerAddr, err := sdk.AccAddressFromBech32(msg.BridgerAddress)
	if err != nil {
		return nil, types.ErrInvalid.Wrapf("oracle address")
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
		return nil, types.ErrInvalid.Wrapf("bridger address is not changed")
	}
	if s.HasOracleAddrByBridgerAddr(ctx, bridgerAddr) {
		return nil, types.ErrInvalid.Wrapf("bridger address is bound to oracle")
	}
	s.DelOracleAddrByBridgerAddr(ctx, oracle.GetBridger())
	oracle.BridgerAddress = msg.BridgerAddress
	s.SetOracle(ctx, oracle)
	s.SetOracleAddrByBridgerAddr(ctx, bridgerAddr, oracleAddr)

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
		return nil, types.ErrInvalid.Wrapf("oracle address")
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
	msgWithdrawDelegatorReward := distributiontypes.NewMsgWithdrawDelegatorReward(delegateAddr.String(), oracle.GetValidator().String())
	if _, err = s.distributionKeeper.WithdrawDelegatorReward(c, msgWithdrawDelegatorReward); err != nil {
		return nil, err
	}
	balances := s.bankKeeper.GetAllBalances(ctx, delegateAddr)
	if !balances.IsAllPositive() {
		return nil, types.ErrInvalid.Wrapf("rewards is empty")
	}
	if err = s.bankKeeper.SendCoins(ctx, delegateAddr, oracleAddr, balances); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.OracleAddress),
	))
	return &types.MsgWithdrawRewardResponse{}, nil
}

func (s MsgServer) UnbondedOracle(c context.Context, msg *types.MsgUnbondedOracle) (*types.MsgUnbondedOracleResponse, error) {
	oracleAddr, err := sdk.AccAddressFromBech32(msg.OracleAddress)
	if err != nil {
		return nil, types.ErrInvalid.Wrapf("oracle address")
	}
	ctx := sdk.UnwrapSDKContext(c)
	if s.IsProposalOracle(ctx, msg.OracleAddress) {
		return nil, types.ErrInvalid.Wrapf("need to pass a proposal to unbind")
	}
	oracle, found := s.GetOracle(ctx, oracleAddr)
	if !found {
		return nil, types.ErrNoFoundOracle
	}
	if oracle.Online {
		return nil, types.ErrInvalid.Wrapf("oracle on line")
	}
	delegateAddr := oracle.GetDelegateAddress(s.moduleName)
	validatorAddr := oracle.GetValidator()
	if _, err = s.stakingKeeper.GetUnbondingDelegation(ctx, delegateAddr, validatorAddr); err != nil {
		return nil, err
	}
	balances := s.bankKeeper.GetAllBalances(ctx, delegateAddr)
	slashAmount := types.NewDelegateAmount(oracle.GetSlashAmount(s.GetSlashFraction(ctx)))
	if slashAmount.IsPositive() {
		if balances.AmountOf(slashAmount.Denom).LT(slashAmount.Amount) {
			return nil, types.ErrInvalid.Wrapf("not sufficient slash amount")
		}
		if err = s.bankKeeper.SendCoinsFromAccountToModule(ctx, delegateAddr, s.moduleName, sdk.NewCoins(slashAmount)); err != nil {
			return nil, err
		}
		if err = s.bankKeeper.BurnCoins(ctx, s.moduleName, sdk.NewCoins(slashAmount)); err != nil {
			return nil, err
		}
	}
	sendCoins := balances.Sub(sdk.NewCoins(slashAmount)...)
	for i := 0; i < len(sendCoins); i++ {
		if !sendCoins[i].IsPositive() {
			sendCoins = append(sendCoins[:i], sendCoins[i+1:]...)
			i--
		}
	}
	if sendCoins.IsAllPositive() {
		if err = s.bankKeeper.SendCoins(ctx, delegateAddr, oracleAddr, sendCoins); err != nil {
			return nil, err
		}
	}

	s.DelOracleAddrByExternalAddr(ctx, oracle.ExternalAddress)
	s.DelOracleAddrByBridgerAddr(ctx, oracle.GetBridger())
	s.DelOracle(ctx, oracle.GetOracle())
	s.DelLastEventNonceByOracle(ctx, oracleAddr)

	return &types.MsgUnbondedOracleResponse{}, nil
}

// Deprecated: SendToExternal Please use precompile BridgeCall
func (s MsgServer) SendToExternal(c context.Context, msg *types.MsgSendToExternal) (*types.MsgSendToExternalResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, types.ErrInvalid.Wrapf("sender address")
	}

	ctx := sdk.UnwrapSDKContext(c)
	batchNonce, err := s.BuildOutgoingTxBatch(ctx, sender, msg.Dest, msg.Amount, msg.BridgeFee)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, msg.ChainName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

	return &types.MsgSendToExternalResponse{
		BatchNonce: batchNonce,
	}, nil
}

// Deprecated: ConfirmBatch Please use Confirm
func (s MsgServer) ConfirmBatch(c context.Context, msg *types.MsgConfirmBatch) (*types.MsgConfirmBatchResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	err := s.ConfirmHandler(ctx, msg)
	return &types.MsgConfirmBatchResponse{}, err
}

// Deprecated: ConfirmBatch Please use Confirm
func (s MsgServer) OracleSetConfirm(c context.Context, msg *types.MsgOracleSetConfirm) (*types.MsgOracleSetConfirmResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	err := s.ConfirmHandler(ctx, msg)
	return &types.MsgOracleSetConfirmResponse{}, err
}

// Deprecated: ConfirmBatch Please use Confirm
func (s MsgServer) BridgeCallConfirm(c context.Context, msg *types.MsgBridgeCallConfirm) (*types.MsgBridgeCallConfirmResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	err := s.ConfirmHandler(ctx, msg)
	return &types.MsgBridgeCallConfirmResponse{}, err
}

func (s MsgServer) UpdateParams(c context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if s.authority != req.Authority {
		return nil, govtypes.ErrInvalidSigner.Wrapf("invalid authority; expected %s, got %s", s.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.SetParams(ctx, &req.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateParamsResponse{}, nil
}

func (s MsgServer) UpdateChainOracles(c context.Context, req *types.MsgUpdateChainOracles) (*types.MsgUpdateChainOraclesResponse, error) {
	if s.authority != req.Authority {
		return nil, govtypes.ErrInvalidSigner.Wrapf("invalid authority; expected %s, got %s", s.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(c)
	if err := s.UpdateProposalOracles(ctx, req.Oracles); err != nil {
		return nil, err
	}
	return &types.MsgUpdateChainOraclesResponse{}, nil
}

func (s MsgServer) Claim(c context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	claim, ok := msg.Claim.GetCachedValue().(types.ExternalClaim)
	if !ok {
		return nil, types.ErrInvalid.Wrapf("invalid claim")
	}

	ctx := sdk.UnwrapSDKContext(c)
	bridgerAddr := claim.GetClaimer()
	oracleAddr, err := s.checkBridgerIsOracle(ctx, bridgerAddr)
	if err != nil {
		return nil, err
	}

	if err = s.claimLogicCheck(ctx, claim); err != nil {
		return nil, err
	}

	// Add the claim to the store
	if _, err = s.Attest(ctx, oracleAddr, claim); err != nil {
		return nil, err
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, s.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, bridgerAddr.String()),
	))
	return &types.MsgClaimResponse{}, nil
}

func (s MsgServer) Confirm(c context.Context, msg *types.MsgConfirm) (*types.MsgConfirmResponse, error) {
	confirm, ok := msg.Confirm.GetCachedValue().(types.Confirm)
	if !ok {
		return nil, types.ErrInvalid.Wrapf("invalid claim")
	}
	ctx := sdk.UnwrapSDKContext(c)
	err := s.ConfirmHandler(ctx, confirm)
	return &types.MsgConfirmResponse{}, err
}

func (s MsgServer) claimLogicCheck(ctx sdk.Context, claim types.ExternalClaim) (err error) {
	if claimMsg, ok := claim.(*types.MsgOracleSetUpdatedClaim); ok {
		for _, member := range claimMsg.Members {
			if !s.HasOracleAddrByExternalAddr(ctx, member.ExternalAddress) {
				return types.ErrInvalid.Wrapf("external address")
			}
		}
	}
	return nil
}

func (s MsgServer) checkBridgerIsOracle(ctx sdk.Context, bridgerAddr sdk.AccAddress) (oracleAddr sdk.AccAddress, err error) {
	oracleAddr, found := s.GetOracleAddrByBridgerAddr(ctx, bridgerAddr)
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
