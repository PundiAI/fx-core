package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/types"

	fxevmtypes "github.com/functionx/fx-core/v8/x/evm/types"
)

var _ types.MsgServer = &Keeper{}

func (k *Keeper) CallContract(goCtx context.Context, msg *fxevmtypes.MsgCallContract) (*fxevmtypes.MsgCallContractResponse, error) {
	if !strings.EqualFold(k.GetAuthority().String(), msg.Authority) {
		return nil, govtypes.ErrInvalidSigner.Wrapf("invalid authority, expected %s, got %s", k.GetAuthority().String(), msg.Authority)
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	contract := common.HexToAddress(msg.ContractAddress)
	account := k.GetAccount(ctx, contract)
	if account == nil || !account.IsContract() {
		return nil, govtypes.ErrInvalidProposalMsg.Wrapf("contract %s not found", contract.Hex())
	}
	nonce, err := k.accountKeeper.GetSequence(ctx, k.module.Bytes())
	if err != nil {
		return nil, err
	}
	_, err = k.callEvm(ctx, k.module, &contract, nil, nonce, common.Hex2Bytes(msg.Data), true)
	if err != nil {
		return nil, err
	}
	return &fxevmtypes.MsgCallContractResponse{}, nil
}
