package v8

var (
	// Deprecated: do not use, remove in v8
	ValidatorOperatorKey = []byte{0x91}
	// Deprecated: do not use, remove in v8
	ConsensusPubKey = []byte{0x92}
	// Deprecated: do not use, remove in v8
	ConsensusProcessKey = []byte{0x93}
)

func GetRemovedStoreKeys() [][]byte {
	return [][]byte{ValidatorOperatorKey, ConsensusPubKey, ConsensusProcessKey}
}
