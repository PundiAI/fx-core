// nolint:staticcheck
package v2

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	"github.com/functionx/fx-core/v3/x/gravity/types"
)

// MigrateStore performs in-place store migrations from v1 to v2.
// migrate data from gravity module
func MigrateStore(cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore, oracleMap map[string]string) {
	// gravity 0x2 -> eth 0x13
	// key                                        			value
	// prefix     external-address                			oracle-address
	// [0x13][0xd98F9E3B1Bc6927700ce4A963429DC157dD4EBDf]   [0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9]
	// gravity 0xe -> eth 0x14
	// key                                        			value
	// prefix     bridger-address                			oracle-address
	// [0x14][0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9]    	[0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9]
	// migrate on MigrateValidatorToOracle

	// gravity 0x3 -> eth 0x15
	// key                                        			value
	// prefix     nonce  			                		OracleSet
	// [0x15][0 0 0 0 0 0 0 1]                           	[object marshal bytes]
	migratePrefix(gravityStore, ethStore, types.ValsetRequestKey, crosschaintypes.OracleSetRequestKey)

	// gravity 0x4 -> eth 0x16
	// key                                       				 			value
	// prefix     nonce  		oracle-address		               			MsgOracleSetConfirm
	// [0x16][0 0 0 0 0 0 0 1][0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9]     [object marshal bytes]
	migrateOracleSetConfirm(cdc, gravityStore, ethStore)

	// gravity 0x5 -> eth 0x17
	// key                                       				 									value
	// prefix     nonce                             claim-details-hash								Attestation
	// [0x5][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]		[object marshal bytes]
	migrateAttestation(cdc, gravityStore, ethStore, oracleMap)

	// gravity 0x6 and 0x7 -> eth 0x18
	// key                                       				 				value
	// prefix            id											 			OutgoingTransferTx
	// [0x6][0 0 0 0 0 0 0 1]													[object marshal bytes]
	// prefix            token-address            		 fee_amount(byte32) 	IDSet -delete-
	// [0x7][0xb4fA5979babd8Bb7e427157d0d353Cf205F43752][1000000000000][id]		[object marshal bytes] delete
	migrateOutgoingTxPool(cdc, gravityStore, ethStore)

	// gravity 0x8 -> eth 0x20
	// key                                       				 			value
	// prefix            token-address            		 nonce		 		OutgoingTxBatch
	// [0x20][0xb4fA5979babd8Bb7e427157d0d353Cf205F43752][0 0 0 0 0 0 0 1]	[object marshal bytes]
	// migratePrefix(gravityStore, ethStore, types.OutgoingTxBatchKey, crosschaintypes.OutgoingTxBatchKey)
	outgoingTxBatches := migrateOutgoingTxBatch(cdc, gravityStore, ethStore)

	// gravity 0x9 -> eth 0x21
	// key                                  value
	// prefix	block-height		 		OutgoingTxBatch
	// [0x21][0 0 0 0 0 0 0 1]				[object marshal bytes]
	migratePrefix(gravityStore, ethStore, types.OutgoingTxBatchBlockKey, crosschaintypes.OutgoingTxBatchBlockKey)

	// gravity 0xa -> eth 0x22
	// key                                  																			value
	// prefix           token-address                		batch-nonce            oracle-address						MsgConfirmBatch
	// [0x22][0xb4fA5979babd8Bb7e427157d0d353Cf205F43752][0 0 0 0 0 0 0 1][fx1mx8euwcmc6f8wqxwf2trg2wuz47af67lads8yg]	[object marshal bytes]
	migrateConfirmBatch(cdc, gravityStore, ethStore, outgoingTxBatches)

	// gravity 0xb -> eth 0x23
	// key                                  					value
	// prefix           oracle-address                			event-nonce
	// [0x23][0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9]			[0 0 0 0 0 0 0 1]
	// migrate on MigrateValidatorToOracle

	// gravity 0xc -> eth 0x24
	// key         		value
	// prefix           event-nonce
	// [0x24]			[0 0 0 0 0 0 0 1]
	migratePrefix(gravityStore, ethStore, types.LastObservedEventNonceKey, crosschaintypes.LastObservedEventNonceKey)

	// gravity 0xd+"lastTxPoolId" -> eth 0x25+"lastTxPoolId"
	// gravity 0xd+"lastBatchId"  -> eth 0x25+"lastBatchId"
	migratePrefix(gravityStore, ethStore, types.SequenceKeyPrefix, crosschaintypes.SequenceKeyPrefix)

	// gravity 0xf -> eth 0x26
	// key         														value
	// prefix  denom		   											BridgeToken
	// [0x26][eth0xb4fA5979babd8Bb7e427157d0d353Cf205F43752]			[object marshal bytes]

	// gravity 0x10 -> eth 0x27
	// key         														value
	// prefix  token-address   											BridgeToken
	// [0x27][0xb4fA5979babd8Bb7e427157d0d353Cf205F43752]				[object marshal bytes]
	migrateBridgeToken(gravityStore, ethStore)

	// gravity 0x11 -> eth 0x28
	// key         		value
	// prefix           oracle-set-nonce
	// [0x28]			[0 0 0 0 0 0 0 1]
	migratePrefix(gravityStore, ethStore, types.LastSlashedValsetNonce, crosschaintypes.LastSlashedOracleSetNonce)

	// gravity 0x12 -> eth 0x29
	// key         		value
	// prefix           oracle-set-nonce
	// [0x29]			[0 0 0 0 0 0 0 1]
	migratePrefix(gravityStore, ethStore, types.LatestValsetNonce, crosschaintypes.LatestOracleSetNonce)

	// gravity 0x13 -> eth 0x30
	// key         		value
	// prefix           block-height
	// [0x30]			[0 0 0 0 0 0 0 1]
	migratePrefix(gravityStore, ethStore, types.LastSlashedBatchBlock, crosschaintypes.LastSlashedBatchBlock)

	// gravity 0x14 -delete-
	// key         		value
	// prefix           block-height
	// [0x14]			[0 0 0 0 0 0 0 1]
	deletePrefixKey(gravityStore, types.LastUnBondingBlockHeight)

	// gravity 0x15 -> eth 0x32
	// key         		value
	// prefix           LastObservedBlockHeight
	// [0x32]			[object marshal bytes]
	migrateLastObservedBlockHeight(cdc, gravityStore, ethStore)

	// gravity 0x16 -> eth 0x33
	// key         		value
	// prefix           OracleSet
	// [0x33]			[object marshal bytes]
	migratePrefix(gravityStore, ethStore, types.LastObservedValsetKey, crosschaintypes.LastObservedOracleSetKey)

	// gravity 0x17 delete
	deletePrefixKey(gravityStore, types.IbcSequenceHeightKey)

	// gravity 0x18 -> eth 0x35
	// key                                  					value
	// prefix           oracle-address                			block-height
	// [0x35][0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9]			[0 0 0 0 0 0 0 1]
	// migrate on MigrateValidatorToOracle
}

