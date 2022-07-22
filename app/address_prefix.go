package app

import (
	"math/big"

	bsctypes "github.com/functionx/fx-core/x/bsc/types"
	"github.com/functionx/fx-core/x/crosschain/types"
	polygontypes "github.com/functionx/fx-core/x/polygon/types"
	trontypes "github.com/functionx/fx-core/x/tron/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func init() {
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/-]{1,127}`
	})

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AddressPrefix, AddressPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(AddressPrefix+sdk.PrefixValidator+sdk.PrefixOperator, AddressPrefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(AddressPrefix+sdk.PrefixValidator+sdk.PrefixConsensus, AddressPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
	config.Seal()

	// votingPower = delegateToken / sdk.PowerReduction  --  sdk.TokensToConsensusPower(tokens Int)
	sdk.PowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil))

	types.RegisterValidateBasic(bsctypes.ModuleName, types.EthereumMsgValidateBasic{})
	types.RegisterValidateBasic(polygontypes.ModuleName, types.EthereumMsgValidateBasic{})
	types.RegisterValidateBasic(trontypes.ModuleName, trontypes.MsgValidateBasic{})
}
