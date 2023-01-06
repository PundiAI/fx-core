package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

// QueryERC20 returns the data of a deployed ERC20 contract
func (k Keeper) QueryERC20(ctx sdk.Context, contract common.Address) (types.ERC20Data, error) {
	erc20 := fxtypes.GetERC20().ABI

	// Name
	res, err := k.CallEVM(ctx, erc20, k.moduleAddress, contract, false, "name")
	if err != nil {
		return types.ERC20Data{}, err
	}
	var nameRes struct{ Value string }
	if err := erc20.UnpackIntoInterface(&nameRes, "name", res.Ret); err != nil {
		return types.ERC20Data{}, sdkerrors.Wrapf(types.ErrABIUnpack, "failed to unpack name: %s", err.Error())
	}

	// Symbol
	res, err = k.CallEVM(ctx, erc20, k.moduleAddress, contract, false, "symbol")
	if err != nil {
		return types.ERC20Data{}, err
	}
	var symbolRes struct{ Value string }
	if err := erc20.UnpackIntoInterface(&symbolRes, "symbol", res.Ret); err != nil {
		return types.ERC20Data{}, sdkerrors.Wrapf(types.ErrABIUnpack, "failed to unpack symbol: %s", err.Error())
	}

	// Decimals
	res, err = k.CallEVM(ctx, erc20, k.moduleAddress, contract, false, "decimals")
	if err != nil {
		return types.ERC20Data{}, err
	}
	var decimalRes struct{ Value uint8 }
	if err := erc20.UnpackIntoInterface(&decimalRes, "decimals", res.Ret); err != nil {
		return types.ERC20Data{}, sdkerrors.Wrapf(types.ErrABIUnpack, "failed to unpack decimals: %s", err.Error())
	}

	return types.NewERC20Data(nameRes.Value, symbolRes.Value, decimalRes.Value), nil
}

// BalanceOf returns the balance of an address for ERC20 contract
func (k Keeper) BalanceOf(ctx sdk.Context, contract, addr common.Address) (*big.Int, error) {
	erc20 := fxtypes.GetERC20().ABI

	res, err := k.CallEVM(ctx, erc20, k.moduleAddress, contract, false, "balanceOf", addr)
	if err != nil {
		return nil, err
	}

	var balanceRes struct{ Value *big.Int }
	if err := erc20.UnpackIntoInterface(&balanceRes, "balanceOf", res.Ret); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrABIUnpack, "failed to unpack balanceOf: %s", err.Error())
	}
	return balanceRes.Value, nil
}

func (k Keeper) DeployUpgradableToken(ctx sdk.Context, from common.Address, name, symbol string, decimals uint8) (common.Address, error) {
	var tokenContract fxtypes.Contract
	if symbol == fxtypes.DefaultDenom {
		tokenContract = fxtypes.GetWFX()
		name = fmt.Sprintf("Wrapped %s", name)
		symbol = fmt.Sprintf("W%s", symbol)
	} else {
		tokenContract = fxtypes.GetERC20()
	}
	k.Logger(ctx).Info("deploy token", "name", name, "symbol", symbol, "decimals", decimals)

	// deploy proxy
	erc1967Proxy := fxtypes.GetERC1967Proxy()
	contract, err := k.DeployContract(ctx, from, erc1967Proxy.ABI, erc1967Proxy.Bin, tokenContract.Address, []byte{})
	if err != nil {
		return common.Address{}, err
	}

	_, err = k.CallEVM(ctx, tokenContract.ABI, from, contract, true, "initialize", name, symbol, decimals, k.moduleAddress)
	if err != nil {
		return common.Address{}, err
	}
	return contract, nil
}

func (k Keeper) DeployContract(ctx sdk.Context, from common.Address, abi abi.ABI, bin []byte, constructorData ...interface{}) (common.Address, error) {
	args, err := abi.Pack("", constructorData...)
	if err != nil {
		return common.Address{}, sdkerrors.Wrap(err, "pack constructor data")
	}
	data := make([]byte, len(bin)+len(args))
	copy(data[:len(bin)], bin)
	copy(data[len(bin):], args)

	nonce, err := k.accountKeeper.GetSequence(ctx, from.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	_, err = k.evmKeeper.CallEVMWithData(ctx, from, nil, data, true)
	if err != nil {
		return common.Address{}, err
	}
	contractAddr := crypto.CreateAddress(from, nonce)
	return contractAddr, nil
}

// CallEVM performs a smart contract method call using given args
func (k Keeper) CallEVM(
	ctx sdk.Context,
	abi abi.ABI,
	from, contract common.Address,
	commit bool,
	method string,
	args ...interface{},
) (*evmtypes.MsgEthereumTxResponse, error) {
	data, err := abi.Pack(method, args...)
	if err != nil {
		return nil, sdkerrors.Wrap(
			types.ErrABIPack,
			sdkerrors.Wrap(err, "failed to create transaction data").Error(),
		)
	}

	resp, err := k.evmKeeper.CallEVMWithData(ctx, from, &contract, data, commit)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "contract call failed: method '%s', contract '%s'", method, contract)
	}
	return resp, nil
}

// monitorApprovalEvent returns an error if the given transactions logs include
// an unexpected `approve` event
func (k Keeper) monitorApprovalEvent(res *evmtypes.MsgEthereumTxResponse) error {
	if res == nil || len(res.Logs) == 0 {
		return nil
	}

	logApprovalSigHash := crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))

	for _, log := range res.Logs {
		if log.Topics[0] == logApprovalSigHash.Hex() {
			return sdkerrors.Wrapf(
				types.ErrUnexpectedEvent, "approval event",
			)
		}
	}

	return nil
}
