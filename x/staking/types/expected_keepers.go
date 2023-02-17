package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/types"
)

type AccountKeeper interface {
	stakingtypes.AccountKeeper
	GetModuleAddress(name string) sdk.AccAddress
	GetSequence(ctx sdk.Context, addr sdk.AccAddress) (uint64, error)
	GetModuleAddressAndPermissions(moduleName string) (addr sdk.AccAddress, permissions []string)
	NewAccount(ctx sdk.Context, acc authtypes.AccountI) authtypes.AccountI
}

type EvmKeeper interface {
	CallEVMWithData(ctx sdk.Context, from common.Address, contract *common.Address, data []byte, commit bool) (*types.MsgEthereumTxResponse, error)
	CreateContractWithCode(ctx sdk.Context, address common.Address, code []byte) error
	DeployUpgradableContract(ctx sdk.Context, from, logic common.Address, logicData []byte, initializeAbi *abi.ABI, initializeArgs ...interface{}) (common.Address, error)
}

type MockEvmKeeper struct{}

func (keeper *MockEvmKeeper) CreateContractWithCode(ctx sdk.Context, address common.Address, code []byte) error {
	return nil
}

var _ EvmKeeper = (*MockEvmKeeper)(nil)

func (keeper *MockEvmKeeper) CallEVMWithData(ctx sdk.Context, from common.Address, contract *common.Address, data []byte, commit bool) (*types.MsgEthereumTxResponse, error) {
	fmt.Printf("call evm with from: %x, to: %x, data: %x", from, contract, data)
	return &types.MsgEthereumTxResponse{
		Hash:    "",
		Logs:    nil,
		Ret:     nil,
		VmError: "",
		GasUsed: 0,
	}, nil
}

func (keeper *MockEvmKeeper) DeployUpgradableContract(ctx sdk.Context, from, logic common.Address, logicData []byte, initializeAbi *abi.ABI, initializeArgs ...interface{}) (common.Address, error) {
	return common.Address{}, nil
}
