// nolint:staticcheck
package v2_test

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	v2 "github.com/functionx/fx-core/v3/x/gravity/legacy/v2"
	"github.com/functionx/fx-core/v3/x/gravity/types"
)

type TestSuite struct {
	suite.Suite
	cdc          codec.Codec
	legacyAmino  *codec.LegacyAmino
	gravityStore store.KVStore
	paramsStore  store.KVStore
	ethStore     store.KVStore
	genesisState types.GenesisState
}

func TestTestSuite(t *testing.T) {
	fxtypes.SetConfig(false)
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) SetupTest() {
	rand.Seed(time.Now().UnixNano())

	gravityStoreKey := sdk.NewKVStoreKey(types.ModuleName)
	paramsStoreKey := sdk.NewKVStoreKey(paramstypes.ModuleName)
	ethStoreKey := sdk.NewKVStoreKey(ethtypes.ModuleName)

	ms := rootmulti.NewStore(dbm.NewMemDB(), log.NewNopLogger())
	ms.MountStoreWithDB(gravityStoreKey, sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(paramsStoreKey, sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(ethStoreKey, sdk.StoreTypeIAVL, nil)
	suite.NoError(ms.LoadLatestVersion())

	suite.gravityStore = ms.GetKVStore(gravityStoreKey)
	suite.paramsStore = ms.GetKVStore(paramsStoreKey)
	suite.ethStore = ms.GetKVStore(ethStoreKey)

	encodingConfig := app.MakeEncodingConfig()
	suite.cdc = encodingConfig.Codec
	suite.legacyAmino = encodingConfig.Amino
}

func (suite *TestSuite) TestMigrateStore() {
	suite.genesisState = types.GenesisState{
		Params:            v2.TestParams(),
		LastObservedNonce: rand.Uint64(),
		LastObservedBlockHeight: types.LastObservedEthereumBlockHeight{
			FxBlockHeight:  rand.Uint64(),
			EthBlockHeight: rand.Uint64(),
		},
		Erc20ToDenoms: []types.ERC20ToDenom{
			{
				Erc20: helpers.GenerateAddress().Hex(),
				Denom: fxtypes.DefaultDenom,
			},
		},
		LastSlashedBatchBlock:  rand.Uint64(),
		LastSlashedValsetNonce: rand.Uint64(),
		LastTxPoolId:           rand.Uint64(),
		LastBatchId:            rand.Uint64(),
	}

	bridgerAddrs := helpers.CreateRandomAccounts(20)
	externals := helpers.CreateRandomAccounts(20)
	valAddrs := helpers.CreateRandomAccounts(20)

	var votes []string
	var members []*types.BridgeValidator
	var delegateKeys []types.MsgSetOrchestratorAddress
	for i, addr := range bridgerAddrs {
		delegateKeys = append(
			delegateKeys,
			types.MsgSetOrchestratorAddress{
				Validator:    sdk.ValAddress(valAddrs[i].Bytes()).String(),
				Orchestrator: addr.String(),
				EthAddress:   common.BytesToAddress(externals[i].Bytes()).String(),
			},
		)
		votes = append(votes, sdk.ValAddress(valAddrs[i].Bytes()).String())
		members = append(members, &types.BridgeValidator{
			Power:      rand.Uint64(),
			EthAddress: common.BytesToAddress(externals[i].Bytes()).String(),
		})
	}
	suite.genesisState.DelegateKeys = delegateKeys

	index := rand.Intn(100)
	for i := 0; i < index; i++ {
		suite.genesisState.Valsets = append(
			suite.genesisState.Valsets,
			types.Valset{
				Nonce:   rand.Uint64(),
				Members: members,
				Height:  rand.Uint64(),
			},
		)

		suite.genesisState.UnbatchedTransfers = append(
			suite.genesisState.UnbatchedTransfers,
			types.OutgoingTransferTx{
				Id:          rand.Uint64(),
				Sender:      sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				DestAddress: helpers.GenerateAddress().Hex(),
				Erc20Token: &types.ERC20Token{
					Contract: helpers.GenerateAddress().Hex(),
					Amount:   sdk.NewInt(rand.Int63() + 1),
				},
				Erc20Fee: &types.ERC20Token{
					Contract: helpers.GenerateAddress().Hex(),
					Amount:   sdk.NewInt(rand.Int63() + 1),
				},
			},
		)
		suite.genesisState.Batches = append(
			suite.genesisState.Batches,
			types.OutgoingTxBatch{
				BatchNonce:   rand.Uint64(),
				BatchTimeout: rand.Uint64(),
				Transactions: []*types.OutgoingTransferTx{
					{
						Id:          rand.Uint64(),
						Sender:      sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
						DestAddress: helpers.GenerateAddress().Hex(),
						Erc20Token: &types.ERC20Token{
							Contract: helpers.GenerateAddress().Hex(),
							Amount:   sdk.NewInt(rand.Int63() + 1),
						},
						Erc20Fee: &types.ERC20Token{
							Contract: helpers.GenerateAddress().Hex(),
							Amount:   sdk.NewInt(rand.Int63() + 1),
						},
					},
					{
						Id:          rand.Uint64(),
						Sender:      sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
						DestAddress: helpers.GenerateAddress().Hex(),
						Erc20Token: &types.ERC20Token{
							Contract: helpers.GenerateAddress().Hex(),
							Amount:   sdk.NewInt(rand.Int63() + 1),
						},
						Erc20Fee: &types.ERC20Token{
							Contract: helpers.GenerateAddress().Hex(),
							Amount:   sdk.NewInt(rand.Int63() + 1),
						},
					},
				},
				TokenContract: helpers.GenerateAddress().Hex(),
				Block:         rand.Uint64(),
				FeeReceive:    helpers.GenerateAddress().Hex(),
			},
		)

		suite.genesisState.BatchConfirms = append(
			suite.genesisState.BatchConfirms,
			types.MsgConfirmBatch{
				Nonce:         rand.Uint64(),
				TokenContract: helpers.GenerateAddress().Hex(),
				EthSigner:     helpers.GenerateAddress().Hex(),
				Orchestrator:  sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				Signature:     hex.EncodeToString(tmrand.Bytes(65)),
			},
		)

		suite.genesisState.ValsetConfirms = append(
			suite.genesisState.ValsetConfirms,
			types.MsgValsetConfirm{
				Nonce:        rand.Uint64(),
				Orchestrator: sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				EthAddress:   helpers.GenerateAddress().Hex(),
				Signature:    hex.EncodeToString(tmrand.Bytes(65)),
			},
		)
	}

	suite.genesisState.Attestations = []types.Attestation{
		{
			Observed: true,
			Votes:    votes,
			Height:   rand.Uint64(),
			Claim: v2.AttClaimToAny(&types.MsgDepositClaim{
				EventNonce:    rand.Uint64(),
				BlockHeight:   rand.Uint64(),
				TokenContract: helpers.GenerateAddress().Hex(),
				Amount:        sdk.NewInt(rand.Int63() + 1),
				EthSender:     helpers.GenerateAddress().Hex(),
				FxReceiver:    sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				TargetIbc:     "",
				Orchestrator:  sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
			}),
		},
		{
			Observed: true,
			Votes:    votes,
			Height:   rand.Uint64(),
			Claim: v2.AttClaimToAny(&types.MsgWithdrawClaim{
				EventNonce:    rand.Uint64(),
				BlockHeight:   rand.Uint64(),
				BatchNonce:    rand.Uint64(),
				TokenContract: helpers.GenerateAddress().Hex(),
				Orchestrator:  sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
			}),
		},
		{
			Observed: true,
			Votes:    votes,
			Height:   rand.Uint64(),
			Claim: v2.AttClaimToAny(&types.MsgValsetUpdatedClaim{
				EventNonce:   rand.Uint64(),
				BlockHeight:  rand.Uint64(),
				ValsetNonce:  rand.Uint64(),
				Members:      members,
				Orchestrator: sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
			}),
		},
	}

	suite.genesisState.LastObservedValset = types.Valset{
		Nonce:   rand.Uint64(),
		Members: members,
		Height:  rand.Uint64(),
	}

	v2.InitTestGravityDB(suite.cdc, suite.legacyAmino, suite.genesisState, suite.paramsStore, suite.gravityStore)
	v2.MigrateStore(suite.cdc, suite.gravityStore, suite.ethStore)
	ctx := sdk.Context{}.WithChainID(fxtypes.TestnetChainId).WithEventManager(sdk.NewEventManager())
	v2.MigrateValidatorToOracle(ctx, suite.cdc, suite.gravityStore, suite.ethStore, testKeeper{}, testKeeper{})

	gravityStoreIter := suite.gravityStore.Iterator(nil, nil)
	defer gravityStoreIter.Close()
	for ; gravityStoreIter.Valid(); gravityStoreIter.Next() {
		suite.T().Log(sdk.ValAddress(gravityStoreIter.Key()[1:]).String())
		suite.Fail(fmt.Sprintf("%x", gravityStoreIter.Key()[0]))
	}
}

func (suite *TestSuite) TestMigrateStoreByExportJson() {
	data, err := os.ReadFile("gravity.json")
	suite.NoError(err)
	suite.cdc.MustUnmarshalJSON(data, &suite.genesisState)

	v2.InitTestGravityDB(suite.cdc, suite.legacyAmino, suite.genesisState, suite.paramsStore, suite.gravityStore)
	v2.MigrateStore(suite.cdc, suite.gravityStore, suite.ethStore)
	ctx := sdk.Context{}.WithChainID(fxtypes.TestnetChainId).WithEventManager(sdk.NewEventManager())
	v2.MigrateValidatorToOracle(ctx, suite.cdc, suite.gravityStore, suite.ethStore, testKeeper{}, testKeeper{})

	suite.Equal(len(ctx.EventManager().Events()), 20)
	oracles := v2.EthInitOracles(fxtypes.TestnetChainId)
	for _, oracleAddr := range oracles {
		oracle := new(crosschaintypes.Oracle)
		value := suite.ethStore.Get(crosschaintypes.GetOracleKey(sdk.MustAccAddressFromBech32(oracleAddr)))
		suite.cdc.MustUnmarshal(value, oracle)
		found := false
		for _, key := range suite.genesisState.DelegateKeys {
			if key.Validator == oracle.DelegateValidator {
				found = true
				suite.Equal(key.EthAddress, oracle.ExternalAddress)
			}
		}
		suite.True(found)
	}

	gravityStoreIter := suite.gravityStore.Iterator(nil, nil)
	defer gravityStoreIter.Close()
	for ; gravityStoreIter.Valid(); gravityStoreIter.Next() {
		suite.T().Log(gravityStoreIter.Key(), gravityStoreIter.Value())
		suite.Fail(fmt.Sprintf("%x", gravityStoreIter.Key()[0]))
	}
}

type testKeeper struct{}

func (s testKeeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

func (s testKeeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error {
	return nil
}

func (s testKeeper) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10_000).MulRaw(1e18)))
}

func (s testKeeper) IterateAllDenomMetaData(ctx sdk.Context, cb func(banktypes.Metadata) bool) {}

func (s testKeeper) GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool) {
	return stakingtypes.Validator{Jailed: false, Status: stakingtypes.Bonded, OperatorAddress: addr.String()}, true
}

func (s testKeeper) Delegate(ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdk.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (newShares sdk.Dec, err error) {
	return sdk.NewDec(1), nil
}
