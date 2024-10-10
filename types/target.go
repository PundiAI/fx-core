package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/contract"
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
	isIBC         bool
	target        string
	Prefix        string
	SourcePort    string
	SourceChannel string
}

func ParseFxTarget(targetStr string, isHex ...bool) FxTarget {
	if len(isHex) > 0 && isHex[0] {
		// ignore hex decode error
		targetByte, _ := hex.DecodeString(targetStr)
		targetStr = string(targetByte)
	}
	// module evm
	if targetStr == LegacyERC20Target {
		return FxTarget{isIBC: false, target: ERC20Target}
	}
	targetStr = strings.TrimPrefix(targetStr, LegacyChainPrefix)
	// cross-chain
	if targetStr == GravityTarget {
		return FxTarget{isIBC: false, target: EthTarget}
	}

	// ibc prefix
	if strings.HasPrefix(targetStr, IBCPrefix) {
		ibcData := strings.Split(targetStr, "/")
		if len(ibcData) == 3 {
			// ibc/{channelId}/{prefix} -> ibc/0/px
			fxTarget := FxTarget{
				isIBC:         true,
				Prefix:        ibcData[2],
				SourcePort:    ibctransfertypes.ModuleName,
				SourceChannel: fmt.Sprintf("%s%s", channeltypes.ChannelPrefix, ibcData[1]),
			}
			if !fxTarget.IBCValidate() {
				return FxTarget{isIBC: false, target: targetStr}
			}
			return fxTarget
		} else if len(ibcData) == 4 {
			// ibc/{prefix}/transfer/channel-{id} -> ibc/px/transfer/channel-0
			targetStr = strings.TrimPrefix(targetStr, IBCPrefix)
		} else {
			return FxTarget{isIBC: false, target: targetStr}
		}
	}

	// px/transfer/channel-0
	ibcData := strings.Split(targetStr, "/")
	if len(ibcData) == 3 {
		fxTarget := FxTarget{
			isIBC:         true,
			Prefix:        ibcData[0],
			SourcePort:    ibcData[1],
			SourceChannel: ibcData[2],
		}
		if !fxTarget.IBCValidate() {
			return FxTarget{isIBC: false, target: targetStr}
		}
		return fxTarget
	}

	return FxTarget{isIBC: false, target: targetStr}
}

func (i FxTarget) GetTarget() string {
	if i.isIBC {
		return fmt.Sprintf("%s/%s", i.SourceChannel, i.Prefix)
	}
	return i.target
}

func (i FxTarget) String() string {
	if i.isIBC {
		return fmt.Sprintf("ibc/%s/%s", strings.TrimPrefix(i.SourceChannel, channeltypes.ChannelPrefix), i.Prefix)
	}
	return i.target
}

func (i FxTarget) IsIBC() bool {
	return i.isIBC
}

func (i FxTarget) IBCValidate() bool {
	if !i.isIBC {
		return false
	}
	if i.SourcePort != ibctransfertypes.ModuleName {
		return false
	}
	if !channeltypes.IsValidChannelID(i.SourceChannel) {
		return false
	}
	if len(strings.TrimSpace(i.Prefix)) == 0 {
		return false
	}
	return true
}

func (i FxTarget) ReceiveAddrToStr(receive sdk.AccAddress) (receiveAddrStr string, err error) {
	if strings.ToLower(i.Prefix) == contract.EthereumAddressPrefix {
		return common.BytesToAddress(receive.Bytes()).String(), nil
	}
	return bech32.ConvertAndEncode(i.Prefix, receive)
}

func GetIbcDenomTrace(denom string, channelIBC string) (ibctransfertypes.DenomTrace, error) {
	channelPath, err := hex.DecodeString(channelIBC)
	if err != nil {
		return ibctransfertypes.DenomTrace{}, fmt.Errorf("invalid channel-ibc: %w", err)
	}

	// transfer/channel-0
	path := string(channelPath)
	if len(path) > 0 {
		pathSplit := strings.Split(path, "/")
		if len(pathSplit) != 2 {
			return ibctransfertypes.DenomTrace{}, errors.New("invalid params channel-ibc")
		}
		if pathSplit[0] != "transfer" {
			return ibctransfertypes.DenomTrace{}, errors.New("invalid source port")
		}
		if !channeltypes.IsValidChannelID(pathSplit[1]) {
			return ibctransfertypes.DenomTrace{}, errors.New("invalid source channel")
		}
	}

	return ibctransfertypes.DenomTrace{
		Path:      path,
		BaseDenom: denom,
	}, nil
}
