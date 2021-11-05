package types

import (
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type OutgoingTransferTxs []*OutgoingTransferTx

func (bs OutgoingTransferTxs) TotalFee() sdk.Int {
	totalFee := sdk.NewInt(0)
	for _, tx := range bs {
		totalFee = totalFee.Add(tx.Fee.Amount)
	}
	return totalFee
}

// GetCheckpoint gets the checkpoint signature from the given outgoing tx batch
func (m OutgoingTxBatch) GetCheckpoint(gravityIDString string) ([]byte, error) {

	abiObj, err := abi.JSON(strings.NewReader(OutgoingBatchTxCheckpointABIJSON))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "bad ABI definition in code")
	}

	// the contract argument is not a arbitrary length array but a fixed length 32 byte
	// array, therefore we have to utf8 encode the string (the default in this case) and
	// then copy the variable length encoded data into a fixed length array. This function
	// will panic if gravityId is too long to fit in 32 bytes
	gravityID, err := StrToFixByteArray(gravityIDString)
	if err != nil {
		return nil, sdkerrors.Wrap(ErrInvalid, "gravity id")
	}

	// Create the methodName argument which salts the signature
	methodNameBytes := []uint8("transactionBatch")
	var batchMethodName [32]uint8
	copy(batchMethodName[:], methodNameBytes[:])

	// Run through the elements of the batch and serialize them
	txAmounts := make([]*big.Int, len(m.Transactions))
	txDestinations := make([]gethcommon.Address, len(m.Transactions))
	txFees := make([]*big.Int, len(m.Transactions))
	for i, tx := range m.Transactions {
		txAmounts[i] = tx.Token.Amount.BigInt()
		txDestinations[i] = gethcommon.HexToAddress(tx.DestAddress)
		txFees[i] = tx.Fee.Amount.BigInt()
	}

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	abiEncodedBatch, err := abiObj.Pack("submitBatch",
		gravityID,
		batchMethodName,
		txAmounts,
		txDestinations,
		txFees,
		big.NewInt(int64(m.BatchNonce)),
		gethcommon.HexToAddress(m.TokenContract),
		big.NewInt(int64(m.BatchTimeout)),
		gethcommon.HexToAddress(m.FeeReceive),
	)

	// this should never happen outside of test since any case that could crash on encoding
	// should be filtered above.
	if err != nil {
		return nil, sdkerrors.Wrap(err, "packing checkpoint")
	}

	// we hash the resulting encoded bytes discarding the first 4 bytes these 4 bytes are the constant
	// method name 'checkpoint'. If you where to replace the checkpoint constant in this code you would
	// then need to adjust how many bytes you truncate off the front to get the output of abi.encode()
	return crypto.Keccak256Hash(abiEncodedBatch[4:]).Bytes(), nil
}
