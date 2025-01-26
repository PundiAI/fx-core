package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "crosschain"

	BridgeCallSender = "bridge_call"

	BridgeFeeCollectorName = "bridge_fee_collector"
)

var (
	// OracleKey key oracle address -> Oracle
	OracleKey = []byte{0x12}

	// OracleAddressByExternalKey key external address -> value oracle address
	OracleAddressByExternalKey = []byte{0x13}

	// OracleAddressByBridgerKey key bridger address -> value oracle address
	OracleAddressByBridgerKey = []byte{0x14}

	// OracleSetRequestKey indexes oracle set requests by nonce
	OracleSetRequestKey = []byte{0x15}

	// OracleSetConfirmKey indexes oracle set confirmations by nonce and the validator account address
	OracleSetConfirmKey = []byte{0x16}

	// OracleAttestationKey attestation details by nonce and validator address
	// An attestation can be thought of as the 'event to be executed' while
	// the Claims are an individual validator saying that they saw an event
	// occur the Attestation is 'the event' that multiple claims vote on and
	// eventually executes
	OracleAttestationKey = []byte{0x17}

	// OutgoingTxBatchKey indexes outgoing tx batches under a nonce and token address
	OutgoingTxBatchKey = []byte{0x20}

	// OutgoingTxBatchBlockKey indexes outgoing tx batches under a block height and token address
	OutgoingTxBatchBlockKey = []byte{0x21}

	// BatchConfirmKey indexes oracle confirmations by token contract address
	BatchConfirmKey = []byte{0x22}

	// LastEventNonceByOracleKey indexes latest event nonce by oracle
	LastEventNonceByOracleKey = []byte{0x23}

	// LastObservedEventNonceKey indexes the latest event nonce
	LastObservedEventNonceKey = []byte{0x24}

	// SequenceKeyPrefix indexes different txIds
	SequenceKeyPrefix = []byte{0x25}

	// KeyLastTxPoolID indexes the lastTxPoolID
	KeyLastTxPoolID = append(SequenceKeyPrefix, []byte("lastTxPoolId")...)

	// KeyLastOutgoingBatchID indexes the lastBatchID
	KeyLastOutgoingBatchID = append(SequenceKeyPrefix, []byte("lastBatchId")...)

	// LastSlashedOracleSetNonce indexes the latest slashed oracleSet nonce
	LastSlashedOracleSetNonce = []byte{0x28}

	// LatestOracleSetNonce indexes the latest oracleSet nonce
	LatestOracleSetNonce = []byte{0x29}

	// LastSlashedBatchBlock indexes the latest slashed batch block height
	LastSlashedBatchBlock = []byte{0x30}

	// Deprecated: LastProposalBlockHeight
	// LastProposalBlockHeight = []byte{0x31}

	// LastObservedBlockHeightKey indexes the latest observed external block height
	LastObservedBlockHeightKey = []byte{0x32}

	// LastObservedOracleSetKey indexes the latest observed OracleSet nonce
	LastObservedOracleSetKey = []byte{0x33}

	// LastEventBlockHeightByOracleKey indexes latest event blockHeight by oracle
	LastEventBlockHeightByOracleKey = []byte{0x35}

	// Deprecated: PastExternalSignatureCheckpointKey indexes eth signature checkpoints that have existed
	// PastExternalSignatureCheckpointKey = []byte{0x36}

	// LastOracleSlashBlockHeight indexes the last oracle slash block height
	LastOracleSlashBlockHeight = []byte{0x37}

	// ProposalOracleKey -> value ProposalOracle
	ProposalOracleKey = []byte{0x38}

	// LastTotalPowerKey oracle set total power
	LastTotalPowerKey = []byte{0x39}

	// ParamsKey is the prefix for params key
	ParamsKey = []byte{0x40}

	// Deprecated: OutgoingTxRelationKey outgoing tx with evm
	// OutgoingTxRelationKey = []byte{0x41}

	BridgeCallConfirmKey = []byte{0x45}

	// LastSlashedBridgeCallNonce indexes the latest slashed bridge call nonce
	LastSlashedBridgeCallNonce = []byte{0x46}

	KeyLastBridgeCallID = append(SequenceKeyPrefix, []byte("bridgeCallId")...)

	OutgoingBridgeCallNonceKey           = []byte{0x48}
	OutgoingBridgeCallAddressAndNonceKey = []byte{0x49}

	PendingExecuteClaimKey = []byte{0x54}

	BridgeCallQuoteKey = []byte{0x55}
)

// GetOracleKey returns the following key format
func GetOracleKey(oracle sdk.AccAddress) []byte {
	return append(OracleKey, oracle.Bytes()...)
}

