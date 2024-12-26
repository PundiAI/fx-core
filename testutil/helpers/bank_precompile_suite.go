package helpers

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/contract"
)

type BankPrecompileSuite struct {
	*ContractBaseSuite
	err error

	contract.BankPrecompileKeeper
}

func NewBankPrecompileSuite(require *require.Assertions, signer *Signer, caller contract.Caller, contractAddr common.Address) BankPrecompileSuite {
	contractBaseSuite := NewContractBaseSuite(require, signer)
	contractBaseSuite.WithContract(contractAddr)
	return BankPrecompileSuite{
		ContractBaseSuite:    contractBaseSuite,
		BankPrecompileKeeper: contract.NewBankPrecompileKeeper(caller, contractAddr),
	}
}

func (s BankPrecompileSuite) WithError(err error) BankPrecompileSuite {
	bankPrecompileKeeper := s
	bankPrecompileKeeper.err = err
	return bankPrecompileKeeper
}

func (s BankPrecompileSuite) requireError(err error) {
	if s.err != nil {
		s.require.ErrorIs(err, evmtypes.ErrVMExecution.Wrap(s.err.Error()))
		return
	}
	s.require.NoError(err)
}

func (s BankPrecompileSuite) TransferFromModuleToAccount(ctx context.Context, from common.Address, args contract.TransferFromModuleToAccountArgs) *evmtypes.MsgEthereumTxResponse {
	transferModuleToAccount, err := s.BankPrecompileKeeper.TransferFromModuleToAccount(ctx, from, args)
	s.requireError(err)
	return transferModuleToAccount
}

func (s BankPrecompileSuite) TransferFromAccountToModule(ctx context.Context, from common.Address, args contract.TransferFromAccountToModuleArgs) *evmtypes.MsgEthereumTxResponse {
	transferAccountToModule, err := s.BankPrecompileKeeper.TransferFromAccountToModule(ctx, from, args)
	s.requireError(err)
	return transferAccountToModule
}
