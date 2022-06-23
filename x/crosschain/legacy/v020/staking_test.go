package v020_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/app/helpers"
	fxtypes "github.com/functionx/fx-core/types"
	bsctypes "github.com/functionx/fx-core/x/bsc/types"
	v020 "github.com/functionx/fx-core/x/crosschain/legacy/v020"
	"github.com/functionx/fx-core/x/crosschain/types"
	polygontypes "github.com/functionx/fx-core/x/polygon/types"
	trontypes "github.com/functionx/fx-core/x/tron/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
	"time"
)

func TestMigrateDepositToStaking(t *testing.T) {
	newOracles := func(number int, validatorAddr string) types.Oracles {
		var oracles = make(types.Oracles, number)
		for i := 0; i < len(oracles); i++ {
			oracles[i] = types.Oracle{
				OracleAddress:     sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes()).String(),
				BridgerAddress:    sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes()).String(),
				ExternalAddress:   sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes()).String(),
				DelegateAmount:    sdk.NewInt(10_000).MulRaw(1e18),
				StartHeight:       100,
				Online:            number%2 == 0,
				DelegateValidator: validatorAddr,
				SlashTimes:        0,
			}
		}
		return oracles
	}
	type args struct {
		moduleName   string
		oracleNumber int
		require      func(*testing.T, *app.App)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "migrate bsc",
			args: args{
				moduleName:   bsctypes.ModuleName,
				oracleNumber: 30,
				require: func(t *testing.T, myApp *app.App) {

				},
			},
		},
		{
			name: "migrate polygon",
			args: args{
				moduleName:   polygontypes.ModuleName,
				oracleNumber: 20,
				require: func(t *testing.T, myApp *app.App) {

				},
			},
		},
		{
			name: "migrate bsc",
			args: args{
				moduleName:   trontypes.ModuleName,
				oracleNumber: 10,
				require: func(t *testing.T, myApp *app.App) {

				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valSet, genAccs, balances := helpers.GenerateGenesisValidator(50, nil)
			myApp := helpers.SetupWithGenesisValSet(t, valSet, genAccs, balances...)

			ctx := myApp.BaseApp.NewContext(false, tmproto.Header{Time: time.Now()})
			validators := myApp.StakingKeeper.GetBondedValidatorsByPower(ctx)
			oracles := newOracles(tt.args.oracleNumber, validators[0].OperatorAddress)
			for _, oracle := range oracles {
				err := myApp.BankKeeper.MintCoins(ctx, tt.args.moduleName, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, oracle.DelegateAmount)))
				require.NoError(t, err)
			}

			err := v020.MigrateDepositToStaking(ctx, tt.args.moduleName, myApp.StakingKeeper, myApp.BankKeeper, oracles, validators[0])
			require.NoError(t, err)

			tt.args.require(t, myApp)
		})
	}
}
