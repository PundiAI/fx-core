package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/suite"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
	helpers2 "github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app    *app.App
	ctx    sdk.Context
	signer *helpers2.Signer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	valNumber := tmrand.Intn(100-1) + 1
	valSet, valAccounts, valBalances := helpers2.GenerateGenesisValidator(valNumber, sdk.Coins{})

	suite.app = helpers2.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.NewContext(false, tmproto.Header{
		ChainID:         fxtypes.MainnetChainId,
		Height:          suite.app.LastBlockHeight() + 1,
		ProposerAddress: valSet.Proposer.Address.Bytes(),
	})

	suite.signer = helpers2.NewSigner(helpers2.NewEthPrivKey())
	helpers2.AddTestAddr(suite.app, suite.ctx, suite.signer.AccAddress(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(100).MulRaw(1e18))))
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
		1558681,
		big.NewInt(500*1e9),
		nil,
		nil,
		append(fxtypes.GetERC20().Bin, args...),
		nil,
		true,
	)

	rsp, err := suite.app.EvmKeeper.ApplyMessage(suite.ctx, msg, nil, true)
	suite.Require().NoError(err)
	suite.Require().False(rsp.Failed(), rsp)
	suite.Equal(rsp.GasUsed, uint64(1558681))
	contractAddress := crypto.CreateAddress(suite.signer.Address(), nonce)

	args, err = fxtypes.GetERC20().ABI.Pack("initialize", "Test Token", "TEST", uint8(18), helpers2.GenerateAddress())
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
	data, err := fxtypes.GetERC20().ABI.Pack("balanceOf", address)
	suite.NoError(err)

	res, err := suite.app.EvmKeeper.CallEVMWithData(suite.ctx, contract, &contract, data, false)
	suite.NoError(err)

	var balanceRes struct {
		Value *big.Int
	}
	err = fxtypes.GetERC20().ABI.UnpackIntoInterface(&balanceRes, "balanceOf", res.Ret)
	suite.NoError(err)
	return balanceRes.Value
}

func (suite *KeeperTestSuite) TestCallEVMWithData() {
	erc20 := fxtypes.GetERC20()
	fromAcc := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, evmtypes.ModuleName)
	from := common.BytesToAddress(fromAcc.GetAddress())
	testCases := []struct {
		name     string
		from     common.Address
		malleate func() ([]byte, *common.Address)
		expPass  bool
	}{
		{
			"unknown method",
			from,
			func() ([]byte, *common.Address) {
				contract := suite.DeployERC20Contract()
				return []byte{1, 2, 3, 4}, &contract
			},
			false,
		},
		{
			"pass",
			from,
			func() ([]byte, *common.Address) {
				contract := suite.DeployERC20Contract()
				data, err := erc20.ABI.Pack("balanceOf", helpers2.GenerateAddress())
				suite.NoError(err)
				return data, &contract
			},
			true,
		},
		{
			"fail empty data",
			from,
			func() ([]byte, *common.Address) {
				contract := suite.DeployERC20Contract()
				return []byte{}, &contract
			},
			false,
		},
		{
			"fail empty sender",
			common.Address{},
			func() ([]byte, *common.Address) {
				contract := suite.DeployERC20Contract()
				return []byte{}, &contract
			},
			false,
		},
		{
			"fail deploy",
			from,
			func() ([]byte, *common.Address) {
				params := suite.app.EvmKeeper.GetParams(suite.ctx)
				params.EnableCreate = false
				suite.app.EvmKeeper.SetParams(suite.ctx, params)
				ctorArgs, err := erc20.ABI.Pack("")
				suite.NoError(err)
				data := append(erc20.Bin, ctorArgs...)
				return data, nil
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			data, contract := tc.malleate()
			_, err := suite.app.EvmKeeper.CallEVMWithData(suite.ctx, tc.from, contract, data, true)
			if tc.expPass {
				suite.NoError(err)
			} else {
				suite.Error(err, err)
			}
		})
	}
}
