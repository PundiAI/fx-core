package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "crosschain"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey is the module name router key
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName
)

var (
	OracleTotalDepositKey = []byte{0x11}
	// LastTotalPowerKey
	LastTotalPowerKey = []byte{0x39}

	// OraclesKey
	OraclesKey = []byte{0x12}

	// OracleAddressByExternalKey key external address -> value oracle address
	// OracleAddressByExternalKey indexes the external keys for an oracle address
	OracleAddressByExternalKey = []byte{0x13}

	// OracleAddressByOrchestratorKey key orchestrator address -> value oracle address
	// OracleAddressByOrchestratorKey indexes the external keys for an oracle address
	OracleAddressByOrchestratorKey = []byte{0x14}

	// OracleSetRequestKey indexes valset requests by nonce
	OracleSetRequestKey = []byte{0x15}

	// OracleSetConfirmKey indexes valset confirmations by nonce and the validator account address
	OracleSetConfirmKey = []byte{0x16}

	// OracleAttestationKey attestation details by nonce and validator address
	// An attestation can be thought of as the 'event to be executed' while
	// the Claims are an individual validator saying that they saw an event
	// occur the Attestation is 'the event' that multiple claims vote on and
	// eventually executes
	OracleAttestationKey = []byte{0x17}

	// OutgoingTXPoolKey indexes the last nonce for the outgoing tx pool
	OutgoingTXPoolKey = []byte{0x18}

	// SecondIndexOutgoingTXFeeKey indexes fee amounts by token contract address
	SecondIndexOutgoingTXFeeKey = []byte{0x19}

	// OutgoingTXBatchKey indexes outgoing tx batches under a nonce and token address
	OutgoingTXBatchKey = []byte{0x20}

	// OutgoingTXBatchBlockKey indexes outgoing tx batches under a block height and token address
	OutgoingTXBatchBlockKey = []byte{0x21}

	// BatchConfirmKey indexes validator confirmations by token contract address
	BatchConfirmKey = []byte{0x22}

	// LastEventNonceByValidatorKey indexes lateset event nonce by validator
	LastEventNonceByValidatorKey = []byte{0x23}

	// LastObservedEventNonceKey indexes the latest event nonce
	LastObservedEventNonceKey = []byte{0x24}

	// SequenceKeyPrefix indexes different txids
	SequenceKeyPrefix = []byte{0x25}

	// KeyLastTXPoolID indexes the lastTxPoolID
	KeyLastTXPoolID = append(SequenceKeyPrefix, []byte("lastTxPoolId")...)

	// KeyLastOutgoingBatchID indexes the lastBatchID
	KeyLastOutgoingBatchID = append(SequenceKeyPrefix, []byte("lastBatchId")...)

	// DenomToTokenKey prefixes the index of asset denoms to external token
	DenomToTokenKey = []byte{0x26}

	// TokenToDenomKey prefixes the index of assets external token to denoms
	TokenToDenomKey = []byte{0x27}

	// LastSlashedOracleSetNonce indexes the latest slashed valset nonce
	LastSlashedOracleSetNonce = []byte{0x28}

	// LatestOracleSetNonce indexes the latest valset nonce
	LatestOracleSetNonce = []byte{0x29}

	// LastSlashedBatchBlock indexes the latest slashed batch block height
	LastSlashedBatchBlock = []byte{0x30}

	// LastProposalBlockHeight indexes the last validator unbonding block height
	LastProposalBlockHeight = []byte{0x31}

	// LastObservedBlockHeightKey indexes the latest Ethereum block height
	LastObservedBlockHeightKey = []byte{0x32}

	// LastObservedOracleSetKey indexes the latest observed valset nonce
	LastObservedOracleSetKey = []byte{0x33}

	// KeyIbcSequenceHeight  indexes the gravity -> ibc sequence block height
	KeyIbcSequenceHeight = []byte{0x34}

	// LastEventBlockHeightByValidatorKey indexes lateset event blockHeight by validator
	LastEventBlockHeightByValidatorKey = []byte{0x35}

	// PastExternalSignatureCheckpointKey indexes eth signature checkpoints that have existed
	PastExternalSignatureCheckpointKey = []byte{0x36}

	// LastOracleSlashBlockHeight indexes the last oracle slash block height
	LastOracleSlashBlockHeight = []byte{0x37}

	// LastOracleSlashBlockHeight indexes the last oracle slash block height
	KeyChainOracles = []byte{0x38}
)

