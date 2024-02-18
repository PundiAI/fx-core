package keeper_test

import (
	"math"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/functionx/fx-core/v7/testutil/helpers"
)

func (suite *KeeperTestSuite) TestKeeper_PastExternalSignatureCheckpoint() {
	var checkpointCache [][]byte
	for i := 0; i < 10; i++ {
		checkpoint := crypto.Keccak256Hash(helpers.GenerateAddress().Bytes()).Bytes()
		checkpointCache = append(checkpointCache, checkpoint)
		suite.Keeper().SetPastExternalSignatureCheckpoint(suite.ctx, checkpoint)
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	}

	index := 0
	suite.Keeper().IteratePastExternalSignatureCheckpoint(suite.ctx, 0, math.MaxUint64, func(checkpoint []byte) bool {
		suite.Equal(checkpointCache[index], checkpoint)
		index = index + 1
		return false
	})
}
