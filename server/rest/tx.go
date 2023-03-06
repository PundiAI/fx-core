package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
)

// ConvertTxToStdTx converts a transaction to the legacy StdTx format
func ConvertTxToStdTx(codec *codec.LegacyAmino, tx signing.Tx) (legacytx.StdTx, error) {
	if stdTx, ok := tx.(legacytx.StdTx); ok {
		return stdTx, nil
	}

	aminoTxConfig := legacytx.StdTxConfig{Cdc: codec}
	builder := aminoTxConfig.NewTxBuilder()

	err := CopyTx(tx, builder, true)
	if err != nil {
		return legacytx.StdTx{}, err
	}

	stdTx, ok := builder.GetTx().(legacytx.StdTx)
	if !ok {
		return legacytx.StdTx{}, fmt.Errorf("expected %T, got %+v", legacytx.StdTx{}, builder.GetTx())
	}

	return stdTx, nil
}

// CopyTx copies a Tx to a new TxBuilder, allowing conversion between
// different transaction formats. If ignoreSignatureError is true, copying will continue
// tx even if the signature cannot be set in the target builder resulting in an unsigned tx.
func CopyTx(tx signing.Tx, builder client.TxBuilder, ignoreSignatureError bool) error {
	err := builder.SetMsgs(tx.GetMsgs()...)
	if err != nil {
		return err
	}

	sigs, err := tx.GetSignaturesV2()
	if err != nil {
		return err
	}

	err = builder.SetSignatures(sigs...)
	if err != nil {
		if ignoreSignatureError {
			// we call SetSignatures() agan with no args to clear any signatures in case the
			// previous call to SetSignatures() had any partial side-effects
			_ = builder.SetSignatures()
		} else {
			return err
		}
	}

	builder.SetMemo(tx.GetMemo())
	builder.SetFeeAmount(tx.GetFee())
	builder.SetGasLimit(tx.GetGas())
	builder.SetTimeoutHeight(tx.GetTimeoutHeight())

	return nil
}

// WriteGeneratedTxResponse writes a generated unsigned transaction to the
// provided http.ResponseWriter. It will simulate gas costs if requested by the
// BaseReq. Upon any error, the error will be written to the http.ResponseWriter.
// Note that this function returns the legacy StdTx Amino JSON format for compatibility
// with legacy clients.
// Deprecated: We are removing Amino soon.
func WriteGeneratedTxResponse(
	clientCtx client.Context, w http.ResponseWriter, br BaseReq, msgs ...sdk.Msg,
) {
	gasAdj, ok := ParseFloat64OrReturnBadRequest(w, br.GasAdjustment, flags.DefaultGasAdjustment)
	if !ok {
		return
	}

	gasSetting, err := flags.ParseGasSetting(br.Gas)
	if CheckBadRequestError(w, err) {
		return
	}

	txf := tx.Factory{}.
		WithFees(br.Fees.String()).
		WithGasPrices(br.GasPrices.String()).
		WithAccountNumber(br.AccountNumber).
		WithSequence(br.Sequence).
		WithGas(gasSetting.Gas).
		WithGasAdjustment(gasAdj).
		WithMemo(br.Memo).
		WithChainID(br.ChainID).
		WithSimulateAndExecute(br.Simulate).
		WithTxConfig(clientCtx.TxConfig).
		WithTimeoutHeight(br.TimeoutHeight)

	if br.Simulate || gasSetting.Simulate {
		if gasAdj < 0 {
			WriteErrorResponse(w, http.StatusBadRequest, errorsmod.ErrorInvalidGasAdjustment.Error())
			return
		}

		_, adjusted, err := tx.CalculateGas(clientCtx, txf, msgs...)
		if CheckInternalServerError(w, err) {
			return
		}

		txf = txf.WithGas(adjusted)

		if br.Simulate {
			WriteSimulationResponse(w, clientCtx.LegacyAmino, txf.Gas())
			return
		}
	}

	tx, err := BuildUnsignedTx(txf, msgs...)
	if CheckBadRequestError(w, err) {
		return
	}

	stdTx, err := ConvertTxToStdTx(clientCtx.LegacyAmino, tx.GetTx())
	if CheckInternalServerError(w, err) {
		return
	}

	output, err := clientCtx.LegacyAmino.MarshalJSON(stdTx)
	if CheckInternalServerError(w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(output)
}

// BuildUnsignedTx builds a transaction to be signed given a set of messages. The
// transaction is initially created via the provided factory's generator. Once
// created, the fee, memo, and messages are set.
func BuildUnsignedTx(txf tx.Factory, msgs ...sdk.Msg) (client.TxBuilder, error) {
	return txf.BuildUnsignedTx(msgs...)
}
