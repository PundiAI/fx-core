package server_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmtypes "github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v4/app"
	"github.com/functionx/fx-core/v4/server"
)

func TestDatabase(t *testing.T) {
	cdc := app.MakeEncodingConfig()
	config := cfg.ResetTestRoot("blockchain_database_test")
	defer os.RemoveAll(config.RootDir)

	database, err := server.NewDatabase(config.RootDir, string(dbm.GoLevelDBBackend), cdc.Codec)
	require.NoError(t, err)

	defer database.Close()
	_, err = database.GetChainId()
	require.Errorf(t, err, errors.New("not found chain id").Error())

	state, err := database.StateStore().LoadFromDBOrGenesisFile(config.GenesisFile())
	require.NoError(t, err)

	for h := int64(1); h <= 5; h++ {
		block, _ := state.MakeBlock(h, nil, new(tmtypes.Commit), nil, state.Validators.GetProposer().Address)
		partSet := block.MakePartSet(2)
		commitSigs := []tmtypes.CommitSig{{
			BlockIDFlag:      tmtypes.BlockIDFlagCommit,
			ValidatorAddress: tmrand.Bytes(crypto.AddressSize),
			Timestamp:        tmtime.Now(),
			Signature:        []byte("Signature"),
		}}
		commit := tmtypes.NewCommit(h, 0, tmtypes.BlockID{Hash: []byte(""), PartSetHeader: tmtypes.PartSetHeader{Hash: []byte(""), Total: 2}}, commitSigs)
		database.BlockStore().SaveBlock(block, partSet, commit)
	}

	chainid, err := database.GetChainId()
	require.NoError(t, err)
	require.Equal(t, chainid, "cometbft_test")

	height, err := database.GetBlockHeight()
	require.NoError(t, err)
	require.Equal(t, height, int64(5))

	_, err = database.GetNodeInfo()
	require.Errorf(t, err, "genesis doc not found")

	genesis, err := os.ReadFile(config.GenesisFile())
	require.NoError(t, err)

	err = database.SetGenesis(genesis)
	require.NoError(t, err)

	_, err = database.GetNodeInfo()
	require.NoError(t, err)
	_, err = database.CurrentPlan()
	require.NoError(t, err)
	_, err = database.GetValidators()
	require.NoError(t, err)
}
