package keeper

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/ethermint/server/config"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	"github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

type Keeper struct {
	*evmkeeper.Keeper

	// access to account state
	accountKeeper types.AccountKeeper

	// has evm hooks
	hasHooks bool
}

func NewKeeper(ek *evmkeeper.Keeper, ak types.AccountKeeper) *Keeper {
	ctx := sdk.Context{}.WithChainID(fxtypes.ChainIdWithEIP155())
	ek.WithChainID(ctx)
	return &Keeper{
		Keeper:        ek,
		accountKeeper: ak,
	}
}

// SetHooks sets the hooks for the EVM module
// It should be called only once during initialization, it panic if called more than once.
func (k *Keeper) SetHooks(eh types.EvmHooks) *Keeper {
	k.Keeper.SetHooks(eh)
	k.hasHooks = true
	return k
}

// CallEVMWithoutGas performs a smart contract method call using contract data without gas
func (k *Keeper) CallEVMWithoutGas(
	ctx sdk.Context,
	from common.Address,
	contract *common.Address,
	data []byte,
	commit bool,
) (*types.MsgEthereumTxResponse, error) {
	gasMeter := ctx.GasMeter()
	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return nil, err
	}

	gasLimit := config.DefaultGasCap
	params := ctx.ConsensusParams()
	if params != nil && params.Block != nil && params.Block.MaxGas > 0 {
		gasLimit = uint64(params.Block.MaxGas)
	}

	msg := ethtypes.NewMessage(
		from,
		contract,
		nonce,
		big.NewInt(0), // amount
		gasLimit,      // gasLimit
		big.NewInt(0), // gasFeeCap
		big.NewInt(0), // gasTipCap
		big.NewInt(0), // gasPrice
		data,
		ethtypes.AccessList{}, // AccessList
		!commit,               // isFake
	)

	res, err := k.ApplyMessage(ctx, msg, types.NewNoOpTracer(), commit)
	if err != nil {
		return nil, err
	}

	if res.Failed() {
		return nil, errorsmod.Wrap(types.ErrVMExecution, res.VmError)
	}

	ctx.WithGasMeter(gasMeter)

	return res, nil
}
