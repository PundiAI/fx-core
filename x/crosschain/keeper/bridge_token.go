package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v8/types"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
)

func (k Keeper) AddBridgeTokenExecuted(ctx sdk.Context, claim *types.MsgBridgeTokenClaim) error {
	k.Logger(ctx).Info("add bridge token claim", "symbol", claim.Symbol, "token", claim.TokenContract)

	// Check if it already exists
	bridgeDenom := types.NewBridgeDenom(k.moduleName, claim.TokenContract)
	has, err := k.erc20Keeper.HasToken(ctx, bridgeDenom)
	if err != nil {
		return err
	}
	if has {
		return types.ErrInvalid.Wrapf("bridge token is exist %s", bridgeDenom)
	}

	if claim.Symbol == fxtypes.DefaultDenom && uint64(fxtypes.DenomUnit) != claim.Decimals {
		return types.ErrInvalid.Wrapf("%s denom decimals not match %d, expect %d",
			fxtypes.DefaultDenom, claim.Decimals, fxtypes.DenomUnit)
	}
	bridgeToken, err := k.erc20Keeper.GetBridgeToken(ctx, claim.Symbol, k.moduleName)
	if err != nil {
		return err
	}
	if bridgeToken.Contract != claim.TokenContract {
		return types.ErrInvalid.Wrapf("bridge token contract not match %s, expect %s", bridgeToken.Contract, claim.TokenContract)
	}
	return nil
}

func (k Keeper) BridgeCoinSupply(ctx context.Context, token, target string) (sdk.Coin, error) {
	baseDenom, err := k.erc20Keeper.GetBaseDenom(ctx, token)
	if err != nil {
		return sdk.Coin{}, err
	}
	fxTarget, err := types.ParseFxTarget(target)
	if err != nil {
		return sdk.Coin{}, err
	}
	var targetDenom string
	if fxTarget.IsIBC() {
		ibcToken, err := k.erc20Keeper.GetIBCToken(ctx, baseDenom, fxTarget.IBCChannel)
		if err != nil {
			return sdk.Coin{}, err
		}
		targetDenom = ibcToken.IbcDenom
	} else {
		bridgeToken, err := k.erc20Keeper.GetBridgeToken(ctx, baseDenom, fxTarget.GetModuleName())
		if err != nil {
			return sdk.Coin{}, err
		}
		targetDenom = bridgeToken.BridgeDenom()
	}

	supply := k.bankKeeper.GetSupply(ctx, targetDenom)
	return supply, nil
}

func (k Keeper) GetBaseDenomByErc20(ctx sdk.Context, erc20Addr common.Address) (erc20types.ERC20Token, error) {
	baseDenom, err := k.erc20Keeper.GetBaseDenom(ctx, erc20Addr.String())
	if err != nil {
		return erc20types.ERC20Token{}, err
	}
	return k.erc20Keeper.GetERC20Token(ctx, baseDenom)
}
