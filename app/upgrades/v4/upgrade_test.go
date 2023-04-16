package v4_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/assert"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v4/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v4/types"
)

func TestAccountConvertModuleAccount(t *testing.T) {
	valSet, genAccs, balances := helpers.GenerateGenesisValidator(1, nil)
	myApp := helpers.SetupWithGenesisValSet(t, valSet, genAccs, balances...)
	ctx := myApp.NewContext(false, tmproto.Header{Height: myApp.LastBlockHeight()})
	moduleAddress, _ := myApp.AccountKeeper.GetModuleAddressAndPermissions(types.ModuleName)
	account := myApp.AccountKeeper.GetAccount(ctx, moduleAddress)
	assert.NotNil(t, account)
	balance := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1000).MulRaw(1e18)))
	err := myApp.BankKeeper.MintCoins(ctx, types.ModuleName, balance)
	assert.NoError(t, err)
	err = myApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, account.GetAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(500).MulRaw(1e18))))
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("%s is not allowed to receive funds: unauthorized", account.GetAddress().String()))
}
