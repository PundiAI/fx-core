package app

import (
	"encoding/json"
	"math/big"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibchost "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/cosmos/ibc-go/v6/modules/core/exported"
	coretypes "github.com/cosmos/ibc-go/v6/modules/core/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	ethtypes "github.com/functionx/fx-core/v7/x/eth/types"
	fxstakingtypes "github.com/functionx/fx-core/v7/x/staking/types"
)

const (
	InitTotalSupply     = "378604525462891000000000000"
	EthModuleInitAmount = "378600525462891000000000000"
)

// GenesisState The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// NewDefAppGenesisByDenom return new genesis state
//
//gocyclo:ignore
func NewDefAppGenesisByDenom(denom string, cdc codec.JSONCodec) GenesisState {
	fxTotalSupply, ok := sdkmath.NewIntFromString(InitTotalSupply)
	if !ok {
		panic("invalid fx total supply")
	}
	ethInitAmount, ok := sdkmath.NewIntFromString(EthModuleInitAmount)
	if !ok {
		panic("invalid eth module init amount")
	}

	genesis := make(map[string]json.RawMessage)
	for _, b := range ModuleBasics {
		switch b.Name() {
		case stakingtypes.ModuleName:
			state := fxstakingtypes.DefaultGenesisState()
			state.Params.BondDenom = denom
			state.Params.MaxValidators = 20
			state.Params.UnbondingTime = time.Hour * 24 * 21
			state.Params.HistoricalEntries = 20000
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case slashingtypes.ModuleName:
			state := slashingtypes.DefaultGenesisState()
			state.Params.MinSignedPerWindow = sdk.NewDecWithPrec(5, 2) // 5%
			state.Params.SignedBlocksWindow = 20000
			state.Params.SlashFractionDoubleSign = sdk.NewDec(1).Quo(sdk.NewDec(20))
			state.Params.SlashFractionDowntime = sdk.NewDec(1).Quo(sdk.NewDec(1000))
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case distributiontypes.ModuleName:
			state := distributiontypes.DefaultGenesisState()
			state.Params.CommunityTax = sdk.NewDecWithPrec(40, 2) // %40
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case govtypes.ModuleName:
			state := govv1.DefaultGenesisState()
			coinOne := sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
			for i := 0; i < len(state.DepositParams.MinDeposit); i++ {
				state.DepositParams.MinDeposit[i].Denom = denom
				state.DepositParams.MinDeposit[i].Amount = coinOne.Mul(sdkmath.NewInt(10000))
			}
			duration := time.Hour * 24 * 14
			state.DepositParams.MaxDepositPeriod = &duration
			state.VotingParams.VotingPeriod = &duration
			state.TallyParams.Quorum = sdk.NewDecWithPrec(4, 1).String() // 40%
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case crisistypes.ModuleName:
			state := crisistypes.DefaultGenesisState()
			coinOne := sdkmath.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
			state.ConstantFee.Denom = denom
			state.ConstantFee.Amount = sdkmath.NewInt(13333).Mul(coinOne)
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case minttypes.ModuleName:
			state := minttypes.DefaultGenesisState()
			state.Params.MintDenom = denom
			state.Params.InflationMin = sdk.NewDecWithPrec(17, 2)        // 17%
			state.Params.InflationMax = sdk.NewDecWithPrec(416762, 6)    // 41.6762%
			state.Params.GoalBonded = sdk.NewDecWithPrec(51, 2)          // 51%
			state.Params.InflationRateChange = sdk.NewDecWithPrec(30, 2) // 30%
			state.Minter.Inflation = sdk.NewDecWithPrec(35, 2)           // 35%
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case banktypes.ModuleName:
			state := banktypes.DefaultGenesisState()
			state.DenomMetadata = []banktypes.Metadata{fxtypes.GetFXMetaData()}

			state.Supply = sdk.NewCoins(sdk.NewCoin(denom, fxTotalSupply))
			state.Balances = append(state.Balances, banktypes.Balance{
				Address: authtypes.NewModuleAddress(ethtypes.ModuleName).String(),
				Coins:   sdk.NewCoins(sdk.NewCoin(denom, ethInitAmount)),
			})
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case paramstypes.ModuleName:
			if state := b.DefaultGenesis(cdc); state == nil {
				genesis[b.Name()] = json.RawMessage("{}")
			} else {
				genesis[b.Name()] = state
			}
		case ibchost.ModuleName:
			state := coretypes.DefaultGenesisState()
			// only allowedClients tendermint
			state.ClientGenesis.Params.AllowedClients = []string{exported.Tendermint}
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case feemarkettypes.ModuleName:
			state := feemarkettypes.DefaultGenesisState()
			state.Params.BaseFee = sdkmath.NewInt(500_000_000_000)
			state.Params.MinGasPrice = sdk.NewDec(500_000_000_000)
			state.Params.MinGasMultiplier = sdk.ZeroDec()
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		default:
			genesis[b.Name()] = b.DefaultGenesis(cdc)
		}
	}
	return genesis
}

func CustomGenesisConsensusParams() *tmproto.ConsensusParams {
	result := tmtypes.DefaultConsensusParams()
	result.Block.MaxBytes = 1048576 // 1M
	result.Block.MaxGas = 30_000_000
	result.Block.TimeIotaMs = 1000
	result.Evidence.MaxAgeNumBlocks = 1000000
	result.Evidence.MaxBytes = 100000
	result.Evidence.MaxAgeDuration = 172800000000000
	return result
}
