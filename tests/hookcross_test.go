package tests

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
)

const (
	bscUSDToken     = "0x0000000000000000000000000000000000000001"
	polygonUSDToken = "0x0000000000000000000000000000000000000002"
	ethUSDToken     = "0x0000000000000000000000000000000000000003"
)

type HookCrossTestSuite struct {
	CrosschainERC20TestSuite
}

func TestERC20TestSuite(t *testing.T) {
	testSuite := NewTestSuite()
	hookCrossTestSuite := &HookCrossTestSuite{
		CrosschainERC20TestSuite: NewCrosschainERC20TestSuite(testSuite),
	}
	suite.Run(t, hookCrossTestSuite)
}

func (suite *HookCrossTestSuite) TestERC20HookCrossChain() {
	suite.InitCrossChain()
	suite.InitRegisterCoinUSDT()

	usdtTokenPair := suite.ERC20.TokenPair("usdt")

	beforeSendToFx := suite.ERC20.BalanceOf(usdtTokenPair.GetERC20Contract(), suite.BSCCrossChain.HexAddr())
	suite.BSCCrossChain.SendToFxClaim(bscUSDToken, sdk.NewInt(100).MulRaw(1e18), "module/evm")
	afterSendToFx := suite.ERC20.BalanceOf(usdtTokenPair.GetERC20Contract(), suite.BSCCrossChain.HexAddr())
	suite.Equal(big.NewInt(0).Sub(afterSendToFx, beforeSendToFx), sdk.NewInt(100).MulRaw(1e18).BigInt())

	suite.PolygonCrossChain.SendToFxClaim(polygonUSDToken, sdk.NewInt(100).MulRaw(1e18), "module/evm")

	beforeTransfer := suite.ERC20.BalanceOf(usdtTokenPair.GetERC20Contract(), suite.BSCCrossChain.HexAddr())
	beforeTransferAddr := suite.ERC20.BalanceOf(usdtTokenPair.GetERC20Contract(), suite.ERC20.HexAddr())
	suite.ERC20.TransferERC20(suite.BSCCrossChain.privKey, usdtTokenPair.GetERC20Contract(), suite.ERC20.HexAddr(), big.NewInt(100))
	afterTransfer := suite.ERC20.BalanceOf(usdtTokenPair.GetERC20Contract(), suite.BSCCrossChain.HexAddr())
	afterTransferAddr := suite.ERC20.BalanceOf(usdtTokenPair.GetERC20Contract(), suite.ERC20.HexAddr())
	suite.Equal(big.NewInt(0).Sub(beforeTransfer, afterTransfer), big.NewInt(0).Sub(afterTransferAddr, beforeTransferAddr))

	suite.ERC20.TransferERC20(suite.PolygonCrossChain.privKey, usdtTokenPair.GetERC20Contract(), suite.ERC20.HexAddr(), big.NewInt(100))

	suite.ERC20.TransferCrossChain(suite.ERC20.privKey, usdtTokenPair.GetERC20Contract(),
		suite.ERC20.HexAddr().String(), big.NewInt(50), big.NewInt(50), "chain/bsc")

	suite.ERC20.TransferCrossChain(suite.ERC20.privKey, usdtTokenPair.GetERC20Contract(),
		suite.ERC20.HexAddr().String(), big.NewInt(50), big.NewInt(50), "chain/polygon")

	resp, err := suite.BSCCrossChain.CrosschainQuery().GetPendingSendToExternal(suite.ctx,
		&crosschaintypes.QueryPendingSendToExternalRequest{ChainName: suite.BSCCrossChain.chainName, SenderAddress: suite.ERC20.Address().String()})
	suite.NoError(err)
	suite.Equal(1, len(resp.UnbatchedTransfers))
	suite.Equal(int64(50), resp.UnbatchedTransfers[0].Token.Amount.Int64())
	suite.Equal(int64(50), resp.UnbatchedTransfers[0].Fee.Amount.Int64())
	suite.Equal(suite.ERC20.Address().String(), resp.UnbatchedTransfers[0].Sender)
	suite.Equal(suite.ERC20.HexAddr().String(), resp.UnbatchedTransfers[0].DestAddress)

	resp, err = suite.PolygonCrossChain.CrosschainQuery().GetPendingSendToExternal(suite.ctx,
		&crosschaintypes.QueryPendingSendToExternalRequest{
			ChainName:     suite.PolygonCrossChain.chainName,
			SenderAddress: suite.ERC20.Address().String(),
		})
	suite.NoError(err)
	suite.Equal(1, len(resp.UnbatchedTransfers))
	suite.Equal(int64(50), resp.UnbatchedTransfers[0].Token.Amount.Int64())
	suite.Equal(int64(50), resp.UnbatchedTransfers[0].Fee.Amount.Int64())
	suite.Equal(suite.ERC20.Address().String(), resp.UnbatchedTransfers[0].Sender)
	suite.Equal(suite.ERC20.HexAddr().String(), resp.UnbatchedTransfers[0].DestAddress)
}
