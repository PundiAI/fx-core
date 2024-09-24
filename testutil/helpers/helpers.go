package helpers

import (
	"encoding/json"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	tmed25519 "github.com/cometbft/cometbft/crypto/ed25519"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/app"
	fxtypes "github.com/functionx/fx-core/v8/types"
)

func newGenesisState(cdc codec.JSONCodec, moduleBasics module.BasicManager) app.GenesisState {
	genesis := app.NewDefAppGenesisByDenom(cdc, moduleBasics)
	bankState := new(banktypes.GenesisState)
	cdc.MustUnmarshalJSON(genesis[banktypes.ModuleName], bankState)
	bankState.Balances = append(bankState.Balances, banktypes.Balance{
		Address: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromUint64(4_000).MulRaw(1e18))),
	})
	genesis[banktypes.ModuleName] = cdc.MustMarshalJSON(bankState)
	return genesis
}

// Deprecated: please use BaseSuite
func GenerateGenesisValidator(validatorNum int, initCoins sdk.Coins) (valSet *tmtypes.ValidatorSet, genAccs authtypes.GenesisAccounts, balances []banktypes.Balance) {
	if initCoins == nil || initCoins.Len() <= 0 {
		initCoins = sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10_000).MulRaw(1e18)))
	}
	validators := make([]*tmtypes.Validator, validatorNum)
	genAccs = make(authtypes.GenesisAccounts, validatorNum)
	balances = make([]banktypes.Balance, validatorNum)
	for i := 0; i < validatorNum; i++ {
		validator := tmtypes.NewValidator(tmed25519.GenPrivKey().PubKey(), 1)
		validators[i] = validator

		senderPrivKey := secp256k1.GenPrivKey()
		acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
		genAccs[i] = acc

		balance := banktypes.Balance{
			Address: acc.GetAddress().String(),
			Coins:   initCoins,
		}
		balances[i] = balance
	}
	return tmtypes.NewValidatorSet(validators), genAccs, balances
}

// Deprecated: please use BaseSuite
func SetupWithGenesisValSet(t *testing.T, valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *app.App {
	t.Helper()

	myApp := NewApp()
	genesisState := newGenesisState(myApp.AppCodec(), myApp.ModuleBasics)

	// set genesis accounts
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
		pk, err := cryptocodec.FromCmtPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(genAccs[i].GetAddress()).String(),
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

// Deprecated: please use BaseSuite.AddTestSigners
func AddTestAddrs(myApp *app.App, ctx sdk.Context, accNum int, coins sdk.Coins) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		testAddrs[i] = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
		AddTestAddr(myApp, ctx, testAddrs[i], coins)
	}
	return testAddrs
}

// Deprecated: please use BaseSuite.MintToken
func AddTestAddr(myApp *app.App, ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) {
	err := myApp.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	if err != nil {
		panic(err)
	}

	err = myApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins)
	if err != nil {
		panic(err)
	}
}
