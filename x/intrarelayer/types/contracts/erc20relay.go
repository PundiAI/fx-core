package contracts

import (
	_ "embed"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
	"github.com/functionx/fx-core/x/intrarelayer/types"
)

var (
	//go:embed ERC20Relay.json
	ERC20RelayJSON []byte // nolint: golint

	// ERC20RelayContract is the compiled erc20 contract
	ERC20RelayContract evmtypes.CompiledContract

	// ERC20RelayAddress is the irm module address
	ERC20RelayAddress common.Address
)

func init() {
	ERC20RelayAddress = types.ModuleAddress

	err := json.Unmarshal(ERC20RelayJSON, &ERC20RelayContract)
	if err != nil {
		panic(err)
	}

	if len(ERC20RelayContract.Bin) == 0 {
		panic("load contract failed")
	}
}
