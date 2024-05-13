package keeper_test

import (
	"bytes"
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/privval"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/functionx/fx-core/v7/client/jsonrpc"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/staking/types"
)

func (suite *KeeperTestSuite) TestValidatorUpdate() {
	valAddr := sdk.ValAddress(suite.valAccounts[1].GetAddress())

	validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	oldPK, err := validator.ConsPubKey()
	suite.Require().NoError(err)
	consAddr := sdk.ConsAddress(oldPK.Address())

	_, tmAny := suite.GenerateConsKey()
	pk, ok := tmAny.GetCachedValue().(cryptotypes.PubKey)
	suite.Require().True(ok)
	err = suite.app.StakingKeeper.SetConsensusPubKey(suite.ctx, valAddr, pk)
	suite.Require().NoError(err)

	// validator update(process start)
	updates := suite.app.StakingKeeper.ConsensusPubKeyUpdate(suite.ctx)
	suite.Require().Len(updates, 2)
	updates = suite.app.StakingKeeper.ConsensusPubKeyUpdate(suite.ctx)
	suite.Require().Len(updates, 0)

	_, found = suite.app.StakingKeeper.GetConsensusPubKey(suite.ctx, valAddr)
	suite.Require().False(found)

	oldPk, err := suite.app.StakingKeeper.GetConsensusProcess(suite.ctx, valAddr, types.ProcessStart)
	suite.Require().NoError(err)
	suite.Require().Equal(consAddr, sdk.ConsAddress(oldPk.Address()))
	nilPk, err := suite.app.StakingKeeper.GetConsensusProcess(suite.ctx, valAddr, types.ProcessEnd)
	suite.Require().NoError(err)
	suite.Require().Nil(nilPk)

	// consensus process(process end)
	suite.app.StakingKeeper.ConsensusProcess(suite.ctx)

	nilPk, err = suite.app.StakingKeeper.GetConsensusProcess(suite.ctx, valAddr, types.ProcessStart)
	suite.Require().NoError(err)
	suite.Require().Nil(nilPk)

	oldPk, err = suite.app.StakingKeeper.GetConsensusProcess(suite.ctx, valAddr, types.ProcessEnd)
	suite.Require().NoError(err)
	suite.Require().Equal(consAddr, sdk.ConsAddress(oldPk.Address()))

	_, err = suite.app.SlashingKeeper.GetPubkey(suite.ctx, oldPk.Address())
	suite.Require().NoError(err)
	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, sdk.ConsAddress(oldPk.Address()))
	suite.Require().True(found)

	// consensus process(process delete)
	suite.app.StakingKeeper.ConsensusProcess(suite.ctx)

	_, found = suite.app.StakingKeeper.GetConsensusPubKey(suite.ctx, valAddr)
	suite.Require().False(found)

	_, err = suite.app.SlashingKeeper.GetPubkey(suite.ctx, oldPk.Address())
	suite.Require().ErrorContains(err, fmt.Sprintf("address %s not found", sdk.ConsAddress(oldPk.Address()).String()))
	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, sdk.ConsAddress(oldPk.Address()))
	suite.Require().False(found)
}

