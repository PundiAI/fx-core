package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/functionx/fx-core/x/crosschain/types"
)

func HandleUpdateCrossChainOraclesProposal(ctx sdk.Context, msgServer types.MsgServer, proposal *types.UpdateCrossChainOraclesProposal) error {
	ethereumMsgServer, ok := msgServer.(*EthereumMsgServer)
	if !ok {
		return sdkerrors.Wrap(types.ErrInvalid, "msg server")
	}
	keeper := ethereumMsgServer.Keeper

	logger := keeper.Logger(ctx)
	logger.Info("handle update crosschain oracles proposal", "proposal", proposal.String())
	if len(proposal.Oracles) > types.MaxOracleSize {
		return sdkerrors.Wrapf(types.ErrInvalid, fmt.Sprintf("oracle length must be less than or equal: %d", types.MaxOracleSize))
	}

	newOracleMap := make(map[string]bool, len(proposal.Oracles))
	for _, oracle := range proposal.Oracles {
		newOracleMap[oracle] = true
	}

	var deleteOracleList []types.Oracle
	var totalPower, deleteTotalPower = sdk.ZeroInt(), sdk.ZeroInt()

	allOracles := keeper.GetAllOracles(ctx, false)
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
	logger.Info("update crosschain oracles proposal", "maxChangePowerThreshold", maxChangePowerThreshold.String(), "deleteTotalPower", deleteTotalPower.String())
	if deleteTotalPower.GT(sdk.ZeroInt()) && deleteTotalPower.GTE(maxChangePowerThreshold) {
		return sdkerrors.Wrapf(types.ErrInvalid, "max change power, maxChangePowerThreshold: %s, deleteTotalPower: %s", maxChangePowerThreshold.String(), deleteTotalPower.String())
	}

	// update proposal oracle
	keeper.SetProposalOracle(ctx, &types.ProposalOracle{Oracles: proposal.Oracles})

	for _, deleteOracle := range deleteOracleList {
		keeper.DelExternalAddressForOracle(ctx, deleteOracle.ExternalAddress)
		bridgerAddr, err := sdk.AccAddressFromBech32(deleteOracle.BridgerAddress)
		if err != nil {
			panic(err)
		}
		keeper.DelOracleByBridger(ctx, bridgerAddr)
		oracleAddr := deleteOracle.GetOracle()
		keeper.DelOracle(ctx, oracleAddr)

		keeper.DelLastEventNonceByOracle(ctx, oracleAddr)
	}
	return nil
}
