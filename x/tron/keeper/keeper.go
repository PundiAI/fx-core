package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	crosschainkeeper "github.com/functionx/fx-core/v2/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v2/x/crosschain/types"
)

type Keeper struct {
	crosschainkeeper.Keeper
}

func NewKeeper(cdc codec.BinaryCodec, moduleName string, storeKey sdk.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper crosschaintypes.StakingKeeper, distributionKeeper crosschaintypes.DistributionKeeper, bankKeeper crosschaintypes.BankKeeper,
	ibcTransferKeeper crosschaintypes.IBCTransferKeeper, channelKeeper crosschaintypes.IBCChannelKeeper, erc20Keeper crosschaintypes.Erc20Keeper) Keeper {
	return Keeper{
		Keeper: crosschainkeeper.NewKeeper(
			cdc, moduleName, storeKey, paramSpace,
			stakingKeeper, distributionKeeper, bankKeeper,
			ibcTransferKeeper, channelKeeper, erc20Keeper,
		),
	}
}
