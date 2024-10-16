package precompile

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	crosschaintypes "github.com/functionx/fx-core/v8/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
)

type BankKeeper interface {
	MintCoins(ctx context.Context, moduleName string, amounts sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type CrosschainKeeper interface {
	ExecuteClaim(ctx sdk.Context, eventNonce uint64) error
	BridgeCoinSupply(ctx context.Context, token, target string) (sdk.Coin, error)
	CrosschainBaseCoin(ctx sdk.Context, from sdk.AccAddress, receipt string, amount, fee sdk.Coin, fxTarget *crosschaintypes.FxTarget, memo string, originToken bool) error
	BridgeCallBaseCoin(ctx sdk.Context, from, refund, to common.Address, coins sdk.Coins, data, memo []byte, fxTarget *crosschaintypes.FxTarget, originTokenAmount sdkmath.Int) (uint64, error)
	GetBaseDenomByErc20(ctx sdk.Context, erc20Addr common.Address) (erc20types.ERC20Token, error)

	HasOracleAddrByExternalAddr(ctx sdk.Context, externalAddress string) bool
	GetOracleAddrByExternalAddr(ctx sdk.Context, externalAddress string) (sdk.AccAddress, bool)
	GetOracle(ctx sdk.Context, oracleAddr sdk.AccAddress) (oracle crosschaintypes.Oracle, found bool)
}

type GovKeeper interface {
	CheckDisabledPrecompiles(ctx sdk.Context, contractAddress common.Address, methodId []byte) error
}