func (suite *KeeperTestSuite) TestEditPubKeyJail() {
	valAddr := sdk.ValAddress(suite.valAccounts[1].GetAddress())

	validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	oldPk, err := validator.ConsPubKey()
	suite.Require().NoError(err)
	oldConsAddr := sdk.ConsAddress(oldPk.Address())

	// new consensus pubkey
	newPriv, _ := suite.GenerateConsKey()
	newPk := newPriv.PubKey()
	newConsAddr := sdk.ConsAddress(newPk.Address())

	suite.Commit()

	// undelegate smaller min self delegate, and edit validator, can not undelegate all, it will be delete in end block
	shares, err := validator.SharesFromTokensTruncated(validator.MinSelfDelegation.Sub(sdkmath.NewInt(1)))
	suite.Require().NoError(err)
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, sdk.AccAddress(valAddr), valAddr, validator.GetDelegatorShares().Sub(shares))
	suite.Require().NoError(err)
	err = suite.app.StakingKeeper.SetConsensusPubKey(suite.ctx, valAddr, newPk)
	suite.Require().NoError(err)

	// check
	found = suite.app.StakingKeeper.HasConsensusPubKey(suite.ctx, valAddr)
	suite.Require().True(found)
	process, err := suite.app.StakingKeeper.GetConsensusProcess(suite.ctx, valAddr, types.ProcessStart)
	suite.Require().NoError(err)
	suite.Require().Nil(process)
	// next block
	suite.Commit(1) // val update, edit skip to next block
	valUpdates := suite.CommitEndBlock()

	// validator jailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().True(validator.IsJailed())

	// check
	found = suite.app.StakingKeeper.HasConsensusPubKey(suite.ctx, valAddr)
	suite.Require().False(found)

	oldSigningInfo, found := suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().True(found)
	newSigningInfo, found := suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, newConsAddr)
	suite.Require().True(found)
	suite.Require().Equal(oldSigningInfo.IndexOffset, newSigningInfo.IndexOffset)

	process, err = suite.app.StakingKeeper.GetConsensusProcess(suite.ctx, valAddr, types.ProcessStart)
	suite.Require().NoError(err)
	suite.Require().NotNil(process)

	// next block
	suite.CommitBeginBlock(valUpdates)
	valUpdates = suite.CommitEndBlock()

	// check
	oldSigningInfo, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().True(found)
	newSigningInfo, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, newConsAddr)
	suite.Require().True(found)
	suite.Require().Equal(oldSigningInfo.IndexOffset, newSigningInfo.IndexOffset)
	process, err = suite.app.StakingKeeper.GetConsensusProcess(suite.ctx, valAddr, types.ProcessEnd)
	suite.Require().NoError(err)
	suite.Require().NotNil(process)

	// next block
	suite.CommitBeginBlock(valUpdates)
	_ = suite.CommitEndBlock()

	// check
	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().False(found)
	newSigningInfo, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, newConsAddr)
	suite.Require().True(found)
	suite.Require().Equal(oldSigningInfo.IndexOffset, newSigningInfo.IndexOffset)
	found = suite.app.StakingKeeper.HasConsensusProcess(suite.ctx, valAddr)
	suite.Require().False(found)
}

func (suite *KeeperTestSuite) TestEditPubKeyJailAndUnjail() {
	valAddr := sdk.ValAddress(suite.valAccounts[1].GetAddress())

	validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	oldPK, err := validator.ConsPubKey()
	suite.Require().NoError(err)

	delAmt := validator.GetTokens().Sub(sdkmath.NewInt(5)) // min self delegate: 10
	delShares, err := validator.SharesFromTokens(delAmt)
	suite.Require().NoError(err)

	// new consensus pubkey
	newPriv, _ := suite.GenerateConsKey()
	newPk := newPriv.PubKey()

	suite.Commit()

	// undelegate all and edit validator
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, sdk.AccAddress(valAddr), valAddr, delShares)
	suite.Require().NoError(err)
	err = suite.app.StakingKeeper.SetConsensusPubKey(suite.ctx, valAddr, newPk)
	suite.Require().NoError(err)

	suite.Commit()
	valUpdates := suite.CommitEndBlock() // edit

	// validator jailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().True(validator.IsJailed())

	// next block
	suite.CommitBeginBlock(valUpdates)

	// unjail
	_, err = suite.app.StakingKeeper.Delegate(suite.ctx, sdk.AccAddress(valAddr), delAmt, stakingtypes.Unbonded, validator, true)
	suite.Require().NoError(err)
	err = suite.app.SlashingKeeper.Unjail(suite.ctx, valAddr)
	suite.Require().NoError(err)

	// end block
	valUpdates = suite.CommitEndBlock() // process start

	// validator unjailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().False(validator.IsJailed())

	// old signing info exist
	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, sdk.ConsAddress(oldPK.Address()))
	suite.Require().True(found)

	// next block
	suite.CommitBeginBlock(valUpdates)
	_ = suite.CommitEndBlock() // process end

	// old signing info deleted
	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, sdk.ConsAddress(oldPK.Address()))
	suite.Require().False(found)
	// new signing info exist
	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, sdk.ConsAddress(newPk.Address()))
	suite.Require().True(found)
}

