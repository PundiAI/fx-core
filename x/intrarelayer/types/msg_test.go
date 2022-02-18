package types

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/functionx/fx-core/tests"

	"github.com/ethereum/go-ethereum/common"
)

type MsgsTestSuite struct {
	suite.Suite
}

func TestMsgsTestSuite(t *testing.T) {
	suite.Run(t, new(MsgsTestSuite))
}

func (suite *MsgsTestSuite) TestMsgConvertCoinGetters() {
	msgInvalid := MsgConvertCoin{}
	msg := NewMsgConvertCoin(
		sdk.NewCoin("test", sdk.NewInt(100)),
		tests.GenerateAddress(),
		sdk.AccAddress(tests.GenerateAddress().Bytes()),
	)
	suite.Require().Equal(RouterKey, msg.Route())
	suite.Require().Equal(TypeMsgConvertCoin, msg.Type())
	suite.Require().NotNil(msgInvalid.GetSignBytes())
	suite.Require().Nil(msgInvalid.GetSigners())
	suite.Require().NotNil(msg.GetSigners())
}

func (suite *MsgsTestSuite) TestMsgConvertCoinNew() {
	testCases := []struct {
		msg        string
		coin       sdk.Coin
		receiver   common.Address
		sender     sdk.AccAddress
		expectPass bool
	}{
		{
			"msg convert coin - pass",
			sdk.NewCoin("test", sdk.NewInt(100)),
			tests.GenerateAddress(),
			sdk.AccAddress(tests.GenerateAddress().Bytes()),
			true,
		},
	}

	for i, tc := range testCases {
		tx := NewMsgConvertCoin(tc.coin, tc.receiver, tc.sender)
		err := tx.ValidateBasic()

		if tc.expectPass {
			suite.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.msg)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.msg)
		}
	}
}

func (suite *MsgsTestSuite) TestMsgConvertCoin() {
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
				Amount: sdk.NewInt(100),
			},
			"0x0000",
			tests.GenerateAddress().String(),
			false,
		},
		{
			"negative coin amount",
			sdk.Coin{
				Denom:  "coin",
				Amount: sdk.NewInt(-100),
			},
			"0x0000",
			tests.GenerateAddress().String(),
			false,
		},
		{
			"msg convert coin - invalid sender",
			sdk.NewCoin("coin", sdk.NewInt(100)),
			tests.GenerateAddress().String(),
			"evmosinvalid",
			false,
		},
		{
			"msg convert coin - invalid receiver",
			sdk.NewCoin("coin", sdk.NewInt(100)),
			"0x0000",
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			false,
		},
		{
			"msg convert coin - pass",
			sdk.NewCoin("coin", sdk.NewInt(100)),
			tests.GenerateAddress().String(),
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			true,
		},
		{
			"msg convert coin - pass with `intrarelayer/` denom",
			sdk.NewCoin("intrarelayer/0xdac17f958d2ee523a2206206994597c13d831ec7", sdk.NewInt(100)),
			tests.GenerateAddress().String(),
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			true,
		},
		{
			"msg convert coin - pass with `ibc/{hash}` denom",
			sdk.NewCoin("ibc/7F1D3FCF4AE79E1554D670D1AD949A9BA4E4A3C76C63093E17E446A46061A7A2", sdk.NewInt(100)),
			tests.GenerateAddress().String(),
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			true,
		},
	}

	for i, tc := range testCases {
		tx := MsgConvertCoin{tc.coin, tc.receiver, tc.sender}
		err := tx.ValidateBasic()

		if tc.expectPass {
			suite.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.msg)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.msg)
		}
	}
}

func (suite *MsgsTestSuite) TestMsgConvertFIP20Getters() {
	msgInvalid := MsgConvertFIP20{}
	priKey := tests.NewPriKey()
	pubKey, _ := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, priKey.PubKey())
	sender := sdk.AccAddress(priKey.PubKey().Address())
	msg := NewMsgConvertFIP20(
		sdk.NewInt(100),
		sender,
		sdk.AccAddress(tests.GenerateAddress().Bytes()),
		tests.GenerateAddress(),
		pubKey,
	)
	suite.Require().Equal(RouterKey, msg.Route())
	suite.Require().Equal(TypeMsgConvertFIP20, msg.Type())
	suite.Require().NotNil(msgInvalid.GetSignBytes())
	suite.Require().NotNil(msg.GetSigners())
}

func (suite *MsgsTestSuite) TestMsgConvertFIP20New() {
	testCases := []struct {
		msg        string
		amount     sdk.Int
		receiver   sdk.AccAddress
		contract   common.Address
		sender     cryptotypes.PrivKey
		expectPass bool
	}{
		{
			"msg convert erc20 - pass",
			sdk.NewInt(100),
			sdk.AccAddress(tests.GenerateAddress().Bytes()),
			tests.GenerateAddress(),
			tests.NewPriKey(),
			true,
		},
	}

	for i, tc := range testCases {
		pubKey, _ := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, tc.sender.PubKey())
		tx := NewMsgConvertFIP20(tc.amount, sdk.AccAddress(tc.sender.PubKey().Address()), tc.receiver, tc.contract, pubKey)
		err := tx.ValidateBasic()

		if tc.expectPass {
			suite.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.msg)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.msg)
		}
	}
}

func (suite *MsgsTestSuite) TestMsgConvertFIP20() {
	testCases := []struct {
		msg        string
		amount     sdk.Int
		receiver   string
		contract   string
		public     []byte
		expectPass bool
	}{
		{
			"invalid contract hex address",
			sdk.NewInt(100),
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			sdk.AccAddress{}.String(),
			tests.NewPriKey().PubKey().Bytes(),
			false,
		},
		{
			"negative coin amount",
			sdk.NewInt(-100),
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			tests.GenerateAddress().String(),
			tests.NewPriKey().PubKey().Bytes(),
			false,
		},
		{
			"invalid receiver address",
			sdk.NewInt(100),
			sdk.AccAddress{}.String(),
			tests.GenerateAddress().String(),
			tests.NewPriKey().PubKey().Bytes(),
			false,
		},
		{
			"invalid sender address",
			sdk.NewInt(100),
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			tests.GenerateAddress().String(),
			[]byte{0x0, 0x1},
			false,
		},
		{
			"msg convert erc20 - pass",
			sdk.NewInt(100),
			sdk.AccAddress(tests.GenerateAddress().Bytes()).String(),
			tests.GenerateAddress().String(),
			tests.NewPriKey().PubKey().Bytes(),
			true,
		},
	}

	for i, tc := range testCases {
		key := &secp256k1.PubKey{Key: tc.public}
		pubKey, _ := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, key)
		sender := string(tc.public)
		if len(tc.public) > 10 {
			sender = sdk.AccAddress(key.Address()).String()
		}
		tx := MsgConvertFIP20{tc.contract, tc.amount, tc.receiver, sender, pubKey}
		err := tx.ValidateBasic()

		if tc.expectPass {
			suite.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.msg)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.msg)
		}
	}
}
