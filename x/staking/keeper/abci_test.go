package keeper_test

import (
	"bytes"
	"fmt"
	"testing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/privval"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/functionx/fx-core/v5/client/jsonrpc"
	"github.com/functionx/fx-core/v5/testutil/helpers"
	"github.com/functionx/fx-core/v5/x/staking/types"
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

	valUpdates := make([]abci.ValidatorUpdate, 0)
	pkPowerUpdate := make(map[string]int64)
	lastVote := make(map[string]bool)
	for _, info := range suite.currentVoteInfo {
		lastVote[string(info.Validator.Address)] = true
	}

	// validator update(process start)
	updates := suite.app.StakingKeeper.ValidatorUpdate(suite.ctx, valUpdates, pkPowerUpdate)
	suite.Require().Len(updates, 2)
	updates = suite.app.StakingKeeper.ValidatorUpdate(suite.ctx, valUpdates, pkPowerUpdate)
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

	// undelegate all and edit validator
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, sdk.AccAddress(valAddr), valAddr, validator.GetDelegatorShares())
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

	suite.True(suite.CurrentVoteFound(oldPk))
	suite.False(suite.NextVoteFound(newPk))

	// check
	oldSigningInfo, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().True(found)
	newSigningInfo, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, newConsAddr)
	suite.Require().True(found)
	suite.Require().Equal(oldSigningInfo.IndexOffset, newSigningInfo.IndexOffset)

	// next block
	suite.CommitBeginBlock(valUpdates)
	_ = suite.CommitEndBlock()

	suite.False(suite.CurrentVoteFound(oldPk))
	suite.False(suite.NextVoteFound(newPk))

	// check
	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, oldConsAddr)
	suite.Require().False(found)
	newSigningInfo, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, newConsAddr)
	suite.Require().True(found)
	suite.Require().Equal(oldSigningInfo.IndexOffset+1, newSigningInfo.IndexOffset)
}