func (suite *KeeperTestSuite) TestEditPubKeyJailNextBlock() {
	valAddr := sdk.ValAddress(suite.valAccounts[1].GetAddress())

	validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	oldPK, err := validator.ConsPubKey()
	suite.Require().NoError(err)
	oldConsAddr := sdk.ConsAddress(oldPK.Address())

	// new consensus pubkey
	newPriv, _ := suite.GenerateConsKey()
	newPk := newPriv.PubKey()
	newConsAddr := sdk.ConsAddress(newPk.Address())

	suite.Commit()

	// edit validator
	err = suite.app.StakingKeeper.SetConsensusPubKey(suite.ctx, valAddr, newPk)
	suite.Require().NoError(err)
	valUpdates := suite.CommitEndBlock()

	// check
	found = suite.app.StakingKeeper.HasConsensusPubKey(suite.ctx, valAddr)
	suite.Require().False(found)

	oldSigningInfo, found := suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().True(found)
	newSigningInfo, found := suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, newConsAddr)
	suite.Require().True(found)
	suite.Require().Equal(oldSigningInfo.IndexOffset, newSigningInfo.IndexOffset)

	process, err := suite.app.StakingKeeper.GetConsensusProcess(suite.ctx, valAddr, types.ProcessStart)
	suite.Require().NoError(err)
	suite.Require().NotNil(process)

	// next block begin
	suite.CommitBeginBlock(valUpdates)
	// undelegate smaller min self delegate, can not undelegate all, it will be delete in end block
	shares, err := validator.SharesFromTokensTruncated(validator.MinSelfDelegation.Sub(sdkmath.NewInt(1)))
	suite.Require().NoError(err)
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, sdk.AccAddress(valAddr), valAddr, validator.GetDelegatorShares().Sub(shares))
	suite.Require().NoError(err)
	valUpdates = suite.CommitEndBlock()

	suite.True(suite.CurrentVoteFound(oldPK))
	suite.True(suite.NextVoteFound(newPk))

	// validator jailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().True(validator.IsJailed())

	oldSigningInfo, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().True(found)
	newSigningInfo, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, newConsAddr)
	suite.Require().True(found)
	suite.Require().Equal(oldSigningInfo.IndexOffset, newSigningInfo.IndexOffset)

	process, err = suite.app.StakingKeeper.GetConsensusProcess(suite.ctx, valAddr, types.ProcessEnd)
	suite.Require().NoError(err)
	suite.Require().NotNil(process)

	// next block begin
	suite.CommitBeginBlock(valUpdates)
	suite.CommitEndBlock()

	suite.False(suite.CurrentVoteFound(oldPK))
	suite.True(suite.CurrentVoteFound(newPk))
	suite.False(suite.NextVoteFound(newPk))

	// check
	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().False(found)
	newSigningInfo, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, newConsAddr)
	suite.Require().True(found)
	suite.Require().Equal(oldSigningInfo.IndexOffset+1, newSigningInfo.IndexOffset)

	found = suite.app.StakingKeeper.HasConsensusProcess(suite.ctx, valAddr)
	suite.Require().False(found)
}

