package types

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
)

type RefundHook interface {
	RefundAfter(ctx sdk.Context, sourceChannel string, sequence uint64, sender sdk.AccAddress, amount sdk.Coin) error
	AckAfter(ctx sdk.Context, sourceChannel string, sequence uint64) error
}

type Erc20Keeper interface {
	ConvertDenomToTarget(ctx sdk.Context, from sdk.AccAddress, coin sdk.Coin, fxTarget fxtypes.FxTarget) (sdk.Coin, error)
	ConvertCoin(goCtx context.Context, msg *erc20types.MsgConvertCoin) (*erc20types.MsgConvertCoinResponse, error)
}

type EvmKeeper interface {
	CallEVM(ctx sdk.Context, from common.Address, contract *common.Address, value *big.Int, gasLimit uint64, data []byte, commit bool) (*evmtypes.MsgEthereumTxResponse, error)
}
