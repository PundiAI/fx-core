package types

import (
	"encoding/hex"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/abi"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

// GetCheckpointOracleSet returns the checkpoint
func GetCheckpointOracleSet(oracleSet *types.OracleSet, gravityIDStr string) ([]byte, error) {
	addresses := make([]string, len(oracleSet.Members))
	powers := make([]*big.Int, len(oracleSet.Members))
	for i, member := range oracleSet.Members {
		addresses[i] = member.ExternalAddress
		powers[i] = big.NewInt(int64(member.Power))
	}

	gravityID, err := fxtypes.StrToByte32(gravityIDStr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse gravity id")
	}
	checkpoint, err := fxtypes.StrToByte32("checkpoint")
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse checkpoint")
	}

	params := []abi.Param{
		{"bytes32": gravityID},
		{"bytes32": checkpoint},
		{"uint256": big.NewInt(int64(oracleSet.Nonce))},
		{"address[]": addresses},
		{"uint256[]": powers},
	}
	encode, err := abi.GetPaddedParam(params)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(encode), nil
}

func GetCheckpointConfirmBatch(txBatch *types.OutgoingTxBatch, gravityIDStr string) ([]byte, error) {
	txCount := len(txBatch.Transactions)
	amounts := make([]*big.Int, txCount)
	destinations := make([]string, txCount)
	fees := make([]*big.Int, txCount)
	for i, transferTx := range txBatch.Transactions {
		amounts[i] = transferTx.Token.Amount.BigInt()
		destinations[i] = transferTx.DestAddress
		fees[i] = transferTx.Fee.Amount.BigInt()
	}

	gravityID, err := fxtypes.StrToByte32(gravityIDStr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse gravity id")
	}
	transactionBatch, err := fxtypes.StrToByte32("transactionBatch")
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse checkpoint")
	}

	params := []abi.Param{
		{"bytes32": gravityID},
		{"bytes32": transactionBatch},
		{"uint256[]": amounts},
		{"address[]": destinations},
		{"uint256[]": fees},
		{"uint256": big.NewInt(int64(txBatch.BatchNonce))},
		{"address": txBatch.TokenContract},
		{"uint256": big.NewInt(int64(txBatch.BatchTimeout))},
		{"address": txBatch.FeeReceive},
	}

	encode, err := abi.GetPaddedParam(params)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(encode), nil
}

func GetCheckpointBridgeCall(bridgeCall *types.OutgoingBridgeCall, gravityIDStr string) ([]byte, error) {
	gravityID, err := fxtypes.StrToByte32(gravityIDStr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse gravity id")
	}
	bridgeCallMethodName, err := fxtypes.StrToByte32("bridgeCall")
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse checkpoint")
	}
	dataBytes, err := hex.DecodeString(bridgeCall.Data)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse data")
	}
	memeBytes, err := hex.DecodeString(bridgeCall.Memo)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse memo")
	}
	contracts := make([]string, 0, len(bridgeCall.Tokens))
	amounts := make([]*big.Int, 0, len(bridgeCall.Tokens))
	for _, token := range bridgeCall.Tokens {
		contracts = append(contracts, token.Contract)
		amounts = append(amounts, token.Amount.BigInt())
	}

	params := []abi.Param{
		{"bytes32": gravityID},
		{"bytes32": bridgeCallMethodName},
		{"address": bridgeCall.Sender},
		{"address": bridgeCall.Receiver},
		{"address[]": contracts},
		{"uint256[]": amounts},
		{"address": bridgeCall.To},
		{"bytes": dataBytes},
		{"bytes": memeBytes},
		{"uint256": big.NewInt(int64(bridgeCall.Nonce))},
		{"uint256": big.NewInt(int64(bridgeCall.Timeout))},
	}

	encode, err := abi.GetPaddedParam(params)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(encode), nil
}
