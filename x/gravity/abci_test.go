package gravity

import (
	"fmt"
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/x/gravity/types"
)

func TestValSetPowerIsChanger(t *testing.T) {
	fmt.Printf("multiple:%v\n", math.MaxUint32)
	testDatas := []testCaseData{
		{
			latestPower:     []uint64{200, 200, 200, 200, 200},
			currentPower:    []uint64{200, 200, 200, 200, 200},
			isChange:        false,
			changeThreshold: sdk.NewDec(1).QuoInt(sdk.NewInt(10)),
		},
		{
			latestPower:     []uint64{12345, 12345, 12345, 12345},
			currentPower:    []uint64{12345, 12345, 12345, 12345, 10000},
			isChange:        true,
			changeThreshold: sdk.NewDec(1).QuoInt(sdk.NewInt(10)),
		},
	}
	for _, testData := range testDatas {
		latestValset := make([]*types.BridgeValidator, len(testData.latestPower))
		totalLatestPower := testData.totalLatestPower()
		for i, power := range testData.latestPower {
			latestValset[i] = &types.BridgeValidator{
				Power:      sdk.NewUint(power).MulUint64(math.MaxUint32).QuoUint64(totalLatestPower).Uint64(),
				EthAddress: fmt.Sprintf("a-%d", i),
			}
		}

		currentValset := make([]*types.BridgeValidator, len(testData.currentPower))
		totalCurrentPower := testData.totalCurrentPower()
		for i, power := range testData.currentPower {
			currentValset[i] = &types.BridgeValidator{
				Power:      sdk.NewUint(power).MulUint64(math.MaxUint32).QuoUint64(totalCurrentPower).Uint64(),
				EthAddress: fmt.Sprintf("a-%d", i),
			}
		}
		powerDiff, isChange := valSetPowerIsChanger(latestValset, currentValset, testData.changeThreshold)
		assert.Equal(t, testData.isChange, isChange, "powerDiff is change", "powerDiff", powerDiff, "changeThreshold", testData.changeThreshold)
	}
}

type testCaseData struct {
	latestPower     []uint64
	currentPower    []uint64
	isChange        bool
	changeThreshold sdk.Dec
}

func (data testCaseData) totalLatestPower() uint64 {
	var result uint64
	for _, power := range data.latestPower {
		result += power
	}
	return result
}

func (data testCaseData) totalCurrentPower() uint64 {
	var result uint64
	for _, power := range data.currentPower {
		result += power
	}
	return result
}
