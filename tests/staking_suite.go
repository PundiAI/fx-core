package tests

import (
	"context"
	"math/big"
	"time"

	sdkmath "cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/functionx/fx-core/v3/client"
	"github.com/functionx/fx-core/v3/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/types/contract"
	fxstakingtypes "github.com/functionx/fx-core/v3/x/staking/types"
)

type StakingTestSuite struct {
	*TestSuite
	privKey cryptotypes.PrivKey
}

func NewStakingTestSuite(ts *TestSuite) StakingTestSuite {
	return StakingTestSuite{
		TestSuite: ts,
		privKey:   helpers.NewEthPrivKey(),
	}
}

func (suite *StakingTestSuite) SetupSuite() {
	suite.TestSuite.SetupSuite()

	// transfer to eth private key
	suite.Send(suite.AccAddress(), suite.NewCoin(sdkmath.NewInt(100).MulRaw(1e18)))
}

func (suite *StakingTestSuite) AccAddress() sdk.AccAddress {
	return sdk.AccAddress(suite.privKey.PubKey().Address())
}

func (suite *StakingTestSuite) EthClient() *ethclient.Client {
	return suite.GetFirstValidator().JSONRPCClient
}

func (suite *StakingTestSuite) Balance(addr common.Address) *big.Int {
	at, err := suite.EthClient().BalanceAt(suite.ctx, addr, nil)
	suite.Require().NoError(err)
	return at
}

func (suite *StakingTestSuite) BalanceOf(contractAddr, address common.Address) *big.Int {
	caller, err := contract.NewLPToken(contractAddr, suite.EthClient())
	suite.NoError(err)
	balance, err := caller.BalanceOf(nil, address)
	suite.NoError(err)
	return balance
}

func (suite *StakingTestSuite) Validators() stakingtypes.Validators {
	resp, err := suite.GRPCClient().StakingQuery().Validators(suite.ctx, &stakingtypes.QueryValidatorsRequest{
		Status: stakingtypes.Bonded.String(),
	})
	suite.Require().NoError(err)

	return resp.Validators
}

func (suite *StakingTestSuite) GetDelegation(delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, sdk.Coin) {
	resp, err := suite.GRPCClient().StakingQuery().Delegation(suite.ctx, &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: valAddr.String(),
	})
	suite.Require().NoError(err)
	return resp.DelegationResponse.Delegation, resp.DelegationResponse.Balance
}

func (suite *StakingTestSuite) Validator(val sdk.ValAddress) stakingtypes.Validator {
	resp, err := suite.GRPCClient().StakingQuery().Validator(suite.ctx, &stakingtypes.QueryValidatorRequest{ValidatorAddr: val.String()})
	suite.Require().NoError(err)

	return resp.Validator
}

func (suite *StakingTestSuite) ValidatorLPToken(val sdk.ValAddress) common.Address {
	resp, err := suite.GRPCClient().FXStakingQuery().ValidatorLPToken(
		suite.ctx, &fxstakingtypes.QueryValidatorLPTokenRequest{ValidatorAddr: val.String()})
	suite.Require().NoError(err)
	return common.HexToAddress(resp.LpToken.Address)
}

func (suite *StakingTestSuite) LPTokenValidator() {}

func (suite *StakingTestSuite) SendTransaction(tx *ethtypes.Transaction) {
	err := suite.EthClient().SendTransaction(suite.ctx, tx)
	suite.Require().NoError(err)

	ctx, cancel := context.WithTimeout(suite.ctx, 5*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, suite.EthClient(), tx)
	suite.Require().NoError(err)
	suite.Require().Equal(receipt.Status, ethtypes.ReceiptStatusSuccessful)
}

func (suite *StakingTestSuite) LPTokenApprove(privateKey cryptotypes.PrivKey, token, recipient common.Address, amount *big.Int) *ethtypes.Transaction {
	pack, err := fxtypes.GetLPToken().ABI.Pack("approve", recipient, amount)
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)
	return ethTx
}

func (suite *StakingTestSuite) LPTokenTransfer(privateKey cryptotypes.PrivKey, token, recipient common.Address, amount *big.Int) *ethtypes.Transaction {
	pack, err := fxtypes.GetLPToken().ABI.Pack("transfer", recipient, amount)
	suite.Require().NoError(err)
	ethTx, err := client.BuildEthTransaction(suite.ctx, suite.EthClient(), privateKey, &token, nil, pack)
	suite.Require().NoError(err)
	suite.SendTransaction(ethTx)
	return ethTx
}
