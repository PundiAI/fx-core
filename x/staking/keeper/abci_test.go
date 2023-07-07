package keeper_test

import (
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (suite *KeeperTestSuite) TestABCIValidatorUpdates() {
	valAddr := sdk.ValAddress(suite.valAccounts[1].GetAddress())
	_, tmAny := suite.GenerateConsKey()
	pk, ok := tmAny.GetCachedValue().(cryptotypes.PubKey)
	suite.Require().True(ok)
	err := suite.app.StakingKeeper.SetValidatorNewConsensusPubKey(suite.ctx, valAddr, pk)
	suite.Require().NoError(err)

	updates := suite.app.StakingKeeper.ABCIValidatorUpdate(suite.ctx, []abci.ValidatorUpdate{})

	suite.Require().Len(updates, 2)
}

func (suite *KeeperTestSuite) TestUpdateValidatorConsensus() {
	valAddr := sdk.ValAddress(suite.valAccounts[1].GetAddress())
	validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	consAddr, err := validator.GetConsAddr()
	suite.Require().NoError(err)
	suite.app.SlashingKeeper.SetValidatorSigningInfo(suite.ctx, consAddr, slashingtypes.NewValidatorSigningInfo(consAddr, 0, 0, time.UnixMicro(0), false, 0))

	_, tmAny := suite.GenerateConsKey()
	pk, ok := tmAny.GetCachedValue().(cryptotypes.PubKey)
	suite.Require().True(ok)
	err = suite.app.StakingKeeper.SetValidatorNewConsensusPubKey(suite.ctx, valAddr, pk)
	suite.Require().NoError(err)
	suite.app.StakingKeeper.SetValidatorOldConsensusAddr(suite.ctx, valAddr, consAddr)

	suite.app.StakingKeeper.UpdateValidatorConsensusKey(suite.ctx)

	_, found = suite.app.StakingKeeper.GetValidatorNewConsensusPubKey(suite.ctx, valAddr)
	suite.Require().False(found)
	valNew, found := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	valPk, err := valNew.ConsPubKey()
	suite.Require().NoError(err)
	suite.Require().Equal(pk.Address(), valPk.Address())
}

func (suite *KeeperTestSuite) TestRemoveValidatorConsensusKey() {
	valAddr := sdk.ValAddress(suite.valAccounts[1].GetAddress())
	validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	consAddr, err := validator.GetConsAddr()
	suite.Require().NoError(err)
	suite.app.StakingKeeper.SetValidatorDelConsensusAddr(suite.ctx, valAddr, consAddr)

	suite.app.StakingKeeper.RemoveValidatorConsensusKey(suite.ctx)

	_, found = suite.app.StakingKeeper.GetValidatorByConsAddr(suite.ctx, consAddr)
	suite.Require().False(found)
}
