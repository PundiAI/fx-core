package types

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

type EvmKeeper interface {
	ExecuteEVM(ctx sdk.Context, from common.Address, contract *common.Address, value *big.Int, gasLimit uint64, data []byte) (*evmtypes.MsgEthereumTxResponse, error)
}

type CrosschainKeeper interface {
	IBCCoinToEvm(ctx sdk.Context, holder sdk.AccAddress, ibcCoin sdk.Coin) error
	IBCCoinToBaseCoin(ctx context.Context, holder sdk.AccAddress, ibcCoin sdk.Coin) (bool, string, error)
	IBCCoinRefund(ctx sdk.Context, holder sdk.AccAddress, ibcCoin sdk.Coin, ibcChannel string, ibcSequence uint64) error
	AfterIBCAckSuccess(ctx sdk.Context, sourceChannel string, sequence uint64) error
}

type CrosschainRouterMsgServer interface {
	SendToExternal(ctx context.Context, msg *types.MsgSendToExternal) (*types.MsgSendToExternalResponse, error)
}

type AccountKeeper interface {
	HasAccount(ctx context.Context, addr sdk.AccAddress) bool
	SetAccount(ctx context.Context, acc sdk.AccountI)
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}
