package server_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmtypes "github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/server"
)

type DatabaseTestSuite struct {
	suite.Suite
	config   *cfg.Config
	database *server.Database
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (suite *DatabaseTestSuite) SetupTest() {
	cdc := app.MakeEncodingConfig()
	newCfg := cfg.ResetTestRootWithChainID("blockchain_database_test", "fxcore")

	database, err := server.NewDatabase(newCfg, cdc.Codec)
	require.NoError(suite.T(), err)

	_, err = database.StateStore().LoadFromDBOrGenesisFile(newCfg.GenesisFile())
	require.NoError(suite.T(), err)
	suite.config = newCfg
	suite.database = database
}

func (suite *DatabaseTestSuite) TearDownSuite() {
	defer os.RemoveAll(suite.config.RootDir)
	suite.database.Close()
}

func (suite *DatabaseTestSuite) TestGetChainID() {
	_, err := suite.database.GetChainId()
	require.Error(suite.T(), err, "not found chain id")

	suite.newBlock(1)

	chainID, err := suite.database.GetChainId()
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), "fxcore", chainID)
}

func (suite *DatabaseTestSuite) TestGetBlockHeight() {
	height, err := suite.database.GetBlockHeight()
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), int64(0), height)

	suite.newBlock(5)

	height, err = suite.database.GetBlockHeight()
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), int64(5), height)
}

func (suite *DatabaseTestSuite) TestGetSyncing() {
	syncing, err := suite.database.GetSyncing()
	require.NoError(suite.T(), err)
	require.True(suite.T(), syncing)
}

func (suite *DatabaseTestSuite) TestGetNodeInfo() {
	nodeInfo, err := suite.database.GetNodeInfo()
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), nodeInfo)
}

func (suite *DatabaseTestSuite) TestCurrentPlan() {
	plan, err := suite.database.CurrentPlan()
	require.NoError(suite.T(), err)
	require.Nil(suite.T(), plan)
}

// test database GetConsensusValidators
func (suite *DatabaseTestSuite) TestGetConsensusValidators() {
	validators, err := suite.database.GetConsensusValidators()
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), validators)
	require.Equal(suite.T(), 2, len(validators))
}

// test database GetLatestHeight
func (suite *DatabaseTestSuite) newBlock(height int64) {
	state, err := suite.database.StateStore().LoadFromDBOrGenesisFile(suite.config.GenesisFile())
	require.NoError(suite.T(), err)

	for h := int64(1); h <= height; h++ {
		block, _ := state.MakeBlock(h, nil, new(tmtypes.Commit), nil, state.Validators.GetProposer().Address)
		partSet := block.MakePartSet(2)
		commitSigs := []tmtypes.CommitSig{{
			BlockIDFlag:      tmtypes.BlockIDFlagCommit,
			ValidatorAddress: tmrand.Bytes(crypto.AddressSize),
			Timestamp:        tmtime.Now(),
			Signature:        []byte("Signature"),
		}}
		commit := tmtypes.NewCommit(h, 0, tmtypes.BlockID{Hash: []byte(""), PartSetHeader: tmtypes.PartSetHeader{Hash: []byte(""), Total: 2}}, commitSigs)
		suite.database.BlockStore().SaveBlock(block, partSet, commit)
	}
}
