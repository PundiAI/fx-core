package ante

import (
	"bytes"
	"fmt"
	"math"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	tmstrings "github.com/tendermint/tendermint/libs/strings"

	fxstakingtypes "github.com/functionx/fx-core/v6/x/staking/types"
)

// DeductFeeDecorator deducts fees from the first signer of the tx
// If the first signer does not have the funds to pay for the fees, return with InsufficientFunds error
// Call next AnteHandler if fees successfully deducted
// CONTRACT: Tx must implement FeeTx interface to use DeductFeeDecorator
type DeductFeeDecorator struct {
	accountKeeper              AccountKeeper
	bankKeeper                 types.BankKeeper
	feegrantKeeper             FeegrantKeeper
	txFeeChecker               ante.TxFeeChecker
	BypassMinFeeMsgTypes       []string
	MaxBypassMinFeeMsgGasUsage uint64
}

func NewDeductFeeDecorator(ak AccountKeeper, bk types.BankKeeper, fk FeegrantKeeper, tfc ante.TxFeeChecker, bypassMsgTypes []string, maxBypassMinFeeMsgGasUsage uint64) DeductFeeDecorator {
	if tfc == nil {
		tfc = checkTxFeeWithValidatorMinGasPrices
	}

	return DeductFeeDecorator{
		accountKeeper:              ak,
		bankKeeper:                 bk,
		feegrantKeeper:             fk,
		txFeeChecker:               tfc,
		BypassMinFeeMsgTypes:       bypassMsgTypes,
		MaxBypassMinFeeMsgGasUsage: maxBypassMinFeeMsgGasUsage,
	}
}

func (dfd DeductFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(errortypes.ErrTxDecode, "Tx must be a FeeTx")
	}

	if !simulate && ctx.BlockHeight() > 0 && feeTx.GetGas() == 0 {
		return ctx, errorsmod.Wrap(errortypes.ErrInvalidGasLimit, "must provide positive gas")
	}

	var (
		priority int64
		err      error
	)

	fee := feeTx.GetFee()
	gas := feeTx.GetGas()
	msgs := feeTx.GetMsgs()
	if ctx.IsCheckTx() && !simulate && !(dfd.bypassMinFeeMsgs(msgs) && gas <= uint64(len(msgs))*dfd.MaxBypassMinFeeMsgGasUsage) {
		fee, priority, err = dfd.txFeeChecker(ctx, tx)
		if err != nil {
			return ctx, err
		}
	}
	if err = dfd.checkDeductFee(ctx, tx, fee); err != nil {
		return ctx, err
	}

	newCtx := ctx.WithPriority(priority)

	return next(newCtx, tx, simulate)
}

func (dfd DeductFeeDecorator) checkDeductFee(ctx sdk.Context, sdkTx sdk.Tx, fee sdk.Coins) error {
	feeTx, ok := sdkTx.(sdk.FeeTx)
	if !ok {
		return errorsmod.Wrap(errortypes.ErrTxDecode, "Tx must be a FeeTx")
	}

	if addr := dfd.accountKeeper.GetModuleAddress(types.FeeCollectorName); addr == nil {
		return fmt.Errorf("fee collector module account (%s) has not been set", types.FeeCollectorName)
	}

	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()
	deductFeesFrom := feePayer

	// if feegranter set deduct fee from feegranter account.
	// this works with only when feegrant enabled.
	if feeGranter != nil {
		if dfd.feegrantKeeper == nil {
			return errortypes.ErrInvalidRequest.Wrap("fee grants are not enabled")
		} else if !feeGranter.Equals(feePayer) {
			err := dfd.feegrantKeeper.UseGrantedFees(ctx, feeGranter, feePayer, fee, sdkTx.GetMsgs())
			if err != nil {
				return errorsmod.Wrapf(err, "%s does not not allow to pay fees for %s", feeGranter, feePayer)
			}
		}

		deductFeesFrom = feeGranter
	}

	deductFeesFromAcc := dfd.accountKeeper.GetAccount(ctx, deductFeesFrom)
	if deductFeesFromAcc == nil {
		return errortypes.ErrUnknownAddress.Wrapf("fee payer address: %s does not exist", deductFeesFrom)
	}

	// deduct the fees
	if !fee.IsZero() {
		err := ante.DeductFees(dfd.bankKeeper, ctx, deductFeesFromAcc, fee)
		if err != nil {
			return err
		}
	}

	events := sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeTx,
			sdk.NewAttribute(sdk.AttributeKeyFee, fee.String()),
			sdk.NewAttribute(sdk.AttributeKeyFeePayer, deductFeesFrom.String()),
		),
	}
	ctx.EventManager().EmitEvents(events)

	return nil
}

func (dfd DeductFeeDecorator) bypassMinFeeMsgs(msgs []sdk.Msg) bool {
	for _, msg := range msgs {
		if tmstrings.StringInSlice(sdk.MsgTypeURL(msg), dfd.BypassMinFeeMsgTypes) {
			continue
		}

		return false
	}

	return true
}

