// nolint:staticcheck
package v2

import (
	"fmt"
	"reflect"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	"github.com/functionx/fx-core/v3/x/gravity/types"
)

// InitTestGravityDB test use
//
//gocyclo:ignore
func InitTestGravityDB(cdc codec.Codec, legacyAmino *codec.LegacyAmino, genesisState types.GenesisState, paramsStore, gravityStore sdk.KVStore) {
	paramsPrefixStore := prefix.NewStore(paramsStore, append([]byte(types.ModuleName), '/'))
	for _, pair := range genesisState.Params.ParamSetPairs() {
		v := reflect.Indirect(reflect.ValueOf(pair.Value)).Interface()
		if err := pair.ValidatorFn(v); err != nil {
			panic(fmt.Sprintf("value from ParamSetPair is invalid: %s", err))
		}
		bz, err := legacyAmino.MarshalJSON(v)
		if err != nil {
			panic(err)
		}
		paramsPrefixStore.Set(pair.Key, bz)
	}

	gravityStore.Set(types.LastObservedEventNonceKey, sdk.Uint64ToBigEndian(genesisState.LastObservedNonce))

	gravityStore.Set(types.LastObservedEthereumBlockHeightKey, cdc.MustMarshal(&genesisState.LastObservedBlockHeight))

	gravityStore.Set(types.LastObservedValsetKey, cdc.MustMarshal(&genesisState.LastObservedValset))

	gravityStore.Set(types.LastSlashedValsetNonce, sdk.Uint64ToBigEndian(genesisState.LastSlashedValsetNonce))

	gravityStore.Set(types.LastSlashedBatchBlock, sdk.Uint64ToBigEndian(genesisState.LastSlashedBatchBlock))

	for _, delKey := range genesisState.DelegateKeys {
		oracleAddr, err := sdk.ValAddressFromBech32(delKey.Validator)
		if err != nil {
			panic(err)
		}
		bridger := sdk.MustAccAddressFromBech32(delKey.Orchestrator)

		gravityStore.Set(append(types.ValidatorAddressByOrchestratorAddress, bridger.Bytes()...), oracleAddr.Bytes())

		gravityStore.Set(append(types.EthAddressByValidatorKey, oracleAddr.Bytes()...), []byte(delKey.EthAddress))

		gravityStore.Set(append(types.ValidatorByEthAddressKey, []byte(delKey.EthAddress)...), oracleAddr.Bytes())
	}

	for _, item := range genesisState.Erc20ToDenoms {
		gravityStore.Set(append(types.DenomToERC20Key, []byte(item.Denom)...), []byte(item.Erc20))
		gravityStore.Set(append(types.ERC20ToDenomKey, []byte(item.Erc20)...), []byte(item.Denom))
	}

	latestValsetNonce := uint64(0)
	for _, vs := range genesisState.Valsets {
		if vs.Nonce > latestValsetNonce {
			latestValsetNonce = vs.Nonce
		}
		gravityStore.Set(append(types.ValsetRequestKey, sdk.Uint64ToBigEndian(vs.Nonce)...), cdc.MustMarshal(&vs))
	}
	gravityStore.Set(types.LatestValsetNonce, sdk.Uint64ToBigEndian(latestValsetNonce))

	for _, conf := range genesisState.ValsetConfirms {
		addr := sdk.MustAccAddressFromBech32(conf.Orchestrator)
		key := append(types.ValsetConfirmKey, append(sdk.Uint64ToBigEndian(conf.Nonce), addr.Bytes()...)...)
		gravityStore.Set(key, cdc.MustMarshal(&conf))
	}

	for _, batch := range genesisState.Batches {
		key := append(append(types.OutgoingTxBatchKey, []byte(batch.TokenContract)...), sdk.Uint64ToBigEndian(batch.BatchNonce)...)
		gravityStore.Set(key, cdc.MustMarshal(&batch))

		blockKey := append(types.OutgoingTxBatchBlockKey, sdk.Uint64ToBigEndian(batch.Block)...)
		gravityStore.Set(blockKey, cdc.MustMarshal(&batch))
	}

	for _, conf := range genesisState.BatchConfirms {
		addr := sdk.MustAccAddressFromBech32(conf.Orchestrator)
		key := append(types.BatchConfirmKey, append([]byte(conf.TokenContract), append(sdk.Uint64ToBigEndian(conf.Nonce), addr.Bytes()...)...)...)
		gravityStore.Set(key, cdc.MustMarshal(&conf))
	}

	for _, tx := range genesisState.UnbatchedTransfers {
		key := append(types.OutgoingTxPoolKey, sdk.Uint64ToBigEndian(tx.Id)...)
		gravityStore.Set(key, cdc.MustMarshal(&tx))

		// add deprecated key
		amount := make([]byte, 32)
		amount = tx.Erc20Fee.Amount.BigInt().FillBytes(amount)
		idxKey := append(types.SecondIndexOutgoingTxFeeKey, append([]byte(tx.Erc20Fee.Contract), amount...)...)
		var idSet crosschaintypes.IDSet
		if gravityStore.Has(idxKey) {
			bz := gravityStore.Get(idxKey)
			cdc.MustUnmarshal(bz, &idSet)
		}
		idSet.Ids = append(idSet.Ids, tx.Id)
		gravityStore.Set(idxKey, cdc.MustMarshal(&idSet))
	}

	for _, att := range genesisState.Attestations {
		claim, err := types.UnpackAttestationClaim(cdc, &att)
		if err != nil {
			panic("couldn't cast to claim")
		}

		aKey := append(types.OracleAttestationKey, append(sdk.Uint64ToBigEndian(claim.GetEventNonce()), claim.ClaimHash()...)...)
		gravityStore.Set(aKey, cdc.MustMarshal(&att))
	}

	getLastEventNonceByValidator := func(val sdk.ValAddress) uint64 {
		bytes := gravityStore.Get(append(types.LastEventNonceByValidatorKey, val.Bytes()...))
		if len(bytes) == 0 {
			lowestObserved := sdk.BigEndianToUint64(gravityStore.Get(types.LastObservedEventNonceKey))
			if lowestObserved > 0 {
				return lowestObserved - 1
			} else {
				return 0
			}
		}
		return sdk.BigEndianToUint64(bytes)
	}

	for _, att := range genesisState.Attestations {
		claim, err := types.UnpackAttestationClaim(cdc, &att)
		if err != nil {
			panic("couldn't cast to claim")
		}
		for _, vote := range att.Votes {
			val, err := sdk.ValAddressFromBech32(vote)
			if err != nil {
				panic(err)
			}
			last := getLastEventNonceByValidator(val)
			if claim.GetEventNonce() > last {
				gravityStore.Set(append(types.LastEventNonceByValidatorKey, val.Bytes()...), sdk.Uint64ToBigEndian(claim.GetEventNonce()))
				gravityStore.Set(append(types.LastEventBlockHeightByValidatorKey, val.Bytes()...), sdk.Uint64ToBigEndian(claim.GetBlockHeight()))
			}
		}
	}

	gravityStore.Set(types.KeyLastTXPoolID, sdk.Uint64ToBigEndian(genesisState.LastTxPoolId))
	gravityStore.Set(types.KeyLastOutgoingBatchID, sdk.Uint64ToBigEndian(genesisState.LastBatchId))

	// add deprecated key
	gravityStore.Set(types.LastUnBondingBlockHeight, sdk.Uint64ToBigEndian(tmrand.Uint64()))
	gravityStore.Set(append(types.IbcSequenceHeightKey, []byte(fmt.Sprintf("%s/%s/%d", "transfer", "channel-1", tmrand.Uint64()))...), sdk.Uint64ToBigEndian(tmrand.Uint64()))
}

func TestParams() types.Params {
	params := types.Params{
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
	}
	return params
}

func AttClaimToAny(msg proto.Message) *codectypes.Any {
	anyMsg, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		panic(err)
	}
	return anyMsg
}
