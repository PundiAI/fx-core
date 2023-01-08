// nolint:staticcheck
package store_test

import (
	"fmt"
	"path/filepath"
	"testing"

	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	bsctypes "github.com/functionx/fx-core/v3/x/bsc/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	erc20types "github.com/functionx/fx-core/v3/x/erc20/types"
	gravitytypes "github.com/functionx/fx-core/v3/x/gravity/types"
	migratetypes "github.com/functionx/fx-core/v3/x/migrate/types"
	polygontypes "github.com/functionx/fx-core/v3/x/polygon/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func TestLocalStoreInV2(t *testing.T) {
	if !helpers.IsLocalTest() {
		t.Skip("skipping local test", t.Name())
	}
	tests := []struct {
		name     string
		testCase func(sdk.Context, *app.App)
	}{
		{
			name: "Iterator gravity module store",
			testCase: func(ctx sdk.Context, myApp *app.App) {
				kvStore := ctx.MultiStore().GetKVStore(myApp.GetKey(gravitytypes.ModuleName))
				checkStoreKey(t, map[byte][][2]int{
					gravitytypes.EthAddressByValidatorKey[0]:              {{20, 0}},
					gravitytypes.ValidatorByEthAddressKey[0]:              {{20, 0}},
					gravitytypes.ValidatorAddressByOrchestratorAddress[0]: {{20, 0}},
					gravitytypes.LastEventBlockHeightByValidatorKey[0]:    {{20, 0}},
					gravitytypes.LastEventNonceByValidatorKey[0]:          {{20, 0}},
					gravitytypes.LastObservedEventNonceKey[0]:             {{1, 0}},
					gravitytypes.SequenceKeyPrefix[0]:                     {{2, 0}},
					gravitytypes.DenomToERC20Key[0]:                       {{1, 0}},
					gravitytypes.ERC20ToDenomKey[0]:                       {{1, 0}},
					gravitytypes.LastSlashedValsetNonce[0]:                {{1, 0}},
					gravitytypes.LatestValsetNonce[0]:                     {{1, 0}},
					gravitytypes.LastSlashedBatchBlock[0]:                 {{0, 0}},
					gravitytypes.LastUnBondingBlockHeight[0]:              {{0, 0}},
					gravitytypes.LastObservedEthereumBlockHeightKey[0]:    {{1, 0}},
					gravitytypes.LastObservedValsetKey[0]:                 {{1, 0}},
					gravitytypes.IbcSequenceHeightKey[0]:                  {{-1, -1}},
					gravitytypes.ValsetRequestKey[0]:                      {{1, 0}},
					gravitytypes.OracleAttestationKey[0]:                  {{1, 0}},
					gravitytypes.BatchConfirmKey[0]:                       {{47180, 0}},
					gravitytypes.ValsetConfirmKey[0]:                      {{1460, 0}},
				}, kvStore)
			},
		},
		{
			name: "Iterator cross-chain module store",
			testCase: func(ctx sdk.Context, myApp *app.App) {
				bscKvStore := ctx.MultiStore().GetKVStore(myApp.GetKey(bsctypes.ModuleName))
				plyKvStore := ctx.MultiStore().GetKVStore(myApp.GetKey(polygontypes.ModuleName))
				tronKvStore := ctx.MultiStore().GetKVStore(myApp.GetKey(trontypes.ModuleName))
				kvStores := []store.KVStore{bscKvStore, plyKvStore, tronKvStore}
				checkStoreKey(t, map[byte][][2]int{
					crosschaintypes.OracleKey[0]:                          {{2, 0}, {10, 0}, {10, 0}},
					crosschaintypes.OracleAddressByExternalKey[0]:         {{2, 0}, {10, 0}, {10, 0}},
					crosschaintypes.OracleAddressByBridgerKey[0]:          {{2, 0}, {10, 0}, {10, 0}},
					crosschaintypes.OracleSetRequestKey[0]:                {{1, 0}, {1, 0}, {1, 0}},
					crosschaintypes.OracleSetConfirmKey[0]:                {{3, 0}, {55, 0}, {55, 0}},
					crosschaintypes.OutgoingTxPoolKey[0]:                  {{0, 0}, {0, 0}, {0, 0}},
					crosschaintypes.OutgoingTxBatchKey[0]:                 {{0, 0}, {0, 0}, {0, 0}},
					crosschaintypes.OutgoingTxBatchBlockKey[0]:            {{0, 0}, {0, 0}, {0, 0}},
					crosschaintypes.LastEventNonceByOracleKey[0]:          {{2, 0}, {10, 0}, {10, 0}},
					crosschaintypes.LastObservedEventNonceKey[0]:          {{1, 0}, {1, 0}, {1, 0}},
					crosschaintypes.SequenceKeyPrefix[0]:                  {{2, 0}, {2, 0}, {2, 0}},
					crosschaintypes.DenomToTokenKey[0]:                    {{1, 0}, {1, 0}, {1, 0}},
					crosschaintypes.TokenToDenomKey[0]:                    {{1, 0}, {1, 0}, {1, 0}},
					crosschaintypes.LastSlashedOracleSetNonce[0]:          {{1, 0}, {1, 0}, {1, 0}},
					crosschaintypes.LatestOracleSetNonce[0]:               {{1, 0}, {1, 0}, {1, 0}},
					crosschaintypes.LastSlashedBatchBlock[0]:              {{0, 0}, {0, 0}, {1, 0}},
					crosschaintypes.LastObservedBlockHeightKey[0]:         {{1, 0}, {1, 0}, {1, 0}},
					crosschaintypes.LastObservedOracleSetKey[0]:           {{1, 0}, {1, 0}, {1, 0}},
					crosschaintypes.LastEventBlockHeightByOracleKey[0]:    {{2, 0}, {10, 0}, {10, 0}},
					crosschaintypes.LastOracleSlashBlockHeight[0]:         {{0, 0}, {0, 0}, {0, 0}},
					crosschaintypes.ProposalOracleKey[0]:                  {{1, 0}, {1, 0}, {1, 0}},
					crosschaintypes.LastTotalPowerKey[0]:                  {{1, 0}, {1, 0}, {1, 0}},
					crosschaintypes.OracleAttestationKey[0]:               {{101, 0}, {103, 0}, {103, 0}},
					crosschaintypes.PastExternalSignatureCheckpointKey[0]: {{505, 0}, {80, 0}, {55, 0}},
					crosschaintypes.BatchConfirmKey[0]:                    {{1006, 0}, {700, 0}, {450, 0}},
					crosschaintypes.LastProposalBlockHeight[0]:            {{-1, -1}, {-1, -1}, {-1, -1}},
				}, kvStores...)
			},
		},
		{
			name: "Iterator erc20 module store",
			testCase: func(ctx sdk.Context, myApp *app.App) {
				kvStore := ctx.MultiStore().GetKVStore(myApp.GetKey(erc20types.ModuleName))
				expected := map[byte][][2]int{
					erc20types.KeyPrefixTokenPair[0]:        {{12, 0}},
					erc20types.KeyPrefixTokenPairByERC20[0]: {{12, 0}},
					erc20types.KeyPrefixTokenPairByDenom[0]: {{12, 0}},
					erc20types.KeyPrefixIBCTransfer[0]:      {{182, 0}},
					erc20types.KeyPrefixAliasDenom[0]:       {{11, 0}},
				}
				checkStoreKey(t, expected, kvStore)
			},
		},
		{
			name: "Iterator migrate module store",
			testCase: func(ctx sdk.Context, myApp *app.App) {
				kvStore := ctx.MultiStore().GetKVStore(myApp.GetKey(migratetypes.ModuleName))
				expected := map[byte][][2]int{
					migratetypes.KeyPrefixMigratedRecord[0]:        {{1242, 0}},
					migratetypes.KeyPrefixMigratedDirectionFrom[0]: {{621, 0}},
					migratetypes.KeyPrefixMigratedDirectionTo[0]:   {{621, 0}},
				}
				checkStoreKey(t, expected, kvStore)
			},
		},
	}
	db, err := sdk.NewLevelDB("application", filepath.Join(fxtypes.GetDefaultNodeHome(), "data"))
	require.NoError(t, err)

	appEncodingCfg := app.MakeEncodingConfig()
	myApp := app.New(log.NewNopLogger(), db,
		nil, true, map[int64]bool{}, fxtypes.GetDefaultNodeHome(), 0,
		appEncodingCfg, app.EmptyAppOptions{},
	)
	ctx := myApp.NewUncachedContext(false, tmproto.Header{Height: myApp.LastBlockHeight()})
	require.Equal(t, ctx.BlockHeight(), int64(7654832))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.testCase(ctx, myApp)
		})
	}
}

func checkStoreKey(t *testing.T, expected map[byte][][2]int, stores ...store.KVStore) {
	for i := 0; i < len(stores); i++ {
		iterator := stores[i].Iterator(nil, nil)
		for ; iterator.Valid(); iterator.Next() {
			x, ok := expected[iterator.Key()[0]]
			assert.True(t, ok, fmt.Sprintf("%x", iterator.Key()[0]), iterator.Value())
			if ok {
				if x[i][0] == -1 && x[i][1] == -1 {
					// ignore
					continue
				}
				// set result
				expected[iterator.Key()[0]][i] = [2]int{x[i][0], x[i][1] + 1}
			}
		}
		iterator.Close()
		for k, x := range expected {
			assert.Equal(t, x[i][0], x[i][1], fmt.Sprintf("%x", k))
		}
	}
}
