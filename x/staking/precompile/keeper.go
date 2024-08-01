package precompile

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Keeper struct {
	bankKeeper       BankKeeper
	distrKeeper      DistrKeeper
	distrMsgServer   distrtypes.MsgServer
	stakingKeeper    StakingKeeper
	stakingMsgServer stakingtypes.MsgServer
	stakingDenom     string
}

func (k Keeper) NewStakingCoin(amount *big.Int) sdk.Coin {
	return sdk.NewCoin(k.stakingDenom, sdkmath.NewIntFromBigInt(amount))
}
