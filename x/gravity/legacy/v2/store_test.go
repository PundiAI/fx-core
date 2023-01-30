// nolint:staticcheck
package v2_test

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

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
	helpers2 "github.com/functionx/fx-core/v3/testutil/helpers"
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
		LastObservedNonce: tmrand.Uint64(),
		LastObservedBlockHeight: types.LastObservedEthereumBlockHeight{
			FxBlockHeight:  tmrand.Uint64(),
			EthBlockHeight: tmrand.Uint64(),
		},
		Erc20ToDenoms: []types.ERC20ToDenom{
			{
				Erc20: helpers2.GenerateAddress().Hex(),
				Denom: fxtypes.DefaultDenom,
			},
		},
		LastSlashedBatchBlock:  tmrand.Uint64(),
		LastSlashedValsetNonce: tmrand.Uint64(),
		LastTxPoolId:           tmrand.Uint64(),
		LastBatchId:            tmrand.Uint64(),
	}

	bridgerAddrs := helpers2.CreateRandomAccounts(20)
	externals := helpers2.CreateRandomAccounts(20)
	valAddrs := helpers2.CreateRandomAccounts(20)

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
			Power:      tmrand.Uint64(),
			EthAddress: common.BytesToAddress(externals[i].Bytes()).String(),
		})
	}
	suite.genesisState.DelegateKeys = delegateKeys

	index := tmrand.Intn(100)
	for i := 0; i < index; i++ {
		suite.genesisState.Valsets = append(
			suite.genesisState.Valsets,
			types.Valset{
				Nonce:   tmrand.Uint64(),
				Members: members,
				Height:  tmrand.Uint64(),
			},
		)

		suite.genesisState.UnbatchedTransfers = append(
			suite.genesisState.UnbatchedTransfers,
			types.OutgoingTransferTx{
				Id:          tmrand.Uint64(),
				Sender:      sdk.AccAddress(helpers2.GenerateAddress().Bytes()).String(),
				DestAddress: helpers2.GenerateAddress().Hex(),
				Erc20Token: &types.ERC20Token{
					Contract: helpers2.GenerateAddress().Hex(),
					Amount:   sdk.NewInt(tmrand.Int63() + 1),
				},
				Erc20Fee: &types.ERC20Token{
					Contract: helpers2.GenerateAddress().Hex(),
					Amount:   sdk.NewInt(tmrand.Int63() + 1),
				},
			},
		)
		suite.genesisState.Batches = append(
			suite.genesisState.Batches,
			types.OutgoingTxBatch{
				BatchNonce:   tmrand.Uint64(),
				BatchTimeout: tmrand.Uint64(),
				Transactions: []*types.OutgoingTransferTx{
					{
						Id:          tmrand.Uint64(),
						Sender:      sdk.AccAddress(helpers2.GenerateAddress().Bytes()).String(),
						DestAddress: helpers2.GenerateAddress().Hex(),
						Erc20Token: &types.ERC20Token{
							Contract: helpers2.GenerateAddress().Hex(),
							Amount:   sdk.NewInt(tmrand.Int63() + 1),
						},
						Erc20Fee: &types.ERC20Token{
							Contract: helpers2.GenerateAddress().Hex(),
							Amount:   sdk.NewInt(tmrand.Int63() + 1),
						},
					},
					{
						Id:          tmrand.Uint64(),
						Sender:      sdk.AccAddress(helpers2.GenerateAddress().Bytes()).String(),
						DestAddress: helpers2.GenerateAddress().Hex(),
						Erc20Token: &types.ERC20Token{
							Contract: helpers2.GenerateAddress().Hex(),
							Amount:   sdk.NewInt(tmrand.Int63() + 1),
						},
						Erc20Fee: &types.ERC20Token{
							Contract: helpers2.GenerateAddress().Hex(),
							Amount:   sdk.NewInt(tmrand.Int63() + 1),
						},
					},
				},
				TokenContract: helpers2.GenerateAddress().Hex(),
				Block:         tmrand.Uint64(),
				FeeReceive:    helpers2.GenerateAddress().Hex(),
			},
		)

		suite.genesisState.BatchConfirms = append(
			suite.genesisState.BatchConfirms,
			types.MsgConfirmBatch{
				Nonce:         tmrand.Uint64(),
				TokenContract: helpers2.GenerateAddress().Hex(),
				EthSigner:     delegateKeys[1%20].EthAddress,
				Orchestrator:  delegateKeys[1%20].Orchestrator,
				Signature:     hex.EncodeToString(tmrand.Bytes(65)),
			},
		)

		suite.genesisState.ValsetConfirms = append(
			suite.genesisState.ValsetConfirms,
			types.MsgValsetConfirm{
				Nonce:        tmrand.Uint64(),
				Orchestrator: delegateKeys[i%20].Orchestrator,
				EthAddress:   delegateKeys[i%20].EthAddress,
				Signature:    hex.EncodeToString(tmrand.Bytes(65)),
			},
		)
	}

	suite.genesisState.Attestations = []types.Attestation{
		{
			Observed: true,
			Votes:    votes,
			Height:   tmrand.Uint64(),
			Claim: v2.AttClaimToAny(&types.MsgDepositClaim{
				EventNonce:    tmrand.Uint64(),
				BlockHeight:   tmrand.Uint64(),
				TokenContract: helpers2.GenerateAddress().Hex(),
				Amount:        sdk.NewInt(tmrand.Int63() + 1),
				EthSender:     helpers2.GenerateAddress().Hex(),
				FxReceiver:    sdk.AccAddress(helpers2.GenerateAddress().Bytes()).String(),
				TargetIbc:     "",
				Orchestrator:  sdk.AccAddress(helpers2.GenerateAddress().Bytes()).String(),
			}),
		},
		{
			Observed: true,
			Votes:    votes,
			Height:   tmrand.Uint64(),
			Claim: v2.AttClaimToAny(&types.MsgWithdrawClaim{
				EventNonce:    tmrand.Uint64(),
				BlockHeight:   tmrand.Uint64(),
				BatchNonce:    tmrand.Uint64(),
				TokenContract: helpers2.GenerateAddress().Hex(),
				Orchestrator:  sdk.AccAddress(helpers2.GenerateAddress().Bytes()).String(),
			}),
		},
		{
			Observed: true,
			Votes:    votes,
			Height:   tmrand.Uint64(),
			Claim: v2.AttClaimToAny(&types.MsgValsetUpdatedClaim{
				EventNonce:   tmrand.Uint64(),
				BlockHeight:  tmrand.Uint64(),
				ValsetNonce:  tmrand.Uint64(),
				Members:      members,
				Orchestrator: sdk.AccAddress(helpers2.GenerateAddress().Bytes()).String(),
			}),
		},
	}

	suite.genesisState.LastObservedValset = types.Valset{
		Nonce:   tmrand.Uint64(),
		Members: members,
		Height:  tmrand.Uint64(),
	}

	v2.InitTestGravityDB(suite.cdc, suite.legacyAmino, suite.genesisState, suite.paramsStore, suite.gravityStore)

	ctx := sdk.Context{}.WithChainID(fxtypes.TestnetChainId).WithEventManager(sdk.NewEventManager()).WithLogger(log.NewNopLogger())
	oracleMap := v2.MigrateValidatorToOracle(ctx, suite.cdc, suite.gravityStore, suite.ethStore, testKeeper{}, testKeeper{})
	v2.MigrateStore(suite.cdc, suite.gravityStore, suite.ethStore, oracleMap)

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

	oracles := v2.GetEthOracleAddrs(fxtypes.TestnetChainId)

	ctx := sdk.Context{}.WithChainID(fxtypes.TestnetChainId).WithEventManager(sdk.NewEventManager()).WithLogger(log.NewNopLogger())
	oracleMap := v2.MigrateValidatorToOracle(ctx, suite.cdc, suite.gravityStore, suite.ethStore, testKeeper{}, testKeeper{})
	suite.Equal(len(oracles), len(oracleMap))

	v2.MigrateStore(suite.cdc, suite.gravityStore, suite.ethStore, oracleMap)

	suite.Equal(len(ctx.EventManager().Events()), 20)

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

func (s testKeeper) SendCoins(sdk.Context, sdk.AccAddress, sdk.AccAddress, sdk.Coins) error {
	return nil
}

func (s testKeeper) SendCoinsFromModuleToModule(sdk.Context, string, string, sdk.Coins) error {
	return nil
}

func (s testKeeper) GetAllBalances(sdk.Context, sdk.AccAddress) sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10_000).MulRaw(1e18)))
}

func (s testKeeper) IterateAllDenomMetaData(sdk.Context, func(banktypes.Metadata) bool) {}

func (s testKeeper) GetValidator(_ sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool) {
	return stakingtypes.Validator{Jailed: false, Status: stakingtypes.Bonded, OperatorAddress: addr.String()}, true
}

func (s testKeeper) Delegate(sdk.Context, sdk.AccAddress, sdk.Int, stakingtypes.BondStatus, stakingtypes.Validator, bool) (newShares sdk.Dec, err error) {
	return sdk.NewDec(1), nil
}
