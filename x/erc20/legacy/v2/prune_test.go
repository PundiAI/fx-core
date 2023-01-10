package v2_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v2 "github.com/functionx/fx-core/v3/x/erc20/legacy/v2"
)

var _ v2.Channelkeeper = &MigrateTestSuite{}

func (suite *MigrateTestSuite) HasPacketCommitment(_ sdk.Context, _, _ string, sequence uint64) bool {
	return sequence%2 == 0
}

func (suite *MigrateTestSuite) TestPruneExpirationIBCTransferRelation() {
	kvStore := suite.ctx.MultiStore().GetKVStore(suite.storeKey)
	v2.PruneExpirationIBCTransferRelation(suite.ctx, kvStore, suite)

	var counts2 int
	iterator := kvStore.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		_, _, sequence, ok := v2.ParseIBCTransferKeyLegacy(string(iterator.Key()[1:]))
		suite.True(ok)
		suite.Equal(sequence%2, uint64(0))
		counts2 = counts2 + 1
	}
	suite.Equal(counts2, suite.count/2+suite.count%2)
}
