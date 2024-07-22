package keeper

import (
	"context"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/types"

	fxevmtypes "github.com/functionx/fx-core/v7/x/evm/types"
)

var _ types.MsgServer = &Keeper{}

func (k *Keeper) CallContract(goCtx context.Context, msg *fxevmtypes.MsgCallContract) (*fxevmtypes.MsgCallContractResponse, error) {
	if !strings.EqualFold(k.GetAuthority().String(), msg.Authority) {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority, expected %s, got %s", k.GetAuthority().String(), msg.Authority)
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	contract := common.HexToAddress(msg.ContractAddress)
	account := k.GetAccount(ctx, contract)
	if account == nil || !account.IsContract() {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidProposalMsg, "contract %s not found", contract.Hex())
	}
	_, err := k.CallEVMWithoutGas(ctx, k.module, &contract, nil, common.Hex2Bytes(msg.Data), true)
	if err != nil {
		return nil, err
	}
	return &fxevmtypes.MsgCallContractResponse{}, nil
}