// GetOracleKey returns the following key format
func GetOracleKey(oracle sdk.AccAddress) []byte {
	return append(OraclesKey, oracle.Bytes()...)
}

// GetOracleAddressByOrchestratorKey returns the following key format
func GetOracleAddressByOrchestratorKey(orchestrator sdk.AccAddress) []byte {
	return append(OracleAddressByOrchestratorKey, orchestrator.Bytes()...)
}

// GetOracleAddressByExternalKey returns the following key format
func GetOracleAddressByExternalKey(externalAddress string) []byte {
	return append(OracleAddressByExternalKey, []byte(externalAddress)...)
}

// GetOracleSetKey returns the following key format
func GetOracleSetKey(nonce uint64) []byte {
	return append(OracleSetRequestKey, UInt64Bytes(nonce)...)
}

// GetOracleSetConfirmKey returns the following key format
// prefix           contract-address                BatchNonce                       external-address
// [0x16][0 0 0 0 0 0 0 1][fx1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetOracleSetConfirmKey(nonce uint64, oracleAddr sdk.AccAddress) []byte {
	return append(OracleSetConfirmKey, append(UInt64Bytes(nonce), oracleAddr.Bytes()...)...)
}

// GetAttestationKey returns the following key format
// prefix     nonce                             claim-details-hash
// [0x0][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]
// An attestation is an event multiple people are voting on, this function needs the claim
// details because each Attestation is aggregating all claims of a specific event, lets say
// validator X and validator y where making different claims about the same event nonce
// Note that the claim hash does NOT include the claimer address and only identifies an event
func GetAttestationKey(eventNonce uint64, claimHash []byte) []byte {
	key := make([]byte, len(OracleAttestationKey)+len(UInt64Bytes(0))+len(claimHash))
	copy(key[0:], OracleAttestationKey)
	copy(key[len(OracleAttestationKey):], UInt64Bytes(eventNonce))
	copy(key[len(OracleAttestationKey)+len(UInt64Bytes(0)):], claimHash)
	return key
}

// GetAttestationKeyWithHash returns the following key format
// prefix     nonce                             claim-details-hash
// [0x0][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]
// An attestation is an event multiple people are voting on, this function needs the claim
// details because each Attestation is aggregating all claims of a specific event, lets say
// validator X and validator y where making different claims about the same event nonce
// Note that the claim hash does NOT include the claimer address and only identifies an event
func GetAttestationKeyWithHash(eventNonce uint64, claimHash []byte) []byte {
	key := make([]byte, len(OracleAttestationKey)+len(UInt64Bytes(0))+len(claimHash))
	copy(key[0:], OracleAttestationKey)
	copy(key[len(OracleAttestationKey):], UInt64Bytes(eventNonce))
	copy(key[len(OracleAttestationKey)+len(UInt64Bytes(0)):], claimHash)
	return key
}

// GetOutgoingTxPoolContractPrefix returns the following key format
// prefix	feeContract
// [0x18][0xc783df8a850f42e7F7e57013759C285caa701eB6]
// This prefix is used for iterating over unbatched transactions for a given contract
func GetOutgoingTxPoolContractPrefix(contractAddress string) []byte {
	return append(OutgoingTXPoolKey, []byte(contractAddress)...)
}

//// GetOutgoingTxPoolKey returns the following key format
//func GetOutgoingTxPoolKey(id uint64) []byte {
//	return append(OutgoingTXPoolKey, sdk.Uint64ToBigEndian(id)...)
//}

