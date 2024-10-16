package precompile

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	crosschainkeeper "github.com/functionx/fx-core/v8/x/crosschain/keeper"
)

type Router struct {
	routes map[string]crosschainkeeper.Keeper
	sealed bool
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]crosschainkeeper.Keeper),
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

func (rtr *Router) AddRoute(module string, hook crosschainkeeper.Keeper) *Router {
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

func (rtr *Router) GetRoute(module string) (CrosschainKeeper, bool) {
	hook, found := rtr.routes[module]
	return hook, found
}