func migratePrefix(gravityStore, ethStore sdk.KVStore, oldPrefix, newPrefix []byte) {
	oldStore := prefix.NewStore(gravityStore, oldPrefix)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		key := oldStoreIter.Key()
		ethStore.Set(append(newPrefix, key...), oldStoreIter.Value())
		oldStore.Delete(key)
	}
}

func MigrateValidatorToOracle(ctx sdk.Context, cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore, stakingKeeper StakingKeeper, bankKeeper BankKeeper) map[string]string {
	oldOracleMap := make(map[string]string)
	chainOracle := new(crosschaintypes.ProposalOracle)
	totalPower := sdk.ZeroInt()

	ethOracles := GetEthOracleAddrs(ctx.ChainID())
	ctx.Logger().Info("migrating validator to oracle", "module", "gravity", "number", len(ethOracles))
	index := 0
	minDelegateAmount := sdk.DefaultPowerReduction.MulRaw(100)

	oldStore := prefix.NewStore(gravityStore, types.ValidatorAddressByOrchestratorAddress)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()
	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		bridgerAddr := sdk.AccAddress(oldStoreIter.Key())
		oldOracleAddress := sdk.AccAddress(oldStoreIter.Value())
		externalAddress := string(gravityStore.Get(append(types.EthAddressByValidatorKey, oldOracleAddress.Bytes()...)))
		validator, found := stakingKeeper.GetValidator(ctx, oldOracleAddress.Bytes())
		if !found {
			panic(fmt.Sprintf("no found validator: %s", sdk.ValAddress(oldOracleAddress).String()))
		}
		oracle := crosschaintypes.Oracle{
			BridgerAddress:    bridgerAddr.String(),
			ExternalAddress:   externalAddress,
			StartHeight:       0,
			DelegateValidator: validator.OperatorAddress,
			DelegateAmount:    sdk.ZeroInt(),
			Online:            false,
			OracleAddress:     oldOracleAddress.String(),
			SlashTimes:        0,
		}
		if len(ethOracles) > index {
			oracle.OracleAddress = ethOracles[index]
			oracleAddr := oracle.GetOracle()
			balances := bankKeeper.GetAllBalances(ctx, oracleAddr)
			if balances.AmountOf(fxtypes.DefaultDenom).GTE(minDelegateAmount) {
				delegateAddr := oracle.GetDelegateAddress(ethtypes.ModuleName)
				if err := bankKeeper.SendCoins(ctx, oracleAddr, delegateAddr,
					sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, minDelegateAmount))); err != nil {
					panic("send to coins error: " + err.Error())
				}
				newShares, err := stakingKeeper.Delegate(ctx,
					delegateAddr, minDelegateAmount, stakingtypes.Unbonded, validator, true)
				if err != nil {
					panic("gravity migrate to eth error: " + err.Error())
				}
				oracle.StartHeight = ctx.BlockHeight()
				oracle.DelegateAmount = minDelegateAmount
				oracle.Online = true
				ctx.EventManager().EmitEvent(sdk.NewEvent(
					stakingtypes.EventTypeDelegate,
					sdk.NewAttribute(stakingtypes.AttributeKeyValidator, oracle.DelegateValidator),
					sdk.NewAttribute(sdk.AttributeKeyAmount, oracle.DelegateAmount.String()),
					sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
				))
			}
		}
		index = index + 1

		if oracle.Online {
			totalPower = totalPower.Add(oracle.GetPower())
		}
		oldOracleMap[sdk.ValAddress(oldStoreIter.Value()).String()] = oracle.OracleAddress

		oracleAddress := oracle.GetOracle()
		ethStore.Set(append(crosschaintypes.OracleAddressByExternalKey, []byte(oracle.ExternalAddress)...), oracleAddress.Bytes())
		ethStore.Set(append(crosschaintypes.OracleAddressByBridgerKey, bridgerAddr.Bytes()...), oracleAddress.Bytes())
		// SetOracle
		ethStore.Set(crosschaintypes.GetOracleKey(oracleAddress), cdc.MustMarshal(&oracle))
		oldStore.Delete(oldStoreIter.Key())

		if value := gravityStore.Get(append(types.LastEventNonceByValidatorKey, oldOracleAddress.Bytes()...)); value != nil {
			ethStore.Set(append(crosschaintypes.LastEventNonceByOracleKey, oracleAddress.Bytes()...), value)
			gravityStore.Delete(append(types.LastEventNonceByValidatorKey, oldOracleAddress.Bytes()...))
		}

		if value := gravityStore.Get(append(types.LastEventBlockHeightByValidatorKey, oldOracleAddress.Bytes()...)); value != nil {
			ethStore.Set(append(crosschaintypes.LastEventBlockHeightByOracleKey, oracleAddress.Bytes()...), value)
			gravityStore.Delete(append(types.LastEventBlockHeightByValidatorKey, oldOracleAddress.Bytes()...))
		}

		chainOracle.Oracles = append(chainOracle.Oracles, oracle.OracleAddress)
	}

	// SetProposalOracle eth 0x38
	if len(chainOracle.Oracles) > 0 {
		ethStore.Set(crosschaintypes.ProposalOracleKey, cdc.MustMarshal(chainOracle))
	}
	// setLastTotalPower eth 0x39
	ethStore.Set(crosschaintypes.LastTotalPowerKey, cdc.MustMarshal(&sdk.IntProto{Int: totalPower}))

	// gravity 0x1 -> eth 0x12
	deletePrefixKey(gravityStore, types.EthAddressByValidatorKey)
	// delete 0x2
	deletePrefixKey(gravityStore, types.ValidatorByEthAddressKey)

	return oldOracleMap
}

