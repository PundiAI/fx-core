package precompile

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v8/contract"
	fxtypes "github.com/functionx/fx-core/v8/types"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
)

type Keeper struct {
	router            *Router
	bankKeeper        BankKeeper
	erc20Keeper       Erc20Keeper
	ibcTransferKeeper IBCTransferKeeper
	accountKeeper     AccountKeeper
}

func (c *Keeper) handlerOriginToken(ctx sdk.Context, _ *vm.EVM, sender common.Address, amount *big.Int) (sdk.Coin, error) {
	totalCoin := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(amount))
	totalCoins := sdk.NewCoins(totalCoin)
	if err := c.bankKeeper.SendCoinsFromAccountToModule(ctx, crosschaintypes.GetAddress().Bytes(), evmtypes.ModuleName, totalCoins); err != nil {
		return sdk.Coin{}, err
	}

	if err := c.bankKeeper.SendCoinsFromModuleToAccount(ctx, evmtypes.ModuleName, sender.Bytes(), totalCoins); err != nil {
		return sdk.Coin{}, err
	}
	return totalCoin, nil
}

func (c *Keeper) handlerERC20Token(ctx sdk.Context, evm *vm.EVM, sender, token common.Address, amount *big.Int) (sdk.Coin, error) {
	tokenPair, found := c.erc20Keeper.GetTokenPairByAddress(ctx, token)
	if !found {
		return sdk.Coin{}, fmt.Errorf("token pair not found: %s", token.String())
	}
	baseDenom := tokenPair.GetDenom()

	// transferFrom to erc20 module
	erc20Call := contract.NewERC20Call(evm, crosschaintypes.GetAddress(), token, 0)
	if err := erc20Call.TransferFrom(sender, c.erc20Keeper.ModuleAddress(), amount); err != nil {
		return sdk.Coin{}, err
	}
	if err := c.convertERC20(ctx, evm, tokenPair, sdk.NewCoin(baseDenom, sdkmath.NewIntFromBigInt(amount)), sender); err != nil {
		return sdk.Coin{}, err
	}
	return sdk.NewCoin(baseDenom, sdkmath.NewIntFromBigInt(amount)), nil
}

func (c *Keeper) convertERC20(
	ctx sdk.Context,
	evm *vm.EVM,
	tokenPair erc20types.TokenPair,
	amount sdk.Coin,
	sender common.Address,
) error {
	if tokenPair.IsNativeCoin() {
		erc20Call := contract.NewERC20Call(evm, c.erc20Keeper.ModuleAddress(), tokenPair.GetERC20Contract(), 0)
		err := erc20Call.Burn(c.erc20Keeper.ModuleAddress(), amount.Amount.BigInt())
		if err != nil {
			return err
		}
		if tokenPair.GetDenom() == fxtypes.DefaultDenom {
			err = c.bankKeeper.SendCoinsFromAccountToModule(ctx, tokenPair.GetERC20Contract().Bytes(), erc20types.ModuleName, sdk.NewCoins(amount))
			if err != nil {
				return err
			}
		}

	} else if tokenPair.IsNativeERC20() {
		if err := c.bankKeeper.MintCoins(ctx, erc20types.ModuleName, sdk.NewCoins(amount)); err != nil {
			return err
		}
	} else {
		return erc20types.ErrUndefinedOwner
	}

	if err := c.bankKeeper.SendCoinsFromModuleToAccount(ctx, erc20types.ModuleName, sender.Bytes(), sdk.NewCoins(amount)); err != nil {
		return err
	}
	return nil
}

// handlerCrossChain cross chain handler
// originToken is true represent cross chain denom(FX)
// when refund it, will not refund to evm token
// NOTE: fip20CrossChain only use for contract token, so origin token flag always false
func (c *Keeper) handlerCrossChain(
	ctx sdk.Context,
	from sdk.AccAddress,
	receipt string,
	amount, fee sdk.Coin,
	fxTarget fxtypes.FxTarget,
	memo string,
	originToken bool,
) error {
	total := sdk.NewCoin(amount.Denom, amount.Amount.Add(fee.Amount))
	// convert denom to target coin
	targetCoin, err := c.erc20Keeper.ConvertDenomToTarget(ctx, from.Bytes(), total, fxTarget)
	if err != nil && !erc20types.IsInsufficientLiquidityErr(err) {
		return fmt.Errorf("convert denom: %s", err.Error())
	}
	amount.Denom = targetCoin.Denom
	fee.Denom = targetCoin.Denom

	if fxTarget.IsIBC() {
		if err != nil {
			return fmt.Errorf("convert denom: %s", err.Error())
		}
		return c.ibcTransfer(ctx, from.Bytes(), receipt, amount, fee, fxTarget, memo, originToken)
	}

	return c.outgoingTransfer(ctx, from.Bytes(), receipt, amount, fee, fxTarget, originToken, err != nil)
}

func (c *Keeper) outgoingTransfer(
	ctx sdk.Context,
	from sdk.AccAddress,
	to string,
	amount, fee sdk.Coin,
	fxTarget fxtypes.FxTarget,
	originToken, insufficientLiquidit bool,
) error {
	if c.router == nil {
		return errors.New("cross chain router empty")
	}
	route, has := c.router.GetRoute(fxTarget.GetTarget())
	if !has {
		return errors.New("invalid target")
	}
	if err := route.TransferAfter(ctx, from, to, amount, fee, originToken, insufficientLiquidit); err != nil {
		return fmt.Errorf("cross chain error: %s", err.Error())
	}
	return nil
}

func (c *Keeper) ibcTransfer(
	ctx sdk.Context,
	from sdk.AccAddress,
	to string,
	amount, fee sdk.Coin,
	fxTarget fxtypes.FxTarget,
	memo string,
	originToken bool,
) error {
	if !fee.IsZero() {
		return fmt.Errorf("ibc transfer fee must be zero: %s", fee.String())
	}
	if strings.ToLower(fxTarget.Prefix) == contract.EthereumAddressPrefix {
		if err := contract.ValidateEthereumAddress(to); err != nil {
			return fmt.Errorf("invalid to address: %s", to)
		}
	} else {
		if _, err := sdk.GetFromBech32(to, fxTarget.Prefix); err != nil {
			return fmt.Errorf("invalid to address: %s", to)
		}
	}

	ibcTimeoutTimestamp := uint64(ctx.BlockTime().UnixNano()) + uint64(c.erc20Keeper.GetIbcTimeout(ctx))
	transferResponse, err := c.ibcTransferKeeper.Transfer(sdk.WrapSDKContext(ctx),
		transfertypes.NewMsgTransfer(
			fxTarget.SourcePort,
			fxTarget.SourceChannel,
			amount,
			from.String(),
			to,
			ibcclienttypes.ZeroHeight(),
			ibcTimeoutTimestamp,
			memo,
		),
	)
	if err != nil {
		return fmt.Errorf("ibc transfer error: %s", err.Error())
	}

	if !originToken {
		c.erc20Keeper.SetIBCTransferRelation(ctx, fxTarget.SourceChannel, transferResponse.GetSequence())
	}
	return nil
}
