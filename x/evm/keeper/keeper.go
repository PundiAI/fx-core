package keeper

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/functionx/fx-core/x/evm/statedb"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ethereum/go-ethereum/params"
	ethermint "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/evm/types"
)

// Keeper grants access to the EVM module state and implements the go-ethereum StateDB interface.
type Keeper struct {
	// Protobuf codec
	cdc codec.BinaryCodec
	// Store key required for the EVM Prefix KVStore. It is required by:
	// - storing account's Storage State
	// - storing account's Code
	// - storing transaction Logs
	// - storing Bloom filters by block height. Needed for the Web3 API.
	storeKey sdk.StoreKey

	// key to access the transient store, which is reset on every block during Commit
	transientKey sdk.StoreKey

	// module specific parameter space that can be configured through governance
	paramSpace paramtypes.Subspace
	// access to account state
	accountKeeper types.AccountKeeper
	// update balance and accounting operations with coins
	bankKeeper types.BankKeeper
	// access historical headers for EVM state transition execution
	stakingKeeper types.StakingKeeper
	// fetch EIP1559 base fee and parameters
	feeMarketKeeper types.FeeMarketKeeper

	// chain ID number obtained from the context's chain id
	eip155ChainID *big.Int

	// Tracer used to collect execution traces from the EVM transaction execution
	tracer string

	// EVM Hooks for tx post-processing
	hooks types.EvmHooks
}

// NewKeeper generates new evm module keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey, transientKey sdk.StoreKey, paramSpace paramtypes.Subspace,
	ak types.AccountKeeper, bankKeeper types.BankKeeper, sk types.StakingKeeper,
	fmk types.FeeMarketKeeper,
	tracer string,
) *Keeper {
	// ensure evm module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the EVM module account has not been set")
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	// NOTE: we pass in the parameter space to the CommitStateDB in order to use custom denominations for the EVM operations
	return &Keeper{
		cdc:             cdc,
		paramSpace:      paramSpace,
		accountKeeper:   ak,
		bankKeeper:      bankKeeper,
		stakingKeeper:   sk,
		feeMarketKeeper: fmk,
		storeKey:        storeKey,
		transientKey:    transientKey,
		tracer:          tracer,
		eip155ChainID:   ethermint.EIP155ChainID(),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// ChainID returns the EIP155 chain ID for the EVM context
func (k Keeper) ChainID() *big.Int {
	return k.eip155ChainID
}

// ----------------------------------------------------------------------------
// Block Bloom
// Required by Web3 API.
// ----------------------------------------------------------------------------

// EmitBlockBloomEvent emit block bloom events
func (k Keeper) EmitBlockBloomEvent(ctx sdk.Context, bloom ethtypes.Bloom) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBlockBloom,
			sdk.NewAttribute(types.AttributeKeyEthereumBloom, string(bloom.Bytes())),
		),
	)
}

// GetBlockBloomTransient returns bloom bytes for the current block height
func (k Keeper) GetBlockBloomTransient(ctx sdk.Context) *big.Int {
	store := prefix.NewStore(ctx.TransientStore(k.transientKey), types.KeyPrefixTransientBloom)
	heightBz := sdk.Uint64ToBigEndian(uint64(ctx.BlockHeight()))
	bz := store.Get(heightBz)
	if len(bz) == 0 {
		return big.NewInt(0)
	}

	return new(big.Int).SetBytes(bz)
}

// SetBlockBloomTransient sets the given bloom bytes to the transient store. This value is reset on
// every block.
func (k Keeper) SetBlockBloomTransient(ctx sdk.Context, bloom *big.Int) {
	store := prefix.NewStore(ctx.TransientStore(k.transientKey), types.KeyPrefixTransientBloom)
	heightBz := sdk.Uint64ToBigEndian(uint64(ctx.BlockHeight()))
	store.Set(heightBz, bloom.Bytes())
}

// ----------------------------------------------------------------------------
// Tx
// ----------------------------------------------------------------------------