func migrateOutgoingTxPool(cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore) {
	oldStore := prefix.NewStore(gravityStore, types.OutgoingTxPoolKey)
	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		// NOTE: migrate key, value is compatible
		var transact crosschaintypes.OutgoingTransferTx
		cdc.MustUnmarshal(oldStoreIter.Value(), &transact)

		ethStore.Set(crosschaintypes.GetOutgoingTxPoolKey(transact.Fee, transact.Id), oldStoreIter.Value())
		oldStore.Delete(oldStoreIter.Key())
	}

	oldStore2 := prefix.NewStore(gravityStore, types.SecondIndexOutgoingTxFeeKey)
	oldStoreIter2 := oldStore2.Iterator(nil, nil)
	defer oldStoreIter2.Close()

	for ; oldStoreIter2.Valid(); oldStoreIter2.Next() {
		oldStore2.Delete(oldStoreIter2.Key())
	}
}

func migrateOracleSetConfirm(cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore) {
	lastSlashedValsetNonce := sdk.BigEndianToUint64(gravityStore.Get(types.LastSlashedValsetNonce))

	oldStore := prefix.NewStore(gravityStore, types.ValsetConfirmKey)
	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var msg types.MsgValsetConfirm
		cdc.MustUnmarshal(oldStoreIter.Value(), &msg)

		key := oldStoreIter.Key()
		nonce := sdk.BigEndianToUint64(key[:8])
		if nonce != msg.Nonce {
			panic(fmt.Sprintf("invalid nonce, expect: %d, actual: %d", nonce, msg.Nonce))
		}
		bridgeAddr := key[8:]
		oracleAddr := ethStore.Get(crosschaintypes.GetOracleAddressByBridgerKey(bridgeAddr))
		if len(oracleAddr) != 20 {
			panic(fmt.Sprintf("invalid oracle address: %v", oracleAddr))
		}

		if msg.Nonce > lastSlashedValsetNonce {
			// Only the Confirm that is being processed is migrated
			ethStore.Set(crosschaintypes.GetOracleSetConfirmKey(msg.Nonce, oracleAddr),
				cdc.MustMarshal(&crosschaintypes.MsgOracleSetConfirm{
					Nonce:           msg.Nonce,
					BridgerAddress:  msg.Orchestrator,
					ExternalAddress: msg.EthAddress,
					Signature:       msg.Signature,
					ChainName:       ethtypes.ModuleName,
				}),
			)
		}

		oldStore.Delete(oldStoreIter.Key())
	}
}

