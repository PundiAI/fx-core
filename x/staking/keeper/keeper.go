package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v3/x/staking/types"
)

// Implements ValidatorSet interface
var _ stakingtypes.ValidatorSet = Keeper{}

// Implements DelegationSet interface
var _ stakingtypes.DelegationSet = Keeper{}

type Keeper struct {
	stakingkeeper.Keeper
	accountKeeper types.AccountKeeper
	evmKeeper     types.EvmKeeper

	lpTokenModuleAddress common.Address
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, ps paramtypes.Subspace) Keeper {
	return Keeper{
		Keeper:               stakingkeeper.NewKeeper(cdc, key, ak, bk, ps),
		accountKeeper:        ak,
		evmKeeper:            &types.MockEvmKeeper{},
		lpTokenModuleAddress: common.BytesToAddress(ak.GetModuleAddress(types.LpTokenName)),
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

	// todo - call evm to mint lpToken
	return newShares, nil
}

func (k Keeper) Undelegate(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, sharesAmount sdk.Dec,
) (time.Time, error) {
	undelegate, err := k.Keeper.Undelegate(ctx, delAddr, valAddr, sharesAmount)
	if err != nil {
		return undelegate, err
	}

	// toto - call evm to burn lpToken
	return undelegate, err
}

// AfterValidatorCreated - call hook if registered
func (k Keeper) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	k.Keeper.AfterValidatorCreated(ctx, valAddr)

	// todo - call evm to create validator lpToken
}

// AfterValidatorRemoved - call hook if registered
func (k Keeper) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	k.Keeper.AfterValidatorCreated(ctx, valAddr)

	// todo -- call evm to destroy validator lpToken
}

func (k *Keeper) SetHooks(sh stakingtypes.StakingHooks) *Keeper {
	k.Keeper.SetHooks(sh)
	return k
}
