package types

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

type TransactionHook interface {
	TransferAfter(ctx sdk.Context, sender sdk.AccAddress, receive string, coins, fee sdk.Coin, originToken bool) error
	PrecompileCancelSendToExternal(ctx sdk.Context, txID uint64, sender sdk.AccAddress) (sdk.Coin, error)
	PrecompileIncreaseBridgeFee(ctx sdk.Context, txID uint64, sender sdk.AccAddress, addBridgeFee sdk.Coin) error
	PrecompileBridgeCall(ctx sdk.Context, dstChainId string, gasLimit uint64, sender, receiver, to common.Address, asset, message []byte, value *big.Int) (uint64, error)
}

type Router struct {
	routes map[string]TransactionHook
	sealed bool
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]TransactionHook),
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

func (rtr *Router) AddRoute(module string, hook TransactionHook) *Router {
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

func (rtr *Router) GetRoute(module string) (TransactionHook, bool) {
	hook, found := rtr.routes[module]
	return hook, found
}
