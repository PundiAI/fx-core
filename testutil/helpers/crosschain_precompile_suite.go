package helpers

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/contract"
)

type CrosschainPrecompileSuite struct {
	*ContractBaseSuite
	err error

	contract.CrosschainPrecompileKeeper
}

func NewCrosschainPrecompileSuite(require *require.Assertions, signer *Signer, caller contract.Caller, contractAddr common.Address) CrosschainPrecompileSuite {
	contractBaseSuite := NewContractBaseSuite(require, signer)
	contractBaseSuite.WithContract(contractAddr)
	return CrosschainPrecompileSuite{
		ContractBaseSuite:          contractBaseSuite,
		CrosschainPrecompileKeeper: contract.NewCrosschainPrecompileKeeper(caller, contractAddr),
	}
}

func (s CrosschainPrecompileSuite) WithError(err error) CrosschainPrecompileSuite {
	crosschainPrecompileKeeper := s
	crosschainPrecompileKeeper.err = err
	return crosschainPrecompileKeeper
}

func (s CrosschainPrecompileSuite) requireError(err error) {
	if s.err != nil {
		s.require.ErrorIs(err, evmtypes.ErrVMExecution.Wrap(s.err.Error()))
		return
	}
	s.require.NoError(err)
}

func (s CrosschainPrecompileSuite) BridgeCoinAmount(ctx context.Context, args contract.BridgeCoinAmountArgs) *big.Int {
	amount, err := s.CrosschainPrecompileKeeper.BridgeCoinAmount(ctx, args)
	s.requireError(err)
	return amount
}

func (s CrosschainPrecompileSuite) HasOracle(ctx context.Context, args contract.HasOracleArgs) bool {
	hasOracle, err := s.CrosschainPrecompileKeeper.HasOracle(ctx, args)
	s.requireError(err)
	return hasOracle
}

func (s CrosschainPrecompileSuite) IsOracleOnline(ctx context.Context, args contract.IsOracleOnlineArgs) bool {
	isOracleOnline, err := s.CrosschainPrecompileKeeper.IsOracleOnline(ctx, args)
	s.requireError(err)
	return isOracleOnline
}

func (s CrosschainPrecompileSuite) BridgeCall(ctx context.Context, from common.Address, args contract.BridgeCallArgs) *evmtypes.MsgEthereumTxResponse {
	res, _, err := s.CrosschainPrecompileKeeper.BridgeCall(ctx, from, args)
	s.requireError(err)
	return res
}

func (s CrosschainPrecompileSuite) ExecuteClaim(ctx context.Context, from common.Address, args contract.ExecuteClaimArgs) *evmtypes.MsgEthereumTxResponse {
	res, err := s.CrosschainPrecompileKeeper.ExecuteClaim(ctx, from, args)
	s.requireError(err)
	return res
}
