package types

import (
	"encoding/hex"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
)

const (
	LegacyERC20Target = "module/evm"
	ERC20Target       = "erc20"
	GravityTarget     = "gravity"
	EthTarget         = "eth"

	LegacyChainPrefix = "chain/"
	IBCPrefix         = "ibc/"
)

type FxTarget struct {
	isIBC        bool
	target       string
	Bech32Prefix string
	IBCChannel   string
}

func ParseFxTarget(targetStr string, isHex ...bool) (*FxTarget, error) {
	if len(isHex) > 0 && isHex[0] {
		// ignore hex decode error
		targetByte, _ := hex.DecodeString(targetStr)
		targetStr = string(targetByte)
	}
	targetStr = strings.TrimPrefix(targetStr, LegacyChainPrefix)

	if targetStr == LegacyERC20Target || targetStr == ERC20Target || targetStr == "" {
		return &FxTarget{isIBC: false, target: ""}, nil
	}

	// ibc prefix
	targetArr := strings.Split(targetStr, "/")
	if len(targetArr) == 1 {
		if targetStr == GravityTarget {
			targetStr = EthTarget
		}
		// target is module name
		return &FxTarget{isIBC: false, target: targetStr}, nil
	}

	if len(targetArr) == 4 {
		// ibc/{prefix}/transfer/channel-{id} -> ibc/px/transfer/channel-0
		if targetArr[2] != transfertypes.ModuleName {
			return nil, fmt.Errorf("invalid target: %s", targetStr)
		}
		targetArr[2] = targetArr[1]
		targetArr[1] = strings.TrimPrefix(targetArr[3], channeltypes.ChannelPrefix)
		targetArr = targetArr[0:3]
	}
	if len(targetArr) == 3 {
		// ibc/{channelId}/{prefix} -> ibc/0/px
		fxTarget := &FxTarget{
			isIBC:        true,
			Bech32Prefix: targetArr[2],
			IBCChannel:   fmt.Sprintf("%s%s", channeltypes.ChannelPrefix, targetArr[1]),
		}
		if err := fxTarget.IBCValidate(); err != nil {
			return nil, err
		}
		return fxTarget, nil
	}

	return nil, fmt.Errorf("invalid target: %s", targetStr)
}

func (i FxTarget) GetModuleName() string {
	if i.isIBC {
		return EthTarget
	}
	return i.target
}

func (i FxTarget) IsIBC() bool {
	return i.isIBC
}

func (i FxTarget) IBCValidate() error {
	if !channeltypes.IsValidChannelID(i.IBCChannel) {
		return fmt.Errorf("invalid channel id: %s", i.IBCChannel)
	}
	if len(strings.TrimSpace(i.Bech32Prefix)) == 0 {
		return fmt.Errorf("empty bech32 prefix: %s", i.Bech32Prefix)
	}
	return nil
}

func (i FxTarget) ReceiveAddrToStr(receive sdk.AccAddress) (receiveAddrStr string, err error) {
	receiveAddrStr, err = bech32.ConvertAndEncode(i.Bech32Prefix, receive)
	if err != nil {
		return "", sdkerrors.ErrInvalidAddress.Wrapf("prefix: %s error: %s", i.Bech32Prefix, err)
	}
	return receiveAddrStr, nil
}

func (i FxTarget) ValidateExternalAddr(receive string) (err error) {
	if i.isIBC {
		_, err = sdk.GetFromBech32(receive, i.Bech32Prefix)
	} else {
		err = ValidateExternalAddr(i.GetModuleName(), receive)
	}
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid receive address: %s", err)
	}
	return nil
}
