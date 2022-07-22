package keeper

import (
	"fmt"

	types2 "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/crosschain/types"
)

var _ ProposalMsgServer = EthereumMsgServer{}

func (s EthereumMsgServer) HandleInitCrossChainParamsProposal(ctx sdk.Context, p *types.InitCrossChainParamsProposal) error {
	// check duplicate init params.
	var gravityId string
	s.paramSpace.GetIfExists(ctx, types.ParamsStoreKeyGravityID, &gravityId)
	if len(gravityId) != 0 {
		return sdkerrors.Wrapf(types.ErrInvalid, "duplicate init params chainName:%s", s.moduleName)
	}

	s.Logger(ctx).Info("handle init cross chain params...", "proposal", p.String())
	// init chain params
	s.SetParams(ctx, *p.Params)

	// FIP: slash fraction cannot greater than one 100%  2021-10-26.
	if ctx.BlockHeight() >= types2.CrossChainSupportTronBlock() {
		if p.Params.SlashFraction.GT(sdk.OneDec()) {
			return sdkerrors.Wrapf(types.ErrInvalid, "slash fraction too large: %s", p.Params.SlashFraction)
		}
		if p.Params.OracleSetUpdatePowerChangePercent.GT(sdk.OneDec()) {
			return sdkerrors.Wrapf(types.ErrInvalid, "oracle set update power change percent too large: %s", p.Params.OracleSetUpdatePowerChangePercent)
		}
	}

	// save chain oracle
	s.SetChainOracles(ctx, &types.ChainOracle{Oracles: p.Params.Oracles})

	s.SetLastProposalBlockHeight(ctx, uint64(ctx.BlockHeight()))

	// init total deposit
	s.SetTotalDeposit(ctx, sdk.NewCoin(p.Params.DepositThreshold.Denom, sdk.ZeroInt()))
	return nil
}

func (s EthereumMsgServer) HandleUpdateChainOraclesProposal(ctx sdk.Context, p *types.UpdateChainOraclesProposal) error {
	logger := s.Logger(ctx)

	logger.Info("handle update cross chain update oracles proposal", "proposal", p.String())
	if len(p.Oracles) > types.MaxOracleSize {
		return sdkerrors.Wrapf(types.ErrInvalid, fmt.Sprintf("oracle length must be less than or equal : %d", types.MaxOracleSize))
	}
	// update chain oracle
	s.SetChainOracles(ctx, &types.ChainOracle{Oracles: p.Oracles})

	newOracleMap := make(map[string]bool, len(p.Oracles))
	for _, oracle := range p.Oracles {
		newOracleMap[oracle] = true
	}

	var deleteOracleList []types.Oracle
	var totalDepositAmount, totalDeleteDepositAmount = sdk.ZeroInt(), sdk.ZeroInt()

	allOracles := s.GetAllOracles(ctx)
	for _, oldOracle := range allOracles {

		if !oldOracle.Jailed {
			totalDepositAmount = totalDepositAmount.Add(oldOracle.DepositAmount.Amount)
		}
		if _, ok := newOracleMap[oldOracle.OracleAddress]; ok {
			continue
		}
		deleteOracleList = append(deleteOracleList, oldOracle)

		if !oldOracle.Jailed {
			totalDeleteDepositAmount = totalDeleteDepositAmount.Add(oldOracle.DepositAmount.Amount)
		}
	}

	maxPowerChangeThreshold := types.AttestationProposalOracleChangePowerThreshold.Mul(totalDepositAmount).Quo(sdk.PowerReduction).Quo(sdk.NewInt(100))
	deleteOraclePower := totalDeleteDepositAmount.Quo(sdk.PowerReduction)
	//maxPowerChangeThreshold := sdk.NewDecFromInt(totalDepositAmount.Quo(sdk.PowerReduction)).Mul(sdk.NewDecFromInt(types.AttestationProposalOracleChangePowerThreshold)).Quo(sdk.NewDec(100))
	//deleteOraclePower := sdk.NewDecFromInt(totalDeleteDepositAmount.Quo(sdk.PowerReduction))
	logger.Info("UpdateChainOraclesProposal", "maxChangePower", maxPowerChangeThreshold.String(), "deleteOraclePower", deleteOraclePower.String())
	if deleteOraclePower.GT(sdk.ZeroInt()) && deleteOraclePower.GTE(maxPowerChangeThreshold) {
		return sdkerrors.Wrapf(types.ErrInvalid, "max change power!maxChangePower:%v,deletePower:%v",
			maxPowerChangeThreshold.String(), deleteOraclePower.String())
	}

	for _, deleteOracle := range deleteOracleList {
		s.DelExternalAddressForOracle(ctx, deleteOracle.ExternalAddress)
		orchestratorAddr, err := sdk.AccAddressFromBech32(deleteOracle.OrchestratorAddress)
		if err != nil {
			panic(err)
		}
		s.DelOracleByOrchestrator(ctx, orchestratorAddr)
		oracleAddr := deleteOracle.GetOracle()
		s.DelOracle(ctx, oracleAddr)

		s.DelLastEventNonceByOracle(ctx, oracleAddr)

		err = s.bankKeeper.SendCoinsFromModuleToAccount(ctx, s.moduleName, oracleAddr, sdk.NewCoins(deleteOracle.DepositAmount))
		if err != nil {
			panic(err)
		}
	}
	return nil
}
