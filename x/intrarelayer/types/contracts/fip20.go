package contracts

import (
	_ "embed"
	"encoding/json"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
)

var (
	//go:embed FIP20.json
	FIP20JSON []byte // nolint: golint

	// FIP20Contract is the compiled fip20 contract
	FIP20Contract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(FIP20JSON, &FIP20Contract)
	if err != nil {
		panic(err)
	}

	if len(FIP20Contract.Bin) == 0 {
		panic("load contract failed")
	}
}
