package keeper

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/functionx/fx-core/contracts"
	"github.com/functionx/fx-core/x/erc20/types"
)

// RegisterCoin deploys an erc20 contract and creates the token pair for the existing cosmos coin
func (k Keeper) RegisterCoin(ctx sdk.Context, coinMetadata banktypes.Metadata) (*types.TokenPair, error) {
	// check if the conversion is globally enabled
	params := k.GetParams(ctx)
	if !params.EnableErc20 {
		return nil, sdkerrors.Wrap(types.ErrERC20Disabled, "registration is currently disabled by governance")
	}

	// prohibit denominations that contain the evm denom
	if strings.Contains(coinMetadata.Base, "evm") {
		return nil, sdkerrors.Wrapf(types.ErrEVMDenom, "cannot register the EVM denomination %s", coinMetadata.Base)
	}

	//description use for name
	name := coinMetadata.Description
	//display use for symbol
	symbol := coinMetadata.Display
	//decimals
	decimals := uint8(0)
	for _, du := range coinMetadata.DenomUnits {
		if du.Denom == symbol {
			decimals = uint8(du.Exponent)
			break
		}
	}
	if decimals == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidMetadata, "invalid display denom exponent")
	}

	// check if the denomination already registered
	if k.IsDenomRegistered(ctx, coinMetadata.Description) {
		return nil, sdkerrors.Wrapf(types.ErrTokenPairAlreadyExists, "coin denomination already registered: %s", coinMetadata.Description)
	}

	// check if the coin exists by ensuring the supply is set
	//if !k.bankKeeper.HasSupply(ctx, coinMetadata.Base) {
	//	return nil, sdkerrors.Wrapf(
	//		sdkerrors.ErrInvalidCoins,
	//		"base denomination '%s' cannot have a supply of 0", coinMetadata.Base,
	//	)
	//}

	if err := k.verifyMetadata(ctx, coinMetadata); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInternalTokenPair, "coin metadata is invalid %s", coinMetadata.Description)
	}

	evmParams := k.evmKeeper.GetParams(ctx)
	addr, err := k.DeployTokenUpgrade(ctx, types.ModuleAddress, name, symbol, decimals, coinMetadata.Base == evmParams.EvmDenom)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to create wrapped coin denom metadata for ERC20")
	}

	pair := types.NewTokenPair(addr, coinMetadata.Base, true, types.OWNER_MODULE)
	k.SetTokenPair(ctx, pair)
	k.SetDenomMap(ctx, pair.Denom, pair.GetID())
	k.SetERC20Map(ctx, common.HexToAddress(pair.Erc20Address), pair.GetID())

	return &pair, nil
}

// verify if the metadata matches the existing one, if not it sets it to the store
func (k Keeper) verifyMetadata(ctx sdk.Context, coinMetadata banktypes.Metadata) error {
	meta := k.bankKeeper.GetDenomMetaData(ctx, coinMetadata.Base)
	if len(meta.Base) == 0 { //new coin metadata
		k.bankKeeper.SetDenomMetaData(ctx, coinMetadata)
		return nil
	}
	// If it already existed, Check that is equal to what is stored
	return types.EqualMetadata(meta, coinMetadata)
}

// DeployERC20Contract creates and deploys an ERC20 contract on the EVM with the
// erc20 module account as owner.
func (k Keeper) DeployERC20Contract(ctx sdk.Context, coinMetadata banktypes.Metadata) (common.Address, error) {
	decimals := uint8(coinMetadata.DenomUnits[0].Exponent)
	erc20 := contracts.GetERC20Config(ctx.BlockHeight())
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
	_, err = k.CallEVMWithData(ctx, types.ModuleAddress, nil, data)
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

	metadata, err := k.CreateCoinMetadata(ctx, contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to create wrapped coin denom metadata for ERC20")
	}

	pair := types.NewTokenPair(contract, metadata.Description, true, types.OWNER_EXTERNAL)
	k.SetTokenPair(ctx, pair)
	k.SetDenomMap(ctx, pair.Denom, pair.GetID())
	k.SetERC20Map(ctx, common.HexToAddress(pair.Erc20Address), pair.GetID())
	return &pair, nil
}

