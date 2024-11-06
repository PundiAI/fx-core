package integration

import (
	sdkmath "cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/functionx/fx-core/v8/testutil/helpers"
)

func (suite *IntegrationTest) ByPassFeeTest() {
	userA := helpers.NewSigner(helpers.NewEthPrivKey())
	userAAddr := userA.AccAddress()
	initBalance := suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18))
	txResponse := suite.Send(userAAddr, initBalance)
	suite.Require().EqualValues(uint32(0), txResponse.Code)

	// zero gasPrices for bypassing fee check
	zeroGasPrice := sdk.Coin{Amount: sdkmath.ZeroInt(), Denom: suite.defDenom}
	for i := 0; i < tmrand.Intn(5); i++ {
		suite.WithGasPrices(zeroGasPrice).BroadcastTx(userA,
			distributiontypes.NewMsgSetWithdrawAddress(userAAddr, userAAddr))
	}

	// check balance
	suite.EqualBalance(userAAddr, initBalance)
}
