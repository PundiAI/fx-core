package v020_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/app/helpers"
	fxtypes "github.com/functionx/fx-core/types"
	bsctypes "github.com/functionx/fx-core/x/bsc/types"
	v020 "github.com/functionx/fx-core/x/crosschain/legacy/v020"
	"github.com/functionx/fx-core/x/crosschain/types"
	polygontypes "github.com/functionx/fx-core/x/polygon/types"
	trontypes "github.com/functionx/fx-core/x/tron/types"
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
			},
		},
		{
			name: "migrate polygon",
			args: args{
				moduleName:   polygontypes.ModuleName,
				oracleNumber: 20,
			},
		},
		{
			name: "migrate bsc",
			args: args{
				moduleName:   trontypes.ModuleName,
				oracleNumber: 10,
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

			allDelegations := myApp.StakingKeeper.GetAllDelegations(ctx)
			require.EqualValues(t, len(allDelegations), oracles.Len()+valSet.Size())

			for _, oracle := range oracles {
				delegateAddr := oracle.GetDelegateAddress(tt.args.moduleName)
				bonded := myApp.StakingKeeper.GetDelegatorBonded(ctx, delegateAddr)
				require.Equal(t, bonded, oracle.DelegateAmount)

				delegation, found := myApp.StakingKeeper.GetDelegation(ctx, delegateAddr, oracle.GetValidator())
				require.True(t, found)
				require.EqualValues(t, delegation.Shares.TruncateInt().Int64(), 100)

				delegations := myApp.StakingKeeper.GetAllDelegatorDelegations(ctx, delegateAddr)
				require.EqualValues(t, len(delegations), 1)
				require.EqualValues(t, delegation, delegations[0])
			}
		})
	}
}
