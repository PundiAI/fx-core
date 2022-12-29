package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

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

// Sealed returns a boolean signifying if the Router is sealed or not.
func (rtr Router) Sealed() bool {
	return rtr.sealed
}

// AddRoute adds IBCModule for a given module name. It returns the Router
// so AddRoute calls can be linked. It will panic if the Router is sealed.
func (rtr *Router) AddRoute(module string, hook TransactionHook) *Router {
	if rtr.sealed {
		panic(fmt.Sprintf("router sealed; cannot register %s route callbacks", module))
	}
	if !sdk.IsAlphaNumeric(module) {
		panic("route expressions can only contain alphanumeric characters")
	}
	if rtr.HasRoute(module) {
		panic(fmt.Sprintf("route %s has already been registered", module))
	}

	rtr.routes[module] = hook
	return rtr
}

// HasRoute returns true if the Router has a module registered or false otherwise.
func (rtr *Router) HasRoute(module string) bool {
	_, ok := rtr.routes[module]
	return ok
}

// GetRoute returns a IBCModule for a given module.
func (rtr *Router) GetRoute(module string) (TransactionHook, bool) {
	if !rtr.HasRoute(module) {
		return nil, false
	}
	return rtr.routes[module], true
}

// RefundHook IBC transfer refund hook
var _ RefundHook = MultiRefundHook{}

// MultiRefundHook multi-refund hook
type MultiRefundHook []RefundHook

func (mrh MultiRefundHook) RefundAfter(ctx sdk.Context, sourcePort, sourceChannel string,
	sequence uint64, sender sdk.AccAddress, amount sdk.Coin) error {
	for i := range mrh {
		if err := mrh[i].RefundAfter(ctx, sourcePort, sourceChannel, sequence, sender, amount); err != nil {
			return sdkerrors.Wrapf(err, "Refund hook %T failed, error %s", mrh[i], err.Error())
		}
	}
	return nil
}

var _ AckHook = MultiAckHook{}

// MultiAckHook multi-ack hook
type MultiAckHook []AckHook

func (mrh MultiAckHook) AckAfter(ctx sdk.Context, sourcePort, sourceChannel string, sequence uint64) error {
	for i := range mrh {
		if err := mrh[i].AckAfter(ctx, sourcePort, sourceChannel, sequence); err != nil {
			return sdkerrors.Wrapf(err, "Ack hook %T failed, error %s", mrh[i], err.Error())
		}
	}
	return nil
}
