package v3

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

func MigrateParams(ctx sdk.Context, legacyAmino *codec.LegacyAmino, paramsKey sdk.StoreKey) {
	paramStoreKeyRejectUnprotectedTx := []byte("RejectUnprotectedTx")

	paramsStore := prefix.NewStore(ctx.KVStore(paramsKey), append([]byte(evmtypes.ModuleName), '/'))
	bzR := paramsStore.Get(paramStoreKeyRejectUnprotectedTx)

	var rejectUnprotectedTx bool
	if err := legacyAmino.UnmarshalJSON(bzR, &rejectUnprotectedTx); err != nil {
		panic(err.Error())
	}

	allowUnprotectedTxs := !rejectUnprotectedTx
	bzA, err := legacyAmino.MarshalJSON(allowUnprotectedTxs)
	if err != nil {
		panic(err.Error())
	}

	ctx.Logger().Info("migrate params", "module", evmtypes.ModuleName,
		"from", fmt.Sprintf("%s:%v", paramStoreKeyRejectUnprotectedTx, rejectUnprotectedTx),
		"to", fmt.Sprintf("%s:%v", evmtypes.ParamStoreKeyAllowUnprotectedTxs, allowUnprotectedTxs))

	paramsStore.Delete(paramStoreKeyRejectUnprotectedTx)
	paramsStore.Set(evmtypes.ParamStoreKeyAllowUnprotectedTxs, bzA)
}
