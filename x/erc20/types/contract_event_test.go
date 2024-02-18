package types_test

import (
	"encoding/json"
	"sync"
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v7/x/erc20/types"
)

func BenchmarkSingleParseEventLog(b *testing.B) {
	logs := getEventLogs()

	b.Run("NoConcurrency", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			moduleAddr := common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes())
			complete := parseEventLogTest(logs, moduleAddr)
			require.True(b, complete)
		}
	})

	b.Run("Concurrency", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			moduleAddr := common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes())
			complete := parseEventLogConcurrencyTest(logs, moduleAddr)
			require.True(b, complete)
		}
	})
}

func BenchmarkMultipleParseEventLog(b *testing.B) {
	var logsAry [][]*ethtypes.Log
	for i := 0; i < 20; i++ {
		logsAry = append(logsAry, getEventLogs())
	}
	b.Run("NoConcurrency", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, logs := range logsAry {
				moduleAddr := common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes())
				complete := parseEventLogTest(logs, moduleAddr)
				require.True(b, complete)
			}
		}
	})

	b.Run("Concurrency", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, logs := range logsAry {
				moduleAddr := common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes())
				complete := parseEventLogConcurrencyTest(logs, moduleAddr)
				require.True(b, complete)
			}
		}
	})
}

func parseEventLogConcurrencyTest(logs []*ethtypes.Log, moduleAddr common.Address) bool {
	complete := true
	wg := sync.WaitGroup{}
	wg.Add(2)

	// parse relay token event
	go func() {
		defer wg.Done()
		for _, log := range logs {
			rt, err := types.ParseTransferEvent(log)
			if err != nil {
				complete = false
				break
			}

			if rt != nil && rt.To == moduleAddr {
				continue
			}
			_ = rt
		}
	}()

	// parse transfer cross chain event
	go func() {
		defer wg.Done()
		for _, log := range logs {
			tc, err := types.ParseTransferCrossChainEvent(log)
			if err != nil {
				complete = false
				break
			}
			if tc == nil {
				continue
			}
		}
	}()

	wg.Wait()

	return complete
}

func parseEventLogTest(logs []*ethtypes.Log, moduleAddress common.Address) bool {
	for _, log := range logs {
		rt, err := types.ParseTransferEvent(log)
		if err != nil {
			return false
		}
		tc, err := types.ParseTransferCrossChainEvent(log)
		if err != nil {
			return false
		}
		if rt != nil && rt.To != moduleAddress && tc == nil {
			continue
		}

		if rt != nil && rt.To == moduleAddress {
			_ = rt
		}
		if tc != nil {
			_ = tc
		}
	}
	return true
}

func getEventLogs() []*ethtypes.Log {
	data := `[
        {
            "address":"0x5fd55a1b9fc24967c4db09c513c3ba0dfa7ff687",
            "topics":[
                "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
                "0x00000000000000000000000049837e63c983fea83ab098e0107f96a714faa52c",
                "0x00000000000000000000000047eeb2eac350e1923b8cbdfa4396a077b36e62a0"
            ],
            "data":"0x0000000000000000000000000000000000000000000000001bc16d674ec80000",
            "blockNumber":"0x2",
            "transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
            "transactionIndex":"0x0",
            "blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
            "logIndex":"0x0",
            "removed":false
        },
        {
            "address":"0x5fd55a1b9fc24967c4db09c513c3ba0dfa7ff687",
            "topics":[
                "0x282dd1817b996776123a00596764d4d54cc16460c9854f7a23f6be020ba0463d",
                "0x00000000000000000000000049837e63c983fea83ab098e0107f96a714faa52c"
            ],
            "data":"0x00000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000000de0b6b3a7640000636861696e2f6273630000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002a30783742374166373731374139354244364631636439303534653964303541623644623962643335386100000000000000000000000000000000000000000000",
            "blockNumber":"0x2",
            "transactionHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
            "transactionIndex":"0x0",
            "blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
            "logIndex":"0x1",
            "removed":false
        }
    ]`
	var logs []*ethtypes.Log
	if err := json.Unmarshal([]byte(data), &logs); err != nil {
		panic(err)
	}
	return logs
}
