package keeper

import (
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/erc20/types"
)

// RegisterNativeCoin deploys an erc20 contract and creates the token pair for the existing cosmos coin
func (k Keeper) RegisterNativeCoin(ctx sdk.Context, coinMetadata banktypes.Metadata) (*types.TokenPair, error) {
	// check if the conversion is globally enabled
	if !k.GetEnableErc20(ctx) {
		return nil, errorsmod.Wrap(types.ErrERC20Disabled, "registration is currently disabled by governance")
	}

	decimals := getErc20Decimals(coinMetadata)

	// check if the denomination already registered
	if k.IsDenomRegistered(ctx, coinMetadata.Base) {
		return nil, errorsmod.Wrapf(types.ErrTokenPairAlreadyExists, "coin denomination already registered: %s", coinMetadata.Base)
	}

	// base not register as alias
	if k.IsAliasDenomRegistered(ctx, coinMetadata.Base) {
		return nil, errorsmod.Wrapf(types.ErrInvalidMetadata, "alias %s already registered", coinMetadata.Base)
	}

	if len(coinMetadata.DenomUnits) > 0 && len(coinMetadata.DenomUnits[0].Aliases) > 0 {
		for _, alias := range coinMetadata.DenomUnits[0].Aliases {
			if alias == coinMetadata.Base || alias == coinMetadata.Display || alias == coinMetadata.Symbol {
				return nil, errorsmod.Wrap(types.ErrInvalidMetadata, "alias can not equal base, display or symbol")
			}
			// alias not register as base
			if k.IsDenomRegistered(ctx, alias) {
				return nil, errorsmod.Wrapf(types.ErrInvalidMetadata, "denom %s already registered", alias)
			}
			// alias must not register
			if k.IsAliasDenomRegistered(ctx, alias) {
				return nil, errorsmod.Wrapf(types.ErrInvalidMetadata, "alias %s already registered", alias)
			}
		}
		k.SetAliasesDenom(ctx, coinMetadata.Base, coinMetadata.DenomUnits[0].Aliases...)
	}

	meta, isExist := k.bankKeeper.GetDenomMetaData(ctx, coinMetadata.Base)
	if isExist {
		if err := types.EqualMetadata(meta, coinMetadata); err != nil {
			return nil, errorsmod.Wrap(types.ErrInvalidMetadata, err.Error())
		}
	} else {
		k.bankKeeper.SetDenomMetaData(ctx, coinMetadata)
	}

	addr, err := k.DeployUpgradableToken(ctx, k.moduleAddress, coinMetadata.Name, coinMetadata.Symbol, decimals)
	if err != nil {
		return nil, err
	}

	pair := types.NewTokenPair(addr, coinMetadata.Base, true, types.OWNER_MODULE)
	k.AddTokenPair(ctx, pair)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRegisterCoin,
		sdk.NewAttribute(types.AttributeKeyDenom, pair.Denom),
		sdk.NewAttribute(types.AttributeKeyTokenAddress, pair.Erc20Address),
	))
	return &pair, nil
}

// RegisterNativeERC20 creates a cosmos coin and registers the token pair between the coin and the ERC20
//
//gocyclo:ignore
func (k Keeper) RegisterNativeERC20(ctx sdk.Context, contract common.Address, aliases ...string) (*types.TokenPair, error) {
	if !k.GetEnableErc20(ctx) {
		return nil, errorsmod.Wrap(types.ErrERC20Disabled, "registration is currently disabled by governance")
	}

	if k.IsERC20Registered(ctx, contract) {
		return nil, errorsmod.Wrapf(types.ErrTokenPairAlreadyExists, "token ERC20 contract already registered: %s", contract.String())
	}

	erc20Data, err := k.QueryERC20(ctx, contract)
	if err != nil {
		return nil, err
	}

	// base denomination
	base := strings.ToLower(erc20Data.Symbol)
	if erc20Data.Symbol == fxtypes.DefaultDenom || k.IsDenomRegistered(ctx, base) {
		return nil, errorsmod.Wrapf(types.ErrInternalTokenPair, "coin denomination already registered: %s", erc20Data.Name)
	}

	// base not register as alias
	if k.IsAliasDenomRegistered(ctx, base) {
		return nil, errorsmod.Wrapf(types.ErrInternalTokenPair, "alias %s already registered", base)
	}

	if len(aliases) > 0 {
		for _, alias := range aliases {
			if alias == base || alias == erc20Data.Symbol {
				return nil, errorsmod.Wrap(types.ErrInvalidAlias, "alias can not equal base, display or symbol")
			}
			// alias not register as base
			if k.IsDenomRegistered(ctx, alias) {
				return nil, errorsmod.Wrapf(types.ErrInvalidAlias, "denom %s already registered", alias)
			}
			// alias must not register
			if k.IsAliasDenomRegistered(ctx, alias) {
				return nil, errorsmod.Wrapf(types.ErrInvalidAlias, "alias %s already registered", alias)
			}
		}
		k.SetAliasesDenom(ctx, base, aliases...)
	}

	_, isExist := k.bankKeeper.GetDenomMetaData(ctx, base) // TODO if register must be equal
	if isExist {
		// metadata already exists; exit
		return nil, errorsmod.Wrap(types.ErrInternalTokenPair, "denom metadata already registered")
	}

	// create a bank denom metadata based on the ERC20 token ABI details
	// metadata name is should always be the contract since it's the key
	// to the bank store
	metadata := banktypes.Metadata{
		Description: types.CreateDenomDescription(contract.String()),
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    base,
				Exponent: 0,
				Aliases:  aliases,
			},
		},
		Base:    base,
		Display: base,
		Name:    erc20Data.Name,
		Symbol:  erc20Data.Symbol,
	}

	// only append metadata if decimals > 0, otherwise validation fails
	if erc20Data.Decimals > 0 {
		metadata.DenomUnits = append(
			metadata.DenomUnits,
			&banktypes.DenomUnit{
				Denom:    erc20Data.Symbol,
				Exponent: uint32(erc20Data.Decimals),
			},
		)
	}

	if err := metadata.Validate(); err != nil {
		return nil, errorsmod.Wrapf(err, "FIP20 token data is invalid for contract %s", contract.String())
	}
	k.bankKeeper.SetDenomMetaData(ctx, metadata)

	pair := types.NewTokenPair(contract, metadata.Base, true, types.OWNER_EXTERNAL)
	k.AddTokenPair(ctx, pair)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRegisterERC20,
		sdk.NewAttribute(types.AttributeKeyDenom, pair.Denom),
		sdk.NewAttribute(types.AttributeKeyTokenAddress, pair.Erc20Address),
	))

	return &pair, nil
}

