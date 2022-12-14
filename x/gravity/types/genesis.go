package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

// Deprecated: after upgrade v3
type GenesisState struct {
	Params                  Params                          `json:"params"`
	LastObservedNonce       uint64                          `json:"last_observed_nonce,omitempty"`
	LastObservedBlockHeight LastObservedEthereumBlockHeight `json:"last_observed_block_height"`
	DelegateKeys            []MsgSetOrchestratorAddress     `json:"delegate_keys"`
	Valsets                 []Valset                        `json:"valsets"`
	Erc20ToDenoms           []ERC20ToDenom                  `json:"erc20_to_denoms"`
	UnbatchedTransfers      []OutgoingTransferTx            `json:"unbatched_transfers"`
	Batches                 []OutgoingTxBatch               `json:"batches"`
	BatchConfirms           []MsgConfirmBatch               `json:"batch_confirms"`
	ValsetConfirms          []MsgValsetConfirm              `json:"valset_confirms"`
	Attestations            []Attestation                   `json:"attestations"`
	LastObservedValset      Valset                          `json:"last_observed_valset"`
	LastSlashedBatchBlock   uint64                          `json:"last_slashed_batch_block,omitempty"`
	LastSlashedValsetNonce  uint64                          `json:"last_slashed_valset_nonce,omitempty"`
	LastTxPoolId            uint64                          `json:"last_tx_pool_id"`
	LastBatchId             uint64                          `json:"last_batch_id"`
}

var (
	// Ensure that params implements the proper interface
	_ paramtypes.ParamSet = &Params{}
)

func (p *Params) ValidateBasic() error {
	if err := validateGravityID(p.GravityId); err != nil {
		return sdkerrors.Wrap(err, "gravity id")
	}
	if err := validateContractHash(p.ContractSourceHash); err != nil {
		return sdkerrors.Wrap(err, "contract hash")
	}
	if err := validateBridgeContractAddress(p.BridgeEthAddress); err != nil {
		return sdkerrors.Wrap(err, "bridge contract address")
	}
	if err := validateBridgeChainID(p.BridgeChainId); err != nil {
		return sdkerrors.Wrap(err, "bridge chain id")
	}
	if err := validateTargetBatchTimeout(p.TargetBatchTimeout); err != nil {
		return sdkerrors.Wrap(err, "Batch timeout")
	}
	if err := validateAverageBlockTime(p.AverageBlockTime); err != nil {
		return sdkerrors.Wrap(err, "Block time")
	}
	if err := validateAverageEthereumBlockTime(p.AverageEthBlockTime); err != nil {
		return sdkerrors.Wrap(err, "Ethereum block time")
	}
	if err := validateSignedValsetsWindow(p.SignedValsetsWindow); err != nil {
		return sdkerrors.Wrap(err, "signed blocks window")
	}
	if err := validateSignedBatchesWindow(p.SignedBatchesWindow); err != nil {
		return sdkerrors.Wrap(err, "signed blocks window")
	}
	if err := validateSignedClaimsWindow(p.SignedClaimsWindow); err != nil {
		return sdkerrors.Wrap(err, "signed blocks window")
	}
	if err := validateSlashFractionValset(p.SlashFractionValset); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	if err := validateSlashFractionBatch(p.SlashFractionBatch); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	if err := validateSlashFractionClaim(p.SlashFractionClaim); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	if err := validateSlashFractionConflictingClaim(p.SlashFractionConflictingClaim); err != nil {
		return sdkerrors.Wrap(err, "slash fraction valset")
	}
	if err := validateUnbondSlashingValsetsWindow(p.UnbondSlashingValsetsWindow); err != nil {
		return sdkerrors.Wrap(err, "unbond Slashing valset window")
	}
	if err := validateValsetUpdatePowerChangePercent(p.ValsetUpdatePowerChangePercent); err != nil {
		return sdkerrors.Wrap(err, "unbond Slashing valset window")
	}
	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	var (
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
	)
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamsStoreKeyGravityID, &p.GravityId, validateGravityID),
		paramtypes.NewParamSetPair(ParamsStoreKeyContractHash, &p.ContractSourceHash, validateContractHash),
		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractAddress, &p.BridgeEthAddress, validateBridgeContractAddress),
		paramtypes.NewParamSetPair(ParamsStoreKeyBridgeContractChainID, &p.BridgeChainId, validateBridgeChainID),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedValsetsWindow, &p.SignedValsetsWindow, validateSignedValsetsWindow),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedBatchesWindow, &p.SignedBatchesWindow, validateSignedBatchesWindow),
		paramtypes.NewParamSetPair(ParamsStoreKeySignedClaimsWindow, &p.SignedClaimsWindow, validateSignedClaimsWindow),
		paramtypes.NewParamSetPair(ParamsStoreKeyAverageBlockTime, &p.AverageBlockTime, validateAverageBlockTime),
		paramtypes.NewParamSetPair(ParamsStoreKeyTargetBatchTimeout, &p.TargetBatchTimeout, validateTargetBatchTimeout),
		paramtypes.NewParamSetPair(ParamsStoreKeyAverageEthereumBlockTime, &p.AverageEthBlockTime, validateAverageEthereumBlockTime),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionValset, &p.SlashFractionValset, validateSlashFractionValset),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionBatch, &p.SlashFractionBatch, validateSlashFractionBatch),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionClaim, &p.SlashFractionClaim, validateSlashFractionClaim),
		paramtypes.NewParamSetPair(ParamsStoreSlashFractionConflictingClaim, &p.SlashFractionConflictingClaim, validateSlashFractionConflictingClaim),
		paramtypes.NewParamSetPair(ParamStoreUnbondSlashingValsetsWindow, &p.UnbondSlashingValsetsWindow, validateUnbondSlashingValsetsWindow),
		paramtypes.NewParamSetPair(ParamStoreIbcTransferTimeoutHeight, &p.IbcTransferTimeoutHeight, validateIbcTransferTimeoutHeight),
		paramtypes.NewParamSetPair(ParamStoreValsetUpdatePowerChangePercent, &p.ValsetUpdatePowerChangePercent, validateValsetUpdatePowerChangePercent),
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
