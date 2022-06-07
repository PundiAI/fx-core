package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "gravity"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey is the module name router key
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)

var (
	// EthAddressByValidatorKey indexes cosmos validator account addresses
	EthAddressByValidatorKey = []byte{0x1}

	// ValidatorByEthAddressKey indexes ethereum addresses
	ValidatorByEthAddressKey = []byte{0x2}

	// ValsetRequestKey indexes valset requests by nonce
	ValsetRequestKey = []byte{0x3}

	// ValsetConfirmKey indexes valset confirmations by nonce and the validator account address
	ValsetConfirmKey = []byte{0x4}

	// OracleAttestationKey attestation details by nonce and validator address
	// An attestation can be thought of as the 'event to be executed' while
	// the Claims are an individual validator saying that they saw an event
	// occur the Attestation is 'the event' that multiple claims vote on and
	// eventually executes
	OracleAttestationKey = []byte{0x5}

	// OutgoingTXPoolKey indexes the last nonce for the outgoing tx pool
	OutgoingTXPoolKey = []byte{0x6}

	// SecondIndexOutgoingTXFeeKey indexes fee amounts by token contract address
	SecondIndexOutgoingTXFeeKey = []byte{0x7}

	// OutgoingTXBatchKey indexes outgoing tx batches under a nonce and token address
	OutgoingTXBatchKey = []byte{0x8}

	// OutgoingTXBatchBlockKey indexes outgoing tx batches under a block height and token address
	OutgoingTXBatchBlockKey = []byte{0x9}

	// BatchConfirmKey indexes validator confirmations by token contract address
	BatchConfirmKey = []byte{0xa}

	// LastEventNonceByValidatorKey indexes lateset event nonce by validator
	LastEventNonceByValidatorKey = []byte{0xb}

	// LastObservedEventNonceKey indexes the latest event nonce
	LastObservedEventNonceKey = []byte{0xc}

	// SequenceKeyPrefix indexes different txids
	SequenceKeyPrefix = []byte{0xd}

	// KeyLastTXPoolID indexes the lastTxPoolID
	KeyLastTXPoolID = append(SequenceKeyPrefix, []byte("lastTxPoolId")...)

	// KeyLastOutgoingBatchID indexes the lastBatchID
	KeyLastOutgoingBatchID = append(SequenceKeyPrefix, []byte("lastBatchId")...)

	// KeyOrchestratorAddress indexes the validator keys for an orchestrator
	KeyOrchestratorAddress = []byte{0xe}

	// DenomToERC20Key prefixes the index of Cosmos originated asset denoms to ERC20s
	DenomToERC20Key = []byte{0xf}

	// ERC20ToDenomKey prefixes the index of Cosmos originated assets ERC20s to denoms
	ERC20ToDenomKey = []byte{0x10}

	// LastSlashedValsetNonce indexes the latest slashed valset nonce
	LastSlashedValsetNonce = []byte{0x11}

	// LatestValsetNonce indexes the latest valset nonce
	LatestValsetNonce = []byte{0x12}

	// LastSlashedBatchBlock indexes the latest slashed batch block height
	LastSlashedBatchBlock = []byte{0x13}

	// LastProposalBlockHeight indexes the last validator unbonding block height
	LastUnBondingBlockHeight = []byte{0x14}

	// LastObservedEthereumBlockHeightKey indexes the latest Ethereum block height
	LastObservedEthereumBlockHeightKey = []byte{0x15}

	// LastObservedValsetKey indexes the latest observed valset nonce
	LastObservedValsetKey = []byte{0x16}

	// KeyIbcSequenceHeight  indexes the gravity -> ibc sequence block height
	// DEPRECATED: delete by v2
	KeyIbcSequenceHeight = []byte{0x17}

	// LastEventBlockHeightByValidatorKey indexes lateset event blockHeight by validator
	LastEventBlockHeightByValidatorKey = []byte{0x18}
)

