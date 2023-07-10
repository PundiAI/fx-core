package keeper_test

import (
	"fmt"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/functionx/fx-core/v5/x/staking/types"
)

func (suite *KeeperTestSuite) TestValidatorUpdate() {
	valAddr := sdk.ValAddress(suite.valAccounts[1].GetAddress())

	validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	consAddr, err := validator.GetConsAddr()
	suite.Require().NoError(err)
	suite.app.SlashingKeeper.SetValidatorSigningInfo(suite.ctx, consAddr, slashingtypes.NewValidatorSigningInfo(consAddr, 0, 0, time.UnixMicro(0), false, 0))

	_, tmAny := suite.GenerateConsKey()
	pk, ok := tmAny.GetCachedValue().(cryptotypes.PubKey)
	suite.Require().True(ok)
	err = suite.app.StakingKeeper.SetConsensusPubKey(suite.ctx, valAddr, pk)
	suite.Require().NoError(err)

	found = suite.app.StakingKeeper.HasConsensusProcess(suite.ctx, valAddr)
	suite.Require().False(found)

	// validator update(process start)
	updates := suite.app.StakingKeeper.ValidatorUpdate(suite.ctx, []abci.ValidatorUpdate{})
	suite.Require().Len(updates, 2)
	updates = suite.app.StakingKeeper.ValidatorUpdate(suite.ctx, []abci.ValidatorUpdate{})
	suite.Require().Len(updates, 0)

	_, found = suite.app.StakingKeeper.GetConsensusPubKey(suite.ctx, valAddr)
	suite.Require().False(found)

	oldConsAddr, found := suite.app.StakingKeeper.GetConsensusProcess(suite.ctx, valAddr, types.ProcessStart)
	suite.Require().True(found)
	suite.Require().Equal(consAddr, oldConsAddr)
	_, found = suite.app.StakingKeeper.GetConsensusProcess(suite.ctx, valAddr, types.ProcessEnd)
	suite.Require().False(found)

	// consensus process(process end)
	suite.app.StakingKeeper.ConsensusProcess(suite.ctx)

	_, found = suite.app.StakingKeeper.GetConsensusProcess(suite.ctx, valAddr, types.ProcessStart)
	suite.Require().False(found)
	oldConsAddr, found = suite.app.StakingKeeper.GetConsensusProcess(suite.ctx, valAddr, types.ProcessEnd)
	suite.Require().True(found)
	suite.Require().Equal(consAddr, oldConsAddr)

	_, err = suite.app.SlashingKeeper.GetPubkey(suite.ctx, oldConsAddr.Bytes())
	suite.Require().NoError(err)
	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().True(found)

	// consensus process(process delete)
	suite.app.StakingKeeper.ConsensusProcess(suite.ctx)

	_, found = suite.app.StakingKeeper.GetConsensusPubKey(suite.ctx, valAddr)
	suite.Require().False(found)

	_, err = suite.app.SlashingKeeper.GetPubkey(suite.ctx, oldConsAddr.Bytes())
	suite.Require().ErrorContains(err, fmt.Sprintf("address %s not found", oldConsAddr.String()))
	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().False(found)
}
