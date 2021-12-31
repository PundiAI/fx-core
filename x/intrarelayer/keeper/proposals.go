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

func (k Keeper) InitIntrarelayer(ctx sdk.Context, p *types.InitIntrarelayerParamsProposal) error {
	if !k.evmKeeper.HasInit(ctx) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "evm module has not init")
	}
	k.SetParams(ctx, *p.Params)

	// ensure intrarelayer module account is set on genesis
	if acc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		panic("the intrarelayer module account has not been set")
	}
	return nil
}

// RegisterCoin deploys an erc20 contract and creates the token pair for the cosmos coin
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

	addr, err := k.DeployERC20Contract(ctx, name, symbol, decimals)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to create wrapped coin denom metadata for ERC20")
	}

	pair := types.NewTokenPair(addr, coinMetadata.Base, true, types.OWNER_MODULE)

	k.SetTokenPair(ctx, pair)
	k.SetDenomMap(ctx, pair.Denom, pair.GetID())
	k.SetERC20Map(ctx, common.HexToAddress(pair.Erc20Address), pair.GetID())

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

// DeployERC20Contract creates and deploys an ERC20 contract on the EVM with the intrarelayer module account as owner
func (k Keeper) DeployERC20Contract(ctx sdk.Context, name, symbol string, decimals uint8) (common.Address, error) {
	ctorArgs, err := contracts.ERC20RelayContract.ABI.Pack("", name, symbol, decimals)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "coin metadata is invalid  %s", name)
	}

	data := make([]byte, len(contracts.ERC20RelayContract.Bin)+len(ctorArgs))
	copy(data[:len(contracts.ERC20RelayContract.Bin)], contracts.ERC20RelayContract.Bin)
	copy(data[len(contracts.ERC20RelayContract.Bin):], ctorArgs)

	nonce, err := k.accountKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleAddress, nonce)
	_, err = k.CallEVMWithPayloadWithModule(ctx, nil, data)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to deploy contract for %s", name)
	}

	return contractAddr, nil
}

// RegisterERC20 creates a cosmos coin and registers the token pair between the coin and the ERC20
func (k Keeper) RegisterERC20(ctx sdk.Context, contract common.Address) (*types.TokenPair, error) {
	params := k.GetParams(ctx)
	if !params.EnableIntrarelayer {
		return nil, sdkerrors.Wrap(types.ErrInternalTokenPair, "intrarelaying is currently disabled by governance")
	}

	if k.IsERC20Registered(ctx, contract) {
		return nil, sdkerrors.Wrapf(types.ErrInternalTokenPair, "token ERC20 contract already registered: %s", contract.String())
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

// CreateCoinMetadata generates the metadata to represent the ERC20 token on evmos.
func (k Keeper) CreateCoinMetadata(ctx sdk.Context, contract common.Address) (*banktypes.Metadata, string, string, error) {
	strContract := contract.String()

	erc20Data, err := k.QueryERC20(ctx, contract)
	if err != nil {
		return nil, "", "", err
	}

	meta := k.bankKeeper.GetDenomMetaData(ctx, types.CreateDenom(strContract))
	if len(meta.Base) > 0 {
		// metadata already exists; exit
		return nil, "", "", sdkerrors.Wrapf(types.ErrInternalTokenPair, "coin denomination already registered")
	}

	if k.IsDenomRegistered(ctx, types.CreateDenom(strContract)) {
		return nil, "", "", sdkerrors.Wrapf(types.ErrInternalTokenPair, "coin denomination already registered: %s", erc20Data.Name)
	}

	baseDenom := types.CreateDenom(strContract)
	// create a bank denom metadata based on the ERC20 token ABI details
	metadata := banktypes.Metadata{
		Description: erc20Data.Name,
		Base:        baseDenom,
		Display:     erc20Data.Symbol,
		// NOTE: Denom units MUST be increasing
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    baseDenom,
				Exponent: 0,
			},
			{
				Denom:    erc20Data.Symbol,
				Exponent: uint32(erc20Data.Decimals),
			},
		},
	}
	symbol := erc20Data.Symbol

	if err := metadata.Validate(); err != nil {
		return nil, "", "", sdkerrors.Wrapf(err, "ERC20 token data is invalid for contract %s", strContract)
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
