package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type EvmKeeper interface {
	ExecuteEVM(ctx sdk.Context, from common.Address, contract *common.Address, value *big.Int, gasLimit uint64, data []byte) (*evmtypes.MsgEthereumTxResponse, error)
}

type CrosschainKeeper interface {
	IBCCoinToEvm(ctx sdk.Context, holder sdk.AccAddress, ibcCoin sdk.Coin) error
	IBCCoinRefund(ctx sdk.Context, holder sdk.AccAddress, ibcCoin sdk.Coin, ibcChannel string, ibcSequence uint64) error
	AfterIBCAckSuccess(ctx sdk.Context, sourceChannel string, sequence uint64) error
}
