package v042

import "fmt"

// KeyIbcSequenceHeight  indexes the gravity -> ibc sequence block height
var KeyIbcSequenceHeight = []byte{0x34}

//GetIbcSequenceHeightKey returns the following key format
func GetIbcSequenceHeightKey(sourcePort, sourceChannel string, sequence uint64) []byte {
	key := fmt.Sprintf("%s/%s/%d", sourcePort, sourceChannel, sequence)
	return append(KeyIbcSequenceHeight, []byte(key)...)
}
