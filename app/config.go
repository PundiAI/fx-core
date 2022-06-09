package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	feemarketkeeper "github.com/tharsis/ethermint/x/feemarket/keeper"

	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ethermint "github.com/tharsis/ethermint/types"
)

func init() {
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/-]{1,127}`
	})

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

	// set calculate base function
	feemarketkeeper.SetCalculateBaseFeeFunc(CalculateBaseFeeLimit)

	// set evm params
	evmtypes.SetDefaultParams(evmParams)
}

func ParseFunctionXChainID(_ string) (*big.Int, error) {
	return fxtypes.EIP155ChainID(), nil
}

func ValidFunctionXChainID(chainID string) bool {
	if len(chainID) > 48 {
		return false
	}
	return true
}

var (
	MinBaseFee      = sdk.NewInt(500 * 1e9)
	MaxBaseFee      = sdk.NewIntFromUint64(math.MaxUint64 - 1).Mul(sdk.NewInt(1e9))
	DefaultGasLimit = sdk.NewIntFromUint64(3 * 1e7)
)

func CalculateBaseFeeLimit(ctx sdk.Context, k feemarketkeeper.Keeper) *big.Int {
	params := k.GetParams(ctx)

	// Ignore the calculation if not enabled
	if !params.IsBaseFeeEnabled(ctx.BlockHeight()) {
		return nil
	}

	// If the current block is the first EIP-1559 block, return the base fee
	// defined in the parameters (DefaultBaseFee if it hasn't been changed by
	// governance).
	if ctx.BlockHeight() == params.EnableHeight {
		return params.BaseFee.BigInt()
	}

	// get the block gas used and the base fee values for the parent block.
	// NOTE: this is not the parent's base fee but the current block's base fee,
	// as it is retrieved from the transient store, which is committed to the
	// persistent KVStore after EndBlock (ABCI Commit).
	parentBaseFee := params.BaseFee.BigInt()
	if parentBaseFee == nil {
		return nil
	}

	parentGasUsed := k.GetBlockGasWanted(ctx)

	gasLimit := DefaultGasLimit.BigInt()
	consParams := ctx.ConsensusParams()
	// NOTE: a MaxGas equal to -1 means that block gas is unlimited
	if consParams != nil && consParams.Block.MaxGas > -1 {
		gasLimit = big.NewInt(consParams.Block.MaxGas)
	}

	// CONTRACT: ElasticityMultiplier cannot be 0 as it's checked in the params
	// validation
	parentGasTargetBig := new(big.Int).Div(gasLimit, new(big.Int).SetUint64(uint64(params.ElasticityMultiplier)))
	if !parentGasTargetBig.IsUint64() {
		return nil
	}

	parentGasTarget := parentGasTargetBig.Uint64()
	baseFeeChangeDenominator := new(big.Int).SetUint64(uint64(params.BaseFeeChangeDenominator))

	// If the parent gasUsed is the same as the target, the baseFee remains unchanged.
	if parentGasUsed == parentGasTarget {
		return new(big.Int).Set(parentBaseFee)
	}

	if parentGasUsed > parentGasTarget {
		// If the parent block used more gas than its target, the baseFee should increase.
		gasUsedDelta := new(big.Int).SetUint64(parentGasUsed - parentGasTarget)
		x := new(big.Int).Mul(parentBaseFee, gasUsedDelta)
		y := x.Div(x, parentGasTargetBig)
		baseFeeDelta := math.BigMax(
			x.Div(y, baseFeeChangeDenominator),
			common.Big1,
		)

		//NOTE: compare with MaxBaseFee, choose the smaller
		return math.BigMin(x.Add(parentBaseFee, baseFeeDelta), MaxBaseFee.BigInt())
	}

	// Otherwise if the parent block used less gas than its target, the baseFee should decrease.
	gasUsedDelta := new(big.Int).SetUint64(parentGasTarget - parentGasUsed)
	x := new(big.Int).Mul(parentBaseFee, gasUsedDelta)
	y := x.Div(x, parentGasTargetBig)
	baseFeeDelta := x.Div(y, baseFeeChangeDenominator)

	//NOTE: compare with MinBaseFee, choose the larger
	return math.BigMax(x.Sub(parentBaseFee, baseFeeDelta), MinBaseFee.BigInt())
}

var (
	evmParams = evmtypes.Params{
		EvmDenom:     fxtypes.DefaultDenom,
		EnableCreate: true,
		EnableCall:   true,
		ExtraEIPs:    nil,
		ChainConfig:  evmtypes.DefaultChainConfig(),
	}
)
