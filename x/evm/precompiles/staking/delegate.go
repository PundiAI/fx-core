package staking

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v3/x/evm/types"
)

var DelegateMethod = abi.NewMethod(DelegateMethodName, DelegateMethodName, abi.Function, "", false, true,
	abi.Arguments{
		abi.Argument{
			Name: "validator",
			Type: types.TypeString,
		},
		abi.Argument{
			Name: "amount",
			Type: types.TypeUint256,
		},
	},
	abi.Arguments{
		abi.Argument{
			Name: "shares",
			Type: types.TypeUint256,
		},
	},
)

func (c *Contract) Delegate(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("delegate method not readonly")
	}
	args, err := DelegateMethod.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, errors.New("failed to unpack input")
	}
	valAddrStr := args[0].(string)
	valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address: %s", valAddrStr)
	}
	amount := args[1].(*big.Int)
	if amount.Sign() <= 0 {
		return nil, fmt.Errorf("invalid amount: %s", amount.String())
	}

	sender := sdk.AccAddress(contract.CallerAddress.Bytes())

	// check contract value
	balance := evm.StateDB.GetBalance(contract.Address())
	if balance.Cmp(amount) != 0 || contract.Value().Cmp(amount) != 0 {
		return nil, fmt.Errorf("invalid msg.value: %s", contract.Value().String())
	}

	val, found := c.stakingKeeper.GetValidator(c.ctx, valAddr)
	if !found {
		return nil, fmt.Errorf("validator not found: %s", valAddr.String())
	}

	snapshot := evm.StateDB.Snapshot()
	cacheCtx, commit := c.ctx.CacheContext()
	bondDenom := c.stakingKeeper.BondDenom(cacheCtx)

	// sub evm balance and mint delegate amount
	evm.StateDB.SubBalance(contract.Address(), amount)
	coins := sdk.NewCoins(sdk.NewCoin(bondDenom, sdk.NewIntFromBigInt(amount)))
	if err = c.bankKeeper.MintCoins(cacheCtx, evmtypes.ModuleName, coins); err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		return nil, fmt.Errorf("mint operation failed: %s", err.Error())
	}
	if err = c.bankKeeper.SendCoinsFromModuleToAccount(
		cacheCtx, evmtypes.ModuleName, sender, coins); err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		return nil, fmt.Errorf("send operation failed: %s", err.Error())
	}

	// withdraw rewards if delegation exist
	if _, found := c.stakingKeeper.GetDelegation(cacheCtx, sender, valAddr); found {
		// receive the commission reward in advance and add it to evm balance
		rewards, err := c.distrKeeper.WithdrawDelegationRewards(cacheCtx, sender, valAddr)
		if err != nil {
			evm.StateDB.RevertToSnapshot(snapshot)
			return nil, fmt.Errorf("withdraw failed: %s", err.Error())
		}
		evm.StateDB.AddBalance(contract.CallerAddress, rewards.AmountOf(bondDenom).BigInt())
	}

	// delegate amount
	shares, err := c.stakingKeeper.Delegate(cacheCtx, sender, sdk.NewIntFromBigInt(amount), stakingtypes.Unbonded, val, true)
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		return nil, err
	}
	commit()

	// todo truncate shares, decimal 18
	return DelegateMethod.Outputs.Pack(shares.TruncateInt().BigInt())
}