func (suite *KeeperTestSuite) TestEditPubKeyUnjailAndJailNextBlock() {
	valAddr := sdk.ValAddress(suite.valAccounts[1].GetAddress())

	validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	oldPk, err := validator.ConsPubKey()
	suite.Require().NoError(err)
	oldConsAddr := sdk.ConsAddress(oldPk.Address())

	delAmt := validator.Tokens.Sub(validator.MinSelfDelegation.Sub(sdkmath.NewInt(1)))
	delShares, err := validator.SharesFromTokens(delAmt)
	suite.Require().NoError(err)

	// new consensus pubkey
	newPriv, _ := suite.GenerateConsKey()
	newPk := newPriv.PubKey()
	suite.Commit(2)

	// undelegate smaller min self delegate, can not undelegate all, it will be delete in end block
	suite.Require().NoError(err)
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, sdk.AccAddress(valAddr), valAddr, delShares)
	suite.Require().NoError(err)

	// validator jailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().True(validator.IsJailed())

	suite.Commit(3)
	valUpdates := suite.CommitEndBlock()

	suite.CommitBeginBlock(valUpdates)
	// edit validator
	err = suite.app.StakingKeeper.SetConsensusPubKey(suite.ctx, valAddr, newPk)
	suite.Require().NoError(err)
	// unjail validator
	validator, _ = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	_, err = suite.app.StakingKeeper.Delegate(suite.ctx, sdk.AccAddress(valAddr), delAmt, stakingtypes.Unbonded, validator, true)
	suite.Require().NoError(err)
	err = suite.app.SlashingKeeper.Unjail(suite.ctx, valAddr)
	suite.Require().NoError(err)
	// end block
	suite.Commit()
	valUpdates = suite.CommitEndBlock() // edit

	// validator unjailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().False(validator.IsJailed())

	consAddr, err := validator.GetConsAddr()
	suite.Require().NoError(err)
	suite.Require().Equal(consAddr, sdk.ConsAddress(newPk.Address()))

	// signing info equal
	oldSigningInfo, found := suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().True(found)
	newSigningInfo, found := suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, sdk.ConsAddress(newPk.Address()))
	suite.Require().True(found)
	suite.Require().Equal(oldSigningInfo.IndexOffset, newSigningInfo.IndexOffset)

	suite.CommitBeginBlock(valUpdates)
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, sdk.AccAddress(valAddr), valAddr, delShares)
	suite.Require().NoError(err)
	valUpdates = suite.CommitEndBlock() // process start

	// validator jailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().True(validator.IsJailed())

	// signing info equal
	oldSigningInfo, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().True(found)
	newSigningInfo, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, sdk.ConsAddress(newPk.Address()))
	suite.Require().True(found)
	suite.Require().Equal(oldSigningInfo.IndexOffset, newSigningInfo.IndexOffset)
	suite.Require().Equal(oldSigningInfo.JailedUntil, newSigningInfo.JailedUntil)

	// next block
	suite.CommitBeginBlock(valUpdates)
	_ = suite.CommitEndBlock() // process end

	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().False(found)
}

func (suite *KeeperTestSuite) TestEditPubKeyJailedPrevBlock() {
	valAddr := sdk.ValAddress(suite.valAccounts[1].GetAddress())

	validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	oldPK, err := validator.ConsPubKey()
	suite.Require().NoError(err)
	oldConsAddr := sdk.ConsAddress(oldPK.Address())

	// new consensus pubkey
	newPriv, _ := suite.GenerateConsKey()
	newPk := newPriv.PubKey()
	newConsAddr := sdk.ConsAddress(newPk.Address())

	suite.Commit()

	// undelegate smaller min self delegate, can not undelegate all, it will be delete in end block
	shares, err := validator.SharesFromTokensTruncated(validator.MinSelfDelegation.Sub(sdkmath.NewInt(1)))
	suite.Require().NoError(err)
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, sdk.AccAddress(valAddr), valAddr, validator.GetDelegatorShares().Sub(shares))
	suite.Require().NoError(err)

	valUpdates := suite.CommitEndBlock()

	// check validator jailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().True(validator.IsJailed())

	suite.CommitBeginBlock(valUpdates)

	// edit validator
	err = suite.app.StakingKeeper.SetConsensusPubKey(suite.ctx, valAddr, newPk)
	suite.Require().NoError(err)

	// check
	found = suite.app.StakingKeeper.HasConsensusPubKey(suite.ctx, valAddr)
	suite.Require().True(found)

	// validator cons address equals to old cons address
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	addr, err := validator.GetConsAddr()
	suite.Require().NoError(err)
	suite.Require().Equal(oldConsAddr, addr)

	// next block
	valUpdates = suite.CommitEndBlock()

	// check consensus pubkey
	found = suite.app.StakingKeeper.HasConsensusPubKey(suite.ctx, valAddr)
	suite.Require().False(found)

	// validator cons address equals to new cons address
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	addr, err = validator.GetConsAddr()
	suite.Require().NoError(err)
	suite.Require().Equal(newConsAddr, addr)

	// signing info equal
	oldSigningInfo, found := suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().True(found)
	newSigningInfo, found := suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, newConsAddr)
	suite.Require().True(found)
	suite.Require().Equal(oldSigningInfo.IndexOffset, newSigningInfo.IndexOffset)

	// next block
	suite.CommitBeginBlock(valUpdates)
	_ = suite.CommitEndBlock()

	// signing info equal
	oldSigningInfo, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().True(found)
	newSigningInfo, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, newConsAddr)
	suite.Require().True(found)
	suite.Require().Equal(oldSigningInfo.IndexOffset, newSigningInfo.IndexOffset)
}

