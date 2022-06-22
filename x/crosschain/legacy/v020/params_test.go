package v020_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/app/helpers"
	fxtypes "github.com/functionx/fx-core/types"
	bsctypes "github.com/functionx/fx-core/x/bsc/types"
	v010 "github.com/functionx/fx-core/x/crosschain/legacy/v010"
	v020 "github.com/functionx/fx-core/x/crosschain/legacy/v020"
	"github.com/functionx/fx-core/x/crosschain/types"
)

func TestMigrateParams(t *testing.T) {
	type args struct {
		moduleName string
		prepare    func(myApp *app.App) error
		require    func(t *testing.T, myApp *app.App)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "migrate bsc",
			args: args{
				moduleName: bsctypes.ModuleName,
				prepare: func(myApp *app.App) error {
					genesisState := helpers.DefGenesisState(myApp.AppCodec())
					genesisState[bsctypes.ModuleName] = json.RawMessage("{}")
					stateBytes, err := json.MarshalIndent(genesisState, "", " ")
					if err != nil {
						return err
					}
					myApp.InitChain(abci.RequestInitChain{
						Validators:      []abci.ValidatorUpdate{},
						ConsensusParams: helpers.DefaultConsensusParams,
						AppStateBytes:   stateBytes,
					})
					ctx := myApp.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

					paramsKey := myApp.GetKey(paramstypes.ModuleName)
					var bscParams = types.Params{
						GravityId:                         "fx-bsc-bridge",
						AverageBlockTime:                  5_000,
						ExternalBatchTimeout:              43_200_000,
						AverageExternalBlockTime:          3_000,
						SignedWindow:                      20_000,
						SlashFraction:                     sdk.MustNewDecFromStr("0.01"),
						OracleSetUpdatePowerChangePercent: sdk.MustNewDecFromStr("0.1"),
						IbcTransferTimeoutHeight:          20_000,
						DelegateThreshold:                 sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10_000).MulRaw(1e18)),
						DelegateMultiple:                  0,
					}
					bscParamSetPairs := v010.GetParamSetPairs(&bscParams)
					bscParamsStore := prefix.NewStore(ctx.KVStore(paramsKey), append([]byte(bsctypes.StoreKey), '/'))
					for _, pair := range bscParamSetPairs {
						bscParamsStore.Set(pair.Key, myApp.LegacyAmino().MustMarshalJSON(pair.Value))
					}
					return nil
				},
				require: func(t *testing.T, myApp *app.App) {
					ctx := myApp.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})
					params := myApp.BscKeeper.GetParams(ctx)
					require.NoError(t, params.ValidateBasic())

					require.EqualValues(t, params.AverageBlockTime, 7_000)
					require.Equal(t, params.DelegateThreshold, sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10_000).MulRaw(1e18)))
					require.EqualValues(t, params.DelegateMultiple, types.DefaultOracleDelegateThreshold)

					defParams := bsctypes.DefaultGenesisState().Params
					defParams.AverageBlockTime = 7_000
					defParams.AverageExternalBlockTime = 3_000
					require.EqualValues(t, &params, defParams)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			myApp := helpers.Setup(true, true)
			require.NoError(t, tt.args.prepare(myApp))

			ctx := myApp.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})
			err := v020.MigrateParams(ctx, tt.args.moduleName, myApp.LegacyAmino(), myApp.GetKey(paramstypes.ModuleName))
			require.NoError(t, err)

			tt.args.require(t, myApp)
		})
	}
}
