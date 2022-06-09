package types

import (
	"bytes"
	"fmt"

	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	MaxOracleSize                  = 150
	DefaultOracleDelegateThreshold = 10
)

var (
	// AttestationVotesPowerThreshold threshold of votes power to succeed
	AttestationVotesPowerThreshold = sdk.NewInt(66)
	// AttestationProposalOracleChangePowerThreshold
	AttestationProposalOracleChangePowerThreshold = sdk.NewInt(30)

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

	// ParamsStoreSlashFraction stores the slash fraction valset
	ParamsStoreSlashFraction = []byte("SlashFraction")

	// ParamStoreOracleSetUpdatePowerChangePercent valset update pwer change percent
	ParamStoreOracleSetUpdatePowerChangePercent = []byte("OracleSetUpdatePowerChangePercent")

	// ParamStoreIbcTransferTimeoutHeight gravity and ibc transfer timeout height
	ParamStoreIbcTransferTimeoutHeight = []byte("IbcTransferTimeoutHeight")

	// ParamStoreOracles stores the module oracles.
	ParamStoreOracles = []byte("Oracles")

	// ParamOracleDelegateThreshold stores the oracle delegate threshold
	ParamOracleDelegateThreshold = []byte("OracleDelegateThreshold")

	// ParamOracleDelegateMultiple stores the oracle delegate multiple
	ParamOracleDelegateMultiple = []byte("OracleDelegateMultiple")
)

var (
	// Ensure that params implements the proper interface
	_ paramtypes.ParamSet = &Params{}
)

// ValidateBasic checks that the parameters have valid values.
func (m Params) ValidateBasic() error {
	if err := validateGravityID(m.GravityId); err != nil {
		return sdkerrors.Wrap(err, "gravity id")
	}
	if err := validateAverageBlockTime(m.AverageBlockTime); err != nil {
		return sdkerrors.Wrap(err, "average block time")
	}
	if err := validateExternalBatchTimeout(m.ExternalBatchTimeout); err != nil {
		return sdkerrors.Wrap(err, "external batch timeout")
	}
	if err := validateAverageExternalBlockTime(m.AverageExternalBlockTime); err != nil {
		return sdkerrors.Wrap(err, "average external block time")
	}
	if err := validateSignedWindow(m.SignedWindow); err != nil {
		return sdkerrors.Wrap(err, "signed blocks window")
	}
	if err := validateSlashFraction(m.SlashFraction); err != nil {
		return sdkerrors.Wrap(err, "slash fraction")
	}
	if err := validateIbcTransferTimeoutHeight(m.IbcTransferTimeoutHeight); err != nil {
		return sdkerrors.Wrap(err, "ibc transfer timeout height")
	}
	if err := validateOracleSetUpdatePowerChangePercent(m.OracleSetUpdatePowerChangePercent); err != nil {
		return sdkerrors.Wrap(err, "power changer update oracle set percent")
	}
	if err := validateOracles(m.Oracles); err != nil {
		return sdkerrors.Wrap(err, "oracles")
	}
	if err := validateOracleDelegateThreshold(m.DelegateThreshold); err != nil {
		return sdkerrors.Wrap(err, "oracle delegate threshold")
	}
	if err := validateOracleDelegateMultiple(m.DelegateMultiple); err != nil {
		return sdkerrors.Wrap(err, "delegate multiple")
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
		paramtypes.NewParamSetPair(ParamStoreOracles, &m.Oracles, validateOracles),
		paramtypes.NewParamSetPair(ParamOracleDelegateThreshold, &m.DelegateThreshold, validateOracleDelegateThreshold),
		paramtypes.NewParamSetPair(ParamOracleDelegateMultiple, &m.DelegateMultiple, validateOracleDelegateMultiple),
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (m Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalLengthPrefixed(&m)
	bz2 := ModuleCdc.MustMarshalLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

func validateGravityID(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if len(v) == 0 {
		return fmt.Errorf("gravityId cannpt be empty")
	}
	if _, err := StrToFixByteArray(v); err != nil {
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
	} else if !c.IsValid() || !c.IsPositive() {
		return fmt.Errorf("invalid delegate threshold")
	}
	if c.Denom != fxtypes.DefaultDenom {
		return fmt.Errorf("")
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
	return nil
}

func validateOracles(i interface{}) error {
	if oracles, ok := i.([]string); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	} else {
		if len(oracles) <= 0 {
			return fmt.Errorf("oracles cannot be empty")
		}
		if len(oracles) > MaxOracleSize {
			return fmt.Errorf("oracle length must be less than or equal: %d", MaxOracleSize)
		}
		oraclesMap := make(map[string]bool)
		for _, addr := range oracles {
			if _, err := sdk.AccAddressFromBech32(addr); err != nil {
				return err
			}
			if oraclesMap[addr] {
				return fmt.Errorf("duplicate oracle: %s", addr)
			}
			oraclesMap[addr] = true
		}
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
	return nil
}

func StrToFixByteArray(s string) ([32]byte, error) {
	var out [32]byte
	if len([]byte(s)) > 32 {
		return out, fmt.Errorf("string too long")
	}
	copy(out[:], s)
	return out, nil
}
