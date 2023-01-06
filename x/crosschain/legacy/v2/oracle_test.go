package v2_test

import (
	"sort"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
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
			var legacyOracles = make([]*crosschainv1.LegacyOracle, len(oracleAddrs))
			for i := 0; i < len(oracleAddrs); i++ {
				legacyOracles[i] = &crosschainv1.LegacyOracle{
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
			ctx := myApp.NewContext(false, tmproto.Header{Time: time.Now()})

			storeKey := myApp.GetKey(tt.args.moduleName)
			store := ctx.KVStore(storeKey)
			for i := 0; i < len(oracleAddrs); i++ {
				store.Set(append(types.OracleKey, oracleAddrs[i].Bytes()...), myApp.AppCodec().MustMarshal(legacyOracles[i]))
			}

			oracles, validatorAddr, err := crosschainv2.MigrateOracle(ctx, myApp.AppCodec(), myApp.GetKey(tt.args.moduleName), myApp.StakingKeeper)
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
					DelegateValidator: validatorAddr.String(),
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
