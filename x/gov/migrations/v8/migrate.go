package v8

var (
	// Deprecated: do not use, remove in v8
	FxBaseParamsKeyPrefix = []byte("0x90")
	// Deprecated: do not use, remove in v8
	FxEGFParamsKey = []byte("0x91")
)

func GetRemovedStoreKeys() [][]byte {
	return [][]byte{FxBaseParamsKeyPrefix, FxEGFParamsKey}
}
