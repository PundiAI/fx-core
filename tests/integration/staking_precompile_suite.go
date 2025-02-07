package integration

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/pundiai/fx-core/v8/contract"
	stakingprecompile "github.com/pundiai/fx-core/v8/precompiles/staking"
	"github.com/pundiai/fx-core/v8/testutil/helpers"
)

type StakingPrecompileSuite struct {
	*EthSuite

	staking *contract.IStaking
	signer  *helpers.Signer
}

func NewStakingSuite(suite *EthSuite, contractAddr common.Address, signer *helpers.Signer) *StakingPrecompileSuite {
	staking, err := contract.NewIStaking(contractAddr, suite.ethCli)
	suite.Require().NoError(err)
	return &StakingPrecompileSuite{
		EthSuite: suite,
		staking:  staking,
		signer:   signer,
	}
}

func (suite *StakingPrecompileSuite) WithSigner(signer *helpers.Signer) *StakingPrecompileSuite {
	return &StakingPrecompileSuite{
		EthSuite: suite.EthSuite,
		staking:  suite.staking,
		signer:   signer,
	}
}

func (suite *StakingPrecompileSuite) DelegateV2(valAddr string, delAmount *big.Int, value ...*big.Int) *ethtypes.Receipt {
	opts := suite.TransactOpts(suite.signer)
	if len(value) > 0 && value[0] != nil {
		opts.Value = value[0]
	}
	ethTx, err := suite.staking.DelegateV2(opts, valAddr, delAmount)
	suite.Require().NoError(err)
	return suite.WaitMined(ethTx)
}

func (suite *StakingPrecompileSuite) RedelegateV2(valSrc, valDst string, amount *big.Int) *ethtypes.Receipt {
	ethTx, err := suite.staking.RedelegateV2(suite.TransactOpts(suite.signer), valSrc, valDst, amount)
	suite.Require().NoError(err)
	return suite.WaitMined(ethTx)
}

func (suite *StakingPrecompileSuite) UnDelegateV2(valAddr string, amount *big.Int) *ethtypes.Receipt {
	ethTx, err := suite.staking.UndelegateV2(suite.TransactOpts(suite.signer), valAddr, amount)
	suite.Require().NoError(err)
	return suite.WaitMined(ethTx)
}

func (suite *StakingPrecompileSuite) Withdraw(valAddr string) *ethtypes.Receipt {
	ethTx, err := suite.staking.Withdraw(suite.TransactOpts(suite.signer), valAddr)
	suite.Require().NoError(err)
	return suite.WaitMined(ethTx)
}

func (suite *StakingPrecompileSuite) Delegation(valAddr string, delAddr common.Address) (*big.Int, *big.Int) {
	delegation, err := suite.staking.Delegation(nil, valAddr, delAddr)
	suite.Require().NoError(err)
	return delegation.Shares, delegation.DelegateAmount
}

func (suite *StakingPrecompileSuite) Rewards(valAddr string, delAddr common.Address) *big.Int {
	rewards, err := suite.staking.DelegationRewards(nil, valAddr, delAddr)
	suite.Require().NoError(err)
	return rewards
}

func (suite *StakingPrecompileSuite) TransferShares(valAddr string, receipt common.Address, shares *big.Int) *ethtypes.Receipt {
	ethTx, err := suite.staking.TransferShares(suite.TransactOpts(suite.signer), valAddr, receipt, shares)
	suite.Require().NoError(err)
	return suite.WaitMined(ethTx)
}

func (suite *StakingPrecompileSuite) TransferFromShares(valAddr string, from, receipt common.Address, shares *big.Int) *ethtypes.Receipt {
	ethTx, err := suite.staking.TransferFromShares(suite.TransactOpts(suite.signer), valAddr, from, receipt, shares)
	suite.Require().NoError(err)
	return suite.WaitMined(ethTx)
}

func (suite *StakingPrecompileSuite) ApproveShares(valAddr string, spender common.Address, shares *big.Int) *ethtypes.Receipt {
	ethTx, err := suite.staking.ApproveShares(suite.TransactOpts(suite.signer), valAddr, spender, shares)
	suite.Require().NoError(err)
	return suite.WaitMined(ethTx)
}

func (suite *StakingPrecompileSuite) AllowanceShares(valAddr string, owner, spender common.Address) *big.Int {
	allowance, err := suite.staking.AllowanceShares(nil, valAddr, owner, spender)
	suite.Require().NoError(err)
	return allowance
}

func (suite *StakingPrecompileSuite) LogReward(logs []*ethtypes.Log, valAddr string, addr common.Address) *big.Int {
	withdrawABI := stakingprecompile.NewWithdrawABI()
	for _, log := range logs {
		if log.Address.String() == contract.StakingAddress &&
			log.Topics[0] == withdrawABI.Event.ID &&
			log.Topics[1] == addr.Hash() {
			unpack, err := withdrawABI.Event.Inputs.NonIndexed().Unpack(log.Data)
			suite.Require().NoError(err)
			suite.Require().Equal(unpack[0].(string), valAddr)
			return unpack[1].(*big.Int)
		}
	}
	return big.NewInt(0)
}
