package contract

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type EvmKeeper interface {
	DeployContract(ctx sdk.Context, from common.Address, abi abi.ABI, bin []byte, args ...interface{}) (common.Address, error)
	CreateContractWithCode(ctx sdk.Context, address common.Address, code []byte) error
	ApplyContract(ctx context.Context, from, contract common.Address, value *big.Int, abi abi.ABI, method string, constructorData ...interface{}) (*evmtypes.MsgEthereumTxResponse, error)
}

type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}
