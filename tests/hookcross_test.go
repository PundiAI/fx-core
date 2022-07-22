package tests

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"

	"github.com/functionx/fx-core/v2/app/helpers"
	fxtypes "github.com/functionx/fx-core/v2/types"
	bsctypes "github.com/functionx/fx-core/v2/x/bsc/types"
	polygontypes "github.com/functionx/fx-core/v2/x/polygon/types"
)

type HookCrossTestSuite struct {
	*TestSuite
	BSCCrossChain     CrosschainTestSuite
	PolygonCrossChain CrosschainTestSuite
	ERC20             ERC20TestSuite
}

func TestERC20TestSuite(t *testing.T) {
	testSuite := NewTestSuite()
	erc20TestSuite := &HookCrossTestSuite{
		TestSuite:         testSuite,
		BSCCrossChain:     NewCrosschainWithTestSuite(bsctypes.ModuleName, testSuite),
		PolygonCrossChain: NewCrosschainWithTestSuite(polygontypes.ModuleName, testSuite),
		ERC20:             NewERC20WithTestSuite(testSuite),
	}
	suite.Run(t, erc20TestSuite)
}

func (suite *HookCrossTestSuite) SetupSuite() {
	err := os.Setenv("GO_ENV", "testing")
	suite.NoError(err)
	fxtypes.SetTestingManyToOneBlock(func() int64 { return 5 })

	suite.TestSuite.SetupSuite()

	suite.Send(suite.BSCCrossChain.OracleAddr(), helpers.NewCoin(sdk.NewInt(10_100).MulRaw(1e18)))
	suite.Send(suite.BSCCrossChain.BridgerFxAddr(), helpers.NewCoin(sdk.NewInt(1_000).MulRaw(1e18)))
	suite.Send(suite.BSCCrossChain.AccAddr(), helpers.NewCoin(sdk.NewInt(1_000).MulRaw(1e18)))
	suite.BSCCrossChain.params = suite.BSCCrossChain.QueryParams()

	suite.Send(suite.PolygonCrossChain.OracleAddr(), helpers.NewCoin(sdk.NewInt(10_100).MulRaw(1e18)))
	suite.Send(suite.PolygonCrossChain.BridgerFxAddr(), helpers.NewCoin(sdk.NewInt(1_000).MulRaw(1e18)))
	suite.Send(suite.PolygonCrossChain.AccAddr(), helpers.NewCoin(sdk.NewInt(1_000).MulRaw(1e18)))
	suite.PolygonCrossChain.params = suite.PolygonCrossChain.QueryParams()

	suite.Send(suite.ERC20.Address(), helpers.NewCoin(sdk.NewInt(10_100).MulRaw(1e18)))
}