// GetOracleAddressByBridgerKey returns the following key format
func GetOracleAddressByBridgerKey(bridger sdk.AccAddress) []byte {
	return append(OracleAddressByBridgerKey, bridger.Bytes()...)
}

// GetOracleAddressByExternalKey returns the following key format
func GetOracleAddressByExternalKey(externalAddress string) []byte {
	return append(OracleAddressByExternalKey, []byte(externalAddress)...)
}

// GetOracleSetKey returns the following key format
func GetOracleSetKey(nonce uint64) []byte {
	return append(OracleSetRequestKey, sdk.Uint64ToBigEndian(nonce)...)
}

// GetOracleSetConfirmKey returns the following key format
func GetOracleSetConfirmKey(nonce uint64, oracleAddr sdk.AccAddress) []byte {
	return append(OracleSetConfirmKey, append(sdk.Uint64ToBigEndian(nonce), oracleAddr.Bytes()...)...)
}

// GetAttestationKey returns the following key format
// An attestation is an event multiple people are voting on, this function needs the claim
// details because each Attestation is aggregating all claims of a specific event, lets say
// validator X and validator y where making different claims about the same event nonce
// Note that the claim hash does NOT include the claimer address and only identifies an event
func GetAttestationKey(eventNonce uint64, claimHash []byte) []byte {
	return append(OracleAttestationKey, append(sdk.Uint64ToBigEndian(eventNonce), claimHash...)...)
}

// GetOutgoingTxBatchKey returns the following key format
func GetOutgoingTxBatchKey(tokenContract string, batchNonce uint64) []byte {
	return append(append(OutgoingTxBatchKey, []byte(tokenContract)...), sdk.Uint64ToBigEndian(batchNonce)...)
}

// GetOutgoingTxBatchBlockKey returns the following key format
func GetOutgoingTxBatchBlockKey(blockHeight, batchNonce uint64) []byte {
	return append(append(OutgoingTxBatchBlockKey, sdk.Uint64ToBigEndian(blockHeight)...), sdk.Uint64ToBigEndian(batchNonce)...)
}

// GetBatchConfirmKey returns the following key format
func GetBatchConfirmKey(tokenContract string, batchNonce uint64, oracleAddr sdk.AccAddress) []byte {
	return append(BatchConfirmKey, append([]byte(tokenContract), append(sdk.Uint64ToBigEndian(batchNonce), oracleAddr.Bytes()...)...)...)
}

// GetLastEventNonceByOracleKey returns the following key format
func GetLastEventNonceByOracleKey(oracleAddr sdk.AccAddress) []byte {
	return append(LastEventNonceByOracleKey, oracleAddr.Bytes()...)
}

// GetLastEventBlockHeightByOracleKey returns the following key format
func GetLastEventBlockHeightByOracleKey(oracleAddr sdk.AccAddress) []byte {
	return append(LastEventBlockHeightByOracleKey, oracleAddr.Bytes()...)
}

func GetBridgeCallConfirmKey(nonce uint64, addr sdk.AccAddress) []byte {
	return append(BridgeCallConfirmKey, append(sdk.Uint64ToBigEndian(nonce), addr.Bytes()...)...)
}

func GetBridgeCallConfirmNonceKey(nonce uint64) []byte {
	return append(BridgeCallConfirmKey, sdk.Uint64ToBigEndian(nonce)...)
}

func GetOutgoingBridgeCallNonceKey(id uint64) []byte {
	return append(OutgoingBridgeCallNonceKey, sdk.Uint64ToBigEndian(id)...)
}

func GetOutgoingBridgeCallAddressAndNonceKey(address string, nonce uint64) []byte {
	return append(GetOutgoingBridgeCallAddressKey(address), sdk.Uint64ToBigEndian(nonce)...)
}

func GetOutgoingBridgeCallAddressKey(address string) []byte {
	return append(OutgoingBridgeCallAddressAndNonceKey, []byte(address)...)
}

func ParseOutgoingBridgeCallNonce(key []byte, address string) (nonce uint64) {
	addrNonce := bytes.TrimPrefix(key, OutgoingBridgeCallAddressAndNonceKey)
	nonceBytes := bytes.TrimPrefix(addrNonce, []byte(address))
	return sdk.BigEndianToUint64(nonceBytes)
}

func GetPendingExecuteClaimKey(nonce uint64) []byte {
	return append(PendingExecuteClaimKey, sdk.Uint64ToBigEndian(nonce)...)
}

func GetBridgeCallQuoteKey(nonce uint64) []byte {
	return append(BridgeCallQuoteKey, sdk.Uint64ToBigEndian(nonce)...)
}
