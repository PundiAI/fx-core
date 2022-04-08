package keeper

import (
	"encoding/hex"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/functionx/fx-core/x/intrarelayer/types"
	"github.com/functionx/fx-core/x/intrarelayer/types/contracts"
)

// ModuleInit export to init module
func (k Keeper) ModuleInit(ctx sdk.Context, enableIntrarelayer, enableEvmHook bool, ibcTransferTimeoutHeight uint64) error {
	k.SetParams(ctx, types.Params{
		EnableIntrarelayer:       enableIntrarelayer,
		EnableEVMHook:            enableEvmHook,
		IbcTransferTimeoutHeight: ibcTransferTimeoutHeight,
	})
	// ensure intrarelayer module account is set on genesis
	if acc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		return errors.New("the intrarelayer module account has not been set")
	}
	//init system contract
	if err := contracts.InitSystemContract(ctx, k.evmKeeper); err != nil {
		return fmt.Errorf("intrarelayer module init system contract error: %v", err)
	}
	return nil
}

// RegisterCoin deploys an fip20 contract and creates the token pair for the cosmos coin
func (k Keeper) RegisterCoin(ctx sdk.Context, coinMetadata banktypes.Metadata) (*types.TokenPair, error) {
	params := k.GetParams(ctx)
	if !params.EnableIntrarelayer {
		return nil, sdkerrors.Wrap(types.ErrInternalTokenPair, "intrarelaying is currently disabled by governance")
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

	if k.IsDenomRegistered(ctx, coinMetadata.Base) {
		return nil, sdkerrors.Wrapf(types.ErrInternalTokenPair, "coin denomination already registered: %s", name)
	}

	if err := k.verifyMetadata(ctx, coinMetadata); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInternalTokenPair, "coin metadata is invalid %s, error %v", name, err)
	}

	evmParams := k.evmKeeper.GetParams(ctx)
	addr, err := k.DeployTokenUpgrade(ctx, types.ModuleAddress, name, symbol, decimals, coinMetadata.Base == evmParams.EvmDenom)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to create wrapped coin denom metadata for FIP20")
	}

	pair := types.NewTokenPair(addr, coinMetadata.Base, true, types.OWNER_MODULE)

	k.SetTokenPair(ctx, pair)
	k.SetDenomMap(ctx, pair.Denom, pair.GetID())
	k.SetFIP20Map(ctx, common.HexToAddress(pair.Fip20Address), pair.GetID())

	return &pair, nil
}

func (k Keeper) verifyMetadata(ctx sdk.Context, coinMetadata banktypes.Metadata) error {
	meta := k.bankKeeper.GetDenomMetaData(ctx, coinMetadata.Base)
	if len(meta.Base) == 0 { //new coin metadata
		k.bankKeeper.SetDenomMetaData(ctx, coinMetadata)
		return nil
	}
	// If it already existed, Check that is equal to what is stored
	return equalMetadata(meta, coinMetadata)
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
	_, err = k.CallEVMWithPayload(ctx, from, nil, data)
	if err != nil {
		return common.Address{}, sdkerrors.Wrap(err, "failed to deploy contract")
	}
	return contractAddr, nil
}

func (k Keeper) DeployTokenUpgrade(ctx sdk.Context, from common.Address, name, symbol string, decimals uint8, origin bool) (common.Address, error) {
	k.Logger(ctx).Info("deploy token upgrade", "name", name, "symbol", symbol, "decimals", decimals)

	ccLogicABI, foundLogic := contracts.GetABI(ctx.BlockHeight(), contracts.FIP20UpgradeType)
	logicAddr := common.HexToAddress(contracts.FIP20UpgradeCodeAddress)
	if origin {
		ccLogicABI, foundLogic = contracts.GetABI(ctx.BlockHeight(), contracts.WFXUpgradeType)
		logicAddr = common.HexToAddress(contracts.WFXUpgradeCodeAddress)
	}
	if !foundLogic {
		return common.Address{}, sdkerrors.Wrapf(types.ErrInvalidContract, "logic contract not found(origin:%v, logic:%v)", origin, foundLogic)
	}

	//deploy proxy
	proxy, err := k.DeployERC1967Proxy(ctx, from, logicAddr)
	if err != nil {
		return common.Address{}, err
	}
	err = k.InitializeUpgradable(ctx, from, proxy, ccLogicABI, name, symbol, decimals, types.ModuleAddress)
	return proxy, err
}

