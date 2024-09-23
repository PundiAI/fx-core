package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
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
	isIBC         bool
	target        string
	Prefix        string
	SourcePort    string
	SourceChannel string
}

func ParseFxTarget(targetStr string) FxTarget {
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

func GetIbcDenomTrace(denom string, channelIBC string) (ibctransfertypes.DenomTrace, error) {
	channelPath, err := hex.DecodeString(channelIBC)
	if err != nil {
		return ibctransfertypes.DenomTrace{}, errorsmod.Wrapf(err, "decode hex channel-ibc err")
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
