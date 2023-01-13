package v3

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschainkeeper "github.com/functionx/fx-core/v3/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
)

func Test_GetMetadata(t *testing.T) {
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/-]{1,127}`
	})

	for _, metadata := range append(GetMetadata(fxtypes.MainnetChainId), GetMetadata(fxtypes.TestnetChainId)...) {
		assert.NoError(t, metadata.Validate())
		assert.NoError(t, fxtypes.ValidateMetadata(metadata))
	}
}

func Test_updateBSCOracles(t *testing.T) {
	fxtypes.SetConfig(false)
	storeKey := sdk.NewKVStoreKey(t.Name())

	ms := rootmulti.NewStore(dbm.NewMemDB(), log.NewNopLogger())
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, nil)
	assert.NoError(t, ms.LoadLatestVersion())

	ctx := sdk.Context{}.
		WithChainID(fxtypes.TestnetChainId).
		WithLogger(log.NewNopLogger()).
		WithMultiStore(ms).
		WithGasMeter(sdk.NewInfiniteGasMeter())
	protoCodec := codec.NewProtoCodec(types.NewInterfaceRegistry())
	subspace := paramtypes.NewSubspace(protoCodec, nil, nil, nil, t.Name())
	keeper := crosschainkeeper.NewKeeper(protoCodec, t.Name(), storeKey,
		subspace, nil, nil, nil, nil, nil, nil)
	updateOracles := getBSCOracleAddrs(ctx.ChainID())
	keeper.SetProposalOracle(ctx, &crosschaintypes.ProposalOracle{
		Oracles: updateOracles,
	})
	for _, oracle := range updateOracles {
		keeper.SetOracle(ctx, crosschaintypes.Oracle{
			OracleAddress: oracle,
			Online:        true,
		})
	}

	updateBSCOracles(ctx, keeper)

	proposalOracle, found := keeper.GetProposalOracle(ctx)
	assert.True(t, found)
	assert.Equal(t, 5, len(proposalOracle.Oracles))
}
