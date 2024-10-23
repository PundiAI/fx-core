package contract

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/types"
)

type ERC20TokenKeeper struct {
	Caller
	abi  abi.ABI
	from common.Address
}

func NewERC20TokenKeeper(caller Caller) ERC20TokenKeeper {
	return ERC20TokenKeeper{
		Caller: caller,
		abi:    GetFIP20().ABI,
		// evm module address
		from: common.BytesToAddress(authtypes.NewModuleAddress(types.ModuleName).Bytes()),
	}
}

func (k ERC20TokenKeeper) Name(ctx context.Context, contractAddr common.Address) (string, error) {
	var nameRes struct{ Value string }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, contractAddr, k.abi, "name", &nameRes); err != nil {
		return "", err
	}
	return nameRes.Value, nil
}

func (k ERC20TokenKeeper) Symbol(ctx context.Context, contractAddr common.Address) (string, error) {
	var symbolRes struct{ Value string }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, contractAddr, k.abi, "symbol", &symbolRes); err != nil {
		return "", err
	}
	return symbolRes.Value, nil
}

func (k ERC20TokenKeeper) Decimals(ctx context.Context, contractAddr common.Address) (uint8, error) {
	var decimalRes struct{ Value uint8 }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, contractAddr, k.abi, "decimals", &decimalRes); err != nil {
		return 0, nil
	}
	return decimalRes.Value, nil
}

func (k ERC20TokenKeeper) BalanceOf(ctx context.Context, contractAddr, addr common.Address) (*big.Int, error) {
	var balanceRes struct {
		Value *big.Int
	}
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, contractAddr, k.abi, "balanceOf", &balanceRes, addr); err != nil {
		return big.NewInt(0), err
	}
	return balanceRes.Value, nil
}

func (k ERC20TokenKeeper) TotalSupply(ctx context.Context, contractAddr common.Address) (*big.Int, error) {
	var totalSupplyRes struct{ Value *big.Int }
	if err := k.QueryContract(sdk.UnwrapSDKContext(ctx), k.from, contractAddr, k.abi, "totalSupply", &totalSupplyRes); err != nil {
		return nil, err
	}
	return totalSupplyRes.Value, nil
}

func (k ERC20TokenKeeper) Mint(ctx context.Context, contractAddr, from, receiver common.Address, amount *big.Int) error {
	_, err := k.ApplyContract(sdk.UnwrapSDKContext(ctx), from, contractAddr, nil, k.abi, "mint", receiver, amount)
	return err
}

func (k ERC20TokenKeeper) Burn(ctx context.Context, contractAddr, from, account common.Address, amount *big.Int) error {
	_, err := k.ApplyContract(sdk.UnwrapSDKContext(ctx), from, contractAddr, nil, k.abi, "burn", account, amount)
	return err
}

func (k ERC20TokenKeeper) Transfer(ctx context.Context, contractAddr, from, receiver common.Address, amount *big.Int) error {
	res, err := k.ApplyContract(sdk.UnwrapSDKContext(ctx), from, contractAddr, nil, k.abi, "transfer", receiver, amount)
	if err != nil {
		return err
	}

	// Check unpackedRet execution
	var unpackedRet struct{ Value bool }
	if err = k.abi.UnpackIntoInterface(&unpackedRet, "transfer", res.Ret); err != nil {
		return sdkerrors.ErrInvalidType.Wrapf("failed to unpack transfer: %s", err.Error())
	}
	if !unpackedRet.Value {
		return sdkerrors.ErrLogic.Wrap("failed to execute transfer")
	}
	return nil
}
