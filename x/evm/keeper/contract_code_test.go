package keeper_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	tmrand "github.com/tendermint/tendermint/libs/rand"

	fxtypes "github.com/functionx/fx-core/v3/types"
)

func (suite *KeeperTestSuite) TestKeeper_CreateContractWithCode() {
	// CreateContractWithCode is in init genesis
	account := suite.app.EvmKeeper.GetAccount(suite.ctx, fxtypes.GetWFX().Address)
	suite.NotNil(account)

	suite.Equal(uint64(0), account.Nonce)
	code := suite.app.EvmKeeper.GetCode(suite.ctx, common.BytesToHash(account.CodeHash))
	suite.Equal(fxtypes.GetWFX().Code, code)
}

func (suite *KeeperTestSuite) TestKeeper_UpdateContractCode() {
	updateCode := []byte{1, 2, 3}
	err := suite.app.EvmKeeper.UpdateContractCode(suite.ctx, fxtypes.GetWFX().Address, updateCode)
	suite.NoError(err)

	account := suite.app.EvmKeeper.GetAccount(suite.ctx, fxtypes.GetWFX().Address)
	suite.NotNil(account)

	suite.Equal(uint64(0), account.Nonce)
	code := suite.app.EvmKeeper.GetCode(suite.ctx, common.BytesToHash(account.CodeHash))
	suite.Equal(updateCode, code)
}

func (suite *KeeperTestSuite) TestKeeper_DeployContract() {
	erc1967Proxy := fxtypes.GetERC1967Proxy()
	erc20 := fxtypes.GetERC20()
	contract, err := suite.app.EvmKeeper.DeployContract(suite.ctx, suite.signer.Address(), erc1967Proxy.ABI, erc1967Proxy.Bin, erc20.Address, []byte{})
	suite.NoError(err)

	nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, suite.signer.Address().Bytes())
	suite.NoError(err)

	contractAddr := crypto.CreateAddress(suite.signer.Address(), nonce-1)
	suite.Equal(contractAddr, contract)
}

func (suite *KeeperTestSuite) TestKeeper_DeployUpgradableContract() {
	erc20 := fxtypes.GetERC20()
	initializeArgs := []interface{}{"FunctionX USD", "fxUSD", uint8(18), suite.app.Erc20Keeper.ModuleAddress()}
	contract, err := suite.app.EvmKeeper.DeployUpgradableContract(suite.ctx, suite.signer.Address(), erc20.Address, nil, &erc20.ABI, initializeArgs...)
	suite.NoError(err)

	nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, suite.signer.Address().Bytes())
	suite.NoError(err)

	contractAddr := crypto.CreateAddress(suite.signer.Address(), nonce-1)
	suite.Equal(contractAddr, contract)
}

func (suite *KeeperTestSuite) TestKeeper_QueryContract() {
	erc20 := fxtypes.GetERC20()
	initializeArgs := []interface{}{"FunctionX USD", "fxUSD", uint8(18), suite.app.Erc20Keeper.ModuleAddress()}
	contract, err := suite.app.EvmKeeper.DeployUpgradableContract(suite.ctx, suite.signer.Address(), erc20.Address, nil, &erc20.ABI, initializeArgs...)
	suite.NoError(err)

	var nameRes struct{ Value string }
	err = suite.app.EvmKeeper.QueryContract(suite.ctx, suite.signer.Address(), contract, erc20.ABI, "name", &nameRes)
	suite.NoError(err)
	suite.Equal(nameRes.Value, "FunctionX USD")

	var balanceRes struct{ Value *big.Int }
	err = suite.app.EvmKeeper.QueryContract(suite.ctx, suite.signer.Address(), contract, erc20.ABI, "balanceOf", &balanceRes, suite.signer.Address())
	suite.NoError(err)
	suite.Equal(big.NewInt(0).String(), balanceRes.Value.String())
}

func (suite *KeeperTestSuite) TestKeeper_ApplyContract() {
	erc20 := fxtypes.GetERC20()
	initializeArgs := []interface{}{"FunctionX USD", "fxUSD", uint8(18), suite.app.Erc20Keeper.ModuleAddress()}
	contract, err := suite.app.EvmKeeper.DeployUpgradableContract(suite.ctx, suite.signer.Address(), erc20.Address, nil, &erc20.ABI, initializeArgs...)
	suite.NoError(err)

	mintAmt := int64(tmrand.Uint32())
	_, err = suite.app.EvmKeeper.ApplyContract(suite.ctx, suite.signer.Address(), contract, erc20.ABI, "mint", suite.signer.Address(), big.NewInt(mintAmt))
	suite.NoError(err)

	var balanceRes struct{ Value *big.Int }
	err = suite.app.EvmKeeper.QueryContract(suite.ctx, suite.signer.Address(), contract, erc20.ABI, "balanceOf", &balanceRes, suite.signer.Address())
	suite.NoError(err)
	suite.Equal(big.NewInt(mintAmt).String(), balanceRes.Value.String())
}
