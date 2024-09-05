package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	tmed25519 "github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v8/app"
	fxtypes "github.com/functionx/fx-core/v8/types"
	fxstakingtypes "github.com/functionx/fx-core/v8/x/staking/types"
)

func Setup(isCheckTx bool, isShowLog bool) *app.App {
	logger := log.NewNopLogger()
	var traceStore io.Writer
	if isShowLog {
		logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
		traceStore = os.Stdout
	}

	myApp := app.New(logger, dbm.NewMemDB(),
		traceStore, true, map[int64]bool{}, os.TempDir(), 1,
		app.MakeEncodingConfig(), app.EmptyAppOptions{}, baseapp.SetChainID(fxtypes.ChainId()),
	)
	if !isCheckTx {
		// InitChain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(DefGenesisState(myApp.AppCodec()), "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		consensusParams := app.CustomGenesisConsensusParams().ToProto()
		myApp.InitChain(abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: &consensusParams,
			AppStateBytes:   stateBytes,
		})
	}

	return myApp
}

func DefGenesisState(cdc codec.Codec) app.GenesisState {
	genesis := app.NewDefAppGenesisByDenom(fxtypes.DefaultDenom, cdc)
	bankState := new(banktypes.GenesisState)
	cdc.MustUnmarshalJSON(genesis[banktypes.ModuleName], bankState)
	bankState.Balances = append(bankState.Balances, banktypes.Balance{
		Address: sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromUint64(4_000).MulRaw(1e18))),
	})
	genesis[banktypes.ModuleName] = cdc.MustMarshalJSON(bankState)
	return genesis
}

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

// SetupWithGenesisValSet initializes a new App with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit (10^6) in the default token of the app from first genesis
// account. A Nop logger is set in App.
func SetupWithGenesisValSet(t *testing.T, valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *app.App {
	myApp := Setup(true, false)
	genesisState := DefGenesisState(myApp.AppCodec())

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
			MinSelfDelegation: sdkmath.OneInt().Mul(sdkmath.NewInt(10)),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[i].GetAddress(), validator.GetOperator(), sdk.NewDecFromInt(bondAmt)))
	}

	var stakingGenesis fxstakingtypes.GenesisState
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
	myApp.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: &consensusParams,
			AppStateBytes:   stateBytes,
			ChainId:         fxtypes.ChainId(),
		},
	)

	// commit genesis changes
	myApp.Commit()
	myApp.BeginBlock(abci.RequestBeginBlock{
		Header: tmproto.Header{
			Height:             myApp.LastBlockHeight() + 1,
			AppHash:            myApp.LastCommitID().Hash,
			ValidatorsHash:     valSet.Hash(),
			NextValidatorsHash: valSet.Hash(),
			ChainID:            fxtypes.ChainId(),
		},
	})

	return myApp
}

// CreateRandomAccounts generated addresses in random order
func CreateRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

// createIncrementalAccounts is a strategy used by addTestAddrs() in order to generated addresses in ascending order.
func createIncrementalAccounts(accNum int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (accNum + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") // base address string

		buffer.WriteString(numString) // adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHexUnsafe(buffer.String())
		bech := res.String()
		addr, _ := TestAddr(buffer.String(), bech)

		addresses = append(addresses, addr)
		buffer.Reset()
	}

	return addresses
}

// AddTestAddrs constructs and returns accNum amount of accounts with an initial balance of accAmt in random order
func AddTestAddrs(myApp *app.App, ctx sdk.Context, accNum int, coins sdk.Coins) []sdk.AccAddress {
	return addTestAddrs(myApp, ctx, accNum, coins, CreateRandomAccounts)
}

func AddTestAddrsIncremental(myApp *app.App, ctx sdk.Context, accNum int, coins sdk.Coins) []sdk.AccAddress {
	return addTestAddrs(myApp, ctx, accNum, coins, createIncrementalAccounts)
}

func addTestAddrs(myApp *app.App, ctx sdk.Context, accNum int, coin sdk.Coins, strategy func(int) []sdk.AccAddress) []sdk.AccAddress {
	testAddrs := strategy(accNum)
	for _, addr := range testAddrs {
		AddTestAddr(myApp, ctx, addr, coin)
	}
	return testAddrs
}

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

func TestAddr(addr string, bech string) (sdk.AccAddress, error) {
	res, err := sdk.AccAddressFromHexUnsafe(addr)
	if err != nil {
		return nil, err
	}
	bechExpected := res.String()
	if bech != bechExpected {
		return nil, fmt.Errorf("bech encoding doesn't match reference")
	}

	accAddr, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(accAddr, res) {
		return nil, err
	}

	return res, nil
}

func MintBlock(myApp *app.App, ctx sdk.Context, block ...int64) sdk.Context {
	nextHeight := ctx.BlockHeight() + 1
	if len(block) > 0 {
		nextHeight = ctx.BlockHeight() + block[0]
	}
	for i := ctx.BlockHeight(); i <= nextHeight; {
		myApp.EndBlock(abci.RequestEndBlock{Height: i})
		myApp.Commit()
		i++
		header := ctx.BlockHeader()
		header.Height = i
		myApp.BeginBlock(abci.RequestBeginBlock{
			Header: header,
		})
		ctx = myApp.NewContext(false, header)
	}
	return ctx
}
