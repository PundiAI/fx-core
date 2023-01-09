package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

// Ensure that params implements the proper interface
var _ paramtypes.ParamSet = &Params{}

func (m *Params) ValidateBasic() error {
	if err := validateGravityID(m.GravityId); err != nil {
		return sdkerrors.Wrap(err, "gravity id")
	}
	if err := validateContractHash(m.ContractSourceHash); err != nil {
		return sdkerrors.Wrap(err, "contract hash")
	}
	if err := validateBridgeContractAddress(m.BridgeEthAddress); err != nil {
		return sdkerrors.Wrap(err, "bridge contract address")
	}
	if err := validateBridgeChainID(m.BridgeChainId); err != nil {
		return sdkerrors.Wrap(err, "bridge chain id")
	}
	if err := validateTargetBatchTimeout(m.TargetBatchTimeout); err != nil {
		return sdkerrors.Wrap(err, "Batch timeout")
	}
	if err := validateAverageBlockTime(m.AverageBlockTime); err != nil {
		return sdkerrors.Wrap(err, "Block time")
	}
	if err := validateAverageEthereumBlockTime(m.AverageEthBlockTime); err != nil {
		return sdkerrors.Wrap(err, "ethereum block time")
	}
	if err := validateSignedValsetsWindow(m.SignedValsetsWindow); err != nil {
		return sdkerrors.Wrap(err, "signed blocks window")
	}
	if err := validateSignedBatchesWindow(m.SignedBatchesWindow); err != nil {
		return sdkerrors.Wrap(err, "signed blocks window")
	}
	if err := validateSignedClaimsWindow(m.SignedClaimsWindow); err != nil {
		return sdkerrors.Wrap(err, "signed blocks window")
	}
	if err := validateSlashFractionValset(m.SlashFractionValset); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	if err := validateSlashFractionBatch(m.SlashFractionBatch); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	if err := validateSlashFractionClaim(m.SlashFractionClaim); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	if err := validateSlashFractionConflictingClaim(m.SlashFractionConflictingClaim); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	if err := validateUnbondSlashingValsetsWindow(m.UnbondSlashingValsetsWindow); err != nil {
		return sdkerrors.Wrap(err, "unbond slashing valset window")
	}
	if err := validateValsetUpdatePowerChangePercent(m.ValsetUpdatePowerChangePercent); err != nil {
		return sdkerrors.Wrap(err, "unbond slashing valset window")
	}
	return nil
}

