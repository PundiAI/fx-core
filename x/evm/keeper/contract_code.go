package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/x/evm/statedb"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/functionx/fx-core/v3/x/evm/types"
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
		return sdkerrors.Wrap(evmtypes.ErrInvalidAccount, address.String())
	}
	codeHash := crypto.Keccak256Hash(contractCode).Bytes()
	if bytes.Equal(codeHash, acc.CodeHash) {
		return fmt.Errorf("update the same code hash: %s", address.String())
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
