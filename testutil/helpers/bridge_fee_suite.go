package helpers

import (
	"context"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/contract"
)

type BridgeFeeSuite struct {
	require *require.Assertions
	err     error

	contract.BridgeFeeQuoteKeeper
	contract.BridgeFeeOracleKeeper
}

func NewBridgeFeeSuite(require *require.Assertions, caller contract.Caller) BridgeFeeSuite {
	return BridgeFeeSuite{
		require:               require,
		BridgeFeeQuoteKeeper:  contract.NewBridgeFeeQuoteKeeper(caller),
		BridgeFeeOracleKeeper: contract.NewBridgeFeeOracleKeeper(caller),
	}
}

func (s BridgeFeeSuite) WithError(err error) BridgeFeeSuite {
	suite := s
	suite.err = err
	return suite
}

func (s BridgeFeeSuite) requireError(err error) {
	if s.err != nil {
		s.require.ErrorIs(err, evmtypes.ErrVMExecution.Wrap(s.err.Error()))
		return
	}
	s.require.NoError(err)
}

func (s BridgeFeeSuite) Quote(ctx context.Context, args contract.IBridgeFeeQuoteQuoteInput) *evmtypes.MsgEthereumTxResponse {
	defOracle, err := s.BridgeFeeOracleKeeper.DefaultOracle(ctx)
	s.require.NoError(err)

	res, err := s.BridgeFeeQuoteKeeper.Quote(ctx, defOracle, []contract.IBridgeFeeQuoteQuoteInput{args})
	s.requireError(err)
	return res
}
