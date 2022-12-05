package v1

import "fmt"

var OracleTotalDepositKey = []byte{0x11}

// IbcSequenceHeightKey  indexes the gravity -> ibc sequence block height
var IbcSequenceHeightKey = []byte{0x34}

// GetIbcSequenceHeightKey returns the following key format
func GetIbcSequenceHeightKey(sourcePort, sourceChannel string, sequence uint64) []byte {
	key := fmt.Sprintf("%s/%s/%d", sourcePort, sourceChannel, sequence)
	return append(IbcSequenceHeightKey, []byte(key)...)
}