// CreateCoinMetadata generates the metadata to represent the ERC20 token on evmos.
func (k Keeper) CreateCoinMetadata(ctx sdk.Context, contract common.Address) (*banktypes.Metadata, error) {
	strContract := contract.String()

	erc20Data, err := k.QueryERC20(ctx, contract)
	if err != nil {
		return nil, err
	}

	meta := k.bankKeeper.GetDenomMetaData(ctx, types.CreateDenom(strContract))
	if len(meta.Base) > 0 {
		// metadata already exists; exit
		return nil, sdkerrors.Wrap(types.ErrInternalTokenPair, "denom metadata already registered")
	}

	if k.IsDenomRegistered(ctx, types.CreateDenom(strContract)) {
		return nil, sdkerrors.Wrapf(types.ErrInternalTokenPair, "coin denomination already registered: %s", erc20Data.Name)
	}

	// base denomination
	base := types.CreateDenom(strContract)

	// create a bank denom metadata based on the ERC20 token ABI details
	// metadata name is should always be the contract since it's the key
	// to the bank store
	metadata := banktypes.Metadata{
		Description: erc20Data.Name,
		Base:        base,
		// NOTE: Denom units MUST be increasing
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    base,
				Exponent: 0,
			},
		},
		Display: erc20Data.Symbol,
	}

	// only append metadata if decimals > 0, otherwise validation fails
	if erc20Data.Decimals > 0 {
		nameSanitized := types.SanitizeERC20Name(erc20Data.Name)
		metadata.DenomUnits = append(
			metadata.DenomUnits,
			&banktypes.DenomUnit{
				Denom:    nameSanitized,
				Exponent: uint32(erc20Data.Decimals),
			},
		)
		metadata.Display = nameSanitized
	}

	if err := metadata.Validate(); err != nil {
		return nil, sdkerrors.Wrapf(err, "ERC20 token data is invalid for contract %s", strContract)
	}

	k.bankKeeper.SetDenomMetaData(ctx, metadata)

	return &metadata, nil
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

// UpdateTokenPairERC20 updates the ERC20 token address for the registered token pair
func (k Keeper) UpdateTokenPairERC20(ctx sdk.Context, erc20Addr, newERC20Addr common.Address) (types.TokenPair, error) {
	id := k.GetERC20Map(ctx, erc20Addr)
	if len(id) == 0 {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrInternalTokenPair, "token %s not registered", erc20Addr)
	}

	pair, found := k.GetTokenPair(ctx, id)
	if !found {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrTokenPairNotFound, "token '%s' not registered", erc20Addr)
	}

	// Get current stored metadata
	metadata := k.bankKeeper.GetDenomMetaData(ctx, pair.Denom)
	if len(metadata.Base) <= 0 {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrInternalTokenPair, "could not get metadata for %s", pair.Denom)
	}

	// safety check
	if len(metadata.DenomUnits) == 0 {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrInternalTokenPair, "metadata denom units for %s cannot be empty", pair.Erc20Address)
	}

	// Get new erc20 values
	erc20Data, err := k.QueryERC20(ctx, newERC20Addr)
	if err != nil {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrInternalTokenPair, "could not get token %s erc20Data", newERC20Addr.String())
	}

	// compare metadata and ERC20 details
	if metadata.Display != erc20Data.Symbol || metadata.Description != erc20Data.Name || metadata.Base != types.CreateDenom(erc20Addr.String()) {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrInternalTokenPair, "metadata details (display, symbol, description) don't match the ERC20 details from %s ", pair.Erc20Address)
	}

	// check that the denom units contain one item with the same
	// name and decimal values as the ERC20
	found = false
	for _, denomUnit := range metadata.DenomUnits {
		// iterate denom units until we found the one with the ERC20 Name
		if denomUnit.Denom != erc20Data.Name {
			continue
		}

		// once found, check it has the same exponent
		if denomUnit.Exponent != uint32(erc20Data.Decimals) {
			return types.TokenPair{}, sdkerrors.Wrapf(
				types.ErrInternalTokenPair, "metadata denom unit exponent doesn't match the ERC20 details from %s, expected %d, got %d",
				pair.Erc20Address, erc20Data.Decimals, denomUnit.Exponent,
			)
		}

		// break as metadata might contain denom units for higher exponents
		found = true
		break
	}

	if !found {
		return types.TokenPair{}, sdkerrors.Wrapf(
			types.ErrInternalTokenPair,
			"metadata doesn't contain denom unit found for ERC20 %s (%s)",
			erc20Data.Name, pair.Erc20Address,
		)
	}

	// Update the metadata description with the new address
	metadata.Base = types.CreateDenom(newERC20Addr.String())
	k.bankKeeper.SetDenomMetaData(ctx, metadata)
	// Delete old token pair (id is changed because the ERC20 address was modifed)
	k.DeleteTokenPair(ctx, pair)
	// Update the address
	pair.Erc20Address = newERC20Addr.Hex()
	newID := pair.GetID()
	// Set the new pair
	k.SetTokenPair(ctx, pair)
	// Overwrite the value because id was changed
	k.SetDenomMap(ctx, pair.Denom, newID)
	// Add the new address
	k.SetERC20Map(ctx, newERC20Addr, newID)
	return pair, nil
}

