package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/types"
)

type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
}

type EvmKeeper interface {
	CallEVMWithData(ctx sdk.Context, from common.Address, contract *common.Address, data []byte, commit bool) (*types.MsgEthereumTxResponse, error)
}

type MockEvmKeeper struct{}

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