func (suite *KeeperTestSuite) TestEditPubKeyJailedPrevBlockAndUnjail() {
	valAddr := sdk.ValAddress(suite.valAccounts[1].GetAddress())

	validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)

	delAmt := validator.GetTokens()
	// new consensus pubkey
	newPriv, _ := suite.GenerateConsKey()
	newPk := newPriv.PubKey()

	suite.Commit()

	// undelegate smaller min self delegate, can not undelegate all, it will be delete in end block
	shares, err := validator.SharesFromTokensTruncated(validator.MinSelfDelegation.Sub(sdkmath.NewInt(1)))
	suite.Require().NoError(err)
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, sdk.AccAddress(valAddr), valAddr, validator.GetDelegatorShares().Sub(shares))
	suite.Require().NoError(err)
	suite.Commit(2)
	valUpdates := suite.CommitEndBlock()

	// check validator jailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().True(validator.IsJailed())

	suite.CommitBeginBlock(valUpdates)
	// edit validator
	err = suite.app.StakingKeeper.SetConsensusPubKey(suite.ctx, valAddr, newPk)
	suite.Require().NoError(err)
	// unjail validator
	_, err = suite.app.StakingKeeper.Delegate(suite.ctx, sdk.AccAddress(valAddr), delAmt, stakingtypes.Unbonded, validator, true)
	suite.Require().NoError(err)
	err = suite.app.SlashingKeeper.Unjail(suite.ctx, valAddr)
	suite.Require().NoError(err)
	suite.Commit()
	valUpdates = suite.CommitEndBlock() // edit

	// validator unjailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().False(validator.IsJailed())

	suite.CommitBeginBlock(valUpdates)
	valUpdates = suite.CommitEndBlock() // process start

	// signing info equal
	oldSigningInfo, found := suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, sdk.ConsAddress(newPk.Address()))
	suite.Require().True(found)
	newSigningInfo, found := suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, sdk.ConsAddress(newPk.Address()))
	suite.Require().True(found)
	suite.Require().Equal(oldSigningInfo.IndexOffset, newSigningInfo.IndexOffset)

	suite.CommitBeginBlock(valUpdates)
	_ = suite.CommitEndBlock() // process end

	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, sdk.ConsAddress(newPk.Address()))
	suite.Require().True(found)
	newSigningInfo, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, sdk.ConsAddress(newPk.Address()))
	suite.Require().True(found)
	suite.Require().Equal(oldSigningInfo.IndexOffset+1, newSigningInfo.IndexOffset)
}

