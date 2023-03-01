package v2_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	crosschainkeeper "github.com/functionx/fx-core/v3/x/crosschain/keeper"
	crosschainv1 "github.com/functionx/fx-core/v3/x/crosschain/migrations/v1"
	crosschainv2 "github.com/functionx/fx-core/v3/x/crosschain/migrations/v2"
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
			params := types.Params{
				GravityId:                         fmt.Sprintf("fx-%s-bridge", tt.args.moduleName),
				AverageBlockTime:                  5_000,
				ExternalBatchTimeout:              43_200_000,
				AverageExternalBlockTime:          3_000,
				SignedWindow:                      30_000,
				SlashFraction:                     sdk.MustNewDecFromStr("0.8"),
				OracleSetUpdatePowerChangePercent: sdk.MustNewDecFromStr("0.1"),
				IbcTransferTimeoutHeight:          20_000,
				DelegateThreshold:                 sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10_000).MulRaw(1e18)),
				DelegateMultiple:                  10,
			}
			genesisState[tt.args.moduleName] = myApp.AppCodec().MustMarshalJSON(&types.GenesisState{
				Params: params,
			})

			valSet, genAccs, balances := helpers.GenerateGenesisValidator(1, sdk.Coins{})
			var authGenesis authtypes.GenesisState
			myApp.AppCodec().MustUnmarshalJSON(genesisState[authtypes.ModuleName], &authGenesis)
			packAccounts, err := authtypes.PackAccounts(genAccs)
			require.NoError(t, err)
			authGenesis.Accounts = packAccounts
			genesisState[authtypes.ModuleName] = myApp.AppCodec().MustMarshalJSON(&authGenesis)

			// set validators and delegations
			validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
			delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

			bondAmt := sdk.DefaultPowerReduction

			for i, val := range valSet.Validators {
				pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
				require.NoError(t, err)
				pkAny, err := codectypes.NewAnyWithValue(pk)
				require.NoError(t, err)
				validator := stakingtypes.Validator{
					OperatorAddress:   sdk.ValAddress(genAccs[i].GetAddress()).String(),
					ConsensusPubkey:   pkAny,
					Jailed:            false,
					Status:            stakingtypes.Bonded,
					Tokens:            bondAmt,
					DelegatorShares:   sdk.NewDecFromInt(bondAmt),
					Description:       stakingtypes.Description{},
					UnbondingHeight:   int64(0),
					UnbondingTime:     time.Unix(0, 0).UTC(),
					Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
					MinSelfDelegation: sdk.ZeroInt(),
				}
				validators = append(validators, validator)
				delegations = append(delegations, stakingtypes.NewDelegation(genAccs[i].GetAddress(), validator.GetOperator(), sdk.NewDecFromInt(bondAmt)))
			}

			var stakingGenesis stakingtypes.GenesisState
			myApp.AppCodec().MustUnmarshalJSON(genesisState[stakingtypes.ModuleName], &stakingGenesis)
			stakingGenesis.Params.MaxValidators = uint32(len(validators))
			stakingGenesis.Validators = validators
			stakingGenesis.Delegations = delegations
			genesisState[stakingtypes.ModuleName] = myApp.AppCodec().MustMarshalJSON(&stakingGenesis)

			// update balances and total supply
			var bankGenesis banktypes.GenesisState
			myApp.AppCodec().MustUnmarshalJSON(genesisState[banktypes.ModuleName], &bankGenesis)
			for _, b := range balances {
				// add genesis acc tokens and delegated tokens to total supply
				bankGenesis.Supply = bankGenesis.Supply.Add(b.Coins.Add()...)
			}
			for range valSet.Validators {
				bankGenesis.Supply = bankGenesis.Supply.Add(sdk.NewCoin(stakingGenesis.Params.BondDenom, bondAmt))
			}
			bankGenesis.Balances = append(bankGenesis.Balances, balances...)
			// add bonded amount to bonded pool module account
			bankGenesis.Balances = append(bankGenesis.Balances, banktypes.Balance{
				Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
				Coins:   sdk.Coins{sdk.NewCoin(stakingGenesis.Params.BondDenom, bondAmt.MulRaw(int64(len(validators))))},
			})
			genesisState[banktypes.ModuleName] = myApp.AppCodec().MustMarshalJSON(&bankGenesis)

			stateBytes, err := json.MarshalIndent(genesisState, "", " ")
			require.NoError(t, err)

			myApp.InitChain(abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: helpers.ABCIConsensusParams,
				AppStateBytes:   stateBytes,
			})
			ctx := myApp.NewContext(false, tmproto.Header{Time: time.Now()})

			paramSpace := myApp.GetSubspace(tt.args.moduleName)
			paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
			paramSpace.SetParamSet(ctx, &params)

			paramsKey := myApp.GetKey(paramstypes.ModuleName)
			paramsStore := prefix.NewStore(ctx.KVStore(paramsKey), append([]byte(tt.args.moduleName), '/'))
			paramsStore.Set(crosschainv1.ParamStoreOracles, myApp.LegacyAmino().MustMarshalJSON([]string{sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Bytes()).String()}))
			paramsStore.Set(crosschainv1.ParamOracleDepositThreshold, myApp.LegacyAmino().MustMarshalJSON(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10_000).MulRaw(1e18))))

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

			require.EqualValues(t, paramsFromDB.AverageBlockTime, 5_000)
			require.Equal(t, paramsFromDB.DelegateThreshold, sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10_000).MulRaw(1e18)))
			require.EqualValues(t, paramsFromDB.DelegateMultiple, types.DefaultOracleDelegateThreshold)

			defParams := types.DefaultParams()
			defParams.GravityId = fmt.Sprintf("fx-%s-bridge", tt.args.moduleName)
			defParams.AverageBlockTime = 5_000
			defParams.AverageExternalBlockTime = 3_000
			require.EqualValues(t, paramsFromDB, defParams)
		})
	}
}
