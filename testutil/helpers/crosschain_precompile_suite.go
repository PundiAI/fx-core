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
	require *require.Assertions
	err     error

	keeper contract.CrosschainPrecompileKeeper
}

func NewCrosschainPrecompileSuite(require *require.Assertions, caller contract.Caller) CrosschainPrecompileSuite {
	address := common.HexToAddress(contract.CrosschainAddress)
	return CrosschainPrecompileSuite{
		require: require,
		keeper:  contract.NewCrosschainPrecompileKeeper(caller, address),
	}
}

func (s CrosschainPrecompileSuite) WithContract(addr common.Address) CrosschainPrecompileSuite {
	suite := s
	suite.keeper = suite.keeper.WithContract(addr)
	return suite
}

func (s CrosschainPrecompileSuite) WithError(err error) CrosschainPrecompileSuite {
	suite := s
	suite.err = err
	return suite
}

func (s CrosschainPrecompileSuite) requireError(err error) {
	if s.err != nil {
		s.require.ErrorIs(err, evmtypes.ErrVMExecution.Wrap(s.err.Error()))
		return
	}
	s.require.NoError(err)
}

func (s CrosschainPrecompileSuite) BridgeCoinAmount(ctx context.Context, args contract.BridgeCoinAmountArgs) *big.Int {
	amount, err := s.keeper.BridgeCoinAmount(ctx, args)
	s.requireError(err)
	return amount
}

func (s CrosschainPrecompileSuite) HasOracle(ctx context.Context, args contract.HasOracleArgs) bool {
	hasOracle, err := s.keeper.HasOracle(ctx, args)
	s.requireError(err)
	return hasOracle
}

func (s CrosschainPrecompileSuite) IsOracleOnline(ctx context.Context, args contract.IsOracleOnlineArgs) bool {
	isOracleOnline, err := s.keeper.IsOracleOnline(ctx, args)
	s.requireError(err)
	return isOracleOnline
}

func (s CrosschainPrecompileSuite) GetERC20Token(ctx context.Context, args contract.GetERC20TokenArgs) (common.Address, bool) {
	token, enable, err := s.keeper.GetERC20Token(ctx, args)
	s.requireError(err)
	return token, enable
}

func (s CrosschainPrecompileSuite) BridgeCall(ctx context.Context, value *big.Int, from common.Address, args contract.BridgeCallArgs) *evmtypes.MsgEthereumTxResponse {
	res, _, err := s.keeper.BridgeCall(ctx, value, from, args)
	s.requireError(err)
	return res
}

func (s CrosschainPrecompileSuite) ExecuteClaim(ctx context.Context, from common.Address, args contract.ExecuteClaimArgs) *evmtypes.MsgEthereumTxResponse {
	res, err := s.keeper.ExecuteClaim(ctx, from, args)
	s.requireError(err)
	return res
}

func (s CrosschainPrecompileSuite) Crosschain(ctx context.Context, value *big.Int, from common.Address, args contract.CrosschainArgs) *evmtypes.MsgEthereumTxResponse {
	res, err := s.keeper.Crosschain(ctx, value, from, args)
	s.requireError(err)
	return res
}
