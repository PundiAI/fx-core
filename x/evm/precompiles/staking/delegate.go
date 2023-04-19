package staking

import (
	"errors"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v4/x/evm/types"
)

func (c *Contract) Delegate(ctx sdk.Context, evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	if readonly {
		return nil, errors.New("delegate method not readonly")
	}
	// parse args
	var args DelegateArgs
	err := types.ParseMethodArgs(DelegateMethod, &args, contract.Input[4:])
	if err != nil {
		return nil, err
	}
	amount := contract.Value()
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("invalid delegate amount: %s", amount.String())
	}
	valAddr := args.GetValidator()
	val, found := c.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, fmt.Errorf("validator not found: %s", valAddr.String())
	}

	sender := sdk.AccAddress(contract.Caller().Bytes())
	evmDenom := c.evmKeeper.GetParams(ctx).EvmDenom

	// sub evm balance and mint delegate amount
	evm.StateDB.SubBalance(contract.Address(), amount)
	coins := sdk.NewCoins(sdk.NewCoin(evmDenom, sdkmath.NewIntFromBigInt(amount)))
	if err = c.bankKeeper.MintCoins(ctx, evmtypes.ModuleName, coins); err != nil {
		return nil, fmt.Errorf("mint operation failed: %s", err.Error())
	}
	if err = c.bankKeeper.SendCoinsFromModuleToAccount(ctx, evmtypes.ModuleName, sender, coins); err != nil {
		return nil, fmt.Errorf("send operation failed: %s", err.Error())
	}

	// withdraw rewards if delegation exist, add reward to evm state balance
	reward := big.NewInt(0)
	if _, found = c.stakingKeeper.GetDelegation(ctx, sender, valAddr); found {
		if reward, err = c.withdraw(ctx, evm, contract.Caller(), valAddr, evmDenom); err != nil {
			return nil, err
		}
	}

	// delegate amount
	shares, err := c.stakingKeeper.Delegate(ctx, sender, sdkmath.NewIntFromBigInt(amount), stakingtypes.Unbonded, val, true)
	if err != nil {
		return nil, err
	}

	// add delegate log
	if err := c.AddLog(DelegateEvent, []common.Hash{contract.Caller().Hash()},
		args.Validator, amount, shares.TruncateInt().BigInt()); err != nil {
		return nil, err
	}

	// TODO truncate shares, decimal 18
	return DelegateMethod.Outputs.Pack(shares.TruncateInt().BigInt(), reward)
}
