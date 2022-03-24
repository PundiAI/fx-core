package types

import (
	fmt "fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store key
var (
	ParamStoreKeyEnableIntrarelayer       = []byte("EnableIntrarelayer")
	ParamStoreKeyEnableEVMHook            = []byte("EnableEVMHook")
	ParamStoreKeyIBCTransferTimeoutHeight = []byte("IBCTransferTimeoutHeight")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	enableIntrarelayer bool,
	enableEVMHook bool,
	ibcTransferTimeoutHeight uint64,
) Params {
	return Params{
		EnableIntrarelayer:       enableIntrarelayer,
		EnableEVMHook:            enableEVMHook,
		IbcTransferTimeoutHeight: ibcTransferTimeoutHeight,
	}
}

func DefaultParams() Params {
	return Params{
		EnableIntrarelayer:       true,
		EnableEVMHook:            true,
		IbcTransferTimeoutHeight: 20000,
	}
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateIbcTransferTimeoutHeight(i interface{}) error {
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyEnableIntrarelayer, &p.EnableIntrarelayer, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyEnableEVMHook, &p.EnableEVMHook, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyIBCTransferTimeoutHeight, &p.IbcTransferTimeoutHeight, validateIbcTransferTimeoutHeight),
	}
}

func (p Params) Validate() error {
	if p.IbcTransferTimeoutHeight == 0 {
		return fmt.Errorf("ibc transfer timeout height cannot be zero: %d", p.IbcTransferTimeoutHeight)
	}
	return nil
}
