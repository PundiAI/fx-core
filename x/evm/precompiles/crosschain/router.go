package crosschain

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

type Precompile interface {
	TransferAfter(ctx sdk.Context, sender sdk.AccAddress, receive string, coins, fee sdk.Coin, originToken bool) error
	PrecompileCancelSendToExternal(ctx sdk.Context, txID uint64, sender sdk.AccAddress) (sdk.Coin, error)
	PrecompileIncreaseBridgeFee(ctx sdk.Context, txID uint64, sender sdk.AccAddress, addBridgeFee sdk.Coin) error
	PrecompileBridgeCall(ctx sdk.Context, sender, receiver, to common.Address, coins sdk.Coins, message []byte, value *big.Int, gasLimit uint64) (uint64, error)
}

type Router struct {
	routes map[string]Precompile
	sealed bool
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]Precompile),
	}
}

// Seal prevents the Router from any subsequent route handlers to be registered.
// Seal will panic if called more than once.
func (rtr *Router) Seal() {
	if rtr.sealed {
		panic("router already sealed")
	}
	rtr.sealed = true
}

func (rtr *Router) Sealed() bool {
	return rtr.sealed
}

func (rtr *Router) AddRoute(module string, hook Precompile) *Router {
	if rtr.sealed {
		panic(fmt.Sprintf("router sealed; cannot register %s route callbacks", module))
	}
	if !sdk.IsAlphaNumeric(module) {
		panic("route expressions can only contain alphanumeric characters")
	}
	if _, found := rtr.GetRoute(module); found {
		panic(fmt.Sprintf("route %s has already been registered", module))
	}

	rtr.routes[module] = hook
	return rtr
}

func (rtr *Router) GetRoute(module string) (Precompile, bool) {
	hook, found := rtr.routes[module]
	return hook, found
}
