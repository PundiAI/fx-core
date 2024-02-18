package tests

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	"github.com/functionx/fx-core/v7/testutil/helpers"
)

func (suite *IntegrationTest) ByPassFeeTest() {
	userA := helpers.NewEthPrivKey()
	userAAddr := userA.PubKey().Address().Bytes()
	initBalance := suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18))
	txResponse := suite.Send(userAAddr, initBalance)
	suite.Require().EqualValues(uint32(0), txResponse.Code)

	// 1. zero gasPrices for bypassing fee check
	minGasPrices := suite.network.Config.MinGasPrices
	config := suite.network.Config
	config.MinGasPrices = suite.NewCoin(sdk.ZeroInt()).String()
	suite.network.Config = config

	for i := 0; i < tmrand.Intn(5); i++ {
		broadcastTx := suite.BroadcastTx(userA, distributiontypes.NewMsgSetWithdrawAddress(userAAddr, userAAddr))
		suite.Require().EqualValues(uint32(0), broadcastTx.Code)
	}

	// check balance
	suite.CheckBalance(userAAddr, initBalance)

	// ok. reset gasPrices
	config.MinGasPrices = minGasPrices
	suite.network.Config = config
}
