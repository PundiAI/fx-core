package types

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type EvmKeeper interface {
	CallEVM(ctx sdk.Context, from common.Address, contract *common.Address, value *big.Int, gasLimit uint64, data []byte, commit bool) (*evmtypes.MsgEthereumTxResponse, error)
}

type CrossChainKeeper interface {
	IBCCoinToEvm(ctx context.Context, coin sdk.Coin, holder sdk.AccAddress) error
	IBCCoinRefund(ctx sdk.Context, coin sdk.Coin, holder sdk.AccAddress, ibcChannel string, ibcSequence uint64) error
	AfterIBCAckSuccess(ctx sdk.Context, sourceChannel string, sequence uint64)
}
