package keeper_test

import (
	"encoding/json"
	"sync"
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/keeper"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

func TestParseEventLog(t *testing.T) {
	moduleAddr := common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes())
	log, complete := keeper.ParseEventLog(getReceipt(), moduleAddr)
	require.True(t, complete)
	require.Equal(t, 1, len(log.RelayToken))
	require.Equal(t, len(log.RelayToken), len(log.TransferCrossChain))
}

func BenchmarkSingleParseEventLog(b *testing.B) {
	receipt := getReceipt()

	b.Run("NoConcurrency", func(b *testing.B) {
		moduleAddr := common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes())
		_, complete := keeper.ParseEventLog(receipt, moduleAddr)
		require.True(b, complete)
	})

	b.Run("Concurrency", func(b *testing.B) {
		moduleAddr := common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes())
		_, complete := ParseEventLogConcurrencyTest(receipt, moduleAddr)
		require.True(b, complete)
	})
}

func BenchmarkMultipleParseEventLog(b *testing.B) {
	var receipts []*ethtypes.Receipt
	for i := 0; i < 20; i++ {
		receipts = append(receipts, getReceipt())
	}
	b.Run("NoConcurrency", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, receipt := range receipts {
				moduleAddr := common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes())
				_, complete := keeper.ParseEventLog(receipt, moduleAddr)
				require.True(b, complete)
			}

		}
	})

	b.Run("Concurrency", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, receipt := range receipts {
				moduleAddr := common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes())
				_, complete := ParseEventLogConcurrencyTest(receipt, moduleAddr)
				require.True(b, complete)
			}
		}
	})
}

func getReceipt() *ethtypes.Receipt {
	ethReceipt := `{"root":"0x","status":"0x1","cumulativeGasUsed":"0xf9a0","logsBloom":"0x00000000000000000000000000000000000000000001000000000000000001000000000100000000000000000000000000000000000000000000000000000000800000000000000000000008000000000000000000000040000000000000000000000000000000000000000000000000020000000000000000000010001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000800000002000000000000000000000000000000000000000000400000000000000000080000000000000000000000000000000000000000000000000000200000","logs":[{"address":"0x5fd55a1b9fc24967c4db09c513c3ba0dfa7ff687","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x00000000000000000000000049837e63c983fea83ab098e0107f96a714faa52c","0x00000000000000000000000047eeb2eac350e1923b8cbdfa4396a077b36e62a0"],"data":"0x0000000000000000000000000000000000000000000000001bc16d674ec80000","blockNumber":"0x2","transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x0","removed":false},{"address":"0x5fd55a1b9fc24967c4db09c513c3ba0dfa7ff687","topics":["0x282dd1817b996776123a00596764d4d54cc16460c9854f7a23f6be020ba0463d","0x00000000000000000000000049837e63c983fea83ab098e0107f96a714faa52c"],"data":"0x00000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000000de0b6b3a7640000636861696e2f6273630000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002a30783742374166373731374139354244364631636439303534653964303541623644623962643335386100000000000000000000000000000000000000000000","blockNumber":"0x2","transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000","transactionIndex":"0x0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","logIndex":"0x1","removed":false}],"transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000","contractAddress":"0x0000000000000000000000000000000000000000","gasUsed":"0xf9a0","blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","blockNumber":"0x2","transactionIndex":"0x0"}`
	r := &ethtypes.Receipt{}
	err := json.Unmarshal([]byte(ethReceipt), r)
	if err != nil {
		panic(err)
	}
	return r
}

func ParseEventLogConcurrencyTest(receipt *ethtypes.Receipt, moduleAddr common.Address) (keeper.EventLog, bool) {
	fip20ABI := fxtypes.GetERC20().ABI

	relayTokenEvents := make([]*keeper.RelayTokenEventLog, 0, len(receipt.Logs))
	transferCrossChainEvents := make([]*keeper.TransferCrossChainEventLog, 0, len(receipt.Logs))

	complete := true
	wg := sync.WaitGroup{}
	wg.Add(2)

	// parse relay token event
	go func() {
		defer wg.Done()
		for _, log := range receipt.Logs {
			rt, toAddr, err := keeper.ParseRelayTokenEvent(fip20ABI, log)
			if err != nil {
				complete = false
				break
			}
			if toAddr == moduleAddr {
				continue
			}
			relayTokenEvents = append(relayTokenEvents, &keeper.RelayTokenEventLog{Event: rt, Log: log})
		}
	}()

	// parse transfer cross chain event
	go func() {
		defer wg.Done()
		for _, log := range receipt.Logs {
			tc, err := fxtypes.ParseTransferCrossChainEvent(fip20ABI, log)
			if err != nil {
				complete = false
				break
			}
			if tc == nil {
				continue
			}
			transferCrossChainEvents = append(transferCrossChainEvents, &keeper.TransferCrossChainEventLog{Event: tc, Log: log})
		}
	}()

	wg.Wait()

	eventLog := keeper.EventLog{RelayToken: relayTokenEvents, TransferCrossChain: transferCrossChainEvents}
	return eventLog, complete
}
