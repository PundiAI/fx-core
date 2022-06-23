package v020_test

import (
	fxtypes "github.com/functionx/fx-core/types"
	crosschainkeeper "github.com/functionx/fx-core/x/crosschain/keeper"
	polygontypes "github.com/functionx/fx-core/x/polygon/types"
	trontypes "github.com/functionx/fx-core/x/tron/types"
	"sort"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	type args struct {
		moduleName   string
		oracleNumber int
		keeper       func(myApp *app.App) crosschainkeeper.Keeper
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "migrate bsc",
			args: args{
				moduleName:   bsctypes.ModuleName,
				oracleNumber: 10,
				keeper: func(myApp *app.App) crosschainkeeper.Keeper {
					return myApp.BscKeeper
				},
			},
		},
		{
			name: "migrate polygon",
			args: args{
				moduleName:   polygontypes.ModuleName,
				oracleNumber: 20,
				keeper: func(myApp *app.App) crosschainkeeper.Keeper {
					return myApp.PolygonKeeper
				},
			},
		},
		{
			name: "migrate tron",
			args: args{
				moduleName:   trontypes.ModuleName,
				oracleNumber: 50,
				keeper: func(myApp *app.App) crosschainkeeper.Keeper {
					return myApp.TronKeeper.Keeper
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			oracleAddrs := make([]sdk.AccAddress, tt.args.oracleNumber)
			for i := 0; i < len(oracleAddrs); i++ {
				oracleAddrs[i] = secp256k1.GenPrivKey().PubKey().Bytes()
			}
			bridgerAddrs := make([]sdk.AccAddress, tt.args.oracleNumber)
			for i := 0; i < len(bridgerAddrs); i++ {
				bridgerAddrs[i] = secp256k1.GenPrivKey().PubKey().Bytes()
			}
			externalAddrs := make([]common.Address, tt.args.oracleNumber)
			for i := 0; i < len(externalAddrs); i++ {
				key, _ := crypto.GenerateKey()
				externalAddrs[i] = crypto.PubkeyToAddress(key.PublicKey)
			}
			var legacyOracles = make([]*v010.LegacyOracle, len(oracleAddrs))
			for i := 0; i < len(oracleAddrs); i++ {
				legacyOracles[i] = &v010.LegacyOracle{
					OracleAddress:       oracleAddrs[i].String(),
					OrchestratorAddress: bridgerAddrs[i].String(),
					ExternalAddress:     externalAddrs[i].String(),
					DepositAmount:       sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10_000).MulRaw(1e18)),
					StartHeight:         100,
					Jailed:              false,
					JailedHeight:        0,
				}
			}
			sort.Slice(legacyOracles, func(i, j int) bool {
				return legacyOracles[i].OracleAddress > legacyOracles[j].OracleAddress
			})

			valSet, genAccs, balances := helpers.GenerateGenesisValidator(50, nil)
			myApp := helpers.SetupWithGenesisValSet(t, valSet, genAccs, balances...)
			ctx := myApp.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})

			storeKey := myApp.GetKey(tt.args.moduleName)
			store := ctx.KVStore(storeKey)
			for i := 0; i < len(oracleAddrs); i++ {
				store.Set(append(types.OracleKey, oracleAddrs[i].Bytes()...), myApp.AppCodec().MustMarshal(legacyOracles[i]))
			}

			oracles, validator, err := v020.MigrateOracle(ctx, myApp.AppCodec(), myApp.GetKey(tt.args.moduleName), myApp.StakingKeeper)
			require.NoError(t, err)

			require.Equal(t, len(oracles), len(legacyOracles))
			sort.Slice(oracles, func(i, j int) bool {
				return oracles[i].OracleAddress > oracles[j].OracleAddress
			})

			oraclesFromDB := tt.args.keeper(myApp).GetAllOracles(ctx, false)
			sort.Slice(oraclesFromDB, func(i, j int) bool {
				return oraclesFromDB[i].OracleAddress > oraclesFromDB[j].OracleAddress
			})
			require.Equal(t, len(oracles), len(oraclesFromDB))

			for i := 0; i < len(oracles); i++ {
				require.Equal(t, tt.args.keeper(myApp).HasOracle(ctx, oracles[i].GetOracle()), true)

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
				oracleFromDB, found := tt.args.keeper(myApp).GetOracle(ctx, oracles[i].GetOracle())
				require.Equal(t, found, true)
				require.EqualValues(t, oracleFromDB, newOracle)

				require.EqualValues(t, oracles[i], newOracle)
				require.EqualValues(t, oraclesFromDB[i], newOracle)
			}
		})
	}
}
