package keeper

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/functionx/fx-core/x/evm/types"

	fxtype "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/intrarelayer/types"
	"github.com/functionx/fx-core/x/intrarelayer/types/contracts"
)

var _ evmtypes.EvmHooks = (*Keeper)(nil)

// PostTxProcessing implements EvmHooks.PostTxProcessing
func (k Keeper) PostTxProcessing(ctx sdk.Context, txHash common.Hash, logs []*ethtypes.Log) error {
	if ctx.BlockHeight() < fxtype.IntrarelayerSupportBlock() || !k.HasInit(ctx) {
		return nil
	}
	params := k.GetParams(ctx)
	if !params.EnableEVMHook {
		return sdkerrors.Wrap(types.ErrInternalTokenPair, "EVM Hook is currently disabled")
	}
	//process relay event
	if err := k.RelayEventProcessing(ctx, txHash, logs); err != nil {
		return err
	}
	return nil
}

func (k *Keeper) RelayEventProcessing(ctx sdk.Context, txHash common.Hash, logs []*ethtypes.Log) error {
	for _, log := range logs {
		if !isRelayEvent(log) {
			continue
		}
		//check contract is registered
		pairID := k.GetERC20Map(ctx, log.Address)
		if len(pairID) == 0 { // contract is not registered coin or erc20
			continue
		}
		pair, found := k.GetTokenPair(ctx, pairID)
		if !found { //token pair info not found
			continue
		}
		//relay amount
		amount, err := parseRelayAmount(log.Data)
		if err != nil {
			k.Logger(ctx).Error("Unpack relay amount", "data", hex.EncodeToString(log.Data), "error", err.Error())
			return errors.New("invalid amount")
		}
		// create the corresponding sdk.Coin that is paired with ERC20
		coins := sdk.Coins{{Denom: pair.Denom, Amount: sdk.NewIntFromBigInt(amount)}}
		//relay from
		sender := common.BytesToAddress(log.Topics[1].Bytes())
		//relay to
		recipient := common.BytesToAddress(log.Topics[2].Bytes())
		k.Logger(ctx).Info("Relay erc20 from evm",
			"coins", coins.String(),
			"contract", pair.Erc20Address,
			"from", sender.String(),
			"recipient", recipient.String(),
			"accAddress", sdk.AccAddress(recipient.Bytes()).String())

		switch pair.ContractOwner {
		case types.OWNER_MODULE:
			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(recipient.Bytes()), coins)
		case types.OWNER_EXTERNAL:
			if err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
				panic(err)
			}
			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(recipient.Bytes()), coins)
		default:
			err = types.ErrUndefinedOwner
		}

		if err != nil {
			k.Logger(ctx).Error(
				"Process EVM hook for ER20 -> coin relay",
				"coin", pair.Denom, "contract", pair.Erc20Address, "error", err.Error(),
			)
			return err
		}
		ctx.EventManager().EmitEvents(
			sdk.Events{
				sdk.NewEvent(
					types.EventTypeERC20Relay,
					sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
					sdk.NewAttribute(types.AttributeKeyReceiver, sdk.AccAddress(recipient.Bytes()).String()),
					sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
					sdk.NewAttribute(types.AttributeKeyCosmosCoin, pair.Denom),
					sdk.NewAttribute(types.AttributeKeyERC20Token, pair.Erc20Address),
					sdk.NewAttribute(types.EventERC20RelayHash, txHash.String()),
				),
			},
		)
	}
	return nil
}

func isRelayEvent(log *ethtypes.Log) bool {
	if len(log.Topics) < 3 {
		return false
	}
	eventID := log.Topics[0] // event ID
	event, err := contracts.ERC20RelayContract.ABI.EventByID(eventID)
	if err != nil {
		// invalid event for ERC20Relay
		return false
	}
	return event.Name == types.ERC20RelayEventRelay
}
func parseRelayAmount(data []byte) (*big.Int, error) {
	//relay amount
	relayEvent, err := contracts.ERC20RelayContract.ABI.Unpack(types.ERC20RelayEventRelay, data)
	if err != nil {
		return nil, fmt.Errorf("unpack relay event error %v", err.Error())
	}
	if len(relayEvent) == 0 {
		return nil, errors.New("invalid relay event")
	}
	amount, ok := relayEvent[0].(*big.Int)
	if !ok || amount == nil {
		return nil, fmt.Errorf("invalid type of relay event")
	}
	if amount.Sign() != 1 {
		return nil, fmt.Errorf("invalid amount %v", amount)
	}
	return amount, nil
}
