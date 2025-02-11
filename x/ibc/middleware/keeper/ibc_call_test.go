package keeper_test

import (
	"encoding/hex"
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
	"github.com/pundiai/fx-core/v8/x/ibc/middleware/types"
)

func (suite *KeeperTestSuite) TestHandlerIbcCall() {
	testCases := []struct {
		name     string
		malleate func(packet *transfertypes.FungibleTokenPacketData) *types.IbcCallEvmPacket
		expPass  bool
		expCheck func(packet transfertypes.FungibleTokenPacketData)
		expError string
	}{
		{
			name: "success - empty",
			malleate: func(packet *transfertypes.FungibleTokenPacketData) *types.IbcCallEvmPacket {
				return &types.IbcCallEvmPacket{To: helpers.GenHexAddress().String(), Data: "", Value: sdkmath.ZeroInt()}
			},
			expPass: true,
		},
		{
			name: "success - origin token",
			malleate: func(packet *transfertypes.FungibleTokenPacketData) *types.IbcCallEvmPacket {
				amt, _ := sdkmath.NewIntFromString(packet.Amount)
				receiver := sdk.MustAccAddressFromBech32(packet.Receiver)
				suite.MintToken(receiver.Bytes(), sdk.NewCoin(packet.Denom, amt))
				return &types.IbcCallEvmPacket{To: helpers.GenHexAddress().String(), Data: "", Value: amt}
			},
			expPass: true,
			expCheck: func(packet transfertypes.FungibleTokenPacketData) {
				var mp types.MemoPacket
				err := suite.App.AppCodec().UnmarshalInterfaceJSON([]byte(packet.Memo), &mp)
				suite.Require().NoError(err)
				callEvmPacket := mp.(*types.IbcCallEvmPacket)

				toAddr := common.HexToAddress(callEvmPacket.To)
				amt, _ := sdkmath.NewIntFromString(packet.Amount)
				suite.AssertBalance(toAddr.Bytes(), sdk.NewCoin(packet.Denom, amt))
			},
		},
		{
			name: "success - erc20 token",
			malleate: func(packet *transfertypes.FungibleTokenPacketData) *types.IbcCallEvmPacket {
				symbol := helpers.NewRandSymbol()
				packet.Denom = strings.ToLower(symbol)

				amt, _ := sdkmath.NewIntFromString(packet.Amount)
				suite.MintTokenToModule(erc20types.ModuleName, sdk.NewCoin(packet.Denom, amt))
				suite.AddBridgeToken(ethtypes.ModuleName, symbol, true, false)

				receiver := sdk.MustAccAddressFromBech32(packet.Receiver)
				erc20Token := suite.GetERC20Token(packet.Denom)
				suite.erc20TokenSuite.WithContract(erc20Token.GetERC20Contract()).
					MintFromERC20Module(suite.Ctx, common.BytesToAddress(receiver.Bytes()), amt.BigInt())

				sender := sdk.MustAccAddressFromBech32(packet.Sender)
				data := helpers.PackERC20Transfer(common.BytesToAddress(sender.Bytes()), amt.BigInt())
				return &types.IbcCallEvmPacket{To: erc20Token.GetErc20Address(), Data: hex.EncodeToString(data), Value: sdkmath.ZeroInt()}
			},
			expPass: true,
			expCheck: func(packet transfertypes.FungibleTokenPacketData) {
				senderAddr := sdk.MustAccAddressFromBech32(packet.Sender)

				erc20Token := suite.GetERC20Token(packet.Denom)
				amt := suite.erc20TokenSuite.WithContract(erc20Token.GetERC20Contract()).
					BalanceOf(suite.Ctx, common.BytesToAddress(senderAddr.Bytes()))
				suite.Require().Equal(packet.Amount, amt.String())
			},
		},
		{
			name: "failed - evm revert",
			malleate: func(packet *transfertypes.FungibleTokenPacketData) *types.IbcCallEvmPacket {
				symbol := helpers.NewRandSymbol()
				packet.Denom = strings.ToLower(symbol)
				suite.AddBridgeToken(ethtypes.ModuleName, symbol, true, false)

				erc20Token := suite.GetERC20Token(packet.Denom)
				data := helpers.PackERC20Mint(helpers.GenHexAddress(), big.NewInt(0))
				return &types.IbcCallEvmPacket{To: erc20Token.GetErc20Address(), Data: hex.EncodeToString(data), Value: sdkmath.ZeroInt()}
			},
			expPass:  false,
			expError: "Ownable: caller is not the owner: evm transaction execution failed",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			sender := helpers.GenAccAddress().String()
			channelID := "channel-0"
			receiver := types.IntermediateSender(transfertypes.PortID, channelID, sender)
			amount := sdkmath.OneInt()
			packet := transfertypes.FungibleTokenPacketData{
				Denom:    fxtypes.DefaultDenom,
				Amount:   amount.String(),
				Sender:   sender,
				Receiver: sdk.AccAddress(receiver.Bytes()).String(),
			}

			memo := tc.malleate(&packet)
			bz, err := suite.App.AppCodec().MarshalInterfaceJSON(memo)
			suite.Require().NoError(err)
			packet.Memo = string(bz)

			err = suite.ibcMiddleware.Keeper.HandlerIbcCall(suite.Ctx, transfertypes.ModuleName, channelID, packet)
			if !tc.expPass {
				suite.Require().Error(err)
				suite.Require().ErrorContains(err, tc.expError)
				return
			}
			suite.Require().NoError(err)

			if tc.expCheck != nil {
				tc.expCheck(packet)
			}
		})
	}
}
