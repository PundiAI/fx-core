package keeper

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/erc20/types"
)

func (k Keeper) RegisterNativeCoin(ctx context.Context, name, symbol string, decimals uint8) (types.ERC20Token, error) {
	if err := k.CheckEnableErc20(ctx); err != nil {
		return types.ERC20Token{}, err
	}

	var tokenContract contract.Contract
	var metadata banktypes.Metadata
	if symbol == fxtypes.DefaultSymbol {
		metadata = fxtypes.NewDefaultMetadata()
		tokenContract = contract.GetWPUNDIAI()
		name = fmt.Sprintf("Wrapped %s", name)
		symbol = fmt.Sprintf("W%s", symbol)
	} else {
		metadata = fxtypes.NewMetadata(name, symbol, uint32(decimals))
		tokenContract = contract.GetERC20()
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	erc20TokenAddr, err := k.evmKeeper.DeployUpgradableContract(sdkCtx, k.contractOwner,
		tokenContract.Address, nil, &tokenContract.ABI, name, symbol, decimals, k.contractOwner)
	if err != nil {
		return types.ERC20Token{}, err
	}

	erc20Token, err := k.AddERC20Token(ctx, metadata, erc20TokenAddr, types.OWNER_MODULE)
	if err != nil {
		return types.ERC20Token{}, err
	}

	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRegisterCoin,
		sdk.NewAttribute(types.AttributeKeyDenom, erc20Token.Denom),
		sdk.NewAttribute(types.AttributeKeyTokenAddress, erc20Token.Erc20Address),
	))
	return erc20Token, nil
}

func (k Keeper) RegisterNativeERC20(ctx context.Context, erc20Addr common.Address) (types.ERC20Token, error) {
	if err := k.CheckEnableErc20(ctx); err != nil {
		return types.ERC20Token{}, err
	}

	name, symbol, decimals, err := k.ERC20BaseInfo(ctx, erc20Addr)
	if err != nil {
		return types.ERC20Token{}, err
	}

	metadata := fxtypes.NewMetadata(name, symbol, uint32(decimals))
	erc20Token, err := k.AddERC20Token(ctx, metadata, erc20Addr, types.OWNER_EXTERNAL)
	if err != nil {
		return types.ERC20Token{}, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRegisterERC20,
		sdk.NewAttribute(types.AttributeKeyDenom, erc20Token.Denom),
		sdk.NewAttribute(types.AttributeKeyTokenAddress, erc20Token.Erc20Address),
	))

	return erc20Token, nil
}

func (k Keeper) RegisterBridgeToken(ctx context.Context, baseDenom, channel, ibcDenom, chainName, contractAddr string, isNative bool) (types.ERC20Token, error) {
	erc20Token, err := k.ERC20Token.Get(ctx, baseDenom)
	if err != nil {
		return types.ERC20Token{}, err
	}

	// add ibc token
	isIBCDenom := strings.HasPrefix(ibcDenom, ibctransfertypes.DenomPrefix+"/")
	if isIBCDenom {
		return erc20Token, k.AddIBCToken(ctx, baseDenom, channel, ibcDenom)
	}

	// add bridge token
	if !fxtypes.IsSupportChain(chainName) {
		return types.ERC20Token{}, sdkerrors.ErrKeyNotFound.Wrapf("chain name %s not found", chainName)
	}
	return erc20Token, k.AddBridgeToken(ctx, baseDenom, chainName, contractAddr, isNative)
}
