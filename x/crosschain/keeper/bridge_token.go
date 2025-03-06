package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (k Keeper) AddBridgeTokenExecuted(ctx sdk.Context, claim *types.MsgBridgeTokenClaim) error {
	k.Logger(ctx).Info("add bridge token claim", "symbol", claim.Symbol, "token", claim.TokenContract)

	baseDenom := claim.GetBaseDenom()
	metaData, found := k.bankKeeper.GetDenomMetaData(ctx, baseDenom)
	if found {
		foundUnit := false
		for _, unit := range metaData.DenomUnits {
			if unit.Exponent == uint32(claim.Decimals) {
				foundUnit = true
			}
		}
		if !foundUnit {
			return types.ErrInvalid.Wrapf("%s denom decimals not match %d", claim.Symbol, claim.Decimals)
		}
	}

	erc20Token, err := k.erc20Keeper.GetERC20Token(ctx, baseDenom)
	if err != nil {
		if !errors.IsOf(err, collections.ErrNotFound) {
			return err
		}
		erc20Token, err = k.erc20Keeper.RegisterNativeCoin(ctx, claim.Name, claim.Symbol, uint8(claim.Decimals))
		if err != nil {
			return err
		}
	}
	_, err = k.erc20Keeper.GetBridgeToken(ctx, k.moduleName, baseDenom)
	if err != nil {
		if !errors.IsOf(err, collections.ErrNotFound) {
			return err
		}
		return k.erc20Keeper.AddBridgeToken(ctx, baseDenom, k.moduleName, claim.TokenContract, erc20Token.IsNativeERC20())
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