func (suite *HookCrossTestSuite) TestERC20ConvertDenom() {
	// bsc crosschain
	const bscUSDToken = "0x0000000000000000000000000000000000000001"
	bscUSDTDenom := fmt.Sprintf("%s%s", suite.BSCCrossChain.chainName, bscUSDToken)

	proposalId := suite.BSCCrossChain.SendUpdateChainOraclesProposal()
	suite.ProposalVote(suite.AdminPrivateKey(), proposalId, govtypes.OptionYes)
	suite.CheckProposal(proposalId, govtypes.StatusPassed)

	suite.BSCCrossChain.BondedOracle()
	suite.BSCCrossChain.SendOracleSetConfirm()

	denom := suite.BSCCrossChain.AddBridgeTokenClaim("Tether USD", "USDT", 18, bscUSDToken, "")
	suite.Equal(denom, bscUSDTDenom)

	suite.BSCCrossChain.SendToFxClaim(bscUSDToken, sdk.NewInt(100).MulRaw(1e18), "")

	// polygon crosschain
	const polygonUSDToken = "0x0000000000000000000000000000000000000002"
	polygonUSDTDenom := fmt.Sprintf("%s%s", suite.PolygonCrossChain.chainName, polygonUSDToken)

	proposalId = suite.PolygonCrossChain.SendUpdateChainOraclesProposal()
	suite.ProposalVote(suite.AdminPrivateKey(), proposalId, govtypes.OptionYes)
	suite.CheckProposal(proposalId, govtypes.StatusPassed)

	suite.PolygonCrossChain.BondedOracle()
	suite.PolygonCrossChain.SendOracleSetConfirm()

	denom = suite.PolygonCrossChain.AddBridgeTokenClaim("Tether USD", "USDT", 18, polygonUSDToken, "")
	suite.Equal(denom, polygonUSDTDenom)

	suite.PolygonCrossChain.SendToFxClaim(polygonUSDToken, sdk.NewInt(100).MulRaw(1e18), "")

	// erc20
	usdtMetadata := banktypes.Metadata{
		Description: "description of the token",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "usdt",
				Exponent: uint32(0),
				Aliases:  []string{bscUSDTDenom, polygonUSDTDenom},
			}, {
				Denom:    "USDT",
				Exponent: uint32(18),
			},
		},
		Base:    "usdt",
		Display: "usdt",
		Name:    "Tether USD",
		Symbol:  "USDT",
	}
	proposalId = suite.ERC20.RegisterCoinProposal(usdtMetadata)
	suite.ProposalVote(suite.AdminPrivateKey(), proposalId, govtypes.OptionYes)
	suite.CheckProposal(proposalId, govtypes.StatusPassed)
	suite.ERC20.CheckRegisterCoin(usdtMetadata.Base, true)

	usdtTokenPair := suite.ERC20.TokenPair("usdt")
	suite.T().Log("token pair", usdtTokenPair.String())

	// bsc -> fx -> evm
	beforeSendToFx := suite.ERC20.BalanceOf(usdtTokenPair.GetERC20Contract(), suite.BSCCrossChain.HexAddr())
	suite.BSCCrossChain.SendToFxClaim(bscUSDToken, sdk.NewInt(100).MulRaw(1e18), "module/evm")
	afterSendToFx := suite.ERC20.BalanceOf(usdtTokenPair.GetERC20Contract(), suite.BSCCrossChain.HexAddr())
	suite.Equal(big.NewInt(0).Sub(afterSendToFx, beforeSendToFx), sdk.NewInt(100).MulRaw(1e18).BigInt())

	beforeBalances := suite.QueryBalances(suite.BSCCrossChain.AccAddr())
	suite.BSCCrossChain.SendToFxClaim(bscUSDToken, sdk.NewInt(100).MulRaw(1e18), "")
	afterSendToFxBalances := suite.QueryBalances(suite.BSCCrossChain.AccAddr())
	suite.Equal(afterSendToFxBalances.AmountOf("usdt").Sub(beforeBalances.AmountOf("usdt")), sdk.NewInt(100).MulRaw(1e18))

	suite.ERC20.ConvertDenom(suite.BSCCrossChain.privKey, suite.BSCCrossChain.AccAddr(), sdk.NewCoin("usdt", sdk.NewInt(100).MulRaw(1e18)), "bsc")
	afterConvertDenomUSDT := suite.QueryBalances(suite.BSCCrossChain.AccAddr())
	suite.Equal(afterConvertDenomUSDT.AmountOf(bscUSDTDenom).Sub(afterSendToFxBalances.AmountOf(bscUSDTDenom)), sdk.NewInt(100).MulRaw(1e18))

	suite.ERC20.ConvertDenom(suite.BSCCrossChain.privKey, suite.BSCCrossChain.AccAddr(), sdk.NewCoin(bscUSDTDenom, sdk.NewInt(100).MulRaw(1e18)), "")
	afterConvertDenomBscUSDT := suite.QueryBalances(suite.BSCCrossChain.AccAddr())
	suite.Equal(afterConvertDenomBscUSDT.AmountOf("usdt").Sub(afterConvertDenomUSDT.AmountOf("usdt")), sdk.NewInt(100).MulRaw(1e18))
}
