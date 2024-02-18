package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/erc20/types"
)

// QueryERC20 returns the data of a deployed ERC20 contract
func (k Keeper) QueryERC20(ctx sdk.Context, contract common.Address) (types.ERC20Data, error) {
	erc20 := fxtypes.GetFIP20().ABI

	// Name
	var nameRes struct{ Value string }
	err := k.evmKeeper.QueryContract(ctx, k.moduleAddress, contract, erc20, "name", &nameRes)
	if err != nil {
		return types.ERC20Data{}, err
	}

	// Symbol
	var symbolRes struct{ Value string }
	err = k.evmKeeper.QueryContract(ctx, k.moduleAddress, contract, erc20, "symbol", &symbolRes)
	if err != nil {
		return types.ERC20Data{}, err
	}

	// Decimals
	var decimalRes struct{ Value uint8 }
	err = k.evmKeeper.QueryContract(ctx, k.moduleAddress, contract, erc20, "decimals", &decimalRes)
	if err != nil {
		return types.ERC20Data{}, err
	}

	return types.NewERC20Data(nameRes.Value, symbolRes.Value, decimalRes.Value), nil
}

func (k Keeper) DeployUpgradableToken(ctx sdk.Context, from common.Address, name, symbol string, decimals uint8) (common.Address, error) {
	var tokenContract fxtypes.Contract
	if symbol == fxtypes.DefaultDenom {
		tokenContract = fxtypes.GetWFX()
		name = fmt.Sprintf("Wrapped %s", name)
		symbol = fmt.Sprintf("W%s", symbol)
	} else {
		tokenContract = fxtypes.GetFIP20()
	}
	k.Logger(ctx).Info("deploy token contract", "name", name, "symbol", symbol, "decimals", decimals)

	return k.evmKeeper.DeployUpgradableContract(ctx, from, tokenContract.Address, nil, &tokenContract.ABI, name, symbol, decimals, k.moduleAddress)
}

// monitorApprovalEvent returns an error if the given transactions logs include
// an unexpected `approve` event
func (k Keeper) monitorApprovalEvent(res *evmtypes.MsgEthereumTxResponse) error {
	if res == nil || len(res.Logs) == 0 {
		return nil
	}

	logApprovalSigHash := crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))

	for _, log := range res.Logs {
		if log.Topics[0] == logApprovalSigHash.Hex() {
			return errorsmod.Wrapf(
				types.ErrUnexpectedEvent, "approval event",
			)
		}
	}

	return nil
}
