package contract

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/types"
)

type AccessControlKeeper struct {
	Caller
	abi      abi.ABI
	from     common.Address
	contract common.Address
}

func NewAccessControlKeeper(caller Caller, contract string) AccessControlKeeper {
	return AccessControlKeeper{
		Caller:   caller,
		abi:      GetAccessControl().ABI,
		from:     common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes()),
		contract: common.HexToAddress(contract),
	}
}

func (k AccessControlKeeper) Initialize(ctx context.Context, admin common.Address) (*types.MsgEthereumTxResponse, error) {
	return k.Caller.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "initialize", admin)
}

func (k AccessControlKeeper) GrantRole(ctx context.Context, role common.Hash, account common.Address) (*types.MsgEthereumTxResponse, error) {
	return k.ApplyContract(ctx, k.from, k.contract, nil, k.abi, "grantRole", role, account)
}

func (k AccessControlKeeper) HasRole(ctx context.Context, role common.Hash, account common.Address) (bool, error) {
	var res struct{ has bool }
	if err := k.QueryContract(ctx, k.from, k.contract, k.abi, "hasRole", &res, role, account); err != nil {
		return false, err
	}
	return res.has, nil
}

func DeployAccessControlContract(
	ctx sdk.Context,
	evmKeeper EvmKeeper,
	evmModuleAddress,
	adminAddress common.Address,
) error {
	if err := DeployProxy(
		ctx,
		evmKeeper,
		GetAccessControl().ABI,
		GetAccessControl().Bin,
		common.HexToAddress(AccessControlAddress),
		evmModuleAddress,
	); err != nil {
		return err
	}
	accessControlKeeper := NewAccessControlKeeper(evmKeeper, AccessControlAddress)
	_, err := accessControlKeeper.Initialize(ctx, adminAddress)
	return err
}
