// nolint:staticcheck
package v2_test

import (
	"encoding/hex"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
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

func (suite *TestSuite) TestPrefixStore() {
	suite.gravityStore.Set([]byte{1, 1}, []byte{1, 1})
	suite.gravityStore.Set([]byte{1, 2}, []byte{2, 2})
	suite.gravityStore.Set([]byte{1}, []byte{3, 3})

	newStore := prefix.NewStore(suite.gravityStore, []byte{1})
	newStore.Set([]byte{4}, []byte{4, 4})
	iter := newStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		suite.T().Log(iter.Key(), iter.Value())
		newStore.Delete(iter.Key())
	}

	iterator := sdk.KVStorePrefixIterator(suite.gravityStore, []byte{1})
	for ; iterator.Valid(); iterator.Next() {
		suite.T().Log(iterator.Key(), iterator.Value())
	}
}

func (suite *TestSuite) TestMigrateStore() {
	suite.genesisState = types.GenesisState{
		Params:            v2.TestParams(),
		LastObservedNonce: rand.Uint64(),
		LastObservedBlockHeight: types.LastObservedEthereumBlockHeight{
			FxBlockHeight:  rand.Uint64(),
			EthBlockHeight: rand.Uint64(),
		},
		DelegateKeys: []types.MsgSetOrchestratorAddress{
			{
				Validator:    sdk.ValAddress(helpers.GenerateAddress().Bytes()).String(),
				Orchestrator: sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				EthAddress:   helpers.GenerateAddress().Hex(),
			},
			{
				Validator:    sdk.ValAddress(helpers.GenerateAddress().Bytes()).String(),
				Orchestrator: sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				EthAddress:   helpers.GenerateAddress().Hex(),
			},
			{
				Validator:    sdk.ValAddress(helpers.GenerateAddress().Bytes()).String(),
				Orchestrator: sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				EthAddress:   helpers.GenerateAddress().Hex(),
			},
		},
		Valsets: []types.Valset{
			{
				Nonce: rand.Uint64(),
				Members: []*types.BridgeValidator{
					{
						Power:      rand.Uint64(),
						EthAddress: helpers.GenerateAddress().Hex(),
					},
					{
						Power:      rand.Uint64(),
						EthAddress: helpers.GenerateAddress().Hex(),
					},
					{
						Power:      rand.Uint64(),
						EthAddress: helpers.GenerateAddress().Hex(),
					},
				},
				Height: rand.Uint64(),
			},
			{
				Nonce: rand.Uint64(),
				Members: []*types.BridgeValidator{
					{
						Power:      rand.Uint64(),
						EthAddress: helpers.GenerateAddress().Hex(),
					},
					{
						Power:      rand.Uint64(),
						EthAddress: helpers.GenerateAddress().Hex(),
					},
					{
						Power:      rand.Uint64(),
						EthAddress: helpers.GenerateAddress().Hex(),
					},
				},
				Height: rand.Uint64(),
			},
		},
		Erc20ToDenoms: []types.ERC20ToDenom{
			{
				Erc20: helpers.GenerateAddress().Hex(),
				Denom: fxtypes.DefaultDenom,
			},
		},
		UnbatchedTransfers: []types.OutgoingTransferTx{
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
		Batches: []types.OutgoingTxBatch{
			{
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
		},
		BatchConfirms: []types.MsgConfirmBatch{
			{
				Nonce:         rand.Uint64(),
				TokenContract: helpers.GenerateAddress().Hex(),
				EthSigner:     helpers.GenerateAddress().Hex(),
				Orchestrator:  sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				Signature:     hex.EncodeToString(tmrand.Bytes(65)),
			},
			{
				Nonce:         rand.Uint64(),
				TokenContract: helpers.GenerateAddress().Hex(),
				EthSigner:     helpers.GenerateAddress().Hex(),
				Orchestrator:  sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				Signature:     hex.EncodeToString(tmrand.Bytes(65)),
			},
		},
		ValsetConfirms: []types.MsgValsetConfirm{
			{
				Nonce:        rand.Uint64(),
				Orchestrator: sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				EthAddress:   helpers.GenerateAddress().Hex(),
				Signature:    hex.EncodeToString(tmrand.Bytes(65)),
			},
			{
				Nonce:        rand.Uint64(),
				Orchestrator: sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				EthAddress:   helpers.GenerateAddress().Hex(),
				Signature:    hex.EncodeToString(tmrand.Bytes(65)),
			},
		},
		Attestations: []types.Attestation{
			{
				Observed: true,
				Votes: []string{
					sdk.ValAddress(helpers.GenerateAddress().Bytes()).String(),
					sdk.ValAddress(helpers.GenerateAddress().Bytes()).String(),
					sdk.ValAddress(helpers.GenerateAddress().Bytes()).String(),
				},
				Height: rand.Uint64(),
				Claim: suite.toAny(&types.MsgDepositClaim{
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
				Votes: []string{
					sdk.ValAddress(helpers.GenerateAddress().Bytes()).String(),
					sdk.ValAddress(helpers.GenerateAddress().Bytes()).String(),
					sdk.ValAddress(helpers.GenerateAddress().Bytes()).String(),
				},
				Height: rand.Uint64(),
				Claim: suite.toAny(&types.MsgWithdrawClaim{
					EventNonce:    rand.Uint64(),
					BlockHeight:   rand.Uint64(),
					BatchNonce:    rand.Uint64(),
					TokenContract: helpers.GenerateAddress().Hex(),
					Orchestrator:  sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				}),
			},
			{
				Observed: true,
				Votes: []string{
					sdk.ValAddress(helpers.GenerateAddress().Bytes()).String(),
					sdk.ValAddress(helpers.GenerateAddress().Bytes()).String(),
				},
				Height: rand.Uint64(),
				Claim: suite.toAny(&types.MsgValsetUpdatedClaim{
					EventNonce:  rand.Uint64(),
					BlockHeight: rand.Uint64(),
					ValsetNonce: rand.Uint64(),
					Members: []*types.BridgeValidator{
						{
							Power:      rand.Uint64(),
							EthAddress: helpers.GenerateAddress().Hex(),
						},
						{
							Power:      rand.Uint64(),
							EthAddress: helpers.GenerateAddress().Hex(),
						},
					},
					Orchestrator: sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
				}),
			},
		},
		LastObservedValset: types.Valset{
			Nonce: rand.Uint64(),
			Members: []*types.BridgeValidator{
				{
					Power:      rand.Uint64(),
					EthAddress: helpers.GenerateAddress().Hex(),
				},
				{
					Power:      rand.Uint64(),
					EthAddress: helpers.GenerateAddress().Hex(),
				},
			},
			Height: rand.Uint64(),
		},
		LastSlashedBatchBlock:  rand.Uint64(),
		LastSlashedValsetNonce: rand.Uint64(),
		LastTxPoolId:           rand.Uint64(),
		LastBatchId:            rand.Uint64(),
	}

	v2.InitTestGravityDB(suite.cdc, suite.legacyAmino, suite.genesisState, suite.paramsStore, suite.gravityStore)
	v2.MigrateStore(suite.cdc, suite.gravityStore, suite.ethStore)

	gravityStoreIter := suite.gravityStore.Iterator(nil, nil)
	defer gravityStoreIter.Close()
	for ; gravityStoreIter.Valid(); gravityStoreIter.Next() {
		suite.T().Log(gravityStoreIter.Key(), gravityStoreIter.Value())
		keys := append(types.ValidatorAddressByOrchestratorAddress, types.EthAddressByValidatorKey...)
		keys = append(keys, types.ValidatorByEthAddressKey...)
		keys = append(keys, append(types.DenomToERC20Key, types.ERC20ToDenomKey...)...)
		suite.Contains(keys, gravityStoreIter.Key()[0])
	}
}

func (suite *TestSuite) TestMigrateStoreByExportJson() {
	data, err := os.ReadFile("gravity.json")
	suite.NoError(err)
	suite.cdc.MustUnmarshalJSON(data, &suite.genesisState)

	v2.InitTestGravityDB(suite.cdc, suite.legacyAmino, suite.genesisState, suite.paramsStore, suite.gravityStore)
	v2.MigrateStore(suite.cdc, suite.gravityStore, suite.ethStore)

	gravityStoreIter := suite.gravityStore.Iterator(nil, nil)
	defer gravityStoreIter.Close()
	for ; gravityStoreIter.Valid(); gravityStoreIter.Next() {
		keys := append(types.ValidatorAddressByOrchestratorAddress, types.EthAddressByValidatorKey...)
		keys = append(keys, types.ValidatorByEthAddressKey...)
		keys = append(keys, append(types.DenomToERC20Key, types.ERC20ToDenomKey...)...)
		suite.Contains(keys, gravityStoreIter.Key()[0])
	}
}

func (suite *TestSuite) toAny(msg proto.Message) *codectypes.Any {
	any, err := codectypes.NewAnyWithValue(msg)
	suite.Require().NoError(err)
	return any
}
