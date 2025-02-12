package contract

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type BankPrecompileKeeper struct {
	Caller
	abi          abi.ABI
	contractAddr common.Address
}

func NewBankPrecompileKeeper(caller Caller, contractAddr common.Address) BankPrecompileKeeper {
	if IsZeroEthAddress(contractAddr) {
		contractAddr = common.HexToAddress(BankAddress)
	}
	return BankPrecompileKeeper{
		Caller:       caller,
		abi:          MustABIJson(IBankMetaData.ABI),
		contractAddr: contractAddr,
	}
}

func (k BankPrecompileKeeper) WithContract(addr common.Address) BankPrecompileKeeper {
	keeper := k
	keeper.contractAddr = addr
	return keeper
}

func (k BankPrecompileKeeper) TransferFromModuleToAccount(ctx context.Context, from common.Address, args TransferFromModuleToAccountArgs) (*evmtypes.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, from, k.contractAddr, nil, k.abi, "transferFromModuleToAccount", args.Module, args.Account, args.Token, args.Amount)
	if err != nil {
		return nil, err
	}
	return unpackRetIsOk(k.abi, "transferFromModuleToAccount", res)
}

func (k BankPrecompileKeeper) TransferFromAccountToModule(ctx context.Context, from common.Address, args TransferFromAccountToModuleArgs) (*evmtypes.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, from, k.contractAddr, nil, k.abi, "transferFromAccountToModule", args.Account, args.Module, args.Token, args.Amount)
	if err != nil {
		return nil, err
	}
	return unpackRetIsOk(k.abi, "transferFromAccountToModule", res)
}
