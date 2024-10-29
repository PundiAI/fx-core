package precompile_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v8/contract"
	testscontract "github.com/functionx/fx-core/v8/tests/contract"
	"github.com/functionx/fx-core/v8/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v8/types"
)

type PrecompileTestSuite struct {
	helpers.BaseSuite

	signer          *helpers.Signer
	stakingTestAddr common.Address
	stakingAddr     common.Address

	helpers.StakingPrecompileSuite
}

func TestPrecompileTestSuite(t *testing.T) {
	fxtypes.SetConfig(true)
	suite.Run(t, new(PrecompileTestSuite))
}

func (suite *PrecompileTestSuite) SetupSubTest() {
	suite.SetupTest()
}

func (suite *PrecompileTestSuite) SetupTest() {
	suite.MintValNumber = 2
	suite.BaseSuite.SetupTest()
	suite.Commit(10)

	suite.signer = suite.NewSigner()
	suite.MintToken(suite.signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10000).Mul(sdkmath.NewInt(1e18))))

	var err error
	suite.stakingTestAddr, err = suite.App.EvmKeeper.DeployContract(suite.Ctx, suite.signer.Address(), contract.MustABIJson(testscontract.StakingTestMetaData.ABI), contract.MustDecodeHex(testscontract.StakingTestMetaData.Bin))
	suite.Require().NoError(err)

	suite.stakingAddr = common.HexToAddress(contract.StakingAddress)

	suite.StakingPrecompileSuite = helpers.NewStakingPrecompileSuite(suite.Require(), suite.signer, suite.App.EvmKeeper, suite.stakingAddr)
}

func (suite *PrecompileTestSuite) DistributionQueryClient(ctx sdk.Context) distributiontypes.QueryClient {
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, suite.App.InterfaceRegistry())
	distributiontypes.RegisterQueryServer(queryHelper, distributionkeeper.NewQuerier(suite.App.DistrKeeper))
	return distributiontypes.NewQueryClient(queryHelper)
}

func (suite *PrecompileTestSuite) PrecompileStakingDelegateV2(signer *helpers.Signer, val sdk.ValAddress, amt *big.Int) *big.Int {
	suite.MintToken(signer.AccAddress(), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewIntFromBigInt(amt)))

	_, amountBefore := suite.Delegation(suite.Ctx, contract.DelegationArgs{
		Validator: val.String(),
		Delegator: signer.Address(),
	})

	suite.WithSigner(signer)
	res := suite.DelegateV2(suite.Ctx, contract.DelegateV2Args{
		Validator: val.String(),
		Amount:    amt,
	})
	suite.Require().False(res.Failed(), res.VmError)

	shares, amount := suite.Delegation(suite.Ctx, contract.DelegationArgs{
		Validator: val.String(),
		Delegator: signer.Address(),
	})

	suite.Require().Equal(amt.String(), big.NewInt(0).Sub(amount, amountBefore).String())
	return shares
}

func (suite *PrecompileTestSuite) PrecompileStakingWithdraw(signer *helpers.Signer, val sdk.ValAddress) sdk.Coins {
	balanceBefore := suite.App.BankKeeper.GetAllBalances(suite.Ctx, signer.AccAddress())

	suite.WithSigner(signer)
	res, _ := suite.Withdraw(suite.Ctx, contract.WithdrawArgs{
		Validator: val.String(),
	})
	suite.Require().False(res.Failed(), res.VmError)

	balanceAfter := suite.App.BankKeeper.GetAllBalances(suite.Ctx, signer.AccAddress())
	return balanceAfter.Sub(balanceBefore...)
}
