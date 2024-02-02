package keeper

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ethermint "github.com/evmos/ethermint/types"
	"github.com/evmos/ethermint/x/evm/types"
	abci "github.com/tendermint/tendermint/abci/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
)

// InitGenesis initializes genesis state based on exported genesis
func (k *Keeper) InitGenesis(ctx sdk.Context, accountKeeper types.AccountKeeper, data types.GenesisState) []abci.ValidatorUpdate {
	if err := k.SetParams(ctx, data.Params); err != nil {
		panic(err)
	}
	// ensure evm module account is set
	if acc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
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
			panic(fmt.Errorf("account %s must be an EthAccount interface, got %T", account.Address, acc))
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

	// init logic contract
	initContract := []fxtypes.Contract{fxtypes.GetFIP20(), fxtypes.GetWFX()}
	for _, contract := range initContract {
		if len(contract.Code) == 0 || contract.Address == common.HexToAddress(fxtypes.EmptyEvmAddress) {
			panic(fmt.Sprintf("invalid contract: %s", contract.Address.String()))
		}
		if err := k.CreateContractWithCode(ctx, contract.Address, contract.Code); err != nil {
			panic(fmt.Sprintf("create contract %s with code error %s", contract.Address.String(), err.Error()))
		}
	}

	return []abci.ValidatorUpdate{}
}
