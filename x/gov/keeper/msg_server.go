package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/functionx/fx-core/v7/x/gov/types"
)

type msgServer struct {
	v1.MsgServer
	*Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(m v1.MsgServer, k *Keeper) types.MsgServerPro {
	return &msgServer{MsgServer: m, Keeper: k}
}

var _ types.MsgServerPro = msgServer{}

func (k msgServer) SubmitProposal(goCtx context.Context, msg *v1.MsgSubmitProposal) (*v1.MsgSubmitProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	msgInitialDeposit := msg.GetInitialDeposit()

	if err := k.validateInitialDeposit(ctx, msgInitialDeposit); err != nil {
		return nil, err
	}

	proposalMsgs, err := msg.GetMsgs()
	if err != nil {
		return nil, err
	}

	proposer, err := sdk.AccAddressFromBech32(msg.GetProposer())
	if err != nil {
		return nil, err
	}

	msgType := ""
	for _, pMsg := range proposalMsgs {
		if msgType != "" && !strings.EqualFold(msgType, sdk.MsgTypeURL(pMsg)) {
			return nil, govtypes.ErrInvalidProposalType.Wrapf("proposal MsgTypeURL is different")
		}
		msgType = sdk.MsgTypeURL(pMsg)
	}

	proposal, err := k.Keeper.SubmitProposal(ctx, proposalMsgs, msg.Metadata, msg.Title, msg.Summary, proposer)
	if err != nil {
		return nil, err
	}

	proposalBytes, err := proposal.Marshal()
	if err != nil {
		return nil, err
	}

	ctx.GasMeter().ConsumeGas(3*ctx.KVGasConfig().WriteCostPerByte*uint64(len(proposalBytes)),
		"submit proposal")

	defer telemetry.IncrCounter(1, govtypes.ModuleName, "proposal")

	minInitialDeposit := k.Keeper.GetMinInitialDeposit(ctx, types.ExtractMsgTypeURL(proposal.Messages))

	if sdk.Coins(msgInitialDeposit).IsAllLT(sdk.NewCoins(minInitialDeposit)) {
		return nil, errorsmod.Wrapf(types.ErrInitialAmountTooLow, "%s is smaller than %s", msgInitialDeposit, minInitialDeposit)
	}

	votingStarted, err := k.Keeper.AddDeposit(ctx, proposal.Id, proposer, msgInitialDeposit)
	if err != nil {
		return nil, err
	}

	if votingStarted {
		submitEvent := sdk.NewEvent(govtypes.EventTypeSubmitProposal,
			sdk.NewAttribute(govtypes.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", proposal.Id)),
		)

		ctx.EventManager().EmitEvent(submitEvent)
	}

	return &v1.MsgSubmitProposalResponse{
		ProposalId: proposal.Id,
	}, nil
}

func (k msgServer) Deposit(goCtx context.Context, msg *v1.MsgDeposit) (*v1.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddr, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	votingStarted, err := k.Keeper.AddDeposit(ctx, msg.ProposalId, accAddr, msg.Amount)
	if err != nil {
		return nil, err
	}

	defer telemetry.IncrCounterWithLabels(
		[]string{govtypes.ModuleName, "deposit"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("proposal_id", strconv.Itoa(int(msg.ProposalId))),
		},
	)

	if votingStarted {
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			govtypes.EventTypeProposalDeposit,
			sdk.NewAttribute(govtypes.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", msg.ProposalId)),
		))
	}

	return &v1.MsgDepositResponse{}, nil
}

func (k msgServer) UpdateFXParams(c context.Context, req *types.MsgUpdateFXParams) (*types.MsgUpdateFXParamsResponse, error) {
	if k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(c)
	if err := k.SetFXParams(ctx, &req.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateFXParamsResponse{}, nil
}

func (k msgServer) UpdateEGFParams(c context.Context, req *types.MsgUpdateEGFParams) (*types.MsgUpdateEGFParamsResponse, error) {
	if k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(c)
	if err := k.SetEGFParams(ctx, &req.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateEGFParamsResponse{}, nil
}

func (k msgServer) UpdateStore(c context.Context, req *types.MsgUpdateStore) (*types.MsgUpdateStoreResponse, error) {
	if k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(c)
	for _, updateStore := range req.UpdateStores {
		key, ok := k.Keeper.storeKeys[updateStore.Space]
		if !ok {
			return nil, errortypes.ErrInvalidRequest.Wrap("invalid store space")
		}
		kvStore := ctx.KVStore(key)
		keyBt := updateStore.KeyToBytes()
		storeValue := kvStore.Get(keyBt)
		if !bytes.Equal(storeValue, updateStore.OldValueToBytes()) {
			return nil, errortypes.ErrInvalidRequest.Wrapf("old value not equal store value: %s", hex.EncodeToString(storeValue))
		}
		kvStore.Set(keyBt, updateStore.ValueToBytes())
	}
	return &types.MsgUpdateStoreResponse{}, nil
}

func (k msgServer) UpdateSwitchParams(c context.Context, req *types.MsgUpdateSwitchParams) (*types.MsgUpdateSwitchParamsResponse, error) {
	if k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(c)
	if err := k.SetSwitchParams(ctx, &req.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateSwitchParamsResponse{}, nil
}

// validateInitialDeposit validates if initial deposit is greater than or equal to the minimum
// required at the time of proposal submission. This threshold amount is determined by
// the deposit parameters. Returns nil on success, error otherwise.
func (keeper Keeper) validateInitialDeposit(ctx sdk.Context, initialDeposit sdk.Coins) error {
	params := keeper.GetParams(ctx)
	minInitialDepositRatio, err := sdk.NewDecFromStr(params.MinInitialDepositRatio)
	if err != nil {
		return err
	}
	if minInitialDepositRatio.IsZero() {
		return nil
	}
	minDepositCoins := params.MinDeposit
	for i := range minDepositCoins {
		minDepositCoins[i].Amount = sdk.NewDecFromInt(minDepositCoins[i].Amount).Mul(minInitialDepositRatio).RoundInt()
	}
	if !initialDeposit.IsAllGTE(minDepositCoins) {
		return govtypes.ErrMinDepositTooSmall.Wrapf("was (%s), need (%s)", initialDeposit, minDepositCoins)
	}
	return nil
}