// GetOutgoingTxPoolKey returns the following key format
// prefix	feeContract		feeAmount     id
// [0x18][0xc783df8a850f42e7F7e57013759C285caa701eB6][1000000000][0 0 0 0 0 0 0 1]
func GetOutgoingTxPoolKey(fee ExternalToken, id uint64) []byte {
	// sdkInts have a size limit of 255 bits or 32 bytes
	// therefore this will never panic and is always safe
	amount := make([]byte, 32)
	amount = fee.Amount.BigInt().FillBytes(amount)

	a := append(amount, UInt64Bytes(id)...)
	b := append([]byte(fee.Contract), a...)
	r := append(OutgoingTXPoolKey, b...)
	return r
}

// GetOutgoingTxBatchKey returns the following key format
func GetOutgoingTxBatchKey(tokenContract string, nonce uint64) []byte {
	return append(append(OutgoingTXBatchKey, []byte(tokenContract)...), UInt64Bytes(nonce)...)
}

// GetOutgoingTxBatchBlockKey returns the following key format
func GetOutgoingTxBatchBlockKey(block uint64) []byte {
	return append(OutgoingTXBatchBlockKey, UInt64Bytes(block)...)
}

// GetBatchConfirmKey returns the following key format
// prefix           contract-address                BatchNonce                       external-address
// [0x0][0xb4fA5979babd8Bb7e427157d0d353Cf205F43752][0 0 0 0 0 0 0 1][fx1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetBatchConfirmKey(tokenContract string, batchNonce uint64, oracleAddr sdk.AccAddress) []byte {
	a := append(UInt64Bytes(batchNonce), oracleAddr.Bytes()...)
	b := append([]byte(tokenContract), a...)
	c := append(BatchConfirmKey, b...)
	return c
}

// GetFeeSecondIndexKey returns the following key format
// prefix            contract-address            fee_amount
// [0x0][0xb4fA5979babd8Bb7e427157d0d353Cf205F43752][1000000000]
func GetFeeSecondIndexKey(fee ExternalToken) []byte {
	res := make([]byte, 1+ExternalContractAddressLen+32)
	// sdkInts have a size limit of 255 bits or 32 bytes
	// therefore this will never panic and is always safe
	amount := make([]byte, 32)
	amount = fee.Amount.BigInt().FillBytes(amount)
	copy(res[0:1], SecondIndexOutgoingTXFeeKey)
	copy(res[1:1+ExternalContractAddressLen], fee.Contract)
	copy(res[1+ExternalContractAddressLen:1+ExternalContractAddressLen+32], amount)
	return res
}

// GetLastEventNonceByOracleKey indexes lateset event nonce by validator
func GetLastEventNonceByOracleKey(validator sdk.AccAddress) []byte {
	return append(LastEventNonceByValidatorKey, validator.Bytes()...)
}

// GetIbcSequenceHeightKey [0xc1][sourcePort/sourceChannel/sequence]
func GetIbcSequenceHeightKey(sourcePort, sourceChannel string, sequence uint64) []byte {
	key := fmt.Sprintf("%s/%s/%d", sourcePort, sourceChannel, sequence)
	return append(KeyIbcSequenceHeight, []byte(key)...)
}

// GetLastEventBlockHeightByOracleKey indexes lateset event blockHeight by validator
func GetLastEventBlockHeightByOracleKey(validator sdk.AccAddress) []byte {
	return append(LastEventBlockHeightByValidatorKey, validator.Bytes()...)
}

func GetDenomToTokenKey(token string) []byte {
	return append(DenomToTokenKey, []byte(token)...)
}

func GetTokenToDenomKey(denom string) []byte {
	return append(TokenToDenomKey, []byte(denom)...)
}

// GetPastExternalSignatureCheckpointKey returns the following key format
// prefix    checkpoint
// [0x0][ checkpoint bytes ]
func GetPastExternalSignatureCheckpointKey(checkpoint []byte) []byte {
	return append(PastExternalSignatureCheckpointKey, checkpoint...)
}
