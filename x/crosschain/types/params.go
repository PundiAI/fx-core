package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
)

const (
	MaxOracleSize                  = 100
	DefaultOracleDelegateThreshold = 10
	OutgoingTxBatchSize            = 100
	MaxResults                     = 100
	MaxOracleSetRequestsResults    = 5
	MaxKeepEventSize               = 100
)

var (
	OracleDelegateDenom = fxtypes.DefaultDenom
	NativeDenom         = fxtypes.DefaultDenom
)

var (
	// AttestationVotesPowerThreshold threshold of votes power to succeed
	AttestationVotesPowerThreshold = sdkmath.NewInt(66)

	AttestationProposalOracleChangePowerThreshold = sdkmath.NewInt(30)

	// ParamsStoreKeyGravityID stores the gravity id
	ParamsStoreKeyGravityID = []byte("GravityID")

	// ParamsStoreKeyAverageBlockTime stores the signed blocks window
	ParamsStoreKeyAverageBlockTime = []byte("AverageBlockTime")

	// ParamsStoreKeyExternalBatchTimeout stores the signed blocks window
	ParamsStoreKeyExternalBatchTimeout = []byte("ExternalBatchTimeout")

	// ParamsStoreKeyAverageExternalBlockTime stores the signed blocks window
	ParamsStoreKeyAverageExternalBlockTime = []byte("AverageExternalBlockTime")

	// ParamsStoreKeySignedWindow stores the signed blocks window
	ParamsStoreKeySignedWindow = []byte("SignedWindow")

	// ParamsStoreSlashFraction stores the slash fraction oracle set
	ParamsStoreSlashFraction = []byte("SlashFraction")

	// ParamStoreOracleSetUpdatePowerChangePercent oracle set update power change percent
	ParamStoreOracleSetUpdatePowerChangePercent = []byte("OracleSetUpdatePowerChangePercent")

	// ParamStoreIbcTransferTimeoutHeight gravity and ibc transfer timeout height
	ParamStoreIbcTransferTimeoutHeight = []byte("IbcTransferTimeoutHeight")

	// ParamOracleDelegateThreshold stores the oracle delegate threshold
	ParamOracleDelegateThreshold = []byte("OracleDelegateThreshold")

	// ParamOracleDelegateMultiple stores the oracle delegate multiple
	ParamOracleDelegateMultiple = []byte("OracleDelegateMultiple")
)

// Ensure that params implements the proper interface
var _ paramtypes.ParamSet = &Params{}

func DefaultParams() Params {
	return Params{
		GravityId:                         "fx-gravity-id",
		AverageBlockTime:                  7_000,
		AverageExternalBlockTime:          5_000,
		ExternalBatchTimeout:              12 * 3600 * 1000,
		SignedWindow:                      30_000,
		SlashFraction:                     sdk.NewDecWithPrec(8, 1), // 80%
		OracleSetUpdatePowerChangePercent: sdk.NewDecWithPrec(1, 1), // 10%
		IbcTransferTimeoutHeight:          20_000,
		DelegateThreshold:                 NewDelegateAmount(sdkmath.NewInt(10_000).MulRaw(1e18)),
		DelegateMultiple:                  DefaultOracleDelegateThreshold,
		Oracles:                           nil,
	}
}

// ValidateBasic checks that the parameters have valid values.
func (m *Params) ValidateBasic() error {
	if err := validateGravityID(m.GravityId); err != nil {
		return errorsmod.Wrap(err, "gravity id")
	}
	if err := validateAverageBlockTime(m.AverageBlockTime); err != nil {
		return errorsmod.Wrap(err, "average block time")
	}
	if err := validateExternalBatchTimeout(m.ExternalBatchTimeout); err != nil {
		return errorsmod.Wrap(err, "external batch timeout")
	}
	if err := validateAverageExternalBlockTime(m.AverageExternalBlockTime); err != nil {
		return errorsmod.Wrap(err, "average external block time")
	}
	if err := validateSignedWindow(m.SignedWindow); err != nil {
		return errorsmod.Wrap(err, "signed blocks window")
	}
	if err := validateSlashFraction(m.SlashFraction); err != nil {
		return errorsmod.Wrap(err, "slash fraction")
	}
	if err := validateIbcTransferTimeoutHeight(m.IbcTransferTimeoutHeight); err != nil {
		return errorsmod.Wrap(err, "ibc transfer timeout height")
	}
	if err := validateOracleSetUpdatePowerChangePercent(m.OracleSetUpdatePowerChangePercent); err != nil {
		return errorsmod.Wrap(err, "power changer update oracle set percent")
	}
	if err := validateOracleDelegateThreshold(m.DelegateThreshold); err != nil {
		return errorsmod.Wrap(err, "oracle delegate threshold")
	}
	if err := validateOracleDelegateMultiple(m.DelegateMultiple); err != nil {
		return errorsmod.Wrap(err, "delegate multiple")
	}
	return nil
}