func (suite *KeeperTestSuite) TestEditPubKeyUnboundValidator() {
	newPriKey := helpers.NewEthPrivKey()
	accAddr := sdk.AccAddress(newPriKey.PubKey().Address())
	newConsPriKey := ed25519.GenPrivKey()
	newConsPubKey := newConsPriKey.PubKey()

	helpers.AddTestAddr(suite.app, suite.ctx, accAddr, sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewInt(10000).Mul(sdk.NewInt(1e18)))))

	// create validator
	valAddr := sdk.ValAddress(accAddr)
	selfDelegateCoin := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1e17))
	des := stakingtypes.Description{Moniker: "test-node"}
	rates := stakingtypes.CommissionRates{
		Rate:          sdk.NewDecWithPrec(1, 2),
		MaxRate:       sdk.NewDecWithPrec(5, 2),
		MaxChangeRate: sdk.NewDecWithPrec(1, 2),
	}
	newValMsg, err := stakingtypes.NewMsgCreateValidator(valAddr, newConsPubKey, selfDelegateCoin, des, rates, sdk.OneInt())
	suite.Require().NoError(err)
	_, err = stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper.Keeper).CreateValidator(suite.ctx, newValMsg)
	suite.Require().NoError(err)

	suite.Commit(3)

	// check validator
	validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().Equal(stakingtypes.Unbonded, validator.GetStatus())

	// edit validator consensus pubkey
	editConsPubKey := ed25519.GenPrivKey().PubKey()
	editPubKeyMsg, err := types.NewMsgEditConsensusPubKey(valAddr, accAddr, editConsPubKey)
	suite.Require().NoError(err)
	_, err = suite.app.StakingKeeper.EditConsensusPubKey(suite.ctx, editPubKeyMsg)
	suite.Require().NoError(err)

	suite.Commit(3)

	// check validator
	editValidator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().Equal(validator, editValidator)
}

func (suite *KeeperTestSuite) TestEditPubKeyDeleteWithoutSigningInfo() {
	initBalance := sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(10000).Mul(sdkmath.NewInt(1e18))))
	accAddr1 := helpers.GenAccAddress()
	helpers.AddTestAddr(suite.app, suite.ctx, accAddr1, initBalance)

	selfDelegateCoin := sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1e17))
	des := stakingtypes.Description{Moniker: "test-node"}
	rates := stakingtypes.CommissionRates{
		Rate:          sdk.NewDecWithPrec(1, 2),
		MaxRate:       sdk.NewDecWithPrec(5, 2),
		MaxChangeRate: sdk.NewDecWithPrec(1, 2),
	}
	oldPubKey := ed25519.GenPrivKey().PubKey()

	// create validator
	newValMsg1, err := stakingtypes.NewMsgCreateValidator(sdk.ValAddress(accAddr1), oldPubKey, selfDelegateCoin, des, rates, sdkmath.OneInt())
	suite.Require().NoError(err)
	_, err = stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper.Keeper).CreateValidator(suite.ctx, newValMsg1)
	suite.Require().NoError(err)

	suite.Commit(3)

	editConsPubKey := ed25519.GenPrivKey().PubKey()
	// 1. edit and delete validator
	editPubKeyMsg, err := types.NewMsgEditConsensusPubKey(sdk.ValAddress(accAddr1), accAddr1, editConsPubKey)
	suite.Require().NoError(err)
	_, err = suite.app.StakingKeeper.EditConsensusPubKey(suite.ctx, editPubKeyMsg)
	suite.Require().NoError(err)
	undelMsg1 := stakingtypes.NewMsgUndelegate(accAddr1, sdk.ValAddress(accAddr1), selfDelegateCoin)
	_, err = stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper.Keeper).Undelegate(suite.ctx, undelMsg1)
	suite.Require().NoError(err)
	suite.CommitEndBlock()

	// block +1, validator delete, process not exist, new signing info not exist, old signing info not exist
	_, found := suite.app.StakingKeeper.GetValidator(suite.ctx, sdk.ValAddress(accAddr1))
	suite.Require().False(found)
	process := suite.app.StakingKeeper.HasConsensusProcess(suite.ctx, sdk.ValAddress(accAddr1))
	suite.Require().False(process)
	found = suite.app.SlashingKeeper.HasValidatorSigningInfo(suite.ctx, sdk.ConsAddress(editConsPubKey.Address()))
	suite.Require().False(found)
	oldConsAddr := sdk.ConsAddress(oldPubKey.Address())
	suite.Require().NoError(err)
	found = suite.app.SlashingKeeper.HasValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().False(found)
}

