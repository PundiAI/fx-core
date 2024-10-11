package precompile

import (
	"context"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v8/types"
	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
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
	SetOutgoingTransferRelation(ctx sdk.Context, moduleName string, txID uint64)
}

type BankKeeper interface {
	GetSupply(ctx context.Context, denom string) sdk.Coin
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	MintCoins(ctx context.Context, moduleName string, amounts sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
}

type IBCTransferKeeper interface {
	Transfer(goCtx context.Context, msg *ibctransfertypes.MsgTransfer) (*ibctransfertypes.MsgTransferResponse, error)
}

type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type CrosschainKeeper interface {
	AddOutgoingBridgeCall(ctx sdk.Context, sender, refundAddr common.Address, baseCoins sdk.Coins, to common.Address, data, memo []byte, eventNonce uint64) (uint64, error)
	EvmToBaseCoin(ctx context.Context, tokenAddr string, amount *big.Int, holder common.Address) (sdk.Coin, error)
	ExecuteClaim(ctx sdk.Context, eventNonce uint64) error

	HasOracleAddrByExternalAddr(ctx sdk.Context, externalAddress string) bool
	GetOracleAddrByExternalAddr(ctx sdk.Context, externalAddress string) (sdk.AccAddress, bool)
	GetOracle(ctx sdk.Context, oracleAddr sdk.AccAddress) (oracle crosschaintypes.Oracle, found bool)

	BaseCoinToIBCCoin(ctx context.Context, coin sdk.Coin, holder sdk.AccAddress, ibcTarget string) (sdk.Coin, error)
	BuildOutgoingTxBatch(ctx sdk.Context, sender sdk.AccAddress, receiver string, amount sdk.Coin, fee sdk.Coin) (uint64, error)
}

type GovKeeper interface {
	CheckDisabledPrecompiles(ctx sdk.Context, contractAddress common.Address, methodId []byte) error
}
