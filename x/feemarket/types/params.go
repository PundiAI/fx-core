package types

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/params"
)

var _ paramtypes.ParamSet = &Params{}

// Parameter keys
var (
	ParamStoreKeyBaseFeeChangeDenominator = []byte("BaseFeeChangeDenominator")
	ParamStoreKeyElasticityMultiplier     = []byte("ElasticityMultiplier")
	ParamStoreKeyBaseFee                  = []byte("BaseFee")
	ParamStoreKeyMinBaseFee               = []byte("MinBaseFee")
	ParamStoreKeyMaxBaseFee               = []byte("MaxBaseFee")
	ParamStoreKeyMaxGas                   = []byte("MaxGas")
)

var (
	MinBaseFee = sdk.NewInt(500 * 1e9)
	MaxBaseFee = sdk.NewIntFromUint64(math.MaxUint64 - 1).Mul(sdk.NewInt(1e9))
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(baseFeeChangeDenom, elasticityMultiplier uint32, baseFee uint64,
	minBaseFee, maxBaseFee sdk.Int, maxGas uint64) Params {
	return Params{
		BaseFeeChangeDenominator: baseFeeChangeDenom,
		ElasticityMultiplier:     elasticityMultiplier,
		BaseFee:                  sdk.NewIntFromUint64(baseFee),
		MinBaseFee:               minBaseFee,
		MaxBaseFee:               maxBaseFee,
		MaxGas:                   sdk.NewIntFromUint64(maxGas),
	}
}

// DefaultParams returns default evm parameters
func DefaultParams() Params {
	return Params{
		BaseFeeChangeDenominator: params.BaseFeeChangeDenominator,
		ElasticityMultiplier:     params.ElasticityMultiplier,
		BaseFee:                  MinBaseFee,
		MinBaseFee:               MinBaseFee,
		MaxBaseFee:               MaxBaseFee,
		MaxGas:                   sdk.NewIntFromUint64(3 * 1e7), //default use block max gas
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyBaseFeeChangeDenominator, &p.BaseFeeChangeDenominator, validateBaseFeeChangeDenominator),
		paramtypes.NewParamSetPair(ParamStoreKeyElasticityMultiplier, &p.ElasticityMultiplier, validateElasticityMultiplier),
		paramtypes.NewParamSetPair(ParamStoreKeyBaseFee, &p.BaseFee, validateBaseFee),
		paramtypes.NewParamSetPair(ParamStoreKeyMinBaseFee, &p.MinBaseFee, validateBaseFee),
		paramtypes.NewParamSetPair(ParamStoreKeyMaxBaseFee, &p.MaxBaseFee, validateBaseFee),
		paramtypes.NewParamSetPair(ParamStoreKeyMaxGas, &p.MaxGas, validateMaxGas),
	}
}

// Validate performs basic validation on fee market parameters.
func (p Params) Validate() error {
	if p.BaseFeeChangeDenominator == 0 {
		return fmt.Errorf("base fee change denominator cannot be 0")
	}
	if p.ElasticityMultiplier == 0 {
		return fmt.Errorf("elasticity multiplier cannot be 0")
	}

	if p.BaseFee.IsNegative() {
		return fmt.Errorf("initial base fee cannot be negative: %s", p.BaseFee)
	}

	if p.MinBaseFee.IsNegative() {
		return fmt.Errorf("min base fee cannot be negative: %s", p.MinBaseFee)
	}

	if p.MaxBaseFee.IsNegative() {
		return fmt.Errorf("max base fee cannot be negative: %s", p.MaxBaseFee)
	}

	if p.MaxBaseFee.GT(sdk.ZeroInt()) && p.MaxBaseFee.LT(MinBaseFee) {
		return fmt.Errorf("if max base fee gt 0, max base fee(%s) must be gte min base fee(%s)", p.MaxBaseFee, p.MinBaseFee)
	}

	if p.MaxGas.IsNegative() {
		return fmt.Errorf("max gas cannot be negative: %s", p.MaxGas)
	}

	return nil
}

func validateBaseFeeChangeDenominator(i interface{}) error {
	value, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value == 0 {
		return fmt.Errorf("base fee change denominator cannot be 0")
	}

	return nil
}

func validateElasticityMultiplier(i interface{}) error {
	value, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if value == 0 {
		return fmt.Errorf("elasticity multiplier cannot be 0")
	}
	return nil
}

func validateBaseFee(i interface{}) error {
	value, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value.IsNegative() {
		return fmt.Errorf("base fee cannot be negative")
	}

	return nil
}

func validateMaxGas(i interface{}) error {
	value, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value.IsNegative() {
		return fmt.Errorf("max gas cannot be negative")
	}
	return nil
}
