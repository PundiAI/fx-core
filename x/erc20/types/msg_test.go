package types_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	_ "github.com/functionx/fx-core/app/fxcore"
	"github.com/functionx/fx-core/x/erc20/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgConvertCoinValidateBasic(t *testing.T) {
	tests := []struct {
		name      string
		msg       types.MsgConvertCoin
		pass      bool
		err       error
		errReason string
	}{
		{
			"valid",
			types.MsgConvertCoin{
				Coin:     sdk.Coin{Denom: "ABB", Amount: sdk.NewInt(1)},
				Receiver: "0xbbB31708Bfe3b271210Ae40b1434FB897409004b",
				Sender:   "fx1hwe3wz9luwe8zgg2us93gd8m396qjqztky35t5",
			},
			true,
			nil,
			"",
		},
		{
			"valid - empty receiver",
			types.MsgConvertCoin{
				Coin:     sdk.Coin{Denom: "ABB", Amount: sdk.NewInt(1)},
				Receiver: "",
				Sender:   "fx1hwe3wz9luwe8zgg2us93gd8m396qjqztky35t5",
			},
			false,
			sdkerrors.ErrInvalidAddress,
			sdkerrors.Wrap(
				sdkerrors.ErrInvalidAddress,
				fmt.Errorf("invalid receiver address empty").Error(),
			).Error(),
		},
		{
			"valid - lowercase receiver",
			types.MsgConvertCoin{
				Coin:     sdk.Coin{Denom: "ABB", Amount: sdk.NewInt(1)},
				Receiver: "0xbbb31708bfe3b271210ae40b1434fb897409004b",
				Sender:   "fx1hwe3wz9luwe8zgg2us93gd8m396qjqztky35t5",
			},
			false,
			sdkerrors.ErrInvalidAddress,
			sdkerrors.Wrap(
				sdkerrors.ErrInvalidAddress,
				fmt.Errorf("invalid receiver address invalid address got:%s, expected:%s",
					"0xbbb31708bfe3b271210ae40b1434fb897409004b",
					"0xbbB31708Bfe3b271210Ae40b1434FB897409004b").Error(),
			).Error(),
		},
		{
			"valid - uppercase receiver",
			types.MsgConvertCoin{
				Coin:     sdk.Coin{Denom: "ABB", Amount: sdk.NewInt(1)},
				Receiver: "0xBBB31708BFE3B271210AE40B1434FB897409004B",
				Sender:   "fx1hwe3wz9luwe8zgg2us93gd8m396qjqztky35t5",
			},
			false,
			sdkerrors.ErrInvalidAddress,
			sdkerrors.Wrap(
				sdkerrors.ErrInvalidAddress,
				fmt.Errorf("invalid receiver address invalid address got:%s, expected:%s",
					"0xBBB31708BFE3B271210AE40B1434FB897409004B",
					"0xbbB31708Bfe3b271210Ae40b1434FB897409004b").Error(),
			).Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.pass {
				require.NoError(t, err)
			} else {
				require.NotNil(t, err)
				require.ErrorIs(t, err, tt.err)
				require.Equal(t, tt.errReason, err.Error())
			}
		})
	}
}

func TestValidateBasic(t *testing.T) {
}
