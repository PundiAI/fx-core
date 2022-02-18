package types

import (
	fmt "fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store key
var (
	ParamStoreKeyEnableIntrarelayer       = []byte("EnableIntrarelayer")
	ParamStoreKeyTokenPairVotingPeriod    = []byte("TokenPairVotingPeriod")
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
	votingPeriod time.Duration,
	enableEVMHook bool,
	ibcTransferTimeoutHeight uint64,
) Params {
	return Params{
		EnableIntrarelayer:       enableIntrarelayer,
		TokenPairVotingPeriod:    votingPeriod,
		EnableEVMHook:            enableEVMHook,
		IbcTransferTimeoutHeight: ibcTransferTimeoutHeight,
	}
}

func DefaultParams() Params {
	return Params{
		EnableIntrarelayer:       true,
		TokenPairVotingPeriod:    time.Hour * 24 * 14,
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

func validatePeriod(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("voting period must be positive: %s", v)
	}

	return nil
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyEnableIntrarelayer, &p.EnableIntrarelayer, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyTokenPairVotingPeriod, &p.TokenPairVotingPeriod, validatePeriod),
		paramtypes.NewParamSetPair(ParamStoreKeyEnableEVMHook, &p.EnableEVMHook, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyIBCTransferTimeoutHeight, &p.IbcTransferTimeoutHeight, validateIbcTransferTimeoutHeight),
	}
}

func (p Params) Validate() error {
	if p.TokenPairVotingPeriod <= 0 {
		return fmt.Errorf("voting period must be positive: %d", p.TokenPairVotingPeriod)
	}

	return nil
}
