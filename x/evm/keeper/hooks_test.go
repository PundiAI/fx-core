package keeper_test

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"

	"github.com/functionx/fx-core/x/evm/types"
)

// LogRecordHook records all the logs
type LogRecordHook struct {
	Logs []*ethtypes.Log
}

func (dh *LogRecordHook) PostTxProcessing(ctx sdk.Context, tx *ethtypes.Transaction, logs []*ethtypes.Log) error {
	dh.Logs = logs
	return nil
}

// FailureHook always fail
type FailureHook struct{}

func (dh FailureHook) PostTxProcessing(ctx sdk.Context, tx *ethtypes.Transaction, logs []*ethtypes.Log) error {
	return errors.New("post tx processing failed")
}

func (suite *KeeperTestSuite) TestEvmHooks() {
	testCases := []struct {
		msg       string
		setupHook func() types.EvmHooks
		expFunc   func(hook types.EvmHooks, result error)
	}{
		{
			"log collect hook",
			func() types.EvmHooks {
				return &LogRecordHook{}
			},
			func(hook types.EvmHooks, result error) {
				suite.Require().NoError(result)
			},
		},
		{
			"always fail hook",
			func() types.EvmHooks {
				return &FailureHook{}
			},
			func(hook types.EvmHooks, result error) {
				suite.Require().NoError(result)
			},
		},
	}

	for _, tc := range testCases {
		suite.SetupTest()
		hook := tc.setupHook()

		k := suite.app.EvmKeeper

		tx := ethtypes.NewTx(&ethtypes.DynamicFeeTx{
			Nonce:      1,
			To:         &common.Address{},
			Value:      big.NewInt(0),
			Gas:        10000,
			AccessList: make(ethtypes.AccessList, 0),
			ChainID:    new(big.Int),
			GasTipCap:  new(big.Int),
			GasFeeCap:  new(big.Int),
			V:          new(big.Int),
			R:          new(big.Int),
			S:          new(big.Int),
		})

		//txHash := common.BigToHash(big.NewInt(1))
		k.SetTxHashTransient(tx.Hash())
		k.AddLog(&ethtypes.Log{
			Topics:  []common.Hash{},
			Address: suite.address,
		})
		logs := k.GetTxLogsTransient(tx.Hash())
		result := k.PostTxProcessing(tx, logs)

		tc.expFunc(hook, result)
	}
}
