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

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
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
	args, err := contract.GetFIP20().ABI.Pack("")
	suite.Require().NoError(err)

	nonce := suite.app.EvmKeeper.GetNonce(suite.ctx, suite.signer.Address())
	msg := ethtypes.NewMessage(
		suite.signer.Address(),
		nil,
		nonce,
		big.NewInt(0),
		1711435,
		big.NewInt(500*1e9),
		nil,
		nil,
		append(contract.GetFIP20().Bin, args...),
		nil,
		true,
	)

	rsp, err := suite.app.EvmKeeper.ApplyMessage(suite.ctx, msg, nil, true)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp.VmError)
	suite.Equal(uint64(1711435), rsp.GasUsed)
	contractAddress := crypto.CreateAddress(suite.signer.Address(), nonce)

	args, err = contract.GetFIP20().ABI.Pack("initialize", "Test Token", "TEST", uint8(18), helpers.GenerateAddress())
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
	args, err = contract.GetFIP20().ABI.Pack("mint", suite.signer.Address(), amount)
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

func (suite *KeeperTestSuite) BalanceOf(contractAddr, address common.Address) *big.Int {
	var balanceRes struct {
		Value *big.Int
	}
	err := suite.app.EvmKeeper.QueryContract(suite.ctx, contractAddr, contractAddr, contract.GetFIP20().ABI, "balanceOf", &balanceRes, address)
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

func (suite *KeeperTestSuite) Owner(contractAddr common.Address) common.Address {
	var ownerRes struct {
		Value common.Address
	}
	err := suite.app.EvmKeeper.QueryContract(suite.ctx, contractAddr, contractAddr, contract.GetFIP20().ABI, "owner", &ownerRes)
	suite.Require().NoError(err)
	return ownerRes.Value
}

func (suite *KeeperTestSuite) Name(contractAddr common.Address) string {
	var nameRes struct {
		Value string
	}
	err := suite.app.EvmKeeper.QueryContract(suite.ctx, contractAddr, contractAddr, contract.GetFIP20().ABI, "name", &nameRes)
	suite.Require().NoError(err)
	return nameRes.Value
}

func (suite *KeeperTestSuite) Symbol(contractAddr common.Address) string {
	var symbolRes struct {
		Value string
	}
	err := suite.app.EvmKeeper.QueryContract(suite.ctx, contractAddr, contractAddr, contract.GetFIP20().ABI, "symbol", &symbolRes)
	suite.Require().NoError(err)
	return symbolRes.Value
}

func (suite *KeeperTestSuite) Decimals(contractAddr common.Address) uint8 {
	var decimalsRes struct {
		Value uint8
	}
	err := suite.app.EvmKeeper.QueryContract(suite.ctx, contractAddr, contractAddr, contract.GetFIP20().ABI, "decimals", &decimalsRes)
	suite.Require().NoError(err)
	return decimalsRes.Value
}

func (suite *KeeperTestSuite) TotalSupply(contractAddr common.Address) *big.Int {
	var totalSupplyRes struct {
		Value *big.Int
	}
	err := suite.app.EvmKeeper.QueryContract(suite.ctx, contractAddr, contractAddr, contract.GetFIP20().ABI, "totalSupply", &totalSupplyRes)
	suite.Require().NoError(err)
	return totalSupplyRes.Value
}

func (suite *KeeperTestSuite) Allowance(contractAddr, owner, spender common.Address) *big.Int {
	var allowanceRes struct {
		Value *big.Int
	}
	err := suite.app.EvmKeeper.QueryContract(suite.ctx, contractAddr, contractAddr, contract.GetFIP20().ABI, "allowance", &allowanceRes, owner, spender)
	suite.Require().NoError(err)
	return allowanceRes.Value
}

func (suite *KeeperTestSuite) Module(contractAddr common.Address) common.Address {
	var moduleRes struct {
		Value common.Address
	}
	err := suite.app.EvmKeeper.QueryContract(suite.ctx, contractAddr, contractAddr, contract.GetFIP20().ABI, "module", &moduleRes)
	suite.Require().NoError(err)
	return moduleRes.Value
}

func (suite *KeeperTestSuite) Approve(contractAddr, spender common.Address, amount *big.Int) {
	rsp, err := suite.app.EvmKeeper.ApplyContract(suite.ctx, spender, contractAddr, nil, contract.GetFIP20().ABI, "approve", suite.signer.Address(), amount)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp)
}

// function transfer(address recipient, uint256 amount) public virtual override returns (bool) {
func (suite *KeeperTestSuite) Transfer(contractAddr, recipient common.Address, amount *big.Int) {
	rsp, err := suite.app.EvmKeeper.ApplyContract(suite.ctx, suite.signer.Address(), contractAddr, nil, contract.GetFIP20().ABI, "transfer", recipient, amount)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp)
}

func (suite *KeeperTestSuite) TransferFrom(contractAddr, sender, recipient common.Address, amount *big.Int) {
	rsp, err := suite.app.EvmKeeper.ApplyContract(suite.ctx, suite.signer.Address(), contractAddr, nil, contract.GetFIP20().ABI, "transferFrom", sender, recipient, amount)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp)
}

