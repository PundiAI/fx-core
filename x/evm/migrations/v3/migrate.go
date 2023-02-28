package v3

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

func MigrateParams(ctx sdk.Context, legacyAmino *codec.LegacyAmino, paramsKey storetypes.StoreKey) {
	paramStoreKeyRejectUnprotectedTx := []byte("RejectUnprotectedTx")

	paramsStore := prefix.NewStore(ctx.KVStore(paramsKey), append([]byte(evmtypes.ModuleName), '/'))

	// NOTE: Default Allow Unprotected txs is false
	allowUnprotectedTxs := evmtypes.DefaultAllowUnprotectedTxs
	bzA, err := legacyAmino.MarshalJSON(allowUnprotectedTxs)
	if err != nil {
		panic(err.Error())
	}

	ctx.Logger().Info("migrating evm module params", "module", "evm",
		"delete", string(paramStoreKeyRejectUnprotectedTx),
		"add", fmt.Sprintf("%s:%v", evmtypes.ParamStoreKeyAllowUnprotectedTxs, allowUnprotectedTxs))

	paramsStore.Delete(paramStoreKeyRejectUnprotectedTx)
	paramsStore.Set(evmtypes.ParamStoreKeyAllowUnprotectedTxs, bzA)
}
