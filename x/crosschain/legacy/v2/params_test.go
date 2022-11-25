package v2_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	crosschainkeeper "github.com/functionx/fx-core/v3/x/crosschain/keeper"
	crosschainv1 "github.com/functionx/fx-core/v3/x/crosschain/legacy/v1"
	crosschainv2 "github.com/functionx/fx-core/v3/x/crosschain/legacy/v2"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func TestMigrateParams(t *testing.T) {
	type args struct {
		moduleName string
		keeper     func(myApp *app.App) crosschainkeeper.Keeper
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "migrate bsc",
			args: args{
				moduleName: bsctypes.ModuleName,
				keeper: func(myApp *app.App) crosschainkeeper.Keeper {
					return myApp.BscKeeper
				},
			},
		},
		{
			name: "migrate polygon",
			args: args{
				moduleName: polygontypes.ModuleName,
				keeper: func(myApp *app.App) crosschainkeeper.Keeper {
					return myApp.PolygonKeeper
				},
			},
		},
		{
			name: "migrate tron",
			args: args{
				moduleName: trontypes.ModuleName,
				keeper: func(myApp *app.App) crosschainkeeper.Keeper {
					return myApp.TronKeeper.Keeper
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			myApp := helpers.Setup(true, true)

			genesisState := helpers.DefGenesisState(myApp.AppCodec())
			genesisState[tt.args.moduleName] = myApp.AppCodec().MustMarshalJSON(&types.GenesisState{
				Params: types.Params{
					GravityId:                         fmt.Sprintf("fx-%s-bridge", tt.args.moduleName),
					AverageBlockTime:                  5_000,
					ExternalBatchTimeout:              43_200_000,
					AverageExternalBlockTime:          3_000,
					SignedWindow:                      20_000,
					SlashFraction:                     sdk.MustNewDecFromStr("0.01"),
					OracleSetUpdatePowerChangePercent: sdk.MustNewDecFromStr("0.1"),
					IbcTransferTimeoutHeight:          20_000,
					DelegateThreshold:                 sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10_000).MulRaw(1e18)),
					DelegateMultiple:                  10,
				},
			})
			stateBytes, err := json.MarshalIndent(genesisState, "", " ")
			require.NoError(t, err)

			myApp.InitChain(abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: helpers.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			})
			ctx := myApp.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

			paramsKey := myApp.GetKey(paramstypes.ModuleName)
			paramsStore := prefix.NewStore(ctx.KVStore(paramsKey), append([]byte(tt.args.moduleName), '/'))
			paramsStore.Set(crosschainv1.ParamStoreOracles, myApp.LegacyAmino().MustMarshalJSON([]string{sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes()).String()}))
			paramsStore.Set(crosschainv1.ParamOracleDepositThreshold, myApp.LegacyAmino().MustMarshalJSON(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10_000).MulRaw(1e18))))

			require.True(t, crosschainv2.CheckInitialize(ctx, tt.args.moduleName, myApp.GetKey(paramstypes.ModuleName)))
			require.NoError(t, crosschainv2.MigrateParams(ctx, tt.args.moduleName, myApp.LegacyAmino(), myApp.GetKey(paramstypes.ModuleName)))

			iterator := paramsStore.Iterator(nil, nil)
			for ; iterator.Valid(); iterator.Next() {
				require.NotEqual(t, iterator.Key(), crosschainv1.ParamOracleDepositThreshold)
				require.NotEqual(t, iterator.Key(), crosschainv1.ParamStoreOracles)
			}
			require.NoError(t, iterator.Close())

			paramsFromDB := tt.args.keeper(myApp).GetParams(ctx)
			require.NoError(t, paramsFromDB.ValidateBasic())

			require.EqualValues(t, paramsFromDB.AverageBlockTime, 7_000)
			require.Equal(t, paramsFromDB.DelegateThreshold, sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10_000).MulRaw(1e18)))
			require.EqualValues(t, paramsFromDB.DelegateMultiple, types.DefaultOracleDelegateThreshold)

			defParams := types.DefaultParams()
			defParams.GravityId = fmt.Sprintf("fx-%s-bridge", tt.args.moduleName)
			defParams.AverageBlockTime = 7_000
			defParams.AverageExternalBlockTime = 3_000
			require.EqualValues(t, paramsFromDB, defParams)
		})
	}
}
