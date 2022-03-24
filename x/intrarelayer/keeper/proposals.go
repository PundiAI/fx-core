package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/functionx/fx-core/x/intrarelayer/types"
	"github.com/functionx/fx-core/x/intrarelayer/types/contracts"
)

// ModuleInit export to init module
func (k Keeper) ModuleInit(ctx sdk.Context, enableIntrarelayer, enableEvmHook bool, ibcTransferTimeoutHeight uint64) {
	k.SetParams(ctx, types.Params{
		EnableIntrarelayer:       enableIntrarelayer,
		EnableEVMHook:            enableEvmHook,
		IbcTransferTimeoutHeight: ibcTransferTimeoutHeight,
	})
	// ensure intrarelayer module account is set on genesis
	if acc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		panic("the intrarelayer module account has not been set")
	}
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
	addr, err := k.DeployFIP20Contract(ctx, name, symbol, decimals, coinMetadata.Base == evmParams.EvmDenom)
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

// DeployFIP20Contract creates and deploys an FIP20 contract on the EVM with the intrarelayer module account as owner
func (k Keeper) DeployFIP20Contract(ctx sdk.Context, name, symbol string, decimals uint8, origin ...bool) (common.Address, error) {
	contract := contracts.FIP20Contract
	if len(origin) > 0 && origin[0] {
		contract = contracts.WFXContract
	}

	ctorArgs, err := contract.ABI.Pack("", name, symbol, decimals)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "coin metadata is invalid  %s", name)
	}

	data := make([]byte, len(contract.Bin)+len(ctorArgs))
	copy(data[:len(contract.Bin)], contract.Bin)
	copy(data[len(contract.Bin):], ctorArgs)

	nonce, err := k.accountKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleAddress, nonce)
	_, err = k.CallEVMWithPayload(ctx, types.ModuleAddress, nil, data)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to deploy contract for %s", name)
	}

	return contractAddr, nil
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
