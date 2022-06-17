package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/functionx/fx-core/x/crosschain/types"
)

func HandleUpdateChainOraclesProposal(ctx sdk.Context, msgServer types.MsgServer, proposal *types.UpdateChainOraclesProposal) error {
	ethereumMsgServer, ok := msgServer.(*EthereumMsgServer)
	if !ok {
		return sdkerrors.Wrap(types.ErrInvalid, "msg server")
	}
	keeper := ethereumMsgServer.Keeper

	logger := keeper.Logger(ctx)
	logger.Info("handle update chain oracles proposal", "proposal", proposal.String())
	if len(proposal.Oracles) > types.MaxOracleSize {
		return sdkerrors.Wrapf(types.ErrInvalid, fmt.Sprintf("oracle length must be less than or equal: %d", types.MaxOracleSize))
	}

	newOracleMap := make(map[string]bool, len(proposal.Oracles))
	for _, oracle := range proposal.Oracles {
		newOracleMap[oracle] = true
	}

	var unbondedOracleList []types.Oracle
	var totalPower, deleteTotalPower = sdk.ZeroInt(), sdk.ZeroInt()

	allOracles := keeper.GetAllOracles(ctx, false)
	proposalOracle, _ := keeper.GetProposalOracle(ctx)

	for _, oldOracle := range allOracles {

		if oldOracle.Online {
			totalPower = totalPower.Add(oldOracle.GetPower())
		}
		if _, ok := newOracleMap[oldOracle.OracleAddress]; ok {
			continue
		}
		for _, oracle := range proposalOracle.Oracles {
			if oracle == oldOracle.OracleAddress {
				unbondedOracleList = append(unbondedOracleList, oldOracle)
				if oldOracle.Online {
					deleteTotalPower = deleteTotalPower.Add(oldOracle.GetPower())
				}
			}
		}
	}

	maxChangePowerThreshold := types.AttestationProposalOracleChangePowerThreshold.Mul(totalPower).Quo(sdk.NewInt(100))
	logger.Info("update chain oracles proposal", "maxChangePowerThreshold", maxChangePowerThreshold.String(), "deleteTotalPower", deleteTotalPower.String())
	if deleteTotalPower.GT(sdk.ZeroInt()) && deleteTotalPower.GTE(maxChangePowerThreshold) {
		return sdkerrors.Wrapf(types.ErrInvalid, "max change power, maxChangePowerThreshold: %s, deleteTotalPower: %s", maxChangePowerThreshold.String(), deleteTotalPower.String())
	}

	// update proposal oracle
	keeper.SetProposalOracle(ctx, &types.ProposalOracle{Oracles: proposal.Oracles})

	var events = make(sdk.Events, 0)
	for _, unbondedOracle := range unbondedOracleList {
		delegateAddr := unbondedOracle.GetDelegateAddress(keeper.moduleName)
		valAddr := unbondedOracle.GetValidator()
		sharesAmount, err := keeper.stakingKeeper.ValidateUnbondAmount(ctx, delegateAddr, valAddr, unbondedOracle.DelegateAmount)
		if err != nil {
			return err
		}
		completionTime, err := keeper.stakingKeeper.Undelegate(ctx, delegateAddr, valAddr, sharesAmount)
		if err != nil {
			return err
		}
		unbondedOracle.Online = false
		keeper.SetOracle(ctx, unbondedOracle)
		events = append(events, sdk.NewEvent(
			stakingtypes.EventTypeUnbond,
			sdk.NewAttribute(stakingtypes.AttributeKeyValidator, unbondedOracle.DelegateValidator),
			sdk.NewAttribute(sdk.AttributeKeyAmount, unbondedOracle.DelegateAmount.String()),
			sdk.NewAttribute(stakingtypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		))
	}
	ctx.EventManager().EmitEvents(events)
	return nil
}
