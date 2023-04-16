package keeper_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/suite"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v4/app"
	"github.com/functionx/fx-core/v4/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v4/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app    *app.App
	ctx    sdk.Context
	signer *helpers.Signer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	valNumber := tmrand.Intn(10) + 1
	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(valNumber, sdk.Coins{})

	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		ChainID:         fxtypes.MainnetChainId,
		Height:          suite.app.LastBlockHeight() + 1,
		ProposerAddress: valSet.Proposer.Address.Bytes(),
	})

	suite.signer = helpers.NewSigner(helpers.NewEthPrivKey())
	helpers.AddTestAddr(suite.app, suite.ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100).MulRaw(1e18))))
}

func (suite *KeeperTestSuite) DeployERC20Contract() common.Address {
	args, err := fxtypes.GetERC20().ABI.Pack("")
	suite.Require().NoError(err)

	nonce := suite.app.EvmKeeper.GetNonce(suite.ctx, suite.signer.Address())
	msg := ethtypes.NewMessage(
		suite.signer.Address(),
		nil,
		nonce,
		big.NewInt(0),
		1706373,
		big.NewInt(500*1e9),
		nil,
		nil,
		append(fxtypes.GetERC20().Bin, args...),
		nil,
		true,
	)

	rsp, err := suite.app.EvmKeeper.ApplyMessage(suite.ctx, msg, nil, true)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp.VmError)
	suite.Equal(rsp.GasUsed, uint64(1706373))
	contractAddress := crypto.CreateAddress(suite.signer.Address(), nonce)

	args, err = fxtypes.GetERC20().ABI.Pack("initialize", "Test Token", "TEST", uint8(18), helpers.GenerateAddress())
	suite.Require().NoError(err)

	msg = ethtypes.NewMessage(
		suite.signer.Address(),
		&contractAddress,
		suite.app.EvmKeeper.GetNonce(suite.ctx, suite.signer.Address()),
		big.NewInt(0),
		165000,
		big.NewInt(500*1e9),
		nil,
		nil,
		args,
		nil,
		true,
	)
	rsp, err = suite.app.EvmKeeper.ApplyMessage(suite.ctx, msg, nil, true)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp)

	amount := new(big.Int).Exp(big.NewInt(10), big.NewInt(20), nil)
	args, err = fxtypes.GetERC20().ABI.Pack("mint", suite.signer.Address(), amount)
	suite.Require().NoError(err)

	msg = ethtypes.NewMessage(
		suite.signer.Address(),
		&contractAddress,
		suite.app.EvmKeeper.GetNonce(suite.ctx, suite.signer.Address()),
		big.NewInt(0),
		71000,
		big.NewInt(500*1e9),
		nil,
		nil,
		args,
		nil,
		true,
	)
	rsp, err = suite.app.EvmKeeper.ApplyMessage(suite.ctx, msg, nil, true)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp)
	return contractAddress
}

func (suite *KeeperTestSuite) BalanceOf(contract, address common.Address) *big.Int {
	var balanceRes struct {
		Value *big.Int
	}
	err := suite.app.EvmKeeper.QueryContract(suite.ctx, contract, contract, fxtypes.GetERC20().ABI, "balanceOf", &balanceRes, address)
	suite.Require().NoError(err)
	return balanceRes.Value
}

func (suite *KeeperTestSuite) MintFeeCollector(coins sdk.Coins) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, evmtypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, evmtypes.ModuleName, authtypes.FeeCollectorName, coins)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) BurnEvmRefundFee(addr sdk.AccAddress, coins sdk.Coins) {
	err := suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, addr, authtypes.FeeCollectorName, coins)
	suite.Require().NoError(err)

	bal := suite.app.BankKeeper.GetBalance(suite.ctx, suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName), fxtypes.DefaultDenom)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, authtypes.FeeCollectorName, evmtypes.ModuleName, sdk.NewCoins(bal))
	suite.Require().NoError(err)

	err = suite.app.BankKeeper.BurnCoins(suite.ctx, evmtypes.ModuleName, sdk.NewCoins(bal))
	suite.Require().NoError(err)
}