// ToggleTokenConvert toggles relaying for a given token pair
func (k Keeper) ToggleTokenConvert(ctx sdk.Context, token string) (types.TokenPair, error) {
	pair, found := k.GetTokenPair(ctx, token)
	if !found {
		return types.TokenPair{}, errorsmod.Wrapf(types.ErrTokenPairNotFound, "token '%s' not registered", token)
	}
	pair.Enabled = !pair.Enabled

	k.SetTokenPair(ctx, pair)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeToggleTokenRelay,
		sdk.NewAttribute(types.AttributeKeyDenom, pair.Denom),
		sdk.NewAttribute(types.AttributeKeyTokenAddress, pair.Erc20Address),
	))

	return pair, nil
}

// UpdateDenomAliases update denom alias
// if alias not registered, add to denom alias
// if alias registered with denom, remove from denom alias
// if alias registered, but not with denom, return error
func (k Keeper) UpdateDenomAliases(ctx sdk.Context, denom, alias string) (bool, error) {
	// check if the denom denomination already registered
	if !k.IsDenomRegistered(ctx, denom) {
		return false, errorsmod.Wrapf(types.ErrInvalidDenom, "coin denomination not registered: %s", denom)
	}
	// check if the alias not registered
	if k.IsDenomRegistered(ctx, alias) {
		return false, errorsmod.Wrapf(types.ErrInvalidDenom, "coin denomination already registered: %s", alias)
	}

	md, found := k.GetValidMetadata(ctx, denom)
	if !found {
		return false, errorsmod.Wrapf(types.ErrInvalidMetadata, "denom %s not support update denom aliases", denom)
	}

	oldAliases := md.DenomUnits[0].Aliases
	newAliases := make([]string, 0, len(oldAliases)+1)

	registeredDenom, found := k.GetAliasDenom(ctx, alias)
	// check if the alias not register denom-alias
	if !found {
		newAliases = append(oldAliases, alias)
		k.SetAliasesDenom(ctx, denom, alias)
	} else if registeredDenom == denom {
		// check if the denom equal alias registered denom
		for _, denomAlias := range oldAliases {
			if denomAlias == alias {
				continue
			}
			newAliases = append(newAliases, denomAlias)
		}
		// NOTE: FX,PUNDIX,PURSE can delete all alias, others must keep at least one
		k.DeleteAliasesDenom(ctx, alias)
	} else {
		// check if denom not equal alias registered denom, return error
		return false, errorsmod.Wrapf(types.ErrInvalidDenom,
			"alias %s already registered, but denom expected: %s, actual: %s",
			alias, registeredDenom, denom)
	}

	md.DenomUnits[0].Aliases = newAliases
	k.bankKeeper.SetDenomMetaData(ctx, md)

	addFlag := len(newAliases) > len(oldAliases)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeToggleTokenRelay,
		sdk.NewAttribute(types.AttributeKeyDenom, denom),
		sdk.NewAttribute(types.AttributeKeyAlias, alias),
		sdk.NewAttribute(types.AttributeKeyUpdateFlag, strconv.FormatBool(addFlag)),
	))

	return addFlag, nil
}

func getErc20Decimals(md banktypes.Metadata) (decimals uint8) {
	decimals = uint8(0)
	for _, du := range md.DenomUnits {
		if du.Denom == md.Symbol {
			decimals = uint8(du.Exponent)
			break
		}
	}
	if md.Base == fxtypes.DefaultDenom {
		decimals = fxtypes.DenomUnit
	}
	return decimals
}
