package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/staking/types"
)

// Implements ValidatorSet interface
var _ stakingtypes.ValidatorSet = Keeper{}

// Implements DelegationSet interface
var _ stakingtypes.DelegationSet = Keeper{}

type Keeper struct {
	stakingkeeper.Keeper
	storeKey      storetypes.StoreKey
	accountKeeper types.AccountKeeper
	bankKeeper    stakingtypes.BankKeeper
	evmKeeper     types.EvmKeeper

	lpTokenModuleAddress common.Address
}

func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, ak types.AccountKeeper, bk stakingtypes.BankKeeper, ps paramtypes.Subspace) *Keeper {
	return &Keeper{
		Keeper:               stakingkeeper.NewKeeper(cdc, key, ak, bk, ps),
		storeKey:             key,
		accountKeeper:        ak,
		bankKeeper:           bk,
		evmKeeper:            nil,
		lpTokenModuleAddress: common.BytesToAddress(ak.GetModuleAddress(types.LPTokenOwnerModuleName)),
	}
}

func (k Keeper) Delegate(
	ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdkmath.Int, tokenSrc stakingtypes.BondStatus,
	validator stakingtypes.Validator, subtractAccount bool,
) (newShares sdk.Dec, err error) {
	newShares, err = k.Keeper.Delegate(ctx, delAddr, bondAmt, tokenSrc, validator, subtractAccount)
	if err != nil {
		return newShares, err
	}

	lpTokenContract, found := k.GetValidatorLPToken(ctx, validator.GetOperator())
	if !found {
		return sdk.ZeroDec(), errortypes.ErrInvalidRequest.Wrapf("lpToken contract not found for validator")
	}

	err = k.MintLPToken(ctx, lpTokenContract, delAddr, newShares)
	return newShares, err
}

func (k Keeper) TransferDelegate(ctx sdk.Context, valAddr sdk.ValAddress, fromAddr, toAddr sdk.AccAddress, shares sdk.Dec) error {
	if k.HasReceivingRedelegation(ctx, fromAddr, valAddr) {
		return stakingtypes.ErrTransitiveRedelegation
	}

	returnAmount, err := k.Keeper.Unbond(ctx, fromAddr, valAddr, shares)
	if err != nil {
		return err
	}
	if returnAmount.IsZero() {
		return types.ErrTinyTransferAmount
	}

	val, found := k.Keeper.GetValidator(ctx, valAddr)
	if !found {
		return stakingtypes.ErrNoValidatorFound
	}

	_, err = k.Keeper.Delegate(ctx, toAddr, returnAmount, val.GetStatus(), val, false)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) DeployLPToken(ctx sdk.Context, valAddr sdk.ValAddress) (common.Address, error) {
	lpToken := fxtypes.GetLPToken()
	contractAddr, err := k.evmKeeper.DeployUpgradableContract(ctx, k.lpTokenModuleAddress, lpToken.Address, nil,
		&lpToken.ABI, valAddr.String(), types.LPTokenSymbol, types.LPTokenDecimals)
	if err != nil {
		return common.Address{}, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventCreateLPToken,
		sdk.NewAttribute(types.AttributeKeyLPTokenAddress, contractAddr.String()),
	))
	k.setLPTokenContract(ctx, valAddr, contractAddr)
	return contractAddr, nil
}

func (k Keeper) MintLPToken(ctx sdk.Context, lpTokenContractAddr common.Address, to sdk.AccAddress, share sdk.Dec) error {
	method := types.MethodMintToken
	args := []interface{}{
		common.BytesToAddress(to.Bytes()),
		share.BigInt(),
	}
	_, err := k.evmKeeper.ApplyContract(ctx, k.lpTokenModuleAddress, lpTokenContractAddr, fxtypes.GetLPToken().ABI, method, args...)
	return err
}

func (k Keeper) BurnLPToken(ctx sdk.Context, lpTokenContractAddr common.Address, to sdk.AccAddress, share sdk.Dec) error {
	method := types.MethodBurnToken
	args := []interface{}{
		common.BytesToAddress(to.Bytes()),
		share.BigInt(),
	}
	_, err := k.evmKeeper.ApplyContract(ctx, k.lpTokenModuleAddress, lpTokenContractAddr, fxtypes.GetLPToken().ABI, method, args...)
	return err
}

func (k Keeper) SelfDestructLPToken(ctx sdk.Context, valAddr sdk.ValAddress) error {
	lpTokenContract, found := k.GetValidatorLPToken(ctx, valAddr)
	if !found {
		return errorsmod.Wrapf(types.ErrLPTokenNotFound, "validator %s", valAddr.String())
	}
	method := types.MethodSelfDestruct
	if _, err := k.evmKeeper.ApplyContract(ctx, k.lpTokenModuleAddress, lpTokenContract, fxtypes.GetLPToken().ABI, method); err != nil {
		return err
	}
	k.deleteLPTokenContract(ctx, valAddr)
	return nil
}

func (k *Keeper) SetHooks(sh stakingtypes.StakingHooks) *Keeper {
	k.Keeper.SetHooks(sh)
	return k
}

func (k *Keeper) SetEvmKeeper(evmKeeper types.EvmKeeper) *Keeper {
	if k.evmKeeper != nil {
		panic("cannot set evm keeper twice")
	}
	k.evmKeeper = evmKeeper
	return k
}

func (k *Keeper) GetValidatorLPToken(ctx sdk.Context, valAddr sdk.ValAddress) (common.Address, bool) {
	kvStore := ctx.KVStore(k.storeKey)
	bz := kvStore.Get(types.GetValidatorLPTokenKey(valAddr))
	return common.BytesToAddress(bz), bz != nil
}

func (k *Keeper) GetLPTokenValidator(ctx sdk.Context, lpTokenContract common.Address) (sdk.ValAddress, bool) {
	kvStore := ctx.KVStore(k.storeKey)
	bz := kvStore.Get(types.GetLPTokenValidatorKey(lpTokenContract))
	return bz, bz != nil
}

func (k *Keeper) setLPTokenContract(ctx sdk.Context, valAddr sdk.ValAddress, lpTokenContract common.Address) {
	kvStore := ctx.KVStore(k.storeKey)
	kvStore.Set(types.GetValidatorLPTokenKey(valAddr), lpTokenContract.Bytes())
	kvStore.Set(types.GetLPTokenValidatorKey(lpTokenContract), valAddr.Bytes())
}

func (k *Keeper) deleteLPTokenContract(ctx sdk.Context, valAddr sdk.ValAddress) {
	kvStore := ctx.KVStore(k.storeKey)
	key := types.GetValidatorLPTokenKey(valAddr)
	lpTokenByte := kvStore.Get(key)
	kvStore.Delete(key)
	kvStore.Delete(types.GetLPTokenValidatorKey(common.BytesToAddress(lpTokenByte)))
}

func (k *Keeper) GetLPTokenModuleAddress() common.Address {
	return k.lpTokenModuleAddress
}

func (k *Keeper) IteratorValidators(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, stakingtypes.ValidatorsKey)
	return iterator
}
