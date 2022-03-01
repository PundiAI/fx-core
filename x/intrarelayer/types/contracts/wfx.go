package contracts

import (
	_ "embed"
	"encoding/json"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
)

var (
	//go:embed WFX.json
	WFXJSON []byte // nolint: golint

	// WFXContract is the compiled wfx contract
	WFXContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(WFXJSON, &WFXContract)
	if err != nil {
		panic(err)
	}

	if len(WFXContract.Bin) == 0 {
		panic("load contract failed")
	}
}
