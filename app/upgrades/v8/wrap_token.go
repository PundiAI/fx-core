package v8

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	fxevmkeeper "github.com/pundiai/fx-core/v8/x/evm/keeper"
)

const (
	WrapName = "Wrapped Pundi AI"
	nameSlot = 201
)

var (
	mainnetWFXAddress = common.HexToAddress("0x80b5a32E4F032B2a058b4F29EC95EEfEEB87aDcd")

	testnetWFXAddress = common.HexToAddress("0x3452e23F9c4cC62c70B7ADAd699B264AF3549C19")
)

func renameWPUNDIAI(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper) {
	wPUNDIAIAddr := GetWFXAddress(ctx.ChainID())

	//  rename wpundiai name
	nameBytes := encodeShortStringWithoutPrefix(WrapName)
	evmKeeper.SetState(ctx, wPUNDIAIAddr, slotToHashByte(nameSlot), nameBytes)

	ctx.Logger().Info("migration rename WPUNDIAI done", "module", "upgrade")
}

func GetWFXAddress(chainID string) common.Address {
	if chainID == fxtypes.MainnetChainId {
		return mainnetWFXAddress
	}
	return testnetWFXAddress
}

func slotToHashByte(slot uint64) common.Hash {
	return common.BigToHash(big.NewInt(int64(slot)))
}

func encodeShortStringWithoutPrefix(s string) []byte {
	if len(s) > 31 {
		panic("String exceeds 31 bytes, requires dynamic encoding!")
	}
	data := []byte(s)
	dataLen := len(data) * 2
	padding := make([]byte, 31-len(data))
	// add padding to the end
	data = append(data, padding...)
	// add length to last byte (add suffix length)
	data = append(data, byte(dataLen))
	return data
}