// ParamKeyTable for auth module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
func (m *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamsStoreKeyGravityID, &m.GravityId, validateGravityID),
		paramtypes.NewParamSetPair(ParamsStoreKeyAverageBlockTime, &m.AverageBlockTime, validateAverageBlockTime),
		paramtypes.NewParamSetPair(ParamsStoreKeyExternalBatchTimeout, &m.ExternalBatchTimeout, validateExternalBatchTimeout),
		paramtypes.NewParamSetPair(ParamsStoreKeyAverageExternalBlockTime, &m.AverageExternalBlockTime, validateAverageExternalBlockTime),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedWindow, &m.SignedWindow, validateSignedWindow),
		paramtypes.NewParamSetPair(ParamsStoreSlashFraction, &m.SlashFraction, validateSlashFraction),
		paramtypes.NewParamSetPair(ParamStoreOracleSetUpdatePowerChangePercent, &m.OracleSetUpdatePowerChangePercent, validateOracleSetUpdatePowerChangePercent),
		paramtypes.NewParamSetPair(ParamStoreIbcTransferTimeoutHeight, &m.IbcTransferTimeoutHeight, validateIbcTransferTimeoutHeight),
		paramtypes.NewParamSetPair(ParamOracleDelegateThreshold, &m.DelegateThreshold, validateOracleDelegateThreshold),
		paramtypes.NewParamSetPair(ParamOracleDelegateMultiple, &m.DelegateMultiple, validateOracleDelegateMultiple),
	}
}

func validateGravityID(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if len(v) == 0 {
		return fmt.Errorf("gravityId cannpt be empty")
	}
	if _, err := fxtypes.StrToByte32(v); err != nil {
		return err
	}
	return nil
}

func validateExternalBatchTimeout(i interface{}) error {
	timeout, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if timeout < 60000 {
		return fmt.Errorf("invalid target batch timeout, less than 60 seconds is too short")
	}
	return nil
}

func validateAverageBlockTime(i interface{}) error {
	time, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if time < 100 {
		return fmt.Errorf("invalid average block time, too short for latency limitations")
	}
	return nil
}

func validateAverageExternalBlockTime(i interface{}) error {
	val, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if val < 100 {
		return fmt.Errorf("invalid average external block time, too short for latency limitations")
	}
	return nil
}

func validateSignedWindow(i interface{}) error {
	window, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if window <= 1 {
		return fmt.Errorf("invalid signed window too short")
	}
	return nil
}

func validateOracleDelegateThreshold(i interface{}) error {
	c, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !c.IsValid() || !c.IsPositive() {
		return fmt.Errorf("invalid delegate threshold")
	}
	if c.Denom != OracleDelegateDenom {
		return fmt.Errorf("oracle delegate denom must FX")
	}
	return nil
}

func validateIbcTransferTimeoutHeight(i interface{}) error {
	timeout, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if timeout <= 1 {
		return fmt.Errorf("invalid ibc transfer timeout too short")
	}
	return nil
}

func validateOracleDelegateMultiple(i interface{}) error {
	multiple, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if multiple <= 0 {
		return fmt.Errorf("invalid delegate multiple")
	}
	return nil
}

func validateOracleSetUpdatePowerChangePercent(i interface{}) error {
	percent, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if percent.IsNegative() {
		return fmt.Errorf("attempted to powet change percent with a negative: %v", percent)
	}
	if percent.GT(sdk.OneDec()) {
		return fmt.Errorf("powet change percent too large: %s", i)
	}
	return nil
}

func validateSlashFraction(i interface{}) error {
	slashFactor, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if slashFactor.IsNegative() {
		return fmt.Errorf("attempted to slash with a negative slash factor: %v", slashFactor)
	}
	if slashFactor.GT(sdk.OneDec()) {
		return fmt.Errorf("slash factor too large: %s", i)
	}
	return nil
}
