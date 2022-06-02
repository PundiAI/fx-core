package keeper

import (
	"encoding/hex"
	"fmt"

	fxtypes "github.com/functionx/fx-core/types"

	"github.com/ethereum/go-ethereum/accounts/abi"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/functionx/fx-core/x/erc20/types"
)

// RegisterCoin deploys an erc20 contract and creates the token pair for the existing cosmos coin
func (k Keeper) RegisterCoin(ctx sdk.Context, coinMetadata banktypes.Metadata) (*types.TokenPair, error) {
	// check if the conversion is globally enabled
	params := k.GetParams(ctx)
	if !params.EnableErc20 {
		return nil, sdkerrors.Wrap(types.ErrERC20Disabled, "registration is currently disabled by governance")
	}

	//description use for name
	name := coinMetadata.Name
	//display use for symbol
	symbol := coinMetadata.Symbol
	//decimals
	decimals := uint8(0)
	for _, du := range coinMetadata.DenomUnits {
		if du.Denom == symbol {
			decimals = uint8(du.Exponent)
			break
		}
	}
	if coinMetadata.Base == fxtypes.DefaultDenom {
		decimals = fxtypes.BaseDenomUnit
	}
	if len(name) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidMetadata, "invalid name")
	}
	if len(symbol) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidMetadata, "invalid symbol")
	}
	if decimals == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidMetadata, "invalid symbol denom exponent")
	}

	// check if the denomination already registered
	if k.IsDenomRegistered(ctx, coinMetadata.Base) {
		return nil, sdkerrors.Wrapf(types.ErrTokenPairAlreadyExists, "coin denomination already registered: %s", coinMetadata.Description)
	}

	meta, isExist := k.bankKeeper.GetDenomMetaData(ctx, coinMetadata.Base)
	if isExist {
		if err := types.EqualMetadata(meta, coinMetadata); err != nil {
			return nil, sdkerrors.Wrap(types.ErrInvalidMetadata, err.Error())
		}
	} else {
		k.bankKeeper.SetDenomMetaData(ctx, coinMetadata)
	}

	addr, err := k.DeployUpgradableToken(ctx, types.ModuleAddress, name, symbol, decimals, coinMetadata.Base == fxtypes.DefaultDenom)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to create wrapped coin denom metadata for ERC20")
	}

	pair := types.NewTokenPair(addr, coinMetadata.Base, true, types.OWNER_MODULE)
	k.SetTokenPair(ctx, pair)
	k.SetDenomMap(ctx, pair.Denom, pair.GetID())
	k.SetERC20Map(ctx, common.HexToAddress(pair.Erc20Address), pair.GetID())

	return &pair, nil
}

// DeployERC20Contract creates and deploys an ERC20 contract on the EVM with the
// erc20 module account as owner.
func (k Keeper) DeployERC20Contract(ctx sdk.Context, coinMetadata banktypes.Metadata) (common.Address, error) {
	decimals := uint8(coinMetadata.DenomUnits[0].Exponent)
	erc20 := fxtypes.GetERC20(ctx.BlockHeight())
	ctorArgs, err := erc20.ABI.Pack(
		"",
		coinMetadata.Description,
		coinMetadata.Display,
		decimals,
	)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIPack, "coin metadata is invalid %s: %s", coinMetadata.Description, err.Error())
	}

	data := make([]byte, len(erc20.Bin)+len(ctorArgs))
	copy(data[:len(erc20.Bin)], erc20.Bin)
	copy(data[len(erc20.Bin):], ctorArgs)

	nonce, err := k.accountKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleAddress, nonce)
	_, err = k.CallEVMWithData(ctx, types.ModuleAddress, nil, data, true)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "failed to deploy contract for %s", coinMetadata.Description)
	}

	return contractAddr, nil
}

// RegisterERC20 creates a cosmos coin and registers the token pair between the coin and the ERC20
func (k Keeper) RegisterERC20(ctx sdk.Context, contract common.Address) (*types.TokenPair, error) {
	params := k.GetParams(ctx)
	if !params.EnableErc20 {
		return nil, sdkerrors.Wrap(types.ErrERC20Disabled, "registration is currently disabled by governance")
	}

	if k.IsERC20Registered(ctx, contract) {
		return nil, sdkerrors.Wrapf(types.ErrTokenPairAlreadyExists, "token ERC20 contract already registered: %s", contract.String())
	}

	_, baseDenom, _, err := k.CreateCoinMetadata(ctx, contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to create wrapped coin denom metadata for ERC20")
	}

	pair := types.NewTokenPair(contract, baseDenom, true, types.OWNER_EXTERNAL)
	k.SetTokenPair(ctx, pair)
	k.SetDenomMap(ctx, pair.Denom, pair.GetID())
	k.SetERC20Map(ctx, common.HexToAddress(pair.Erc20Address), pair.GetID())
	return &pair, nil
}

