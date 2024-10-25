package v8

var (
	KeyPrefixTokenPair        = []byte{0x01}
	KeyPrefixTokenPairByERC20 = []byte{0x02}
	KeyPrefixTokenPairByDenom = []byte{0x03}
	KeyPrefixIBCTransfer      = []byte{0x04}
	KeyPrefixAliasDenom       = []byte{0x05}
	ParamsKey                 = []byte{0x06}
	KeyPrefixOutgoingTransfer = []byte{0x07}
)

func GetRemovedStoreKeys() [][]byte {
	return [][]byte{
		KeyPrefixTokenPair, KeyPrefixTokenPairByERC20, KeyPrefixTokenPairByDenom, KeyPrefixAliasDenom,
	}
}
