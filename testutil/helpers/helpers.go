package helpers

import (
	"encoding/json"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	tmed25519 "github.com/cometbft/cometbft/crypto/ed25519"
	tmtypes "github.com/cometbft/cometbft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/app"
	fxtypes "github.com/pundiai/fx-core/v8/types"
)

func generateGenesisValidator(validatorNum int) (*tmtypes.ValidatorSet, []*Signer) {
	validators := make([]*tmtypes.Validator, validatorNum)
	valPrivs := make([]*Signer, validatorNum)
	for i := 0; i < validatorNum; i++ {
		validator := tmtypes.NewValidator(tmed25519.GenPrivKey().PubKey(), 1)
		validators[i] = validator

		valPrivs[i] = NewSigner(NewEthPrivKey())
	}
	return tmtypes.NewValidatorSet(validators), valPrivs
}

func setupWithGenesisValSet(t *testing.T, valSet *tmtypes.ValidatorSet, valPrivs []*Signer, initCoins ...sdk.Coin) *app.App {
	t.Helper()
	if len(initCoins) == 0 {
		initCoins = sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10_000).MulRaw(1e18)))
	}

	// set validators and delegations
	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction
	genAccs := make([]authtypes.GenesisAccount, 0, len(valSet.Validators))
	genAccBalances := make([]banktypes.Balance, 0, len(valSet.Validators))

	for i, val := range valSet.Validators {
		privKey := valPrivs[i]
		acc := authtypes.NewBaseAccount(privKey.AccAddress(), privKey.PubKey(), 0, 0)
		genAccs = append(genAccs, acc)
		genAccBalances = append(genAccBalances, banktypes.Balance{
			Address: acc.GetAddress().String(),
			Coins:   initCoins,
		})

		pk, err := cryptocodec.FromCmtPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(acc.GetAddress()).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdkmath.LegacyNewDecFromInt(bondAmt),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec()),
			MinSelfDelegation: sdkmath.OneInt().Mul(sdkmath.NewInt(10)),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[i].GetAddress().String(), validator.GetOperator(), sdkmath.LegacyNewDecFromInt(bondAmt)))
	}

	myApp := NewApp()
	genesisState := myApp.DefaultGenesis()

	var stakingGenesis stakingtypes.GenesisState
	myApp.AppCodec().MustUnmarshalJSON(genesisState[stakingtypes.ModuleName], &stakingGenesis)
	stakingGenesis.Params.MaxValidators = uint32(len(validators))
	stakingGenesis.Validators = validators
	stakingGenesis.Delegations = delegations
	genesisState[stakingtypes.ModuleName] = myApp.AppCodec().MustMarshalJSON(&stakingGenesis)

	// set genesis accounts
	var authGenesis authtypes.GenesisState
	myApp.AppCodec().MustUnmarshalJSON(genesisState[authtypes.ModuleName], &authGenesis)
	packAccounts, err := authtypes.PackAccounts(genAccs)
	require.NoError(t, err)
	authGenesis.Accounts = packAccounts
	genesisState[authtypes.ModuleName] = myApp.AppCodec().MustMarshalJSON(&authGenesis)

	// update balances and total supply
	var bankGenesis banktypes.GenesisState
	myApp.AppCodec().MustUnmarshalJSON(genesisState[banktypes.ModuleName], &bankGenesis)
	for _, b := range genAccBalances {
		// add genesis acc tokens and delegated tokens to total supply
		bankGenesis.Supply = bankGenesis.Supply.Add(b.Coins.Add()...)
	}
	for range valSet.Validators {
		bankGenesis.Supply = bankGenesis.Supply.Add(sdk.NewCoin(stakingGenesis.Params.BondDenom, bondAmt))
	}
	bankGenesis.Balances = append(bankGenesis.Balances, genAccBalances...)
	// add bonded amount to bonded pool module account
	bankGenesis.Balances = append(bankGenesis.Balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(stakingGenesis.Params.BondDenom, bondAmt.MulRaw(int64(len(validators))))},
	})
	genesisState[banktypes.ModuleName] = myApp.AppCodec().MustMarshalJSON(&bankGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	consensusParams := app.CustomGenesisConsensusParams().ToProto()
	// init chain will set the validator set and initialize the genesis accounts
	_, err = myApp.InitChain(&abci.RequestInitChain{
		ConsensusParams: &consensusParams,
		AppStateBytes:   stateBytes,
		InitialHeight:   1,
	})
	require.NoError(t, err)

	return myApp
}
