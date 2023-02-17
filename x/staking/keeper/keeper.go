package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	storeKey      sdk.StoreKey
	accountKeeper types.AccountKeeper
	evmKeeper     types.EvmKeeper

	lpTokenModuleAddress common.Address
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, ak types.AccountKeeper, bk stakingtypes.BankKeeper, ps paramtypes.Subspace) *Keeper {
	return &Keeper{
		Keeper:               stakingkeeper.NewKeeper(cdc, key, ak, bk, ps),
		storeKey:             key,
		accountKeeper:        ak,
		evmKeeper:            nil,
		lpTokenModuleAddress: common.BytesToAddress(ak.GetModuleAddress(types.LPTokenOwnerModuleName)),
	}
}

func (k Keeper) Delegate(
	ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdk.Int, tokenSrc stakingtypes.BondStatus,
	validator stakingtypes.Validator, subtractAccount bool,
) (newShares sdk.Dec, err error) {
	newShares, err = k.Keeper.Delegate(ctx, delAddr, bondAmt, tokenSrc, validator, subtractAccount)
	if err != nil {
		return newShares, err
	}

	// todo - call evm contract
	// lpTokenContract, found := k.GetLPTokenContract(ctx, validator.GetOperator())
	// if !found {
	// 	return sdk.ZeroDec(), sdkerrors.ErrInvalidRequest.Wrapf("lpToken contract not found for validator")
	// }
	//
	// err = k.MintLPToken(ctx, lpTokenContract, delAddr, newShares)
	return newShares, err
}

func (k Keeper) MintLPToken(ctx sdk.Context, lpTokenContractAddr common.Address, to sdk.AccAddress, share sdk.Dec) error {
	erc20 := fxtypes.GetLPToken().ABI
	data, err := erc20.Pack("mint", common.BytesToAddress(to.Bytes()), share.BigInt())
	if err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("failed to pack data: %s", err.Error())
	}

	return k.callEVM(ctx, &lpTokenContractAddr, data)
}

func (k Keeper) Undelegate(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, sharesAmount sdk.Dec,
) (time.Time, error) {
	undelegate, err := k.Keeper.Undelegate(ctx, delAddr, valAddr, sharesAmount)
	if err != nil {
		return undelegate, err
	}

	// todo - call evm contract
	// lpTokenContract, found := k.GetLPTokenContract(ctx, valAddr)
	// if !found {
	// 	return undelegate, sdkerrors.ErrInvalidRequest.Wrapf("lpToken contract not found for validator")
	// }
	//
	// data, err := fxtypes.GetLPToken().ABI.Pack("burn", common.BytesToAddress(delAddr.Bytes()), sharesAmount.BigInt())
	// if err != nil {
	// 	return undelegate, sdkerrors.ErrInvalidRequest.Wrapf("failed to pack data: %s", err.Error())
	// }
	//
	// err = k.callEVM(ctx, &lpTokenContract, data)
	return undelegate, err
}

func (k Keeper) DeployLPToken(ctx sdk.Context, valAddr sdk.ValAddress) (common.Address, error) {
	_, err := k.accountKeeper.GetSequence(ctx, k.lpTokenModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, sdkerrors.ErrInvalidRequest.Wrapf("failed to get lpToken module address nonce")
	}

	// todo - call evm
	// erc1967Proxy := fxtypes.GetERC1967Proxy()
	//
	// lpToken := fxtypes.GetLPToken()
	// deployInputData, err := erc1967Proxy.ABI.Pack("", lpToken.Address, []byte{})
	// if err != nil {
	// 	return sdkerrors.ErrInvalidAddress.Wrapf("failed to pack deploy contract data")
	// }
	// data := make([]byte, len(erc1967Proxy.Bin)+len(deployInputData))
	// copy(data[:len(erc1967Proxy.Bin)], erc1967Proxy.Bin)
	// copy(data[len(erc1967Proxy.Bin):], deployInputData)
	//
	// if err = k.callEVM(ctx, nil, data); err != nil {
	// 	return sdkerrors.ErrInvalidRequest.Wrapf("failed to deploy lpToken contract: %s", err.Error())
	// }
	//
	// contractAddr := crypto.CreateAddress(k.lpTokenModuleAddress, nonce)
	// initializeInputData, err := lpToken.ABI.Pack("initialize", valAddr.String(), types.LPTokenSymbol, types.LPTokenDecimals, k.lpTokenModuleAddress)
	// if err != nil {
	// 	return sdkerrors.ErrInvalidRequest.Wrapf("failed to pack call initialize method data")
	// }
	//
	// err = k.callEVM(ctx, &contractAddr, initializeInputData)
	// if err != nil {
	// 	return sdkerrors.ErrInvalidRequest.Wrapf("failed to call initialize method: %s", err.Error())
	// }
	// k.setLPTokenContract(ctx, valAddr, contractAddr)
	return common.Address{}, nil
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

func (k *Keeper) GetLPTokenContract(ctx sdk.Context, valAddr sdk.ValAddress) (common.Address, bool) {
	kvStore := ctx.KVStore(k.storeKey)
	bz := kvStore.Get(types.GetLPTokenKey(valAddr))
	return common.BytesToAddress(bz), bz == nil
}

// func (k *Keeper) setLPTokenContract(ctx sdk.Context, valAddr sdk.ValAddress, lpTokenContract common.Address) {
// 	kvStore := ctx.KVStore(k.storeKey)
// 	kvStore.Set(types.GetLPTokenKey(valAddr), lpTokenContract.Bytes())
// }

func (k *Keeper) deleteLPTokenContract(ctx sdk.Context, valAddr sdk.ValAddress) {
	kvStore := ctx.KVStore(k.storeKey)
	kvStore.Delete(types.GetLPTokenKey(valAddr))
}

func (k *Keeper) callEVM(ctx sdk.Context, contract *common.Address, data []byte) error {
	k.Logger(ctx).Info("evmKeeper", "key", k.evmKeeper)
	_, err := k.evmKeeper.CallEVMWithData(ctx, k.lpTokenModuleAddress, contract, data, true)
	if err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("call evm failed: %s", err.Error())
	}
	return nil
}

func (k *Keeper) GetLPTokenModuleAddress() common.Address {
	return k.lpTokenModuleAddress
}

func (k *Keeper) IteratorValidators(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, stakingtypes.ValidatorsKey)
	return iterator
}
