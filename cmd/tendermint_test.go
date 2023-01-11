package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"
)

func Test_genesis_hash(t *testing.T) {
	genesisFile := filepath.Join("../public", "mainnet", "genesis.json")
	genesisDoc, err := tmtypes.GenesisDocFromFile(genesisFile)
	assert.NoError(t, err)
	data, err := tmjson.Marshal(genesisDoc)
	assert.NoError(t, err)
	hash := md5.Sum(data)
	assert.Equal(t, hex.EncodeToString(hash[:]), mainnetGenesisHash)
}