// checkTxFeeWithValidatorMinGasPrices implements the default fee logic, where the minimum price per
// unit of gas is fixed and set by each validator, can the tx priority is computed from the gas price.
func checkTxFeeWithValidatorMinGasPrices(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return nil, 0, errorsmod.Wrap(errortypes.ErrTxDecode, "Tx must be a FeeTx")
	}

	feeCoins := feeTx.GetFee()
	gas := feeTx.GetGas()

	// Ensure that the provided fees meet a minimum threshold for the validator,
	// if this is a CheckTx. This is only for local mempool purposes, and thus
	// is only ran on check tx.
	if ctx.IsCheckTx() {
		minGasPrices := ctx.MinGasPrices()
		if !minGasPrices.IsZero() {
			requiredFees := make(sdk.Coins, len(minGasPrices))

			// Determine the required fees by multiplying each required minimum gas
			// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
			glDec := sdk.NewDec(int64(gas))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec)
				requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}

			if !feeCoins.IsAnyGTE(requiredFees) {
				return nil, 0, errorsmod.Wrapf(errortypes.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeCoins, requiredFees)
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

// EthMinGasPriceDecorator will check if the transaction's fee is at least as large
// as the MinGasPrices param. If fee is too low, decorator returns error and tx
// is rejected. This applies to both CheckTx and DeliverTx and regardless
// if London hard fork or fee market params (EIP-1559) are enabled.
// If fee is high enough, then call next AnteHandler
type EthMinGasPriceDecorator struct {
	feesKeeper FeeMarketKeeper
	evmKeeper  EVMKeeper
}

// NewEthMinGasPriceDecorator creates a new MinGasPriceDecorator instance used only for
// Ethereum transactions.
func NewEthMinGasPriceDecorator(fk FeeMarketKeeper, ek EVMKeeper) EthMinGasPriceDecorator {
	return EthMinGasPriceDecorator{feesKeeper: fk, evmKeeper: ek}
}

// AnteHandle ensures that the that the effective fee from the transaction is greater than the
// minimum global fee, which is defined by the  MinGasPrice (parameter) * GasLimit (tx argument).
func (empd EthMinGasPriceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	minGasPrice := empd.feesKeeper.GetParams(ctx).MinGasPrice

	// short-circuit if min gas price is 0
	if minGasPrice.IsZero() {
		return next(ctx, tx, simulate)
	}

	evmParams := empd.evmKeeper.GetParams(ctx)
	chainCfg := evmParams.GetChainConfig()
	ethCfg := chainCfg.EthereumConfig(empd.evmKeeper.ChainID())
	baseFee := empd.evmKeeper.GetBaseFee(ctx, ethCfg)

	for _, msg := range tx.GetMsgs() {
		ethMsg, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return ctx, errorsmod.Wrapf(
				errortypes.ErrUnknownRequest,
				"invalid message type %T, expected %T",
				msg, (*evmtypes.MsgEthereumTx)(nil),
			)
		}

		feeAmt := ethMsg.GetFee()

		// For dynamic transactions, GetFee() uses the GasFeeCap value, which
		// is the maximum gas price that the signer can pay. In practice, the
		// signer can pay less, if the block's BaseFee is lower. So, in this case,
		// we use the EffectiveFee. If the feemarket formula results in a BaseFee
		// that lowers EffectivePrice until it is < MinGasPrices, the users must
		// increase the GasTipCap (priority fee) until EffectivePrice > MinGasPrices.
		// Transactions with MinGasPrices * gasUsed < tx fees < EffectiveFee are rejected
		// by the feemarket AnteHandle

		txData, err := evmtypes.UnpackTxData(ethMsg.Data)
		if err != nil {
			return ctx, errorsmod.Wrapf(err, "failed to unpack tx data %s", ethMsg.Hash)
		}

		if txData.TxType() != ethtypes.LegacyTxType {
			feeAmt = ethMsg.GetEffectiveFee(baseFee)
		}

		gasLimit := sdk.NewDecFromBigInt(new(big.Int).SetUint64(ethMsg.GetGas()))

		requiredFee := minGasPrice.Mul(gasLimit)
		fee := sdk.NewDecFromBigInt(feeAmt)

		if fee.LT(requiredFee) {
			return ctx, errorsmod.Wrapf(
				errortypes.ErrInsufficientFee,
				"provided fee < minimum global fee (%d < %d). Please increase the priority tip (for EIP-1559 txs) or the gas prices (for access list or legacy txs)", //nolint:lll
				fee.TruncateInt().Int64(), requiredFee.TruncateInt().Int64(),
			)
		}
	}

	return next(ctx, tx, simulate)
}

func DeductTxCostsFromUserBalance(ctx sdk.Context, ak AccountKeeper, bk types.BankKeeper, fees sdk.Coins, from common.Address) error {
	// fetch sender account
	signerAcc, err := ante.GetSignerAcc(ctx, ak, from.Bytes())
	if err != nil {
		return errorsmod.Wrapf(err, "account not found for sender %s", from)
	}

	// check if the account is disabled
	if signerAcc.GetPubKey() != nil && bytes.Equal(signerAcc.GetPubKey().Bytes(), fxstakingtypes.DisablePKBytes[:]) {
		return errorsmod.Wrap(errortypes.ErrInvalidAddress, "account disabled")
	}

	// deduct the full gas cost from the user balance
	if err := ante.DeductFees(bk, ctx, signerAcc, fees); err != nil {
		return errorsmod.Wrapf(err, "failed to deduct full gas cost %s from the user %s balance", fees, from)
	}
	return nil
}