func (k Keeper) HandleInitEvmProposal(ctx sdk.Context, p *types.InitEvmProposal) error {
	//init fee market
	k.Logger(ctx).Info("init fee market", "params", p.FeemarketParams.String())
	if p.FeemarketParams.BaseFee.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "base fee cannot be negative")
	}
	// set feeMarket baseFee
	k.feeMarketKeeper.SetBaseFee(ctx, p.FeemarketParams.BaseFee.BigInt())
	// set feeMarket blockGasUsed
	k.feeMarketKeeper.SetBlockGasUsed(ctx, 0)
	// init feeMarket module params
	k.feeMarketKeeper.SetParams(ctx, *p.FeemarketParams)

	//init evm
	k.Logger(ctx).Info("init evm", "params", p.EvmParams.String())
	k.evmKeeper.SetParams(ctx, *p.EvmParams)

	//init erc20
	k.Logger(ctx).Info("init erc20", "params", p.Erc20Params.String())

	//if err := k.ModuleInit(ctx, p.Erc20Params.EnableErc20,
	//	p.Erc20Params.EnableEVMHook, p.Erc20Params.IbcTransferTimeoutHeight); err != nil {
	//	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	//}

	////init register coin
	//events := make([]sdk.Event, 0, len(p.Metadata))
	//for _, metadata := range p.Metadata {
	//	k.Logger(ctx).Info("register coin", "coin", metadata.String())
	//	pair, err := k.erc20Keeper.RegisterCoin(ctx, metadata)
	//	if err != nil {
	//		return sdkerrors.Wrapf(erc20types.ErrInvalidMetadata, fmt.Sprintf("base %s, display %s, error %s",
	//			metadata.Base, metadata.Display, err.Error()))
	//	}
	//	event := sdk.NewEvent(
	//		erc20types.EventTypeRegisterCoin,
	//		sdk.NewAttribute(erc20types.AttributeKeyCosmosCoin, pair.Denom),
	//		sdk.NewAttribute(erc20types.AttributeKeyFIP20Token, pair.Fip20Address),
	//	)
	//	events = append(events, event)
	//}
	//ctx.EventManager().EmitEvents(events)
	return nil
}

func (k Keeper) DeployTokenUpgrade(ctx sdk.Context, from common.Address, name, symbol string, decimals uint8, origin bool) (common.Address, error) {
	k.Logger(ctx).Info("deploy token upgrade", "name", name, "symbol", symbol, "decimals", decimals)

	logicConfig := contracts.GetERC20Config(ctx.BlockHeight())
	logicAddr := common.HexToAddress(contracts.FIP20UpgradeCodeAddress)
	if origin {
		logicConfig = contracts.GetWFXConfig(ctx.BlockHeight())
		logicAddr = common.HexToAddress(contracts.WFXUpgradeCodeAddress)
	}

	//deploy proxy
	proxy, err := k.DeployERC1967Proxy(ctx, from, logicAddr)
	if err != nil {
		return common.Address{}, err
	}
	err = k.InitializeUpgradable(ctx, from, proxy, logicConfig.ABI, name, symbol, decimals, types.ModuleAddress)
	return proxy, err
}

func (k Keeper) DeployERC1967Proxy(ctx sdk.Context, from, logicAddr common.Address, logicData ...byte) (common.Address, error) {
	k.Logger(ctx).Info("deploy erc1967 proxy", "logic", logicAddr.String(), "data", hex.EncodeToString(logicData))
	erc1967ProxyConfig := contracts.GetERC1967ProxyConfig(ctx.BlockHeight())

	if len(logicData) == 0 {
		logicData = []byte{}
	}
	return k.DeployContract(ctx, from, erc1967ProxyConfig.ABI, erc1967ProxyConfig.Bin, logicAddr, logicData)
}

func (k Keeper) InitializeUpgradable(ctx sdk.Context, from, contract common.Address, abi abi.ABI, data ...interface{}) error {
	k.Logger(ctx).Info("initialize upgradable", "contract", contract.Hex())
	_, err := k.CallEVM(ctx, abi, from, contract, "initialize", data...)
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
	_, err = k.CallEVMWithData(ctx, from, nil, data)
	if err != nil {
		return common.Address{}, sdkerrors.Wrap(err, "failed to deploy contract")
	}
	return contractAddr, nil
}