func (suite *KeeperTestSuite) TestEditPubKeyDelete() {
	valUpdates, accAddr, oldConsAddr := suite.CreateValidatorJailed()
	editConsPubKey := ed25519.GenPrivKey().PubKey()

	suite.CommitBeginBlock(valUpdates)
	// edit and delete validator
	editPubKeyMsg, err := types.NewMsgEditConsensusPubKey(sdk.ValAddress(accAddr), accAddr, editConsPubKey)
	suite.Require().NoError(err)
	_, err = suite.app.StakingKeeper.EditConsensusPubKey(suite.ctx, editPubKeyMsg)
	suite.Require().NoError(err)
	undelMsg1 := stakingtypes.NewMsgUndelegate(accAddr, sdk.ValAddress(accAddr), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1e18).Mul(sdkmath.NewInt(10))))
	_, err = stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper.Keeper).Undelegate(suite.ctx, undelMsg1)
	suite.Require().NoError(err)
	_ = suite.CommitEndBlock()

	// validator delete, new signing info not exist, old signing info exist
	_, found := suite.app.StakingKeeper.GetValidator(suite.ctx, sdk.ValAddress(accAddr))
	suite.Require().False(found)
	found = suite.app.SlashingKeeper.HasValidatorSigningInfo(suite.ctx, sdk.ConsAddress(editConsPubKey.Address()))
	suite.Require().False(found)
	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().True(found)
}

func (suite *KeeperTestSuite) TestEditPubKeyDeleteNextBlock() {
	valUpdates, accAddr, oldConsAddr := suite.CreateValidatorJailed()
	editConsPubKey := ed25519.GenPrivKey().PubKey()

	// edit pubkey
	suite.CommitBeginBlock(valUpdates)
	editPubKeyMsg, err := types.NewMsgEditConsensusPubKey(sdk.ValAddress(accAddr), accAddr, editConsPubKey)
	suite.Require().NoError(err)
	_, err = suite.app.StakingKeeper.EditConsensusPubKey(suite.ctx, editPubKeyMsg)
	suite.Require().NoError(err)
	valUpdates = suite.CommitEndBlock()

	process := suite.app.StakingKeeper.HasConsensusProcess(suite.ctx, sdk.ValAddress(accAddr))
	suite.Require().True(process)

	// next block, validator delete
	suite.CommitBeginBlock(valUpdates)
	undelMsg := stakingtypes.NewMsgUndelegate(accAddr, sdk.ValAddress(accAddr), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1e18).Mul(sdkmath.NewInt(10))))
	_, err = stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper.Keeper).Undelegate(suite.ctx, undelMsg)
	suite.Require().NoError(err)
	valUpdates = suite.CommitEndBlock()

	// validator delete and old and new signing info exist
	_, found := suite.app.StakingKeeper.GetValidator(suite.ctx, sdk.ValAddress(accAddr))
	suite.Require().False(found)
	process = suite.app.StakingKeeper.HasConsensusProcess(suite.ctx, sdk.ValAddress(accAddr))
	suite.Require().True(process)
	found = suite.app.SlashingKeeper.HasValidatorSigningInfo(suite.ctx, sdk.ConsAddress(editConsPubKey.Address()))
	suite.Require().True(found)
	found = suite.app.SlashingKeeper.HasValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().True(found)

	// block +1, old signing info deleted
	suite.CommitBeginBlock(valUpdates)
	_ = suite.CommitEndBlock()

	// old signing info deleted, process not exist
	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().False(found)
	process = suite.app.StakingKeeper.HasConsensusProcess(suite.ctx, sdk.ValAddress(accAddr))
	suite.Require().False(process)
}

