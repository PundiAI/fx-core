package bank_test

import (
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/precompiles/bank"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/erc20/types"
)

type BankPrecompileTestSuite struct {
	helpers.BaseSuite

	keeper *bank.Keeper
	helpers.BankPrecompileSuite
}

func TestBankPrecompileTestSuite(t *testing.T) {
	testingSuite := new(BankPrecompileTestSuite)
	suite.Run(t, testingSuite)
}

func (suite *BankPrecompileTestSuite) SetupTest() {
	suite.BaseSuite.SetupTest()
	suite.Ctx = suite.Ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(fxtypes.DefaultDenom, sdkmath.OneInt())))
	suite.Ctx = suite.Ctx.WithBlockGasMeter(storetypes.NewGasMeter(1e18))

	suite.keeper = bank.NewKeeper(suite.App.BankKeeper, suite.App.Erc20Keeper)
	suite.BankPrecompileSuite = helpers.NewBankPrecompileSuite(suite.Require(), suite.App.EvmKeeper)
}

func (suite *BankPrecompileTestSuite) SetupSubTest() {
	suite.SetupTest()
}

func (suite *BankPrecompileTestSuite) SetErc20Token(name string, token common.Address) {
	metadata := fxtypes.NewMetadata(name, strings.ToUpper(name), 18)
	_, err := suite.App.Erc20Keeper.AddERC20Token(suite.Ctx, metadata, token, types.OWNER_EXTERNAL)
	suite.Require().NoError(err)
}
