package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (k Keeper) AddBridgeTokenExecuted(ctx sdk.Context, claim *types.MsgBridgeTokenClaim) error {
	k.Logger(ctx).Info("add bridge token claim", "symbol", claim.Symbol, "token", claim.TokenContract)

	baseDenom := strings.ToLower(claim.Symbol)
	if claim.Symbol == fxtypes.DefaultSymbol {
		if uint64(fxtypes.DenomUnit) != claim.Decimals {
			return types.ErrInvalid.Wrapf("%s denom decimals not match %d, expect %d",
				fxtypes.DefaultSymbol, claim.Decimals, fxtypes.DenomUnit)
		}
		baseDenom = fxtypes.DefaultDenom
	}

	err := k.erc20Keeper.AddBridgeToken(ctx, baseDenom, k.moduleName, claim.TokenContract, false)
	if err != nil {
		return err
	}

	hasToken, err := k.erc20Keeper.HasERC20Token(ctx, baseDenom)
	if err != nil {
		return err
	}
	if hasToken {
		return nil
	}
	_, err = k.erc20Keeper.RegisterNativeCoin(ctx, claim.Name, claim.Symbol, uint8(claim.Decimals))
	return err
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
		bridgeToken, err := k.erc20Keeper.GetBridgeToken(ctx, fxTarget.GetModuleName(), baseDenom)
		if err != nil {
			return sdk.Coin{}, err
		}
		targetDenom = bridgeToken.BridgeDenom()
	}

	supply := k.bankKeeper.GetSupply(ctx, targetDenom)
	return supply, nil
}

func (k Keeper) GetERC20TokenByAddr(ctx sdk.Context, erc20Addr common.Address) (erc20types.ERC20Token, error) {
	baseDenom, err := k.erc20Keeper.GetBaseDenom(ctx, erc20Addr.String())
	if err != nil {
		return erc20types.ERC20Token{}, err
	}
	return k.erc20Keeper.GetERC20Token(ctx, baseDenom)
}
