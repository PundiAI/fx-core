package keeper

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
	feemarkettypes "github.com/functionx/fx-core/x/feemarket/types"
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
	if k.IsDenomRegistered(ctx, coinMetadata.Base) {
		return nil, sdkerrors.Wrapf(types.ErrTokenPairAlreadyExists, "coin denomination already registered: %s", coinMetadata.Description)
	}

	////check if the coin exists by ensuring the supply is set
	//if !k.bankKeeper.HasSupply(ctx, coinMetadata.Base) {
	//	return nil, sdkerrors.Wrapf(
	//		sdkerrors.ErrInvalidCoins,
	//		"base denomination '%s' cannot have a supply of 0", coinMetadata.Base,
	//	)
	//}

	meta := k.bankKeeper.GetDenomMetaData(ctx, coinMetadata.Base)
	if len(meta.Base) > 0 {
		if err := types.EqualMetadata(meta, coinMetadata); err != nil {
			return nil, sdkerrors.Wrap(types.ErrInvalidMetadata, err.Error())
		}
	} else {
		k.bankKeeper.SetDenomMetaData(ctx, coinMetadata)
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
	erc20 := contracts.GetERC20(ctx.BlockHeight())
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

	meta := k.bankKeeper.GetDenomMetaData(ctx, types.CreateDenom(strContract))
	if len(meta.Base) > 0 {
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
		metadata.DenomUnits = append(
			metadata.DenomUnits,
			&banktypes.DenomUnit{
				Denom:    erc20Data.Symbol,
				Exponent: uint32(erc20Data.Decimals),
			},
		)
		metadata.Display = erc20Data.Symbol
	}
	symbol := erc20Data.Symbol

	if err := metadata.Validate(); err != nil {
		return nil, "", "", sdkerrors.Wrapf(err, "FIP20 token data is invalid for contract %s", strContract)
	}

	k.bankKeeper.SetDenomMetaData(ctx, metadata)

	return &metadata, base, symbol, nil
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

func (k Keeper) HandleInitEvmProposal(ctx sdk.Context, erc20Params types.Params, feemarketParams feemarkettypes.Params, evmParams evmtypes.Params, metadataList []banktypes.Metadata) error {
	// init fee market
	k.Logger(ctx).Info("init fee market", "erc20Params", feemarketParams.String())
	// set feeMarket baseFee
	k.feeMarketKeeper.SetBaseFee(ctx, feemarketParams.BaseFee.BigInt())
	// set feeMarket blockGasUsed
	k.feeMarketKeeper.SetBlockGasUsed(ctx, 0)
	// init feeMarket module erc20Params
	k.feeMarketKeeper.SetParams(ctx, feemarketParams)

	// init evm
	k.Logger(ctx).Info("init evm", "erc20Params", evmParams.String())
	k.evmKeeper.SetParams(ctx, evmParams)

	// init erc20
	k.Logger(ctx).Info("init erc20", "erc20Params", erc20Params.String())
	k.SetParams(ctx, erc20Params)

	// init contract
	if err := k.initSystemContract(ctx); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// init coin
	for _, metadata := range metadataList {
		k.Logger(ctx).Info("register coin", "coin", metadata.String())
		_, err := k.RegisterCoin(ctx, metadata)
		if err != nil {
			return sdkerrors.Wrapf(types.ErrInvalidMetadata, fmt.Sprintf("base %s, display %s, error %s",
				metadata.Base, metadata.Display, err.Error()))
		}
		//ctx.EventManager().EmitEvent(sdk.NewEvent(
		//	erc20types.EventTypeRegisterCoin,
		//	sdk.NewAttribute(erc20types.AttributeKeyCosmosCoin, pair.Denom),
		//	sdk.NewAttribute(erc20types.AttributeKeyFIP20Token, pair.Fip20Address),
		//))
	}
	return nil
}

func (k Keeper) initSystemContract(ctx sdk.Context) error {
	// ensure erc20 module account is set on genesis
	if acc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		return errors.New("the erc20 module account has not been set")
	}
	for _, contract := range contracts.GetInitContracts() {
		if len(contract.Code) <= 0 || contract.Address == common.HexToAddress(contracts.EmptyEvmAddress) {
			return errors.New("invalid contract")
		}
		if err := k.evmKeeper.CreateContractWithCode(ctx, contract.Address, contract.Code); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) DeployTokenUpgrade(ctx sdk.Context, from common.Address, name, symbol string, decimals uint8, origin bool) (common.Address, error) {
	k.Logger(ctx).Info("deploy token upgrade", "name", name, "symbol", symbol, "decimals", decimals)

	tokenContract := contracts.GetERC20(ctx.BlockHeight())
	if origin {
		tokenContract = contracts.GetWFX(ctx.BlockHeight())
	}

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
	erc1967Proxy := contracts.GetERC1967Proxy(ctx.BlockHeight())

	if len(logicData) == 0 {
		logicData = []byte{}
	}
	return k.DeployContract(ctx, from, erc1967Proxy.ABI, erc1967Proxy.Bin, logicAddr, logicData)
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
