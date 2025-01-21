package crosschain_test

import (
	"math/big"
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/pundiai/fx-core/v8/contract"
	"github.com/pundiai/fx-core/v8/precompiles/crosschain"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	ethtypes "github.com/pundiai/fx-core/v8/x/eth/types"
)

func TestCrosschainABI(t *testing.T) {
	crosschainABI := crosschain.NewCrosschainABI()

	require.Len(t, crosschainABI.Method.Inputs, 6)
	require.Len(t, crosschainABI.Method.Outputs, 1)

	require.Len(t, crosschainABI.Event.Inputs, 8)
}

func (suite *CrosschainPrecompileTestSuite) TestContract_Crosschain() {
	suite.AddBridgeToken(fxtypes.DefaultSymbol, true)

	suite.App.CrosschainKeepers.GetKeeper(suite.chainName).
		SetLastObservedBlockHeight(suite.Ctx, 100, 100)

	balance := suite.Balance(suite.signer.AccAddress())

	txResponse := suite.Crosschain(suite.Ctx, big.NewInt(2), suite.signer.Address(),
		contract.CrosschainArgs{
			Token:   common.Address{},
			Receipt: helpers.GenExternalAddr(suite.chainName),
			Amount:  big.NewInt(1),
			Fee:     big.NewInt(1),
			Target:  contract.MustStrToByte32(suite.chainName),
			Memo:    "",
		},
	)
	suite.NotNil(txResponse)
	suite.Len(txResponse.Logs, 1)

	transferCoin := helpers.NewStakingCoin(2, 0)
	suite.AssertBalance(suite.signer.AccAddress(), balance.Sub(transferCoin)...)
	suite.AssertBalance(authtypes.NewModuleAddress(ethtypes.ModuleName), transferCoin)
}
