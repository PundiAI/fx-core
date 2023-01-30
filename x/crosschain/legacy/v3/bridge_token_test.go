// nolint:staticcheck
package v3_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/testutil/helpers"
	v3 "github.com/functionx/fx-core/v3/x/crosschain/legacy/v3"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
)

func TestMigrateBridgeToken(t *testing.T) {
	moduleName := ethtypes.ModuleName
	storeKey := sdk.NewKVStoreKey(moduleName)
	ms := rootmulti.NewStore(dbm.NewMemDB(), log.NewNopLogger())
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, nil)
	assert.NoError(t, ms.LoadLatestVersion())

	store := ms.GetKVStore(storeKey)
	encodingConfig := app.MakeEncodingConfig()
	cdc := encodingConfig.Codec

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(100)

	var bridgeTokens []types.BridgeToken
	for i := 0; i < index; i++ {
		bridgeToken := types.BridgeToken{
			Token: helpers.GenerateAddress().Hex(),
		}
		bridgeToken.Denom = fmt.Sprintf("%s%s", moduleName, bridgeToken.Token)
		if i%5 == 0 {
			bridgeToken.ChannelIbc = "transfer/channel-0"
			bridgeToken.Denom = ibctransfertypes.DenomTrace{
				Path:      bridgeToken.ChannelIbc,
				BaseDenom: bridgeToken.Denom,
			}.IBCDenom()
		}
		store.Set(types.GetTokenToDenomKey(bridgeToken.Denom),
			cdc.MustMarshal(&types.BridgeToken{
				Token:      bridgeToken.Token,
				ChannelIbc: bridgeToken.ChannelIbc,
			}),
		)
		store.Set(types.GetDenomToTokenKey(bridgeToken.Token),
			cdc.MustMarshal(&types.BridgeToken{
				Denom:      bridgeToken.Denom,
				ChannelIbc: bridgeToken.ChannelIbc,
			}),
		)
		bridgeTokens = append(bridgeTokens, bridgeToken)
	}

	v3.MigrateBridgeToken(cdc, store)

	for _, bridgeToken := range bridgeTokens {
		assert.Equal(t, store.Get(types.GetTokenToDenomKey(bridgeToken.Denom)), []byte(bridgeToken.Token))
		assert.Equal(t, store.Get(types.GetDenomToTokenKey(bridgeToken.Token)), []byte(bridgeToken.Denom))
	}
}
