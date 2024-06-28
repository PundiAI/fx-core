package precompile

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v7/types"
	erc20types "github.com/functionx/fx-core/v7/x/erc20/types"
)

type Erc20Keeper interface {
	ModuleAddress() common.Address
	GetTokenPairByAddress(ctx sdk.Context, address common.Address) (erc20types.TokenPair, bool)
	ConvertDenomToTarget(ctx sdk.Context, from sdk.AccAddress, coin sdk.Coin, fxTarget fxtypes.FxTarget) (sdk.Coin, error)
	GetIbcTimeout(ctx sdk.Context) time.Duration
	SetIBCTransferRelation(ctx sdk.Context, channel string, sequence uint64)
	HasOutgoingTransferRelation(ctx sdk.Context, moduleName string, txID uint64) bool
	ToTargetDenom(ctx sdk.Context, denom, base string, aliases []string, fxTarget fxtypes.FxTarget) string
	GetTokenPair(ctx sdk.Context, tokenOrDenom string) (erc20types.TokenPair, bool)
	IsOriginDenom(ctx sdk.Context, denom string) bool
	HasDenomAlias(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
}

type BankKeeper interface {
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	MintCoins(ctx sdk.Context, moduleName string, amounts sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amounts sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
}

type IBCTransferKeeper interface {
	Transfer(goCtx context.Context, msg *ibctransfertypes.MsgTransfer) (*ibctransfertypes.MsgTransferResponse, error)
}

type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type CrosschainKeeper interface {
	TransferAfter(ctx sdk.Context, sender sdk.AccAddress, receive string, coins, fee sdk.Coin, originToken, insufficientLiquidity bool) error
	PrecompileCancelSendToExternal(ctx sdk.Context, txID uint64, sender sdk.AccAddress) (sdk.Coin, error)
	PrecompileIncreaseBridgeFee(ctx sdk.Context, txID uint64, sender sdk.AccAddress, addBridgeFee sdk.Coin) error
	PrecompileBridgeCall(ctx sdk.Context, sender, refund common.Address, coins sdk.Coins, to common.Address, data, memo []byte) (uint64, error)
	PrecompileCancelPendingBridgeCall(ctx sdk.Context, nonce uint64, sender sdk.AccAddress) (sdk.Coins, error)
}