// GetOrchestratorAddressKey returns the following key format
func GetOrchestratorAddressKey(orc sdk.AccAddress) []byte {
	return append(KeyOrchestratorAddress, orc.Bytes()...)
}

// GetEthAddressByValidatorKey returns the following key format
func GetEthAddressByValidatorKey(validator sdk.ValAddress) []byte {
	return append(EthAddressByValidatorKey, validator.Bytes()...)
}

// GetValidatorByEthAddressKey returns the following key format
func GetValidatorByEthAddressKey(ethAddress string) []byte {
	return append(ValidatorByEthAddressKey, []byte(ethAddress)...)
}

// GetValsetKey returns the following key format
func GetValsetKey(nonce uint64) []byte {
	return append(ValsetRequestKey, UInt64Bytes(nonce)...)
}

// GetValsetConfirmKey returns the following key format
func GetValsetConfirmKey(nonce uint64, validator sdk.AccAddress) []byte {
	return append(ValsetConfirmKey, append(UInt64Bytes(nonce), validator.Bytes()...)...)
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

// GetOutgoingTxPoolKey returns the following key format
func GetOutgoingTxPoolKey(id uint64) []byte {
	return append(OutgoingTXPoolKey, sdk.Uint64ToBigEndian(id)...)
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
// prefix           eth-contract-address                BatchNonce                       Validator-address
// [0x0][0xb4fA5979babd8Bb7e427157d0d353Cf205F43752][0 0 0 0 0 0 0 1][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func GetBatchConfirmKey(tokenContract string, batchNonce uint64, validator sdk.AccAddress) []byte {
	a := append(UInt64Bytes(batchNonce), validator.Bytes()...)
	b := append([]byte(tokenContract), a...)
	c := append(BatchConfirmKey, b...)
	return c
}

// GetFeeSecondIndexKey returns the following key format
// prefix            eth-contract-address            fee_amount
// [0x0][0xb4fA5979babd8Bb7e427157d0d353Cf205F43752][1000000000]
func GetFeeSecondIndexKey(fee ERC20Token) []byte {
	res := make([]byte, 1+ETHContractAddressLen+32)
	// sdkInts have a size limit of 255 bits or 32 bytes
	// therefore this will never panic and is always safe
	amount := make([]byte, 32)
	amount = fee.Amount.BigInt().FillBytes(amount)
	copy(res[0:1], SecondIndexOutgoingTXFeeKey)
	copy(res[1:1+ETHContractAddressLen], fee.Contract)
	copy(res[1+ETHContractAddressLen:1+ETHContractAddressLen+32], amount)
	return res
}

// GetLastEventNonceByOracleKey indexes lateset event nonce by validator
func GetLastEventNonceByValidatorKey(validator sdk.ValAddress) []byte {
	return append(LastEventNonceByValidatorKey, validator.Bytes()...)
}

// GetDenomToERC20Key denom -> erc20
func GetDenomToERC20Key(denom string) []byte {
	return append(DenomToERC20Key, []byte(denom)...)
}

// GetERC20ToDenomKey erc20 -> denom
func GetERC20ToDenomKey(erc20 string) []byte {
	return append(ERC20ToDenomKey, []byte(erc20)...)
}

//GetIbcSequenceHeightKey [0xc1][sourcePort/sourceChannel/sequence]
// DEPRECATED: delete by v2
func GetIbcSequenceHeightKey(sourcePort, sourceChannel string, sequence uint64) []byte {
	key := fmt.Sprintf("%s/%s/%d", sourcePort, sourceChannel, sequence)
	return append(KeyIbcSequenceHeight, []byte(key)...)
}

// GetLastEventBlockHeightByValidatorKey indexes lateset event blockHeight by validator
func GetLastEventBlockHeightByValidatorKey(validator sdk.ValAddress) []byte {
	return append(LastEventBlockHeightByValidatorKey, validator.Bytes()...)
}
