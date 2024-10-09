package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/functionx/fx-core/v8/x/gov/types"
)

type msgServer struct {
	v1.MsgServer
	*Keeper
}

// NewMsgServerImpl returns an implementation of the gov msgServer interface
// for the provided Keeper.
func NewMsgServerImpl(m v1.MsgServer, k *Keeper) types.MsgServerPro {
	return &msgServer{MsgServer: m, Keeper: k}
}

var _ types.MsgServerPro = msgServer{}

//nolint:gocyclo
func (k msgServer) SubmitProposal(goCtx context.Context, msg *v1.MsgSubmitProposal) (*v1.MsgSubmitProposalResponse, error) {
	if msg.Title == "" {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("proposal title cannot be empty")
	}
	if msg.Summary == "" {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("proposal summary cannot be empty")
	}

	proposer, err := k.authKeeper.AddressCodec().StringToBytes(msg.GetProposer())
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid proposer address: %s", err)
	}

	// check that either metadata or Msgs length is non nil.
	if len(msg.Messages) == 0 && len(msg.Metadata) == 0 {
		return nil, govtypes.ErrNoProposalMsgs.Wrap("either metadata or Msgs length must be non-nil")
	}

	// verify that if present, the metadata title and summary equals the proposal title and summary
	if len(msg.Metadata) != 0 {
		proposalMetadata := govtypes.ProposalMetadata{}
		if err := json.Unmarshal([]byte(msg.Metadata), &proposalMetadata); err == nil {
			if proposalMetadata.Title != msg.Title {
				return nil, govtypes.ErrInvalidProposalContent.Wrapf("metadata title '%s' must equal proposal title '%s'", proposalMetadata.Title, msg.Title)
			}

			if proposalMetadata.Summary != msg.Summary {
				return nil, govtypes.ErrInvalidProposalContent.Wrapf("metadata summary '%s' must equal proposal summary '%s'", proposalMetadata.Summary, msg.Summary)
			}
		}

		// if we can't unmarshal the metadata, this means the client didn't use the recommended metadata format
		// nothing can be done here, and this is still a valid case, so we ignore the error
	}

	proposalMsgs, err := msg.GetMsgs()
	if err != nil {
		return nil, err
	}

	// check that all proposal messages are of the same type
	if err = checkProposalMsgs(proposalMsgs); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	initialDeposit := msg.GetInitialDeposit()

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get governance parameters: %w", err)
	}

	if err = k.validateInitialDeposit(params, initialDeposit, msg.Expedited); err != nil {
		return nil, err
	}

	if err = k.validateDepositDenom(params, initialDeposit); err != nil {
		return nil, err
	}

	proposal, err := k.Keeper.SubmitProposal(ctx, proposalMsgs, msg.Metadata, msg.Title, msg.Summary, proposer, msg.Expedited)
	if err != nil {
		return nil, err
	}

	proposalBytes, err := proposal.Marshal()
	if err != nil {
		return nil, err
	}

	// ref: https://github.com/cosmos/cosmos-sdk/issues/9683
	ctx.GasMeter().ConsumeGas(
		3*ctx.KVGasConfig().WriteCostPerByte*uint64(len(proposalBytes)),
		"submit proposal",
	)

	votingStarted, err := k.AddDeposit(ctx, proposal.Id, proposer, msg.GetInitialDeposit())
	if err != nil {
		return nil, err
	}

	if votingStarted {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(govtypes.EventTypeSubmitProposal,
				sdk.NewAttribute(govtypes.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", proposal.Id)),
			),
		)
	}

	return &v1.MsgSubmitProposalResponse{
		ProposalId: proposal.Id,
	}, nil
}

func (k msgServer) Deposit(goCtx context.Context, msg *v1.MsgDeposit) (*v1.MsgDepositResponse, error) {
	accAddr, err := k.authKeeper.AddressCodec().StringToBytes(msg.Depositor)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid depositor address: %s", err)
	}

	if err := validateDeposit(msg.Amount); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	votingStarted, err := k.AddDeposit(ctx, msg.ProposalId, accAddr, msg.Amount)
	if err != nil {
		return nil, err
	}

	if votingStarted {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				govtypes.EventTypeProposalDeposit,
				sdk.NewAttribute(govtypes.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", msg.ProposalId)),
			),
		)
	}

	return &v1.MsgDepositResponse{}, nil
}

func (k msgServer) UpdateStore(c context.Context, req *types.MsgUpdateStore) (*types.MsgUpdateStoreResponse, error) {
	if k.authority != req.Authority {
		return nil, govtypes.ErrInvalidSigner.Wrapf("invalid authority; expected %s, got %s", k.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(c)
	for _, updateStore := range req.UpdateStores {
		key, ok := k.Keeper.storeKeys[updateStore.Space]
		if !ok {
			return nil, sdkerrors.ErrInvalidRequest.Wrap("invalid store space")
		}
		kvStore := ctx.KVStore(key)
		keyBt := updateStore.KeyToBytes()
		storeValue := kvStore.Get(keyBt)
		if !bytes.Equal(storeValue, updateStore.OldValueToBytes()) {
			return nil, sdkerrors.ErrInvalidRequest.Wrapf("old value not equal store value: %s", hex.EncodeToString(storeValue))
		}
		kvStore.Set(keyBt, updateStore.ValueToBytes())
	}
	return &types.MsgUpdateStoreResponse{}, nil
}

func (k msgServer) UpdateSwitchParams(c context.Context, req *types.MsgUpdateSwitchParams) (*types.MsgUpdateSwitchParamsResponse, error) {
	if k.authority != req.Authority {
		return nil, govtypes.ErrInvalidSigner.Wrapf("invalid authority; expected %s, got %s", k.authority, req.Authority)
	}
	ctx := sdk.UnwrapSDKContext(c)
	if err := k.SetSwitchParams(ctx, &req.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateSwitchParamsResponse{}, nil
}

func (k msgServer) UpdateCustomParams(ctx context.Context, req *types.MsgUpdateCustomParams) (*types.MsgUpdateCustomParamsResponse, error) {
	if k.authority != req.Authority {
		return nil, govtypes.ErrInvalidSigner.Wrapf("invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	// delete the message params if the params are empty
	if req.GetCustomParams() == (types.CustomParams{}) {
		if err := k.CustomerParams.Remove(ctx, req.MsgUrl); err != nil {
			return nil, err
		}

		return &types.MsgUpdateCustomParamsResponse{}, nil
	}

	if err := req.CustomParams.ValidateBasic(); err != nil {
		return nil, err
	}

	if err := k.CustomerParams.Set(ctx, req.MsgUrl, req.CustomParams); err != nil {
		return nil, err
	}

	return &types.MsgUpdateCustomParamsResponse{}, nil
}

func checkProposalMsgs(proposalMsgs []sdk.Msg) error {
	msgType := ""
	for _, pMsg := range proposalMsgs {
		if msgType != "" && !strings.EqualFold(msgType, sdk.MsgTypeURL(pMsg)) {
			return govtypes.ErrInvalidProposalType.Wrapf("proposal MsgTypeURL is different")
		}
		msgType = sdk.MsgTypeURL(pMsg)
	}
	return nil
}

// validateDeposit validates the deposit amount, do not use for initial deposit.
func validateDeposit(amount sdk.Coins) error {
	if !amount.IsValid() || !amount.IsAllPositive() {
		return sdkerrors.ErrInvalidCoins.Wrap(amount.String())
	}

	return nil
}
