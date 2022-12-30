package evm

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ethermint "github.com/evmos/ethermint/types"
	"github.com/evmos/ethermint/x/evm/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/functionx/fx-core/v3/x/evm/keeper"
)

// InitGenesis initializes genesis state based on exported genesis
func InitGenesis(
	ctx sdk.Context,
	k *keeper.Keeper,
	accountKeeper types.AccountKeeper,
	data types.GenesisState,
) []abci.ValidatorUpdate {
	k.SetParams(ctx, data.Params)

	// ensure evm module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the EVM module account has not been set")
	}

	for _, account := range data.Accounts {
		address := common.HexToAddress(account.Address)
		accAddress := sdk.AccAddress(address.Bytes())
		// check that the EVM balance the matches the account balance
		acc := accountKeeper.GetAccount(ctx, accAddress)
		if acc == nil {
			panic(fmt.Errorf("account not found for address %s", account.Address))
		}

		ethAcct, ok := acc.(ethermint.EthAccountI)
		if !ok {
			panic(
				fmt.Errorf("account %s must be an EthAccount interface, got %T",
					account.Address, acc,
				),
			)
		}

		code := common.Hex2Bytes(account.Code)
		codeHash := crypto.Keccak256Hash(code)
		if !bytes.Equal(ethAcct.GetCodeHash().Bytes(), codeHash.Bytes()) {
			panic("code don't match codeHash")
		}

		k.SetCode(ctx, codeHash.Bytes(), code)

		for _, storage := range account.Storage {
			k.SetState(ctx, address, common.HexToHash(storage.Key), common.HexToHash(storage.Value).Bytes())
		}
	}

	return []abci.ValidatorUpdate{}
}
