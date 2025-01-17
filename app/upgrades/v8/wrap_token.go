package v8

import (
	"bytes"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/pundiai/fx-core/v8/types"
	fxevmkeeper "github.com/pundiai/fx-core/v8/x/evm/keeper"
)

const (
	WrapTokenSymbol = "WPUNDIAI"
	// #nosec G101
	WrapTokenName = "Wrapped Pundi AIFX Token"
)

const (
	nameSlot        = 201
	symbolSlot      = 202
	totalSupplySlot = 204
	lastSlot        = 300
)

var (
	// skip eip1967.proxy.implementation 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc
	eip1967ProxyImplKey = common.HexToHash("0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc")
	// skip eip1967.proxy.admin
	eip1967ProxyAdminKey = common.HexToHash("0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103")
	// skip eip1967.proxy.rollback
	eip1967ProxyRollbackKey = common.HexToHash("0x4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd9143")
	// skip eip1967.proxy.beacon
	eip1967ProxyBeaconKey = common.HexToHash("0xa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d50")

	mainnetWFXAddress = common.HexToAddress("0x80b5a32E4F032B2a058b4F29EC95EEfEEB87aDcd")

	testnetWFXAddress = common.HexToAddress("0x3452e23F9c4cC62c70B7ADAd699B264AF3549C19")
)

func migrateWFXToWPUNDIAI(ctx sdk.Context, evmKeeper *fxevmkeeper.Keeper) {
	wfxAddr := GetWFXAddress(ctx.ChainID())

	totalUpdateState := 0
	evmKeeper.ForEachStorage(ctx, wfxAddr, func(key, value common.Hash) bool {
		slot := big.NewInt(0).SetBytes(key.Bytes())
		if slot.Cmp(big.NewInt(lastSlot)) < 0 || isEip1967ProxyKey(key) {
			return true
		}
		evmKeeper.SetState(ctx, wfxAddr, key, scaleDownBy100(value.Bytes()))
		totalUpdateState++
		return true
	})

	ctx.Logger().Info("total update state", "total", totalUpdateState)

	// set totalSupply
	totalSupplyByte := evmKeeper.GetState(ctx, wfxAddr, slotToHashByte(totalSupplySlot))
	evmKeeper.SetState(ctx, wfxAddr, slotToHashByte(totalSupplySlot), scaleDownBy100(totalSupplyByte.Bytes()))

	// test update storage
	nameBytes := encodeShortStringWithoutPrefix(WrapTokenName)
	evmKeeper.SetState(ctx, wfxAddr, slotToHashByte(nameSlot), nameBytes)

	symbolBytes := encodeShortStringWithoutPrefix(WrapTokenSymbol)
	evmKeeper.SetState(ctx, wfxAddr, slotToHashByte(symbolSlot), symbolBytes)
	ctx.Logger().Info("migration WFX to WPUNDIAI done")
}

func GetWFXAddress(chainID string) common.Address {
	if chainID == fxtypes.MainnetChainId {
		return mainnetWFXAddress
	}
	return testnetWFXAddress
}

func isEip1967ProxyKey(key common.Hash) bool {
	return bytes.Equal(key.Bytes(), eip1967ProxyImplKey.Bytes()) ||
		bytes.Equal(key.Bytes(), eip1967ProxyAdminKey.Bytes()) ||
		bytes.Equal(key.Bytes(), eip1967ProxyRollbackKey.Bytes()) ||
		bytes.Equal(key.Bytes(), eip1967ProxyBeaconKey.Bytes())
}

func slotToHashByte(slot uint64) common.Hash {
	return common.BigToHash(big.NewInt(int64(slot)))
}

func scaleDownBy100(valueByte []byte) []byte {
	value := big.NewInt(0).SetBytes(valueByte)
	newValue := value.Div(value, big.NewInt(100))
	return common.BigToHash(newValue).Bytes()
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
