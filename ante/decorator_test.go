package ante_test

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v2/ante"
	fxtypes "github.com/functionx/fx-core/v2/types"
	"github.com/functionx/fx-core/v2/x/erc20/types"
)

func (suite *AnteTestSuite) TestMsgInterceptDecorator() {
	clientCtx := NewClientCtx()
	txBuilder := clientCtx.TxConfig.NewTxBuilder()

	mfd := ante.NewMsgInterceptDecorator(map[int64][]string{
		10: fxtypes.SupportDenomManyToOneMsgTypes(),
	})

	antehandler := sdk.ChainAnteDecorators(mfd)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{suite.privateKey}, []uint64{0}, []uint64{0}
	ctx := suite.GetContext(3)

	suite.Require().NoError(txBuilder.SetMsgs(testdata.NewTestMsg(suite.GetAccAddress())))
	tx, err := suite.CreateEmptyTestTx(txBuilder, privs, accNums, accSeqs)
	suite.Require().NoError(err)

	_, err = antehandler(ctx, tx, false)
	suite.Require().NoError(err)

	suite.Require().NoError(txBuilder.SetMsgs(
		types.NewMsgConvertDenom(suite.GetAccAddress(), suite.GetAccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(100)), ""),
	))

	tx2, err := suite.CreateEmptyTestTx(txBuilder, privs, accNums, accSeqs)
	suite.Require().NoError(err)
	_, err = antehandler(ctx, tx2, false)
	suite.Require().Error(err)
}
