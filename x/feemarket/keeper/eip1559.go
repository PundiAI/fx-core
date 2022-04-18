package keeper

import (
	"math/big"

	fxtypes "github.com/functionx/fx-core/types"

	"github.com/functionx/fx-core/x/feemarket/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

// CalculateBaseFee calculates the base fee for the current block. This is only calculated once per
// block during BeginBlock. If the NoBaseFee parameter is enabled or below activation height, this function returns nil.
// NOTE: This code is inspired from the go-ethereum EIP1559 implementation and adapted to Cosmos SDK-based
// chains. For the canonical code refer to: https://github.com/ethereum/go-ethereum/blob/master/consensus/misc/eip1559.go
func (k Keeper) CalculateBaseFee(ctx sdk.Context) *big.Int {
	params := k.GetParams(ctx)

	// If the current block is the first EIP-1559 block, return the InitialBaseFee.
	if ctx.BlockHeight() == fxtypes.EvmSupportBlock() {
		return params.BaseFee.BigInt()
	}

	// get the block gas used and the base fee values for the parent block.
	parentBaseFee := params.BaseFee.BigInt()
	if parentBaseFee == nil {
		return nil
	}

	minBaseFee := params.MinBaseFee.BigInt()
	if minBaseFee == nil {
		minBaseFee = types.MinBaseFee.BigInt()
	}

	maxBaseFee := params.MaxBaseFee.BigInt()
	if maxBaseFee == nil || maxBaseFee.Cmp(big.NewInt(0)) <= 0 {
		maxBaseFee = types.MaxBaseFee.BigInt()
	}

	gasLimit := new(big.Int).SetUint64(math.MaxUint64)
	if !params.MaxGas.IsNil() && params.MaxGas.GT(sdk.ZeroInt()) {
		gasLimit = params.MaxGas.BigInt()
	} else {
		consParams := ctx.ConsensusParams()
		if consParams != nil && consParams.Block.MaxGas > -1 {
			gasLimit = big.NewInt(consParams.Block.MaxGas)
		}
	}

	parentGasTargetBig := new(big.Int).Div(gasLimit, new(big.Int).SetUint64(uint64(params.ElasticityMultiplier)))
	if !parentGasTargetBig.IsUint64() {
		return nil
	}

	parentGasTarget := parentGasTargetBig.Uint64()
	baseFeeChangeDenominator := new(big.Int).SetUint64(uint64(params.BaseFeeChangeDenominator))

	parentGasUsed := k.GetBlockGasUsed(ctx)
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
		return math.BigMin(x.Add(parentBaseFee, baseFeeDelta), maxBaseFee)
	}

	// Otherwise if the parent block used less gas than its target, the baseFee should decrease.
	gasUsedDelta := new(big.Int).SetUint64(parentGasTarget - parentGasUsed)
	x := new(big.Int).Mul(parentBaseFee, gasUsedDelta)
	y := x.Div(x, parentGasTargetBig)
	baseFeeDelta := x.Div(y, baseFeeChangeDenominator)

	return math.BigMax(x.Sub(parentBaseFee, baseFeeDelta), minBaseFee)
}