func (suite *KeeperTestSuite) Mint(contractAddr, to common.Address, amount *big.Int) {
	rsp, err := suite.app.EvmKeeper.ApplyContract(suite.ctx, suite.signer.Address(), contractAddr, nil, contract.GetFIP20().ABI, "mint", to, amount)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp)
}

func (suite *KeeperTestSuite) Burn(contractAddr, from common.Address, amount *big.Int) {
	rsp, err := suite.app.EvmKeeper.ApplyContract(suite.ctx, suite.signer.Address(), contractAddr, nil, contract.GetFIP20().ABI, "burn", from, amount)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp)
}

func (suite *KeeperTestSuite) TransferOwnership(contractAddr, newOwner common.Address) {
	rsp, err := suite.app.EvmKeeper.ApplyContract(suite.ctx, suite.signer.Address(), contractAddr, nil, contract.GetFIP20().ABI, "transferOwnership", newOwner)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp)
}

func (suite *KeeperTestSuite) Deposit(contractAddr common.Address, amount *big.Int) {
	args, err := contract.GetWFX().ABI.Pack("deposit")
	suite.Require().NoError(err)

	msg := ethtypes.NewMessage(
		suite.signer.Address(),
		&contractAddr,
		suite.app.EvmKeeper.GetNonce(suite.ctx, suite.signer.Address()),
		amount,
		80000,
		big.NewInt(500*1e9),
		nil,
		nil,
		args,
		nil,
		true,
	)

	rsp, err := suite.app.EvmKeeper.ApplyMessage(suite.ctx, msg, evmtypes.NewNoOpTracer(), true)
	suite.Require().NoError(err)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp)
}

func (suite *KeeperTestSuite) Withdraw(contractAddr common.Address, amount *big.Int) {
	rsp, err := suite.app.EvmKeeper.ApplyContract(suite.ctx, suite.signer.Address(), contractAddr, nil, contract.GetWFX().ABI, "withdraw0", suite.signer.Address(), amount)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp)
}

func (suite *KeeperTestSuite) WithdrawSelf(contractAddr common.Address, amount *big.Int) {
	rsp, err := suite.app.EvmKeeper.ApplyContract(suite.ctx, suite.signer.Address(), contractAddr, nil, contract.GetWFX().ABI, "withdraw", amount)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp)
}

func (suite *KeeperTestSuite) CheckFIP20Method(contract common.Address, name, symbol string, decimals uint8, totalSupply *big.Int, owner common.Address) {
	account := helpers.NewSigner(helpers.NewEthPrivKey())
	helpers.AddTestAddr(suite.app, suite.ctx, account.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(100).MulRaw(1e18))))

	// name
	suite.Equal(name, suite.Name(contract))
	// symbol
	suite.Equal(symbol, suite.Symbol(contract))
	// decimals
	suite.Equal(decimals, suite.Decimals(contract))
	// totalSupply
	suite.True(suite.TotalSupply(contract).Cmp(totalSupply) == 0)
	// balanceOf
	suite.True(suite.BalanceOf(contract, suite.signer.Address()).Cmp(big.NewInt(0)) == 0)
	// allowance
	suite.True(suite.Allowance(contract, suite.signer.Address(), account.Address()).Cmp(big.NewInt(0)) == 0)
	// owner
	suite.Equal(owner, suite.Owner(contract))
	// approve
	suite.Approve(contract, account.Address(), big.NewInt(100))
	suite.True(suite.Allowance(contract, account.Address(), suite.signer.Address()).Cmp(big.NewInt(100)) == 0)
	// mint
	suite.Mint(contract, suite.signer.Address(), big.NewInt(100))
	suite.True(suite.BalanceOf(contract, suite.signer.Address()).Cmp(big.NewInt(100)) == 0)
	// burn
	suite.Burn(contract, suite.signer.Address(), big.NewInt(1))
	suite.True(suite.BalanceOf(contract, suite.signer.Address()).Cmp(big.NewInt(99)) == 0)
	suite.Burn(contract, suite.signer.Address(), big.NewInt(99))
	suite.True(suite.BalanceOf(contract, suite.signer.Address()).Cmp(big.NewInt(0)) == 0)
	// transfer
	suite.Mint(contract, suite.signer.Address(), big.NewInt(100))
	suite.Transfer(contract, account.Address(), big.NewInt(100))
	suite.True(suite.BalanceOf(contract, suite.signer.Address()).Cmp(big.NewInt(0)) == 0)
	suite.True(suite.BalanceOf(contract, account.Address()).Cmp(big.NewInt(100)) == 0)

	// transferFrom
	suite.Approve(contract, account.Address(), big.NewInt(100))
	suite.TransferFrom(contract, account.Address(), suite.signer.Address(), big.NewInt(100))
	suite.True(suite.BalanceOf(contract, suite.signer.Address()).Cmp(big.NewInt(100)) == 0)
	suite.True(suite.BalanceOf(contract, account.Address()).Cmp(big.NewInt(0)) == 0)
	suite.Burn(contract, suite.signer.Address(), big.NewInt(100))

	// module
	suite.Equal(common.BytesToAddress(suite.app.AccountKeeper.GetModuleAddress(evmtypes.ModuleName)), suite.Module(contract))
}
