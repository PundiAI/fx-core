package keeper

import (
	"math/big"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func (k Keeper) HookRelayToken(ctx sdk.Context, rtels []*RelayTokenEventLog, receipt *ethtypes.Receipt) error {
	fip20ABI := fxtypes.GetERC20().ABI
	for _, rtel := range rtels {
		k.Logger(ctx).Info("relay token", "hash", receipt.TxHash.String(), "from", rtel.Event.From.Hex(),
			"amount", rtel.Event.Value.String(), "denom", rtel.Pair.Denom, "token", rtel.Pair.Erc20Address)

		if err := k.ProcessRelayToken(ctx, fip20ABI, receipt.TxHash, rtel.Pair, rtel.Event.From, rtel.Event.Value); err != nil {
			k.Logger(ctx).Error("failed to relay token", "hash", receipt.TxHash.String(), "error", err.Error())
			return err
		}
		k.Logger(ctx).Info("relay transfer token success", "hash", receipt.TxHash.Hex())
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "relay_token"},
			1,
			[]metrics.Label{
				telemetry.NewLabel("erc20", rtel.Pair.Erc20Address),
				telemetry.NewLabel("denom", rtel.Pair.Denom),
				telemetry.NewLabel("amount", rtel.Event.Value.String()),
			},
		)
	}
	return nil
}

func (k Keeper) ProcessRelayToken(ctx sdk.Context, fip20ABI abi.ABI, txHash common.Hash, pair *types.TokenPair, from common.Address, amount *big.Int) error {
	var err error
	// create the corresponding sdk.Coin that is paired with FIP20
	coins := sdk.Coins{{Denom: pair.Denom, Amount: sdk.NewIntFromBigInt(amount)}}

	switch pair.ContractOwner {
	case types.OWNER_MODULE:
		if _, err = k.CallEVM(ctx, fip20ABI, k.moduleAddress, pair.GetERC20Contract(),
			true, "burn", k.moduleAddress, amount); err != nil {
			return err
		}

		if pair.Denom == fxtypes.DefaultDenom {
			if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, pair.GetERC20Contract().Bytes(), types.ModuleName, coins); err != nil {
				return err
			}
		}
	case types.OWNER_EXTERNAL:
		if err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
			return err
		}
	default:
		return types.ErrUndefinedOwner
	}

	//sender receive relay amount
	recipient := sdk.AccAddress(from.Bytes())
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins); err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeRelayToken,
				sdk.NewAttribute(sdk.AttributeKeySender, from.String()),
				sdk.NewAttribute(types.AttributeKeyReceiver, sdk.AccAddress(recipient.Bytes()).String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
				sdk.NewAttribute(types.AttributeKeyDenom, pair.Denom),
				sdk.NewAttribute(types.AttributeKeyTokenAddress, pair.Erc20Address),
				sdk.NewAttribute(types.AttributeKeyEvmTxHash, txHash.String()),
			),
		},
	)
	k.Logger(ctx).Info("relay token from evm success", "amount", amount.String(), "coins", coins.String(),
		"contract", pair.Erc20Address, "from", from.String(), "recipient", sdk.AccAddress(recipient.Bytes()).String())
	return nil
}

type RelayTokenEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

// ParseRelayTokenEvent transfer event ---> event Transfer(address indexed from, address indexed to, uint256 value);
func ParseRelayTokenEvent(fip20ABI abi.ABI, log *ethtypes.Log) (*RelayTokenEvent, common.Address, error) {
	// Note: the `Transfer` event contains 3 topics (id, from, to)
	if len(log.Topics) != 3 {
		return nil, common.Address{}, nil
	}
	eventID := log.Topics[0] // event ID
	event, err := fip20ABI.EventByID(eventID)
	if err != nil {
		return nil, common.Address{}, nil
	}
	if !(event.Name == types.ERC20EventTransfer) {
		return nil, common.Address{}, nil
	}
	toAddr := common.BytesToAddress(log.Topics[2].Bytes())

	relayTokenEvent := new(RelayTokenEvent)
	if log.Topics[0] != fip20ABI.Events[types.ERC20EventTransfer].ID {
		return nil, toAddr, nil
	}
	if len(log.Data) > 0 {
		if err := fip20ABI.UnpackIntoInterface(relayTokenEvent, types.ERC20EventTransfer, log.Data); err != nil {
			return nil, toAddr, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range fip20ABI.Events[types.ERC20EventTransfer].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(relayTokenEvent, indexed, log.Topics[1:]); err != nil {
		return nil, toAddr, err
	}
	return relayTokenEvent, toAddr, nil
}
