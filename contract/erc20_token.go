package contract

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/types"
)

type ERC20TokenKeeper struct {
	Caller
	abi abi.ABI
}

func NewERC20TokenKeeper(caller Caller) ERC20TokenKeeper {
	return ERC20TokenKeeper{
		Caller: caller,
		abi:    GetWPUNDIAI().ABI,
	}
}

func (k ERC20TokenKeeper) Owner(ctx context.Context, contractAddr common.Address) (common.Address, error) {
	var ownerRes struct{ Value common.Address }
	if err := k.QueryContract(ctx, common.Address{}, contractAddr, k.abi, "owner", &ownerRes); err != nil {
		return common.Address{}, err
	}
	return ownerRes.Value, nil
}

func (k ERC20TokenKeeper) Name(ctx context.Context, contractAddr common.Address) (string, error) {
	var nameRes struct{ Value string }
	if err := k.QueryContract(ctx, common.Address{}, contractAddr, k.abi, "name", &nameRes); err != nil {
		return "", err
	}
	return nameRes.Value, nil
}

func (k ERC20TokenKeeper) Symbol(ctx context.Context, contractAddr common.Address) (string, error) {
	var symbolRes struct{ Value string }
	if err := k.QueryContract(ctx, common.Address{}, contractAddr, k.abi, "symbol", &symbolRes); err != nil {
		return "", err
	}
	return symbolRes.Value, nil
}

func (k ERC20TokenKeeper) Decimals(ctx context.Context, contractAddr common.Address) (uint8, error) {
	var decimalRes struct{ Value uint8 }
	if err := k.QueryContract(ctx, common.Address{}, contractAddr, k.abi, "decimals", &decimalRes); err != nil {
		return 0, err
	}
	return decimalRes.Value, nil
}

func (k ERC20TokenKeeper) BalanceOf(ctx context.Context, contractAddr, addr common.Address) (*big.Int, error) {
	var balanceRes struct {
		Value *big.Int
	}
	if err := k.QueryContract(ctx, common.Address{}, contractAddr, k.abi, "balanceOf", &balanceRes, addr); err != nil {
		return big.NewInt(0), err
	}
	return balanceRes.Value, nil
}

func (k ERC20TokenKeeper) TotalSupply(ctx context.Context, contractAddr common.Address) (*big.Int, error) {
	var totalSupplyRes struct{ Value *big.Int }
	if err := k.QueryContract(ctx, common.Address{}, contractAddr, k.abi, "totalSupply", &totalSupplyRes); err != nil {
		return nil, err
	}
	return totalSupplyRes.Value, nil
}

func (k ERC20TokenKeeper) Allowance(ctx context.Context, contractAddr, owner, spender common.Address) (*big.Int, error) {
	var allowanceRes struct{ Value *big.Int }
	if err := k.QueryContract(ctx, owner, contractAddr, k.abi, "allowance", &allowanceRes, owner, spender); err != nil {
		return big.NewInt(0), err
	}
	return allowanceRes.Value, nil
}

func (k ERC20TokenKeeper) Approve(ctx context.Context, contractAddr, from, spender common.Address, amount *big.Int) (*types.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, from, contractAddr, nil, k.abi, "approve", spender, amount)
	if err != nil {
		return nil, err
	}
	return unpackRetIsOk(k.abi, "approve", res)
}

func (k ERC20TokenKeeper) Mint(ctx context.Context, contractAddr, from, receiver common.Address, amount *big.Int) (*types.MsgEthereumTxResponse, error) {
	return k.ApplyContract(ctx, from, contractAddr, nil, k.abi, "mint", receiver, amount)
}

func (k ERC20TokenKeeper) Burn(ctx context.Context, contractAddr, from, account common.Address, amount *big.Int) (*types.MsgEthereumTxResponse, error) {
	return k.ApplyContract(ctx, from, contractAddr, nil, k.abi, "burn", account, amount)
}

func (k ERC20TokenKeeper) Transfer(ctx context.Context, contractAddr, from, receiver common.Address, amount *big.Int) (*types.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, from, contractAddr, nil, k.abi, "transfer", receiver, amount)
	if err != nil {
		return nil, err
	}
	return unpackRetIsOk(k.abi, "transfer", res)
}

func (k ERC20TokenKeeper) TransferFrom(ctx context.Context, contractAddr, from, sender, receiver common.Address, amount *big.Int) (*types.MsgEthereumTxResponse, error) {
	res, err := k.ApplyContract(ctx, from, contractAddr, nil, k.abi, "transferFrom", sender, receiver, amount)
	if err != nil {
		return nil, err
	}
	return unpackRetIsOk(k.abi, "transferFrom", res)
}

func (k ERC20TokenKeeper) TransferOwnership(ctx context.Context, contractAddr, owner, newOwner common.Address) (*types.MsgEthereumTxResponse, error) {
	return k.ApplyContract(ctx, owner, contractAddr, nil, k.abi, "transferOwnership", newOwner)
}

func (k ERC20TokenKeeper) Withdraw(ctx context.Context, contractAddr, from, receiver common.Address, amount *big.Int) (*types.MsgEthereumTxResponse, error) {
	return k.ApplyContract(ctx, from, contractAddr, nil, k.abi, "withdraw0", receiver, amount)
}

func (k ERC20TokenKeeper) WithdrawToSelf(ctx context.Context, contractAddr, from common.Address, amount *big.Int) (*types.MsgEthereumTxResponse, error) {
	return k.ApplyContract(ctx, from, contractAddr, nil, k.abi, "withdraw", from, amount)
}

func (k ERC20TokenKeeper) Deposit(ctx context.Context, contractAddr, from common.Address, amount *big.Int) (*types.MsgEthereumTxResponse, error) {
	return k.ApplyContract(ctx, from, contractAddr, amount, k.abi, "deposit")
}
