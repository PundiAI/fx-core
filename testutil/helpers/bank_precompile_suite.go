package helpers

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/contract"
)

type BankPrecompileSuite struct {
	require *require.Assertions
	err     error

	keeper contract.BankPrecompileKeeper
}

func NewBankPrecompileSuite(require *require.Assertions, caller contract.Caller) BankPrecompileSuite {
	address := common.HexToAddress(contract.BankAddress)
	return BankPrecompileSuite{
		require: require,
		keeper:  contract.NewBankPrecompileKeeper(caller, address),
	}
}

func (s BankPrecompileSuite) WithContract(addr common.Address) BankPrecompileSuite {
	suite := s
	suite.keeper = suite.keeper.WithContract(addr)
	return suite
}

func (s BankPrecompileSuite) WithError(err error) BankPrecompileSuite {
	suite := s
	suite.err = err
	return suite
}

func (s BankPrecompileSuite) requireError(err error) {
	if s.err != nil {
		s.require.ErrorIs(err, evmtypes.ErrVMExecution.Wrap(s.err.Error()))
		return
	}
	s.require.NoError(err)
}

func (s BankPrecompileSuite) TransferFromModuleToAccount(ctx context.Context, from common.Address, args contract.TransferFromModuleToAccountArgs) *evmtypes.MsgEthereumTxResponse {
	transferModuleToAccount, err := s.keeper.TransferFromModuleToAccount(ctx, from, args)
	s.requireError(err)
	return transferModuleToAccount
}

func (s BankPrecompileSuite) TransferFromAccountToModule(ctx context.Context, from common.Address, args contract.TransferFromAccountToModuleArgs) *evmtypes.MsgEthereumTxResponse {
	transferAccountToModule, err := s.keeper.TransferFromAccountToModule(ctx, from, args)
	s.requireError(err)
	return transferAccountToModule
}
