package keeper_test

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/testutil/helpers"
	"github.com/functionx/fx-core/v8/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestBridgeCallHandler() {
	testCases := []struct {
		Name              string
		Msg               types.MsgBridgeCallClaim
		TokenIsNativeCoin []bool
		Success           bool
		CallContract      bool
	}{
		{
			Name: "success - token",
			Msg: types.MsgBridgeCallClaim{
				ChainName:      suite.chainName,
				BridgerAddress: helpers.GenAccAddress().String(),
				EventNonce:     1,
				BlockHeight:    1,
				Sender:         helpers.GenExternalAddr(suite.chainName),
				Refund:         helpers.GenExternalAddr(suite.chainName),
				TokenContracts: []string{
					helpers.GenExternalAddr(suite.chainName),
					helpers.GenExternalAddr(suite.chainName),
				},
				Amounts: []sdkmath.Int{
					helpers.NewRandAmount(),
					helpers.NewRandAmount(),
				},
				To:       helpers.GenExternalAddr(suite.chainName),
				Data:     "",
				Value:    sdkmath.ZeroInt(),
				Memo:     "",
				TxOrigin: helpers.GenExternalAddr(suite.chainName),
			},
			TokenIsNativeCoin: []bool{true, true},
			Success:           true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.Name, func() {
			_, _, erc20Addrs := suite.BridgeCallClaimInitialize(tc.Msg, tc.TokenIsNativeCoin)

			err := suite.Keeper().BridgeCallHandler(suite.Ctx, &tc.Msg)
			if tc.Success {
				suite.Require().NoError(err)
				if !tc.CallContract {
					for i, addr := range erc20Addrs {
						suite.CheckBalanceOf(addr, tc.Msg.GetToAddr(), tc.Msg.Amounts[i].BigInt())
					}
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) BridgeCallClaimInitialize(msg types.MsgBridgeCallClaim, tokenIsNativeCoin []bool) (baseDenoms, bridgeDenoms []string, erc20Addrs []common.Address) {
	suite.Require().Equal(len(tokenIsNativeCoin), len(msg.TokenContracts))

	baseDenoms = make([]string, 0, len(msg.TokenContracts))
	bridgeDenoms = make([]string, 0, len(msg.TokenContracts))
	erc20Addrs = make([]common.Address, 0, len(msg.TokenContracts))
	for i, c := range msg.TokenContracts {
		baseDenom := helpers.NewRandDenom()
		bridgeDenom := types.NewBridgeDenom(suite.chainName, c)
		suite.SetToken(strings.ToUpper(baseDenom), bridgeDenom)
		suite.AddBridgeToken(c, strings.ToLower(baseDenom))
		erc20Addr := suite.AddTokenPair(baseDenom, tokenIsNativeCoin[i])

		baseDenoms = append(baseDenoms, baseDenom)
		bridgeDenoms = append(bridgeDenoms, bridgeDenom)
		erc20Addrs = append(erc20Addrs, erc20Addr)

		if !tokenIsNativeCoin[i] {
			suite.MintTokenToModule(suite.chainName, sdk.NewCoin(bridgeDenom, msg.Amounts[i]))
		}
	}
	return baseDenoms, bridgeDenoms, erc20Addrs
}
