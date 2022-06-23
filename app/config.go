package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"

	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ethermint "github.com/evmos/ethermint/types"
)

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(fxtypes.AddressPrefix, fxtypes.AddressPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(fxtypes.AddressPrefix+sdk.PrefixValidator+sdk.PrefixOperator, fxtypes.AddressPrefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(fxtypes.AddressPrefix+sdk.PrefixValidator+sdk.PrefixConsensus, fxtypes.AddressPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
	config.SetCoinType(60)
	config.Seal()

	// votingPower = delegateToken / sdk.PowerReduction  --  sdk.TokensToConsensusPower(tokens Int)
	sdk.DefaultPowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil))

	if err := sdk.RegisterDenom(fxtypes.DefaultDenom, sdk.NewDec(18)); err != nil {
		panic(err)
	}

	// set chain id function
	ethermint.SetParseChainIDFunc(ParseFunctionXChainID)
	ethermint.SetValidChainIDFunc(ValidFunctionXChainID)

	// set calculate base fee function
	feemarketkeeper.SetCalculateBaseFeeFunc(CalculateBaseFee)
}

func ParseFunctionXChainID(_ string) (*big.Int, error) {
	return fxtypes.EIP155ChainID(), nil
}

func ValidFunctionXChainID(chainID string) bool {
	return fxtypes.ChainId() == chainID
}

func CalculateBaseFee(ctx sdk.Context, k feemarketkeeper.Keeper) *big.Int {
	// default calculate base fee
	baseFee := feemarketkeeper.DefaultCalculateBaseFee(ctx, k)
	if baseFee == nil {
		return nil
	}
	// check min gas price
	params := k.GetParams(ctx)
	if params.MinGasPrice.IsNil() || params.MinGasPrice.IsZero() {
		return baseFee
	}
	return math.BigMax(baseFee, params.MinGasPrice.TruncateInt().BigInt())
}