func (k Keeper) DeployERC1967Proxy(ctx sdk.Context, from, logicAddr common.Address, logicData ...byte) (common.Address, error) {
	k.Logger(ctx).Info("deploy erc1967 proxy", "logic", logicAddr.String(), "data", hex.EncodeToString(logicData))
	proxyABI, found := contracts.GetABI(ctx.BlockHeight(), contracts.ERC1967ProxyType)
	if !found {
		return common.Address{}, sdkerrors.Wrap(types.ErrInvalidContract, "erc1967 proxy abi not found")
	}
	proxyBin, found := contracts.GetBin(ctx.BlockHeight(), contracts.ERC1967ProxyType)
	if !found {
		return common.Address{}, sdkerrors.Wrap(types.ErrInvalidContract, "erc1967 proxy bin not found")
	}

	if len(logicData) == 0 {
		logicData = []byte{}
	}
	return k.DeployContract(ctx, from, proxyABI, proxyBin, logicAddr, logicData)
}

func (k Keeper) InitializeUpgradable(ctx sdk.Context, from, contract common.Address, abi abi.ABI, data ...interface{}) error {
	k.Logger(ctx).Info("initialize upgradable", "contract", contract.Hex())
	_, err := k.CallEVM(ctx, abi, from, contract, "initialize", data...)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to initialize contract")
	}
	return nil
}

// RegisterFIP20 creates a cosmos coin and registers the token pair between the coin and the FIP20
func (k Keeper) RegisterFIP20(ctx sdk.Context, contract common.Address) (*types.TokenPair, error) {
	params := k.GetParams(ctx)
	if !params.EnableIntrarelayer {
		return nil, sdkerrors.Wrap(types.ErrInternalTokenPair, "intrarelaying is currently disabled by governance")
	}

	if k.IsFIP20Registered(ctx, contract) {
		return nil, sdkerrors.Wrapf(types.ErrInternalTokenPair, "token FIP20 contract already registered: %s", contract.String())
	}

	_, baseDenom, _, err := k.CreateCoinMetadata(ctx, contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to create wrapped coin denom metadata for FIP20")
	}

	pair := types.NewTokenPair(contract, baseDenom, true, types.OWNER_EXTERNAL)
	k.SetTokenPair(ctx, pair)
	k.SetDenomMap(ctx, pair.Denom, pair.GetID())
	k.SetFIP20Map(ctx, common.HexToAddress(pair.Fip20Address), pair.GetID())
	return &pair, nil
}

// CreateCoinMetadata generates the metadata to represent the FIP20 token on evmos.
func (k Keeper) CreateCoinMetadata(ctx sdk.Context, contract common.Address) (*banktypes.Metadata, string, string, error) {
	strContract := contract.String()

	fip20Data, err := k.QueryFIP20(ctx, contract)
	if err != nil {
		return nil, "", "", err
	}

	meta := k.bankKeeper.GetDenomMetaData(ctx, types.CreateDenom(strContract))
	if len(meta.Base) > 0 {
		// metadata already exists; exit
		return nil, "", "", sdkerrors.Wrapf(types.ErrInternalTokenPair, "coin denomination already registered")
	}

	if k.IsDenomRegistered(ctx, types.CreateDenom(strContract)) {
		return nil, "", "", sdkerrors.Wrapf(types.ErrInternalTokenPair, "coin denomination already registered: %s", fip20Data.Name)
	}

	baseDenom := types.CreateDenom(strContract)
	// create a bank denom metadata based on the FIP20 token ABI details
	metadata := banktypes.Metadata{
		Description: fip20Data.Name,
		Base:        baseDenom,
		Display:     fip20Data.Symbol,
		// NOTE: Denom units MUST be increasing
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    baseDenom,
				Exponent: 0,
			},
			{
				Denom:    fip20Data.Symbol,
				Exponent: uint32(fip20Data.Decimals),
			},
		},
	}
	symbol := fip20Data.Symbol

	if err := metadata.Validate(); err != nil {
		return nil, "", "", sdkerrors.Wrapf(err, "FIP20 token data is invalid for contract %s", strContract)
	}

	k.bankKeeper.SetDenomMetaData(ctx, metadata)

	return &metadata, baseDenom, symbol, nil
}

// ToggleRelay toggles relaying for a given token pair
func (k Keeper) ToggleRelay(ctx sdk.Context, token string) (types.TokenPair, error) {
	id := k.GetTokenPairID(ctx, token)

	if len(id) == 0 {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrInternalTokenPair, "token %s not registered", token)
	}

	pair, found := k.GetTokenPair(ctx, id)
	if !found {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrInternalTokenPair, "not registered")
	}

	pair.Enabled = !pair.Enabled

	k.SetTokenPair(ctx, pair)
	return pair, nil
}
