package ante

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"

	"github.com/evmos/ethermint/x/evm/statedb"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
)

// AccountKeeper defines the expected account keeper interface
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetAllAccounts(ctx sdk.Context) (accounts []authtypes.AccountI)
	IterateAccounts(ctx sdk.Context, cb func(account authtypes.AccountI) bool)
	GetSequence(sdk.Context, sdk.AccAddress) (uint64, error)
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, account authtypes.AccountI)
	RemoveAccount(ctx sdk.Context, account authtypes.AccountI)
	GetParams(ctx sdk.Context) (params authtypes.Params)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}

// FeegrantKeeper defines the expected feegrant keeper.
type FeegrantKeeper interface {
	UseGrantedFees(ctx sdk.Context, granter, grantee sdk.AccAddress, fee sdk.Coins, msgs []sdk.Msg) error
}

// EVMKeeper defines the expected keeper interface used on the Eth AnteHandler
type EVMKeeper interface {
	statedb.Keeper

	ChainID() *big.Int
	GetParams(ctx sdk.Context) evmtypes.Params
	NewEVM(ctx sdk.Context, msg core.Message, cfg *evmtypes.EVMConfig, tracer vm.EVMLogger, stateDB vm.StateDB) *vm.EVM
	DeductTxCostsFromUserBalance(
		ctx sdk.Context, msgEthTx evmtypes.MsgEthereumTx, txData evmtypes.TxData, denom string, homestead, istanbul, london bool,
	) (sdk.Coins, error)
	GetBaseFee(ctx sdk.Context, ethCfg *params.ChainConfig) *big.Int
	GetBalance(ctx sdk.Context, addr common.Address) *big.Int
	ResetTransientGasUsed(ctx sdk.Context)
	GetTxIndexTransient(ctx sdk.Context) uint64
}

type protoTxProvider interface {
	GetProtoTx() *tx.Tx
}

// FeeMarketKeeper defines the expected keeper interface used on the AnteHandler
type FeeMarketKeeper interface {
	GetParams(ctx sdk.Context) (params feemarkettypes.Params)
	AddTransientGasWanted(ctx sdk.Context, gasWanted uint64) (uint64, error)
}