// CreateCoinMetadata generates the metadata to represent the ERC20 token on functionX.
func (k Keeper) CreateCoinMetadata(ctx sdk.Context, contract common.Address) (*banktypes.Metadata, string, string, error) {
	strContract := contract.String()

	erc20Data, err := k.QueryERC20(ctx, contract)
	if err != nil {
		return nil, "", "", err
	}

	_, isExist := k.bankKeeper.GetDenomMetaData(ctx, types.CreateDenom(strContract))
	if isExist {
		// metadata already exists; exit
		return nil, "", "", sdkerrors.Wrap(types.ErrInternalTokenPair, "denom metadata already registered")
	}

	if k.IsDenomRegistered(ctx, types.CreateDenom(strContract)) {
		return nil, "", "", sdkerrors.Wrapf(types.ErrInternalTokenPair, "coin denomination already registered: %s", erc20Data.Name)
	}

	// base denomination
	base := types.CreateDenom(strContract)

	// create a bank denom metadata based on the ERC20 token ABI details
	// metadata name is should always be the contract since it's the key
	// to the bank store
	metadata := banktypes.Metadata{
		Description: types.CreateDenomDescription(strContract),
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    base,
				Exponent: 0,
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
		return nil, "", "", sdkerrors.Wrapf(err, "FIP20 token data is invalid for contract %s", strContract)
	}

	k.bankKeeper.SetDenomMetaData(ctx, metadata)

	return &metadata, base, metadata.Symbol, nil
}

// ToggleRelay toggles relaying for a given token pair
func (k Keeper) ToggleRelay(ctx sdk.Context, token string) (types.TokenPair, error) {
	id := k.GetTokenPairID(ctx, token)
	if len(id) == 0 {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrTokenPairNotFound, "token '%s' not registered by id", token)
	}

	pair, found := k.GetTokenPair(ctx, id)
	if !found {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrTokenPairNotFound, "token '%s' not registered", token)
	}

	pair.Enabled = !pair.Enabled

	k.SetTokenPair(ctx, pair)
	return pair, nil
}

func (k Keeper) DeployUpgradableToken(ctx sdk.Context, from common.Address, name, symbol string, decimals uint8, origin bool) (common.Address, error) {
	tokenContract := fxtypes.GetERC20(ctx.BlockHeight())
	if origin {
		tokenContract, name, symbol = WrappedDenom(ctx.BlockHeight(), name, symbol)
	}
	k.Logger(ctx).Info("deploy token", "name", name, "symbol", symbol, "decimals", decimals, "origin", origin)
	//deploy proxy
	proxy, err := k.DeployERC1967Proxy(ctx, from, tokenContract.Address)
	if err != nil {
		return common.Address{}, err
	}
	err = k.InitializeUpgradable(ctx, from, proxy, tokenContract.ABI, name, symbol, decimals, types.ModuleAddress)
	return proxy, err
}

func (k Keeper) DeployERC1967Proxy(ctx sdk.Context, from, logicAddr common.Address, logicData ...byte) (common.Address, error) {
	k.Logger(ctx).Info("deploy erc1967 proxy", "logic", logicAddr.String(), "data", hex.EncodeToString(logicData))
	erc1967Proxy := fxtypes.GetERC1967Proxy(ctx.BlockHeight())

	if len(logicData) == 0 {
		logicData = []byte{}
	}
	return k.DeployContract(ctx, from, erc1967Proxy.ABI, erc1967Proxy.Bin, logicAddr, logicData)
}

func (k Keeper) InitializeUpgradable(ctx sdk.Context, from, contract common.Address, abi abi.ABI, data ...interface{}) error {
	k.Logger(ctx).Info("initialize upgradable", "contract", contract.Hex())
	_, err := k.CallEVM(ctx, abi, from, contract, true, "initialize", data...)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to initialize contract")
	}
	return nil
}

func (k Keeper) DeployContract(ctx sdk.Context, from common.Address, abi abi.ABI, bin []byte, constructorData ...interface{}) (common.Address, error) {
	ctorArgs, err := abi.Pack("", constructorData...)
	if err != nil {
		return common.Address{}, sdkerrors.Wrap(err, "pack constructor data")
	}
	data := make([]byte, len(bin)+len(ctorArgs))
	copy(data[:len(bin)], bin)
	copy(data[len(bin):], ctorArgs)

	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(from, nonce)
	_, err = k.CallEVMWithData(ctx, from, nil, data, true)
	if err != nil {
		return common.Address{}, sdkerrors.Wrap(err, "failed to deploy contract")
	}
	return contractAddr, nil
}

func WrappedDenom(height int64, name, symbol string) (fxtypes.Contract, string, string) {
	contract := fxtypes.GetWFX(height)
	wrappedName := fmt.Sprintf("Wrapped %s", name)
	wrappedSymbol := fmt.Sprintf("W%s", symbol)

	return contract, wrappedName, wrappedSymbol
}
