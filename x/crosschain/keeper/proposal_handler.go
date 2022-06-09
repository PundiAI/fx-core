package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/crosschain/types"
)

var _ ProposalMsgServer = EthereumMsgServer{}

func (s EthereumMsgServer) HandleInitCrossChainParamsProposal(ctx sdk.Context, proposal *types.InitCrossChainParamsProposal) error {
	// check duplicate init params.
	var gravityId string
	s.paramSpace.GetIfExists(ctx, types.ParamsStoreKeyGravityID, &gravityId)
	if len(gravityId) != 0 {
		return sdkerrors.Wrapf(types.ErrInvalid, "duplicate init params chainName:%s", s.moduleName)
	}

	s.Logger(ctx).Info("handle init cross chain params...", "proposal", proposal.String())
	// init chain params
	s.SetParams(ctx, *proposal.Params)

	// FIP: slash fraction cannot greater than one 100%  2021-10-26.
	if proposal.Params.SlashFraction.GT(sdk.OneDec()) {
		return sdkerrors.Wrapf(types.ErrInvalid, "slash fraction too large: %s", proposal.Params.SlashFraction)
	}
	if proposal.Params.OracleSetUpdatePowerChangePercent.GT(sdk.OneDec()) {
		return sdkerrors.Wrapf(types.ErrInvalid, "oracle set update power change percent too large: %s", proposal.Params.OracleSetUpdatePowerChangePercent)
	}

	// save chain oracle
	s.SetChainOracles(ctx, &types.ChainOracle{Oracles: proposal.Params.Oracles})

	s.SetLastProposalBlockHeight(ctx, uint64(ctx.BlockHeight()))

	// init total stake
	s.SetTotalDelegate(ctx, sdk.NewCoin(proposal.Params.DelegateThreshold.Denom, sdk.ZeroInt()))
	return nil
}

func (s EthereumMsgServer) HandleUpdateChainOraclesProposal(ctx sdk.Context, proposal *types.UpdateChainOraclesProposal) error {
	logger := s.Logger(ctx)

	logger.Info("handle update cross chain update oracles proposal", "proposal", proposal.String())
	if len(proposal.Oracles) > types.MaxOracleSize {
		return sdkerrors.Wrapf(types.ErrInvalid, fmt.Sprintf("oracle length must be less than or equal: %d", types.MaxOracleSize))
	}
	// update chain oracle
	s.SetChainOracles(ctx, &types.ChainOracle{Oracles: proposal.Oracles})

	newOracleMap := make(map[string]bool, len(proposal.Oracles))
	for _, oracle := range proposal.Oracles {
		newOracleMap[oracle] = true
	}

	var deleteOracleList []types.Oracle
	var totalPower, deleteTotalPower = sdk.ZeroInt(), sdk.ZeroInt()

	allOracles := s.GetAllOracles(ctx)
	for _, oldOracle := range allOracles {

		if !oldOracle.Jailed {
			totalPower = totalPower.Add(oldOracle.GetPower())
		}
		if _, ok := newOracleMap[oldOracle.OracleAddress]; ok {
			continue
		}
		deleteOracleList = append(deleteOracleList, oldOracle)

		if !oldOracle.Jailed {
			deleteTotalPower = deleteTotalPower.Add(oldOracle.GetPower())
		}
	}

	maxChangePowerThreshold := types.AttestationProposalOracleChangePowerThreshold.Mul(totalPower).Quo(sdk.NewInt(100))
	logger.Info("UpdateChainOraclesProposal", "maxChangePowerThreshold", maxChangePowerThreshold.String(), "deleteTotalPower", deleteTotalPower.String())
	if deleteTotalPower.GT(sdk.ZeroInt()) && deleteTotalPower.GTE(maxChangePowerThreshold) {
		return sdkerrors.Wrapf(types.ErrInvalid, "max change power, maxChangePowerThreshold: %s, deleteTotalPower: %s", maxChangePowerThreshold.String(), deleteTotalPower.String())
	}

	for _, deleteOracle := range deleteOracleList {
		s.DelExternalAddressForOracle(ctx, deleteOracle.ExternalAddress)
		bridgerAddr, err := sdk.AccAddressFromBech32(deleteOracle.BridgerAddress)
		if err != nil {
			panic(err)
		}
		s.DelOracleByBridger(ctx, bridgerAddr)
		oracleAddr := deleteOracle.GetOracle()
		s.DelOracle(ctx, oracleAddr)

		s.DelLastEventNonceByOracle(ctx, oracleAddr)
	}
	return nil
}
