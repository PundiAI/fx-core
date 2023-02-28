package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
)

type proposalServer interface {
	UpdateChainOraclesProposal(ctx sdk.Context, proposal *types.UpdateChainOraclesProposal) error
}

var _ proposalServer = Keeper{}

func (k Keeper) UpdateChainOraclesProposal(ctx sdk.Context, proposal *types.UpdateChainOraclesProposal) error {
	k.Logger(ctx).Info("handle update chain oracles proposal", "proposal", proposal.String())
	if len(proposal.Oracles) > types.MaxOracleSize {
		return errorsmod.Wrapf(types.ErrInvalid,
			fmt.Sprintf("oracle length must be less than or equal: %d", types.MaxOracleSize))
	}

	newOracleMap := make(map[string]bool, len(proposal.Oracles))
	for _, oracle := range proposal.Oracles {
		newOracleMap[oracle] = true
	}

	var unbondedOracleList []types.Oracle
	totalPower, deleteTotalPower := sdkmath.ZeroInt(), sdkmath.ZeroInt()

	allOracles := k.GetAllOracles(ctx, false)
	proposalOracle, _ := k.GetProposalOracle(ctx)
	oldOracleMap := make(map[string]bool, len(proposal.Oracles))
	for _, oracle := range proposalOracle.Oracles {
		oldOracleMap[oracle] = true
	}

	for _, oracle := range allOracles {
		if oracle.Online {
			totalPower = totalPower.Add(oracle.GetPower())
		}
		// oracle in new proposal
		if _, ok := newOracleMap[oracle.OracleAddress]; ok {
			continue
		}
		// oracle not in new proposal and oracle in old proposal
		if _, ok := oldOracleMap[oracle.OracleAddress]; ok {
			unbondedOracleList = append(unbondedOracleList, oracle)
			if oracle.Online {
				deleteTotalPower = deleteTotalPower.Add(oracle.GetPower())
			}
		}
	}

	maxChangePowerThreshold := types.AttestationProposalOracleChangePowerThreshold.Mul(totalPower).Quo(sdkmath.NewInt(100))
	k.Logger(ctx).Info("update chain oracles proposal",
		"maxChangePowerThreshold", maxChangePowerThreshold.String(), "deleteTotalPower", deleteTotalPower.String())
	if deleteTotalPower.GT(sdkmath.ZeroInt()) && deleteTotalPower.GTE(maxChangePowerThreshold) {
		return errorsmod.Wrapf(types.ErrInvalid, "max change power, "+
			"maxChangePowerThreshold: %s, deleteTotalPower: %s", maxChangePowerThreshold.String(), deleteTotalPower.String())
	}

	// update proposal oracle
	k.SetProposalOracle(ctx, &types.ProposalOracle{Oracles: proposal.Oracles})

	for _, unbondedOracle := range unbondedOracleList {
		if err := k.UnbondedOracleFromProposal(ctx, unbondedOracle); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) UnbondedOracleFromProposal(ctx sdk.Context, oracle types.Oracle) error {
	delegateAddr := oracle.GetDelegateAddress(k.moduleName)
	valAddr := oracle.GetValidator()
	getOracleDelegateToken, err := k.GetOracleDelegateToken(ctx, delegateAddr, valAddr)
	if err != nil {
		return err
	}
	msgUndelegate := stakingtypes.NewMsgUndelegate(delegateAddr, valAddr, sdk.NewCoin(fxtypes.DefaultDenom, getOracleDelegateToken))
	if _, err = k.stakingMsgServer.Undelegate(sdk.WrapSDKContext(ctx), msgUndelegate); err != nil {
		return err
	}

	oracle.Online = false
	k.SetOracle(ctx, oracle)

	return nil
}
