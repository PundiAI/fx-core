package ante

import (
	"math"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type CheckTxFeees struct {
	bypassMsgTypesMap    map[string]bool
	maxBypassMsgGasUsage uint64
}

func NewCheckTxFeees(bypassMinFeeMsgTypes []string, maxBypassMinFeeMsgGasUsage uint64) CheckTxFeees {
	bypassMinFeeMsgTypesMap := make(map[string]bool, len(bypassMinFeeMsgTypes))
	for _, msgType := range bypassMinFeeMsgTypes {
		bypassMinFeeMsgTypesMap[msgType] = true
	}
	return CheckTxFeees{
		bypassMsgTypesMap:    bypassMinFeeMsgTypesMap,
		maxBypassMsgGasUsage: maxBypassMinFeeMsgGasUsage,
	}
}

func (ctf CheckTxFeees) Check(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
	return ctf.checkTxFeeWithValidatorMinGasPrices(ctx, tx)
}

func (ctf CheckTxFeees) isByPassMinFee(msgs []sdk.Msg, gas uint64) bool {
	return ctf.bypassMinFeeMsgs(msgs) && ctf.isBypassMinFeeMsgGasUsage(msgs, gas)
}

func (ctf CheckTxFeees) bypassMinFeeMsgs(msgs []sdk.Msg) bool {
	result := false
	for _, msg := range msgs {
		result = true
		if _, ok := ctf.bypassMsgTypesMap[sdk.MsgTypeURL(msg)]; !ok {
			return false
		}
	}

	return result
}

func (ctf CheckTxFeees) isBypassMinFeeMsgGasUsage(msgs []sdk.Msg, gas uint64) bool {
	return uint64(len(msgs))*ctf.maxBypassMsgGasUsage >= gas
}

// checkTxFeeWithValidatorMinGasPrices implements the default fee logic, where the minimum price per
// unit of gas is fixed and set by each validator, can the tx priority is computed from the gas price.
func (ctf CheckTxFeees) checkTxFeeWithValidatorMinGasPrices(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return nil, 0, sdkerrors.ErrTxDecode.Wrap("Tx must be a FeeTx")
	}

	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()

	// Ensure that the provided fees meet a minimum threshold for the validator,
	// if this is a CheckTx. This is only for local mempool purposes, and thus
	// is only ran on check tx.
	if ctx.IsCheckTx() {
		// Begin ==========>
		if ctf.isByPassMinFee(feeTx.GetMsgs(), gas) {
			priority := getTxPriority(feeCoins, int64(gas))
			return feeCoins, priority, nil
		}
		// <========== End
		minGasPrices := ctx.MinGasPrices()
		if !minGasPrices.IsZero() {
			requiredFees := make(sdk.Coins, len(minGasPrices))

			// Determine the required fees by multiplying each required minimum gas
			// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
			glDec := sdkmath.LegacyNewDec(int64(gas))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec)
				requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}

			if !feeCoins.IsAnyGTE(requiredFees) {
				return nil, 0, sdkerrors.ErrInsufficientFee.Wrapf("insufficient fees; got: %s required: %s", feeCoins, requiredFees)
			}
		}
	}

	priority := getTxPriority(feeCoins, int64(gas))
	return feeCoins, priority, nil
}

// getTxPriority returns a naive tx priority based on the amount of the smallest denomination of the gas price
// provided in a transaction.
// NOTE: This implementation should be used with a great consideration as it opens potential attack vectors
// where txs with multiple coins could not be prioritize as expected.
func getTxPriority(fee sdk.Coins, gas int64) int64 {
	var priority int64
	for _, c := range fee {
		p := int64(math.MaxInt64)
		gasPrice := c.Amount.QuoRaw(gas)
		if gasPrice.IsInt64() {
			p = gasPrice.Int64()
		}
		if priority == 0 || p < priority {
			priority = p
		}
	}

	return priority
}
