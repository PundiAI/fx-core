// nolint
package v2_test

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"reflect"
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
		Params: types.Params{
			GravityId:                      "fx-bridge-eth",
			BridgeChainId:                  1,
			SignedValsetsWindow:            10000,
			SignedBatchesWindow:            10000,
			SignedClaimsWindow:             10000,
			TargetBatchTimeout:             43200000,
			AverageBlockTime:               5000,
			AverageEthBlockTime:            15000,
			SlashFractionValset:            sdk.NewDec(1).Quo(sdk.NewDec(1000)),
			SlashFractionBatch:             sdk.NewDec(1).Quo(sdk.NewDec(1000)),
			SlashFractionClaim:             sdk.NewDec(1).Quo(sdk.NewDec(1000)),
			SlashFractionConflictingClaim:  sdk.NewDec(1).Quo(sdk.NewDec(1000)),
			UnbondSlashingValsetsWindow:    10000,
			IbcTransferTimeoutHeight:       10000,
			ValsetUpdatePowerChangePercent: sdk.NewDec(1).Quo(sdk.NewDec(10)),
		},
		LastObservedNonce: rand.Uint64(),
		LastObservedBlockHeight: types.LastObservedEthereumBlockHeight{
			FxBlockHeight:  rand.Uint64(),
			EthBlockHeight: rand.Uint64(),
		},
		DelegateKeys: []types.MsgSetOrchestratorAddress{
			{
				Validator:    sdk.ValAddress(helpers.NewPriKey().PubKey().Address()).String(),
				Orchestrator: sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
				EthAddress:   helpers.GenerateAddress().Hex(),
			},
			{
				Validator:    sdk.ValAddress(helpers.NewPriKey().PubKey().Address()).String(),
				Orchestrator: sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
				EthAddress:   helpers.GenerateAddress().Hex(),
			},
			{
				Validator:    sdk.ValAddress(helpers.NewPriKey().PubKey().Address()).String(),
				Orchestrator: sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
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
				Sender:      sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
				DestAddress: helpers.GenerateAddress().Hex(),
				Erc20Token: &types.ERC20Token{
					Contract: helpers.GenerateAddress().Hex(),
					Amount:   sdk.NewInt(rand.Int63()),
				},
				Erc20Fee: &types.ERC20Token{
					Contract: helpers.GenerateAddress().Hex(),
					Amount:   sdk.NewInt(rand.Int63()),
				},
			},
			{
				Id:          rand.Uint64(),
				Sender:      sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
				DestAddress: helpers.GenerateAddress().Hex(),
				Erc20Token: &types.ERC20Token{
					Contract: helpers.GenerateAddress().Hex(),
					Amount:   sdk.NewInt(rand.Int63()),
				},
				Erc20Fee: &types.ERC20Token{
					Contract: helpers.GenerateAddress().Hex(),
					Amount:   sdk.NewInt(rand.Int63()),
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
						Sender:      sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
						DestAddress: helpers.GenerateAddress().Hex(),
						Erc20Token: &types.ERC20Token{
							Contract: helpers.GenerateAddress().Hex(),
							Amount:   sdk.NewInt(rand.Int63()),
						},
						Erc20Fee: &types.ERC20Token{
							Contract: helpers.GenerateAddress().Hex(),
							Amount:   sdk.NewInt(rand.Int63()),
						},
					},
					{
						Id:          rand.Uint64(),
						Sender:      sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
						DestAddress: helpers.GenerateAddress().Hex(),
						Erc20Token: &types.ERC20Token{
							Contract: helpers.GenerateAddress().Hex(),
							Amount:   sdk.NewInt(rand.Int63()),
						},
						Erc20Fee: &types.ERC20Token{
							Contract: helpers.GenerateAddress().Hex(),
							Amount:   sdk.NewInt(rand.Int63()),
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
				Orchestrator:  sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
				Signature:     hex.EncodeToString(tmrand.Bytes(65)),
			},
			{
				Nonce:         rand.Uint64(),
				TokenContract: helpers.GenerateAddress().Hex(),
				EthSigner:     helpers.GenerateAddress().Hex(),
				Orchestrator:  sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
				Signature:     hex.EncodeToString(tmrand.Bytes(65)),
			},
		},
		ValsetConfirms: []types.MsgValsetConfirm{
			{
				Nonce:        rand.Uint64(),
				Orchestrator: sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
				EthAddress:   helpers.GenerateAddress().Hex(),
				Signature:    hex.EncodeToString(tmrand.Bytes(65)),
			},
			{
				Nonce:        rand.Uint64(),
				Orchestrator: sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
				EthAddress:   helpers.GenerateAddress().Hex(),
				Signature:    hex.EncodeToString(tmrand.Bytes(65)),
			},
		},
		Attestations: []types.Attestation{
			{
				Observed: true,
				Votes: []string{
					sdk.ValAddress(helpers.NewPriKey().PubKey().Address()).String(),
					sdk.ValAddress(helpers.NewPriKey().PubKey().Address()).String(),
					sdk.ValAddress(helpers.NewPriKey().PubKey().Address()).String(),
				},
				Height: rand.Uint64(),
				Claim: suite.toAny(&types.MsgDepositClaim{
					EventNonce:    rand.Uint64(),
					BlockHeight:   rand.Uint64(),
					TokenContract: helpers.GenerateAddress().Hex(),
					Amount:        sdk.NewInt(rand.Int63()),
					EthSender:     helpers.GenerateAddress().Hex(),
					FxReceiver:    sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
					TargetIbc:     "",
					Orchestrator:  sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
				}),
			},
			{
				Observed: true,
				Votes: []string{
					sdk.ValAddress(helpers.NewPriKey().PubKey().Address()).String(),
					sdk.ValAddress(helpers.NewPriKey().PubKey().Address()).String(),
					sdk.ValAddress(helpers.NewPriKey().PubKey().Address()).String(),
				},
				Height: rand.Uint64(),
				Claim: suite.toAny(&types.MsgWithdrawClaim{
					EventNonce:    rand.Uint64(),
					BlockHeight:   rand.Uint64(),
					BatchNonce:    rand.Uint64(),
					TokenContract: helpers.GenerateAddress().Hex(),
					Orchestrator:  sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
				}),
			},
			{
				Observed: true,
				Votes: []string{
					sdk.ValAddress(helpers.NewPriKey().PubKey().Address()).String(),
					sdk.ValAddress(helpers.NewPriKey().PubKey().Address()).String(),
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
					Orchestrator: sdk.AccAddress(helpers.NewPriKey().PubKey().Address()).String(),
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
	suite.InitGravityDB()
	v2.MigrateStore(suite.cdc, suite.gravityStore, suite.ethStore)

	gravityStoreIter := suite.gravityStore.Iterator(nil, nil)
	defer gravityStoreIter.Close()
	for ; gravityStoreIter.Valid(); gravityStoreIter.Next() {
		suite.Contains(append(types.EthAddressByValidatorKey, append(types.DenomToERC20Key, types.ERC20ToDenomKey...)...), gravityStoreIter.Key()[0])
	}
}

func (suite *TestSuite) toAny(msg proto.Message) *codectypes.Any {
	any, err := codectypes.NewAnyWithValue(msg)
	suite.Require().NoError(err)
	return any
}

func (suite *TestSuite) InitGravityDB() {
	paramsPrefixStore := prefix.NewStore(suite.paramsStore, append([]byte(types.ModuleName), '/'))
	for _, pair := range suite.genesisState.Params.ParamSetPairs() {
		v := reflect.Indirect(reflect.ValueOf(pair.Value)).Interface()
		if err := pair.ValidatorFn(v); err != nil {
			panic(fmt.Sprintf("value from ParamSetPair is invalid: %s", err))
		}
		bz, err := suite.legacyAmino.MarshalJSON(v)
		if err != nil {
			panic(err)
		}
		paramsPrefixStore.Set(pair.Key, bz)
	}

	suite.gravityStore.Set(types.LastObservedEventNonceKey, sdk.Uint64ToBigEndian(suite.genesisState.LastObservedNonce))

	suite.gravityStore.Set(types.LastObservedEthereumBlockHeightKey, suite.cdc.MustMarshal(&suite.genesisState.LastObservedBlockHeight))

	suite.gravityStore.Set(types.LastObservedValsetKey, suite.cdc.MustMarshal(&suite.genesisState.LastObservedValset))

	suite.gravityStore.Set(types.LastSlashedValsetNonce, sdk.Uint64ToBigEndian(suite.genesisState.LastSlashedValsetNonce))

	suite.gravityStore.Set(types.LastSlashedBatchBlock, sdk.Uint64ToBigEndian(suite.genesisState.LastSlashedBatchBlock))

	for _, delKey := range suite.genesisState.DelegateKeys {
		oracleAddr, err := sdk.ValAddressFromBech32(delKey.Validator)
		if err != nil {
			panic(err)
		}
		bridger := sdk.MustAccAddressFromBech32(delKey.Orchestrator)

		suite.gravityStore.Set(append(types.ValidatorAddressByOrchestratorAddress, bridger.Bytes()...), oracleAddr.Bytes())

		suite.gravityStore.Set(append(types.EthAddressByValidatorKey, oracleAddr.Bytes()...), []byte(delKey.EthAddress))

		suite.gravityStore.Set(append(types.ValidatorByEthAddressKey, []byte(delKey.EthAddress)...), oracleAddr.Bytes())
	}

	for _, item := range suite.genesisState.Erc20ToDenoms {
		suite.gravityStore.Set(append(types.DenomToERC20Key, []byte(item.Denom)...), []byte(item.Erc20))
		suite.gravityStore.Set(append(types.ERC20ToDenomKey, []byte(item.Erc20)...), []byte(item.Denom))
	}

	for _, vs := range suite.genesisState.Valsets {
		suite.gravityStore.Set(append(types.ValsetRequestKey, sdk.Uint64ToBigEndian(vs.Nonce)...), suite.cdc.MustMarshal(&vs))
		suite.gravityStore.Set(types.LatestValsetNonce, sdk.Uint64ToBigEndian(vs.Nonce))
	}

	for _, conf := range suite.genesisState.ValsetConfirms {
		addr := sdk.MustAccAddressFromBech32(conf.Orchestrator)
		key := append(types.ValsetConfirmKey, append(sdk.Uint64ToBigEndian(conf.Nonce), addr.Bytes()...)...)
		suite.gravityStore.Set(key, suite.cdc.MustMarshal(&conf))
	}

	for _, batch := range suite.genesisState.Batches {
		key := append(append(types.OutgoingTxBatchKey, []byte(batch.TokenContract)...), sdk.Uint64ToBigEndian(batch.BatchNonce)...)
		suite.gravityStore.Set(key, suite.cdc.MustMarshal(&batch))

		blockKey := append(types.OutgoingTxBatchBlockKey, sdk.Uint64ToBigEndian(batch.Block)...)
		suite.gravityStore.Set(blockKey, suite.cdc.MustMarshal(&batch))
	}

	for _, conf := range suite.genesisState.BatchConfirms {
		addr := sdk.MustAccAddressFromBech32(conf.Orchestrator)
		key := append(types.BatchConfirmKey, append([]byte(conf.TokenContract), append(sdk.Uint64ToBigEndian(conf.Nonce), addr.Bytes()...)...)...)
		suite.gravityStore.Set(key, suite.cdc.MustMarshal(&conf))
	}

	for _, tx := range suite.genesisState.UnbatchedTransfers {
		key := append(types.OutgoingTxPoolKey, sdk.Uint64ToBigEndian(tx.Id)...)
		suite.gravityStore.Set(key, suite.cdc.MustMarshal(&tx))

		// add deprecated key
		amount := make([]byte, 32)
		amount = tx.Erc20Fee.Amount.BigInt().FillBytes(amount)
		idxKey := append(types.SecondIndexOutgoingTxFeeKey, append([]byte(tx.Erc20Fee.Contract), amount...)...)
		var idSet crosschaintypes.IDSet
		if suite.gravityStore.Has(idxKey) {
			bz := suite.gravityStore.Get(idxKey)
			suite.cdc.MustUnmarshal(bz, &idSet)
		}
		idSet.Ids = append(idSet.Ids, tx.Id)
		suite.gravityStore.Set(idxKey, suite.cdc.MustMarshal(&idSet))
	}

	attMap := make(map[uint64][]types.Attestation)
	for _, att := range suite.genesisState.Attestations {
		claim, err := types.UnpackAttestationClaim(suite.cdc, &att)
		if err != nil {
			panic("couldn't cast to claim")
		}
		if val, ok := attMap[claim.GetEventNonce()]; !ok {
			attMap[claim.GetEventNonce()] = []types.Attestation{att}
		} else {
			attMap[claim.GetEventNonce()] = append(val, att)
		}

		aKey := append(types.OracleAttestationKey, append(sdk.Uint64ToBigEndian(claim.GetEventNonce()), claim.ClaimHash()...)...)
		suite.gravityStore.Set(aKey, suite.cdc.MustMarshal(&att))
	}

	for _, att := range suite.genesisState.Attestations {
		claim, err := types.UnpackAttestationClaim(suite.cdc, &att)
		if err != nil {
			panic("couldn't cast to claim")
		}
		for _, vote := range att.Votes {
			val, err := sdk.ValAddressFromBech32(vote)
			if err != nil {
				panic(err)
			}
			last := suite.GetLastEventNonceByValidator(val, attMap)
			if claim.GetEventNonce() > last {
				suite.gravityStore.Set(append(types.LastEventNonceByValidatorKey, val.Bytes()...), sdk.Uint64ToBigEndian(claim.GetEventNonce()))
				suite.gravityStore.Set(append(types.LastEventBlockHeightByValidatorKey, val.Bytes()...), sdk.Uint64ToBigEndian(claim.GetBlockHeight()))
			}
		}
	}

	suite.gravityStore.Set(types.KeyLastTXPoolID, sdk.Uint64ToBigEndian(suite.genesisState.LastTxPoolId))
	suite.gravityStore.Set(types.KeyLastOutgoingBatchID, sdk.Uint64ToBigEndian(suite.genesisState.LastBatchId))

	// add deprecated key
	suite.gravityStore.Set(types.LastUnBondingBlockHeight, sdk.Uint64ToBigEndian(rand.Uint64()))
	suite.gravityStore.Set(append(types.IbcSequenceHeightKey, []byte(fmt.Sprintf("%s/%s/%d", "transfer", "channel-1", rand.Uint64()))...), sdk.Uint64ToBigEndian(rand.Uint64()))
}

func (suite *TestSuite) GetLastEventNonceByValidator(val sdk.ValAddress, attMap map[uint64][]types.Attestation) uint64 {
	bytes := suite.gravityStore.Get(append(types.LastEventNonceByValidatorKey, val.Bytes()...))
	if len(bytes) == 0 {
		lowestObserved := sdk.BigEndianToUint64(suite.gravityStore.Get(types.LastObservedEventNonceKey))
		if len(attMap) == 0 {
			return lowestObserved
		}
		for nonce, atts := range attMap {
			for att := range atts {
				if atts[att].Observed && nonce < lowestObserved {
					lowestObserved = nonce
				}
			}
		}
		if lowestObserved > 0 {
			return lowestObserved - 1
		} else {
			return 0
		}
	}
	return sdk.BigEndianToUint64(bytes)
}
