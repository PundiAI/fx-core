package types

import (
	"encoding/hex"
	"fmt"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
)

const (
	LegacyERC20Target = "module/evm"
	LegacyChainPrefix = "chain/"

	ERC20Target = "erc20"
	IBCPrefix   = "ibc/"
)

type FxTarget struct {
	isIBC         bool
	target        string
	Prefix        string
	SourcePort    string
	SourceChannel string
}

func ParseFxTarget(targetStr string) FxTarget {
	// module evm
	if targetStr == LegacyERC20Target {
		targetStr = ERC20Target
	}
	targetStr = strings.TrimPrefix(targetStr, LegacyChainPrefix)

	// ibc prefix
	if strings.HasPrefix(targetStr, IBCPrefix) {
		ibcData := strings.Split(targetStr, "/")
		if len(ibcData) == 3 {
			// ibc/{channelId}/{prefix} -> ibc/0/px
			return FxTarget{
				isIBC:         true,
				Prefix:        ibcData[2],
				SourcePort:    ibctransfertypes.ModuleName,
				SourceChannel: fmt.Sprintf("%s%s", channeltypes.ChannelPrefix, ibcData[1]),
			}
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
		return FxTarget{
			isIBC:         true,
			Prefix:        ibcData[0],
			SourcePort:    ibcData[1],
			SourceChannel: ibcData[2],
		}
	}

	return FxTarget{isIBC: false, target: targetStr}
}

func (i FxTarget) GetTarget() string {
	if i.isIBC {
		return fmt.Sprintf("%s/%s", i.SourceChannel, i.Prefix)
	}
	return i.target
}

func (i FxTarget) IsIBC() bool {
	return i.isIBC
}

func GetIbcDenomTrace(denom string, channelIBC string) (ibctransfertypes.DenomTrace, error) {
	channelPath, err := hex.DecodeString(channelIBC)
	if err != nil {
		return ibctransfertypes.DenomTrace{}, sdkerrors.Wrapf(err, "decode channel ibc err")
	}

	// todo need check path
	path := string(channelPath)
	return ibctransfertypes.DenomTrace{
		Path:      path,
		BaseDenom: denom,
	}, nil
}