func (suite *KeeperTestSuite) TestEditPubKeyJailAndUnjail() {
	valAddr := sdk.ValAddress(suite.valAccounts[1].GetAddress())

	validator, found := suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	oldPK, err := validator.ConsPubKey()
	suite.Require().NoError(err)

	delAmt := validator.GetTokens()
	delShares := validator.GetDelegatorShares()

	// new consensus pubkey
	newPriv, _ := suite.GenerateConsKey()
	newPk := newPriv.PubKey()

	suite.Commit()

	// undelegate all and edit validator
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, sdk.AccAddress(valAddr), valAddr, delShares)
	suite.Require().NoError(err)
	err = suite.app.StakingKeeper.SetConsensusPubKey(suite.ctx, valAddr, newPk)
	suite.Require().NoError(err)

	// end block
	valUpdates := suite.CommitEndBlock()

	// validator jailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().True(validator.IsJailed())

	// next block
	suite.CommitBeginBlock(valUpdates)

	suite.True(suite.CurrentVoteFound(oldPK))
	suite.False(suite.NextVoteFound(newPk))

	// unjail
	_, err = suite.app.StakingKeeper.Delegate(suite.ctx, sdk.AccAddress(valAddr), delAmt, stakingtypes.Unbonded, validator, true)
	suite.Require().NoError(err)
	err = suite.app.SlashingKeeper.Unjail(suite.ctx, valAddr)
	suite.Require().NoError(err)

	// end block
	valUpdates = suite.CommitEndBlock()

	suite.True(suite.CurrentVoteFound(oldPK))
	suite.False(suite.CurrentVoteFound(newPk))
	suite.False(suite.NextVoteFound(newPk))

	// validator unjailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().False(validator.IsJailed())

	// old signing info exist
	_, found = suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, sdk.ConsAddress(oldPK.Address()))
	suite.Require().True(found)

	// next block
	suite.CommitBeginBlock(valUpdates)
	_ = suite.CommitEndBlock()

	suite.False(suite.CurrentVoteFound(newPk))
	suite.True(suite.NextVoteFound(newPk))

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
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, sdk.AccAddress(valAddr), valAddr, validator.GetDelegatorShares())
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

	delAmt := validator.GetTokens()
	delShares := validator.GetDelegatorShares()

	// new consensus pubkey
	newPriv, _ := suite.GenerateConsKey()
	newPk := newPriv.PubKey()

	suite.Commit()

	// jailed by undelegate all
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, sdk.AccAddress(valAddr), valAddr, delShares)
	suite.Require().NoError(err)

	// validator jailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().True(validator.IsJailed())

	suite.Commit()
	suite.Commit()
	valUpdates := suite.CommitEndBlock()

	suite.False(suite.CurrentVoteFound(oldPk))
	suite.False(suite.NextVoteFound(oldPk))

	suite.CommitBeginBlock(valUpdates)
	// edit validator
	err = suite.app.StakingKeeper.SetConsensusPubKey(suite.ctx, valAddr, newPk)
	suite.Require().NoError(err)
	// unjail validator
	_, err = suite.app.StakingKeeper.Delegate(suite.ctx, sdk.AccAddress(valAddr), delAmt, stakingtypes.Unbonded, validator, true)
	suite.Require().NoError(err)
	err = suite.app.SlashingKeeper.Unjail(suite.ctx, valAddr)
	suite.Require().NoError(err)
	// end block
	valUpdates = suite.CommitEndBlock()

	// validator unjailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().False(validator.IsJailed())

	consAddr, err := validator.GetConsAddr()
	suite.Require().NoError(err)
	suite.Require().Equal(consAddr, sdk.ConsAddress(newPk.Address()))

	suite.CommitBeginBlock(valUpdates)
	// validator jailed
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, sdk.AccAddress(valAddr), valAddr, delShares)
	suite.Require().NoError(err)
	_ = suite.CommitEndBlock()

	suite.True(suite.NextVoteFound(newPk))

	// validator jailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().True(validator.IsJailed())
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

	// jailed by undelegate all
	_, err = suite.app.StakingKeeper.Undelegate(suite.ctx, sdk.AccAddress(valAddr), valAddr, validator.GetDelegatorShares())
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
	delShares := validator.GetDelegatorShares()

	// new consensus pubkey
	newPriv, _ := suite.GenerateConsKey()
	newPk := newPriv.PubKey()

	suite.Commit()

	// jailed by undelegate all
	_, err := suite.app.StakingKeeper.Undelegate(suite.ctx, sdk.AccAddress(valAddr), valAddr, delShares)
	suite.Require().NoError(err)
	valUpdates := suite.CommitEndBlock()

	suite.CommitBeginBlock(valUpdates)
	valUpdates = suite.CommitEndBlock()
	suite.CommitBeginBlock(valUpdates)
	valUpdates = suite.CommitEndBlock()

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
	valUpdates = suite.CommitEndBlock()

	// validator unjailed
	validator, found = suite.app.StakingKeeper.GetValidator(suite.ctx, valAddr)
	suite.Require().True(found)
	suite.Require().False(validator.IsJailed())

	suite.False(suite.CurrentVoteFound(newPk))
	suite.False(suite.NextVoteFound(newPk))

	suite.CommitBeginBlock(valUpdates)
	_ = suite.CommitEndBlock()

	suite.False(suite.CurrentVoteFound(newPk))
	suite.True(suite.NextVoteFound(newPk))
}

func TestValidatorUpdateEvidence(t *testing.T) {
	helpers.SkipTest(t, "Skipping local test:", t.Name())

	keyFilePath := ""
	height := int64(100)

	rpc := jsonrpc.NewNodeRPC(jsonrpc.NewClient("http://localhost:26657"))
	key := privval.LoadFilePVEmptyState(keyFilePath, "")
	t.Log("address:", key.GetAddress().String())

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
