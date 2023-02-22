package keeper_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/staking/keeper"
	"github.com/functionx/fx-core/v3/x/staking/types"
)

func (suite *KeeperTestSuite) TestLPTokenTransferHandler() {
	validator, found := suite.app.StakingKeeper.GetValidatorByConsAddr(suite.ctx, suite.ctx.BlockHeader().ProposerAddress)
	suite.Require().True(found)

	del1, val, lpToken, _, share := suite.RandDelegates(validator)
	del2 := suite.RandSigner()

	delegate, found := suite.app.StakingKeeper.GetDelegation(suite.ctx, del1.AccAddress(), val)
	suite.Require().True(found)
	suite.Require().Equal(share, delegate.Shares)
	delegate, found = suite.app.StakingKeeper.GetDelegation(suite.ctx, del2.AccAddress(), val)
	suite.Require().False(found)

	h := keeper.NewLPTokenTransferHandler(suite.app.StakingKeeper)

	// 1. transfer lp token
	topics := []common.Hash{h.EventID(), common.BytesToHash(del1.Address().Bytes()), common.BytesToHash(del2.Address().Bytes())}
	data, err := fxtypes.GetLPToken().ABI.Events[types.LPTokenTransferEventName].Inputs.NonIndexed().Pack(share.BigInt())
	suite.Require().NoError(err)
	log := &ethtypes.Log{Address: lpToken, Topics: topics, Data: data}

	err = h.Handle(suite.ctx, nil, log)
	suite.Require().NoError(err)

	delegate, found = suite.app.StakingKeeper.GetDelegation(suite.ctx, del1.AccAddress(), val)
	suite.Require().False(found)
	delegate, found = suite.app.StakingKeeper.GetDelegation(suite.ctx, del2.AccAddress(), val)
	suite.Require().True(found)
	suite.Require().Equal(share, delegate.Shares)

	// 2. transfer zero lp token
	topics = []common.Hash{h.EventID(), common.BytesToHash(del1.Address().Bytes()), common.BytesToHash(del2.Address().Bytes())}
	data, err = fxtypes.GetLPToken().ABI.Events[types.LPTokenTransferEventName].Inputs.NonIndexed().Pack(big.NewInt(0))
	suite.Require().NoError(err)
	log = &ethtypes.Log{Address: lpToken, Topics: topics, Data: data}

	err = h.Handle(suite.ctx, nil, log)
	suite.NoError(err)

	// 3. failed to transfer lp token
	topics = []common.Hash{h.EventID(), common.BytesToHash(del1.Address().Bytes()), common.BytesToHash(del2.Address().Bytes())}
	data, err = fxtypes.GetLPToken().ABI.Events[types.LPTokenTransferEventName].Inputs.NonIndexed().Pack(big.NewInt(1))
	suite.Require().NoError(err)
	log = &ethtypes.Log{Address: lpToken, Topics: topics, Data: data}

	err = h.Handle(suite.ctx, nil, log)
	suite.Error(err)
	suite.Equal(err.Error(), "delegator does not contain delegation")
}