func (suite *KeeperTestSuite) TestEditPubKeyDeleteNextNextBlock() {
	valUpdates, accAddr, oldConsAddr := suite.CreateValidatorJailed()
	editConsPubKey := ed25519.GenPrivKey().PubKey()

	// edit pubkey
	suite.CommitBeginBlock(valUpdates)
	editPubKeyMsg, err := types.NewMsgEditConsensusPubKey(sdk.ValAddress(accAddr), accAddr, editConsPubKey)
	suite.Require().NoError(err)
	_, err = suite.app.StakingKeeper.EditConsensusPubKey(suite.ctx, editPubKeyMsg)
	suite.Require().NoError(err)
	valUpdates = suite.CommitEndBlock()

	process := suite.app.StakingKeeper.HasConsensusProcess(suite.ctx, sdk.ValAddress(accAddr))
	suite.Require().True(process)

	suite.CommitBeginBlock(valUpdates)
	valUpdates = suite.CommitEndBlock()

	process = suite.app.StakingKeeper.HasConsensusProcess(suite.ctx, sdk.ValAddress(accAddr))
	suite.Require().True(process)
	// validator exist, old and new signing info exist, new signing info exit
	_, found := suite.app.StakingKeeper.GetValidator(suite.ctx, sdk.ValAddress(accAddr))
	suite.Require().True(found)
	found = suite.app.SlashingKeeper.HasValidatorSigningInfo(suite.ctx, sdk.ConsAddress(editConsPubKey.Address()))
	suite.Require().True(found)
	found = suite.app.SlashingKeeper.HasValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().True(found)

	suite.CommitBeginBlock(valUpdates)
	// validator delete
	undelMsg := stakingtypes.NewMsgUndelegate(accAddr, sdk.ValAddress(accAddr), sdk.NewCoin(fxtypes.DefaultDenom, sdkmath.NewInt(1e18).Mul(sdkmath.NewInt(10))))
	_, err = stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper.Keeper).Undelegate(suite.ctx, undelMsg)
	suite.Require().NoError(err)
	_ = suite.CommitEndBlock()

	// validator delete, new signing info exist, old signing info not exist
	_, found = suite.app.StakingKeeper.GetValidator(suite.ctx, sdk.ValAddress(accAddr))
	suite.Require().False(found)
	found = suite.app.SlashingKeeper.HasValidatorSigningInfo(suite.ctx, sdk.ConsAddress(editConsPubKey.Address()))
	suite.Require().True(found)
	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().False(found)
}

func TestValidatorUpdateEvidence(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	keyFilePath := ""
	height := int64(100)

	rpc := jsonrpc.NewNodeRPC(jsonrpc.NewClient("http://localhost:26657"))
	key := privval.LoadFilePVEmptyState(keyFilePath, "")

	// val power
	valResp, err := rpc.Validators(height, 1, 100)
	require.NoError(t, err)
	valPowers := make([]*tmtypes.Validator, 0, len(valResp.Validators))
	for _, val := range valResp.Validators {
		valPowers = append(valPowers, &tmtypes.Validator{
			Address:          val.Address,
			PubKey:           val.PubKey,
			VotingPower:      val.VotingPower,
			ProposerPriority: val.ProposerPriority,
		})
	}

	// val commit
	commitResp, err := rpc.Commit(height)
	require.NoError(t, err)
	commit := commitResp.Commit
	cs := tmtypes.CommitSig{}
	for _, sig := range commitResp.Commit.Signatures {
		if !bytes.Equal(sig.ValidatorAddress, key.GetAddress()) {
			continue
		}
		if sig.BlockIDFlag == tmtypes.BlockIDFlagNil {
			t.Fatalf("block %d validator BlockIDFlag is nil, change block", height)
		}
		cs = sig
		break
	}

	voteA := tmtypes.Vote{
		ValidatorAddress: cs.ValidatorAddress,
		Height:           commit.Height,
		Round:            commit.Round,
		Timestamp:        cs.Timestamp,
		Type:             tmproto.SignedMsgType(commit.Type()),
		BlockID:          cs.BlockID(commit.BlockID),
		Signature:        cs.Signature,
	}

	voteB := voteA
	voteB.BlockID.Hash = tmhash.Sum([]byte("test"))
	voteB.Signature, err = key.Key.PrivKey.Sign(tmtypes.VoteSignBytes(commitResp.ChainID, voteB.ToProto()))
	require.NoError(t, err)

	// evidence
	evidence := tmtypes.NewDuplicateVoteEvidence(&voteA, &voteB, commitResp.Header.Time, tmtypes.NewValidatorSet(valPowers))
	_, err = rpc.BroadcastEvidence(evidence)
	require.NoError(t, err)
}
