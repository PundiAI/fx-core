package ante_test

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"

	"github.com/functionx/fx-core/v3/ante"
	fxtypes "github.com/functionx/fx-core/v3/types"
)

func (suite *AnteTestSuite) TestMempoolFeeDecorator() {

	clientCtx := NewClientCtx()
	txBuilder := clientCtx.TxConfig.NewTxBuilder()
	mfd := ante.NewMempoolFeeDecorator([]string{
		sdk.MsgTypeURL(&ibcchanneltypes.MsgRecvPacket{}),
		sdk.MsgTypeURL(&ibcchanneltypes.MsgAcknowledgement{}),
		sdk.MsgTypeURL(&ibcclienttypes.MsgUpdateClient{}),
	}, 300_000)
	antehandler := sdk.ChainAnteDecorators(mfd)
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	msg := testdata.NewTestMsg(addr1)
	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()

	suite.Require().NoError(txBuilder.SetMsgs(msg))
	txBuilder.SetFeeAmount(feeAmount)
	txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	ctx := suite.GetContext(3)

	tx, err := suite.CreateEmptyTestTx(txBuilder, privs, accNums, accSeqs)
	suite.Require().NoError(err)

	// Set high gas price so standard test fee fails
	minGasPrice := []sdk.DecCoin{sdk.NewDecCoinFromDec(fxtypes.DefaultDenom, sdk.NewDec(200))}
	ctx = ctx.WithMinGasPrices(minGasPrice).WithIsCheckTx(true)

	// antehandler errors with insufficient fees
	_, err = antehandler(ctx, tx, false)
	suite.Require().Error(err, "expected error due to low fee")

	// ensure no fees for certain IBC msgs
	suite.Require().NoError(txBuilder.SetMsgs(
		ibcchanneltypes.NewMsgRecvPacket(ibcchanneltypes.Packet{}, nil, ibcclienttypes.Height{}, sdk.AccAddress{}.String()),
	))

	oracleTx, err := suite.CreateEmptyTestTx(txBuilder, privs, accNums, accSeqs)
	suite.Require().NoError(err)
	_, err = antehandler(ctx, oracleTx, false)
	suite.Require().NoError(err, "expected min fee bypass for IBC messages")

	ctx = ctx.WithIsCheckTx(false)

	// antehandler should not error since we do not check min gas prices in DeliverTx
	_, err = antehandler(ctx, tx, false)
	suite.Require().NoError(err, "unexpected error during DeliverTx")
}
