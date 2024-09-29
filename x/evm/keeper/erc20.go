package keeper

import (
	"context"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v8/contract"
	"github.com/functionx/fx-core/v8/x/evm/types"
)

func (k *Keeper) ERE20Name(ctx context.Context, contractAddr common.Address) (string, error) {
	var nameRes struct{ Value string }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.module, contractAddr, contract.GetFIP20().ABI, "name", &nameRes); err != nil {
		return "", err
	}
	return nameRes.Value, nil
}

func (k *Keeper) ERE20Symbol(ctx context.Context, contractAddr common.Address) (string, error) {
	var symbolRes struct{ Value string }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.module, contractAddr, contract.GetFIP20().ABI, "symbol", &symbolRes); err != nil {
		return "", err
	}
	return symbolRes.Value, nil
}

func (k *Keeper) ERE20Decimals(ctx context.Context, contractAddr common.Address) (uint8, error) {
	var decimalRes struct{ Value uint8 }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.module, contractAddr, contract.GetFIP20().ABI, "decimals", &decimalRes); err != nil {
		return 0, nil
	}
	return decimalRes.Value, nil
}

func (k *Keeper) ERC20BalanceOf(ctx context.Context, contractAddr, addr common.Address) (*big.Int, error) {
	var balanceRes struct {
		Value *big.Int
	}
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.module, contractAddr, contract.GetFIP20().ABI, "balanceOf", &balanceRes, addr); err != nil {
		return big.NewInt(0), err
	}
	return balanceRes.Value, nil
}

func (k *Keeper) ERE20Mint(ctx context.Context, contractAddr common.Address, from, receiver common.Address, amount *big.Int) error {
	_, err := k.ApplyContract(sdk.UnwrapSDKContext(ctx), from, contractAddr, nil, contract.GetFIP20().ABI, "mint", receiver, amount)
	return err
}

func (k *Keeper) ERE20Burn(ctx context.Context, contractAddr common.Address, from, account common.Address, amount *big.Int) error {
	_, err := k.ApplyContract(sdk.UnwrapSDKContext(ctx), from, contractAddr, nil, contract.GetFIP20().ABI, "burn", account, amount)
	return err
}

func (k *Keeper) ERE20Transfer(ctx context.Context, contractAddr common.Address, from, receiver common.Address, amount *big.Int) error {
	erc20ABI := contract.GetFIP20().ABI
	res, err := k.ApplyContract(sdk.UnwrapSDKContext(ctx), from, contractAddr, nil, erc20ABI, "transfer", receiver, amount)
	if err != nil {
		return err
	}

	// Check unpackedRet execution
	var unpackedRet struct{ Value bool }
	if err = erc20ABI.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
		return errorsmod.Wrapf(types.ErrABIUnpack, "failed to unpack transfer: %s", err.Error())
	}
	if !unpackedRet.Value {
		return errorsmod.Wrap(errortypes.ErrLogic, "failed to execute transfer")
	}
	return nil
}
