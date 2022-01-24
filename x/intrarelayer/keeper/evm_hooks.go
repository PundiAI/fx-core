package keeper

import (
	"bytes"
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
	if err := k.RelayTokenProcessing(ctx, txHash, logs); err != nil {
		return err
	}
	return nil
}

// RelayTokenProcessing relay token from evm contract to chain address
func (k Keeper) RelayTokenProcessing(ctx sdk.Context, txHash common.Hash, logs []*ethtypes.Log) error {
	for _, log := range logs {
		if !isRelayTokenEvent(log) {
			continue
		}
		pair, found := k.GetTokenPairByAddress(ctx, log.Address)
		if !found {
			continue
		}
		amount, err := parseTransferAmount(log.Data)
		if err != nil {
			return fmt.Errorf("parse transfer amount error %v", err.Error())
		}
		from := common.BytesToAddress(log.Topics[1].Bytes())
		err = k.ProcessRelayToken(ctx, txHash, pair, from, amount)
		if err != nil {
			k.Logger(ctx).Error("Process EVM hook -> Relay token from evm", "amount", amount.String(),
				"coin", pair.Denom, "contract", pair.Erc20Address, "error", err.Error())
			return err
		}
	}
	return nil
}
func (k Keeper) GetTokenPairByAddress(ctx sdk.Context, address common.Address) (types.TokenPair, bool) {
	//check contract is registered
	pairID := k.GetERC20Map(ctx, address)
	if len(pairID) == 0 {
		// contract is not registered coin or erc20
		return types.TokenPair{}, false
	}
	return k.GetTokenPair(ctx, pairID)
}
func (k Keeper) ProcessRelayToken(ctx sdk.Context, txHash common.Hash, pair types.TokenPair, from common.Address, amount *big.Int) error {
	var err error
	// create the corresponding sdk.Coin that is paired with ERC20
	coins := sdk.Coins{{Denom: pair.Denom, Amount: sdk.NewIntFromBigInt(amount)}}

	switch pair.ContractOwner {
	case types.OWNER_MODULE:
		_, err = k.CallEVMWithModule(ctx, contracts.ERC20RelayContract.ABI, pair.GetERC20Contract(), "burn", amount)
	case types.OWNER_EXTERNAL:
		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	default:
		err = types.ErrUndefinedOwner
	}

	//sender receive relay amount
	recipient := sdk.AccAddress(from.Bytes())
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeRelayToken,
				sdk.NewAttribute(sdk.AttributeKeySender, from.String()),
				sdk.NewAttribute(types.AttributeKeyReceiver, sdk.AccAddress(recipient.Bytes()).String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
				sdk.NewAttribute(types.AttributeKeyCosmosCoin, pair.Denom),
				sdk.NewAttribute(types.AttributeKeyERC20Token, pair.Erc20Address),
				sdk.NewAttribute(types.EventEthereumTxHash, txHash.String()),
			),
		},
	)
	k.Logger(ctx).Info("Process EVM hook -> Relay token from evm", "amount", amount.String(), "coins", coins.String(),
		"contract", pair.Erc20Address, "from", from.String(), "recipient", sdk.AccAddress(recipient.Bytes()).String())
	return nil
}

//isRelayTokenEvent check transfer event is relay token
//transfer event ---> event Transfer(address indexed from, address indexed to, uint256 value);
//address to must be equal ModuleAddress
func isRelayTokenEvent(log *ethtypes.Log) bool {
	if len(log.Topics) < 3 {
		return false
	}
	eventID := log.Topics[0] // event ID
	event, err := contracts.ERC20RelayContract.ABI.EventByID(eventID)
	if err != nil {
		return false
	}
	if !(event.Name == types.ERC20EventTransfer) {
		return false
	}
	//transfer to module address
	to := common.BytesToAddress(log.Topics[2].Bytes())
	return bytes.Equal(to.Bytes(), types.ModuleAddress.Bytes())
}

//parseTransferAmount parse transfer event data to big int
func parseTransferAmount(data []byte) (*big.Int, error) {
	//relay amount
	transferEvent, err := contracts.ERC20RelayContract.ABI.Unpack(types.ERC20EventTransfer, data)
	if err != nil {
		return nil, fmt.Errorf("unpack transfer event error %v", err.Error())
	}
	if len(transferEvent) == 0 {
		return nil, errors.New("invalid transfer event")
	}
	amount, ok := transferEvent[0].(*big.Int)
	if !ok || amount == nil {
		return nil, fmt.Errorf("invalid type of transfer event")
	}
	if amount.Sign() != 1 {
		return nil, fmt.Errorf("invalid transfer amount %v", amount)
	}
	return amount, nil
}
