package contract

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// Deprecated: please use ERC20TokenKeeper
type ERC20ABI struct {
	abi abi.ABI
}

func NewERC20ABI() ERC20ABI {
	return ERC20ABI{}
}

func (e ERC20ABI) PackName() (data []byte, err error) {
	data, err = e.abi.Pack("name")
	if err != nil {
		return nil, fmt.Errorf("pack name: %s", err.Error())
	}
	return data, err
}

func (e ERC20ABI) UnpackName(ret []byte) (string, error) {
	var unpackedRet struct{ Value string }
	if err := e.abi.UnpackIntoInterface(&unpackedRet, "name", ret); err != nil {
		return "", fmt.Errorf("unpack name: %s", err.Error())
	}
	return unpackedRet.Value, nil
}

func (e ERC20ABI) PackSymbol() (data []byte, err error) {
	data, err = e.abi.Pack("symbol")
	if err != nil {
		return nil, fmt.Errorf("pack symbol: %s", err.Error())
	}
	return data, err
}

func (e ERC20ABI) UnpackSymbol(ret []byte) (string, error) {
	var unpackedRet struct{ Value string }
	if err := e.abi.UnpackIntoInterface(&unpackedRet, "symbol", ret); err != nil {
		return "", fmt.Errorf("unpack symbol: %s", err.Error())
	}
	return unpackedRet.Value, nil
}

func (e ERC20ABI) PackDecimals() (data []byte, err error) {
	data, err = e.abi.Pack("decimals")
	if err != nil {
		return nil, fmt.Errorf("pack decimals: %s", err.Error())
	}
	return data, err
}

func (e ERC20ABI) UnpackDecimals(ret []byte) (uint8, error) {
	var unpackedRet struct{ Value uint8 }
	if err := e.abi.UnpackIntoInterface(&unpackedRet, "decimals", ret); err != nil {
		return 0, fmt.Errorf("unpack decimals: %s", err.Error())
	}
	return unpackedRet.Value, nil
}

func (e ERC20ABI) PackBalanceOf(account common.Address) (data []byte, err error) {
	data, err = e.abi.Pack("balanceOf", account)
	if err != nil {
		return nil, fmt.Errorf("pack balanceOf: %s", err.Error())
	}
	return data, err
}

func (e ERC20ABI) UnpackBalanceOf(ret []byte) (*big.Int, error) {
	var unpackedRet struct{ Value *big.Int }
	if err := e.abi.UnpackIntoInterface(&unpackedRet, "balanceOf", ret); err != nil {
		return nil, fmt.Errorf("unpack balanceOf: %s", err.Error())
	}
	return unpackedRet.Value, nil
}

func (e ERC20ABI) PackTotalSupply() (data []byte, err error) {
	data, err = e.abi.Pack("totalSupply")
	if err != nil {
		return nil, fmt.Errorf("pack totalSupply: %s", err.Error())
	}
	return data, err
}

func (e ERC20ABI) UnpackTotalSupply(ret []byte) (*big.Int, error) {
	var unpackedRet struct{ Value *big.Int }
	if err := e.abi.UnpackIntoInterface(&unpackedRet, "totalSupply", ret); err != nil {
		return nil, fmt.Errorf("unpack totalSupply: %s", err.Error())
	}
	return unpackedRet.Value, nil
}

func (e ERC20ABI) PackApprove(spender common.Address, amount *big.Int) (data []byte, err error) {
	data, err = e.abi.Pack("approve", spender, amount)
	if err != nil {
		return nil, fmt.Errorf("pack approve: %s", err.Error())
	}
	return data, err
}

func (e ERC20ABI) PackAllowance(owner, spender common.Address) (data []byte, err error) {
	data, err = e.abi.Pack("allowance", owner, spender)
	if err != nil {
		return nil, fmt.Errorf("pack allowance: %s", err.Error())
	}
	return data, err
}

func (e ERC20ABI) PackTransferFrom(sender, to common.Address, amount *big.Int) (data []byte, err error) {
	data, err = e.abi.Pack("transferFrom", sender, to, amount)
	if err != nil {
		return nil, fmt.Errorf("pack transferFrom: %s", err.Error())
	}
	return data, err
}

func (e ERC20ABI) UnpackTransferFrom(ret []byte) (bool, error) {
	var unpackedRet struct{ Value bool }
	if err := e.abi.UnpackIntoInterface(&unpackedRet, "transferFrom", ret); err != nil {
		return false, fmt.Errorf("unpack transferFrom: %s", err.Error())
	}
	return unpackedRet.Value, nil
}

func (e ERC20ABI) PackTransfer(to common.Address, amount *big.Int) (data []byte, err error) {
	data, err = e.abi.Pack("transfer", to, amount)
	if err != nil {
		return nil, fmt.Errorf("pack transfer: %s", err.Error())
	}
	return data, err
}

func (e ERC20ABI) PackBurn(account common.Address, amount *big.Int) (data []byte, err error) {
	data, err = e.abi.Pack("burn", account, amount)
	if err != nil {
		return nil, fmt.Errorf("pack burn: %s", err.Error())
	}
	return data, err
}

func (e ERC20ABI) PackMint(account common.Address, amount *big.Int) (data []byte, err error) {
	data, err = e.abi.Pack("mint", account, amount)
	if err != nil {
		return nil, fmt.Errorf("pack mint: %s", err.Error())
	}
	return data, err
}

func (e ERC20ABI) PackDeposit() (data []byte, err error) {
	data, err = e.abi.Pack("deposit")
	if err != nil {
		return nil, fmt.Errorf("pack deposit: %s", err.Error())
	}
	return data, err
}

func (e ERC20ABI) PackWithdraw(recipient common.Address, amount *big.Int) (data []byte, err error) {
	data, err = e.abi.Pack("withdraw0", recipient, amount)
	if err != nil {
		return nil, fmt.Errorf("pack withdraw: %s", err.Error())
	}
	return data, err
}

func (e ERC20ABI) PackTransferOwnership(newOwner common.Address) (data []byte, err error) {
	data, err = e.abi.Pack("transferOwnership", newOwner)
	if err != nil {
		return nil, fmt.Errorf("pack transferOwnership: %s", err.Error())
	}
	return data, err
}
