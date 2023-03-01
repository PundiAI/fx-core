package v3_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	dbm "github.com/tendermint/tm-db"

	erc20keeper "github.com/functionx/fx-core/v3/x/erc20/keeper"
	v3 "github.com/functionx/fx-core/v3/x/erc20/migrations/v3"
	"github.com/functionx/fx-core/v3/x/erc20/types"
)

type MigrateTestSuite struct {
	suite.Suite
	ctx      sdk.Context
	storeKey *storetypes.KVStoreKey
	count    int
}

func TestMigrate(t *testing.T) {
	suite.Run(t, new(MigrateTestSuite))
}

func (suite *MigrateTestSuite) SetupTest() {
	suite.storeKey = sdk.NewKVStoreKey(suite.T().Name())
	ms := rootmulti.NewStore(dbm.NewMemDB(), log.NewNopLogger())
	ms.MountStoreWithDB(suite.storeKey, storetypes.StoreTypeIAVL, nil)
	suite.NoError(ms.LoadLatestVersion())

	suite.ctx = sdk.Context{}.
		WithLogger(log.NewNopLogger()).
		WithMultiStore(ms).
		WithGasMeter(sdk.NewInfiniteGasMeter())

	kvStore := ms.GetKVStore(suite.storeKey)
	suite.count = tmrand.Intn(100)
	for i := 0; i < suite.count; i++ {
		key := append(types.KeyPrefixIBCTransfer, []byte(fmt.Sprintf("transfer/channel-%d/%d", i%10, i))...)
		kvStore.Set(key, []byte{})
	}
}

func (suite *MigrateTestSuite) TestMigrateIBCTransferRelation() {
	protoCodec := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

	keeper := erc20keeper.NewKeeper(suite.storeKey, protoCodec, suite, nil, nil, nil, "")

	kvStore := suite.ctx.MultiStore().GetKVStore(suite.storeKey)
	v3.MigrateIBCTransferRelation(suite.ctx, kvStore, keeper)

	var counts2 int
	iterator := kvStore.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		split := strings.Split(string(iterator.Key()[1:]), "/")
		suite.Equal(2, len(split))
		_, err := strconv.ParseUint(split[1], 10, 64)
		suite.NoError(err)
		counts2 = counts2 + 1
	}
	suite.Equal(counts2, suite.count)
}

var _ types.AccountKeeper = &MigrateTestSuite{}

func (suite *MigrateTestSuite) GetModuleAccount(sdk.Context, string) authtypes.ModuleAccountI {
	return &authtypes.ModuleAccount{}
}

func (suite *MigrateTestSuite) GetModuleAddress(string) sdk.AccAddress {
	return sdk.AccAddress{}
}

func (suite *MigrateTestSuite) GetSequence(sdk.Context, sdk.AccAddress) (uint64, error) {
	return 0, nil
}
