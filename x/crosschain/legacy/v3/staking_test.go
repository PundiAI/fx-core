package v3_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	crosschainv2 "github.com/functionx/fx-core/v3/x/crosschain/legacy/v2"
	crosschainv3 "github.com/functionx/fx-core/v3/x/crosschain/legacy/v3"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func TestMigrateDepositToStaking(t *testing.T) {
	proposalOracle := types.ProposalOracle{}
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
			proposalOracle.Oracles = append(proposalOracle.Oracles, oracles[i].OracleAddress)
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
			valSet, genAccs, balances := helpers.GenerateGenesisValidator(20, nil)
			myApp := helpers.SetupWithGenesisValSet(t, valSet, genAccs, balances...)

			ctx := myApp.NewContext(false, tmproto.Header{Time: time.Now(), Height: 1})
			validators := myApp.StakingKeeper.GetBondedValidatorsByPower(ctx)
			require.Equal(t, len(validators), 20)
			delegatorValidator := validators[0]

			totalTokens := sdk.NewInt(0)
			totalDelegatorShares := sdk.NewDec(0)
			for _, validator := range validators {
				totalTokens = totalTokens.Add(validator.Tokens)
				totalDelegatorShares = totalDelegatorShares.Add(validator.DelegatorShares)
			}
			require.EqualValues(t, totalTokens.String(), totalDelegatorShares.RoundInt().String())

			bondedPool := myApp.StakingKeeper.GetBondedPool(ctx)
			bondedPoolAllBalances := myApp.BankKeeper.GetAllBalances(ctx, bondedPool.GetAddress())
			require.Equal(t, bondedPoolAllBalances.String(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, totalTokens)).String())

			notBondedPool := myApp.StakingKeeper.GetNotBondedPool(ctx)
			notBondedPoolAllBalances := myApp.BankKeeper.GetAllBalances(ctx, notBondedPool.GetAddress())
			require.Equal(t, notBondedPoolAllBalances.String(), sdk.Coins{}.String())

			validator, found := myApp.StakingKeeper.GetValidator(ctx, delegatorValidator.GetOperator())
			require.True(t, found)
			require.EqualValues(t, validator, delegatorValidator)

			oracles := newOracles(tt.args.oracleNumber, delegatorValidator.OperatorAddress)
			for _, oracle := range oracles {
				delegateCoins := sdk.Coins{sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: oracle.DelegateAmount}}
				err := myApp.BankKeeper.MintCoins(ctx, tt.args.moduleName, delegateCoins)
				require.NoError(t, err)
			}

			require.NoError(t, crosschainv2.MigrateDepositToStaking(ctx, tt.args.moduleName, myApp.StakingKeeper,
				myApp.BankKeeper, oracles, delegatorValidator.GetOperator()))
			require.NoError(t, crosschainv3.MigrateDepositToStaking(ctx, tt.args.moduleName, myApp.StakingKeeper,
				myApp.BankKeeper, oracles, delegatorValidator.GetOperator()))

			allDelegations := myApp.StakingKeeper.GetAllDelegations(ctx)
			require.EqualValues(t, len(allDelegations), oracles.Len()+valSet.Size())

			for _, oracle := range oracles {
				delegateAddr := oracle.GetDelegateAddress(tt.args.moduleName)
				bonded := myApp.StakingKeeper.GetDelegatorBonded(ctx, delegateAddr)
				require.Equal(t, bonded.String(), oracle.DelegateAmount.String())

				delegation, found := myApp.StakingKeeper.GetDelegation(ctx, delegateAddr, oracle.GetValidator())
				require.True(t, found)
				require.Equal(t, delegation.Shares.TruncateInt().String(), oracle.DelegateAmount.String())

				delegations := myApp.StakingKeeper.GetAllDelegatorDelegations(ctx, delegateAddr)
				require.EqualValues(t, len(delegations), 1)
				require.EqualValues(t, delegation, delegations[0])

				bondedPoolAllBalances = bondedPoolAllBalances.Add(sdk.NewCoin(fxtypes.DefaultDenom, oracle.DelegateAmount))
				totalTokens = totalTokens.Add(oracle.DelegateAmount)
				totalDelegatorShares = totalDelegatorShares.Add(oracle.DelegateAmount.ToDec())
			}

			validators1 := myApp.StakingKeeper.GetBondedValidatorsByPower(ctx)
			totalTokens1 := sdk.NewInt(0)
			totalDelegatorShares1 := sdk.NewDec(0)
			for _, validator := range validators1 {
				totalTokens1 = totalTokens1.Add(validator.Tokens)
				totalDelegatorShares1 = totalDelegatorShares1.Add(validator.DelegatorShares)
			}
			require.Equal(t, totalTokens1.String(), totalDelegatorShares1.RoundInt().String())
			require.Equal(t, totalTokens.String(), totalTokens1.String())
			require.Equal(t, totalDelegatorShares, totalDelegatorShares1)

			require.Equal(t, bondedPoolAllBalances, myApp.BankKeeper.GetAllBalances(ctx, bondedPool.GetAddress()))
			require.Equal(t, sdk.Coins{}, myApp.BankKeeper.GetAllBalances(ctx, notBondedPool.GetAddress()))
		})
	}
}
