package v020_test

import (
	fxtypes "github.com/functionx/fx-core/types"
	"sort"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/app/helpers"
	bsctypes "github.com/functionx/fx-core/x/bsc/types"
	v010 "github.com/functionx/fx-core/x/crosschain/legacy/v010"
	v020 "github.com/functionx/fx-core/x/crosschain/legacy/v020"
	"github.com/functionx/fx-core/x/crosschain/types"
)

func TestMigrateOracle(t *testing.T) {
	oracleAddr := make([]sdk.AccAddress, 20)
	for i := 0; i < len(oracleAddr); i++ {
		oracleAddr[i] = secp256k1.GenPrivKey().PubKey().Bytes()
	}
	bridgerAddr := make([]sdk.AccAddress, 20)
	for i := 0; i < len(oracleAddr); i++ {
		bridgerAddr[i] = secp256k1.GenPrivKey().PubKey().Bytes()
	}
	externalAddr := make([]common.Address, 20)
	for i := 0; i < len(oracleAddr); i++ {
		key, _ := crypto.GenerateKey()
		externalAddr[i] = crypto.PubkeyToAddress(key.PublicKey)
	}
	var legacyOracles = make([]*v010.LegacyOracle, len(oracleAddr))
	for i := 0; i < len(oracleAddr); i++ {
		legacyOracles[i] = &v010.LegacyOracle{
			OracleAddress:       oracleAddr[i].String(),
			OrchestratorAddress: bridgerAddr[i].String(),
			ExternalAddress:     externalAddr[i].String(),
			DepositAmount:       sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10_000).MulRaw(1e18)),
			StartHeight:         100,
			Jailed:              false,
			JailedHeight:        0,
		}
	}
	sort.Slice(legacyOracles, func(i, j int) bool {
		return legacyOracles[i].OracleAddress > legacyOracles[j].OracleAddress
	})
	type args struct {
		moduleName string
		prepare    func(*app.App) error
		require    func(*testing.T, *app.App, types.Oracles, stakingtypes.Validator)
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
					ctx := myApp.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})
					storeKey := myApp.GetKey(bsctypes.ModuleName)
					store := ctx.KVStore(storeKey)
					for i := 0; i < len(oracleAddr); i++ {
						store.Set(append(types.OracleKey, oracleAddr[i].Bytes()...), myApp.AppCodec().MustMarshal(legacyOracles[i]))
					}
					return nil
				},
				require: func(t *testing.T, myApp *app.App, oracles types.Oracles, validator stakingtypes.Validator) {
					require.Equal(t, len(oracles), len(legacyOracles))
					sort.Slice(oracles, func(i, j int) bool {
						return oracles[i].OracleAddress > oracles[j].OracleAddress
					})
					ctx := myApp.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})
					oraclesFromDB := myApp.BscKeeper.GetAllOracles(ctx, false)
					sort.Slice(oraclesFromDB, func(i, j int) bool {
						return oraclesFromDB[i].OracleAddress > oracles[j].OracleAddress
					})
					require.Equal(t, len(oracles), len(oraclesFromDB))

					for i := 0; i < len(oracles); i++ {
						require.Equal(t, myApp.BscKeeper.HasOracle(ctx, oracles[i].GetOracle()), true)

						newOracle := types.Oracle{
							OracleAddress:     legacyOracles[i].OracleAddress,
							BridgerAddress:    legacyOracles[i].OrchestratorAddress,
							ExternalAddress:   legacyOracles[i].ExternalAddress,
							DelegateAmount:    legacyOracles[i].DepositAmount.Amount,
							StartHeight:       legacyOracles[i].StartHeight,
							Online:            legacyOracles[i].Jailed == false,
							DelegateValidator: validator.OperatorAddress,
							SlashTimes:        0,
						}
						oracleFromDB, found := myApp.BscKeeper.GetOracle(ctx, oracles[i].GetOracle())
						require.Equal(t, found, true)
						require.EqualValues(t, oracleFromDB, newOracle)

						require.EqualValues(t, oracles[i], newOracle)
						require.EqualValues(t, oraclesFromDB[i], newOracle)
					}
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valSet, genAccs, balances := helpers.GenerateGenesisValidator(50, nil)
			myApp := helpers.SetupWithGenesisValSet(t, valSet, genAccs, balances...)
			require.NoError(t, tt.args.prepare(myApp))

			ctx := myApp.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})
			oracle, validator, err := v020.MigrateOracle(ctx, myApp.AppCodec(), myApp.GetKey(tt.args.moduleName), myApp.StakingKeeper)
			require.NoError(t, err)

			tt.args.require(t, myApp, oracle, validator)
		})
	}
}
