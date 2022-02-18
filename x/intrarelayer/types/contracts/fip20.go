package contracts

import (
	_ "embed"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
	"github.com/functionx/fx-core/x/intrarelayer/types"
)

var (
	//go:embed FIP20.json
	FIP20JSON []byte // nolint: golint

	// FIP20Contract is the compiled fip20 contract
	FIP20Contract evmtypes.CompiledContract

	// FIP20Address is the irm module address
	FIP20Address common.Address
)

func init() {
	FIP20Address = types.ModuleAddress

	err := json.Unmarshal(FIP20JSON, &FIP20Contract)
	if err != nil {
		panic(err)
	}

	if len(FIP20Contract.Bin) == 0 {
		panic("load contract failed")
	}
}
