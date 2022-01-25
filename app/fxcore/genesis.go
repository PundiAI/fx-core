package fxcore

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibchost "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	gravitytypes "github.com/functionx/fx-core/x/gravity/types"
)

const (
	BankModuleTotalSupply   = "378604525462891000000000000"
	GravityModuleInitAmount = "378600525462891000000000000"
)

// AppGenesisState The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type AppGenesisState map[string]json.RawMessage

func NewDefAppGenesisByDenom(denom string, cdc codec.JSONCodec) map[string]json.RawMessage {
	genesis := make(map[string]json.RawMessage)
	for _, b := range ModuleBasics {
		switch b.Name() {
		case stakingtypes.ModuleName:
			state := stakingtypes.DefaultGenesisState()
			state.Params.BondDenom = denom
			state.Params.MaxValidators = 20
			state.Params.UnbondingTime = time.Hour * 24 * 21
			state.Params.HistoricalEntries = 20000
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case slashingtypes.ModuleName:
			state := slashingtypes.DefaultGenesisState()
			state.Params.MinSignedPerWindow = sdk.NewDecWithPrec(5, 2)
			state.Params.SignedBlocksWindow = 20000
			state.Params.SlashFractionDoubleSign = sdk.NewDec(1).Quo(sdk.NewDec(20))
			state.Params.SlashFractionDowntime = sdk.NewDec(1).Quo(sdk.NewDec(1000))
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case distributiontypes.ModuleName:
			state := distributiontypes.DefaultGenesisState()
			state.Params.CommunityTax = sdk.NewDecWithPrec(40, 2) // %40
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case govtypes.ModuleName:
			state := govtypes.DefaultGenesisState()
			coinOne := sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
			for i := 0; i < state.DepositParams.MinDeposit.Len(); i++ {
				state.DepositParams.MinDeposit[i].Denom = denom
				state.DepositParams.MinDeposit[i].Amount = coinOne.Mul(sdk.NewInt(10000))
			}
			state.DepositParams.MaxDepositPeriod = time.Hour * 24 * 14
			state.VotingParams.VotingPeriod = time.Hour * 24 * 14
			state.TallyParams.Quorum = sdk.NewDecWithPrec(4, 1)
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case crisistypes.ModuleName:
			state := crisistypes.DefaultGenesisState()
			coinOne := sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
			state.ConstantFee.Denom = denom
			state.ConstantFee.Amount = sdk.NewInt(13333).Mul(coinOne)
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case minttypes.ModuleName:
			state := minttypes.DefaultGenesisState()
			state.Params.MintDenom = denom
			state.Params.InflationMin = sdk.NewDecWithPrec(17, 2)
			state.Params.InflationMax = sdk.NewDecWithPrec(416762, 6)
			state.Params.GoalBonded = sdk.NewDecWithPrec(51, 2)
			state.Params.InflationRateChange = sdk.NewDecWithPrec(30, 2)
			state.Minter.Inflation = sdk.NewDecWithPrec(35, 2)
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case banktypes.ModuleName:
			state := banktypes.DefaultGenesisState()
			state.DenomMetadata = []banktypes.Metadata{GetFxBankMetaData(denom)}
			fxTotalSupply, ok := sdk.NewIntFromString(BankModuleTotalSupply)
			if !ok {
				panic("invalid fx total supply")
			}
			state.Supply = sdk.NewCoins(sdk.NewCoin(denom, fxTotalSupply))
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case gravitytypes.ModuleName:
			state := gravitytypes.DefaultGenesisState()
			state.Params.SignedValsetsWindow = 20000
			state.Params.SignedBatchesWindow = 20000
			state.Params.SignedClaimsWindow = 20000
			state.Params.UnbondSlashingValsetsWindow = 20000
			state.Params.IbcTransferTimeoutHeight = 20000

			initAmount, ok := sdk.NewIntFromString(GravityModuleInitAmount)
			if !ok {
				panic(fmt.Errorf("gravity module init amount err!!!amount:[%v]", GravityModuleInitAmount))
			}
			state.ModuleCoins = sdk.NewCoins(sdk.NewCoin(denom, initAmount))
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		case ibchost.ModuleName:
			state := types.DefaultGenesisState()
			// only allowedClients tendermint
			state.ClientGenesis.Params.AllowedClients = []string{exported.Tendermint}
			genesis[b.Name()] = cdc.MustMarshalJSON(state)
		default:
			genesis[b.Name()] = b.DefaultGenesis(cdc)
		}
	}
	return genesis
}

func GetFxBankMetaData(denom string) banktypes.Metadata {
	return banktypes.Metadata{
		Description: "Function X",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    strings.ToLower(denom),
				Exponent: 0,
				Aliases:  nil,
			},
			{
				Denom:    denom,
				Exponent: 18,
				Aliases:  nil,
			},
		},
		Base:    strings.ToLower(denom),
		Display: denom,
	}
}
