package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/ethereum/go-ethereum/common"

	"github.com/pundiai/fx-core/v8/precompiles/staking"
	fxtypes "github.com/pundiai/fx-core/v8/types"
)

const RefundContractAddress = "0xa6084D3f37236cCCF84368a623741da50f221156"

func refundDelegate(ctx sdk.Context, bankKeeper bankkeeper.Keeper) error {
	stakingAddr := new(staking.Contract).Address()
	refundContract := common.HexToAddress(RefundContractAddress)

	bal := bankKeeper.GetBalance(ctx, stakingAddr.Bytes(), fxtypes.DefaultDenom)
	return bankKeeper.SendCoins(ctx, stakingAddr.Bytes(), refundContract.Bytes(), sdk.NewCoins(bal))
}
