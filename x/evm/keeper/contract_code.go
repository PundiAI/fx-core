package keeper

import (
	"bytes"
	"encoding/hex"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/x/evm/statedb"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/x/evm/types"
)

// CreateContractWithCode create contract account and set code
func (k *Keeper) CreateContractWithCode(ctx sdk.Context, address common.Address, code []byte) error {
	codeHash := crypto.Keccak256Hash(code)
	k.Logger(ctx).Debug("create contract with code", "address", address.String(), "code-hash", codeHash)

	acc := k.GetAccount(ctx, address)
	if acc == nil {
		acc = statedb.NewEmptyAccount()
	}
	acc.CodeHash = codeHash.Bytes()
	k.SetCode(ctx, acc.CodeHash, code)
	if err := k.SetAccount(ctx, address, *acc); err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventContractCode,
		sdk.NewAttribute(types.AttributeKeyContract, address.String()),
		sdk.NewAttribute(types.AttributeKeyCodeHash, hex.EncodeToString(acc.CodeHash)),
	))
	return nil
}

// UpdateContractCode update contract code and code-hash
func (k *Keeper) UpdateContractCode(ctx sdk.Context, address common.Address, contractCode []byte) error {
	acc := k.GetAccount(ctx, address)
	if acc == nil {
		return evmtypes.ErrInvalidAccount.Wrap(address.String())
	}
	codeHash := crypto.Keccak256Hash(contractCode).Bytes()
	if bytes.Equal(codeHash, acc.CodeHash) {
		return sdkerrors.ErrInvalidRequest.Wrapf("update the same code hash: %s", address.String())
	}

	acc.CodeHash = codeHash
	k.SetCode(ctx, acc.CodeHash, contractCode)
	if err := k.SetAccount(ctx, address, *acc); err != nil {
		return err
	}

	k.Logger(ctx).Info("update contract code", "address", address.String(), "code-hash", hex.EncodeToString(acc.CodeHash))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventContractCode,
		sdk.NewAttribute(types.AttributeKeyContract, address.String()),
		sdk.NewAttribute(types.AttributeKeyCodeHash, hex.EncodeToString(acc.CodeHash)),
	))
	return nil
}

// DeployContract deploy contract with args
func (k *Keeper) DeployContract(ctx sdk.Context, from common.Address, abi abi.ABI, bin []byte, args ...interface{}) (common.Address, error) {
	data, err := abi.Pack("", args...)
	if err != nil {
		return common.Address{}, types.ErrABIPack.Wrap(err.Error())
	}
	data = append(bin, data...)

	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	_, err = k.CallEVMWithoutGas(ctx, from, nil, nil, data, true)
	if err != nil {
		return common.Address{}, err
	}
	contractAddr := crypto.CreateAddress(from, nonce)
	return contractAddr, nil
}

// DeployUpgradableContract deploy upgrade contract and initialize it
func (k *Keeper) DeployUpgradableContract(ctx sdk.Context, from, logic common.Address, logicData []byte, abi *abi.ABI, args ...interface{}) (common.Address, error) {
	// deploy proxy
	erc1967Proxy := contract.GetERC1967Proxy()
	if logicData == nil {
		logicData = []byte{}
	}
	proxyContract, err := k.DeployContract(ctx, from, erc1967Proxy.ABI, erc1967Proxy.Bin, logic, logicData)
	if err != nil {
		return common.Address{}, err
	}

	// initialize contract
	if abi != nil {
		_, err = k.ApplyContract(ctx, from, proxyContract, nil, *abi, "initialize", args...)
		if err != nil {
			return common.Address{}, err
		}
	}
	return proxyContract, nil
}

// QueryContract query contract with args and res
func (k *Keeper) QueryContract(ctx sdk.Context, from, contract common.Address, abi abi.ABI, method string, res interface{}, args ...interface{}) error {
	data, err := abi.Pack(method, args...)
	if err != nil {
		return types.ErrABIPack.Wrap(err.Error())
	}
	resp, err := k.CallEVMWithoutGas(ctx, from, &contract, nil, data, false)
	if err != nil {
		return err
	}
	if err = abi.UnpackIntoInterface(res, method, resp.Ret); err != nil {
		return types.ErrABIUnpack.Wrap(err.Error())
	}
	return nil
}

// ApplyContract apply contract with args
func (k *Keeper) ApplyContract(ctx sdk.Context, from, contract common.Address, value *big.Int, abi abi.ABI, method string, constructorData ...interface{}) (*evmtypes.MsgEthereumTxResponse, error) {
	args, err := abi.Pack(method, constructorData...)
	if err != nil {
		return nil, types.ErrABIPack.Wrap(err.Error())
	}
	resp, err := k.CallEVMWithoutGas(ctx, from, &contract, value, args, true)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