func migrateOutgoingTxBatch(cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore) []*types.OutgoingTxBatch {
	oldStore := prefix.NewStore(gravityStore, types.OutgoingTxBatchKey)
	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	outgoingTxBatchs := make([]*types.OutgoingTxBatch, 0)
	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var batch types.OutgoingTxBatch
		cdc.MustUnmarshal(oldStoreIter.Value(), &batch)
		outgoingTxBatchs = append(outgoingTxBatchs, &batch)

		// NOTE: value is compatible
		ethStore.Set(crosschaintypes.GetOutgoingTxBatchKey(batch.TokenContract, batch.BatchNonce), oldStoreIter.Value())
		oldStore.Delete(oldStoreIter.Key())
	}
	return outgoingTxBatchs
}

func migrateConfirmBatch(cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore, outgoingTxBaths []*types.OutgoingTxBatch) {
	oldStore := prefix.NewStore(gravityStore, types.BatchConfirmKey)
	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var msg types.MsgConfirmBatch
		cdc.MustUnmarshal(oldStoreIter.Value(), &msg)

		key := oldStoreIter.Key()
		token := string(key[:len(msg.TokenContract)])
		if token != msg.TokenContract {
			panic(fmt.Sprintf("invalid token contract, expect: %s, actual: %s", token, msg.TokenContract))
		}
		nonce := sdk.BigEndianToUint64(key[len(msg.TokenContract) : len(msg.TokenContract)+8])
		if nonce != msg.Nonce {
			panic(fmt.Sprintf("invalid nonce, expect: %d, actual: %d", nonce, msg.Nonce))
		}
		bridgeAddr := key[len(msg.TokenContract)+8:]
		oracleAddr := ethStore.Get(crosschaintypes.GetOracleAddressByBridgerKey(bridgeAddr))
		if len(oracleAddr) != 20 {
			panic(fmt.Sprintf("invalid oracle address: %v", oracleAddr))
		}
		for _, bath := range outgoingTxBaths {
			// Only the Confirm that is being processed is migrated
			if bath.BatchNonce == msg.Nonce && bath.TokenContract == msg.TokenContract {
				ethStore.Set(crosschaintypes.GetBatchConfirmKey(msg.TokenContract, msg.Nonce, oracleAddr),
					cdc.MustMarshal(&crosschaintypes.MsgConfirmBatch{
						Nonce:           msg.Nonce,
						TokenContract:   msg.TokenContract,
						BridgerAddress:  msg.Orchestrator,
						ExternalAddress: msg.EthSigner,
						Signature:       msg.Signature,
						ChainName:       ethtypes.ModuleName,
					}),
				)
			}
		}

		oldStore.Delete(oldStoreIter.Key())
	}
}

