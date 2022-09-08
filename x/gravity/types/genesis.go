package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	// AttestationVotesPowerThreshold threshold of votes power to succeed
	AttestationVotesPowerThreshold = sdk.NewInt(66)

	// ParamsStoreKeyGravityID stores the gravity id
	ParamsStoreKeyGravityID = []byte("GravityID")

	// ParamsStoreKeyContractHash stores the contract hash
	ParamsStoreKeyContractHash = []byte("ContractHash")

	// ParamsStoreKeyStartThreshold stores the start threshold
	//ParamsStoreKeyStartThreshold = []byte("StartThreshold")

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

	//ParamStoreValsetUpdatePowerChangePercent valset update pwer change percent
	ParamStoreValsetUpdatePowerChangePercent = []byte("ParamStoreValsetUpdatePowerChangePercent")

	// Ensure that params implements the proper interface
	_ paramtypes.ParamSet = &Params{}
)

// DefaultParams returns a copy of the default params
func DefaultParams() Params {
	return Params{
		GravityId:                      "fx-bridge-eth",
		BridgeChainId:                  1,
		SignedValsetsWindow:            10000,
		SignedBatchesWindow:            10000,
		SignedClaimsWindow:             10000,
		TargetBatchTimeout:             43200000,
		AverageBlockTime:               5000,
		AverageEthBlockTime:            15000,
		SlashFractionValset:            sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashFractionBatch:             sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashFractionClaim:             sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		SlashFractionConflictingClaim:  sdk.NewDec(1).Quo(sdk.NewDec(1000)),
		UnbondSlashingValsetsWindow:    10000,
		IbcTransferTimeoutHeight:       10000,
		ValsetUpdatePowerChangePercent: sdk.NewDec(1).Quo(sdk.NewDec(10)),
	}
}

func (p Params) ValidateBasic() error {
	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamsStoreKeyGravityID, &p.GravityId, validate),
		paramtypes.NewParamSetPair(ParamsStoreKeyContractHash, &p.ContractSourceHash, validate),
		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractAddress, &p.BridgeEthAddress, validate),
		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractChainID, &p.BridgeChainId, validate),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedValsetsWindow, &p.SignedValsetsWindow, validate),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedBatchesWindow, &p.SignedBatchesWindow, validate),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedClaimsWindow, &p.SignedClaimsWindow, validate),
		paramtypes.NewParamSetPair(ParamsStoreKeyAverageBlockTime, &p.AverageBlockTime, validate),
		paramtypes.NewParamSetPair(ParamsStoreKeyTargetBatchTimeout, &p.TargetBatchTimeout, validate),
		paramtypes.NewParamSetPair(ParamsStoreKeyAverageEthereumBlockTime, &p.AverageEthBlockTime, validate),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionValset, &p.SlashFractionValset, validate),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionBatch, &p.SlashFractionBatch, validate),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionClaim, &p.SlashFractionClaim, validate),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionConflictingClaim, &p.SlashFractionConflictingClaim, validate),
		paramtypes.NewParamSetPair(ParamStoreUnbondSlashingValsetsWindow, &p.UnbondSlashingValsetsWindow, validate),
		paramtypes.NewParamSetPair(ParamStoreIbcTransferTimeoutHeight, &p.IbcTransferTimeoutHeight, validate),
		paramtypes.NewParamSetPair(ParamStoreValsetUpdatePowerChangePercent, &p.ValsetUpdatePowerChangePercent, validate),
	}
}

func validate(_ interface{}) error {
	return nil
}
