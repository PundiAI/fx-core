package types

const (
	// ModuleName is the name of the module
	ModuleName = "gravity"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey is the module name router key
	RouterKey = ModuleName
)

var (
	// Deprecated: EthAddressByValidatorKey indexes cosmos validator account addresses
	EthAddressByValidatorKey = []byte{0x1}

	// Deprecated: ValidatorByEthAddressKey indexes ethereum addresses
	ValidatorByEthAddressKey = []byte{0x2}

	// Deprecated: ValsetRequestKey indexes valset requests by nonce
	ValsetRequestKey = []byte{0x3}

	// Deprecated: ValsetConfirmKey indexes valset confirmations by nonce and the validator account address
	ValsetConfirmKey = []byte{0x4}

	// Deprecated: OracleAttestationKey attestation details by nonce and validator address
	// An attestation can be thought of as the 'event to be executed' while
	// the Claims are an individual validator saying that they saw an event
	// occur the Attestation is 'the event' that multiple claims vote on and
	// eventually executes
	OracleAttestationKey = []byte{0x5}

	// Deprecated: OutgoingTxPoolKey indexes the last nonce for the outgoing tx pool
	OutgoingTxPoolKey = []byte{0x6}

	// Deprecated: SecondIndexOutgoingTxFeeKey indexes fee amounts by token contract address
	SecondIndexOutgoingTxFeeKey = []byte{0x7}

	// Deprecated: OutgoingTxBatchKey indexes outgoing tx batches under a nonce and token address
	OutgoingTxBatchKey = []byte{0x8}

	// Deprecated: OutgoingTxBatchBlockKey indexes outgoing tx batches under a block height and token address
	OutgoingTxBatchBlockKey = []byte{0x9}

	// Deprecated: BatchConfirmKey indexes validator confirmations by token contract address
	BatchConfirmKey = []byte{0xa}

	// Deprecated: LastEventNonceByValidatorKey indexes latest event nonce by validator
	LastEventNonceByValidatorKey = []byte{0xb}

	// Deprecated: LastObservedEventNonceKey indexes the latest event nonce
	LastObservedEventNonceKey = []byte{0xc}

	// Deprecated: SequenceKeyPrefix indexes different txids
	SequenceKeyPrefix = []byte{0xd}

	// Deprecated: KeyLastTXPoolID indexes the lastTxPoolID
	KeyLastTXPoolID = append(SequenceKeyPrefix, []byte("lastTxPoolId")...)

	// Deprecated: KeyLastOutgoingBatchID indexes the lastBatchID
	KeyLastOutgoingBatchID = append(SequenceKeyPrefix, []byte("lastBatchId")...)

	// Deprecated: ValidatorAddressByOrchestratorAddress indexes the validator keys for an orchestrator
	ValidatorAddressByOrchestratorAddress = []byte{0xe}

	// Deprecated: DenomToERC20Key prefixes the index of Cosmos originated asset denoms to ERC20s
	DenomToERC20Key = []byte{0xf}

	// Deprecated: ERC20ToDenomKey prefixes the index of Cosmos originated assets ERC20s to denoms
	ERC20ToDenomKey = []byte{0x10}

	// Deprecated: LastSlashedValsetNonce indexes the latest slashed valset nonce
	LastSlashedValsetNonce = []byte{0x11}

	// Deprecated: LatestValsetNonce indexes the latest valset nonce
	LatestValsetNonce = []byte{0x12}

	// Deprecated: LastSlashedBatchBlock indexes the latest slashed batch block height
	LastSlashedBatchBlock = []byte{0x13}

	// Deprecated: LastUnBondingBlockHeight indexes the last validator unbonding block height
	LastUnBondingBlockHeight = []byte{0x14}

	// Deprecated: LastObservedEthereumBlockHeightKey indexes the latest Ethereum block height
	LastObservedEthereumBlockHeightKey = []byte{0x15}

	// Deprecated: LastObservedValsetKey indexes the latest observed valset nonce
	LastObservedValsetKey = []byte{0x16}

	// Deprecated: IbcSequenceHeightKey  indexes the gravity -> ibc sequence block height
	IbcSequenceHeightKey = []byte{0x17}

	// Deprecated: LastEventBlockHeightByValidatorKey indexes latest event blockHeight by validator
	LastEventBlockHeightByValidatorKey = []byte{0x18}
)
