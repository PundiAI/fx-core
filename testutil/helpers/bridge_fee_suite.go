package helpers

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
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

func (s *BridgeFeeSuite) MockQuote(ctx context.Context, chainName, denom string) contract.IBridgeFeeQuoteQuoteInfo {
	input := contract.IBridgeFeeQuoteQuoteInput{
		Cap:       0,
		GasLimit:  contract.DefaultGasCap,
		Expiry:    uint64(time.Now().Add(time.Hour).Unix()),
		ChainName: contract.MustStrToByte32(chainName),
		TokenName: contract.MustStrToByte32(denom),
		Amount:    big.NewInt(1),
	}

	// add token if not exist
	tokens, err := s.GetTokens(ctx, input.ChainName)
	s.require.NoError(err)
	found := false
	for _, token := range tokens {
		if token == input.TokenName {
			found = true
		}
	}
	if !found {
		_, err = s.AddToken(ctx, input.ChainName, []common.Hash{input.TokenName})
		s.require.NoError(err)
	}

	s.Quote(ctx, input)

	quote, err := s.GetQuoteById(ctx, big.NewInt(1))
	s.require.NoError(err)
	s.require.Equal(input.ChainName, quote.ChainName)
	s.require.Equal(input.TokenName, quote.TokenName)
	s.require.Equal(input.Amount, quote.Amount)
	s.require.Equal(input.GasLimit, quote.GasLimit)
	s.require.Equal(input.Expiry, quote.Expiry)
	return quote
}
