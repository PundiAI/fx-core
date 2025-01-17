package app

import (
	"encoding/json"
	"time"

	sdkmath "cosmossdk.io/math"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	coretypes "github.com/cosmos/ibc-go/v8/modules/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	fxtypes "github.com/pundiai/fx-core/v8/types"
)

// GenesisState The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// newDefAppGenesisByDenom return new genesis state
//
//nolint:gocyclo // a lot of modules need to be modified
func newDefAppGenesisByDenom(cdc codec.JSONCodec, moduleBasics module.BasicManager) GenesisState {
	denom := fxtypes.DefaultDenom
	genesis := make(map[string]json.RawMessage)
	for _, m := range moduleBasics {
		switch m.Name() {
		case stakingtypes.ModuleName:
			state := stakingtypes.DefaultGenesisState()
			state.Params.BondDenom = denom
			state.Params.MaxValidators = 20
			state.Params.UnbondingTime = time.Hour * 24 * 21
			state.Params.HistoricalEntries = 20000
			genesis[m.Name()] = cdc.MustMarshalJSON(state)
		case slashingtypes.ModuleName:
			state := slashingtypes.DefaultGenesisState()
			state.Params.MinSignedPerWindow = sdkmath.LegacyNewDecWithPrec(5, 2) // 5%
			state.Params.SignedBlocksWindow = 20000
			state.Params.SlashFractionDoubleSign = sdkmath.LegacyNewDec(1).Quo(sdkmath.LegacyNewDec(20))
			state.Params.SlashFractionDowntime = sdkmath.LegacyNewDec(1).Quo(sdkmath.LegacyNewDec(1000))
			genesis[m.Name()] = cdc.MustMarshalJSON(state)
		case distributiontypes.ModuleName:
			state := distributiontypes.DefaultGenesisState()
			state.Params.CommunityTax = sdkmath.LegacyNewDecWithPrec(40, 2) // %40
			state.Params.BaseProposerReward = sdkmath.LegacyNewDecWithPrec(1, 2)
			state.Params.BonusProposerReward = sdkmath.LegacyNewDecWithPrec(4, 2)
			genesis[m.Name()] = cdc.MustMarshalJSON(state)
		case govtypes.ModuleName:
			state := govv1.DefaultGenesisState()
			minDepositAmount := sdkmath.NewInt(1e18).MulRaw(30)
			state.Params.MinDeposit = sdk.NewCoins(sdk.NewCoin(denom, minDepositAmount))
			state.Params.ExpeditedMinDeposit = sdk.NewCoins(sdk.NewCoin(denom, minDepositAmount.MulRaw(govv1.DefaultMinExpeditedDepositTokensRatio)))
			state.Params.MinInitialDepositRatio = sdkmath.LegacyMustNewDecFromStr("0.33").String()
			state.Params.MinDepositRatio = sdkmath.LegacyMustNewDecFromStr("0").String()

			duration := time.Hour * 24 * 14
			state.Params.MaxDepositPeriod = &duration
			state.Params.VotingPeriod = &duration
			state.Params.Quorum = sdkmath.LegacyNewDecWithPrec(4, 1).String() // 40%
			genesis[m.Name()] = cdc.MustMarshalJSON(state)
		case crisistypes.ModuleName:
			state := crisistypes.DefaultGenesisState()
			state.ConstantFee.Denom = denom
			state.ConstantFee.Amount = sdkmath.NewInt(133).MulRaw(1e18)
			genesis[m.Name()] = cdc.MustMarshalJSON(state)
		case minttypes.ModuleName:
			state := minttypes.DefaultGenesisState()
			state.Params.MintDenom = denom
			state.Params.InflationMin = sdkmath.LegacyNewDecWithPrec(17, 2)        // 17%
			state.Params.InflationMax = sdkmath.LegacyNewDecWithPrec(416762, 6)    // 41.6762%
			state.Params.GoalBonded = sdkmath.LegacyNewDecWithPrec(51, 2)          // 51%
			state.Params.InflationRateChange = sdkmath.LegacyNewDecWithPrec(30, 2) // 30%
			state.Minter.Inflation = sdkmath.LegacyNewDecWithPrec(35, 2)           // 35%
			genesis[m.Name()] = cdc.MustMarshalJSON(state)
		case banktypes.ModuleName:
			state := banktypes.DefaultGenesisState()
			state.DenomMetadata = []banktypes.Metadata{fxtypes.NewDefaultMetadata()}
			genesis[m.Name()] = cdc.MustMarshalJSON(state)
		case paramstypes.ModuleName:
			if mod, ok := m.(module.HasGenesisBasics); ok {
				if state := mod.DefaultGenesis(cdc); state == nil {
					genesis[m.Name()] = json.RawMessage("{}")
				} else {
					genesis[m.Name()] = state
				}
			}
		case ibcexported.ModuleName:
			state := coretypes.DefaultGenesisState()
			// only allowedClients tendermint
			state.ClientGenesis.Params.AllowedClients = []string{ibcexported.Tendermint}
			genesis[m.Name()] = cdc.MustMarshalJSON(state)
		case evmtypes.ModuleName:
			state := evmtypes.DefaultGenesisState()
			state.Params.EvmDenom = denom
			genesis[m.Name()] = cdc.MustMarshalJSON(state)
		case feemarkettypes.ModuleName:
			state := feemarkettypes.DefaultGenesisState()
			state.Params.BaseFee = sdkmath.NewInt(fxtypes.DefaultGasPrice)
			state.Params.MinGasPrice = sdkmath.LegacyNewDec(fxtypes.DefaultGasPrice)
			state.Params.MinGasMultiplier = sdkmath.LegacyZeroDec()
			genesis[m.Name()] = cdc.MustMarshalJSON(state)
		default:
			if mod, ok := m.(module.HasGenesisBasics); ok {
				genesis[m.Name()] = mod.DefaultGenesis(cdc)
			}
		}
	}
	return genesis
}

func CustomGenesisConsensusParams() *tmtypes.ConsensusParams {
	result := tmtypes.DefaultConsensusParams()
	result.Block.MaxBytes = 1048576 // 1M
	result.Block.MaxGas = 30_000_000
	result.Evidence.MaxAgeNumBlocks = 1000000
	result.Evidence.MaxBytes = 100000
	result.Evidence.MaxAgeDuration = 172800000000000
	return result
}
