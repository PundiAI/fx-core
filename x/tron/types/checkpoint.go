package types

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/abi"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
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
		return nil, fmt.Errorf("parse gravity id: %w", err)
	}
	checkpoint, err := fxtypes.StrToByte32("checkpoint")
	if err != nil {
		return nil, fmt.Errorf("parse checkpoint: %w", err)
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
		return nil, fmt.Errorf("parse gravity id: %w", err)
	}
	transactionBatch, err := fxtypes.StrToByte32("transactionBatch")
	if err != nil {
		return nil, fmt.Errorf("parse transaction batch: %w", err)
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
		return nil, fmt.Errorf("parse gravity id: %w", err)
	}
	bridgeCallMethodName, err := fxtypes.StrToByte32("bridgeCall")
	if err != nil {
		return nil, fmt.Errorf("parse bridge call method name: %w", err)
	}
	dataBytes, err := hex.DecodeString(bridgeCall.Data)
	if err != nil {
		return nil, fmt.Errorf("parse data: %w", err)
	}
	memeBytes, err := hex.DecodeString(bridgeCall.Memo)
	if err != nil {
		return nil, fmt.Errorf("parse memo: %w", err)
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
		{"address": bridgeCall.Refund},
		{"address[]": contracts},
		{"uint256[]": amounts},
		{"address": bridgeCall.To},
		{"bytes": dataBytes},
		{"bytes": memeBytes},
		{"uint256": big.NewInt(int64(bridgeCall.Nonce))},
		{"uint256": big.NewInt(int64(bridgeCall.Timeout))},
		{"uint256": big.NewInt(int64(bridgeCall.GasLimit))},
		{"uint256": big.NewInt(int64(bridgeCall.EventNonce))},
	}

	encode, err := abi.GetPaddedParam(params)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(encode), nil
}
