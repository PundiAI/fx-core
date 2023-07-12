package keeper_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
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
