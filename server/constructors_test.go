package server

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func Test_openDB(t *testing.T) {
	t.Parallel()
	_, err := openDB(t.TempDir(), dbm.GoLevelDBBackend)
	require.NoError(t, err)
}

func Test_openTraceWriter(t *testing.T) {
	t.Parallel()

	fname := filepath.Join(t.TempDir(), "logfile")
	w, err := openTraceWriter(fname)
	require.NoError(t, err)
	require.NotNil(t, w)

	// test no-op
	w, err = openTraceWriter("")
	require.NoError(t, err)
	require.Nil(t, w)
}

func Test_genesis_hash(t *testing.T) {
	genesisFile := filepath.Join("../public", "mainnet", "genesis.json")
	genesisDoc, err := tmtypes.GenesisDocFromFile(genesisFile)
	assert.NoError(t, err)
	genesisBytes, err := tmjson.Marshal(genesisDoc)
	assert.NoError(t, err)
	assert.Equal(t, sha256Hex(genesisBytes), mainnetGenesisHash)
}
