package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/functionx/fx-core/server/config"
	evmtypes "github.com/functionx/fx-core/x/evm/types"

	"github.com/functionx/fx-core/x/intrarelayer/types"
	"github.com/functionx/fx-core/x/intrarelayer/types/contracts"
)

// QueryFIP20 returns the data of a deployed FIP20 contract
func (k Keeper) QueryFIP20(ctx sdk.Context, contract common.Address) (types.FIP20Data, error) {
	var (
		nameRes    types.FIP20StringResponse
		symbolRes  types.FIP20StringResponse
		decimalRes types.FIP20Uint8Response
	)

	fip20 := contracts.FIP20Contract.ABI

	// Name
	res, err := k.CallEVMWithModule(ctx, fip20, contract, "name")
	if err != nil {
		return types.FIP20Data{}, err
	}

	if err := fip20.UnpackIntoInterface(&nameRes, "name", res.Ret); err != nil {
		return types.FIP20Data{}, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to unpack name: %s", err.Error())
	}

	// Symbol
	res, err = k.CallEVMWithModule(ctx, fip20, contract, "symbol")
	if err != nil {
		return types.FIP20Data{}, err
	}

	if err := fip20.UnpackIntoInterface(&symbolRes, "symbol", res.Ret); err != nil {
		return types.FIP20Data{}, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to unpack symbol: %s", err.Error())
	}

	// Decimals
	res, err = k.CallEVMWithModule(ctx, fip20, contract, "decimals")
	if err != nil {
		return types.FIP20Data{}, err
	}

	if err := fip20.UnpackIntoInterface(&decimalRes, "decimals", res.Ret); err != nil {
		return types.FIP20Data{}, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to unpack decimals: %s", err.Error())
	}

	return types.NewFIP20Data(nameRes.Value, symbolRes.Value, decimalRes.Value), nil
}

func (k Keeper) QueryFIP20BalanceOf(ctx sdk.Context, contract, addr common.Address) (*big.Int, error) {
	fip20 := contracts.FIP20Contract.ABI
	res, err := k.CallEVMWithModule(ctx, fip20, contract, "balanceOf", addr)
	if err != nil {
		return nil, err
	}

	var balanceRes types.FIP20Uint256Response
	if err := fip20.UnpackIntoInterface(&balanceRes, "balanceOf", res.Ret); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to unpack balanceOf: %s", err.Error())
	}
	return balanceRes.Value, nil
}

// CallEVM performs a smart contract method call using  given args
func (k Keeper) CallEVM(ctx sdk.Context, abi abi.ABI, from, contract common.Address, method string, args ...interface{}) (*evmtypes.MsgEthereumTxResponse, error) {
	payload, err := abi.Pack(method, args...)
	if err != nil {
		return nil, sdkerrors.Wrap(
			types.ErrWritingEthTxPayload,
			sdkerrors.Wrap(err, "failed to create transaction payload").Error(),
		)
	}

	resp, err := k.CallEVMWithPayload(ctx, from, &contract, payload)
	if err != nil {
		return nil, fmt.Errorf("contract call failed: method '%s' %s, %s", method, contract, err)
	}
	return resp, nil
}

// CallEVMWithPayload performs a smart contract method call using contract data
func (k Keeper) CallEVMWithPayload(ctx sdk.Context, from common.Address, contract *common.Address, transferData []byte) (*evmtypes.MsgEthereumTxResponse, error) {
	k.evmKeeper.WithContext(ctx)

	nonce := k.evmKeeper.GetNonce(from)

	msg := ethtypes.NewMessage(
		from,
		contract,
		nonce,
		big.NewInt(0),        // amount
		config.DefaultGasCap, // gasLimit
		big.NewInt(0),        // gasFeeCap
		big.NewInt(0),        // gasTipCap
		big.NewInt(0),        // gasPrice
		transferData,
		ethtypes.AccessList{}, // AccessList
		true,                  // checkNonce
	)

	res, err := k.evmKeeper.ApplyMessage(msg, evmtypes.NewNoOpTracer(), true)
	if err != nil {
		return nil, err
	}

	k.evmKeeper.SetNonce(from, nonce+1)

	if res.Failed() {
		return nil, sdkerrors.Wrap(evmtypes.ErrVMExecution, res.VmError)
	}

	return res, nil
}

// CallEVMWithModule performs a smart contract method call using  given args with module
func (k Keeper) CallEVMWithModule(ctx sdk.Context, abi abi.ABI, contract common.Address, method string, args ...interface{}) (*evmtypes.MsgEthereumTxResponse, error) {
	payload, err := abi.Pack(method, args...)
	if err != nil {
		return nil, sdkerrors.Wrap(
			types.ErrWritingEthTxPayload,
			sdkerrors.Wrap(err, "failed to create transaction payload").Error(),
		)
	}

	resp, err := k.CallEVMWithPayloadWithModule(ctx, &contract, payload)
	if err != nil {
		return nil, fmt.Errorf("contract call failed: method '%s' %s, %s", method, contract, err)
	}
	return resp, nil
}

// CallEVMWithPayloadWithModule performs a smart contract method call using contract data with module
func (k Keeper) CallEVMWithPayloadWithModule(ctx sdk.Context, contract *common.Address, transferData []byte) (*evmtypes.MsgEthereumTxResponse, error) {
	k.evmKeeper.WithContext(ctx)

	nonce, err := k.accountKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return nil, err
	}

	msg := ethtypes.NewMessage(
		types.ModuleAddress,
		contract,
		nonce,
		big.NewInt(0),        // amount
		config.DefaultGasCap, // gasLimit
		big.NewInt(0),        // gasFeeCap
		big.NewInt(0),        // gasTipCap
		big.NewInt(0),        // gasPrice
		transferData,
		ethtypes.AccessList{}, // AccessList
		true,                  // checkNonce
	)

	res, err := k.evmKeeper.ApplyMessage(msg, evmtypes.NewNoOpTracer(), true)
	if err != nil {
		return nil, err
	}

	k.evmKeeper.SetNonce(types.ModuleAddress, nonce+1)

	if res.Failed() {
		return nil, sdkerrors.Wrap(evmtypes.ErrVMExecution, res.VmError)
	}

	return res, nil
}