func migrateLastObservedBlockHeight(cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore) {
	var msg types.LastObservedEthereumBlockHeight
	cdc.MustUnmarshal(gravityStore.Get(types.LastObservedEthereumBlockHeightKey), &msg)

	ethStore.Set(crosschaintypes.LastObservedBlockHeightKey,
		cdc.MustMarshal(&crosschaintypes.LastObservedBlockHeight{
			ExternalBlockHeight: msg.EthBlockHeight,
			BlockHeight:         msg.FxBlockHeight,
		}),
	)
	gravityStore.Delete(types.LastObservedEthereumBlockHeightKey)
}

func migrateAttestation(cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore, oracleMap map[string]string) {
	oldStore := prefix.NewStore(gravityStore, types.OracleAttestationKey)
	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var att types.Attestation
		cdc.MustUnmarshal(oldStoreIter.Value(), &att)

		for i := 0; i < len(att.Votes); i++ {
			if newOracle, ok := oracleMap[att.Votes[i]]; ok {
				att.Votes[i] = newOracle
			}
		}

		claim, err := types.UnpackAttestationClaim(cdc, &att)
		if err != nil {
			panic(err.Error())
		}
		var newClaim crosschaintypes.ExternalClaim
		switch c := claim.(type) {
		case *types.MsgDepositClaim:
			newClaim = &crosschaintypes.MsgSendToFxClaim{
				EventNonce:     c.EventNonce,
				BlockHeight:    c.BlockHeight,
				TokenContract:  c.TokenContract,
				Amount:         c.Amount,
				Sender:         c.EthSender,
				Receiver:       c.FxReceiver,
				TargetIbc:      c.TargetIbc,
				BridgerAddress: c.Orchestrator,
				ChainName:      ethtypes.ModuleName,
			}
		case *types.MsgWithdrawClaim:
			newClaim = &crosschaintypes.MsgSendToExternalClaim{
				EventNonce:     c.EventNonce,
				BlockHeight:    c.BlockHeight,
				BatchNonce:     c.BatchNonce,
				TokenContract:  c.TokenContract,
				BridgerAddress: c.Orchestrator,
				ChainName:      ethtypes.ModuleName,
			}
		case *types.MsgFxOriginatedTokenClaim:
			// ignore
		case *types.MsgValsetUpdatedClaim:
			myClaim := &crosschaintypes.MsgOracleSetUpdatedClaim{
				EventNonce:     c.EventNonce,
				BlockHeight:    c.BlockHeight,
				OracleSetNonce: c.ValsetNonce,
				Members:        make([]crosschaintypes.BridgeValidator, len(c.Members)),
				BridgerAddress: c.Orchestrator,
				ChainName:      ethtypes.ModuleName,
			}
			for i := 0; i < len(c.Members); i++ {
				myClaim.Members[i] = crosschaintypes.BridgeValidator{
					Power:           c.Members[i].Power,
					ExternalAddress: c.Members[i].EthAddress,
				}
			}
			newClaim = myClaim
		}

		anyMsg, err := codectypes.NewAnyWithValue(newClaim)
		if err != nil {
			panic(err.Error())
		}

		// new claim hash
		ethStore.Set(crosschaintypes.GetAttestationKey(newClaim.GetEventNonce(), newClaim.ClaimHash()),
			cdc.MustMarshal(&crosschaintypes.Attestation{
				Observed: att.Observed,
				Votes:    att.Votes,
				Height:   att.Height,
				Claim:    anyMsg,
			}),
		)
		oldStore.Delete(oldStoreIter.Key())
	}
}