func (m *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	var (
		// ParamsStoreKeyGravityID stores the gravity id
		ParamsStoreKeyGravityID = []byte("GravityID")

		// ParamsStoreKeyContractHash stores the contract hash
		ParamsStoreKeyContractHash = []byte("ContractHash")

		// ParamsStoreKeyStartThreshold stores the start threshold
		// ParamsStoreKeyStartThreshold = []byte("StartThreshold")

		// ParamsStoreKeyBridgeContractAddress stores the contract address
		ParamsStoreKeyBridgeContractAddress = []byte("BridgeContractAddress")

		// ParamsStoreKeyBridgeContractChainID stores the bridge chain id
		ParamsStoreKeyBridgeContractChainID = []byte("BridgeChainID")

		// ParamsStoreKeySignedValsetsWindow stores the signed blocks window
		ParamsStoreKeySignedValsetsWindow = []byte("SignedValsetsWindow")

		// ParamsStoreKeySignedBatchesWindow stores the signed blocks window
		ParamsStoreKeySignedBatchesWindow = []byte("SignedBatchesWindow")

		// ParamsStoreKeySignedClaimsWindow stores the signed blocks window
		ParamsStoreKeySignedClaimsWindow = []byte("SignedClaimsWindow")

		// ParamsStoreKeyTargetBatchTimeout stores the signed blocks window
		ParamsStoreKeyTargetBatchTimeout = []byte("TargetBatchTimeout")

		// ParamsStoreKeyAverageBlockTime stores the signed blocks window
		ParamsStoreKeyAverageBlockTime = []byte("AverageBlockTime")

		// ParamsStoreKeyAverageEthereumBlockTime stores the signed blocks window
		ParamsStoreKeyAverageEthereumBlockTime = []byte("AverageEthereumBlockTime")

		// ParamsStoreSlashFractionValset stores the slash fraction valset
		ParamsStoreSlashFractionValset = []byte("SlashFractionValset")

		// ParamsStoreSlashFractionBatch stores the slash fraction Batch
		ParamsStoreSlashFractionBatch = []byte("SlashFractionBatch")

		// ParamsStoreSlashFractionClaim stores the slash fraction Claim
		ParamsStoreSlashFractionClaim = []byte("SlashFractionClaim")

		// ParamsStoreSlashFractionConflictingClaim stores the slash fraction ConflictingClaim
		ParamsStoreSlashFractionConflictingClaim = []byte("SlashFractionConflictingClaim")

		// ParamStoreUnbondSlashingValsetsWindow stores unbond slashing valset window
		ParamStoreUnbondSlashingValsetsWindow = []byte("UnbondSlashingValsetsWindow")

		// ParamStoreIbcTransferTimeoutHeight gravity and ibc transfer timeout height
		ParamStoreIbcTransferTimeoutHeight = []byte("ParamStoreIbcTransferTimeoutHeight")

		// ParamStoreValsetUpdatePowerChangePercent valset update power change percent
		ParamStoreValsetUpdatePowerChangePercent = []byte("ParamStoreValsetUpdatePowerChangePercent")
	)
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamsStoreKeyGravityID, &m.GravityId, validateGravityID),
		paramtypes.NewParamSetPair(ParamsStoreKeyContractHash, &m.ContractSourceHash, validateContractHash),
		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractAddress, &m.BridgeEthAddress, validateBridgeContractAddress),
		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractChainID, &m.BridgeChainId, validateBridgeChainID),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedValsetsWindow, &m.SignedValsetsWindow, validateSignedValsetsWindow),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedBatchesWindow, &m.SignedBatchesWindow, validateSignedBatchesWindow),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedClaimsWindow, &m.SignedClaimsWindow, validateSignedClaimsWindow),
		paramtypes.NewParamSetPair(ParamsStoreKeyAverageBlockTime, &m.AverageBlockTime, validateAverageBlockTime),
		paramtypes.NewParamSetPair(ParamsStoreKeyTargetBatchTimeout, &m.TargetBatchTimeout, validateTargetBatchTimeout),
		paramtypes.NewParamSetPair(ParamsStoreKeyAverageEthereumBlockTime, &m.AverageEthBlockTime, validateAverageEthereumBlockTime),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionValset, &m.SlashFractionValset, validateSlashFractionValset),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionBatch, &m.SlashFractionBatch, validateSlashFractionBatch),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionClaim, &m.SlashFractionClaim, validateSlashFractionClaim),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionConflictingClaim, &m.SlashFractionConflictingClaim, validateSlashFractionConflictingClaim),
		paramtypes.NewParamSetPair(ParamStoreUnbondSlashingValsetsWindow, &m.UnbondSlashingValsetsWindow, validateUnbondSlashingValsetsWindow),
		paramtypes.NewParamSetPair(ParamStoreIbcTransferTimeoutHeight, &m.IbcTransferTimeoutHeight, validateIbcTransferTimeoutHeight),
		paramtypes.NewParamSetPair(ParamStoreValsetUpdatePowerChangePercent, &m.ValsetUpdatePowerChangePercent, validateValsetUpdatePowerChangePercent),
	}
}

func validateGravityID(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if _, err := strToFixByteArray(v); err != nil {
		return err
	}
	return nil
}

func validateContractHash(i interface{}) error {
	if _, ok := i.(string); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if err := fxtypes.ValidateEthereumAddress(i.(string)); err != nil {
		if err.Error() != "empty" {
			return err
		}
	}
	return nil
}

func validateBridgeChainID(i interface{}) error {
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateTargetBatchTimeout(i interface{}) error {
	val, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	} else if val < 60000 {
		return fmt.Errorf("invalid target batch timeout, less than 60 seconds is too short")
	}
	return nil
}

func validateAverageBlockTime(i interface{}) error {
	val, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	} else if val < 100 {
		return fmt.Errorf("invalid average Cosmos block time, too short for latency limitations")
	}
	return nil
}

func validateAverageEthereumBlockTime(i interface{}) error {
	val, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	} else if val < 100 {
		return fmt.Errorf("invalid average Ethereum block time, too short for latency limitations")
	}
	return nil
}

func validateBridgeContractAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if err := fxtypes.ValidateEthereumAddress(v); err != nil {
		if err.Error() != "empty" {
			return err
		}
	}
	return nil
}

func validateSignedValsetsWindow(i interface{}) error {
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateUnbondSlashingValsetsWindow(i interface{}) error {
	if _, ok := i.(uint64); !ok {
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

func validateValsetUpdatePowerChangePercent(i interface{}) error {
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionValset(i interface{}) error {
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSignedBatchesWindow(i interface{}) error {
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSignedClaimsWindow(i interface{}) error {
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionBatch(i interface{}) error {
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionClaim(i interface{}) error {
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSlashFractionConflictingClaim(i interface{}) error {
	if _, ok := i.(sdk.Dec); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func strToFixByteArray(s string) ([32]byte, error) {
	var out [32]byte
	if len([]byte(s)) > 32 {
		return out, fmt.Errorf("string too long")
	}
	copy(out[:], s)
	return out, nil
}
