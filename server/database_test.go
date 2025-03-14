package server_test

import (
	"os"
	"testing"

	tmcfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/crypto"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	tmtypes "github.com/cometbft/cometbft/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/stretchr/testify/suite"

	"github.com/pundiai/fx-core/v8/server"
	"github.com/pundiai/fx-core/v8/testutil"
)

type DatabaseTestSuite struct {
	suite.Suite
	config   *tmcfg.Config
	database *server.Database
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (suite *DatabaseTestSuite) SetupTest() {
	newCfg := testutil.ResetTestRootWithChainID("blockchain_database_test", "fxcore")

	database, err := server.NewDatabase(newCfg)
	suite.Require().NoError(err)

	_, err = database.StateStore().LoadFromDBOrGenesisFile(newCfg.GenesisFile())
	suite.Require().NoError(err)
	suite.config = newCfg
	suite.database = database
}

func (suite *DatabaseTestSuite) TearDownSuite() {
	defer os.RemoveAll(suite.config.RootDir)
	suite.database.Close()
}

func (suite *DatabaseTestSuite) TestGetChainID() {
	_, err := suite.database.GetChainId()
	suite.Require().Error(err, "not found chain id")

	suite.newBlock(1)

	chainID, err := suite.database.GetChainId()
	suite.Require().NoError(err)
	suite.Equal("fxcore", chainID)
}

func (suite *DatabaseTestSuite) TestGetBlockHeight() {
	height, err := suite.database.GetBlockHeight()
	suite.Require().NoError(err)
	suite.Equal(int64(0), height)

	suite.newBlock(5)

	height, err = suite.database.GetBlockHeight()
	suite.Require().NoError(err)
	suite.Equal(int64(5), height)
}

func (suite *DatabaseTestSuite) TestGetSyncing() {
	syncing, err := suite.database.GetSyncing()
	suite.Require().NoError(err)
	suite.True(syncing)
}

func (suite *DatabaseTestSuite) TestGetNodeInfo() {
	nodeInfo, err := suite.database.GetNodeInfo()
	suite.Require().NoError(err)
	suite.NotNil(nodeInfo)
}

func (suite *DatabaseTestSuite) TestCurrentPlan() {
	plan, err := suite.database.CurrentPlan()
	suite.Require().NoError(err)
	suite.Nil(plan)
}

func (suite *DatabaseTestSuite) TestGetConsensusValidators() {
	validators, err := suite.database.GetConsensusValidators()
	suite.Require().NoError(err)
	suite.NotNil(validators)
	suite.Len(validators, 1)
}

func (suite *DatabaseTestSuite) newBlock(height int64) {
	state, err := suite.database.StateStore().LoadFromDBOrGenesisFile(suite.config.GenesisFile())
	suite.Require().NoError(err)

	for h := int64(1); h <= height; h++ {
		block := state.MakeBlock(h, nil, new(tmtypes.Commit), nil, state.Validators.GetProposer().Address)
		partSet, err := block.MakePartSet(2)
		suite.Require().NoError(err)
		commitSigs := []tmtypes.CommitSig{{
			BlockIDFlag:      tmtypes.BlockIDFlagCommit,
			ValidatorAddress: tmrand.Bytes(crypto.AddressSize),
			Timestamp:        tmtime.Now(),
			Signature:        []byte("Signature"),
		}}
		commit := &tmtypes.Commit{
			Height:     h,
			Round:      0,
			BlockID:    tmtypes.BlockID{Hash: []byte(""), PartSetHeader: tmtypes.PartSetHeader{Hash: []byte(""), Total: 2}},
			Signatures: commitSigs,
		}
		suite.database.BlockStore().SaveBlock(block, partSet, commit)
	}
}
