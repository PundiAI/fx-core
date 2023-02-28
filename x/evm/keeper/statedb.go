package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/evmos/ethermint/types"
	"github.com/evmos/ethermint/x/evm/statedb"
	"github.com/evmos/ethermint/x/evm/types"
)

// SetAccount updates nonce/balance/codeHash together.
func (k *Keeper) SetAccount(ctx sdk.Context, addr common.Address, account statedb.Account) error {
	// update account
	cosmosAddr := sdk.AccAddress(addr.Bytes())
	acct := k.accountKeeper.GetAccount(ctx, cosmosAddr)
	if acct == nil {
		acct = k.accountKeeper.NewAccountWithAddress(ctx, cosmosAddr)
	}

	if err := acct.SetSequence(account.Nonce); err != nil {
		return err
	}

	codeHash := common.BytesToHash(account.CodeHash)
	ethAcct, ok := acct.(ethermint.EthAccountI)
	if ok {
		if err := ethAcct.SetCodeHash(codeHash); err != nil {
			return err
		}
	} else {
		if account.IsContract() {
			if baseAcct, isBaseAccount := acct.(*authtypes.BaseAccount); isBaseAccount {
				acct = &ethermint.EthAccount{
					BaseAccount: baseAcct,
					CodeHash:    codeHash.Hex(),
				}
			} else {
				return errorsmod.Wrapf(types.ErrInvalidAccount, "type %T, address %s", acct, addr)
			}
		}
	}

	k.accountKeeper.SetAccount(ctx, acct)

	if err := k.SetBalance(ctx, addr, account.Balance); err != nil {
		return err
	}

	k.Logger(ctx).Debug(
		"account updated",
		"ethereum-address", addr.Hex(),
		"nonce", account.Nonce,
		"codeHash", codeHash.Hex(),
		"balance", account.Balance,
	)
	return nil
}
