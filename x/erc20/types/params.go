package types

import (
	"fmt"
	"time"
)

// NewParams creates a new Params object
func NewParams(enableErc20 bool, enableEVMHook bool, ibcTimeout time.Duration) Params {
	return Params{
		EnableErc20:   enableErc20,
		EnableEVMHook: enableEVMHook,
		IbcTimeout:    ibcTimeout,
	}
}

func DefaultParams() Params {
	return Params{
		EnableErc20:   true,
		EnableEVMHook: true,
		IbcTimeout:    12 * time.Hour,
	}
}

func (p *Params) Validate() error {
	if p.IbcTimeout <= 0 {
		return fmt.Errorf("ibc timeout cannot be 0")
	}
	return nil
}
