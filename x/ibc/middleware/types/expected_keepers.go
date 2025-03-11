package types

import (
	"context"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/pundiai/fx-core/v8/x/crosschain/types"
	erc20types "github.com/pundiai/fx-core/v8/x/erc20/types"
)

type EvmKeeper interface {
	ExecuteEVM(ctx sdk.Context, from common.Address, contract *common.Address, value *big.Int, gasLimit uint64, data []byte) (*evmtypes.MsgEthereumTxResponse, error)
}

type CrosschainKeeper interface {
	IBCCoinToEvm(ctx sdk.Context, holder string, ibcCoin sdk.Coin) error
	IBCCoinToBaseCoin(ctx context.Context, holder sdk.AccAddress, ibcCoin sdk.Coin) (bool, string, error)
	IBCCoinRefund(ctx sdk.Context, holder sdk.AccAddress, ibcCoin sdk.Coin, ibcChannel string, ibcSequence uint64) error
	AfterIBCAckSuccess(ctx sdk.Context, sourceChannel string, sequence uint64) error
	BridgeTokenToBaseCoin(ctx context.Context, holder sdk.AccAddress, amount sdkmath.Int, bridgeToken erc20types.BridgeToken) (sdk.Coin, error)
}

type CrosschainRouterMsgServer interface {
	SendToExternal(ctx context.Context, msg *types.MsgSendToExternal) (*types.MsgSendToExternalResponse, error)
}

type AccountKeeper interface {
	HasAccount(ctx context.Context, addr sdk.AccAddress) bool
	SetAccount(ctx context.Context, acc sdk.AccountI)
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

type Erc20Keeper interface {
	GetBaseDenom(ctx context.Context, token string) (string, error)
	GetBridgeToken(ctx context.Context, chainName, baseDenom string) (erc20types.BridgeToken, error)
}
