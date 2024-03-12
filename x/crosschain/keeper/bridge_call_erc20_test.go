package keeper_test

import (
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/keeper"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestKeeper_BridgeCallERC20Handler() {
	suite.TestMsgSetOracleSetConfirm()

	bridgeTokenAddr := helpers.GenerateAddress()
	asset, err := types.PackERC20AssetWithType(
		[]common.Address{
			bridgeTokenAddr,
		},
		[]*big.Int{
			big.NewInt(100),
		},
	)
	suite.NoError(err)
	_, assetBytes, err := types.UnpackAssetType(asset)
	suite.NoError(err)
	eventNonce := uint64(1)
	claim := types.MsgBridgeCallClaim{
		ChainName:  suite.chainName,
		Sender:     helpers.GenerateAddressByModule(suite.chainName),
		To:         helpers.GenerateAddressByModule(suite.chainName),
		Receiver:   helpers.GenerateAddressByModule(suite.chainName),
		DstChainId: types.FxcoreChainID,
		Message:    "",
		Value:      sdkmath.NewInt(0),
		GasLimit:   3000000,
		EventNonce: eventNonce,
	}
	err = suite.Keeper().BridgeCallERC20Handler(
		suite.ctx,
		assetBytes,
		claim.MustSender(),
		claim.MustTo(),
		claim.MustReceiver(),
		claim.DstChainId,
		claim.MustMessage(),
		claim.Value,
		claim.GasLimit,
		claim.EventNonce,
	)
	suite.NoError(err)
	refundRecord, found := suite.Keeper().GetRefundRecord(suite.ctx, eventNonce)
	suite.True(found)
	suite.Equal(refundRecord.Receiver, claim.Sender)

	event := suite.FindEvent(types.EventTypeBridgeCallRefund)
	suite.NotNil(event)
	suite.Equal(string(event.Attributes[0].GetKey()), types.AttributeKeyErrCause)
	suite.Equal(string(event.Attributes[0].GetValue()), "bridge token is not exist: invalid")
	suite.Equal(string(event.Attributes[1].GetKey()), types.AttributeKeyEventNonce)
	suite.Equal(string(event.Attributes[1].GetValue()), fmt.Sprint(eventNonce))
	suite.Equal(string(event.Attributes[2].GetKey()), types.AttributeKeyRefundAddress)
	suite.Equal(string(event.Attributes[2].GetValue()), claim.Sender)
	suite.ctx = suite.ctx.WithEventManager(sdk.NewEventManager())

	// add bridge token
	bridgeToken := fxtypes.AddressToStr(bridgeTokenAddr.Bytes(), suite.chainName)

	suite.addBridgeToken(bridgeToken, fxtypes.GetCrossChainMetadataManyToOne("test token", "TT", 18))

	suite.registerCoin(keeper.NewBridgeDenom(suite.chainName, bridgeToken))

	claim.EventNonce = 2
	err = suite.Keeper().BridgeCallERC20Handler(
		suite.ctx,
		assetBytes,
		claim.MustSender(),
		claim.MustTo(),
		claim.MustReceiver(),
		claim.DstChainId,
		claim.MustMessage(),
		claim.Value,
		claim.GasLimit,
		claim.EventNonce,
	)
	suite.NoError(err)
	refundRecord, found = suite.Keeper().GetRefundRecord(suite.ctx, claim.EventNonce)
	suite.False(found)
	suite.Nil(refundRecord)
}

func (suite *KeeperTestSuite) FindEvent(tp string) *sdk.Event {
	for _, event := range suite.ctx.EventManager().Events() {
		if event.Type == tp {
			return &event
		}
	}
	return nil
}
