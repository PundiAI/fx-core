package keeper

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	storetypes "cosmossdk.io/store/types"
	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	"github.com/evmos/ethermint/x/evm/types"

	fxcontract "github.com/functionx/fx-core/v8/contract"
	fxevmtypes "github.com/functionx/fx-core/v8/x/evm/types"
)

type Keeper struct {
	*evmkeeper.Keeper

	// access to account state
	accountKeeper fxevmtypes.AccountKeeper
	module        common.Address
}

func NewKeeper(ek *evmkeeper.Keeper, ak fxevmtypes.AccountKeeper) *Keeper {
	addr := ak.GetModuleAddress(types.ModuleName)
	if addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	return &Keeper{
		Keeper:        ek,
		accountKeeper: ak,
		module:        common.BytesToAddress(addr),
	}
}

// Deprecated: please use callEvm todo: remove this
func (k *Keeper) CallEVMWithoutGas(
	ctx sdk.Context,
	from common.Address,
	contract *common.Address,
	value *big.Int,
	data []byte,
	commit bool,
) (*types.MsgEthereumTxResponse, error) {
	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return nil, err
	}
	return k.callEvm(ctx, from, contract, value, nonce, data, commit)
}

func (k *Keeper) callEvm(
	ctx sdk.Context,
	from common.Address,
	contract *common.Address,
	value *big.Int,
	nonce uint64,
	data []byte,
	commit bool,
) (*types.MsgEthereumTxResponse, error) {
	gasMeter := ctx.GasMeter()
	ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())

	gasLimit := fxcontract.DefaultGasCap
	params := ctx.ConsensusParams()
	if params.Block != nil && params.Block.MaxGas > 0 {
		gasLimit = uint64(params.Block.MaxGas)
	}

	if value == nil {
		value = big.NewInt(0)
	}
	msg := &core.Message{
		From:              from,
		To:                contract,
		Nonce:             nonce,
		Value:             value,
		GasLimit:          gasLimit,
		GasPrice:          big.NewInt(0),
		GasFeeCap:         big.NewInt(0),
		GasTipCap:         big.NewInt(0),
		Data:              data,
		AccessList:        ethtypes.AccessList{},
		SkipAccountChecks: false,
	}

	res, err := k.ApplyMessage(ctx, msg, types.NewNoOpTracer(), commit)
	if err != nil {
		return nil, err
	}

	if res.Failed() {
		errStr := res.VmError
		if res.VmError == vm.ErrExecutionReverted.Error() {
			if cause, err := abi.UnpackRevert(common.CopyBytes(res.Ret)); err == nil {
				errStr = cause
			}
		}
		return res, types.ErrVMExecution.Wrap(errStr)
	}

	ctx.WithGasMeter(gasMeter)

	return res, nil
}

func (k *Keeper) ExecuteEVM(
	ctx sdk.Context,
	from common.Address,
	contract *common.Address,
	value *big.Int,
	gasLimit uint64,
	data []byte,
) (*types.MsgEthereumTxResponse, error) {
	gasMeter := ctx.GasMeter()
	ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())

	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return nil, err
	}

	if value == nil {
		value = big.NewInt(0)
	}
	msg := &core.Message{
		To:                contract,
		From:              from,
		Nonce:             nonce,
		Value:             value,
		GasLimit:          gasLimit,
		GasPrice:          big.NewInt(0),
		GasFeeCap:         big.NewInt(0),
		GasTipCap:         big.NewInt(0),
		Data:              data,
		AccessList:        ethtypes.AccessList{},
		SkipAccountChecks: false,
	}

	res, err := k.ApplyMessage(ctx, msg, types.NewNoOpTracer(), true)
	if err != nil {
		return nil, err
	}

	attrs := []sdk.Attribute{
		sdk.NewAttribute(sdk.AttributeKeyAmount, value.String()),
		// add event for ethereum transaction hash format
		sdk.NewAttribute(types.AttributeKeyEthereumTxHash, res.Hash),
		// add event for index of valid ethereum tx
		sdk.NewAttribute(types.AttributeKeyTxIndex, strconv.FormatUint(0, 10)),
		// add event for eth tx gas used, we can't get it from cosmos tx result when it contains multiple eth tx msgs.
		sdk.NewAttribute(types.AttributeKeyTxGasUsed, strconv.FormatUint(res.GasUsed, 10)),
	}

	if len(ctx.TxBytes()) > 0 {
		// add event for tendermint transaction hash format
		hash := tmbytes.HexBytes(tmtypes.Tx(ctx.TxBytes()).Hash())
		attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyTxHash, hash.String()))
	}

	attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyRecipient, contract.Hex()))

	if res.Failed() {
		attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyEthereumTxFailed, res.VmError))
	}

	txLogAttrs := make([]sdk.Attribute, len(res.Logs))
	for i, log := range res.Logs {
		logStr, e := json.Marshal(log)
		if e != nil {
			return nil, sdkerrors.ErrJSONMarshal.Wrapf("failed to marshal tx log: %s", e.Error())
		}
		txLogAttrs[i] = sdk.NewAttribute(types.AttributeKeyTxLog, string(logStr))
	}

	// emit events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEthereumTx,
			attrs...,
		),
		sdk.NewEvent(
			types.EventTypeTxLog,
			txLogAttrs...,
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, from.String()),
		),
	})

	ctx.WithGasMeter(gasMeter)

	return res, nil
}

func (k *Keeper) IsContract(ctx sdk.Context, account common.Address) bool {
	acc := k.GetAccount(ctx, account)
	return acc != nil && acc.IsContract()
}
