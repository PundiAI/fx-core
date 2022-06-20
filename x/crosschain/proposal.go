package crosschain

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/functionx/fx-core/x/crosschain/keeper"
	"github.com/functionx/fx-core/x/crosschain/types"
	tronkeeper "github.com/functionx/fx-core/x/tron/keeper"
)

func NewChainProposalHandler(k keeper.RouterKeeper) govtypes.Handler {
	moduleHandlerRouter := k.Router()
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.UpdateChainOraclesProposal:
			if !moduleHandlerRouter.HasRoute(c.ChainName) {
				return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized cross chain type: %s", c.ChainName))
			}
			return HandleUpdateChainOraclesProposal(ctx, k.Router().GetRoute(c.ChainName).MsgServer, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Unrecognized %s proposal content type: %T", types.ModuleName, c)
		}
	}
}

func HandleUpdateChainOraclesProposal(ctx sdk.Context, msgServer types.MsgServer, proposal *types.UpdateChainOraclesProposal) error {
	var k keeper.Keeper
	switch server := msgServer.(type) {
	case *keeper.EthereumMsgServer:
		k = server.Keeper
	case *tronkeeper.TronMsgServer:
		k = server.Keeper
	default:
		return sdkerrors.Wrapf(types.ErrInvalid, "msg server: %T", msgServer)
	}

	logger := k.Logger(ctx)
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

	allOracles := k.GetAllOracles(ctx, false)
	proposalOracle, _ := k.GetProposalOracle(ctx)

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
	k.SetProposalOracle(ctx, &types.ProposalOracle{Oracles: proposal.Oracles})

	for _, unbondedOracle := range unbondedOracleList {
		if err := k.UnbondedOracle(ctx, unbondedOracle); err != nil {
			return err
		}
	}
	return nil
}