// SetTxIndexTransient set the index of processing transaction
func (k Keeper) SetTxIndexTransient(ctx sdk.Context, index uint64) {
	store := ctx.TransientStore(k.transientKey)
	store.Set(types.KeyPrefixTransientTxIndex, sdk.Uint64ToBigEndian(index))
}

// GetTxIndexTransient returns EVM transaction index on the current block.
func (k Keeper) GetTxIndexTransient(ctx sdk.Context) uint64 {
	store := ctx.TransientStore(k.transientKey)
	bz := store.Get(types.KeyPrefixTransientTxIndex)
	if len(bz) == 0 {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

// ----------------------------------------------------------------------------
// Log
// ----------------------------------------------------------------------------

// GetLogSizeTransient returns EVM log index on the current block.
func (k Keeper) GetLogSizeTransient(ctx sdk.Context) uint64 {
	store := ctx.TransientStore(k.transientKey)
	bz := store.Get(types.KeyPrefixTransientLogSize)
	if len(bz) == 0 {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

// SetLogSizeTransient fetches the current EVM log index from the transient store, increases its
// value by one and then sets the new index back to the transient store.
func (k Keeper) SetLogSizeTransient(ctx sdk.Context, logSize uint64) {
	store := ctx.TransientStore(k.transientKey)
	store.Set(types.KeyPrefixTransientLogSize, sdk.Uint64ToBigEndian(logSize))
}

// ----------------------------------------------------------------------------
// Storage
// ----------------------------------------------------------------------------

// GetAccountStorage return state storage associated with an account
func (k Keeper) GetAccountStorage(ctx sdk.Context, address common.Address) types.Storage {
	storage := types.Storage{}

	k.ForEachStorage(ctx, address, func(key, value common.Hash) bool {
		storage = append(storage, types.NewState(key, value))
		return true
	})

	return storage
}

// ----------------------------------------------------------------------------
// Account
// ----------------------------------------------------------------------------

// SetHooks sets the hooks for the EVM module
// It should be called only once during initialization, it panic if called more than once.
func (k *Keeper) SetHooks(eh types.EvmHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set evm hooks twice")
	}

	k.hooks = eh
	return k
}

// PostTxProcessing delegate the call to the hooks. If no hook has been registered, this function returns with a `nil` error
func (k *Keeper) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	if k.hooks == nil {
		return nil
	}
	return k.hooks.PostTxProcessing(ctx, msg, receipt)
}

// Tracer return a default vm.Tracer based on current keeper state
func (k Keeper) Tracer(ctx sdk.Context, msg core.Message, ethCfg *params.ChainConfig) vm.EVMLogger {
	return types.NewTracer(k.tracer, msg, ethCfg, ctx.BlockHeight())
}

// GetAccountWithoutBalance load nonce and codehash without balance,
// more efficient in cases where balance is not needed.
func (k *Keeper) GetAccountWithoutBalance(ctx sdk.Context, addr common.Address) *statedb.Account {
	accAddr := sdk.AccAddress(addr.Bytes())
	acct := k.accountKeeper.GetAccount(ctx, accAddr)
	if acct == nil {
		return nil
	}
	codeHash := types.EmptyCodeHash
	if bz, found := k.GetAddressCodeHash(ctx, common.BytesToAddress(acct.GetAddress())); found {
		codeHash = bz
	}
	return &statedb.Account{
		Nonce:    acct.GetSequence(),
		CodeHash: codeHash,
	}
}

// GetAccountOrEmpty returns empty account if not exist, returns error if it's not `EthAccount`
func (k *Keeper) GetAccountOrEmpty(ctx sdk.Context, addr common.Address) statedb.Account {
	acct := k.GetAccount(ctx, addr)
	if acct != nil {
		return *acct
	}

	// empty account
	return statedb.Account{
		Balance:  new(big.Int),
		CodeHash: types.EmptyCodeHash,
	}
}

// GetNonce returns the sequence number of an account, returns 0 if not exists.
func (k *Keeper) GetNonce(ctx sdk.Context, addr common.Address) uint64 {
	accAddr := sdk.AccAddress(addr.Bytes())
	acct := k.accountKeeper.GetAccount(ctx, accAddr)
	if acct == nil {
		return 0
	}

	return acct.GetSequence()
}

// GetBalance load account's balance of gas token
func (k *Keeper) GetBalance(ctx sdk.Context, addr common.Address) *big.Int {
	accAddr := sdk.AccAddress(addr.Bytes())
	params := k.GetParams(ctx)
	coin := k.bankKeeper.GetBalance(ctx, accAddr, params.EvmDenom)
	return coin.Amount.BigInt()
}

// BaseFee returns current base fee, return values:
// - `nil`: london hardfork not enabled.
// - `0`: london hardfork enabled but feemarket is not enabled.
// - `n`: both london hardfork and feemarket are enabled.
func (k Keeper) BaseFee(ctx sdk.Context, ethCfg *params.ChainConfig) *big.Int {
	if !types.IsLondon(ethCfg, ctx.BlockHeight()) {
		return nil
	}
	baseFee := k.feeMarketKeeper.GetBaseFee(ctx)
	if baseFee == nil {
		// return 0 if feemarket not enabled.
		baseFee = big.NewInt(0)
	}
	return baseFee
}

// ResetTransientGasUsed reset gas used to prepare for execution of current fx tx, called in ante handler.
func (k Keeper) ResetTransientGasUsed(ctx sdk.Context) {
	store := ctx.TransientStore(k.transientKey)
	store.Delete(types.KeyPrefixTransientGasUsed)
}

// GetTransientGasUsed returns the gas used by current fx tx.
func (k Keeper) GetTransientGasUsed(ctx sdk.Context) uint64 {
	store := ctx.TransientStore(k.transientKey)
	bz := store.Get(types.KeyPrefixTransientGasUsed)
	if len(bz) == 0 {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

// SetTransientGasUsed sets the gas used by current fx tx.
func (k Keeper) SetTransientGasUsed(ctx sdk.Context, gasUsed uint64) {
	store := ctx.TransientStore(k.transientKey)
	bz := sdk.Uint64ToBigEndian(gasUsed)
	store.Set(types.KeyPrefixTransientGasUsed, bz)
}

// AddTransientGasUsed accumulate gas used by each eth msgs included in current fx tx.
func (k Keeper) AddTransientGasUsed(ctx sdk.Context, gasUsed uint64) (uint64, error) {
	result := k.GetTransientGasUsed(ctx) + gasUsed
	if result < gasUsed {
		return 0, sdkerrors.Wrap(types.ErrGasOverflow, "transient gas used")
	}
	k.SetTransientGasUsed(ctx, result)
	return result, nil
}

func (k Keeper) HasInit(ctx sdk.Context) bool {
	return k.paramSpace.Has(ctx, types.ParamStoreKeyEVMDenom)
}

func (k Keeper) GetAddressCodeHash(ctx sdk.Context, addr common.Address) ([]byte, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAddressCode)
	data := store.Get(addr.Bytes())
	return data, len(data) != 0
}

func (k Keeper) SetAddressCode(ctx sdk.Context, addr common.Address, codeHash []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAddressCode)
	store.Set(addr.Bytes(), codeHash)
}

func (k Keeper) CreateContractWithCode(ctx sdk.Context, addr common.Address, code []byte) error {
	codeHash := crypto.Keccak256(code)
	acc := k.GetAccount(ctx, addr)
	if acc == nil {
		k.Logger(ctx).Info("create contract with code", "address", addr.String(), "code", hex.EncodeToString(code))
		acc = statedb.NewEmptyAccount()
		acc.CodeHash = codeHash
		k.SetCode(ctx, acc.CodeHash, code)
		return k.SetAccount(ctx, addr, *acc)
	}
	if !acc.IsContract() {
		return sdkerrors.Wrapf(types.ErrInvalidAccount, "address %s not contract, can not update code", addr.Hex())
	}
	k.Logger(ctx).Info("create contract with code", "address", addr.String(), "code", hex.EncodeToString(code))
	k.SetCode(ctx, acc.CodeHash, nil)
	acc.CodeHash = codeHash
	k.SetCode(ctx, acc.CodeHash, code)
	return k.SetAccount(ctx, addr, *acc)
}
