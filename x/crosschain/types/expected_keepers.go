package types

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/ibc-go/v3/modules/core/exported"
	"github.com/ethereum/go-ethereum/common"

	tranfsertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

type StakingKeeper interface {
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)

	GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, found bool)
	Delegate(ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdk.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (newShares sdk.Dec, err error)
	Undelegate(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, sharesAmount sdk.Dec) (time.Time, error)
	BeginRedelegation(ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, sharesAmount sdk.Dec) (completionTime time.Time, err error)
	GetUnbondingDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (ubd stakingtypes.UnbondingDelegation, found bool)

	// migrate v3
	RemoveDelegation(ctx sdk.Context, delegation stakingtypes.Delegation)
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator
}

type DistributionKeeper interface {
	WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error)
}

// BankKeeper defines the expected bank keeper methods
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetDenomMetaData(ctx sdk.Context, denom string) (bank.Metadata, bool)
}

type Erc20Keeper interface {
	RelayConvertCoin(ctx sdk.Context, sender sdk.AccAddress, receiver common.Address, coin sdk.Coin) error
	RelayConvertDenomToOne(ctx sdk.Context, from sdk.AccAddress, coin sdk.Coin) (sdk.Coin, error)
	RelayConvertDenomToMany(ctx sdk.Context, from sdk.AccAddress, coin sdk.Coin, target string) (sdk.Coin, error)
}

type IBCChannelKeeper interface {
	GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, exported.ClientState, error)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
}

type IBCTransferKeeper interface {
	Transfer(goCtx context.Context, msg *tranfsertypes.MsgTransfer) (*tranfsertypes.MsgTransferResponse, error)

	SetDenomTrace(ctx sdk.Context, denomTrace tranfsertypes.DenomTrace)
}
