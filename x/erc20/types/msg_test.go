package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/erc20/types"
)

type MsgTestSuite struct {
	suite.Suite
}

func TestMsgTestSuite(t *testing.T) {
	suite.Run(t, new(MsgTestSuite))
}

func (suite *MsgTestSuite) TestMsgConvertCoin() {
	testCases := []struct {
		msg        string
		coin       sdk.Coin
		receiver   string
		sender     string
		expectPass bool
	}{
		{
			"invalid denom",
			sdk.Coin{
				Denom:  "",
				Amount: sdkmath.NewInt(100),
			},
			"0x0000",
			helpers.GenHexAddress().String(),
			false,
		},
		{
			"negative coin amount",
			sdk.Coin{
				Denom:  "coin",
				Amount: sdkmath.NewInt(-100),
			},
			"0x0000",
			helpers.GenHexAddress().String(),
			false,
		},
		{
			"msg convert coin - invalid sender",
			sdk.NewCoin("coin", sdkmath.NewInt(100)),
			helpers.GenHexAddress().String(),
			"evmosinvalid",
			false,
		},
		{
			"msg convert coin - invalid receiver",
			sdk.NewCoin("coin", sdkmath.NewInt(100)),
			"0x0000",
			helpers.GenAccAddress().String(),
			false,
		},
		{
			"msg convert coin - pass",
			sdk.NewCoin("coin", sdkmath.NewInt(100)),
			helpers.GenHexAddress().String(),
			helpers.GenAccAddress().String(),
			true,
		},
		{
			"msg convert coin - pass with `erc20/` denom",
			sdk.NewCoin("erc20/0xdac17f958d2ee523a2206206994597c13d831ec7", sdkmath.NewInt(100)),
			helpers.GenHexAddress().String(),
			helpers.GenAccAddress().String(),
			true,
		},
		{
			"msg convert coin - pass with `ibc/{hash}` denom",
			sdk.NewCoin("ibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A2", sdkmath.NewInt(100)),
			helpers.GenHexAddress().String(),
			helpers.GenAccAddress().String(),
			true,
		},
	}

	for i, tc := range testCases {
		tx := types.MsgConvertCoin{Coin: tc.coin, Receiver: tc.receiver, Sender: tc.sender}
		err := tx.ValidateBasic()

		if tc.expectPass {
			suite.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.msg)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.msg)
		}
	}
}