func migrateBridgeToken(gravityStore, ethStore sdk.KVStore) {
	token := gravityStore.Get(append(types.DenomToERC20Key, []byte(fxtypes.DefaultDenom)...))
	if token == nil {
		return
	}
	ethStore.Set(crosschaintypes.GetTokenToDenomKey(fxtypes.DefaultDenom), token)
	ethStore.Set(crosschaintypes.GetDenomToTokenKey(string(token)), []byte(fxtypes.DefaultDenom))
	gravityStore.Delete(append(types.DenomToERC20Key, []byte(fxtypes.DefaultDenom)...))
	gravityStore.Delete(append(types.ERC20ToDenomKey, token...))
}

func MigrateBridgeTokenFromMetadatas(metadatas []banktypes.Metadata, ethStore sdk.KVStore) {
	for _, data := range metadatas {
		if len(data.DenomUnits) > 0 && len(data.DenomUnits[0].Aliases) > 0 {
			for i := 0; i < len(data.DenomUnits[0].Aliases); i++ {
				denom := data.DenomUnits[0].Aliases[i]
				if strings.HasPrefix(denom, ethtypes.ModuleName) {
					token := strings.TrimPrefix(denom, ethtypes.ModuleName)
					ethStore.Set(crosschaintypes.GetTokenToDenomKey(denom), []byte(token))
					ethStore.Set(crosschaintypes.GetDenomToTokenKey(token), []byte(denom))
				}
			}
		} else {
			if strings.HasPrefix(data.Base, ethtypes.ModuleName) {
				token := strings.TrimPrefix(data.Base, ethtypes.ModuleName)
				ethStore.Set(crosschaintypes.GetTokenToDenomKey(data.Base), []byte(token))
				ethStore.Set(crosschaintypes.GetDenomToTokenKey(token), []byte(data.Base))
			}
		}
	}
}

func deletePrefixKey(gravityStore sdk.KVStore, prefixKey []byte) {
	oldStore := prefix.NewStore(gravityStore, prefixKey)
	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		oldStore.Delete(oldStoreIter.Key())
	}
}
